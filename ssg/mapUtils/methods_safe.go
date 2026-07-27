package mapUtils

import (
	"github.com/ALiwoto/ssg/ssg/internal"
	"github.com/ALiwoto/ssg/ssg/listUtils"
)

func (s *SafeMap[TKey, TValue]) lock() {
	s.mut.Lock()
}

func (s *SafeMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *SafeMap[TKey, TValue]) rLock() {
	s.mut.RLock()
}

func (s *SafeMap[TKey, TValue]) rUnlock() {
	s.mut.RUnlock()
}

func (s *SafeMap[TKey, TValue]) Exists(key TKey) bool {
	s.rLock()
	defer s.rUnlock()
	_, b := s.values[key]

	return b
}

func (s *SafeMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	defer s.unlock()

	if s.disabled {
		return
	}
	s.values[key] = value
}

// GetOrCreate returns the value of the key if it exists, otherwise it creates a new
// value using the provided createFn function and adds it to the map.
// If the map is disabled, it won't call createFn and will return nil instead.
func (s *SafeMap[TKey, TValue]) GetOrCreate(key TKey, createFn func() *TValue) *TValue {
	if createFn == nil {
		return s.Get(key)
	}

	s.lock()
	defer s.unlock()

	value, exists := s.values[key]
	if exists {
		return value
	}

	if s.disabled {
		return nil
	}

	newValue := createFn()
	s.values[key] = newValue
	return newValue
}

// GetOrCreateDefault will call GetOrCreate with a default initializer.
func (s *SafeMap[TKey, TValue]) GetOrCreateDefault(key TKey) *TValue {
	return s.GetOrCreate(key, internal.DefaultInitializer)
}

// ForEach calls fn for each entry while holding the map's write lock.
// The callback must not call another method on this map or wait for a goroutine
// that does so, because the callback would prevent the lock from being released
// and cause a deadlock. The callback may start asynchronous map operations as
// long as it returns without waiting for them. Use the returned ForEachOperation
// to remove the current entry or stop the iteration.
func (s *SafeMap[TKey, TValue]) ForEach(fn func(TKey, *TValue) ForEachOperation) {
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

// ForEachReadOnly calls fn for each entry while holding the map's read lock.
// The callback must not call another method on this map or wait for a goroutine
// that does so, because nested locking can deadlock, particularly when a writer
// is waiting. The callback may start asynchronous map operations as long as it
// returns without waiting for them. Remove operations returned by the callback
// are treated as continue or break operations and do not modify the map.
func (s *SafeMap[TKey, TValue]) ForEachReadOnly(fn func(TKey, *TValue) ForEachOperation) {
	if fn == nil {
		return
	}
	s.rLock()
	defer s.rUnlock()

myFor:
	for key, value := range s.values {
		switch fn(key, value) {
		case ForEachOperationContinue, ForEachOperationRemove:
			continue
		case ForEachOperationBreak, ForEachOperationRemoveBreak:
			break myFor
		}
	}
}

func (s *SafeMap[TKey, TValue]) ToArray() []TValue {
	s.rLock()
	defer s.rUnlock()

	result := make([]TValue, len(s.values))
	i := 0
	for _, v := range s.values {
		if v == nil {
			result[i] = s.defaultValues
			i++
			continue
		}

		result[i] = *v
		i++
	}

	return result
}

func (s *SafeMap[TKey, TValue]) ToPointerArray() []*TValue {
	s.rLock()
	defer s.rUnlock()

	result := make([]*TValue, len(s.values))
	i := 0
	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		result[i] = v
		i++
	}

	return result
}

func (s *SafeMap[TKey, TValue]) ToList() listUtils.GenericList[*TValue] {
	return listUtils.GetListFromArray(s.ToPointerArray())
}

func (s *SafeMap[TKey, TValue]) AddList(keyGetter func(*TValue) TKey, elements ...TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(&current), &current)
	}
}

func (s *SafeMap[TKey, TValue]) AddPointerList(keyGetter func(*TValue) TKey, elements ...*TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(current), current)
	}
}

func (s *SafeMap[TKey, TValue]) delete(key TKey, useLock bool) {
	if useLock {
		s.lock()
		defer s.unlock()
	}

	if s.disabled {
		return
	}

	delete(s.values, key)
}

func (s *SafeMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
}

// DeleteIf deletes key when condFn returns true for its non-nil value.
// condFn runs while the map is locked; re-using this map inside condFn will
// result in a deadlock.
func (s *SafeMap[TKey, TValue]) DeleteIf(key TKey, condFn func(*TValue) bool) {
	if condFn == nil {
		return
	}

	s.lock()
	defer s.unlock()

	if s.disabled {
		return
	}

	value := s.values[key]
	if value == nil {
		return
	}

	if condFn(value) {
		s.delete(key, false)
	}
}

func (s *SafeMap[TKey, TValue]) Get(key TKey) *TValue {
	s.rLock()
	defer s.rUnlock()

	value := s.values[key]
	return value
}

func (s *SafeMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.rLock()
	defer s.rUnlock()

	value := s.values[key]
	if value == nil {
		return s.defaultValues
	}

	return *value
}

func (s *SafeMap[TKey, TValue]) SetDefault(value TValue) {
	s.lock()
	defer s.unlock()

	s.defaultValues = value
}

// Set function sets the key of type TKey in this safe map to the value.
// the value should be of type TValue or *TValue, otherwise this function won't
// do anything at all.
func (s *SafeMap[TKey, TValue]) Set(key TKey, value any) {
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
func (s *SafeMap[TKey, TValue]) Clear() {
	s.lock()
	defer s.unlock()

	if s.disabled {
		return
	}

	if len(s.values) != 0 {
		s.values = make(map[TKey]*TValue)
	}
}

func (s *SafeMap[TKey, TValue]) Length() int {
	s.rLock()
	defer s.rUnlock()

	return len(s.values)
}

func (s *SafeMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *SafeMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	normalMap := make(map[TKey]TValue)
	s.rLock()
	defer s.rUnlock()

	for k, v := range s.values {
		if v == nil {
			normalMap[k] = s.defaultValues
			continue
		}

		normalMap[k] = *v
	}

	return normalMap
}

func (s *SafeMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *SafeMap[TKey, TValue]) IsValid() bool {
	if s == nil || s.mut == nil {
		return false
	}

	s.rLock()
	defer s.rUnlock()

	return s.values != nil
}

// IsDisabled reports whether the map's entries are frozen. A disabled map
// remains readable, but its entries cannot be added, replaced, or removed.
func (s *SafeMap[TKey, TValue]) IsDisabled() bool {
	s.rLock()
	defer s.rUnlock()

	return s.disabled
}

// Disable freezes the map's entries. Existing entries remain readable, but
// calls that would add, replace, delete, or clear entries have no effect until
// Enable is called.
func (s *SafeMap[TKey, TValue]) Disable() {
	s.lock()
	defer s.unlock()

	s.disabled = true
}

// Enable unfreezes the map, allowing its entries to be modified again.
func (s *SafeMap[TKey, TValue]) Enable() {
	s.lock()
	defer s.unlock()

	s.disabled = false
}
