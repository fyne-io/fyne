package app

import "C"
import "fyne.io/fyne/v2/internal/driver/mobile/event/key"

// KeyboardType represents the type of a keyboard
type KeyboardType int32

const (
	// DefaultKeyboard is the keyboard with default input style and "return" return key
	DefaultKeyboard KeyboardType = iota
	// SingleLineKeyboard is the keyboard with default input style and "Done" return key
	SingleLineKeyboard
	// NumberKeyboard is the keyboard with number input style and "Done" return key
	NumberKeyboard
)

//export keyboardTyped
func keyboardTyped(str *C.char) {
	for _, r := range C.GoString(str) {
		k := key.Event{
			Rune:      r,
			Code:      getCodeFromRune(r),
			Direction: key.DirPress,
		}
		theApp.events.In() <- k

		k.Direction = key.DirRelease
		theApp.events.In() <- k
	}
}

//export keyboardDelete
func keyboardDelete() {
	theApp.events.In() <- key.Event{
		Code:      key.CodeDeleteBackspace,
		Direction: key.DirPress,
	}
	theApp.events.In() <- key.Event{
		Code:      key.CodeDeleteBackspace,
		Direction: key.DirRelease,
	}
}

var codeRune = map[rune]key.Code{
	'0':  key.Code0,
	'1':  key.Code1,
	'2':  key.Code2,
	'3':  key.Code3,
	'4':  key.Code4,
	'5':  key.Code5,
	'6':  key.Code6,
	'7':  key.Code7,
	'8':  key.Code8,
	'9':  key.Code9,
	'a':  key.CodeA,
	'b':  key.CodeB,
	'c':  key.CodeC,
	'd':  key.CodeD,
	'e':  key.CodeE,
	'f':  key.CodeF,
	'g':  key.CodeG,
	'h':  key.CodeH,
	'i':  key.CodeI,
	'j':  key.CodeJ,
	'k':  key.CodeK,
	'l':  key.CodeL,
	'm':  key.CodeM,
	'n':  key.CodeN,
	'o':  key.CodeO,
	'p':  key.CodeP,
	'q':  key.CodeQ,
	'r':  key.CodeR,
	's':  key.CodeS,
	't':  key.CodeT,
	'u':  key.CodeU,
	'v':  key.CodeV,
	'w':  key.CodeW,
	'x':  key.CodeX,
	'y':  key.CodeY,
	'z':  key.CodeZ,
	'A':  key.CodeA,
	'B':  key.CodeB,
	'C':  key.CodeC,
	'D':  key.CodeD,
	'E':  key.CodeE,
	'F':  key.CodeF,
	'G':  key.CodeG,
	'H':  key.CodeH,
	'I':  key.CodeI,
	'J':  key.CodeJ,
	'K':  key.CodeK,
	'L':  key.CodeL,
	'M':  key.CodeM,
	'N':  key.CodeN,
	'O':  key.CodeO,
	'P':  key.CodeP,
	'Q':  key.CodeQ,
	'R':  key.CodeR,
	'S':  key.CodeS,
	'T':  key.CodeT,
	'U':  key.CodeU,
	'V':  key.CodeV,
	'W':  key.CodeW,
	'X':  key.CodeX,
	'Y':  key.CodeY,
	'Z':  key.CodeZ,
	',':  key.CodeComma,
	'.':  key.CodeFullStop,
	' ':  key.CodeSpacebar,
	'\n': key.CodeReturnEnter,
	'`':  key.CodeGraveAccent,
	'-':  key.CodeHyphenMinus,
	'=':  key.CodeEqualSign,
	'[':  key.CodeLeftSquareBracket,
	']':  key.CodeRightSquareBracket,
	'\\': key.CodeBackslash,
	';':  key.CodeSemicolon,
	'\'': key.CodeApostrophe,
	'/':  key.CodeSlash,
}

func getCodeFromRune(r rune) key.Code {
	if code, ok := codeRune[r]; ok {
		return code
	}
	return key.CodeUnknown
}
