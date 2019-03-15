package widget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const defaultText = "%d%%"

type progressRenderer struct {
	objects []fyne.CanvasObject

	bar   *canvas.Rectangle
	label *canvas.Text

	progress *ProgressBar
}

// MinSize calculates the minimum size of a progress bar.
// This is simply the "100%" label size plus padding.
func (p *progressRenderer) MinSize() fyne.Size {
	text := textMinSize("100%", p.label.TextSize, p.label.TextStyle)

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
	ratio := float32(p.progress.Value-p.progress.Min) / float32(delta)

	p.label.Text = fmt.Sprintf(defaultText, int(ratio*100))

	size := p.progress.Size()
	p.bar.Resize(fyne.NewSize(int(float32(size.Width)*ratio), size.Height))
}

// Layout the components of the check widget
func (p *progressRenderer) Layout(size fyne.Size) {
	p.label.Resize(size)
	p.updateBar()
}

// ApplyTheme is called when the progress bar may need to update it's look
func (p *progressRenderer) ApplyTheme() {
	p.bar.FillColor = theme.PrimaryColor()
	p.label.Color = theme.TextColor()

	p.Refresh()
}

func (p *progressRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

func (p *progressRenderer) Refresh() {
	p.updateBar()

	canvas.Refresh(p.progress)
}

func (p *progressRenderer) Objects() []fyne.CanvasObject {
	return p.objects
}

func (p *progressRenderer) Destroy() {
}

// ProgressBar widget creates a horizontal panel that indicates progress
type ProgressBar struct {
	baseWidget

	Min, Max, Value float64
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (p *ProgressBar) Resize(size fyne.Size) {
	p.resize(size, p)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (p *ProgressBar) Move(pos fyne.Position) {
	p.move(pos, p)
}

// MinSize returns the smallest size this widget can shrink to
func (p *ProgressBar) MinSize() fyne.Size {
	return p.minSize(p)
}

// Show this widget, if it was previously hidden
func (p *ProgressBar) Show() {
	p.show(p)
}

// Hide this widget, if it was previously visible
func (p *ProgressBar) Hide() {
	p.hide(p)
}

// SetValue changes the current value of this progress bar (from p.Min to p.Max).
// The widget will be refreshed to indicate the change.
func (p *ProgressBar) SetValue(v float64) {
	p.Value = v
	Renderer(p).Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (p *ProgressBar) CreateRenderer() fyne.WidgetRenderer {
	if p.Min == 0 && p.Max == 0 {
		p.Max = 1.0
	}

	bar := canvas.NewRectangle(theme.PrimaryColor())
	label := canvas.NewText("0%", theme.TextColor())
	label.Alignment = fyne.TextAlignCenter
	return &progressRenderer{[]fyne.CanvasObject{bar, label}, bar, label, p}
}

// NewProgressBar creates a new progress bar widget.
// The default Min is 0 and Max is 1, Values set should be between those numbers.
// The display will convert this to a percentage.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{Min: 0, Max: 1}

	Renderer(p).Layout(p.MinSize())
	return p
}
