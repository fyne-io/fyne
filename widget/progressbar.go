package widget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const defaultText = "%d%%"

type progressRenderer struct {
	baseRenderer
	bar      *canvas.Rectangle
	label    *canvas.Text
	progress *ProgressBar
}

// MinSize calculates the minimum size of a progress bar.
// This is simply the "100%" label size plus padding.
func (p *progressRenderer) MinSize() fyne.Size {
	text := fyne.MeasureText("100%", p.label.TextSize, p.label.TextStyle)

	return fyne.NewSize(text.Width+theme.Padding()*4, text.Height+theme.Padding()*2)
}

func (p *progressRenderer) updateBar() {
	if p.progress.Value < p.progress.Min {
		p.progress.Value = p.progress.Min
	}
	if p.progress.Value > p.progress.Max {
		p.progress.Value = p.progress.Max
	}

	delta := float32(p.progress.Max - p.progress.Min)
	ratio := float32(p.progress.Value-p.progress.Min) / delta

	p.label.Text = fmt.Sprintf(defaultText, int(ratio*100))

	size := p.progress.Size()
	p.bar.Resize(fyne.NewSize(int(float32(size.Width)*ratio), size.Height))
}

// Layout the components of the check widget
func (p *progressRenderer) Layout(size fyne.Size) {
	p.label.Resize(size)
	p.updateBar()
}

// applyTheme updates the progress bar to match the current theme
func (p *progressRenderer) applyTheme() {
	p.bar.FillColor = theme.PrimaryColor()
	p.label.Color = theme.TextColor()
	p.label.TextSize = theme.TextSize()
}

func (p *progressRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

func (p *progressRenderer) Refresh() {
	p.applyTheme()
	p.updateBar()

	canvas.Refresh(p.progress)
}

// ProgressBar widget creates a horizontal panel that indicates progress
type ProgressBar struct {
	BaseWidget

	Min, Max, Value float64

	maxBind, minBind, valueBind       binding.Float64
	maxNotify, minNotify, valueNotify binding.Notifiable
}

// SetMin changes the current minimum of this progress bar.
// The widget will be refreshed to indicate the change.
func (p *ProgressBar) SetMin(m float64) {
	if p.Min == m {
		return
	}

	p.Min = m

	p.Refresh()
}

// SetMax changes the current maximum of this progress bar.
// The widget will be refreshed to indicate the change.
func (p *ProgressBar) SetMax(m float64) {
	if p.Max == m {
		return
	}

	p.Max = m

	p.Refresh()
}

// SetValue changes the current value of this progress bar (from p.Min to p.Max).
// The widget will be refreshed to indicate the change.
func (p *ProgressBar) SetValue(v float64) {
	if p.Value == v {
		return
	}

	p.Value = v

	p.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (p *ProgressBar) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *ProgressBar) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	if p.Min == 0 && p.Max == 0 {
		p.Max = 1.0
	}

	bar := canvas.NewRectangle(theme.PrimaryColor())
	label := canvas.NewText("0%", theme.TextColor())
	label.Alignment = fyne.TextAlignCenter
	return &progressRenderer{baseRenderer{[]fyne.CanvasObject{bar, label}}, bar, label, p}
}

// BindMin binds the ProgressBar's Min to the given data binding.
// Returns the ProgressBar for chaining.
func (p *ProgressBar) BindMin(data binding.Float64) *ProgressBar {
	p.minBind = data
	p.minNotify = data.AddFloat64Listener(p.SetMin)
	return p
}

// UnbindMin unbinds the ProgressBar's Min from the data binding (if any).
// Returns the ProgressBar for chaining.
func (p *ProgressBar) UnbindMin() *ProgressBar {
	if p.minBind != nil {
		p.minBind.DeleteListener(p.minNotify)
	}
	p.minBind = nil
	p.minNotify = nil
	return p
}

// BindMax binds the ProgressBar's Max to the given data binding.
// Returns the ProgressBar for chaining.
func (p *ProgressBar) BindMax(data binding.Float64) *ProgressBar {
	p.maxBind = data
	p.maxNotify = data.AddFloat64Listener(p.SetMax)
	return p
}

// UnbindMax unbinds the ProgressBar's Max from the data binding (if any).
// Returns the ProgressBar for chaining.
func (p *ProgressBar) UnbindMax() *ProgressBar {
	if p.maxBind != nil {
		p.maxBind.DeleteListener(p.maxNotify)
	}
	p.maxBind = nil
	p.maxNotify = nil
	return p
}

// BindValue binds the ProgressBar's Value to the given data binding.
// Returns the ProgressBar for chaining.
func (p *ProgressBar) BindValue(data binding.Float64) *ProgressBar {
	p.valueBind = data
	p.valueNotify = data.AddFloat64Listener(p.SetValue)
	return p
}

// UnbindValue unbinds the ProgressBar's Value from the data binding (if any).
// Returns the ProgressBar for chaining.
func (p *ProgressBar) UnbindValue() *ProgressBar {
	if p.valueBind != nil {
		p.valueBind.DeleteListener(p.valueNotify)
	}
	p.valueBind = nil
	p.valueNotify = nil
	return p
}

// NewProgressBar creates a new progress bar widget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{Min: 0, Max: 1}

	Renderer(p).Layout(p.MinSize())
	return p
}
