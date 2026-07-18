package mapUtils

import "sync"

// SafeMap is a safe map of type TIndex to pointers of type TValue.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
type SafeMap[TKey comparable, TValue any] struct {
	mut    *sync.RWMutex
	values map[TKey]*TValue

	// defaultValues field is the default value this map has to return in GetValue
	// method when the key is not found.
	// this field is only for value types, not pointers, we would still return nil for pointers.
	defaultValues TValue

	// disabled determines whether the map is disabled or not.
	disabled bool
}
