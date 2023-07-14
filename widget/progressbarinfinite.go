package widget

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

const (
	infiniteRefreshRate              = 50 * time.Millisecond
	maxProgressBarInfiniteWidthRatio = 1.0 / 5
	minProgressBarInfiniteWidthRatio = 1.0 / 20
	progressBarInfiniteStepSizeRatio = 1.0 / 50
)

type infProgressRenderer struct {
	widget.BaseRenderer
	background, bar *canvas.Rectangle
	animation       *fyne.Animation
	running         bool
	progress        *ProgressBarInfinite
}

// MinSize calculates the minimum size of a progress bar.
func (p *infProgressRenderer) MinSize() fyne.Size {
	// this is to create the same size infinite progress bar as regular progress bar
	text := fyne.MeasureText("100%", theme.TextSize(), fyne.TextStyle{})

	return fyne.NewSize(text.Width+theme.InnerPadding()*2, text.Height+theme.InnerPadding()*2)
}

func (p *infProgressRenderer) updateBar(done float32) {
	progressWidth := p.progress.size.Width
	spanWidth := progressWidth + (progressWidth * (maxProgressBarInfiniteWidthRatio / 2))
	maxBarWidth := progressWidth * maxProgressBarInfiniteWidthRatio

	barCenterX := spanWidth*done - (spanWidth-progressWidth)/2
	barPos := fyne.NewPos(barCenterX-maxBarWidth/2, 0)
	barSize := fyne.NewSize(maxBarWidth, p.progress.size.Height)
	if barPos.X < 0 {
		barSize.Width += barPos.X
		barPos.X = 0
	}
	if barPos.X+barSize.Width > progressWidth {
		barSize.Width = progressWidth - barPos.X
	}

	p.bar.Resize(barSize)
	p.bar.Move(barPos)
	canvas.Refresh(p.bar)
}

// Layout the components of the infinite progress bar
func (p *infProgressRenderer) Layout(size fyne.Size) {
	p.background.Resize(size)
}

// Refresh updates the size and position of the horizontal scrolling infinite progress bar
func (p *infProgressRenderer) Refresh() {
	if p.isRunning() {
		return // we refresh from the goroutine
	}

	p.background.FillColor = progressBackgroundColor()
	p.background.CornerRadius = theme.InputRadiusSize()
	p.bar.FillColor = theme.PrimaryColor()
	p.bar.CornerRadius = theme.InputRadiusSize()
	p.background.Refresh()
	p.bar.Refresh()
	canvas.Refresh(p.progress.super())
}

func (p *infProgressRenderer) isRunning() bool {
	p.progress.propertyLock.RLock()
	defer p.progress.propertyLock.RUnlock()

	return p.running
}

// Start the infinite progress bar background thread to update it continuously
func (p *infProgressRenderer) start() {
	if p.isRunning() {
		return
	}

	p.progress.propertyLock.Lock()
	defer p.progress.propertyLock.Unlock()
	p.animation = fyne.NewAnimation(time.Second*3, p.updateBar)
	p.animation.Curve = fyne.AnimationLinear
	p.animation.RepeatCount = fyne.AnimationRepeatForever
	p.running = true

	p.animation.Start()
}

// Stop the background thread from updating the infinite progress bar
func (p *infProgressRenderer) stop() {
	p.progress.propertyLock.Lock()
	defer p.progress.propertyLock.Unlock()

	p.running = false
	p.animation.Stop()
}

func (p *infProgressRenderer) Destroy() {
	p.stop()
}

// ProgressBarInfinite widget creates a horizontal panel that indicates waiting indefinitely
// An infinite progress bar loops 0% -> 100% repeatedly until Stop() is called
type ProgressBarInfinite struct {
	BaseWidget
}

// Show this widget, if it was previously hidden
func (p *ProgressBarInfinite) Show() {
	p.Start()
	p.BaseWidget.Show()
}

// Hide this widget, if it was previously visible
func (p *ProgressBarInfinite) Hide() {
	p.Stop()
	p.BaseWidget.Hide()
}

// Start the infinite progress bar animation
func (p *ProgressBarInfinite) Start() {
	cache.Renderer(p).(*infProgressRenderer).start()
}

// Stop the infinite progress bar animation
func (p *ProgressBarInfinite) Stop() {
	cache.Renderer(p).(*infProgressRenderer).stop()
}

// Running returns the current state of the infinite progress animation
func (p *ProgressBarInfinite) Running() bool {
	if !cache.IsRendered(p) {
		return false
	}

	return cache.Renderer(p).(*infProgressRenderer).isRunning()
}

// MinSize returns the size that this widget should not shrink below
func (p *ProgressBarInfinite) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *ProgressBarInfinite) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	background := canvas.NewRectangle(progressBackgroundColor())
	background.CornerRadius = theme.InputRadiusSize()
	bar := canvas.NewRectangle(theme.PrimaryColor())
	bar.CornerRadius = theme.InputRadiusSize()
	render := &infProgressRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{background, bar}),
		background:   background,
		bar:          bar,
		progress:     p,
	}
	render.start()
	return render
}

// NewProgressBarInfinite creates a new progress bar widget that loops indefinitely from 0% -> 100%
// SetValue() is not defined for infinite progress bar
// To stop the looping progress and set the progress bar to 100%, call ProgressBarInfinite.Stop()
func NewProgressBarInfinite() *ProgressBarInfinite {
	p := &ProgressBarInfinite{}
	cache.Renderer(p).Layout(p.MinSize())
	return p
}
