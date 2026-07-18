package listUtils

import "github.com/ALiwoto/ssg/ssg/commonUtils"

type ListW[T comparable] struct {
	_values []T
}

type ListLike interface {
	IsEmpty() bool
	Length() int
}

type GenericList[T comparable] interface {
	ListLike
	commonUtils.Validator

	Find(element T) int
	Count(element T) int
	Counts(element ...T) int
	Contains(element T) bool
	ContainsAll(elements ...T) bool
	ContainsOne(elements ...T) bool
	Change(index int, element T)
	Exists(element T) bool
	Append(elements ...T)
	Add(elements ...T)
	RemoveAt(index int)
	RemoveOnce(element T)
	RemoveAll(element ...T)
	Remove(element T)
	AsArray() []T
	ToArray() []T
	Clear()
	Get(index int) T
}
