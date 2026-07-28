package tests

import (
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

func TestSafeEMapGetOrCreateMaintainsKeyIndexes(t *testing.T) {
	m := ssg.NewSafeEMap[string, valuesContainer]()
	m.SetExpiration(time.Hour)

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

	m.Delete("created")
	if m.Exists("created") || m.Length() != 0 {
		t.Fatalf("Delete did not remove GetOrCreate entry; length is %d", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after Delete", key)
	}

	m.Set("expired", valuesContainer{Value1: 1, Value2: "expired"})
	m.SetExpiration(-time.Nanosecond)
	replaced := m.GetOrCreate("expired", func() (*valuesContainer, bool) {
		return &valuesContainer{Value1: 2, Value2: "replacement"}, true
	})
	if replaced == nil || replaced.Value1 != 2 || replaced.Value2 != "replacement" {
		t.Fatalf("GetOrCreate returned %v for expired entry, want 2", replaced)
	}
	if m.Length() != 1 {
		t.Fatalf("Length after replacing expired entry is %d, want 1", m.Length())
	}

	m.Delete("expired")
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned duplicate key %q after replacing and deleting expired entry", key)
	}
}

func TestSafeEMapClearClearsKeyIndexes(t *testing.T) {
	m := ssg.NewSafeEMap[string, valuesContainer]()
	m.Set("old", valuesContainer{Value1: 1, Value2: "old"})
	m.Clear()

	if m.Length() != 0 {
		t.Fatalf("Length after Clear is %d, want 0", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after Clear", key)
	}

	m.Set("new", valuesContainer{Value1: 2, Value2: "new"})
	key, ok := m.GetRandomKey()
	if !ok || key != "new" {
		t.Fatalf("GetRandomKey returned (%q, %v), want (%q, true)", key, ok, "new")
	}
	m.Delete("new")
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after deleting post-Clear entry", key)
	}
}

func TestSafeEMapDoCheckRemovesKeyIndexes(t *testing.T) {
	m := ssg.NewSafeEMap[string, valuesContainer]()
	m.SetExpiration(-time.Nanosecond)
	m.Set("expired-1", valuesContainer{Value1: 1, Value2: "expired one"})
	m.Set("expired-2", valuesContainer{Value1: 2, Value2: "expired two"})

	m.DoCheck()

	if m.Length() != 0 {
		t.Fatalf("Length after DoCheck is %d, want 0", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned expired key %q after DoCheck", key)
	}

	m.Set("new", valuesContainer{Value1: 3, Value2: "new"})
	if value := m.GetRandom(); value == nil || value.Value1 != 3 || value.Value2 != "new" {
		t.Fatalf("GetRandom returned %v after DoCheck, want 3", value)
	}
	m.Delete("new")
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after deleting post-check entry", key)
	}
}
