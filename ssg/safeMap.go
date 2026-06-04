package ssg

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
	_, b := s.values[key]
	s.rUnlock()
	return b
}

func (s *SafeMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	defer s.unlock()

	if s._disabled {
		return
	}
	s.values[key] = value
}

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
			result[i] = s._default
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

func (s *SafeMap[TKey, TValue]) ToList() GenericList[*TValue] {
	return GetListFromArray(s.ToPointerArray())
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

	delete(s.values, key)
}

func (s *SafeMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
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
		return s._default
	}

	return *value
}

func (s *SafeMap[TKey, TValue]) SetDefault(value TValue) {
	s._default = value
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
			normalMap[k] = s._default
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
	return s.Length() > 0
}

// IsDisabled returns true if this map is disabled.
// Disabled maps won't be able to add new values, but will still be able to
// delete/read values.
func (s *SafeMap[TKey, TValue]) IsDisabled() bool {
	s.rLock()
	defer s.rUnlock()

	return s._disabled
}

// Disable will disable this map, meaning that it won't be able to add new
// values, but will still be able to delete/read values.
func (s *SafeMap[TKey, TValue]) Disable() {
	s.lock()
	defer s.unlock()

	s._disabled = true
}

// Enable will enable this map, meaning that it will be able to add new values.
func (s *SafeMap[TKey, TValue]) Enable() {
	s.lock()
	defer s.unlock()

	s._disabled = false
}
