package tests

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/ALiwoto/ssg/ssg"
)

func TestAdvancedMapGetOrCreateAndClearBookkeeping(t *testing.T) {
	m := ssg.NewAdvancedMap[string, valuesContainer]()

	created := m.GetOrCreate("created", func() (*valuesContainer, bool) {
		return &valuesContainer{Value1: 42, Value2: "created"}, true
	})
	if created == nil || created.Value1 != 42 || created.Value2 != "created" {
		t.Fatalf("GetOrCreate returned %v, want 42", created)
	}

	key, ok := m.GetRandomKey()
	if !ok || key != "created" {
		t.Fatalf("GetRandomKey returned (%q, %v), want (%q, true)", key, ok, "created")
	}
	if value := m.GetRandom(); value == nil || value.Value1 != 42 || value.Value2 != "created" {
		t.Fatalf("GetRandom returned %v, want 42", value)
	}

	createCalled := false
	existing := m.GetOrCreate("created", func() (*valuesContainer, bool) {
		createCalled = true
		return &valuesContainer{Value1: 99, Value2: "replacement"}, true
	})
	if createCalled || existing != created {
		t.Fatal("GetOrCreate replaced an existing value")
	}
	if m.Length() != 1 {
		t.Fatalf("Length after repeated GetOrCreate is %d, want 1", m.Length())
	}

	m.Delete("created")
	if m.Exists("created") || m.Length() != 0 {
		t.Fatalf("Delete did not remove GetOrCreate entry; length is %d", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after Delete", key)
	}

	m.Set("old-1", valuesContainer{Value1: 1, Value2: "old one"})
	m.Set("old-2", valuesContainer{Value1: 2, Value2: "old two"})
	m.Clear()
	if m.Length() != 0 {
		t.Fatalf("Length after Clear is %d, want 0", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after Clear", key)
	}
	if value := m.GetRandom(); value != nil {
		t.Fatalf("GetRandom returned stale value %v after Clear", *value)
	}

	m.Set("new", valuesContainer{Value1: 3, Value2: "new"})
	for range 100 {
		key, ok := m.GetRandomKey()
		if !ok || key != "new" {
			t.Fatalf("GetRandomKey returned (%q, %v), want (%q, true)", key, ok, "new")
		}
		if value := m.GetRandom(); value == nil || value.Value1 != 3 || value.Value2 != "new" {
			t.Fatalf("GetRandom returned %v after Clear, want 3", value)
		}
	}
}

func TestAdvancedMapAggressiveBookkeepingModel(t *testing.T) {
	const operationCount = 20_000

	rng := rand.New(rand.NewSource(1))
	m := ssg.NewAdvancedMap[int, valuesContainer]()
	m.SetDefault(valuesContainer{Value1: -1, Value2: "default"})
	model := make(map[int]valuesContainer)

	for step := range operationCount {
		key := rng.Intn(128)
		value := valuesContainer{Value1: step + 1, Value2: "model value"}

		switch rng.Intn(12) {
		case 0, 1, 2:
			m.Set(key, value)
			model[key] = value
		case 3:
			m.Add(key, &value)
			model[key] = value
		case 4, 5:
			expected, exists := model[key]
			createCalled := false
			got := m.GetOrCreate(key, func() (*valuesContainer, bool) {
				createCalled = true
				return &value, true
			})
			if !exists {
				expected = value
				model[key] = value
			}
			if createCalled == exists {
				t.Fatalf("step %d: create callback state is %v, key existed: %v", step, createCalled, exists)
			}
			if got == nil || *got != expected {
				t.Fatalf("step %d: GetOrCreate returned %v, want %v", step, got, expected)
			}
		case 6, 7:
			m.Delete(key)
			delete(model, key)
		case 8:
			m.Clear()
			clear(model)
		case 9:
			expected, exists := model[key]
			if got := m.Exists(key); got != exists {
				t.Fatalf("step %d: Exists(%d) is %v, want %v", step, key, got, exists)
			}
			if exists {
				if got := m.Get(key); got == nil || *got != expected {
					t.Fatalf("step %d: Get(%d) returned %v, want %v", step, key, got, expected)
				}
			}
		case 10, 11:
			// Random-access behavior is checked below after every operation.
		}

		assertAdvancedMapMatchesModel(t, step, m, model)
	}
}

func assertAdvancedMapMatchesModel(
	t *testing.T,
	step int,
	m *ssg.AdvancedMap[int, valuesContainer],
	model map[int]valuesContainer,
) {
	t.Helper()

	if got := m.Length(); got != len(model) {
		t.Fatalf("step %d: Length is %d, want %d", step, got, len(model))
	}

	key, ok := m.GetRandomKey()
	if len(model) == 0 {
		if ok {
			t.Fatalf("step %d: GetRandomKey returned stale key %d", step, key)
		}
		if value := m.GetRandom(); value != nil {
			t.Fatalf("step %d: GetRandom returned stale value %v", step, *value)
		}
		if value := m.GetRandomValue(); value.Value1 != -1 || value.Value2 != "default" {
			t.Fatalf("step %d: GetRandomValue returned %v for empty map, want -1", step, value)
		}
		return
	}

	expected, exists := model[key]
	if !ok || !exists {
		t.Fatalf("step %d: GetRandomKey returned (%d, %v), which is absent from the model", step, key, ok)
	}
	if got := m.GetValue(key); got != expected {
		t.Fatalf("step %d: random key %d resolves to %v, want %v", step, key, got, expected)
	}

	if value := m.GetRandom(); value == nil || !advancedModelContainsValue(model, *value) {
		t.Fatalf("step %d: GetRandom returned value outside the model: %v", step, value)
	}
	if value := m.GetRandomValue(); !advancedModelContainsValue(model, value) {
		t.Fatalf("step %d: GetRandomValue returned value outside the model: %v", step, value)
	}
}

func advancedModelContainsValue(model map[int]valuesContainer, value valuesContainer) bool {
	for _, current := range model {
		if current == value {
			return true
		}
	}

	return false
}

func TestAdvancedMapConcurrentBookkeepingStress(t *testing.T) {
	const (
		workerCount     = 32
		operationsPerGo = 2_000
	)

	m := ssg.NewAdvancedMap[int, valuesContainer]()
	m.SetDefault(valuesContainer{Value1: -1, Value2: "default"})
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for worker := range workerCount {
		go func(worker int) {
			defer wg.Done()
			<-start

			rng := rand.New(rand.NewSource(int64(worker + 1)))
			for operation := range operationsPerGo {
				key := rng.Intn(256)
				value := valuesContainer{
					Value1: worker*operationsPerGo + operation + 1,
					Value2: "worker value",
				}

				switch rng.Intn(8) {
				case 0, 1:
					m.Set(key, value)
				case 2:
					m.Add(key, &value)
				case 3:
					m.GetOrCreate(key, func() (*valuesContainer, bool) { return &value, true })
				case 4:
					m.Delete(key)
				case 5:
					m.GetRandomKey()
				case 6:
					m.GetRandom()
				case 7:
					m.Clear()
				}
			}
		}(worker)
	}

	close(start)
	wg.Wait()

	m.Clear()
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %d after final Clear", key)
	}

	m.Set(999, valuesContainer{Value1: 123, Value2: "final"})
	for range 1_000 {
		key, ok := m.GetRandomKey()
		if !ok || key != 999 {
			t.Fatalf("GetRandomKey returned (%d, %v), want (999, true)", key, ok)
		}
		if value := m.GetRandom(); value == nil || value.Value1 != 123 || value.Value2 != "final" {
			t.Fatalf("GetRandom returned %v, want 123", value)
		}
	}

	m.Delete(999)
	if !m.IsEmpty() {
		t.Fatalf("map is not empty after deleting final entry; length is %d", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %d after deleting final entry", key)
	}
}
