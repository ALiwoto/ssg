package strongParser

import "reflect"

type rValue = reflect.Value

// InvalidParseError describes an invalid argument passed to ParseConfig.
// (The argument to ParseConfig must be a non-nil pointer.)
type InvalidParseError struct {
	Type reflect.Type
}

type Section struct {
	Name    string
	options Dict
	lookup  Dict
}

// Dict is a simple string->string map.
type Dict map[string]string

// Config represents a Python style configuration file.
type Config map[string]*Section

// ConfigParser ties together a Config and default values for use in
// interpolated configuration values.
type ConfigParser struct {
	config   Config
	defaults *Section
}

type ChainMap struct {
	maps []Dict
}
