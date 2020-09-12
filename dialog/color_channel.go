package dialog

import (
	"fmt"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	internalwidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var _ fyne.Widget = (*colorChannel)(nil)

// colorChannel controls a channel of a color and triggers the callback when changed.
type colorChannel struct {
	widget.BaseWidget
	name      string
	value     float64
	onChanged func(float64)
}

// newColorChannel returns a new color channel control for the channel with the given name.
func newColorChannel(name string, value float64, onChanged func(float64)) *colorChannel {
	c := &colorChannel{
		name:      name,
		value:     colorClamp(value),
		onChanged: onChanged,
	}
	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *colorChannel) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(c.name)
	entry := &widget.Entry{
		Text: "0.00",
		OnChanged: func(text string) {
			var value float64
			n, err := fmt.Sscanf(text, "%f", &value)
			if err != nil {
				fyne.LogError("Couldn't parse value", err)
			} else if n == 1 {
				c.SetValue(value)
			}
		},
		// TODO add number 0.0-1.0 validator
	}
	slider := &widget.Slider{
		Value:       0,
		Min:         0,
		Max:         1.0,
		Step:        0.000001,
		Orientation: widget.Horizontal,
		OnChanged: func(value float64) {
			c.SetValue(value)
		},
	}
	contents := widget.NewVBox(
		widget.NewHBox(
			label,
			layout.NewSpacer(),
			entry,
		),
		slider,
	)
	r := &colorChannelRenderer{
		BaseRenderer: internalwidget.NewBaseRenderer([]fyne.CanvasObject{contents}),
		control:      c,
		label:        label,
		entry:        entry,
		slider:       slider,
		contents:     contents,
	}
	r.updateObjects()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (c *colorChannel) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetValue updates the value in this color widget
func (c *colorChannel) SetValue(value float64) {
	value = colorClamp(value)
	if math.Abs(c.value-value) < 0.000001 {
		return
	}
	c.value = value
	c.Refresh()
	if f := c.onChanged; f != nil {
		f(value)
	}
}

type colorChannelRenderer struct {
	internalwidget.BaseRenderer
	control  *colorChannel
	label    *widget.Label
	entry    *widget.Entry
	slider   *widget.Slider
	contents fyne.CanvasObject
}

func (r *colorChannelRenderer) Layout(size fyne.Size) {
	r.contents.Move(fyne.NewPos(0, 0))
	r.contents.Resize(fyne.NewSize(size.Width, size.Height))
}

func (r *colorChannelRenderer) MinSize() fyne.Size {
	return r.contents.MinSize()
}

func (r *colorChannelRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.control.Size())
	canvas.Refresh(r.control)
}

func (r *colorChannelRenderer) updateObjects() {
	r.entry.SetText(fmt.Sprintf("%.6f", r.control.value))
	r.slider.Value = r.control.value
	r.slider.Refresh()
}
