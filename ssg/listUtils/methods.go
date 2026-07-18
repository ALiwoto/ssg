package listUtils

func (l *ListW[T]) Find(element T) int {
	for i, v := range l._values {
		if v == element {
			return i
		}
	}

	return -1
}

func (l *ListW[T]) Count(element T) int {
	count := 0
	for _, v := range l._values {
		if v == element {
			count++
		}
	}

	return count
}

func (l *ListW[T]) Counts(element ...T) int {
	count := 0
	for _, v := range l._values {
		for _, current := range element {
			if v == current {
				count++
			}
		}
	}

	return count
}

func (l *ListW[T]) Contains(element T) bool {
	return l.Find(element) != -1
}

func (l *ListW[T]) ContainsAll(elements ...T) bool {
	for _, current := range elements {
		if !l.Contains(current) {
			return false
		}
	}

	return true
}

func (l *ListW[T]) ContainsOne(elements ...T) bool {
	for _, current := range elements {
		if l.Contains(current) {
			return true
		}
	}

	return false
}

func (l *ListW[T]) Change(index int, element T) {
	if index < 0 || index >= len(l._values) {
		return
	}

	l._values[index] = element
}

func (l *ListW[T]) Exists(element T) bool {
	return l.Find(element) != -1
}

func (l *ListW[T]) Append(elements ...T) {
	l._values = append(l._values, elements...)
}

func (l *ListW[T]) Add(elements ...T) {
	l._values = append(l._values, elements...)
}

func (l *ListW[T]) RemoveAt(index int) {
	l._values = append(l._values[:index], l._values[index+1:]...)
}

func (l *ListW[T]) RemoveOnce(element T) {
	index := l.Find(element)
	if index != -1 {
		l.RemoveAt(index)
	}
}

func (l *ListW[T]) RemoveAll(element ...T) {
	var newVal []T
	for _, current := range element {
		for _, v := range l._values {
			if v != current {
				newVal = append(newVal, v)
			}
		}
	}

	l._values = newVal
}

func (l *ListW[T]) Remove(element T) {
	l.RemoveOnce(element)
}

// AsArray returns a copy of the value of this list as an array.
// please do notice that if you make changes to the underlying values of
// that array, change won't be applied to the list.
func (l *ListW[T]) AsArray() []T {
	var arr = make([]T, len(l._values))
	copy(arr, l._values)
	return arr
}

// ToArray is equivalent to AsArray method in any way.
// it returns a copy of the value of this list as an array.
// please do notice that if you make changes to the underlying values of
// that array, change won't be applied to the list.
func (l *ListW[T]) ToArray() []T {
	return l.AsArray()
}

// Clear method clears the whole list.
func (l *ListW[T]) Clear() {
	l._values = nil
}

func (l *ListW[T]) Get(index int) T {
	return l._values[index]
}

func (l *ListW[T]) IsThreadSafe() bool {
	return true
}

func (l *ListW[T]) IsEmpty() bool {
	return len(l._values) == 0
}

func (l *ListW[T]) Length() int {
	return len(l._values)
}

func (l *ListW[T]) IsValid() bool {
	return len(l._values) > 0
}

//---------------------------------------------------------
