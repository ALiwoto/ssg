package mapUtils

import (
	"math/rand"

	"github.com/ALiwoto/ssg/ssg/listUtils"
)

func (s *AdvancedMap[TKey, TValue]) lock() {
	s.mut.Lock()
}
func (s *AdvancedMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *AdvancedMap[TKey, TValue]) rLock() {
	s.mut.RLock()
}
func (s *AdvancedMap[TKey, TValue]) rUnlock() {
	s.mut.RUnlock()
}

func (s *AdvancedMap[TKey, TValue]) Exists(key TKey) bool {
	s.rLock()
	defer s.rUnlock()

	return len(s.values) != 0 && s.values[key] != nil
}

func (s *AdvancedMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	defer s.unlock()

	_, exists := s.values[key]
	s.values[key] = value
	if exists {
		return
	}

	s.keys = append(s.keys, key)

	// store the index of the map key
	index := len(s.keys) - 1
	s.sliceKeyIndex[key] = index
}
func (s *AdvancedMap[TKey, TValue]) GetRandom() *TValue {
	s.rLock()
	defer s.rUnlock()

	if len(s.keys) == 0 || len(s.values) == 0 {
		return nil
	}

	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]

	return value
}

func (s *AdvancedMap[TKey, TValue]) ForEach(fn func(TKey, *TValue) ForEachOperation) {
	if fn == nil {
		return
	}
	s.lock()
	defer s.unlock()

myFor:
	for key, value := range s.values {
		switch fn(key, value) {
		case ForEachOperationContinue:
			continue
		case ForEachOperationBreak:
			break myFor
		case ForEachOperationRemove:
			s.delete(key, false)
		case ForEachOperationRemoveBreak:
			s.delete(key, false)
			break myFor
		}
	}
}

func (s *AdvancedMap[TKey, TValue]) GetRandomValue() TValue {
	s.rLock()
	defer s.rUnlock()

	if len(s.keys) == 0 {
		return s.defaultValue
	}

	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]

	if value == nil {
		return s.defaultValue
	}

	return *value
}

func (s *AdvancedMap[TKey, TValue]) GetRandomKey() (key TKey, ok bool) {
	s.rLock()
	defer s.rUnlock()

	if len(s.keys) == 0 {
		return
	}

	key = s.keys[rand.Intn(len(s.keys))]
	ok = true
	return
}

func (s *AdvancedMap[TKey, TValue]) ToArray() []TValue {
	var array []TValue
	s.rLock()
	defer s.rUnlock()

	for _, v := range s.values {
		if v == nil {
			array = append(array, s.defaultValue)
			continue
		}

		array = append(array, *v)
	}

	return array
}

func (s *AdvancedMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue
	s.rLock()
	defer s.rUnlock()

	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v)
	}

	return array
}

func (s *AdvancedMap[TKey, TValue]) ToList() listUtils.GenericList[*TValue] {
	list := listUtils.GetEmptyList[*TValue]()
	s.rLock()
	defer s.rUnlock()

	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v)
	}

	return list
}

func (s *AdvancedMap[TKey, TValue]) AddList(keyGetter func(*TValue) TKey, elements ...TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(&current), &current)
	}
}

func (s *AdvancedMap[TKey, TValue]) AddPointerList(keyGetter func(*TValue) TKey, elements ...*TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(current), current)
	}
}

func (s *AdvancedMap[TKey, TValue]) delete(key TKey, useLock bool) {
	if useLock {
		s.lock()
		defer s.unlock()
	}

	// get index in key slice for key
	index, exists := s.sliceKeyIndex[key]
	if !exists {
		// item does not exist
		return
	}

	delete(s.sliceKeyIndex, key)

	wasLastIndex := len(s.keys)-1 == index

	// remove key from slice of keys
	s.keys[index] = s.keys[len(s.keys)-1]
	s.keys = s.keys[:len(s.keys)-1]

	// we just swapped the last element to another position.
	// so we need to update it's index (if it was not in last position)
	if !wasLastIndex {
		otherKey := s.keys[index]
		s.sliceKeyIndex[otherKey] = index
	}

	delete(s.values, key)
}

func (s *AdvancedMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
}

func (s *AdvancedMap[TKey, TValue]) Get(key TKey) *TValue {
	s.rLock()
	defer s.rUnlock()

	value := s.values[key]
	return value
}

func (s *AdvancedMap[TKey, TValue]) GetOrCreate(key TKey, createFn func() *TValue) *TValue {
	if createFn == nil {
		return nil
	}

	s.lock()
	defer s.unlock()

	value := s.values[key]
	if value == nil {
		value = createFn()
		s.values[key] = value
	}

	return value
}

func (s *AdvancedMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.rLock()
	defer s.rUnlock()

	value := s.values[key]
	if value == nil {
		return s.defaultValue
	}

	return *value
}

func (s *AdvancedMap[TKey, TValue]) SetDefault(value TValue) {
	s.lock()
	defer s.unlock()

	s.defaultValue = value
}

// Set function sets the key of type TKey in this safe map to the value.
// the value should be of type TValue or *TValue, otherwise this function won't
// do anything at all.
func (s *AdvancedMap[TKey, TValue]) Set(key TKey, value any) {
	correctValue, ok := value.(*TValue)
	if !ok {
		anotherValue, ok := value.(TValue)
		if !ok {
			return
		}

		correctValue = &anotherValue
	}

	s.Add(key, correctValue)
}

// Clear will clear the whole map.
func (s *AdvancedMap[TKey, TValue]) Clear() {
	s.lock()
	defer s.unlock()

	if len(s.values) != 0 {
		s.values = make(map[TKey]*TValue)
	}
}

func (s *AdvancedMap[TKey, TValue]) Length() int {
	s.rLock()
	defer s.rUnlock()

	return len(s.values)
}

func (s *AdvancedMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *AdvancedMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)
	s.rLock()
	defer s.rUnlock()

	for k, v := range s.values {
		if v == nil {
			m[k] = s.defaultValue
			continue
		}

		m[k] = *v
	}

	return m
}

func (s *AdvancedMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *AdvancedMap[TKey, TValue]) IsValid() bool {
	return s != nil && s.Length() > 0
}

//---------------------------------------------------------
