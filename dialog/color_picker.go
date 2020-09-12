package dialog

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	internalwidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// newColorBasicPicker returns a component for selecting basic colors.
func newColorBasicPicker(callback func(color.Color)) fyne.CanvasObject {
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

// newColorGreyscalePicker returns a component for selecting greyscale colors.
func newColorGreyscalePicker(callback func(color.Color)) fyne.CanvasObject {
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

// newColorRecentPicker returns a component for selecting recent colors.
func newColorRecentPicker(callback func(color.Color)) fyne.CanvasObject {
	return newColorButtonBox(stringsToColors(readRecentColors()...), theme.HistoryIcon(), callback)
}

var _ fyne.Widget = (*colorAdvancedPicker)(nil)

// colorAdvancedPicker widget is a component for selecting a color.
type colorAdvancedPicker struct {
	widget.BaseWidget
	Red, Green, Blue, Alpha    float64
	Hue, Saturation, Lightness float64
	ColorModel                 string

	onChange func(color.Color)
}

// newColorAdvancedPicker returns a new color widget set to the given color.
func newColorAdvancedPicker(color color.Color, onChange func(color.Color)) *colorAdvancedPicker {
	c := &colorAdvancedPicker{
		onChange: onChange,
	}
	c.ExtendBaseWidget(c)
	c.updateColor(color)
	return c
}

// Color returns the currently selected color.
func (p *colorAdvancedPicker) Color() color.Color {
	return &color.NRGBA{
		uint8(p.Red * 255),
		uint8(p.Green * 255),
		uint8(p.Blue * 255),
		uint8(p.Alpha * 255),
	}
}

// SetColor updates the color selected in this color widget.
func (p *colorAdvancedPicker) SetColor(color color.Color) {
	if p.updateColor(color) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(color)
		}
	}
}

// SetHSLA updated the Hue, Saturation, Lightness, and Alpha components of the currently selected color.
func (p *colorAdvancedPicker) SetHSLA(h, s, l, a float64) {
	if p.updateHSLA(h, s, l, a) {
		p.Refresh()
		if f := p.onChange; f != nil {
			f(p.Color())
		}
	}
}

// SetRGBA updated the Red, Green, Blue, and Alpha components of the currently selected color.
func (p *colorAdvancedPicker) SetRGBA(r, g, b, a float64) {
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
	preview := &canvas.Rectangle{}

	// HSL
	hueChannel := newColorChannel("H:", p.Hue, func(h float64) {
		p.SetHSLA(h, p.Saturation, p.Lightness, p.Alpha)
	})
	saturationChannel := newColorChannel("S:", p.Saturation, func(s float64) {
		p.SetHSLA(p.Hue, s, p.Lightness, p.Alpha)
	})
	lightnessChannel := newColorChannel("L:", p.Lightness, func(l float64) {
		p.SetHSLA(p.Hue, p.Saturation, l, p.Alpha)
	})
	hslBox := widget.NewVBox(
		hueChannel,
		saturationChannel,
		lightnessChannel,
	)
	hslTab := widget.NewTabItem(
		"HSL",
		hslBox,
	)

	// RGB
	redChannel := newColorChannel("R:", p.Red, func(r float64) {
		p.SetRGBA(r, p.Green, p.Blue, p.Alpha)
	})
	greenChannel := newColorChannel("G:", p.Green, func(g float64) {
		p.SetRGBA(p.Red, g, p.Blue, p.Alpha)
	})
	blueChannel := newColorChannel("B:", p.Blue, func(b float64) {
		p.SetRGBA(p.Red, p.Green, b, p.Alpha)
	})
	rgbBox := widget.NewVBox(
		redChannel,
		greenChannel,
		blueChannel,
	)
	rgbTab := widget.NewTabItem(
		"RGB",
		rgbBox,
	)

	modelTabContainer := widget.NewTabContainer(
		hslTab,
		rgbTab,
	)
	modelTabContainer.OnChanged = func(item *widget.TabItem) {
		p.ColorModel = item.Text
	}

	// Wheel
	wheel := newColorWheel(func(hue, saturation, lightness, alpha float64) {
		p.SetHSLA(hue, saturation, lightness, alpha)
	})

	// Alpha
	alphaChannel := newColorChannel("A:", p.Alpha, func(a float64) {
		p.SetRGBA(p.Red, p.Green, p.Blue, a)
	})

	hex := &widget.Entry{
		OnChanged: func(text string) {
			c, err := stringToColor(text)
			if err != nil {
				// TODO trigger entry invalid state
			} else {
				p.SetColor(c)
			}
		},
	}

	contents := widget.NewVBox(
		layout.NewGridContainerWithColumns(3,
			fyne.NewContainerWithLayout(layout.NewPaddedLayout(), wheel),
			hslBox,
			rgbBox),
		layout.NewGridContainerWithColumns(3,
			fyne.NewContainerWithLayout(layout.NewPaddedLayout(), preview),
			fyne.NewContainerWithLayout(layout.NewCenterLayout(), hex),
			alphaChannel,
		),
	)
	/*
		contents := widget.NewHBox(
			widget.NewVBox(
				fyne.NewContainerWithLayout(layout.NewPaddedLayout(), wheel),
				fyne.NewContainerWithLayout(layout.NewPaddedLayout(), newCheckeredBackground(), preview),
				fyne.NewContainerWithLayout(layout.NewPaddedLayout(), hex),
			),
			widget.NewVBox(
				modelTabContainer,
				alphaChannel,
			),
		)
	*/
	r := &colorPickerRenderer{
		BaseRenderer:      internalwidget.NewBaseRenderer([]fyne.CanvasObject{contents}),
		picker:            p,
		modelTabContainer: modelTabContainer,
		rgbTab:            rgbTab,
		redChannel:        redChannel,
		greenChannel:      greenChannel,
		blueChannel:       blueChannel,
		hslTab:            hslTab,
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
	r, g, b, a := colorToRGBA(color)
	if p.Red == r && p.Green == g && p.Blue == b && p.Alpha == a {
		return false
	}
	return p.updateRGBA(r, g, b, a)
}

func (p *colorAdvancedPicker) updateHSLA(h, s, l, a float64) bool {
	h = colorClamp(h)
	s = colorClamp(s)
	l = colorClamp(l)
	a = colorClamp(a)
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

func (p *colorAdvancedPicker) updateRGBA(r, g, b, a float64) bool {
	r = colorClamp(r)
	g = colorClamp(g)
	b = colorClamp(b)
	a = colorClamp(a)
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
	internalwidget.BaseRenderer
	picker            *colorAdvancedPicker
	modelTabContainer *widget.TabContainer
	rgbTab            *widget.TabItem
	redChannel        *colorChannel
	greenChannel      *colorChannel
	blueChannel       *colorChannel
	hslTab            *widget.TabItem
	hueChannel        *colorChannel
	saturationChannel *colorChannel
	lightnessChannel  *colorChannel
	wheel             *colorWheel
	preview           *canvas.Rectangle
	alphaChannel      *colorChannel
	hex               *widget.Entry
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
	canvas.Refresh(r.picker)
}

func (r *colorPickerRenderer) updateObjects() {
	if r.picker.ColorModel != r.modelTabContainer.CurrentTab().Text {
		switch r.picker.ColorModel {
		case "HSL":
			r.modelTabContainer.SelectTab(r.hslTab)
		case "RGB":
			r.modelTabContainer.SelectTab(r.rgbTab)
		}
	}

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
	r.preview.FillColor = color
	r.preview.Refresh()

	// Alpha
	r.alphaChannel.SetValue(r.picker.Alpha)

	r.hex.SetText(colorToString(color))
}
