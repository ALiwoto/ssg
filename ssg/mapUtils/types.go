package mapUtils

import "github.com/ALiwoto/ssg/ssg/commonUtils"

// ForEachOperation describes an operation that has to be returned
// from a ForEach method.
type ForEachOperation int

type checkAction uint8

type GetOptions[TKey comparable, TValue any] struct {
	// CreateFn will be called when the value is missing inside of the function.
	// the second return value, `ok bool` is existing flag; if it's true it means
	// the value must be injected into the map (even if it's nil), otherwise
	// it won't.
	//
	// IMPORTANT: DO **NOT** re-use the map inside of this function, otherwise it will
	// result in a deadlock.
	CreateFn commonUtils.PtrCreatorFunc[TValue]

	// DoFn will be called only when the value is found, even if the value is nil.
	// It runs while the map's read lock protects the entry from removal or replacement.
	// Multiple DoFn callbacks may run concurrently, including for the same key, so
	// TValue must provide any synchronization required for its own mutable state.
	//
	// IMPORTANT: DO **NOT** re-use the map inside of this function, otherwise it will
	// result in a deadlock.
	//
	// If DoFn acquires a lock that remains held after DoFn returns, code holding that
	// lock must not call this map's methods. Another GetWithOptions call can hold the
	// map lock while waiting for the value lock, resulting in a lock-order deadlock.
	DoFn func(*TValue)
}
