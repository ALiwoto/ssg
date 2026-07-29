package strongParser

import (
	"reflect"
	"regexp"
)

var (
	sectionHeader      = regexp.MustCompile(`\[([^]]+)\]`)
	keyValue           = regexp.MustCompile(`([^:=\s][^:=]*)\s*(?P<vi>[:=])\s*(.*)$`)
	DefaultMainSection = "main"
)

var invalidReflectValue = reflect.ValueOf(nil)
