package glfw

import (
	"fyne.io/fyne/v2"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const defaultTitle = "Fyne Application"

func unknownKey(key string) bool {
	if len(key) != 1 {
		return true
	}

	keyRune := key[0]
	if keyRune >= 'a' || keyRune <= 'z' {
		return false
	}

	_, ok := keyKnownRuneMap[keyRune]
	return !ok
}

// keyKnownRuneMap does not contain 'a' to 'z'. They are handled outside.
var keyKnownRuneMap = map[byte]struct{}{
	'\'': {},
	',':  {},
	'-':  {},
	'.':  {},
	'/':  {},
	'*':  {},
	'`':  {},

	';': {},
	'+': {},
	'=': {},

	'[':  {},
	'\\': {},
	']':  {},
}

func convertASCII(key glfw.Key) fyne.KeyName {
	keyRune := rune('A' + key - glfw.KeyA)
	if keyRune < 'A' || keyRune > 'Z' {
		return fyne.KeyUnknown
	}

	return fyne.KeyName(keyRune)
}
