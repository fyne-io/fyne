package glfw

const defaultTitle = "Fyne Application"

func unknownKey(key string) bool {
	if len(key) != 1 {
		return true
	}

	keyRune := key[0]
	if keyRune >= 'a' && keyRune <= 'z' {
		return false
	}

	// This catches [ \ ] per ASCII table.
	if keyRune >= '[' && keyRune <= ']' {
		return false
	}

	// This catches * + , - . / per ASCII table.
	if keyRune >= '*' && keyRune <= '/' {
		return false
	}

	return keyRune == '\'' || keyRune == '`'
}
