package listUtils

func GetEmptyList[T comparable]() GenericList[T] {
	return &ListW[T]{}
}

func GetListFromArray[T comparable](array []T) GenericList[T] {
	return &ListW[T]{
		_values: array,
	}
}
