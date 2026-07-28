// Package mapUtils provides generic, concurrency-safe map implementations.
//
// # Copy safety
//
// The map types store pointers to TValue, but some APIs copy TValue by accepting
// or returning it by value. If TValue must not be copied after first use, as with
// sync.Mutex, sync.RWMutex, sync.Once, sync.WaitGroup, and the sync/atomic types,
// avoid value-oriented APIs such as GetValue, GetRandomValue, ToArray,
// ToNormalMap, AddList, value-form Set, SetDefault, and SetOnExpired.
//
// Use pointer-oriented APIs instead: Add, AddPointerList, Get, GetWithOptions,
// GetRandom, ToPointerArray, ToList, pointer-form Set, and SetOnExpiredPtr. Go does
// not prevent copying synchronization values at runtime; go vet only provides
// best-effort static detection.
package mapUtils
