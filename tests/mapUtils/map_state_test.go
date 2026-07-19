package tests

import (
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

func TestMapIsValidHandlesNilAndZeroValues(t *testing.T) {
	var nilSafeMap *ssg.SafeMap[int, int]
	if nilSafeMap.IsValid() {
		t.Fatal("nil SafeMap is valid")
	}
	var zeroSafeMap ssg.SafeMap[int, int]
	if zeroSafeMap.IsValid() {
		t.Fatal("zero-value SafeMap is valid")
	}
	if !ssg.NewSafeMap[int, int]().IsValid() {
		t.Fatal("constructed SafeMap is invalid")
	}

	var nilSafeEMap *ssg.SafeEMap[int, int]
	if nilSafeEMap.IsValid() {
		t.Fatal("nil SafeEMap is valid")
	}
	var zeroSafeEMap ssg.SafeEMap[int, int]
	if zeroSafeEMap.IsValid() {
		t.Fatal("zero-value SafeEMap is valid")
	}
	validSafeEMap := ssg.NewSafeEMap[int, int]()
	validSafeEMap.SetExpiration(time.Second)
	validSafeEMap.SetInterval(2 * time.Second)
	if !validSafeEMap.IsValid() {
		t.Fatal("configured SafeEMap is invalid")
	}

	var nilAdvancedMap *ssg.AdvancedMap[int, int]
	if nilAdvancedMap.IsValid() {
		t.Fatal("nil AdvancedMap is valid")
	}
	var zeroAdvancedMap ssg.AdvancedMap[int, int]
	if zeroAdvancedMap.IsValid() {
		t.Fatal("zero-value AdvancedMap is valid")
	}
	validAdvancedMap := ssg.NewAdvancedMap[int, int]()
	validAdvancedMap.Set(1, 1)
	if !validAdvancedMap.IsValid() {
		t.Fatal("populated AdvancedMap is invalid")
	}
}

func TestSafeMapDisableFreezesEntries(t *testing.T) {
	m := ssg.NewSafeMap[string, int]()
	m.Set("existing", 1)
	m.Disable()

	if !m.IsDisabled() {
		t.Fatal("map is not disabled after Disable")
	}

	m.Set("existing", 2)
	m.Set("new", 3)
	m.Delete("existing")
	m.Clear()
	m.ForEach(func(_ string, _ *int) mapUtils.ForEachOperation {
		return mapUtils.ForEachOperationRemove
	})

	createCalled := false
	if value := m.GetOrCreate("missing", func() *int {
		createCalled = true
		created := 4
		return &created
	}); value != nil {
		t.Fatalf("GetOrCreate added missing value while disabled: %v", *value)
	}
	if createCalled {
		t.Fatal("GetOrCreate called createFn while disabled")
	}

	existingCreateCalled := false
	existing := m.GetOrCreate("existing", func() *int {
		existingCreateCalled = true
		created := 5
		return &created
	})
	if existingCreateCalled || existing == nil || *existing != 1 {
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
	m := ssg.NewSafeEMap[string, int]()
	m.SetExpiration(time.Hour)
	m.Set("existing", 1)
	m.Disable()

	if !m.IsDisabled() {
		t.Fatal("map is not disabled after Disable")
	}

	m.Set("existing", 2)
	m.Set("new", 3)
	m.Delete("existing")
	m.Clear()
	m.ForEach(func(_ string, _ *int) mapUtils.ForEachOperation {
		return mapUtils.ForEachOperationRemove
	})

	createCalled := false
	if value := m.GetOrCreate("missing", func() *int {
		createCalled = true
		created := 4
		return &created
	}); value != nil {
		t.Fatalf("GetOrCreate added missing value while disabled: %v", *value)
	}
	if createCalled {
		t.Fatal("GetOrCreate called createFn while disabled")
	}

	existingCreateCalled := false
	existing := m.GetOrCreate("existing", func() *int {
		existingCreateCalled = true
		created := 5
		return &created
	})
	if existingCreateCalled || existing == nil || *existing != 1 {
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
