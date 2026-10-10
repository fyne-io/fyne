package dialog

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/internal/painter/geom"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// newColorBasicPicker returns a component for selecting basic colors.
func newColorBasicPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(
		stringsToColors(
			"#f44336", // red
			"#ff9800", // orange
			"#ffeb3b", // yellow
			"#8bc34a", // green
			"#296ff6", // blue
			"#9c27b0", // purple
			"#795548", // brown
		),
		theme.ColorChromaticIcon(),
		callback,
	)
}

// newColorGreyscalePicker returns a component for selecting greyscale colors.
func newColorGreyscalePicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(
		stringsToColors(
			"#ffffff",
			"#cccccc",
			"#aaaaaa",
			"#808080",
			"#555555",
			"#333333",
			"#000000",
		),
		theme.ColorAchromaticIcon(),
		callback,
	)
}

// newColorRecentPicker returns a component for selecting recent colors.
func newColorRecentPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors(readRecentColors()...), theme.HistoryIcon(), callback)
}

var _ fyne.Widget = (*colorAdvancedPicker)(nil)

// colorAdvancedPicker widget is a component for selecting a color.
type colorAdvancedPicker struct {
	widget.BaseWidget
	Red, Green, Blue, Alpha uint8 // Range 0-255
	Hue                     int   // Range 0-360 (degrees)
	Saturation, Lightness   int   // Range 0-100 (percent)
	ColorModel              string
	previousColor           color.Color

	onChange func(color.Color)
}

// newColorAdvancedPicker returns a new color widget set to the given color.
func newColorAdvancedPicker(c color.Color, onChange func(color.Color)) *colorAdvancedPicker {
	p := &colorAdvancedPicker{
		onChange: onChange,
	}
	p.ExtendBaseWidget(p)
	p.previousColor = c
	p.updateColor(c)
	return p
}

// Color returns the currently selected color.
func (p *colorAdvancedPicker) Color() color.Color {
	return &color.NRGBA{
		p.Red,
		p.Green,
		p.Blue,
		p.Alpha,
	}
}

// SetColor updates the color selected in this color widget.
func (p *colorAdvancedPicker) SetColor(c color.Color) {
	p.previousColor = c
	if p.updateColor(c) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(c)
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
	hueChannel := newColorChannel("H", 0, geom.AngleFull, p.Hue, func(h int) {
		p.setHSLA(h, p.Saturation, p.Lightness, p.Alpha)
	})
	saturationChannel := newColorChannel("S", 0, 100, p.Saturation, func(s int) {
		p.setHSLA(p.Hue, s, p.Lightness, p.Alpha)
	})
	lightnessChannel := newColorChannel("L", 0, 100, p.Lightness, func(l int) {
		p.setHSLA(p.Hue, p.Saturation, l, p.Alpha)
	})
	hslBox := container.NewVBox(
		hueChannel,
		saturationChannel,
		lightnessChannel,
	)

	// RGB
	redChannel := newColorChannel("R", 0, math.MaxUint8, int(p.Red), func(r int) {
		p.setRGBA(uint8(r), p.Green, p.Blue, p.Alpha) //gosec:disable G115 -- r’s value is limited by newColorChannel
	})
	greenChannel := newColorChannel("G", 0, math.MaxUint8, int(p.Green), func(g int) {
		p.setRGBA(p.Red, uint8(g), p.Blue, p.Alpha) //gosec:disable G115 -- g’s value is limited by newColorChannel
	})
	blueChannel := newColorChannel("B", 0, math.MaxUint8, int(p.Blue), func(b int) {
		p.setRGBA(p.Red, p.Green, uint8(b), p.Alpha) //gosec:disable G115 -- b’s value is limited by newColorChannel
	})
	rgbBox := container.NewVBox(
		redChannel,
		greenChannel,
		blueChannel,
	)

	// Wheel
	wheel := newColorWheel(func(hue, saturation, lightness int, alpha uint8) {
		p.setHSLA(hue, saturation, lightness, alpha)
	})

	// Alpha
	alphaChannel := newColorChannel("A", 0, math.MaxUint8, int(p.Alpha), func(a int) {
		p.setRGBA(p.Red, p.Green, p.Blue, uint8(a)) //gosec:disable G115 -- a’s value is limited by newColorChannel
	})

	// Hex
	hex := newUserChangeEntry("")
	hex.setOnChanged(func(text string) {
		c, err := col.Parse(text)
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
		container.NewGridWithColumns(
			3,
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

// setHSLA updates the Hue, Saturation, Lightness, and Alpha components of the currently selected color.
func (p *colorAdvancedPicker) setHSLA(h, s, l int, a uint8) {
	if p.updateHSLA(h, s, l, a) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(p.Color())
		}
	}
}

// setRGBA updates the Red, Green, Blue, and Alpha components of the currently selected color.
func (p *colorAdvancedPicker) setRGBA(r, g, b, a uint8) {
	if p.updateRGBA(r, g, b, a) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(p.Color())
		}
	}
}

func (p *colorAdvancedPicker) updateColor(c color.Color) bool {
	r, g, b, a := col.ToNRGBA(c)
	return p.updateRGBA(r, g, b, a)
}

func (p *colorAdvancedPicker) updateHSLA(h, s, l int, a uint8) bool {
	h = wrapHue(h)
	s = clamp(s, 0, 100)
	l = clamp(l, 0, 100)
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

func (p *colorAdvancedPicker) updateRGBA(r, g, b, a uint8) bool {
	if p.Red == r && p.Green == g && p.Blue == b && p.Alpha == a {
		return false
	}

	p.Red = r
	p.Green = g
	p.Blue = b
	p.Alpha = a
	p.Hue, p.Saturation, p.Lightness = rgbToHsl(r, g, b)
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
	r.redChannel.SetValue(int(r.picker.Red))
	r.greenChannel.SetValue(int(r.picker.Green))
	r.blueChannel.SetValue(int(r.picker.Blue))

	// Wheel
	r.wheel.SetHSLA(r.picker.Hue, r.picker.Saturation, r.picker.Lightness, r.picker.Alpha)

	c := r.picker.Color()

	// Preview
	r.preview.SetColor(c)

	// Alpha
	r.alphaChannel.SetValue(int(r.picker.Alpha))

	// Hex
	r.hex.SetText(colorToString(c))
}
