package widget

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/internal/cache"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

type progressRenderer struct {
	widget.BaseRenderer
	background, bar *canvas.Rectangle
	label           *canvas.Text
	progress        *ProgressBar
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

	padding := theme.InnerPadding() * 2
	return fyne.NewSize(tsize.Width+padding, tsize.Height+padding)
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
		p.label.Text = strconv.Itoa(int(ratio*100)) + "%"
	}

	size := p.progress.Size()
	p.bar.Resize(fyne.NewSize(size.Width*ratio, size.Height))
}

// Layout the components of the check widget
func (p *progressRenderer) Layout(size fyne.Size) {
	p.background.Resize(size)
	p.label.Resize(size)
	p.updateBar()
}

// applyTheme updates the progress bar to match the current theme
func (p *progressRenderer) applyTheme() {
	p.background.FillColor = progressBackgroundColor()
	p.background.CornerRadius = theme.InputRadiusSize()
	p.bar.FillColor = theme.PrimaryColor()
	p.bar.CornerRadius = theme.InputRadiusSize()
	p.label.Color = theme.BackgroundColor()
	p.label.TextSize = theme.TextSize()
}

func (p *progressRenderer) Refresh() {
	p.applyTheme()
	p.updateBar()
	p.background.Refresh()
	p.bar.Refresh()
	p.label.Refresh()
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
	TextFormatter func() string `json:"-"`

	binder basicBinder
}

// Bind connects the specified data source to this ProgressBar.
// The current value will be displayed and any changes in the data will cause the widget to update.
//
// Since: 2.0
func (p *ProgressBar) Bind(data binding.Float) {
	p.binder.SetCallback(p.updateFromData)
	p.binder.Bind(data)
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

	background := canvas.NewRectangle(progressBackgroundColor())
	background.CornerRadius = theme.InputRadiusSize()
	bar := canvas.NewRectangle(theme.PrimaryColor())
	bar.CornerRadius = theme.InputRadiusSize()
	label := canvas.NewText("0%", theme.BackgroundColor())
	label.Alignment = fyne.TextAlignCenter
	return &progressRenderer{widget.NewBaseRenderer([]fyne.CanvasObject{background, bar, label}), background, bar, label, p}
}

// Unbind disconnects any configured data source from this ProgressBar.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (p *ProgressBar) Unbind() {
	p.binder.Unbind()
}

// NewProgressBar creates a new progress bar widget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{Min: 0, Max: 1}

	cache.Renderer(p).Layout(p.MinSize())
	return p
}

// NewProgressBarWithData returns a progress bar connected with the specified data source.
//
// Since: 2.0
func NewProgressBarWithData(data binding.Float) *ProgressBar {
	p := NewProgressBar()
	p.Bind(data)

	return p
}

func progressBackgroundColor() color.Color {
	r, g, b, a := col.ToNRGBA(theme.PrimaryColor())
	faded := uint8(a) / 2
	return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: faded}
}

func (p *ProgressBar) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatSource, ok := data.(binding.Float)
	if !ok {
		return
	}

	val, err := floatSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	p.SetValue(val)
}
