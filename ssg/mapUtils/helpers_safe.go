package mapUtils

import "sync"

func NewSafeMap[TKey comparable, TValue any]() *SafeMap[TKey, TValue] {
	return &SafeMap[TKey, TValue]{
		mut:    &sync.RWMutex{},
		values: make(map[TKey]*TValue),
	}
}
