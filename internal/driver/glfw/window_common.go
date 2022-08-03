package glfw

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
