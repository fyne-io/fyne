//go:build !wasm && !test_web_driver && !no_glfw && !mobile

package glfw

import (
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/stretchr/testify/assert"
)

const hugeWindowDim = 1 << 30

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

func TestClampToMonitorSize(t *testing.T) {
	w, h := clampToMonitorSize(100, 80, 1920, 1080)
	assert.Equal(t, 100, w)
	assert.Equal(t, 80, h)

	w, h = clampToMonitorSize(1920, 1080, 1920, 1080)
	assert.Equal(t, 1920, w)
	assert.Equal(t, 1080, h)

	w, h = clampToMonitorSize(4000, 80, 1920, 1080)
	assert.Equal(t, 1920, w)
	assert.Equal(t, 80, h)

	w, h = clampToMonitorSize(100, 4000, 1920, 1080)
	assert.Equal(t, 100, w)
	assert.Equal(t, 1080, h)

	w, h = clampToMonitorSize(4000, 3000, 0, 0)
	assert.Equal(t, 4000, w)
	assert.Equal(t, 3000, h)

	w, h = clampToMonitorSize(4000, 3000, -1, -2)
	assert.Equal(t, 4000, w)
	assert.Equal(t, 3000, h)

	w, h = clampToMonitorSize(0, 0, 1920, 1080)
	assert.Equal(t, 0, w)
	assert.Equal(t, 0, h)
}

func TestWindow_BoundToMonitorSize_OversizedContent(t *testing.T) {
	cases := []struct {
		name string
		min  fyne.Size
	}{
		{"huge width", fyne.NewSize(100000, 40)},
		{"huge height", fyne.NewSize(40, 100000)},
	}

	for _, tt := range cases {
		t.Run(tt.name+" before create", func(t *testing.T) {
			r := canvas.NewRectangle(color.White)
			r.SetMinSize(tt.min)
			var w *window
			var maxW, maxH, nativeW, nativeH int
			runOnMain(func() {
				w = d.CreateWindow("Huge").(*window)
				w.SetContent(r)
				w.create()
				w.view().SetSizeCallback(func(*glfw.Window, int, int) {})
				maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
				nativeW, nativeH = w.view().GetSize()
			})
			if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
				t.Skip("no usable monitor bounds")
			}

			assert.Greater(t, nativeW, 0)
			assert.Greater(t, nativeH, 0)
			assert.LessOrEqual(t, nativeW, maxW)
			assert.LessOrEqual(t, nativeH, maxH)
			assert.LessOrEqual(t, w.shouldWidth, maxW)
			assert.LessOrEqual(t, w.shouldHeight, maxH)
			assert.LessOrEqual(t, w.requestedWidth, maxW)
			assert.LessOrEqual(t, w.requestedHeight, maxH)
		})

		t.Run(tt.name+" after show", func(t *testing.T) {
			w := createWindow("Huge")
			r := canvas.NewRectangle(color.White)
			r.SetMinSize(fyne.NewSize(40, 40))
			w.SetContent(r)
			runOnMain(func() {
				w.Show()
			})

			var maxW, maxH, nativeW, nativeH int
			var expanded bool
			runOnMain(func() {
				maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
				r.SetMinSize(tt.min)
				w.fitContent()
				expanded = w.shouldExpand
				w.view().SetSize(w.shouldWidth, w.shouldHeight)
				nativeW, nativeH = w.view().GetSize()
			})
			if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
				t.Skip("no usable monitor bounds")
			}

			assert.True(t, expanded)
			assert.LessOrEqual(t, w.shouldWidth, maxW)
			assert.LessOrEqual(t, w.shouldHeight, maxH)
			assert.LessOrEqual(t, nativeW, maxW)
			assert.LessOrEqual(t, nativeH, maxH)
		})
	}
}

func TestWindow_BoundToMonitorSize_ResizeFixedAndScale(t *testing.T) {
	t.Run("explicit resize", func(t *testing.T) {
		w := createWindow("Resize")
		r := canvas.NewRectangle(color.White)
		r.SetMinSize(fyne.NewSize(40, 40))
		w.SetContent(r)
		w.Resize(fyne.NewSize(100000, 50))

		var maxW, maxH, nativeW, nativeH int
		runOnMain(func() {
			maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
			nativeW, nativeH = w.view().GetSize()
		})
		if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
			t.Skip("no usable monitor bounds")
		}

		assert.LessOrEqual(t, w.requestedWidth, maxW)
		assert.LessOrEqual(t, w.requestedHeight, maxH)
		assert.LessOrEqual(t, w.width, maxW)
		assert.LessOrEqual(t, w.height, maxH)
		assert.LessOrEqual(t, nativeW, maxW)
		assert.LessOrEqual(t, nativeH, maxH)
	})

	t.Run("fixed size create and fitContent", func(t *testing.T) {
		r := canvas.NewRectangle(color.White)
		r.SetMinSize(fyne.NewSize(100000, 50))
		var w *window
		var maxW, maxH, nativeW, nativeH int
		runOnMain(func() {
			w = d.CreateWindow("Fixed").(*window)
			w.SetFixedSize(true)
			w.SetContent(r)
			w.create()
			w.view().SetSizeCallback(func(*glfw.Window, int, int) {})
			maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
			w.fitContent()
			nativeW, nativeH = w.view().GetSize()
		})
		if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
			t.Skip("no usable monitor bounds")
		}

		assert.LessOrEqual(t, w.shouldWidth, maxW)
		assert.LessOrEqual(t, w.shouldHeight, maxH)
		assert.Equal(t, w.shouldWidth, w.requestedWidth)
		assert.Equal(t, w.shouldHeight, w.requestedHeight)
		assert.LessOrEqual(t, nativeW, maxW)
		assert.LessOrEqual(t, nativeH, maxH)
	})

	t.Run("below and exact bound", func(t *testing.T) {
		w := createWindow("Normal")
		r := canvas.NewRectangle(color.White)
		minSize := fyne.NewSize(100, 80)
		r.SetMinSize(minSize)
		w.SetPadded(true)
		w.SetContent(r)

		var maxW, maxH, minW, minH, expectedW, expectedH, exactW, exactH int
		runOnMain(func() {
			maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
			minW, minH = w.minSizeOnScreen()
			expectedW = scale.ToScreenCoordinate(w.canvas, minSize.Width+theme.Padding()*2)
			expectedH = scale.ToScreenCoordinate(w.canvas, minSize.Height+theme.Padding()*2)
			exactW, exactH = w.boundToMonitorSize(maxW, maxH)
		})
		if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
			t.Skip("no usable monitor bounds")
		}

		assert.Equal(t, expectedW, minW)
		assert.Equal(t, expectedH, minH)
		assert.Less(t, minW, maxW)
		assert.Less(t, minH, maxH)
		assert.Equal(t, maxW, exactW)
		assert.Equal(t, maxH, exactH)
	})

	t.Run("scales use window coordinates", func(t *testing.T) {
		w := createWindow("Scale")
		var maxW, maxH, scaledW, scaledH, highScaleW, highScaleH int
		runOnMain(func() {
			maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
			w.canvas.scale = 2
			w.canvas.detectedScale = 2
			w.canvas.texScale = 2
			scaledW, scaledH = w.boundToMonitorSize(maxW+1, maxH+1)
			w.RescaleContext()
			highScaleW, highScaleH = w.view().GetSize()
		})
		if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
			t.Skip("no usable monitor bounds")
		}

		assert.Equal(t, maxW, scaledW)
		assert.Equal(t, maxH, scaledH)
		assert.LessOrEqual(t, highScaleW, maxW)
		assert.LessOrEqual(t, highScaleH, maxH)
	})

	t.Run("nil viewport uses primary monitor", func(t *testing.T) {
		var boundedW, boundedH int
		runOnMain(func() {
			w := d.CreateWindow("NoView").(*window)
			boundedW, boundedH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
		})
		assert.Greater(t, boundedW, 0)
		assert.Greater(t, boundedH, 0)
		assert.Less(t, boundedW, hugeWindowDim)
		assert.Less(t, boundedH, hugeWindowDim)
	})

	t.Run("long label text preserved", func(t *testing.T) {
		text := strings.Repeat("A", 10000)
		label := widget.NewLabel(text)
		var w *window
		var maxW, maxH, nativeW, nativeH int
		runOnMain(func() {
			w = d.CreateWindow("Label").(*window)
			w.SetContent(label)
			w.create()
			w.view().SetSizeCallback(func(*glfw.Window, int, int) {})
			maxW, maxH = w.boundToMonitorSize(hugeWindowDim, hugeWindowDim)
			nativeW, nativeH = w.view().GetSize()
		})
		assert.Equal(t, text, label.Text)
		if maxW <= 0 || maxH <= 0 || maxW >= hugeWindowDim || maxH >= hugeWindowDim {
			t.Skip("no usable monitor bounds")
		}
		assert.LessOrEqual(t, nativeW, maxW)
		assert.LessOrEqual(t, nativeH, maxH)
	})
}
