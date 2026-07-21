package dialog

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynecolor "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestColorDialog_Resize(t *testing.T) {
	test.NewTempApp(t)

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(1000, 800))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Resize(fyne.NewSize(800, 600))
	d.Show()

	size := d.content.Size()

	d.Resize(fyne.NewSize(900, 600))
	newSize := d.content.Size()

	assert.Equal(t, float32(100), newSize.Width-size.Width)
}

func TestColorDialog_Theme(t *testing.T) {
	test.NewTempApp(t)

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(1000, 800))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Advanced = true
	d.Show()

	test.AssertRendersToImage(t, "color/dialog_theme_default.png", w.Canvas())

	test.ApplyTheme(t, test.NewTheme())
	d.Resize(d.MinSize())
	test.AssertRendersToImage(t, "color/dialog_theme_ugly.png", w.Canvas())

	d.advanced.Open(0)

	test.ApplyTheme(t, test.Theme())
	d.Resize(d.MinSize())

	test.AssertRendersToImage(t, "color/dialog_expanded_theme_default.png", w.Canvas())

	test.ApplyTheme(t, test.NewTheme())
	d.Resize(d.MinSize())
	test.AssertRendersToImage(t, "color/dialog_expanded_theme_ugly.png", w.Canvas())
}

func TestColorDialog_Recents(t *testing.T) {
	a := test.NewTempApp(t)
	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#2196f3,#4caf50,#f44336")

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(800, 600))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Advanced = true
	d.Refresh()
	d.Show()

	test.AssertRendersToImage(t, "color/dialog_recents_theme_default.png", w.Canvas())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertRendersToImage(t, "color/dialog_recents_theme_ugly.png", w.Canvas())
}

func TestColorDialog_SetColor(t *testing.T) {
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(800, 600))

	col := color.RGBA{R: 70, G: 210, B: 200, A: 255}

	d := NewColorPicker("pick colour", "select colour", func(c color.Color) {
		col = color.RGBAModel.Convert(c).(color.RGBA)
	}, w)
	d.Advanced = true

	assert.Nil(t, d.picker)
	d.SetColor(col)
	assert.NotNil(t, d.picker)

	assert.Equal(t, uint8(70), d.picker.Red)
	assert.Equal(t, uint8(210), d.picker.Green)
	assert.Equal(t, uint8(200), d.picker.Blue)
	assert.Equal(t, uint8(255), d.picker.Alpha)

	col = color.RGBA{R: 244, G: 40, B: 70, A: 244}
	assert.NotEqual(t, col.R, d.picker.Red)
	assert.NotEqual(t, col.G, d.picker.Green)
	assert.NotEqual(t, col.B, d.picker.Blue)
	assert.NotEqual(t, col.A, d.picker.Alpha)

	d.SetColor(col)
	assert.Equal(t, uint8(255), d.picker.Red)
	assert.Equal(t, uint8(41), d.picker.Green)
	assert.Equal(t, uint8(73), d.picker.Blue)
	assert.Equal(t, uint8(244), d.picker.Alpha)

	d.Show()
}

func TestColorDialogSimple_Theme(t *testing.T) {
	test.NewTempApp(t)

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(600, 400))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	test.AssertRendersToImage(t, "color/dialog_simple_theme_default.png", w.Canvas())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertRendersToImage(t, "color/dialog_simple_theme_ugly.png", w.Canvas())
}

func TestColorDialogSimple_Recents(t *testing.T) {
	a := test.NewTempApp(t)

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#2196f3,#4caf50,#f44336")

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(600, 400))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	test.AssertRendersToImage(t, "color/dialog_simple_recents_theme_default.png", w.Canvas())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertRendersToImage(t, "color/dialog_simple_recents_theme_ugly.png", w.Canvas())
}

func Test_recent_color(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		test.NewTempApp(t)
		colors := readRecentColors()
		assert.Empty(t, colors)
	})
	t.Run("Single", func(t *testing.T) {
		test.NewTempApp(t)
		writeRecentColor("#ff0000") // Red
		colors := readRecentColors()
		assert.Len(t, colors, 1)
		assert.Equal(t, "#ff0000", colors[0])
	})
	t.Run("Order", func(t *testing.T) {
		test.NewTempApp(t)
		// Recents are last in, first out
		writeRecentColor("#ff0000") // Red
		writeRecentColor("#00ff00") // Green
		writeRecentColor("#0000ff") // Blue
		colors := readRecentColors()
		assert.Len(t, colors, 3)
		assert.Equal(t, "#0000ff", colors[0])
		assert.Equal(t, "#00ff00", colors[1])
		assert.Equal(t, "#ff0000", colors[2])
	})
	t.Run("Deduplicate", func(t *testing.T) {
		test.NewTempApp(t)
		// Ensure no duplicates
		writeRecentColor("#ff0000") // Red
		writeRecentColor("#00ff00") // Green
		writeRecentColor("#0000ff") // Blue
		writeRecentColor("#ff0000") // Red again
		colors := readRecentColors()
		assert.Len(t, colors, 3)
		assert.Equal(t, "#ff0000", colors[0]) // Red
		assert.Equal(t, "#0000ff", colors[1]) // Blue
		assert.Equal(t, "#00ff00", colors[2]) // Green
	})
	t.Run("Limit", func(t *testing.T) {
		test.NewTempApp(t)
		// Max recents is 7
		writeRecentColor("#000000") // Black
		writeRecentColor("#bbbbbb") // Dark Grey
		writeRecentColor("#444444") // Light Grey
		writeRecentColor("#ffffff") // White
		writeRecentColor("#ff0000") // Red
		writeRecentColor("#00ff00") // Green
		writeRecentColor("#0000ff") // Blue
		writeRecentColor("#ffff00") // Yellow
		writeRecentColor("#00ffff") // Cyan
		writeRecentColor("#ff00ff") // Magenta
		colors := readRecentColors()
		assert.Len(t, colors, 7)
		assert.Equal(t, "#ff00ff", colors[0]) // Magenta
		assert.Equal(t, "#00ffff", colors[1]) // Cyan
		assert.Equal(t, "#ffff00", colors[2]) // Yellow
		assert.Equal(t, "#0000ff", colors[3]) // Blue
		assert.Equal(t, "#00ff00", colors[4]) // Green
		assert.Equal(t, "#ff0000", colors[5]) // Red
		assert.Equal(t, "#ffffff", colors[6]) // White
	})
}

func Test_clamp(t *testing.T) {
	// No Change
	assert.Equal(t, 5, clamp(5, 0, 5))
	// Clamp min
	assert.Equal(t, 0, clamp(-1, 0, 5))
	// Clamp max
	assert.Equal(t, 5, clamp(6, 0, 5))
}

func Test_wrapHue(t *testing.T) {
	// No Change
	assert.Equal(t, 180, wrapHue(180))
	// Wrap to 359
	assert.Equal(t, 359, wrapHue(-1))
	// Wrap to 1
	assert.Equal(t, 1, wrapHue(361))
	// Wrap to 359
	assert.Equal(t, 359, wrapHue(-361))
	// Wrap to 1
	assert.Equal(t, 1, wrapHue(721))
}

type rgbhsl struct {
	hex     string
	r, g, b uint8
	h, s, l int
}

var rgbhslMap = map[string]rgbhsl{
	"black": {
		hex: "#000000",
	},
	"white": {
		hex: "#ffffff",
		r:   255,
		g:   255,
		b:   255,
		l:   100,
	},
	"red": {
		hex: "#ff0000",
		r:   255,
		s:   100,
		l:   50,
	},
	"green": {
		hex: "#00ff00",
		g:   255,
		h:   120,
		s:   100,
		l:   50,
	},
	"blue": {
		hex: "#0000ff",
		b:   255,
		h:   240,
		s:   100,
		l:   50,
	},
	"yellow": {
		hex: "#ffff00",
		r:   255,
		g:   255,
		h:   60,
		s:   100,
		l:   50,
	},
	"cyan": {
		hex: "#00ffff",
		g:   255,
		b:   255,
		h:   180,
		s:   100,
		l:   50,
	},
	"magenta": {
		hex: "#ff00ff",
		r:   255,
		b:   255,
		h:   300,
		s:   100,
		l:   50,
	},
}

func Test_colorToString(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			hex := colorToString(&color.NRGBA{
				R: tt.r,
				G: tt.g,
				B: tt.b,
				A: 0xff,
			})
			assert.Equal(t, tt.hex, hex)
		})
	}
}

func Test_colorToHSLA(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			var c color.Color = &color.NRGBA{
				R: tt.r,
				G: tt.g,
				B: tt.b,
				A: 0xff,
			}
			r, g, b, a := fynecolor.ToNRGBA(c)
			h, s, l := rgbToHsl(r, g, b)
			assert.Equal(t, tt.h, h)
			assert.Equal(t, tt.s, s)
			assert.Equal(t, tt.l, l)
			assert.Equal(t, uint8(255), a)
		})
	}
}

func Test_toRGBA(t *testing.T) {
	// Test_toRGBA is only still here instead of with ToNRGBA because it uses rgbhslMap
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			r, g, b, a := fynecolor.ToNRGBA(&color.NRGBA{
				R: tt.r,
				G: tt.g,
				B: tt.b,
				A: 0xff,
			})
			assert.Equal(t, tt.r, r)
			assert.Equal(t, tt.g, g)
			assert.Equal(t, tt.b, b)
			assert.Equal(t, uint8(0xff), a)
		})
	}
}

func Test_rgbToHsl(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			h, s, l := rgbToHsl(tt.r, tt.g, tt.b)
			assert.Equal(t, tt.h, h)
			assert.Equal(t, tt.s, s)
			assert.Equal(t, tt.l, l)
		})
	}
}

func Test_hslToRgb(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			r, g, b := hslToRgb(tt.h, tt.s, tt.l)
			assert.Equal(t, tt.r, r)
			assert.Equal(t, tt.g, g)
			assert.Equal(t, tt.b, b)
		})
	}
}

func Test_hueToChannel(t *testing.T) {
	for name, tt := range map[string]struct {
		h, v1, v2 float64
		expected  float64
	}{
		"red": {
			h:        0,
			v1:       1,
			expected: 0,
		},
		"green": {
			h:        0.3333333333333333,
			v1:       1,
			expected: 1,
		},
		"blue": {
			h:        0.6666666666666666,
			v1:       1,
			expected: 0,
		},
		"cyan": {
			h:        0.5,
			v1:       1,
			expected: 1,
		},
		"yellow": {
			h:        0.16666666666666666,
			v1:       1,
			expected: 1,
		},
		"magenta": {
			h:        0.8333333333333334,
			v1:       1,
			expected: 0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, hueToChannel(tt.h, tt.v1, tt.v2), 0.000000000000001)
		})
	}
}
