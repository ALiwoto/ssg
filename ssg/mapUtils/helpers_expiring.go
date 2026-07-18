package mapUtils

import (
	"sync"
	"time"
)

func NewEValue[T any](value T) *ExpiringValue[T] {
	return &ExpiringValue[T]{
		mut:       &sync.Mutex{},
		value:     value,
		timestamp: time.Now(),
	}
}

func NewSafeEMap[TKey comparable, TValue any]() *SafeEMap[TKey, TValue] {
	return &SafeEMap[TKey, TValue]{
		mut:           &sync.RWMutex{},
		values:        make(map[TKey]*ExpiringValue[*TValue]),
		sliceKeyIndex: make(map[TKey]int),
	}
}
