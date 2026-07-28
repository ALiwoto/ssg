package tests

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

type getOptionsMap[TValue any] struct {
	add func(string, *TValue)
	get func(
		string,
		*mapUtils.GetOptions[string, TValue],
	) *TValue
}

func TestGetWithOptionsAllowsConcurrentKeys(t *testing.T) {
	tests := []struct {
		name   string
		newMap func() getOptionsMap[sync.Mutex]
	}{
		{
			name: "SafeMap",
			newMap: func() getOptionsMap[sync.Mutex] {
				m := mapUtils.NewSafeMap[string, sync.Mutex]()
				return getOptionsMap[sync.Mutex]{add: m.Add, get: m.GetWithOptions}
			},
		},
		{
			name: "AdvancedMap",
			newMap: func() getOptionsMap[sync.Mutex] {
				m := mapUtils.NewAdvancedMap[string, sync.Mutex]()
				return getOptionsMap[sync.Mutex]{add: m.Add, get: m.GetWithOptions}
			},
		},
		{
			name: "SafeEMap",
			newMap: func() getOptionsMap[sync.Mutex] {
				m := mapUtils.NewSafeEMap[string, sync.Mutex]()
				m.SetExpiration(time.Hour)
				return getOptionsMap[sync.Mutex]{add: m.Add, get: m.GetWithOptions}
			},
		},
	}

	for _, current := range tests {
		t.Run(current.name, func(t *testing.T) {
			testGetWithOptionsAllowsConcurrentKeys(t, current.newMap())
		})
	}
}

func testGetWithOptionsAllowsConcurrentKeys(
	t *testing.T,
	m getOptionsMap[sync.Mutex],
) {
	t.Helper()

	blockedValue := &sync.Mutex{}
	blockedValue.Lock()
	m.add("blocked", blockedValue)
	m.add("available", &sync.Mutex{})

	blockedStarted := make(chan struct{})
	blockedDone := make(chan struct{})
	go func() {
		value := m.get("blocked", &mapUtils.GetOptions[string, sync.Mutex]{
			DoFn: func(value *sync.Mutex) {
				close(blockedStarted)
				value.Lock()
			},
		})
		value.Unlock()
		close(blockedDone)
	}()

	select {
	case <-blockedStarted:
	case <-time.After(time.Second):
		t.Fatal("blocked DoFn did not start")
	}

	availableDone := make(chan struct{})
	go func() {
		value := m.get("available", &mapUtils.GetOptions[string, sync.Mutex]{
			DoFn: func(value *sync.Mutex) {
				value.Lock()
			},
		})
		value.Unlock()
		close(availableDone)
	}()

	select {
	case <-availableDone:
	case <-time.After(time.Second):
		blockedValue.Unlock()
		t.Fatal("a blocked key prevented an available key from running")
	}

	blockedValue.Unlock()
	select {
	case <-blockedDone:
	case <-time.After(time.Second):
		t.Fatal("blocked DoFn did not finish after the value was unlocked")
	}
}

func TestGetWithOptionsDoubleChecksConcurrentCreation(t *testing.T) {
	for _, current := range newContainerGetOptionsMaps() {
		t.Run(current.name, func(t *testing.T) {
			const workerCount = 32

			start := make(chan struct{})
			results := make(chan *valuesContainer, workerCount)
			var createCalls atomic.Int64
			var wg sync.WaitGroup
			wg.Add(workerCount)

			for range workerCount {
				go func() {
					defer wg.Done()
					<-start
					results <- current.m.get(
						"created",
						&mapUtils.GetOptions[string, valuesContainer]{
							CreateFn: func() (*valuesContainer, bool) {
								createCalls.Add(1)
								return &valuesContainer{
									Value1: 1,
									Value2: "created",
								}, true
							},
						},
					)
				}()
			}

			close(start)
			wg.Wait()
			close(results)

			var retained *valuesContainer
			for value := range results {
				if retained == nil {
					retained = value
					continue
				}
				if value != retained {
					t.Fatal("concurrent callers returned different retained values")
				}
			}
			if createCalls.Load() != 1 {
				t.Fatalf("CreateFn was called %d times, want 1", createCalls.Load())
			}
		})
	}
}

func TestGetWithOptionsStoredNilRemainsFound(t *testing.T) {
	for _, current := range newContainerGetOptionsMaps() {
		t.Run(current.name, func(t *testing.T) {
			value := current.m.get("nil", &mapUtils.GetOptions[string, valuesContainer]{
				CreateFn: func() (*valuesContainer, bool) {
					return nil, true
				},
			})
			if value != nil {
				t.Fatal("accepted nil creation returned a non-nil value")
			}

			createCalled := false
			doCalled := false
			value = current.m.get("nil", &mapUtils.GetOptions[string, valuesContainer]{
				CreateFn: func() (*valuesContainer, bool) {
					createCalled = true
					return &valuesContainer{Value1: 2, Value2: "replacement"}, true
				},
				DoFn: func(value *valuesContainer) {
					doCalled = true
					if value != nil {
						t.Error("stored nil changed before DoFn")
					}
				},
			})
			if value != nil || createCalled || !doCalled {
				t.Fatal("stored nil was treated as a missing value")
			}
		})
	}
}

func TestGetWithOptionsPanicsReleaseLocks(t *testing.T) {
	for _, current := range newContainerGetOptionsMaps() {
		t.Run(current.name+"/DoFn", func(t *testing.T) {
			current.m.add("existing", &valuesContainer{Value1: 1, Value2: "existing"})
			assertGetWithOptionsPanics(t, func() {
				current.m.get("existing", &mapUtils.GetOptions[string, valuesContainer]{
					DoFn: func(*valuesContainer) {
						panic("DoFn panic")
					},
				})
			})
			assertMapWriteCompletes(t, current.m, "after-do-panic")
		})

		t.Run(current.name+"/CreateFn", func(t *testing.T) {
			assertGetWithOptionsPanics(t, func() {
				current.m.get("missing", &mapUtils.GetOptions[string, valuesContainer]{
					CreateFn: func() (*valuesContainer, bool) {
						panic("CreateFn panic")
					},
				})
			})
			assertMapWriteCompletes(t, current.m, "after-create-panic")
		})
	}
}

func TestSafeEMapPreExpiringConditionPanicReleasesLock(t *testing.T) {
	m := mapUtils.NewSafeEMap[string, valuesContainer]()
	m.SetExpiration(-time.Nanosecond)
	m.Set("expired", valuesContainer{Value1: 1, Value2: "expired"})
	m.SetPreExpiringConditionFn(func(string, *valuesContainer) bool {
		panic("pre-expiring condition panic")
	})

	assertGetWithOptionsPanics(t, func() {
		m.GetWithOptions("expired", &mapUtils.GetOptions[string, valuesContainer]{
			CreateFn: func() (*valuesContainer, bool) {
				return &valuesContainer{Value1: 2, Value2: "replacement"}, true
			},
		})
	})

	done := make(chan struct{})
	go func() {
		m.SetPreExpiringConditionFn(nil)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("pre-expiring condition panic left the write lock held")
	}
}

type namedContainerGetOptionsMap struct {
	name string
	m    getOptionsMap[valuesContainer]
}

func newContainerGetOptionsMaps() []namedContainerGetOptionsMap {
	safeMap := mapUtils.NewSafeMap[string, valuesContainer]()
	advancedMap := mapUtils.NewAdvancedMap[string, valuesContainer]()
	expiringMap := mapUtils.NewSafeEMap[string, valuesContainer]()
	expiringMap.SetExpiration(time.Hour)

	return []namedContainerGetOptionsMap{
		{
			name: "SafeMap",
			m: getOptionsMap[valuesContainer]{
				add: safeMap.Add,
				get: safeMap.GetWithOptions,
			},
		},
		{
			name: "AdvancedMap",
			m: getOptionsMap[valuesContainer]{
				add: advancedMap.Add,
				get: advancedMap.GetWithOptions,
			},
		},
		{
			name: "SafeEMap",
			m: getOptionsMap[valuesContainer]{
				add: expiringMap.Add,
				get: expiringMap.GetWithOptions,
			},
		},
	}
}

func assertGetWithOptionsPanics(t *testing.T, fn func()) {
	t.Helper()

	didPanic := false
	func() {
		defer func() {
			didPanic = recover() != nil
		}()
		fn()
	}()
	if !didPanic {
		t.Fatal("GetWithOptions did not propagate the callback panic")
	}
}

func assertMapWriteCompletes(
	t *testing.T,
	m getOptionsMap[valuesContainer],
	key string,
) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		m.add(key, &valuesContainer{Value1: 2, Value2: "after panic"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("callback panic left a map lock held")
	}
}
