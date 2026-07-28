package mapUtils

import (
	"sync"
	"time"
)

// ExpiringValue represents a value of type T that has an expiration time.
// This struct is thread safe.
type ExpiringValue[T any] struct {
	mut       *sync.Mutex
	value     T
	timestamp time.Time
}

// SafeEMap is a safe map of type TIndex to pointers of type TValue.
// this map is completely thread safe and is using internal lock when
// getting and setting variables.
// the difference of SafeEMap and SafeMap is that SafeEMap is using a checker loop
// for removing the expired values from itself.
type SafeEMap[TKey comparable, TValue any] struct {
	checkingEnabled bool
	isInCheckLoop   bool

	checkInterval time.Duration
	expiration    time.Duration
	mut           *sync.RWMutex
	values        map[TKey]*ExpiringValue[*TValue]
	// keys field is a slice of the map keys used in the map above. We put them in a slice
	// so that we can get a random key by choosing a random index.
	keys []TKey
	// We store the index of each key, so that when we remove an item, we can
	// quickly remove it from the slice above.
	sliceKeyIndex map[TKey]int
	defaultValue  TValue

	// disabled determines whether the map is disabled or not.
	disabled bool

	// preExpiringConditionFn is called when we want to cleanup a value due to its
	// expiration reaching. it can return false to cancel the cleanup of this value.
	// Calling the map's methods inside of this function will result in a deadlock.
	preExpiringConditionFn func(key TKey, value *TValue) bool

	// onExpired is the event function that will be called when a value with the certain
	// key on the map is expired. this event function will be called in a new goroutine.
	onExpired func(key TKey, value TValue)

	// onExpiredPtr is the event function that will be called when a value with the certain
	// key on the map is expired. this event function will NOT be called in a new goroutine.
	onExpiredPtr func(key TKey, value *TValue)
}
