package tests

import (
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

func TestSafeEMapGetOrCreateMaintainsKeyIndexes(t *testing.T) {
	m := ssg.NewSafeEMap[string, int]()
	m.SetExpiration(time.Hour)

	created := m.GetOrCreate("created", func() *int {
		value := 42
		return &value
	})
	if created == nil || *created != 42 {
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

	m.Set("expired", 1)
	m.SetExpiration(-time.Nanosecond)
	replaced := m.GetOrCreate("expired", func() *int {
		value := 2
		return &value
	})
	if replaced == nil || *replaced != 2 {
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
	m := ssg.NewSafeEMap[string, int]()
	m.Set("old", 1)
	m.Clear()

	if m.Length() != 0 {
		t.Fatalf("Length after Clear is %d, want 0", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after Clear", key)
	}

	m.Set("new", 2)
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
	m := ssg.NewSafeEMap[string, int]()
	m.SetExpiration(-time.Nanosecond)
	m.Set("expired-1", 1)
	m.Set("expired-2", 2)

	m.DoCheck()

	if m.Length() != 0 {
		t.Fatalf("Length after DoCheck is %d, want 0", m.Length())
	}
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned expired key %q after DoCheck", key)
	}

	m.Set("new", 3)
	if value := m.GetRandom(); value == nil || *value != 3 {
		t.Fatalf("GetRandom returned %v after DoCheck, want 3", value)
	}
	m.Delete("new")
	if key, ok := m.GetRandomKey(); ok {
		t.Fatalf("GetRandomKey returned stale key %q after deleting post-check entry", key)
	}
}
