package app

import "C"
import "github.com/fyne-io/mobile/event/key"

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
			Rune: r,
			Code: getCodeFromRune(r),
			Direction: key.DirPress,
		}
		theApp.eventsIn <- k

		k.Direction = key.DirRelease
		theApp.eventsIn <- k
	}
}

//export keyboardDelete
func keyboardDelete() {
	theApp.eventsIn <- key.Event{
		Code: key.CodeDeleteBackspace,
		Direction: key.DirPress,
	}
	theApp.eventsIn <- key.Event{
		Code: key.CodeDeleteBackspace,
		Direction: key.DirRelease,
	}
}

func getCodeFromRune(r rune) key.Code {
	switch r {
	case '0':
		return key.Code0
	case '1':
		return key.Code1
	case '2':
		return key.Code2
	case '3':
		return key.Code3
	case '4':
		return key.Code4
	case '5':
		return key.Code5
	case '6':
		return key.Code6
	case '7':
		return key.Code7
	case '8':
		return key.Code8
	case '9':
		return key.Code9
	case 'a', 'A':
		return key.CodeA
	case 'b', 'B':
		return key.CodeB
	case 'c', 'C':
		return key.CodeC
	case 'd', 'D':
		return key.CodeD
	case 'e', 'E':
		return key.CodeE
	case 'f', 'F':
		return key.CodeF
	case 'g', 'G':
		return key.CodeG
	case 'h', 'H':
		return key.CodeH
	case 'i', 'I':
		return key.CodeI
	case 'j', 'J':
		return key.CodeJ
	case 'k', 'K':
		return key.CodeK
	case 'l', 'L':
		return key.CodeL
	case 'm', 'M':
		return key.CodeM
	case 'n', 'N':
		return key.CodeN
	case 'o', 'O':
		return key.CodeO
	case 'p', 'P':
		return key.CodeP
	case 'q', 'Q':
		return key.CodeQ
	case 'r', 'R':
		return key.CodeR
	case 's', 'S':
		return key.CodeS
	case 't', 'T':
		return key.CodeT
	case 'u', 'U':
		return key.CodeU
	case 'v', 'V':
		return key.CodeV
	case 'w', 'W':
		return key.CodeW
	case 'x', 'X':
		return key.CodeX
	case 'y', 'Y':
		return key.CodeY
	case 'z', 'Z':
		return key.CodeZ
	case ',':
		return key.CodeComma
	case '.':
		return key.CodeFullStop
	case ' ':
		return key.CodeSpacebar
	case '\n':
		return key.CodeReturnEnter
	case '`':
		return key.CodeGraveAccent
	case '-':
		return key.CodeHyphenMinus
	case '=':
		return key.CodeEqualSign
	case '[':
		return key.CodeLeftSquareBracket
	case ']':
		return key.CodeRightSquareBracket
	case '\\':
		return key.CodeBackslash
	case ';':
		return key.CodeSemicolon
	case '\'':
		return key.CodeApostrophe
	case '/':
		return key.CodeSlash
	}
	return key.CodeUnknown
}
