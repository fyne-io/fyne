package dialog

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ColorPickerDialog is a simple dialog window that displays a color picker.
//
// Since: 1.4
type ColorPickerDialog struct {
	*dialog
	Advanced bool
	color    color.Color
	callback func(c color.Color)
	advanced *widget.Accordion
	picker   *colorAdvancedPicker
}

// NewColorPicker creates a color dialog and returns the handle.
// Using the returned type you should call Show() and then set its color through SetColor().
// The callback is triggered when the user selects a color.
//
// Since: 1.4
func NewColorPicker(title, message string, callback func(c color.Color), parent fyne.Window) *ColorPickerDialog {
	p := &ColorPickerDialog{
		dialog:   newDialog(title, message, theme.ColorPaletteIcon(), nil /*cancel?*/, parent),
		color:    theme.PrimaryColor(),
		callback: callback,
	}
	p.dialog.layout = &dialogLayout{d: p.dialog}

	return p
}

// ShowColorPicker creates and shows a color dialog.
// The callback is triggered when the user selects a color.
//
// Since: 1.4
func ShowColorPicker(title, message string, callback func(c color.Color), parent fyne.Window) {
	NewColorPicker(title, message, callback, parent).Show()
}

// Refresh causes this dialog to be updated
func (p *ColorPickerDialog) Refresh() {
	p.updateUI()
}

// SetColor updates the color of the color picker.
func (p *ColorPickerDialog) SetColor(c color.Color) {
	p.picker.SetColor(c)
}

// Show causes this dialog to be displayed
func (p *ColorPickerDialog) Show() {
	if p.win == nil || p.Advanced != (p.advanced != nil) {
		p.updateUI()
	}
	p.dialog.Show()
}

func (p *ColorPickerDialog) createSimplePickers() (contents []fyne.CanvasObject) {
	contents = append(contents, newColorBasicPicker(p.selectColor), newColorGreyscalePicker(p.selectColor))
	if recent := newColorRecentPicker(p.selectColor); len(recent.(*fyne.Container).Objects) > 0 {
		// Add divider and recents if there are any
		contents = append(contents, canvas.NewLine(theme.ShadowColor()), recent)
	}
	return
}

func (p *ColorPickerDialog) selectColor(c color.Color) {
	p.dialog.Hide()
	writeRecentColor(colorToString(p.color))
	if f := p.callback; f != nil {
		f(c)
	}
}

func (p *ColorPickerDialog) updateUI() {
	if w := p.win; w != nil {
		w.Hide()
	}
	p.dialog.dismiss = &widget.Button{Text: "Cancel", Icon: theme.CancelIcon(),
		OnTapped: p.dialog.Hide,
	}
	if p.Advanced {
		p.picker = newColorAdvancedPicker(p.color, func(c color.Color) {
			p.color = c
		})
		p.advanced = widget.NewAccordion(widget.NewAccordionItem("Advanced", p.picker))

		p.dialog.content = fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
			fyne.NewContainerWithLayout(layout.NewCenterLayout(),
				fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
					p.createSimplePickers()...,
				),
			),
			widget.NewSeparator(),
			p.advanced,
		)

		confirm := &widget.Button{Text: "Confirm", Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
			OnTapped: func() {
				p.selectColor(p.color)
			},
		}
		p.dialog.setButtons(newButtonList(p.dialog.dismiss, confirm))
	} else {
		p.dialog.content = fyne.NewContainerWithLayout(layout.NewVBoxLayout(), p.createSimplePickers()...)
		p.dialog.setButtons(newButtonList(p.dialog.dismiss))
	}
}

func clamp(value, min, max int) int {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}

func wrapHue(hue int) int {
	for hue < 0 {
		hue += 360
	}
	for hue > 360 {
		hue -= 360
	}
	return hue
}

func newColorButtonBox(colors []color.Color, icon fyne.Resource, callback func(color.Color)) fyne.CanvasObject {
	var objects []fyne.CanvasObject
	if icon != nil && len(colors) > 0 {
		objects = append(objects, widget.NewIcon(icon))
	}
	for _, c := range colors {
		objects = append(objects, newColorButton(c, callback))
	}
	return fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(8), objects...)
}

func newCheckeredBackground() *canvas.Raster {
	return canvas.NewRasterWithPixels(func(x, y, _, _ int) color.Color {
		const boxSize = 10

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
	red, green, blue, alpha := col.ToNRGBA(c)
	if alpha == 0xff {
		return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", red, green, blue, alpha)
}

func stringToColor(s string) (color.Color, error) {
	var c color.NRGBA
	var err error
	if len(s) == 7 {
		c.A = 0xFF
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	} else {
		_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	}
	return c, err
}

func stringsToColors(ss ...string) (colors []color.Color) {
	for _, s := range ss {
		if s == "" {
			continue
		}
		c, err := stringToColor(s)
		if err != nil {
			fyne.LogError("Couldn't parse color:", err)
		} else {
			colors = append(colors, c)
		}
	}
	return
}

func colorToHSLA(c color.Color) (int, int, int, int) {
	r, g, b, a := col.ToNRGBA(c)
	h, s, l := rgbToHsl(r, g, b)
	return h, s, l, a
}

// https://www.niwa.nu/2013/05/math-behind-colorspace-conversions-rgb-hsl/

func rgbToHsl(r, g, b int) (int, int, int) {
	red := float64(r) / 255.0
	green := float64(g) / 255.0
	blue := float64(b) / 255.0

	min := math.Min(red, math.Min(green, blue))
	max := math.Max(red, math.Max(green, blue))

	lightness := (max + min) / 2.0

	delta := max - min

	if delta == 0.0 {
		// Achromatic
		return 0, 0, int(lightness * 100.0)
	}

	// Chromatic

	var saturation float64

	if lightness < 0.5 {
		saturation = (max - min) / (max + min)
	} else {
		saturation = (max - min) / (2.0 - max - min)
	}

	var hue float64

	if red == max {
		hue = (green - blue) / delta
	} else if green == max {
		hue = 2.0 + (blue-red)/delta
	} else if blue == max {
		hue = 4.0 + (red-green)/delta
	}

	h := wrapHue(int(hue * 60.0))
	s := int(saturation * 100.0)
	l := int(lightness * 100.0)
	return h, s, l
}

func hslToRgb(h, s, l int) (int, int, int) {
	hue := float64(h) / 360.0
	saturation := float64(s) / 100.0
	lightness := float64(l) / 100.0

	if saturation == 0.0 {
		// Greyscale
		g := int(lightness * 255.0)
		return g, g, g
	}

	var v1 float64
	if lightness < 0.5 {
		v1 = lightness * (1.0 + saturation)
	} else {
		v1 = (lightness + saturation) - (lightness * saturation)
	}

	v2 := 2.0*lightness - v1

	red := hueToChannel(hue+(1.0/3.0), v1, v2)
	green := hueToChannel(hue, v1, v2)
	blue := hueToChannel(hue-(1.0/3.0), v1, v2)

	r := int(math.Round(255.0 * red))
	g := int(math.Round(255.0 * green))
	b := int(math.Round(255.0 * blue))

	return r, g, b
}

func hueToChannel(h, v1, v2 float64) float64 {
	for h < 0.0 {
		h += 1.0
	}
	for h > 1.0 {
		h -= 1.0
	}
	if 6.0*h < 1.0 {
		return v2 + (v1-v2)*6*h
	}
	if 2.0*h < 1.0 {
		return v1
	}
	if 3.0*h < 2.0 {
		return v2 + (v1-v2)*6*((2.0/3.0)-h)
	}
	return v2
}
