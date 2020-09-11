package dialog

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ColorPickerDialog is a simple dialog window that displays a color picker.
type ColorPickerDialog struct {
	*dialog
	color    color.Color
	advanced *widget.AccordionContainer
	picker   *colorAdvancedPicker
}

// SetColor updates the color of the color picker.
func (p *ColorPickerDialog) SetColor(c color.Color) {
	p.picker.SetColor(c)
}

// NewColorPicker creates a color dialog and returns the handle.
// Using the returned type you should call Show() and then set its color through SetColor().
// The callback is triggered when the user selects a color.
func NewColorPicker(title, message string, callback func(c color.Color), parent fyne.Window) *ColorPickerDialog {
	d := newDialog(title, message, theme.ColorPaletteIcon(), nil /*cancel?*/, parent)

	cb := func(c color.Color) {
		d.Hide()
		// TODO FIXME writeRecentColor(colorToString(p.color))
		callback(c)
	}

	basic := newColorBasicPicker(cb)

	greyscale := newColorGreyscalePicker(cb)

	recent := newColorRecentPicker(cb)

	p := &ColorPickerDialog{
		dialog: d,
		color:  theme.PrimaryColor(),
	}

	p.picker = newColorAdvancedPicker(theme.PrimaryColor(), func(c color.Color) {
		p.color = c
	})

	p.advanced = widget.NewAccordionContainer(widget.NewAccordionItem("Advanced", p.picker))

	contents := []fyne.CanvasObject{
		basic,
		greyscale,
	}
	if !recent.MinSize().IsZero() {
		// Add divider and recents if there are any
		contents = append(contents, canvas.NewLine(theme.ShadowColor()), recent)
	}

	d.content = widget.NewVBox(append(contents, p.advanced)...)

	d.dismiss = &widget.Button{Text: "Cancel", Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	confirm := &widget.Button{Text: "Confirm", Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
		OnTapped: func() {
			cb(p.color)
		},
	}
	d.setButtons(newButtonList(d.dismiss, confirm))
	return p
}

// ShowColorPicker creates and shows a color dialog.
// The callback is triggered when the user selects a color.
func ShowColorPicker(title, message string, callback func(c color.Color), parent fyne.Window) {
	NewColorPicker(title, message, callback, parent).Show()
}

func colorClamp(value float64) float64 {
	if value < 0.0 {
		value = 0.0
	}
	if value > 1.0 {
		value = 1.0
	}
	return value
}

func newColorButtonBox(colors []color.Color, icon fyne.Resource, callback func(color.Color)) fyne.CanvasObject {
	box := widget.NewHBox()
	if icon != nil && len(colors) > 0 {
		box.Children = append(box.Children, widget.NewIcon(icon))
	}
	for _, c := range colors {
		box.Children = append(box.Children, newColorButton(c, callback))
	}
	return box
}

func newCheckeredBackground() *canvas.Raster {
	return canvas.NewRasterWithPixels(func(x, y, _, _ int) color.Color {
		const boxSize = 25

		if (x/boxSize)%2 == (y/boxSize)%2 {
			return color.Gray{Y: 58}
		}

		return color.Gray{Y: 84}
	})
}

const (
	preferenceRecents    = "color_recents"
	preferenceMaxRecents = 8
)

func readRecentColors() (recents []string) {
	for _, r := range strings.Split(fyne.CurrentApp().Preferences().String(preferenceRecents), ",") {
		if r != "" {
			recents = append(recents, r)
		}
	}
	return
}

func writeRecentColor(color string) {
	recents := []string{color}
	for _, r := range readRecentColors() {
		if r == color {
			continue // Color already in recents
		}
		recents = append(recents, r)
	}
	if len(recents) > preferenceMaxRecents {
		recents = recents[:preferenceMaxRecents]
	}
	fyne.CurrentApp().Preferences().SetString(preferenceRecents, strings.Join(recents, ","))
}

func colorToString(c color.Color) string {
	var red, green, blue, alpha uint8
	switch col := c.(type) {
	case *color.NRGBA:
		red, green, blue, alpha = col.R, col.G, col.B, col.A
	default:
		r, g, b, a := c.RGBA() // TODO FIXME this returns alpha-pre-multiplied
		// TODO FIXME possible loss of precision, should this be x/65535*255?
		red = uint8(r)
		green = uint8(g)
		blue = uint8(b)
		alpha = uint8(a)
	}
	if alpha == 0xff {
		return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", red, green, blue, alpha)
}

func stringToColor(s string) color.Color {
	var c color.NRGBA
	var err error
	if len(s) == 7 {
		c.A = 0xFF
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	} else {
		_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	}
	if err != nil {
		fyne.LogError("Couldn't parse color:", err)
	}
	return c
}

func stringsToColors(ss ...string) (colors []color.Color) {
	for _, s := range ss {
		if s == "" {
			continue
		}
		colors = append(colors, stringToColor(s))
	}
	return
}

func colorToHSLA(c color.Color) (float64, float64, float64, float64) {
	r, g, b, a := colorToRGBA(c)
	h, s, l := rgbToHsl(r, g, b)
	return h, s, l, a
}

func colorToRGBA(c color.Color) (r, g, b, a float64) {
	switch col := c.(type) {
	case *color.NRGBA:
		// Convert to range 0.0-1.0
		r = float64(col.R) / 255.0
		g = float64(col.G) / 255.0
		b = float64(col.B) / 255.0
		a = float64(col.A) / 255.0
	default:
		red, green, blue, alpha := c.RGBA() // TODO FIXME this returns alpha-pre-multiplied
		// Convert to range 0.0-1.0
		r = float64(red) / 65535.0
		g = float64(green) / 65535.0
		b = float64(blue) / 65535.0
		a = float64(alpha) / 65535.0
	}
	return
}

func cartesianToPolar(x, y, width, height float64) (angle, radius, limit float64) {
	dx := float64(x - (width / 2))
	dy := float64(y - (height / 2))
	angle = math.Atan2(dy, dx)
	radius = math.Sqrt(dx*dx + dy*dy)
	limit = math.Min(width, height) / 2.0
	return
}

func polarToCartesian(angle, radius, width, height float64) (x, y float64) {
	x = (radius * math.Cos(angle)) + (width / 2)
	y = (radius * math.Sin(angle)) + (height / 2)
	return
}

func polarToHS(angle, radius, limit float64) (hue, saturation float64) {
	hue = angle / (2.0 * math.Pi)
	if hue < 0.0 {
		hue += 1.0
	}
	if hue > 1.0 {
		hue -= 1.0
	}
	saturation = radius / limit
	return
}

// The next three functions are adpated from https://play.golang.org/p/9q5yBNDh3W

func rgbToHsl(r, g, b float64) (float64, float64, float64) {
	min := math.Min(r, math.Min(g, b)) //Min. value of RGB
	max := math.Max(r, math.Max(g, b)) //Max. value of RGB
	del := max - min                   //Delta RGB value

	l := (max + min) / 2.0

	var h, s float64

	if del == 0.0 { // Achromatic
		h = 0.0
		s = 0.0
	} else { // Chromatic
		if l < 0.5 {
			s = del / (max + min)
		} else {
			s = del / (2.0 - max - min)
		}

		delR := (((max - r) / 6.0) + (del / 2.0)) / del
		delG := (((max - g) / 6.0) + (del / 2.0)) / del
		delB := (((max - b) / 6.0) + (del / 2.0)) / del

		if r == max {
			h = delB - delG
		} else if g == max {
			h = (1.0 / 3.0) + delR - delB
		} else if b == max {
			h = (2.0 / 3.0) + delG - delR
		}

		if h < 0.0 {
			h += 1.0
		}
		if h > 1.0 {
			h -= 1.0
		}

	}
	return h, s, l
}

func hslToRgb(h, s, l float64) (float64, float64, float64) {
	var r, g, b float64
	if s == 0.0 {
		r = l
		g = l
		b = l
	} else {
		var v1, v2 float64
		if l < 0.5 {
			v2 = l * (1.0 + s)
		} else {
			v2 = (l + s) - (s * l)
		}

		v1 = 2.0*l - v2

		r = hueToRgb(v1, v2, h+(1.0/3.0))
		g = hueToRgb(v1, v2, h)
		b = hueToRgb(v1, v2, h-(1.0/3.0))
	}
	return r, g, b
}

func hueToRgb(v1, v2, vH float64) float64 {
	if vH < 0.0 {
		vH += 1.0
	}
	if vH > 1.0 {
		vH -= 1.0
	}
	if (6.0 * vH) < 1.0 {
		return (v1 + (v2-v1)*6.0*vH)
	}
	if (2.0 * vH) < 1.0 {
		return v2
	}
	if (3.0 * vH) < 2.0 {
		return (v1 + (v2-v1)*((2.0/3.0)-vH)*6.0)
	}
	return v1
}
