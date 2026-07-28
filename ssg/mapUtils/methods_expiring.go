package mapUtils

import (
	"math/rand"
	"time"

	"github.com/ALiwoto/ssg/ssg/commonUtils"
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
	}

	s.setNewValue(key, value)
}

// setNewValue replaces the value for an existing key or registers a new key in
// all of the map's internal indexes. The caller must hold the map's write lock.
func (s *SafeEMap[TKey, TValue]) setNewValue(key TKey, value *TValue) {
	_, exists := s.values[key]
	s.values[key] = NewEValue(value)
	if exists {
		return
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

	if s.disabled {
		return
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

// DeleteIf deletes key when condFn returns true for its non-nil value.
// condFn runs while the map is locked; re-using this map inside condFn will
// result in a deadlock.
func (s *SafeEMap[TKey, TValue]) DeleteIf(key TKey, condFn func(*TValue) bool) {
	if condFn == nil {
		return
	}

	s.lock()
	defer s.unlock()

	if s.disabled {
		return
	}

	expiringValue := s.values[key]
	if expiringValue == nil {
		return
	}

	value := expiringValue.GetValue(false)
	if value == nil {
		return
	}

	if condFn(value) {
		s.delete(key, false)
	}
}

// ForEach calls fn for each entry while holding the map's write lock.
// The callback must not call another method on this map or wait for a goroutine
// that does so, because the callback would prevent the lock from being released
// and cause a deadlock. The callback may start asynchronous map operations as
// long as it returns without waiting for them. Use the returned ForEachOperation
// to remove the current entry or stop the iteration.
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

func (s *SafeEMap[TKey, TValue]) GetWithOptions(
	key TKey,
	opts *GetOptions[TKey, TValue],
) *TValue {
	if opts == nil {
		s.rLock()
		defer s.rUnlock()

		value := s.values[key]
		if value == nil {
			return nil
		}

		return value.GetValue(true)
	}

	s.lock()
	defer s.unlock()

	expiringValue, exists := s.values[key]
	found := exists && expiringValue != nil
	if found && opts.CreateFn != nil {
		found = !expiringValue.IsExpired(s.expiration)
	}

	var value *TValue
	if found {
		value = expiringValue.GetValue(true)
	} else {
		if s.disabled || opts.CreateFn == nil {
			return nil
		}

		var ok bool
		value, ok = opts.CreateFn()
		if !ok {
			return nil
		}

		s.setNewValue(key, value)
	}

	if opts.DoFn != nil {
		opts.DoFn(value)
	}

	return value
}

func (s *SafeEMap[TKey, TValue]) Get(key TKey) *TValue {
	return s.GetWithOptions(key, nil)
}

func (s *SafeEMap[TKey, TValue]) GetOrCreate(
	key TKey,
	createFn commonUtils.PtrCreatorFunc[TValue],
) *TValue {
	return s.GetWithOptions(
		key,
		&GetOptions[TKey, TValue]{
			CreateFn: createFn,
		},
	)
}

// GetOrCreateDefault will call GetOrCreate with a default initializer.
func (s *SafeEMap[TKey, TValue]) GetOrCreateDefault(key TKey) *TValue {
	return s.GetOrCreate(key, commonUtils.DefaultPtrInitializer)
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

	if s.disabled {
		return
	}

	s.values = make(map[TKey]*ExpiringValue[*TValue])
	s.keys = nil
	s.sliceKeyIndex = make(map[TKey]int)
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
	if s == nil || s.mut == nil {
		return false
	}

	s.rLock()
	defer s.rUnlock()

	return s.values != nil && s.hasValidTimings()
}

// IsDisabled reports whether the map's entries are frozen. A disabled map
// remains readable, but its entries cannot be added, replaced, or removed.
// Expiration checks also leave entries untouched while the map is disabled.
func (s *SafeEMap[TKey, TValue]) IsDisabled() bool {
	s.rLock()
	defer s.rUnlock()

	return s.disabled
}

// Disable freezes the map's entries. Existing entries remain readable, but
// calls that would add, replace, delete, clear, or expire entries have no effect
// until Enable is called.
func (s *SafeEMap[TKey, TValue]) Disable() {
	s.lock()
	defer s.unlock()

	s.disabled = true
}

// Enable unfreezes the map, allowing its entries to be modified again.
func (s *SafeEMap[TKey, TValue]) Enable() {
	s.lock()
	defer s.unlock()

	s.disabled = false
}

func (s *SafeEMap[TKey, TValue]) hasValidTimings() bool {
	return s.expiration > time.Microsecond && s.checkInterval > time.Second
}

func (s *SafeEMap[TKey, TValue]) EnableChecking() {
	s.lock()
	defer s.unlock()

	s.checkingEnabled = true

	if s.isInCheckLoop {
		return
	}

	s.isInCheckLoop = true
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

func (s *SafeEMap[TKey, TValue]) SetPreExpiringConditionFn(fn func(key TKey, value *TValue) bool) {
	s.lock()
	defer s.unlock()

	s.preExpiringConditionFn = fn
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

	if s.disabled || len(s.values) == 0 {
		return
	}

	for key, current := range s.values {
		if current == nil || current.IsExpired(s.expiration) {
			if current != nil &&
				s.preExpiringConditionFn != nil &&
				!s.preExpiringConditionFn(key, current.GetValue(false)) {
				continue
			}

			s.delete(key, false)
			if s.onExpired != nil {
				go s.onExpired(key, s.getRealValue(current, false))
			}

			if s.onExpiredPtr != nil {
				go s.onExpiredPtr(key, current.GetValue(false))
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
	defer s.onCheckLoopFinished()

	for {
		time.Sleep(max(s.getCheckInterval(), time.Microsecond))

		status := s.getCheckStatus()
		if status == checkActionReturn {
			return
		} else if status == checkActionContinue {
			continue
		}

		s.DoCheck()
	}
}

func (s *SafeEMap[TKey, TValue]) onCheckLoopFinished() {
	s.lock()
	defer s.unlock()

	s.isInCheckLoop = false

	// EnableChecking may have run while this loop was exiting.
	if s.checkingEnabled {
		s.isInCheckLoop = true
		go s.checkLoop()
	}
}
