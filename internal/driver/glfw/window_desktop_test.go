//go:build !wasm && !test_web_driver

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/stretchr/testify/assert"
)

var keyCodeMap = map[glfw.Key]fyne.KeyName{
	// non-printable
	glfw.KeyEscape:    fyne.KeyEscape,
	glfw.KeyEnter:     fyne.KeyReturn,
	glfw.KeyTab:       fyne.KeyTab,
	glfw.KeyBackspace: fyne.KeyBackspace,
	glfw.KeyInsert:    fyne.KeyInsert,
	glfw.KeyDelete:    fyne.KeyDelete,
	glfw.KeyRight:     fyne.KeyRight,
	glfw.KeyLeft:      fyne.KeyLeft,
	glfw.KeyDown:      fyne.KeyDown,
	glfw.KeyUp:        fyne.KeyUp,
	glfw.KeyPageUp:    fyne.KeyPageUp,
	glfw.KeyPageDown:  fyne.KeyPageDown,
	glfw.KeyHome:      fyne.KeyHome,
	glfw.KeyEnd:       fyne.KeyEnd,

	glfw.KeySpace:   fyne.KeySpace,
	glfw.KeyKPEnter: fyne.KeyEnter,

	// functions
	glfw.KeyF1:  fyne.KeyF1,
	glfw.KeyF2:  fyne.KeyF2,
	glfw.KeyF3:  fyne.KeyF3,
	glfw.KeyF4:  fyne.KeyF4,
	glfw.KeyF5:  fyne.KeyF5,
	glfw.KeyF6:  fyne.KeyF6,
	glfw.KeyF7:  fyne.KeyF7,
	glfw.KeyF8:  fyne.KeyF8,
	glfw.KeyF9:  fyne.KeyF9,
	glfw.KeyF10: fyne.KeyF10,
	glfw.KeyF11: fyne.KeyF11,
	glfw.KeyF12: fyne.KeyF12,

	// numbers - lookup by code to avoid AZERTY using the symbol name instead of number
	glfw.Key0:   fyne.Key0,
	glfw.KeyKP0: fyne.Key0,
	glfw.Key1:   fyne.Key1,
	glfw.KeyKP1: fyne.Key1,
	glfw.Key2:   fyne.Key2,
	glfw.KeyKP2: fyne.Key2,
	glfw.Key3:   fyne.Key3,
	glfw.KeyKP3: fyne.Key3,
	glfw.Key4:   fyne.Key4,
	glfw.KeyKP4: fyne.Key4,
	glfw.Key5:   fyne.Key5,
	glfw.KeyKP5: fyne.Key5,
	glfw.Key6:   fyne.Key6,
	glfw.KeyKP6: fyne.Key6,
	glfw.Key7:   fyne.Key7,
	glfw.KeyKP7: fyne.Key7,
	glfw.Key8:   fyne.Key8,
	glfw.KeyKP8: fyne.Key8,
	glfw.Key9:   fyne.Key9,
	glfw.KeyKP9: fyne.Key9,

	// desktop
	glfw.KeyLeftShift:    desktop.KeyShiftLeft,
	glfw.KeyRightShift:   desktop.KeyShiftRight,
	glfw.KeyLeftControl:  desktop.KeyControlLeft,
	glfw.KeyRightControl: desktop.KeyControlRight,
	glfw.KeyLeftAlt:      desktop.KeyAltLeft,
	glfw.KeyRightAlt:     desktop.KeyAltRight,
	glfw.KeyLeftSuper:    desktop.KeySuperLeft,
	glfw.KeyRightSuper:   desktop.KeySuperRight,
	glfw.KeyMenu:         desktop.KeyMenu,
	glfw.KeyPrintScreen:  desktop.KeyPrintScreen,
	glfw.KeyCapsLock:     desktop.KeyCapsLock,
}

func TestGlfwKeyToKeyName(t *testing.T) {
	for key, value := range keyCodeMap {
		translated := glfwKeyToKeyName(key)
		assert.Equal(t, value, translated)
	}

	invalid := glfwKeyToKeyName(glfw.Key(-1))
	assert.Equal(t, fyne.KeyUnknown, invalid)
}

func TestConvertASCII(t *testing.T) {
	for i := 0; i <= 'Z'-'A'; i++ {
		translated := convertASCII(glfw.KeyA + glfw.Key(i))
		expected := fyne.KeyName(rune(fyne.KeyA[0] + byte(i)))
		assert.Equal(t, expected, translated)
	}

	invalid := convertASCII(glfw.Key(-1))
	assert.Equal(t, fyne.KeyUnknown, invalid)
}

func TestDesktopModifierCorrected(t *testing.T) {
	for _, tt := range []struct {
		name     string
		key      glfw.Key
		modifier glfw.ModifierKey
		want     fyne.KeyModifier
	}{
		{"left control", glfw.KeyLeftControl, glfw.ModControl, fyne.KeyModifierControl},
		{"right control", glfw.KeyRightControl, glfw.ModControl, fyne.KeyModifierControl},
		{"left shift", glfw.KeyLeftShift, glfw.ModShift, fyne.KeyModifierShift},
		{"right shift", glfw.KeyRightShift, glfw.ModShift, fyne.KeyModifierShift},
		{"left alt", glfw.KeyLeftAlt, glfw.ModAlt, fyne.KeyModifierAlt},
		{"right alt", glfw.KeyRightAlt, glfw.ModAlt, fyne.KeyModifierAlt},
		{"left super", glfw.KeyLeftSuper, glfw.ModSuper, fyne.KeyModifierSuper},
		{"right super", glfw.KeyRightSuper, glfw.ModSuper, fyne.KeyModifierSuper},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, desktopModifierCorrected(0, tt.key, glfw.Press))
			assert.Equal(t, tt.want, desktopModifierCorrected(tt.modifier, tt.key, glfw.Press))
			assert.Equal(t, tt.want, desktopModifierCorrected(tt.modifier, tt.key, glfw.Repeat))
			assert.Equal(t, fyne.KeyModifier(0), desktopModifierCorrected(0, tt.key, glfw.Release))
			assert.Equal(t, fyne.KeyModifier(0), desktopModifierCorrected(tt.modifier, tt.key, glfw.Release))

			allModifiers := glfw.ModControl | glfw.ModShift | glfw.ModAlt | glfw.ModSuper
			allWant := fyne.KeyModifierControl | fyne.KeyModifierShift | fyne.KeyModifierAlt | fyne.KeyModifierSuper
			assert.Equal(t, allWant, desktopModifierCorrected(allModifiers&^tt.modifier, tt.key, glfw.Press))
			assert.Equal(t, allWant, desktopModifierCorrected(allModifiers, tt.key, glfw.Repeat))
			assert.Equal(t, allWant&^tt.want, desktopModifierCorrected(allModifiers, tt.key, glfw.Release))
		})
	}
}

func TestDesktopModifierCorrected_NonModifierKey(t *testing.T) {
	for _, action := range []glfw.Action{glfw.Press, glfw.Repeat, glfw.Release} {
		assert.Equal(t, fyne.KeyModifierControl|fyne.KeyModifierShift,
			desktopModifierCorrected(glfw.ModControl|glfw.ModShift, glfw.KeySpace, action))
	}
}

var keyNameMapSpecialCharacters = map[string]fyne.KeyName{
	"'": fyne.KeyApostrophe,
	",": fyne.KeyComma,
	"-": fyne.KeyMinus,
	".": fyne.KeyPeriod,
	"/": fyne.KeySlash,
	"*": fyne.KeyAsterisk,
	"`": fyne.KeyBackTick,

	";": fyne.KeySemicolon,
	"+": fyne.KeyPlus,
	"=": fyne.KeyEqual,

	"[":  fyne.KeyLeftBracket,
	"\\": fyne.KeyBackslash,
	"]":  fyne.KeyRightBracket,
}

func TestKeyCodeToKeyName(t *testing.T) {
	for key, value := range keyNameMapSpecialCharacters {
		translated := keyCodeToKeyName(key)
		assert.Equal(t, value, translated)
	}

	for i := rune(0); i <= 'z'-'a'; i++ {
		translated := keyCodeToKeyName(string('a' + i))
		expected := fyne.KeyName(rune(fyne.KeyA[0]) + i)
		assert.Equal(t, expected, translated)
	}

	invalid := keyCodeToKeyName("@")
	assert.Equal(t, fyne.KeyUnknown, invalid)

	invalid = keyCodeToKeyName("invalid")
	assert.Equal(t, fyne.KeyUnknown, invalid)
}
