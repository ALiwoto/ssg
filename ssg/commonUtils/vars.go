package commonUtils

// BoolMapping is a map of strings to bool.
// WARNING: This is not a safe map, it should remain read-only.
var BoolMapping = map[string]bool{
	"1":        true,
	"true":     true,
	"on":       true,
	"yes":      true,
	"y":        true,
	"enable":   true,
	"enabled":  true,
	"0":        false,
	"false":    false,
	"off":      false,
	"no":       false,
	"n":        false,
	"disable":  false,
	"disabled": false,
}
