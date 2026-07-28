package tests

import (
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

func TestAdvancedMapGetWithOptions(t *testing.T) {
	m := mapUtils.NewAdvancedMap[string, valuesContainer]()
	m.Set("existing", valuesContainer{Value1: 1, Value2: "existing"})

	createCalled := false
	value := m.GetWithOptions("existing", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) {
			createCalled = true
			return nil, false
		},
		DoFn: func(value *valuesContainer) {
			value.Value1 = 2
			value.Value2 = "handled"
		},
	})
	if createCalled || value == nil ||
		value.Value1 != 2 || value.Value2 != "handled" {
		t.Fatalf("existing value was not handled correctly: %v", value)
	}

	doCalled := false
	value = m.GetWithOptions("rejected", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) { return nil, false },
		DoFn:     func(*valuesContainer) { doCalled = true },
	})
	if value != nil || doCalled || m.Exists("rejected") {
		t.Fatal("rejected creation was treated as a found value")
	}

	value = m.GetWithOptions("nil", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) { return nil, true },
		DoFn: func(value *valuesContainer) {
			if value != nil {
				t.Fatal("created nil value changed before DoFn")
			}
			doCalled = true
		},
	})
	if value != nil || !doCalled || !m.Exists("nil") {
		t.Fatal("accepted nil creation was not registered as an existing value")
	}
}

func TestSafeEMapGetWithOptionsExpiration(t *testing.T) {
	m := mapUtils.NewSafeEMap[string, valuesContainer]()
	m.SetExpiration(time.Hour)
	m.Set("value", valuesContainer{Value1: 1, Value2: "original"})

	createCalled := false
	value := m.GetWithOptions("value", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) {
			createCalled = true
			return nil, false
		},
		DoFn: func(value *valuesContainer) {
			value.Value1 = 2
			value.Value2 = "handled"
		},
	})
	if createCalled || value == nil || value.Value1 != 2 || value.Value2 != "handled" {
		t.Fatalf("non-expired value was not handled correctly: %v", value)
	}

	m.SetExpiration(-time.Nanosecond)
	doCalled := false
	value = m.GetWithOptions("value", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) { return nil, false },
		DoFn:     func(*valuesContainer) { doCalled = true },
	})
	if value != nil || doCalled || m.Length() != 1 {
		t.Fatal("failed creation replaced or exposed an expired value")
	}

	value = m.GetWithOptions("value", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) {
			return &valuesContainer{Value1: 3, Value2: "replacement"}, true
		},
		DoFn: func(value *valuesContainer) {
			value.Value1 = 4
			value.Value2 = "replaced and handled"
		},
	})
	if value == nil || value.Value1 != 4 || value.Value2 != "replaced and handled" || m.Length() != 1 {
		t.Fatalf("expired value was not replaced correctly: %v", value)
	}

	value = m.GetWithOptions("value", &mapUtils.GetOptions[string, valuesContainer]{
		DoFn: func(value *valuesContainer) {
			value.Value1 = 5
			value.Value2 = "read without creator"
		},
	})
	if value == nil || value.Value1 != 5 || value.Value2 != "read without creator" {
		t.Fatalf("expired value was not readable without a creator: %v", value)
	}
	key, ok := m.GetRandomKey()
	if !ok || key != "value" {
		t.Fatalf("replacement corrupted key indexes: (%q, %v)", key, ok)
	}
}

func TestSafeEMapGetWithOptionsHonorsPreExpiringCondition(t *testing.T) {
	const key = "value"

	m := mapUtils.NewSafeEMap[string, valuesContainer]()
	m.SetExpiration(-time.Nanosecond)
	existing := &valuesContainer{Value1: 1, Value2: "existing"}
	m.Add(key, existing)

	conditionCalls := 0
	m.SetPreExpiringConditionFn(func(gotKey string, value *valuesContainer) bool {
		conditionCalls++
		if gotKey != key || value != existing {
			t.Errorf("condition received (%q, %p), want (%q, %p)", gotKey, value, key, existing)
		}
		return false
	})

	createCalled := false
	value := m.GetWithOptions(key, &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) {
			createCalled = true
			return &valuesContainer{Value1: 2, Value2: "replacement"}, true
		},
		DoFn: func(value *valuesContainer) {
			value.Value1 = 3
			value.Value2 = "kept"
		},
	})
	if createCalled || conditionCalls != 1 || value != existing {
		t.Fatal("a rejected expiration replaced the existing value")
	}
	if value.Value1 != 3 || value.Value2 != "kept" {
		t.Fatalf("DoFn did not receive the kept value: %+v", value)
	}

	m.SetPreExpiringConditionFn(func(_ string, value *valuesContainer) bool {
		conditionCalls++
		return value == existing
	})
	replacement := m.GetWithOptions(key, &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) {
			return &valuesContainer{Value1: 4, Value2: "replacement"}, true
		},
	})
	if conditionCalls != 2 || replacement == nil || replacement == existing {
		t.Fatal("an accepted expiration did not replace the existing value")
	}
	if replacement.Value1 != 4 || replacement.Value2 != "replacement" {
		t.Fatalf("unexpected replacement value: %+v", replacement)
	}
}

func TestSafeMapGetWithOptionsHonorsDisabledState(t *testing.T) {
	m := mapUtils.NewSafeMap[string, valuesContainer]()
	m.Disable()

	createCalled := false
	doCalled := false
	value := m.GetWithOptions("missing", &mapUtils.GetOptions[string, valuesContainer]{
		CreateFn: func() (*valuesContainer, bool) {
			createCalled = true
			return &valuesContainer{Value1: 1, Value2: "created"}, true
		},
		DoFn: func(*valuesContainer) { doCalled = true },
	})
	if value != nil || createCalled || doCalled || m.Exists("missing") {
		t.Fatal("disabled map created or handled a missing value")
	}
}
