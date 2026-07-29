// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"strconv"
)

//---------------------------------------------------------

func (e *EndpointError) Error() string {
	return strconv.Itoa(e.ErrorCode) + ": " + e.Message
}

//---------------------------------------------------------

// Next method returns the next id and will move the current cursor to
// the next id and will return that.
// current integer won't be available after this method call.
func (n *NumIdGenerator[T]) Next() T {
	n.mut.Lock()
	defer n.mut.Unlock()
	n.current++
	return n.current
}

// Set method will set the current value to the given value.
// please do notice that with calling n.Next() method, you will get
// the next id and not the current value.
func (n *NumIdGenerator[T]) Set(current T) {
	n.mut.Lock()
	defer n.mut.Unlock()
	n.current = current
}

// SafeSet method will set the current value to the given value
// if and only if the passed value is more than current value of
// the generator.
// please do notice that with calling n.Next() method, you will get
// the next id and not the current value.
func (n *NumIdGenerator[T]) SafeSet(current T) {
	n.mut.Lock()
	defer n.mut.Unlock()
	if current > n.current {
		n.current = current
	}
}

// Reset method will set the current value to 0.
func (n *NumIdGenerator[T]) Reset() {
	n.mut.Lock()
	defer n.mut.Unlock()
	n.current = 0
}

//---------------------------------------------------------
