package timeUtils

// These are borrowed from unicode/utf8 and strconv and replicate behavior in
// that package, since we can't take a dependency on either.
const (
	lowerHex  = "0123456789abcdef"
	runeSelf  = 0x80
	runeError = '\uFFFD'
)
