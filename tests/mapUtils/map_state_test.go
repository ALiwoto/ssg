package tests

import (
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

func TestMapIsValidHandlesNilAndZeroValues(t *testing.T) {
	var nilSafeMap *ssg.SafeMap[int, valuesContainer]
	if nilSafeMap.IsValid() {
		t.Fatal("nil SafeMap is valid")
	}
	var zeroSafeMap ssg.SafeMap[int, valuesContainer]
	if zeroSafeMap.IsValid() {
		t.Fatal("zero-value SafeMap is valid")
	}
	if !ssg.NewSafeMap[int, valuesContainer]().IsValid() {
		t.Fatal("constructed SafeMap is invalid")
	}

	var nilSafeEMap *ssg.SafeEMap[int, valuesContainer]
	if nilSafeEMap.IsValid() {
		t.Fatal("nil SafeEMap is valid")
	}
	var zeroSafeEMap ssg.SafeEMap[int, valuesContainer]
	if zeroSafeEMap.IsValid() {
		t.Fatal("zero-value SafeEMap is valid")
	}
	validSafeEMap := ssg.NewSafeEMap[int, valuesContainer]()
	validSafeEMap.SetExpiration(time.Second)
	validSafeEMap.SetInterval(2 * time.Second)
	if !validSafeEMap.IsValid() {
		t.Fatal("configured SafeEMap is invalid")
	}

	var nilAdvancedMap *ssg.AdvancedMap[int, valuesContainer]
	if nilAdvancedMap.IsValid() {
		t.Fatal("nil AdvancedMap is valid")
	}
	var zeroAdvancedMap ssg.AdvancedMap[int, valuesContainer]
	if zeroAdvancedMap.IsValid() {
		t.Fatal("zero-value AdvancedMap is valid")
	}
	validAdvancedMap := ssg.NewAdvancedMap[int, valuesContainer]()
	validAdvancedMap.Set(1, valuesContainer{Value1: 1, Value2: "valid"})
	if !validAdvancedMap.IsValid() {
		t.Fatal("populated AdvancedMap is invalid")
	}
}

func TestSafeMapDisableFreezesEntries(t *testing.T) {
	m := ssg.NewSafeMap[string, valuesContainer]()
	m.Set("existing", valuesContainer{Value1: 1, Value2: "existing"})
	m.Disable()

	if !m.IsDisabled() {
		t.Fatal("map is not disabled after Disable")
	}

	m.Set("existing", valuesContainer{Value1: 2, Value2: "replacement"})
	m.Set("new", valuesContainer{Value1: 3, Value2: "new"})
	m.Delete("existing")
	m.Clear()
	m.ForEach(func(_ string, _ *valuesContainer) mapUtils.ForEachOperation {
		return mapUtils.ForEachOperationRemove
	})

	createCalled := false
	if value := m.GetOrCreate("missing", func() (*valuesContainer, bool) {
		createCalled = true
		return &valuesContainer{Value1: 4, Value2: "created"}, true
	}); value != nil {
		t.Fatalf("GetOrCreate added missing value while disabled: %v", *value)
	}
	if createCalled {
		t.Fatal("GetOrCreate called createFn while disabled")
	}

	existingCreateCalled := false
	existing := m.GetOrCreate("existing", func() (*valuesContainer, bool) {
		existingCreateCalled = true
		return &valuesContainer{Value1: 5, Value2: "replacement"}, true
	})
	if existingCreateCalled || existing == nil ||
		existing.Value1 != 1 || existing.Value2 != "existing" {
		t.Fatalf("GetOrCreate did not return frozen existing value: %v", existing)
	}
	if m.Length() != 1 || !m.Exists("existing") || m.Exists("new") {
		t.Fatalf("disabled map membership changed: %#v", m.ToNormalMap())
	}

	m.Enable()
	if m.IsDisabled() {
		t.Fatal("map is disabled after Enable")
	}
	m.Delete("existing")
	if !m.IsEmpty() {
		t.Fatalf("enabled map did not delete existing entry: %#v", m.ToNormalMap())
	}
}

func TestSafeEMapDisableFreezesEntriesAndExpiration(t *testing.T) {
	m := ssg.NewSafeEMap[string, valuesContainer]()
	m.SetExpiration(time.Hour)
	m.Set("existing", valuesContainer{Value1: 1, Value2: "existing"})
	m.Disable()

	if !m.IsDisabled() {
		t.Fatal("map is not disabled after Disable")
	}

	m.Set("existing", valuesContainer{Value1: 2, Value2: "replacement"})
	m.Set("new", valuesContainer{Value1: 3, Value2: "new"})
	m.Delete("existing")
	m.Clear()
	m.ForEach(func(_ string, _ *valuesContainer) mapUtils.ForEachOperation {
		return mapUtils.ForEachOperationRemove
	})

	createCalled := false
	if value := m.GetOrCreate("missing", func() (*valuesContainer, bool) {
		createCalled = true
		return &valuesContainer{Value1: 4, Value2: "created"}, true
	}); value != nil {
		t.Fatalf("GetOrCreate added missing value while disabled: %v", *value)
	}
	if createCalled {
		t.Fatal("GetOrCreate called createFn while disabled")
	}

	existingCreateCalled := false
	existing := m.GetOrCreate("existing", func() (*valuesContainer, bool) {
		existingCreateCalled = true
		return &valuesContainer{Value1: 5, Value2: "replacement"}, true
	})
	if existingCreateCalled || existing == nil ||
		existing.Value1 != 1 || existing.Value2 != "existing" {
		t.Fatalf("GetOrCreate did not return frozen existing value: %v", existing)
	}

	m.SetExpiration(-time.Nanosecond)
	m.DoCheck()
	if m.Length() != 1 || !m.Exists("existing") || m.Exists("new") {
		t.Fatalf("disabled map membership changed: %#v", m.ToNormalMap())
	}

	m.Enable()
	if m.IsDisabled() {
		t.Fatal("map is disabled after Enable")
	}
	m.DoCheck()
	if !m.IsEmpty() {
		t.Fatalf("enabled map did not expire existing entry: %#v", m.ToNormalMap())
	}
}
