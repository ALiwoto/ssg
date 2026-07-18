package mapUtils

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ALiwoto/ssg/ssg/listUtils"
)

func (e *ExpiringValue[T]) SetTime(t time.Time) {
	e.mut.Lock()
	e.timestamp = t
	e.mut.Unlock()
}

func (e *ExpiringValue[T]) GetTime() time.Time {
	e.mut.Lock()
	defer e.mut.Unlock()

	return e.timestamp
}

func (e *ExpiringValue[T]) Reset() {
	e.SetTime(time.Now())
}

func (e *ExpiringValue[T]) IsExpired(duration time.Duration) bool {
	e.mut.Lock()
	defer e.mut.Unlock()

	return time.Since(e.timestamp) > duration
}

func (e *ExpiringValue[T]) SetValue(value T) {
	e.mut.Lock()
	e.value = value
	e.mut.Unlock()
}

func (e *ExpiringValue[T]) GetValue(shouldReset bool) T {
	e.mut.Lock()
	defer e.mut.Unlock()

	if shouldReset {
		e.timestamp = time.Now()
	}
	return e.value
}

//---------------------------------------------------------

func (s *SafeEMap[TKey, TValue]) lock() {
	s.mut.Lock()
}

func (s *SafeEMap[TKey, TValue]) unlock() {
	s.mut.Unlock()
}

func (s *SafeEMap[TKey, TValue]) rLock() {
	s.mut.RLock()
}

func (s *SafeEMap[TKey, TValue]) rUnlock() {
	s.mut.RUnlock()
}

func (s *SafeEMap[TKey, TValue]) Exists(key TKey) bool {
	s.rLock()
	defer s.rUnlock()

	_, b := s.values[key]
	return b
}

func (s *SafeEMap[TKey, TValue]) AddList(keyGetter func(*TValue) TKey, elements ...TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(&current), &current)
	}
}

func (s *SafeEMap[TKey, TValue]) ToList() listUtils.GenericList[*TValue] {
	list := listUtils.GetEmptyList[*TValue]()
	s.rLock()
	defer s.rUnlock()

	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		list.Add(v.GetValue(false))
	}

	return list
}

func (s *SafeEMap[TKey, TValue]) AddPointerList(keyGetter func(*TValue) TKey, elements ...*TValue) {
	if len(elements) == 0 || keyGetter == nil {
		return
	}

	for _, current := range elements {
		s.Add(keyGetter(current), current)
	}
}

func (s *SafeEMap[TKey, TValue]) Add(key TKey, value *TValue) {
	s.lock()
	defer s.unlock()

	if s.disabled {
		return
	}

	old := s.values[key]
	if old != nil {
		// don't allocate new memory if we already have the expiring-value struct in
		// the map... just set the new value and reset the time
		old.SetValue(value)
		old.Reset()
		return
	} else {
		s.values[key] = NewEValue(value)
	}

	s.keys = append(s.keys, key)

	// store the index of the map key
	index := len(s.keys) - 1
	s.sliceKeyIndex[key] = index
}

func (s *SafeEMap[TKey, TValue]) delete(key TKey, useLock bool) {
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
	// so we need to update its index (if it was not in last position)
	if !wasLastIndex {
		otherKey := s.keys[index]
		s.sliceKeyIndex[otherKey] = index
	}

	delete(s.values, key)
}

func (s *SafeEMap[TKey, TValue]) Delete(key TKey) {
	s.delete(key, true)
}

func (s *SafeEMap[TKey, TValue]) ForEach(fn func(TKey, *TValue) ForEachOperation) {
	if fn == nil {
		return
	}
	s.lock()
	defer s.unlock()

	var tmpValue *TValue

myFor:
	for key, value := range s.values {
		if value == nil {
			tmpValue = nil
		} else {
			tmpValue = value.GetValue(true)
		}

		switch fn(key, tmpValue) {
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

func (s *SafeEMap[TKey, TValue]) GetRandom() *TValue {
	s.rLock()
	defer s.rUnlock()

	if len(s.keys) == 0 || len(s.values) == 0 {
		return nil
	}

	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]

	return value.GetValue(true)
}

func (s *SafeEMap[TKey, TValue]) GetRandomValue() TValue {
	s.rLock()
	defer s.rUnlock()

	if len(s.keys) == 0 || len(s.values) == 0 {
		return s.defaultValue
	}

	randomIndex := rand.Intn(len(s.keys))
	key := s.keys[randomIndex]
	value := s.values[key]

	return s.getRealValue(value, true)
}

func (s *SafeEMap[TKey, TValue]) GetRandomKey() (key TKey, ok bool) {
	s.rLock()
	defer s.rUnlock()

	if len(s.keys) == 0 {
		return
	}

	ok = true
	key = s.keys[rand.Intn(len(s.keys))]
	return
}

func (s *SafeEMap[TKey, TValue]) Get(key TKey) *TValue {
	s.rLock()
	defer s.rUnlock()

	value := s.values[key]
	if value == nil {
		return nil
	}

	return value.GetValue(true)
}

func (s *SafeEMap[TKey, TValue]) GetOrCreate(key TKey, createFn func() *TValue) *TValue {
	if createFn == nil {
		return s.Get(key)
	}

	s.lock()
	defer s.unlock()

	if s.disabled {
		return nil
	}

	value, exists := s.values[key]
	if exists && !value.IsExpired(s.expiration) {
		return value.GetValue(true)
	}

	newValue := createFn()
	s.values[key] = NewEValue(newValue)
	return newValue
}

func (s *SafeEMap[TKey, TValue]) GetValue(key TKey) TValue {
	s.rLock()
	defer s.rUnlock()

	value := s.values[key]
	return s.getRealValue(value, true)
}

func (s *SafeEMap[TKey, TValue]) SetDefault(value TValue) {
	s.lock()
	defer s.unlock()

	s.defaultValue = value
}

// Set function sets the key of type TKey in this safe map to the value.
// the value should be of type TValue or *TValue, otherwise this function won't
// do anything at all.
func (s *SafeEMap[TKey, TValue]) Set(key TKey, value any) {
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
func (s *SafeEMap[TKey, TValue]) Clear() {
	s.lock()
	defer s.unlock()

	if len(s.values) != 0 {
		s.values = make(map[TKey]*ExpiringValue[*TValue])
	}

}

func (s *SafeEMap[TKey, TValue]) Length() int {
	s.rLock()
	defer s.rUnlock()

	return len(s.values)
}

func (s *SafeEMap[TKey, TValue]) IsEmpty() bool {
	return s.Length() == 0
}

func (s *SafeEMap[TKey, TValue]) ToNormalMap() map[TKey]TValue {
	m := make(map[TKey]TValue)

	s.rLock()
	defer s.rUnlock()

	for k, v := range s.values {
		if v == nil {
			m[k] = s.defaultValue
			continue
		}

		realValue := v.GetValue(false)
		if realValue == nil {
			m[k] = s.defaultValue
			continue
		}

		m[k] = *realValue
	}

	return m
}

func (s *SafeEMap[TKey, TValue]) ToArray() []TValue {
	var array []TValue

	s.rLock()
	defer s.rUnlock()

	for _, v := range s.values {
		if v == nil {
			array = append(array, s.defaultValue)
			continue
		}

		realValue := v.GetValue(false)
		if realValue == nil {
			array = append(array, s.defaultValue)
			continue
		}

		array = append(array, *realValue)
	}

	return array
}

func (s *SafeEMap[TKey, TValue]) ToPointerArray() []*TValue {
	var array []*TValue

	s.rLock()
	defer s.rUnlock()

	for _, v := range s.values {
		if v == nil {
			// most likely impossible, this checker is here just for more safety.
			continue
		}

		array = append(array, v.GetValue(false))
	}

	return array
}

func (s *SafeEMap[TKey, TValue]) IsThreadSafe() bool {
	return true
}

func (s *SafeEMap[TKey, TValue]) IsValid() bool {
	s.rLock()
	defer s.rUnlock()

	return s != nil && len(s.values) > 0 && s.hasValidTimings()
}

// IsDisabled returns true if this map is disabled.
// Disabled maps won't be able to add new values, but will still be able to
// delete/read values.
func (s *SafeEMap[TKey, TValue]) IsDisabled() bool {
	s.rLock()
	defer s.rUnlock()

	return s.disabled
}

// Disable will disable this map, turning it into a read-only state.
func (s *SafeEMap[TKey, TValue]) Disable() {
	s.lock()
	defer s.unlock()

	s.disabled = true
}

// Enable will enable this map, meaning that it will be able to add new values.
func (s *SafeEMap[TKey, TValue]) Enable() {
	s.lock()
	defer s.unlock()

	s.disabled = false
}

func (s *SafeEMap[TKey, TValue]) hasValidTimings() bool {
	return s.expiration > time.Microsecond && s.checkInterval > time.Second
}

func (s *SafeEMap[TKey, TValue]) EnableChecking() {
	if s.checkerMut == nil {
		s.checkerMut = &sync.Mutex{}
	}

	// this lock here makes sure that only 1 checkLoop is running at a time.
	s.checkerMut.Lock()
	defer s.checkerMut.Unlock()

	s.lock()
	defer s.unlock()

	if s.checkingEnabled {
		return
	}

	s.checkingEnabled = true
	go s.checkLoop()
}

func (s *SafeEMap[TKey, TValue]) DisableChecking() {
	s.lock()
	defer s.unlock()

	if !s.checkingEnabled {
		return
	}

	s.checkingEnabled = false
}

func (s *SafeEMap[TKey, TValue]) IsChecking() bool {
	s.rLock()
	defer s.rUnlock()

	return s.checkingEnabled
}

func (s *SafeEMap[TKey, TValue]) SetExpiration(duration time.Duration) {
	s.lock()
	defer s.unlock()

	s.expiration = duration
}

func (s *SafeEMap[TKey, TValue]) SetOnExpired(event func(key TKey, value TValue)) {
	s.lock()
	defer s.unlock()

	s.onExpired = event
}

func (s *SafeEMap[TKey, TValue]) SetOnExpiredPtr(event func(key TKey, value *TValue)) {
	s.lock()
	defer s.unlock()

	s.onExpiredPtr = event
}

func (s *SafeEMap[TKey, TValue]) SetInterval(duration time.Duration) {
	s.lock()
	defer s.unlock()

	s.checkInterval = duration
}

// getRealValue returns the real value of the ExpiringValue struct,
// or the default value if the ExpiringValue is nil or has a nil value.
// IMPORTANT: this function does not lock the map, so it should be called within a lock.
func (s *SafeEMap[TKey, TValue]) getRealValue(
	eValue *ExpiringValue[*TValue],
	shouldReset bool,
) TValue {
	if eValue == nil {
		return s.defaultValue
	}

	realValue := eValue.GetValue(shouldReset)
	if realValue == nil {
		return s.defaultValue
	}

	return *realValue
}

// DoCheck iterates over the map and checks for expired variables and removes them.
// if the `onExpired` member of the map is set, it will call them.
func (s *SafeEMap[TKey, TValue]) DoCheck() {
	s.lock()
	defer s.unlock()

	if len(s.values) == 0 {
		return
	}

	for i, current := range s.values {
		if current == nil || current.IsExpired(s.expiration) {
			delete(s.values, i)
			if s.onExpired != nil {
				go s.onExpired(i, s.getRealValue(current, false))
			}

			if s.onExpiredPtr != nil {
				go s.onExpiredPtr(i, current.GetValue(false))
			}
		}
	}
}

func (s *SafeEMap[TKey, TValue]) getCheckStatus() checkAction {
	if s == nil {
		return checkActionReturn
	}

	s.rLock()
	defer s.rUnlock()

	if !s.checkingEnabled {
		return checkActionReturn
	}

	if len(s.values) == 0 {
		return checkActionContinue
	}

	return checkActionNormal
}

func (s *SafeEMap[TKey, TValue]) getCheckInterval() time.Duration {
	s.rLock()
	defer s.rUnlock()

	return s.checkInterval
}

func (s *SafeEMap[TKey, TValue]) checkLoop() {
	for {
		time.Sleep(max(s.getCheckInterval(), time.Microsecond))

		status := s.getCheckStatus()
		if status == checkActionReturn {
			s.DisableChecking()
			return
		} else if status == checkActionContinue {
			continue
		}

		s.DoCheck()
	}
}
