package commonUtils

type Validator interface {
	IsValid() bool
}

// PtrCreatorFunc is a pointer value creator function.
type PtrCreatorFunc[TValue any] = func() (value *TValue, ok bool)
