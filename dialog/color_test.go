package dialog

import (
	"image/color"
	"testing"

	"math"

	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

func TestColorDialog_Theme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(500, 300))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	test.AssertImageMatches(t, "color/light.png", w.Canvas().Capture())

	test.ApplyTheme(t, theme.DarkTheme())
	test.AssertImageMatches(t, "color/dark.png", w.Canvas().Capture())

	w.Close()
}

func TestColorDialog_Advanced_Theme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(800, 800))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	d.advanced.Open(0)
	d.Resize(d.win.MinSize()) // TODO FIXME Hack

	test.AssertImageMatches(t, "color/advanced_light.png", w.Canvas().Capture())

	test.ApplyTheme(t, theme.DarkTheme())
	test.AssertImageMatches(t, "color/advanced_dark.png", w.Canvas().Capture())

	w.Close()
}

func TestColorDialog_Recents(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#2196f3,#4caf50,#f44336")

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(500, 300))

	d := NewColorPicker("Color Picker", "Pick a Color", nil, w)
	d.Show()

	test.AssertImageMatches(t, "color/recents.png", w.Canvas().Capture())

	w.Close()
}

func Test_recent_color(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		test.NewApp()
		defer test.NewApp()
		colors := readRecentColors()
		assert.Equal(t, 0, len(colors))
	})
	t.Run("Single", func(t *testing.T) {
		test.NewApp()
		defer test.NewApp()
		writeRecentColor("#ff0000") // Red
		colors := readRecentColors()
		assert.Equal(t, 1, len(colors))
		assert.Equal(t, "#ff0000", colors[0])
	})
	t.Run("Order", func(t *testing.T) {
		test.NewApp()
		defer test.NewApp()
		// Recents are last in, first out
		writeRecentColor("#ff0000") // Red
		writeRecentColor("#00ff00") // Green
		writeRecentColor("#0000ff") // Blue
		colors := readRecentColors()
		assert.Equal(t, 3, len(colors))
		assert.Equal(t, "#0000ff", colors[0])
		assert.Equal(t, "#00ff00", colors[1])
		assert.Equal(t, "#ff0000", colors[2])
	})
	t.Run("Deduplicate", func(t *testing.T) {
		test.NewApp()
		defer test.NewApp()
		// Ensure no duplicates
		writeRecentColor("#ff0000") // Red
		writeRecentColor("#00ff00") // Green
		writeRecentColor("#0000ff") // Blue
		writeRecentColor("#ff0000") // Red again
		colors := readRecentColors()
		assert.Equal(t, 3, len(colors))
		assert.Equal(t, "#ff0000", colors[0]) // Red
		assert.Equal(t, "#0000ff", colors[1]) // Blue
		assert.Equal(t, "#00ff00", colors[2]) // Green
	})
	t.Run("Limit", func(t *testing.T) {
		test.NewApp()
		defer test.NewApp()
		// Max recents is 8
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
		assert.Equal(t, 8, len(colors))
		assert.Equal(t, "#ff00ff", colors[0]) // Magenta
		assert.Equal(t, "#00ffff", colors[1]) // Cyan
		assert.Equal(t, "#ffff00", colors[2]) // Yellow
		assert.Equal(t, "#0000ff", colors[3]) // Blue
		assert.Equal(t, "#00ff00", colors[4]) // Green
		assert.Equal(t, "#ff0000", colors[5]) // Red
		assert.Equal(t, "#ffffff", colors[6]) // White
		assert.Equal(t, "#444444", colors[7]) // Light Grey
	})
}

func Test_colorClamp(t *testing.T) {
	// No Change
	assert.Equal(t, 0.5, colorClamp(0.5))
	// Clamp 0.0
	assert.Equal(t, 0.0, colorClamp(-0.1))
	// Clamp 1.0
	assert.Equal(t, 1.0, colorClamp(1.1))
}

type rgbhsl struct {
	hex     string
	r, g, b float64
	h, s, l float64
}

var rgbhslMap = map[string]rgbhsl{
	"black": {
		hex: "#000000",
	},
	"white": {
		hex: "#ffffff",
		r:   1.0,
		g:   1.0,
		b:   1.0,
		l:   1.0,
	},
	"red": {
		hex: "#ff0000",
		r:   1.0,
		s:   1.0,
		l:   0.5,
	},
	"green": {
		hex: "#00ff00",
		g:   1.0,
		h:   120.0 / 360.0,
		s:   1.0,
		l:   0.5,
	},
	"blue": {
		hex: "#0000ff",
		b:   1.0,
		h:   240.0 / 360.0,
		s:   1.0,
		l:   0.5,
	},
	"yellow": {
		hex: "#ffff00",
		r:   1.0,
		g:   1.0,
		h:   60.0 / 360.0,
		s:   1.0,
		l:   0.5,
	},
	"cyan": {
		hex: "#00ffff",
		g:   1.0,
		b:   1.0,
		h:   180.0 / 360.0,
		s:   1.0,
		l:   0.5,
	},
	"magenta": {
		hex: "#ff00ff",
		r:   1.0,
		b:   1.0,
		h:   300.0 / 360.0,
		s:   1.0,
		l:   0.5,
	},
}

func Test_colorToString(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			hex := colorToString(&color.NRGBA{
				R: uint8(tt.r * 255),
				G: uint8(tt.g * 255),
				B: uint8(tt.b * 255),
				A: 0xff,
			})
			assert.Equal(t, tt.hex, hex)
		})
	}
}

func Test_stringToColor(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			c := stringToColor(tt.hex)
			r, g, b, _ := c.RGBA()
			assert.InDelta(t, tt.r*65535, r, 0.0001)
			assert.InDelta(t, tt.g*65535, g, 0.0001)
			assert.InDelta(t, tt.b*65535, b, 0.0001)
		})
	}
}

func Test_colorToHSLA(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			h, s, l, a := colorToHSLA(&color.NRGBA{
				R: uint8(tt.r * 255),
				G: uint8(tt.g * 255),
				B: uint8(tt.b * 255),
				A: 0xff,
			})
			assert.InDelta(t, tt.h, h, 0.0001)
			assert.InDelta(t, tt.s, s, 0.0001)
			assert.InDelta(t, tt.l, l, 0.0001)
			assert.InDelta(t, 1.0, a, 0.0001)
		})
	}
}

func Test_colorToRGBA(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			r, g, b, a := colorToRGBA(&color.NRGBA{
				R: uint8(tt.r * 255),
				G: uint8(tt.g * 255),
				B: uint8(tt.b * 255),
				A: 0xff,
			})
			assert.InDelta(t, tt.r, r, 0.0001)
			assert.InDelta(t, tt.g, g, 0.0001)
			assert.InDelta(t, tt.b, b, 0.0001)
			assert.InDelta(t, 1.0, a, 0.0001)
		})
	}
}

type xyarl struct {
	x, y    int
	a, r, l float64
}

var (
	width    = 100.0
	height   = 100.0
	xyarlMap = map[string]xyarl{
		"north": {
			x: 50,
			y: 90,
			a: math.Pi / 2,
			r: 40,
			l: 50,
		},
		"northeast": {
			x: 90,
			y: 90,
			a: math.Pi / 4,
			r: math.Sqrt(40 * 40 * 2),
			l: 50,
		},
		"east": {
			x: 90,
			y: 50,
			r: 40,
			l: 50,
		},
		"southeast": {
			x: 90,
			y: 10,
			a: -math.Pi / 4,
			r: math.Sqrt(40 * 40 * 2),
			l: 50,
		},
		"south": {
			x: 50,
			y: 10,
			a: -math.Pi / 2,
			r: 40,
			l: 50,
		},
		"southwest": {
			x: 10,
			y: 10,
			a: -math.Pi * 3 / 4,
			r: math.Sqrt(40 * 40 * 2),
			l: 50,
		},
		"west": {
			x: 10,
			y: 50,
			a: math.Pi,
			r: 40,
			l: 50,
		},
		"northwest": {
			x: 10,
			y: 90,
			a: math.Pi * 3 / 4,
			r: math.Sqrt(40 * 40 * 2),
			l: 50,
		},
	}
)

func Test_cartesianToPolar(t *testing.T) {
	for name, tt := range xyarlMap {
		t.Run(name, func(t *testing.T) {
			a, r, l := cartesianToPolar(float64(tt.x), float64(tt.y), width, height)
			assert.InDelta(t, tt.a, a, 0.0001)
			assert.InDelta(t, tt.r, r, 0.0001)
			assert.InDelta(t, tt.l, l, 0.0001)
		})
	}
}

func Test_polarToCartesian(t *testing.T) {
	for name, tt := range xyarlMap {
		t.Run(name, func(t *testing.T) {
			x, y := polarToCartesian(tt.a, tt.r, width, height)
			assert.InDelta(t, tt.x, x, 0.0001)
			assert.InDelta(t, tt.y, y, 0.0001)
		})
	}
}

func Test_polarToHS(t *testing.T) {
	// TODO
}

func Test_rgbToHsl(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			h, s, l := rgbToHsl(tt.r, tt.g, tt.b)
			assert.InDelta(t, tt.h, h, 0.0001)
			assert.InDelta(t, tt.s, s, 0.0001)
			assert.InDelta(t, tt.l, l, 0.0001)
		})
	}
}

func Test_hslToRgb(t *testing.T) {
	for name, tt := range rgbhslMap {
		t.Run(name, func(t *testing.T) {
			r, g, b := hslToRgb(tt.h, tt.s, tt.l)
			assert.InDelta(t, tt.r, r, 0.0001)
			assert.InDelta(t, tt.g, g, 0.0001)
			assert.InDelta(t, tt.b, b, 0.0001)
		})
	}
}

func Test_hueToRgb(t *testing.T) {
	// TODO
}
