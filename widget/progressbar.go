package widget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

const defaultText = "%d%%"

type progressRenderer struct {
	widget.BaseRenderer
	bar      *canvas.Rectangle
	label    *canvas.Text
	progress *ProgressBar
}

// MinSize calculates the minimum size of a progress bar.
// This is simply the "100%" label size plus padding.
func (p *progressRenderer) MinSize() fyne.Size {
	var tsize fyne.Size
	if text := p.progress.TextFormatter; text != nil {
		tsize = fyne.MeasureText(text(), p.label.TextSize, p.label.TextStyle)
	} else {
		tsize = fyne.MeasureText("100%", p.label.TextSize, p.label.TextStyle)
	}

	return fyne.NewSize(tsize.Width+theme.Padding()*4, tsize.Height+theme.Padding()*2)
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

	if text := p.progress.TextFormatter; text != nil {
		p.label.Text = text()
	} else {
		p.label.Text = fmt.Sprintf(defaultText, int(ratio*100))
	}

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
	return theme.ShadowColor()
}

func (p *progressRenderer) Refresh() {
	p.applyTheme()
	p.updateBar()

	canvas.Refresh(p.progress.super())
}

// ProgressBar widget creates a horizontal panel that indicates progress
type ProgressBar struct {
	BaseWidget

	Min, Max, Value float64

	// TextFormatter can be used to have a custom format of progress text.
	// If set, it overrides the percentage readout and runs each time the value updates.
	//
	// Since: 1.4
	TextFormatter func() string
}

// SetValue changes the current value of this progress bar (from p.Min to p.Max).
// The widget will be refreshed to indicate the change.
func (p *ProgressBar) SetValue(v float64) {
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
	return &progressRenderer{widget.NewBaseRenderer([]fyne.CanvasObject{bar, label}), bar, label, p}
}

// NewProgressBar creates a new progress bar widget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{Min: 0, Max: 1}

	Renderer(p).Layout(p.MinSize())
	return p
}
