// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package tests

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/ALiwoto/ssg/ssg/legacyQs"
)

func TestAppendUnique(t *testing.T) {
	myArray := []int32{1, 2, 3, 4, 5}
	myArray = ssg.AppendUnique(myArray, 1, 10, 22)

	if len(myArray) != 7 {
		t.Error("Expected 7, got", len(myArray))
		return
	}

	anotherArray := []int32{1, 2, 3, 4, 6}
	anotherArray = ssg.AppendUnique(anotherArray, myArray...)

	if len(anotherArray) != 8 {
		t.Error("Expected 8, got", len(anotherArray))
		return
	}
}

func TestClone(t *testing.T) {
	myValue := &struct {
		Name string
	}{}

	add1 := fmt.Sprintf("%p", myValue)
	add2 := fmt.Sprintf("%p", ssg.Clone(myValue))
	add3 := fmt.Sprintf("%p", myValue)
	if add1 == add2 {
		t.Error("Expected different pointers")
		return
	}

	if add1 != add3 {
		t.Error("Expected same pointers")
		return
	}
}

func TestTitleCase(t *testing.T) {
	const (
		str1 = "string1"
		str2 = "thisIsString2"
		str3 = "HelloThere"
	)

	tmp := ssg.Title(str1)
	if tmp != "String1" {
		t.Errorf("Expected %s, got %s", "String1", tmp)
		return
	}

	tmp = ssg.Title(str2)
	if tmp != "ThisIsString2" {
		t.Errorf("Expected %s, got %s", "ThisIsString2", tmp)
		return
	}

	tmp = ssg.Title(str3)
	if tmp != "HelloThere" {
		t.Errorf("Expected %s, got %s", "HelloThere", tmp)
		return
	}
}

func TestStrong(t *testing.T) {
	LogStr("Hi")
	LogInt(5)
	s := legacyQs.Qss("hello!; how are you? () are you okay?")
	if s == nil {
		t.FailNow()
	} else {
		array := s.SplitStr("; ", "() ")
		LogStr("real: " + s.GetValue())
		for i, str := range array {
			LogStr("NOW " + strconv.Itoa(i) + ": " + str.GetValue())
		}
	}
}

func TestToBool(t *testing.T) {
	s := legacyQs.Qss("true")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("false")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = legacyQs.Qss("TRUE")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("FALSE")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = legacyQs.Qss("True")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("False")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = legacyQs.Qss("on")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("off")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = legacyQs.Qss("ON")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("OFF")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = legacyQs.Qss("yes")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("no")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = legacyQs.Qss("YES")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = legacyQs.Qss("NO")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}
}

func TestIntegerHelpers(t *testing.T) {
	const (
		i1       = int64(1)
		i294     = int64(294)
		i356     = int64(356)
		i487     = int64(487)
		i5900    = int64(5900)
		i0x98760 = int64(0x98760)
	)

	s := ssg.ToBase10(i1)
	i := ssg.ToInt64(s)
	if ssg.ToInt64(s) != i1 {
		t.Errorf("Expected %d, got %d", i1, i)
		return
	}

	s = ssg.ToBase10(i294)
	i = ssg.ToInt64(s)
	if ssg.ToInt64(s) != i294 {
		t.Errorf("Expected %d, got %d", i294, i)
		return
	}

	s = ssg.ToBase10(i356)
	i = ssg.ToInt64(s)
	if ssg.ToInt64(s) != i356 {
		t.Errorf("Expected %d, got %d", i356, i)
		return
	}

	s = ssg.ToBase10(i487)
	i = ssg.ToInt64(s)
	if ssg.ToInt64(s) != i487 {
		t.Errorf("Expected %d, got %d", i487, i)
		return
	}

	s = ssg.ToBase10(i5900)
	i = ssg.ToInt64(s)
	if ssg.ToInt64(s) != i5900 {
		t.Errorf("Expected %d, got %d", i5900, i)
		return
	}

	s = ssg.ToBase10(i0x98760)
	i = ssg.ToInt64(s)
	if ssg.ToInt64(s) != i0x98760 {
		t.Errorf("Expected %d, got %d", i0x98760, i)
		return
	}
}

func LogStr(s string) {
	log.Println(s)
}

func LogInt(i int) {
	log.Println(i)
}
