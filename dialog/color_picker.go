package dialog

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// newColorBasicPicker returns a component for selecting basic colors.
func newColorBasicPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox([]color.Color{
		theme.PrimaryColorNamed(theme.ColorRed),
		theme.PrimaryColorNamed(theme.ColorOrange),
		theme.PrimaryColorNamed(theme.ColorYellow),
		theme.PrimaryColorNamed(theme.ColorGreen),
		theme.PrimaryColorNamed(theme.ColorBlue),
		theme.PrimaryColorNamed(theme.ColorPurple),
		theme.PrimaryColorNamed(theme.ColorBrown),
		// theme.PrimaryColorNamed(theme.ColorGray),
	}, theme.ColorChromaticIcon(), callback)
}

// newColorGreyscalePicker returns a component for selecting greyscale colors.
func newColorGreyscalePicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors([]string{
		"#ffffff",
		"#cccccc",
		"#aaaaaa",
		"#808080",
		"#555555",
		"#333333",
		"#000000",
	}...), theme.ColorAchromaticIcon(), callback)
}

// newColorRecentPicker returns a component for selecting recent colors.
func newColorRecentPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors(readRecentColors()...), theme.HistoryIcon(), callback)
}

var _ fyne.Widget = (*colorAdvancedPicker)(nil)

// colorAdvancedPicker widget is a component for selecting a color.
type colorAdvancedPicker struct {
	widget.BaseWidget
	Red, Green, Blue, Alpha int // Range 0-255
	Hue                     int // Range 0-360 (degrees)
	Saturation, Lightness   int // Range 0-100 (percent)
	ColorModel              string
	previousColor           color.Color

	onChange func(color.Color)
}

// newColorAdvancedPicker returns a new color widget set to the given color.
func newColorAdvancedPicker(color color.Color, onChange func(color.Color)) *colorAdvancedPicker {
	c := &colorAdvancedPicker{
		onChange: onChange,
	}
	c.ExtendBaseWidget(c)
	c.previousColor = color
	c.updateColor(color)
	return c
}

// Color returns the currently selected color.
func (p *colorAdvancedPicker) Color() color.Color {
	return &color.NRGBA{
		uint8(p.Red),
		uint8(p.Green),
		uint8(p.Blue),
		uint8(p.Alpha),
	}
}

// SetColor updates the color selected in this color widget.
func (p *colorAdvancedPicker) SetColor(color color.Color) {
	p.previousColor = color
	if p.updateColor(color) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(color)
		}
	}
}

// SetHSLA updated the Hue, Saturation, Lightness, and Alpha components of the currently selected color.
func (p *colorAdvancedPicker) SetHSLA(h, s, l, a int) {
	if p.updateHSLA(h, s, l, a) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(p.Color())
		}
	}
}

// SetRGBA updated the Red, Green, Blue, and Alpha components of the currently selected color.
func (p *colorAdvancedPicker) SetRGBA(r, g, b, a int) {
	if p.updateRGBA(r, g, b, a) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(p.Color())
		}
	}
}

// MinSize returns the size that this widget should not shrink below.
func (p *colorAdvancedPicker) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (p *colorAdvancedPicker) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)

	// Preview
	preview := newColorPreview(p.previousColor)

	// HSL
	hueChannel := newColorChannel("H", 0, 360, p.Hue, func(h int) {
		p.SetHSLA(h, p.Saturation, p.Lightness, p.Alpha)
	})
	saturationChannel := newColorChannel("S", 0, 100, p.Saturation, func(s int) {
		p.SetHSLA(p.Hue, s, p.Lightness, p.Alpha)
	})
	lightnessChannel := newColorChannel("L", 0, 100, p.Lightness, func(l int) {
		p.SetHSLA(p.Hue, p.Saturation, l, p.Alpha)
	})
	hslBox := container.NewVBox(
		hueChannel,
		saturationChannel,
		lightnessChannel,
	)

	// RGB
	redChannel := newColorChannel("R", 0, 255, p.Red, func(r int) {
		p.SetRGBA(r, p.Green, p.Blue, p.Alpha)
	})
	greenChannel := newColorChannel("G", 0, 255, p.Green, func(g int) {
		p.SetRGBA(p.Red, g, p.Blue, p.Alpha)
	})
	blueChannel := newColorChannel("B", 0, 255, p.Blue, func(b int) {
		p.SetRGBA(p.Red, p.Green, b, p.Alpha)
	})
	rgbBox := container.NewVBox(
		redChannel,
		greenChannel,
		blueChannel,
	)

	// Wheel
	wheel := newColorWheel(func(hue, saturation, lightness, alpha int) {
		p.SetHSLA(hue, saturation, lightness, alpha)
	})

	// Alpha
	alphaChannel := newColorChannel("A", 0, 255, p.Alpha, func(a int) {
		p.SetRGBA(p.Red, p.Green, p.Blue, a)
	})

	// Hex
	hex := newUserChangeEntry("")
	hex.setOnChanged(func(text string) {
		c, err := stringToColor(text)
		if err != nil {
			fyne.LogError("Error parsing color: "+text, err)
			// TODO trigger entry invalid state
		} else {
			p.SetColor(c)
		}
	})

	contents := container.NewPadded(container.NewVBox(
		container.NewGridWithColumns(3,
			container.NewPadded(wheel),
			hslBox,
			rgbBox),
		container.NewGridWithColumns(3,
			container.NewPadded(preview),

			hex,
			alphaChannel,
		),
	))

	r := &colorPickerRenderer{
		WidgetRenderer:    widget.NewSimpleRenderer(contents),
		picker:            p,
		redChannel:        redChannel,
		greenChannel:      greenChannel,
		blueChannel:       blueChannel,
		hueChannel:        hueChannel,
		saturationChannel: saturationChannel,
		lightnessChannel:  lightnessChannel,
		wheel:             wheel,
		preview:           preview,
		alphaChannel:      alphaChannel,
		hex:               hex,
		contents:          contents,
	}
	r.updateObjects()
	return r
}

func (p *colorAdvancedPicker) updateColor(color color.Color) bool {
	r, g, b, a := col.ToNRGBA(color)
	if p.Red == r && p.Green == g && p.Blue == b && p.Alpha == a {
		return false
	}
	return p.updateRGBA(r, g, b, a)
}

func (p *colorAdvancedPicker) updateHSLA(h, s, l, a int) bool {
	h = wrapHue(h)
	s = clamp(s, 0, 100)
	l = clamp(l, 0, 100)
	a = clamp(a, 0, 255)
	if p.Hue == h && p.Saturation == s && p.Lightness == l && p.Alpha == a {
		return false
	}
	p.Hue = h
	p.Saturation = s
	p.Lightness = l
	p.Alpha = a
	p.Red, p.Green, p.Blue = hslToRgb(p.Hue, p.Saturation, p.Lightness)
	return true
}

func (p *colorAdvancedPicker) updateRGBA(r, g, b, a int) bool {
	r = clamp(r, 0, 255)
	g = clamp(g, 0, 255)
	b = clamp(b, 0, 255)
	a = clamp(a, 0, 255)
	if p.Red == r && p.Green == g && p.Blue == b && p.Alpha == a {
		return false
	}
	p.Red = r
	p.Green = g
	p.Blue = b
	p.Alpha = a
	p.Hue, p.Saturation, p.Lightness = rgbToHsl(p.Red, p.Green, p.Blue)
	return true
}

var _ fyne.WidgetRenderer = (*colorPickerRenderer)(nil)

type colorPickerRenderer struct {
	fyne.WidgetRenderer
	picker            *colorAdvancedPicker
	redChannel        *colorChannel
	greenChannel      *colorChannel
	blueChannel       *colorChannel
	hueChannel        *colorChannel
	saturationChannel *colorChannel
	lightnessChannel  *colorChannel
	wheel             *colorWheel
	preview           *colorPreview
	alphaChannel      *colorChannel
	hex               *userChangeEntry
	contents          fyne.CanvasObject
}

func (r *colorPickerRenderer) Refresh() {
	r.updateObjects()
	r.WidgetRenderer.Refresh()
}

func (r *colorPickerRenderer) updateObjects() {
	// HSL
	r.hueChannel.SetValue(r.picker.Hue)
	r.saturationChannel.SetValue(r.picker.Saturation)
	r.lightnessChannel.SetValue(r.picker.Lightness)

	// RGB
	r.redChannel.SetValue(r.picker.Red)
	r.greenChannel.SetValue(r.picker.Green)
	r.blueChannel.SetValue(r.picker.Blue)

	// Wheel
	r.wheel.SetHSLA(r.picker.Hue, r.picker.Saturation, r.picker.Lightness, r.picker.Alpha)

	color := r.picker.Color()

	// Preview
	r.preview.SetColor(color)

	// Alpha
	r.alphaChannel.SetValue(r.picker.Alpha)

	// Hex
	r.hex.SetText(colorToString(color))
}
