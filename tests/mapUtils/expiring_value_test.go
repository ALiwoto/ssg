package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

func TestExpiringValueGetValueResetOption(t *testing.T) {
	value := mapUtils.NewEValue("value")
	originalTime := time.Unix(1, 0)
	value.SetTime(originalTime)

	if got := value.GetValue(false); got != "value" {
		t.Fatalf("GetValue(false) returned %q, want %q", got, "value")
	}
	if got := value.GetTime(); !got.Equal(originalTime) {
		t.Fatalf("GetValue(false) changed timestamp to %v, want %v", got, originalTime)
	}

	value.GetValue(true)
	if got := value.GetTime(); !got.After(originalTime) {
		t.Fatalf("GetValue(true) did not reset timestamp: got %v", got)
	}
}

func TestSafeEMapConcurrentGetResetsTimestampSafely(t *testing.T) {
	const (
		workerCount = 32
		readCount   = 500
	)

	m := ssg.NewSafeEMap[string, valuesContainer]()
	m.Set("key", valuesContainer{Value1: 42, Value2: "concurrent"})

	start := make(chan struct{})
	results := make(chan bool, workerCount)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for range workerCount {
		go func() {
			defer wg.Done()
			<-start

			for range readCount {
				value := m.Get("key")
				if value == nil || value.Value1 != 42 || value.Value2 != "concurrent" {
					results <- false
					return
				}
			}

			results <- true
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	for result := range results {
		if !result {
			t.Fatal("concurrent Get returned an unexpected value")
		}
	}
}
