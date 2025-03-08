package timeUtils

import "time"

var unitMap = map[string]uint64{
	"ns": uint64(time.Nanosecond),
	"us": uint64(time.Microsecond),
	"µs": uint64(time.Microsecond), // U+00B5 = micro symbol
	"μs": uint64(time.Microsecond), // U+03BC = Greek letter mu
	"ms": uint64(time.Millisecond),
	"s":  uint64(time.Second),
	"m":  uint64(time.Minute),
	"h":  uint64(time.Hour),
	"d":  uint64(time.Hour * 24),
	"w":  uint64(time.Hour * 24 * 7),
	"mo": uint64(time.Hour * 24 * 30),
	"y":  uint64(time.Hour * 24 * 365),
}
