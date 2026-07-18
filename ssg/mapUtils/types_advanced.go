package mapUtils

import "sync"

// AdvancedMap is a safe map of type TIndex to pointers of type TValue with
// extra advanced features that you can't find in safe-map types.
// obviously, because of its extra features, it's slightly slower than other
// normal safe-map types.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
type AdvancedMap[TKey comparable, TValue any] struct {
	mut    *sync.RWMutex
	values map[TKey]*TValue
	// keys field is a slice of the map keys used in the map above. We put them in a slice
	// so that we can get a random key by choosing a random index.
	keys []TKey
	// We store the index of each key, so that when we remove an item, we can
	// quickly remove it from the slice above.
	sliceKeyIndex map[TKey]int

	// defaultValue field is the default value this map has to return in GetValue
	// method when the key is not found. (only for value, not pointers, we would still
	// return nil for pointers)
	defaultValue TValue
}
