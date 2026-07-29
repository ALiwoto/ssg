package legacyQs

import "github.com/ALiwoto/ssg/ssg/listUtils"

// StrongString is the legacy type kept here for legacy purposes.
type StrongString struct {
	_value []rune
}

type QString interface {
	listUtils.ListLike

	GetValue() string
	GetIndexV(int) rune
	IsEqual(QString) bool
	Split(...QString) []QString
	SplitN(int, ...QString) []QString
	SplitFirst(qs ...QString) []QString
	SplitStr(...string) []QString
	SplitStrN(int, ...string) []QString
	SplitStrFirst(...string) []QString
	Contains(...QString) bool
	ContainsStr(...string) bool
	ContainsAll(...QString) bool
	ContainsStrAll(...string) bool
	TrimPrefix(...QString) QString
	TrimPrefixStr(...string) QString
	TrimSuffix(...QString) QString
	TrimSuffixStr(...string) QString
	Trim(qs ...QString) QString
	TrimStr(...string) QString
	Replace(qs, newS QString) QString
	ReplaceStr(qs, newS string) QString
	LockSpecial()
	UnlockSpecial()
	ToBool() bool
}
