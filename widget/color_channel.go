package widget

import (
	"fmt"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*ColorChannel)(nil)

// ColorChannel controls a channel of a color and triggers the callback when changed.
type ColorChannel struct {
	BaseWidget
	name      string
	value     float64
	onChanged func(float64)
}

// NewColorChannel returns a new color channel control for the channel with the given name.
func NewColorChannel(name string, value float64, onChanged func(float64)) *ColorChannel {
	c := &ColorChannel{
		name:      name,
		value:     colorClamp(value),
		onChanged: onChanged,
	}
	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (c *ColorChannel) CreateRenderer() fyne.WidgetRenderer {
	label := NewLabel(c.name)
	entry := &Entry{
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
	slider := &Slider{
		Value:       0,
		Min:         0,
		Max:         1.0,
		Step:        0.000001,
		Orientation: Horizontal,
		OnChanged: func(value float64) {
			c.SetValue(value)
		},
	}
	contents := NewVBox(
		NewHBox(
			label,
			layout.NewSpacer(),
			entry,
		),
		slider,
	)
	r := &colorChannelRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{contents}),
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
func (c *ColorChannel) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetValue updates the value in this color widget
func (c *ColorChannel) SetValue(value float64) {
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
	widget.BaseRenderer
	control  *ColorChannel
	label    *Label
	entry    *Entry
	slider   *Slider
	contents fyne.CanvasObject
}

func (r *colorChannelRenderer) Layout(size fyne.Size) {
	r.contents.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	r.contents.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height-2*theme.Padding()))
}

func (r *colorChannelRenderer) MinSize() (min fyne.Size) {
	min = r.contents.MinSize()
	min = min.Add(fyne.NewSize(2*theme.Padding(), 2*theme.Padding()))
	return
}

func (r *colorChannelRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.control.Size())
	canvas.Refresh(r.control.super())
}

func (r *colorChannelRenderer) updateObjects() {
	r.entry.SetText(fmt.Sprintf("%.6f", r.control.value))
	r.slider.Value = r.control.value
	r.slider.Refresh()
}
