package timeUtils

import "errors"

var (
	errLeadingInt = errors.New("time: bad [0-9]*") // never printed
)
