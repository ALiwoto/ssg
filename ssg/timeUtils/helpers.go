package timeUtils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// leadingInt consumes the leading [0-9]* from s.
func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > 1<<63 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

func quote(s string) string {
	buf := make([]byte, 1, len(s)+2) // slice will be at least len(s) + quotes
	buf[0] = '"'
	for i, c := range s {
		if c >= runeSelf || c < ' ' {
			// This means you are asking us to parse a time.Duration or
			// time.Location with unprintable or non-ASCII characters in it.
			// We don't expect to hit this case very often. We could try to
			// reproduce strconv.Quote's behavior with full fidelity but
			// given how rarely we expect to hit these edge cases, speed and
			// conciseness are better.
			var width int
			if c == runeError {
				width = 1
				if i+2 < len(s) && s[i:i+3] == string(runeError) {
					width = 3
				}
			} else {
				width = len(string(c))
			}
			for j := 0; j < width; j++ {
				buf = append(buf, `\x`...)
				buf = append(buf, lowerHex[s[i+j]>>4])
				buf = append(buf, lowerHex[s[i+j]&0xF])
			}
		} else {
			if c == '"' || c == '\\' {
				buf = append(buf, '\\')
			}
			buf = append(buf, string(c)...)
		}
	}
	buf = append(buf, '"')
	return string(buf)
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// If you want to define custom units, you can use the RegisterUnit function.
// This code is a modified version of time.ParseDuration from the Go standard library.
// The original code can be found at https://golang.org/src/time/format.go
// or https://github.com/golang/go/blob/master/src/time/format.go
func ParseDuration(s string) (Duration, error) {
	s = strings.ReplaceAll(s, " ", "")
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, errors.New("time: invalid duration " + quote(orig))
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, errors.New("time: invalid duration " + quote(orig))
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New("time: invalid duration " + quote(orig))
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, errors.New("time: invalid duration " + quote(orig))
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, errors.New("time: missing unit in duration " + quote(orig))
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, errors.New("time: unknown unit " + quote(u) + " in duration " + quote(orig))
		}
		if v > 1<<63/unit {
			// overflow
			return 0, errors.New("time: invalid duration " + quote(orig))
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return 0, errors.New("time: invalid duration " + quote(orig))
			}
		}
		d += v
		if d > 1<<63 {
			return 0, errors.New("time: invalid duration " + quote(orig))
		}
	}
	if neg {
		return -Duration(d), nil
	}
	if d > 1<<63-1 {
		return 0, errors.New("time: invalid duration " + quote(orig))
	}
	return Duration(d), nil
}

// ParseDurationWithDefault parses a duration string and returns the duration.
// If the string is invalid, it will return the default value.
// If the string is a number without a unit, it will multiply the number by the default value.
func ParseDurationWithDefault(s string, defaultValue Duration) Duration {
	if s == "" {
		return defaultValue
	}

	duration, err := ParseDuration(s)
	if err != nil {
		if strings.Contains(err.Error(), "missing unit in duration") {
			converted, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				return Duration(converted) * defaultValue
			}
		}
		return defaultValue
	}
	return duration
}

// RegisterUnit registers a new unit with the given value.
// WARNING: This function is not thread-safe, it's recommended to use
// it only during the initialization phase of the program.
// If the unit is already registered, the value will be updated.
// If the unit contains whitespaces, they will be removed.
// If the unit is an empty string, nothing will happen.
func RegisterUnit(unit string, value uint64) {
	unit = strings.ReplaceAll(unit, " ", "")
	if unit == "" {
		return
	}
	unitMap[unit] = value
}

// GetPrettyTimeDuration returns a human-readable string representing the time duration.
// If shorten is true, the output will be formatted as "1d2h3m4s" instead of
// "1 day 2 hours 3 minutes 4 seconds".
func GetPrettyTimeDuration(d time.Duration, shorten bool) string {
	var result string
	totalSeconds := int(d.Seconds())

	year := totalSeconds / (60 * 60 * 24 * 365)
	totalSeconds -= year * (60 * 60 * 24 * 365)

	month := totalSeconds / (60 * 60 * 24 * 30)
	totalSeconds -= month * (60 * 60 * 24 * 30)

	day := totalSeconds / (60 * 60 * 24)
	totalSeconds -= day * (60 * 60 * 24)

	hour := totalSeconds / (60 * 60)
	totalSeconds -= hour * (60 * 60)

	minute := totalSeconds / 60
	totalSeconds -= minute * 60

	seconds := totalSeconds

	yBool := year > 0
	mBool := month > 0 || yBool
	shorten = !mBool && shorten
	dBool := day > 0 || mBool
	hBool := hour > 0 || dBool
	if yBool {
		result += strconv.Itoa(year) + " year"
		if year > 1 {
			result += "s"
		}
		result += " "
	}
	if mBool {
		result += " " + strconv.Itoa(month) + " month"
		if month > 1 {
			result += "s"
		}
		result += " "
	}
	if dBool {
		result += strconv.Itoa(day)
		if shorten {
			result += "d"
		} else {
			result += " day"
			if day > 1 {
				result += "s"
			}
		}
		result += " "
	}
	if hBool {
		result += strconv.Itoa(hour)
		if shorten {
			result += "h"
		} else {
			result += " hour"
			if hour > 1 {
				result += "s"
			}
		}
		result += " "
	}
	result += strconv.Itoa(minute)
	if shorten {
		result += "m"
	} else {
		result += " minute"
		if minute > 1 {
			result += "s"
		}
	}

	result += " " + strconv.Itoa(seconds)
	if shorten {
		result += "s"
	} else {
		result += " second"
		if seconds > 1 {
			result += "s"
		}
	}
	return result
}
