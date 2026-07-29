package legacyQs

import "github.com/ALiwoto/ssg/ssg/internal"

const (
	LineChar         = '-' // line : '-'
	EqualChar        = '=' // equal: '='
	SpaceChar        = ' ' // space: ' '
	DPointChar       = ':' // double point: ':'
	BracketOpenChar  = '[' // bracket open: '['
	BracketCloseChar = ']' // bracket close: ']'
	ComaChar         = ',' // coma: ','
)

const (
	JA_FLAG       = "〰〰"
	JA_STR        = "❞" // start character (") for string in japanese.
	JA_EQUALITY   = "＝" // equal character (＝) for string in japanese.
	JA_DDOT       = "：" // equal character (＝) for string in japanese.
	BACK_STR      = "\\\""
	BACK_FLAG     = "\\--"
	BACK_EQUALITY = "\\="
	BACK_DDOT     = "\\:"
)

// the prefix values for commands.
const (
	COMMAND_PREFIX1 = "!"
	COMMAND_PREFIX2 = "/"
	SUDO_PREFIX1    = ">"
	FLAG_PREFIX     = "--"
)

// the base constant values.
const (
	BaseIndex      = 0  // number 0
	BaseOneIndex   = 1  // number 1
	BaseTwoIndex   = 2  // number 2
	BaseThreeIndex = 3  // number 2
	Base4Bit       = 4  // number 8
	Base8Bit       = 8  // number 8
	Base16Bit      = 16 // number 16
	Base32Bit      = 32 // number 32
	Base64Bit      = 64 // number 64
	BaseTimeOut    = 40 // 40 seconds
	BaseTen        = 10 // 10 seconds
)

// additional constants which are not actually used in
// this package, but may be useful in another packages.
const (
	BaseIndexStr    = "0" // number 0
	BaseOneIndexStr = "1" // number 1
	DotStr          = "." // dot : .
	LineStr         = "-" // line : -
	EMPTY           = ""  //an empty string.
	UNDER           = "_" // an underscore : _
	STR_SIGN        = `"` // the string sign : "
	CHAR_STR        = '"' // the string sign : '"'
)

const (
	LineEscape      = "\n"
	R_ESCAPE        = "\r"
	Coma            = ","
	ParaOpen        = "("
	ParaClose       = ")"
	NullStr         = "null"
	DoubleQ         = "\""
	SingleQ         = "'"
	DoubleQJ        = "”"
	BracketOpen     = "["
	BracketClose    = "]"
	Star            = "*"
	BackSlash       = "\\"
	DoubleBackSlash = "\\\\"
	Point           = "."
	AutoStr         = "auto"
	AtSign          = "@"
	EqualStr        = "="
	DdotSign        = ":"
	Yes             = "Yes"
	No              = "No"
	TrueStr         = "True"
	FalseStr        = "False"
	OnStr           = "On"
	OffStr          = "Off"
	LowerYes        = "yes"
	LowerNo         = "no"
	LowerTrueStr    = "true"
	LowerFalseStr   = "false"
	LowerOnStr      = "on"
	LowerOffStr     = "off"
	OrRegexp        = internal.OrRegexp // the or string sign: "|"
)

const (
	LIST_INDEX_NOTFOUND = -1
)
