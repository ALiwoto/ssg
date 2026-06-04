package rangeValues

import (
	"strconv"
	"strings"

	"github.com/ALiwoto/ssg/ssg/internal"
)

func ParseIntContainer[T Integer](value string) *IntContainer[T] {
	parts := strings.Split(value, ":")
	if len(parts) < 1 {
		return nil
	}

	intValue, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil
	}

	return &IntContainer[T]{
		Value: T(intValue),
		Flags: parts[1:],
	}
}

func ParseIntArray[T Integer](value string) []T {
	arr := internal.SplitN(value, -1, ",", " ", "[", "]")
	var myInts []T

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, T(theValue))
	}

	return myInts
}
