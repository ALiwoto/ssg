package tests

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

func TestSafeEMapCheckerLoopRapidToggleUniqueness(t *testing.T) {
	const (
		workerCount     = 16
		togglesPerGo    = 250
		checkerInterval = 100 * time.Microsecond
	)

	m := ssg.NewSafeEMap[int, valuesContainer]()
	m.SetInterval(checkerInterval)
	m.EnableChecking()
	t.Cleanup(m.DisableChecking)

	waitForStableCheckerLoopCount(t, 1, 20*time.Millisecond, 2*time.Second)

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for worker := range workerCount {
		go func(worker int) {
			defer wg.Done()
			<-start

			for toggle := range togglesPerGo {
				m.DisableChecking()
				runtime.Gosched()
				m.EnableChecking()

				if (worker+toggle)%8 == 0 {
					runtime.Gosched()
				}
			}
		}(worker)
	}

	close(start)
	wg.Wait()

	// Establish the requested final state after all concurrent togglers finish.
	m.EnableChecking()
	if !m.IsChecking() {
		t.Fatal("checking is disabled after the final EnableChecking call")
	}
	waitForStableCheckerLoopCount(t, 1, 50*time.Millisecond, 3*time.Second)

	m.DisableChecking()
	if m.IsChecking() {
		t.Fatal("checking is enabled after DisableChecking")
	}
	waitForStableCheckerLoopCount(t, 0, 20*time.Millisecond, 3*time.Second)
}

func TestSafeEMapCheckerLoopProcessesEntriesOnceAfterToggleStress(t *testing.T) {
	const (
		entryCount       = 1_000
		togglerCount     = 12
		togglesPerGo     = 200
		checkerInterval  = 100 * time.Microsecond
		expirationPeriod = 500 * time.Microsecond
	)

	m := ssg.NewSafeEMap[int, valuesContainer]()
	m.SetInterval(checkerInterval)
	m.SetExpiration(expirationPeriod)

	var expiredCount atomic.Int64
	m.SetOnExpired(func(_ int, _ valuesContainer) {
		expiredCount.Add(1)
	})

	m.EnableChecking()
	t.Cleanup(m.DisableChecking)

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(togglerCount + 1)

	for range togglerCount {
		go func() {
			defer wg.Done()
			<-start

			for range togglesPerGo {
				m.DisableChecking()
				runtime.Gosched()
				m.EnableChecking()
			}
		}()
	}

	go func() {
		defer wg.Done()
		<-start

		for key := range entryCount {
			m.Set(key, valuesContainer{Value1: key, Value2: "expiring"})
			if key%16 == 0 {
				runtime.Gosched()
			}
		}
	}()

	close(start)
	wg.Wait()
	m.EnableChecking()

	waitForCondition(t, 5*time.Second, func() bool {
		return m.Length() == 0 && expiredCount.Load() == entryCount
	}, func() string {
		return "length=" + strconv.Itoa(m.Length()) +
			", expired callbacks=" + strconv.FormatInt(expiredCount.Load(), 10)
	})

	// Give any incorrectly duplicated callback goroutines time to report.
	time.Sleep(20 * time.Millisecond)
	if got := expiredCount.Load(); got != entryCount {
		t.Fatalf("expired callback count is %d, want exactly %d", got, entryCount)
	}
	waitForStableCheckerLoopCount(t, 1, 20*time.Millisecond, 3*time.Second)

	m.DisableChecking()
	waitForStableCheckerLoopCount(t, 0, 20*time.Millisecond, 3*time.Second)
}

func waitForStableCheckerLoopCount(
	t *testing.T,
	want int,
	stableFor time.Duration,
	timeout time.Duration,
) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	stableSince := time.Time{}
	lastCount := -1
	for time.Now().Before(deadline) {
		lastCount = safeEMapCheckerLoopCount()
		if lastCount == want {
			if stableSince.IsZero() {
				stableSince = time.Now()
			} else if time.Since(stableSince) >= stableFor {
				return
			}
		} else {
			stableSince = time.Time{}
		}

		time.Sleep(time.Millisecond)
	}

	t.Fatalf("checker loop count did not remain at %d: last count was %d", want, lastCount)
}

func safeEMapCheckerLoopCount() int {
	bufferSize := 64 << 10
	for {
		stack := make([]byte, bufferSize)
		length := runtime.Stack(stack, true)
		if length == len(stack) {
			bufferSize *= 2
			continue
		}

		count := 0
		for _, line := range bytes.Split(stack[:length], []byte{'\n'}) {
			if bytes.Contains(line, []byte("mapUtils.(*SafeEMap[")) &&
				bytes.Contains(line, []byte(").checkLoop(")) {
				count++
			}
		}
		return count
	}
}

func waitForCondition(
	t *testing.T,
	timeout time.Duration,
	condition func() bool,
	state func() string,
) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatalf("condition was not met before timeout: %s", state())
}
