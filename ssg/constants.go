// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

const (
	// ForEachOperationBreak will just continue the loop without doing anything.
	ForEachOperationContinue = mapUtils.ForEachOperationContinue

	// ForEachOperationBreak will just break the loop without doing anything.
	ForEachOperationBreak = mapUtils.ForEachOperationBreak

	// ForEachOperationBreak will just remove the current item from the list
	// and continue the loop.
	ForEachOperationRemove = mapUtils.ForEachOperationRemove

	// ForEachOperationBreak will remove the current item from the list
	// and break the loop.
	ForEachOperationRemoveBreak = mapUtils.ForEachOperationRemoveBreak
)

const (
	MaxInt = int(^uint(0) >> 1)
)
