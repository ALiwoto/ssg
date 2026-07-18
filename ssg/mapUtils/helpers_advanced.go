package mapUtils

import "sync"

func NewAdvancedMap[TKey comparable, TValue any]() *AdvancedMap[TKey, TValue] {
	return &AdvancedMap[TKey, TValue]{
		mut:           &sync.RWMutex{},
		values:        make(map[TKey]*TValue),
		sliceKeyIndex: make(map[TKey]int),
	}
}
