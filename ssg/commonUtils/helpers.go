package commonUtils

// DefaultInitializer will initialize a new value of type T and
// return a pointer to it.
func DefaultInitializer[T any]() *T {
	var value T
	return &value
}

// DefaultPtrInitializer
func DefaultPtrInitializer[T any]() (*T, bool) {
	var value T
	return &value, true
}
