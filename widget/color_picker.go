package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// NewColorBasicPicker returns a component for selecting basic colors.
func NewColorBasicPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors([]string{
		"#f44336", //Red
		"#ff9800", //Orange
		"#ffeb3b", //Yellow
		"#4caf50", //Green
		"#2196f3", //Blue
		"#3f51b5", //Indigo
		"#9c27b0", //Purple
		"#795548", //Brown
	}...), theme.ColorChromaticIcon(), callback)
}

// NewColorGreyscalePicker returns a component for selecting greyscale colors.
func NewColorGreyscalePicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors([]string{
		"#ffffff",
		"#dddddd",
		"#bbbbbb",
		"#999999",
		"#666666",
		"#444444",
		"#222222",
		"#000000",
	}...), theme.ColorAchromaticIcon(), callback)
}

// NewColorRecentPicker returns a component for selecting recent colors.
func NewColorRecentPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors(readRecentColors()...), theme.HistoryIcon(), callback)
}

var _ fyne.Widget = (*ColorAdvancedPicker)(nil)

// ColorAdvancedPicker widget is a component for selecting a color.
type ColorAdvancedPicker struct {
	BaseWidget
	Red, Green, Blue, Alpha    float64
	Hue, Saturation, Lightness float64
	ColorModel                 string

	onChange func(color.Color)
}

// NewColorAdvancedPicker returns a new color widget set to the given color.
func NewColorAdvancedPicker(color color.Color, onChange func(color.Color)) *ColorAdvancedPicker {
	c := &ColorAdvancedPicker{
		onChange: onChange,
	}
	c.ExtendBaseWidget(c)
	c.updateColor(color)
	return c
}

// Color returns the currently selected color.
func (p *ColorAdvancedPicker) Color() color.Color {
	return &color.NRGBA{
		uint8(p.Red * 255),
		uint8(p.Green * 255),
		uint8(p.Blue * 255),
		uint8(p.Alpha * 255),
	}
}

// SetColor updates the color selected in this color widget.
func (p *ColorAdvancedPicker) SetColor(color color.Color) {
	if p.updateColor(color) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(color)
		}
	}
}

func (p *ColorAdvancedPicker) updateColor(color color.Color) bool {
	r, g, b, a := colorToRGBA(color)
	if p.Red == r && p.Green == g && p.Blue == b && p.Alpha == a {
		return false
	}
	p.Red = r
	p.Green = g
	p.Blue = b
	p.Alpha = a
	p.rgbToHsl()
	return true
}

// RGBA return the Red, Green, Blue, and Alpha components of the currently selected color.
func (p *ColorAdvancedPicker) RGBA() (float64, float64, float64, float64) {
	return p.Red, p.Green, p.Blue, p.Alpha
}

// SetRGBA updated the Red, Green, Blue, and Alpha components of the currently selected color.
func (p *ColorAdvancedPicker) SetRGBA(r, g, b, a float64) {
	r = colorClamp(r)
	g = colorClamp(g)
	b = colorClamp(b)
	a = colorClamp(a)
	if p.Red == r && p.Green == g && p.Blue == b && p.Alpha == a {
		return
	}
	p.Red = r
	p.Green = g
	p.Blue = b
	p.Alpha = a
	p.rgbToHsl()
	p.Refresh()
	if f := p.onChange; f != nil {
		f(p.Color())
	}
}

// HSLA return the Hue, Saturation, Lightness, and Alpha components of the currently selected color.
func (p *ColorAdvancedPicker) HSLA() (float64, float64, float64, float64) {
	return p.Hue, p.Saturation, p.Lightness, p.Alpha
}

// SetHSLA updated the Hue, Saturation, Lightness, and Alpha components of the currently selected color.
func (p *ColorAdvancedPicker) SetHSLA(h, s, l, a float64) {
	h = colorClamp(h)
	s = colorClamp(s)
	l = colorClamp(l)
	a = colorClamp(a)
	if p.Hue == h && p.Saturation == s && p.Lightness == l && p.Alpha == a {
		return
	}
	p.Hue = h
	p.Saturation = s
	p.Lightness = l
	p.Alpha = a
	p.hslToRgb()
	p.Refresh()
	if f := p.onChange; f != nil {
		f(p.Color())
	}
}

// MinSize returns the size that this widget should not shrink below.
func (p *ColorAdvancedPicker) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (p *ColorAdvancedPicker) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)

	p.propertyLock.RLock()
	defer p.propertyLock.RUnlock()

	// RGB
	rgbArea := NewRGBAColorArea(p)
	rgbArea.SetMinSize(fyne.NewSize(128, 128))
	redChannel := NewColorChannel("R:", p.Red, func(r float64) {
		p.SetRGBA(r, p.Green, p.Blue, p.Alpha)
	})
	greenChannel := NewColorChannel("G:", p.Green, func(g float64) {
		p.SetRGBA(p.Red, g, p.Blue, p.Alpha)
	})
	blueChannel := NewColorChannel("B:", p.Blue, func(b float64) {
		p.SetRGBA(p.Red, p.Green, b, p.Alpha)
	})
	rgbTab := NewTabItem(
		"RGB",
		NewHBox(
			rgbArea,
			NewVBox(
				redChannel,
				greenChannel,
				blueChannel,
			),
		),
	)

	// HSL
	hslArea := NewHSLAColorArea(p)
	hslArea.SetMinSize(fyne.NewSize(128, 128))
	hueChannel := NewColorChannel("H:", p.Hue, func(h float64) {
		p.SetHSLA(h, p.Saturation, p.Lightness, p.Alpha)
	})
	saturationChannel := NewColorChannel("S:", p.Saturation, func(s float64) {
		p.SetHSLA(p.Hue, s, p.Lightness, p.Alpha)
	})
	lightnessChannel := NewColorChannel("L:", p.Lightness, func(l float64) {
		p.SetHSLA(p.Hue, p.Saturation, l, p.Alpha)
	})
	hslTab := NewTabItem(
		"HSL",
		NewHBox(
			hslArea,
			NewVBox(
				hueChannel,
				saturationChannel,
				lightnessChannel,
			),
		),
	)

	modelTabContainer := NewTabContainer(
		rgbTab,
		hslTab,
	)
	modelTabContainer.OnChanged = func(item *TabItem) {
		p.ColorModel = item.Text
	}

	// Preview
	preview := &canvas.Rectangle{}
	preview.SetMinSize(fyne.NewSize(128, 128))

	// Alpha
	alphaChannel := NewColorChannel("A:", p.Alpha, func(a float64) {
		p.SetRGBA(p.Red, p.Green, p.Blue, a)
	})

	// TODO Hex color code entry

	contents := NewVBox(
		modelTabContainer,
		NewHBox(
			fyne.NewContainerWithLayout(layout.NewPaddedLayout(), newCheckeredBackground(), preview),
			NewVBox(
				alphaChannel,
				// TODO Hex color code entry
			),
		),
	)
	r := &colorPickerRenderer{
		BaseRenderer:      widget.NewBaseRenderer([]fyne.CanvasObject{contents}),
		picker:            p,
		modelTabContainer: modelTabContainer,
		rgbTab:            rgbTab,
		rgbArea:           rgbArea,
		redChannel:        redChannel,
		greenChannel:      greenChannel,
		blueChannel:       blueChannel,
		hslTab:            hslTab,
		hslArea:           hslArea,
		hueChannel:        hueChannel,
		saturationChannel: saturationChannel,
		lightnessChannel:  lightnessChannel,
		preview:           preview,
		alphaChannel:      alphaChannel,
		contents:          contents,
	}
	r.updateObjects()
	return r
}

func (p *ColorAdvancedPicker) rgbToHsl() {
	p.Hue, p.Saturation, p.Lightness = rgbToHsl(p.Red, p.Green, p.Blue)
}

func (p *ColorAdvancedPicker) hslToRgb() {
	p.Red, p.Green, p.Blue = hslToRgb(p.Hue, p.Saturation, p.Lightness)
}

var _ fyne.WidgetRenderer = (*colorPickerRenderer)(nil)

type colorPickerRenderer struct {
	widget.BaseRenderer
	picker            *ColorAdvancedPicker
	modelTabContainer *TabContainer
	rgbTab            *TabItem
	rgbArea           *ColorArea
	redChannel        *ColorChannel
	greenChannel      *ColorChannel
	blueChannel       *ColorChannel
	hslTab            *TabItem
	hslArea           *ColorArea
	hueChannel        *ColorChannel
	saturationChannel *ColorChannel
	lightnessChannel  *ColorChannel
	preview           *canvas.Rectangle
	alphaChannel      *ColorChannel
	contents          fyne.CanvasObject
}

func (r *colorPickerRenderer) Layout(size fyne.Size) {
	r.contents.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	r.contents.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height-2*theme.Padding()))
}

func (r *colorPickerRenderer) MinSize() (min fyne.Size) {
	min = r.contents.MinSize()
	min = min.Add(fyne.NewSize(2*theme.Padding(), 2*theme.Padding()))
	return
}

func (r *colorPickerRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.picker.Size())
	canvas.Refresh(r.picker.super())
}

func (r *colorPickerRenderer) updateObjects() {
	if r.picker.ColorModel != r.modelTabContainer.CurrentTab().Text {
		switch r.picker.ColorModel {
		case "hsl":
			r.modelTabContainer.SelectTab(r.hslTab)
		case "rgb":
			r.modelTabContainer.SelectTab(r.rgbTab)
		}
	}

	// RGB
	r.rgbArea.Refresh()
	r.redChannel.SetValue(r.picker.Red)
	r.greenChannel.SetValue(r.picker.Green)
	r.blueChannel.SetValue(r.picker.Blue)

	// HSL
	r.hslArea.Refresh()
	r.hueChannel.SetValue(r.picker.Hue)
	r.saturationChannel.SetValue(r.picker.Saturation)
	r.lightnessChannel.SetValue(r.picker.Lightness)

	// Preview
	r.preview.FillColor = r.picker.Color()
	r.preview.Refresh()

	// Alpha
	r.alphaChannel.SetValue(r.picker.Alpha)

	// TODO Hex color code entry
}
