// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package ssg

import (
	"sync"

	"github.com/ALiwoto/ssg/ssg/listUtils"
	"github.com/ALiwoto/ssg/ssg/mapUtils"
	"github.com/ALiwoto/ssg/ssg/rangeValues"
	"github.com/ALiwoto/ssg/ssg/shellUtils"
)

type (
	ExpiringValue[T any]                     = mapUtils.ExpiringValue[T]
	ForEachOperation                         = mapUtils.ForEachOperation
	AdvancedMap[TKey comparable, TValue any] = mapUtils.AdvancedMap[TKey, TValue]
	SafeEMap[TKey comparable, TValue any]    = mapUtils.SafeEMap[TKey, TValue]
	SafeMap[TKey comparable, TValue any]     = mapUtils.SafeMap[TKey, TValue]
)

type (
	ListW[T comparable]       = listUtils.ListW[T]
	GenericList[T comparable] = listUtils.GenericList[T]
)

// the StrongString used in the program for additional usage.
type StrongString struct {
	_value []rune
}

type NumIdGenerator[T rangeValues.Integer] struct {
	current T
	mut     *sync.Mutex
}

type (
	RangeInt     = rangeValues.IntegerRange[int]
	RangeInt32   = rangeValues.IntegerRange[int32]
	RangeInt64   = rangeValues.IntegerRange[int64]
	RangeFloat64 = rangeValues.RangeFloat64

	Int64Container  = rangeValues.IntContainer[int64]
	UInt64Container = rangeValues.IntContainer[uint64]
	Int32Container  = rangeValues.IntContainer[int32]
	UInt32Container = rangeValues.IntContainer[uint32]
	Int16Container  = rangeValues.IntContainer[int16]
	UInt16Container = rangeValues.IntContainer[uint16]
	Int8Container   = rangeValues.IntContainer[int8]
	UInt8Container  = rangeValues.IntContainer[uint8]
)

type ExecuteCommandResult = shellUtils.ExecuteCommandResult

//type safeList[T any] #TODO: implement safe-list

type (
	StringUniqueIdContainer = UniqueIdContainer[string]
	Int64UniqueIdContainer  = UniqueIdContainer[int64]
)

type UniqueIdContainer[T comparable] interface {
	GetUniqueId() T
	SetAsUniqueId(value T)
	HasValidUniqueId() bool
}
