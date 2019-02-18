package widget

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

const infiniteRefreshRate = 50 * time.Millisecond

type infProgressRenderer struct {
	objects  []fyne.CanvasObject
	bar      *canvas.Rectangle
	progress *InfProgressBar
}

// MinSize calculates the minimum size of a progress bar.
func (p *infProgressRenderer) MinSize() fyne.Size {
	// this is to create the same size infinite progress bar as regular progress bar
	text := textMinSize("100%", theme.TextSize(), fyne.TextStyle{})

	return fyne.NewSize(text.Width+theme.Padding()*4, text.Height+theme.Padding()*2)
}

func (p *infProgressRenderer) updateBar() {
	barSize := p.bar.Size()
	progressSize := p.progress.Size()
	barWidth := barSize.Width
	maxBarWidth := progressSize.Width / 5 // 1/5 of bar
	barPos := p.bar.Position()
	progPos := p.progress.Position()

	stepSize := (int)(progressSize.Width / 50)

	// check to make sure inner bar is sized correctly
	if barWidth < maxBarWidth || barSize.Height < progressSize.Height {
		// set width
		newBoxSize := fyne.Size{Width: maxBarWidth, Height: progressSize.Height}
		p.bar.Resize(newBoxSize)
	}

	// move bar to the right by stepSize
	barPos.X += stepSize

	// loop around to start when bar goes to end
	if barWidth+barPos.X >= progPos.X+progressSize.Width {
		barPos.X = 0
	}

	p.bar.Move(fyne.Position{X: barPos.X, Y: barPos.Y})
}

// Layout the components of the progress bar widget
func (p *infProgressRenderer) Layout(size fyne.Size) {
	// set height of progress bar
	p.updateBar()
}

// ApplyTheme is called when the progress bar may need to update it's look
func (p *infProgressRenderer) ApplyTheme() {
	p.bar.FillColor = theme.PrimaryColor()

	p.Refresh()
}

func (p *infProgressRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

func (p *infProgressRenderer) Refresh() {
	p.updateBar()
	canvas.Refresh(p.progress)
}

func (p *infProgressRenderer) Objects() []fyne.CanvasObject {
	return p.objects
}

// InfProgressBar widget creates a horizontal panel that indicates waiting indefinitely
// An infinite progress bar loops 0% -> 100% until Stop() is called
type InfProgressBar struct {
	baseWidget

	ticker *time.Ticker
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (p *InfProgressBar) Resize(size fyne.Size) {
	p.resize(size, p)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (p *InfProgressBar) Move(pos fyne.Position) {
	p.move(pos, p)
}

// MinSize returns the smallest size this widget can shrink to
func (p *InfProgressBar) MinSize() fyne.Size {
	return p.minSize(p)
}

// Show this widget, if it was previously hidden
func (p *InfProgressBar) Show() {
	p.Start()
	p.show(p)
}

// Hide this widget, if it was previously visible
func (p *InfProgressBar) Hide() {
	p.Stop()
	p.hide(p)
}

// Start the infinite progress bar background thread to update it continuously
func (p *InfProgressBar) Start() {
	if p.ticker == nil {
		go p.infiniteProgressLoop()
	}
}

// Stop the infinite progress goroutine and sets value to the Max
func (p *InfProgressBar) Stop() {
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}
}

// internal loop called with `go infiniteProgressLoop()`
// updates the infinite-style progress bar
// can be exited by calling ProgressBar.StopInfiniteProgress()
func (p *InfProgressBar) infiniteProgressLoop() {
	defer p.Stop()
	p.ticker = time.NewTicker(infiniteRefreshRate)

	for range p.ticker.C {
		Renderer(p).Refresh()
	}
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (p *InfProgressBar) CreateRenderer() fyne.WidgetRenderer {
	bar := canvas.NewRectangle(theme.PrimaryColor())

	return &infProgressRenderer{[]fyne.CanvasObject{bar}, bar, p}
}

// NewInfiniteProgressBar creates a new progress bar widget that loops indefinitely from 0% -> 100%
// SetValue() should not be called when using an infinite progress bar
// To stop the looping progress and set the progress bar to 100%, call ProgressBar.StopInfiniteProgress()
func NewInfiniteProgressBar() *InfProgressBar {
	p := &InfProgressBar{}
	Renderer(p).Layout(p.MinSize())
	return p
}
