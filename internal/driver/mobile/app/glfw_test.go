//go:build (linux && !android) || freebsd || openbsd

package app

import (
	"testing"

	"fyne.io/fyne/v2/internal/driver/mobile/event/key"

	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/stretchr/testify/assert"
)

func TestGLFWKeyToCode_Letters(t *testing.T) {
	// GLFW reports the uppercase code point for a letter, the case that is
	// actually typed comes from the character callback.
	for i := 0; i < 26; i++ {
		assert.Equal(t, key.CodeA+key.Code(i), glfwKeyToCode(glfw.Key('A'+i)), "letter %c", 'A'+i)
	}
}

func TestGLFWKeyToCode_Digits(t *testing.T) {
	// the 0 key is out of order in the USB HID spec that the key package uses
	for r := '0'; r <= '9'; r++ {
		assert.Equal(t, getCodeFromRune(r), glfwKeyToCode(glfw.Key(r)), "digit %c", r)
	}
	assert.Equal(t, key.Code0, glfwKeyToCode(glfw.Key0))
	assert.Equal(t, key.Code9, glfwKeyToCode(glfw.Key9))
}

func TestGLFWKeyToCode_Punctuation(t *testing.T) {
	for k, code := range map[glfw.Key]key.Code{
		glfw.KeySpace:        key.CodeSpacebar,
		glfw.KeyApostrophe:   key.CodeApostrophe,
		glfw.KeyComma:        key.CodeComma,
		glfw.KeyEqual:        key.CodeEqualSign,
		glfw.KeyGraveAccent:  key.CodeGraveAccent,
		glfw.KeyLeftBracket:  key.CodeLeftSquareBracket,
		glfw.KeyMinus:        key.CodeHyphenMinus,
		glfw.KeyPeriod:       key.CodeFullStop,
		glfw.KeyRightBracket: key.CodeRightSquareBracket,
		glfw.KeySemicolon:    key.CodeSemicolon,
		glfw.KeySlash:        key.CodeSlash,
		glfw.KeyBackslash:    key.CodeBackslash,
	} {
		assert.Equal(t, code, glfwKeyToCode(k), "key %d", k)
	}
}

func TestGLFWKeyToCode_NonPrintable(t *testing.T) {
	for k, code := range map[glfw.Key]key.Code{
		glfw.KeyBackspace:  key.CodeDeleteBackspace,
		glfw.KeyDelete:     key.CodeDeleteForward,
		glfw.KeyDown:       key.CodeDownArrow,
		glfw.KeyEnd:        key.CodeEnd,
		glfw.KeyEnter:      key.CodeReturnEnter,
		glfw.KeyEscape:     key.CodeEscape,
		glfw.KeyHome:       key.CodeHome,
		glfw.KeyInsert:     key.CodeInsert,
		glfw.KeyLeft:       key.CodeLeftArrow,
		glfw.KeyPageDown:   key.CodePageDown,
		glfw.KeyPageUp:     key.CodePageUp,
		glfw.KeyRight:      key.CodeRightArrow,
		glfw.KeyTab:        key.CodeTab,
		glfw.KeyUp:         key.CodeUpArrow,
		glfw.KeyKPEnter:    key.CodeKeypadEnter,
		glfw.KeyKPDecimal:  key.CodeKeypadFullStop,
		glfw.KeyKPSubtract: key.CodeKeypadHyphenMinus,
		glfw.KeyKP0:        key.CodeKeypad0,
		glfw.KeyKP9:        key.CodeKeypad9,
		glfw.KeyMenu:       key.CodeCompose,
	} {
		assert.Equal(t, code, glfwKeyToCode(k), "key %d", k)
	}
}

func TestGLFWKeyToCode_Modifiers(t *testing.T) {
	for k, code := range map[glfw.Key]key.Code{
		glfw.KeyLeftAlt:      key.CodeLeftAlt,
		glfw.KeyLeftControl:  key.CodeLeftControl,
		glfw.KeyLeftShift:    key.CodeLeftShift,
		glfw.KeyLeftSuper:    key.CodeLeftGUI,
		glfw.KeyRightAlt:     key.CodeRightAlt,
		glfw.KeyRightControl: key.CodeRightControl,
		glfw.KeyRightShift:   key.CodeRightShift,
		glfw.KeyRightSuper:   key.CodeRightGUI,
	} {
		assert.Equal(t, code, glfwKeyToCode(k), "key %d", k)
	}
}

func TestGLFWKeyToCode_FunctionKeys(t *testing.T) {
	for i := 1; i <= 12; i++ {
		assert.Equal(t, key.CodeF1+key.Code(i-1), glfwKeyToCode(glfw.KeyF1+glfw.Key(i-1)), "F%d", i)
	}
}

func TestGLFWKeyToCode_Unknown(t *testing.T) {
	// keys that the key package has no code for must not fall through to the
	// printable range and invent one
	for _, k := range []glfw.Key{glfw.KeyUnknown, glfw.KeyF25, glfw.KeyCapsLock, glfw.KeyNumLock, glfw.KeyPause} {
		assert.Equal(t, key.CodeUnknown, glfwKeyToCode(k), "key %d", k)
	}
}

func TestGLFWModifiers(t *testing.T) {
	assert.Equal(t, key.Modifiers(0), modifiers(0))
	assert.Equal(t, key.ModShift, modifiers(glfw.ModShift))
	assert.Equal(t, key.ModControl, modifiers(glfw.ModControl))
	assert.Equal(t, key.ModAlt, modifiers(glfw.ModAlt))
	assert.Equal(t, key.ModMeta, modifiers(glfw.ModSuper))
	assert.Equal(t, key.ModShift|key.ModControl|key.ModAlt|key.ModMeta,
		modifiers(glfw.ModShift|glfw.ModControl|glfw.ModAlt|glfw.ModSuper))
	// the lock keys have no equivalent in the mobile key package
	assert.Equal(t, key.Modifiers(0), modifiers(glfw.ModCapsLock|glfw.ModNumLock))
}

func TestModifiersCorrected(t *testing.T) {
	// X11 does not report a modifier key that is being pressed or released in
	// the modifier mask, so it has to be applied by hand.
	// See https://github.com/glfw/glfw/issues/1630
	assert.Equal(t, key.ModShift, modifiersCorrected(0, glfw.KeyLeftShift, glfw.Press))
	assert.Equal(t, key.Modifiers(0), modifiersCorrected(glfw.ModShift, glfw.KeyLeftShift, glfw.Release))

	assert.Equal(t, key.ModControl, modifiersCorrected(0, glfw.KeyRightControl, glfw.Press))
	assert.Equal(t, key.Modifiers(0), modifiersCorrected(glfw.ModControl, glfw.KeyRightControl, glfw.Release))

	assert.Equal(t, key.ModAlt, modifiersCorrected(0, glfw.KeyLeftAlt, glfw.Press))
	assert.Equal(t, key.ModMeta, modifiersCorrected(0, glfw.KeyLeftSuper, glfw.Press))

	// a repeat keeps the modifier that the initial press added
	assert.Equal(t, key.ModShift, modifiersCorrected(glfw.ModShift, glfw.KeyLeftShift, glfw.Repeat))

	// a key that is not a modifier is left alone
	assert.Equal(t, key.Modifiers(0), modifiersCorrected(0, glfw.KeyA, glfw.Press))
	assert.Equal(t, key.Modifiers(0), modifiersCorrected(0, glfw.KeyA, glfw.Release))
}
