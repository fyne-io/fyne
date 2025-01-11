package widget

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	background, bar canvas.Rectangle
	animation       fyne.Animation
	wasRunning      bool
	progress        *ProgressBarInfinite
}

// MinSize calculates the minimum size of a progress bar.
func (p *infProgressRenderer) MinSize() fyne.Size {
	th := p.progress.Theme()
	innerPad2 := th.Size(theme.SizeNameInnerPadding) * 2
	// this is to create the same size infinite progress bar as regular progress bar
	text := fyne.MeasureText("100%", th.Size(theme.SizeNameText), fyne.TextStyle{})

	return fyne.NewSize(text.Width+innerPad2, text.Height+innerPad2)
}

func (p *infProgressRenderer) updateBar(done float32) {
	size := p.progress.Size()
	progressWidth := size.Width
	spanWidth := progressWidth + (progressWidth * (maxProgressBarInfiniteWidthRatio / 2))
	maxBarWidth := progressWidth * maxProgressBarInfiniteWidthRatio

	barCenterX := spanWidth*done - (spanWidth-progressWidth)/2
	barPos := fyne.NewPos(barCenterX-maxBarWidth/2, 0)
	barSize := fyne.NewSize(maxBarWidth, size.Height)
	if barPos.X < 0 {
		barSize.Width += barPos.X
		barPos.X = 0
	}
	if barPos.X+barSize.Width > progressWidth {
		barSize.Width = progressWidth - barPos.X
	}

	p.bar.Resize(barSize)
	p.bar.Move(barPos)
	canvas.Refresh(&p.bar)
}

// Layout the components of the infinite progress bar
func (p *infProgressRenderer) Layout(size fyne.Size) {
	p.background.Resize(size)
}

// Refresh updates the size and position of the horizontal scrolling infinite progress bar
func (p *infProgressRenderer) Refresh() {
	running := p.progress.Running()
	if running {
		if !p.wasRunning {
			p.start()
		}
		return // we refresh from the goroutine
	} else if p.wasRunning {
		p.stop()
		return
	}

	th := p.progress.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	cornerRadius := th.Size(theme.SizeNameInputRadius)
	primaryColor := th.Color(theme.ColorNamePrimary, v)

	p.background.FillColor = progressBlendColor(primaryColor)
	p.background.CornerRadius = cornerRadius
	p.bar.FillColor = primaryColor
	p.bar.CornerRadius = cornerRadius
	p.background.Refresh()
	p.bar.Refresh()
	canvas.Refresh(p.progress.super())
}

// Start the infinite progress bar background thread to update it continuously
func (p *infProgressRenderer) start() {
	p.animation.Duration = time.Second * 3
	p.animation.Tick = p.updateBar
	p.animation.Curve = fyne.AnimationLinear
	p.animation.RepeatCount = fyne.AnimationRepeatForever

	p.wasRunning = true
	p.animation.Start()
}

// Stop the background thread from updating the infinite progress bar
func (p *infProgressRenderer) stop() {
	p.wasRunning = false
	p.animation.Stop()
}

func (p *infProgressRenderer) Destroy() {
	p.progress.running = false

	p.stop()
}

// ProgressBarInfinite widget creates a horizontal panel that indicates waiting indefinitely
// An infinite progress bar loops 0% -> 100% repeatedly until Stop() is called
type ProgressBarInfinite struct {
	BaseWidget
	running bool
}

// Show this widget, if it was previously hidden
func (p *ProgressBarInfinite) Show() {
	p.running = true

	p.BaseWidget.Show()
}

// Hide this widget, if it was previously visible
func (p *ProgressBarInfinite) Hide() {
	p.running = false

	p.BaseWidget.Hide()
}

// Start the infinite progress bar animation
func (p *ProgressBarInfinite) Start() {
	if p.running {
		return
	}

	p.running = true
	p.BaseWidget.Refresh()
}

// Stop the infinite progress bar animation
func (p *ProgressBarInfinite) Stop() {
	if !p.running {
		return
	}

	p.running = false
	p.BaseWidget.Refresh()
}

// Running returns the current state of the infinite progress animation
func (p *ProgressBarInfinite) Running() bool {
	return p.running
}

// MinSize returns the size that this widget should not shrink below
func (p *ProgressBarInfinite) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *ProgressBarInfinite) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	primaryColor := th.Color(theme.ColorNamePrimary, v)
	cornerRadius := th.Size(theme.SizeNameInputRadius)

	render := &infProgressRenderer{
		background: canvas.Rectangle{
			FillColor:    progressBlendColor(primaryColor),
			CornerRadius: cornerRadius,
		},
		bar: canvas.Rectangle{
			FillColor:    primaryColor,
			CornerRadius: cornerRadius,
		},
		progress: p,
	}

	render.SetObjects([]fyne.CanvasObject{&render.background, &render.bar})

	p.running = true
	return render
}

// NewProgressBarInfinite creates a new progress bar widget that loops indefinitely from 0% -> 100%
// SetValue() is not defined for infinite progress bar
// To stop the looping progress and set the progress bar to 100%, call ProgressBarInfinite.Stop()
func NewProgressBarInfinite() *ProgressBarInfinite {
	bar := &ProgressBarInfinite{}
	bar.ExtendBaseWidget(bar)
	return bar
}
