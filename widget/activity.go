package widget

import (
	"image/color"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Activity)(nil)

// Activity is used to indicate that something is happening that should be waited for,
// or is in the background (depending on usage).
//
// Since: 2.5
type Activity struct {
	BaseWidget

	started atomic.Bool
}

// NewActivity returns a widget for indicating activity
//
// Since: 2.5
func NewActivity() *Activity {
	a := &Activity{}
	a.ExtendBaseWidget(a)
	return a
}

func (a *Activity) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// Start the activity indicator animation
func (a *Activity) Start() {
	if !a.started.CompareAndSwap(false, true) {
		return // already started
	}

	a.Refresh()
}

// Stop the activity indicator animation
func (a *Activity) Stop() {
	if !a.started.CompareAndSwap(true, false) {
		return // already stopped
	}

	a.Refresh()
}

func (a *Activity) CreateRenderer() fyne.WidgetRenderer {
	dots := make([]fyne.CanvasObject, 3)
	v := fyne.CurrentApp().Settings().ThemeVariant()
	for i := range dots {
		dots[i] = canvas.NewCircle(a.Theme().Color(theme.ColorNameForeground, v))
	}
	r := &activityRenderer{dots: dots, parent: a}
	r.anim = &fyne.Animation{
		Duration:    time.Second * 2,
		RepeatCount: fyne.AnimationRepeatForever,
		Tick:        r.animate}
	r.updateColor()

	if a.started.Load() {
		r.start()
	}

	return r
}

var _ fyne.WidgetRenderer = (*activityRenderer)(nil)

type activityRenderer struct {
	anim   *fyne.Animation
	dots   []fyne.CanvasObject
	parent *Activity

	bound      fyne.Size
	maxCol     color.NRGBA
	maxRad     float32
	wasStarted bool
}

func (a *activityRenderer) Destroy() {
	a.parent.started.Store(false)
	a.stop()
}

func (a *activityRenderer) Layout(size fyne.Size) {
	a.maxRad = fyne.Min(size.Width, size.Height) / 2
	a.bound = size
}

func (a *activityRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(a.parent.Theme().Size(theme.SizeNameInlineIcon))
}

func (a *activityRenderer) Objects() []fyne.CanvasObject {
	return a.dots
}

func (a *activityRenderer) Refresh() {
	started := a.parent.started.Load()
	if started {
		if !a.wasStarted {
			a.start()
		}
		return
	} else if a.wasStarted {
		a.stop()
		return
	}

	a.updateColor()
}

func (a *activityRenderer) animate(done float32) {
	off := done * 2
	if off > 1 {
		off = 2 - off
	}

	off1 := (done + 0.25) * 2
	if done >= 0.75 {
		off1 = (done - 0.75) * 2
	}
	if off1 > 1 {
		off1 = 2 - off1
	}

	off2 := (done + 0.75) * 2
	if done >= 0.25 {
		off2 = (done - 0.25) * 2
	}
	if off2 > 1 {
		off2 = 2 - off2
	}

	a.scaleDot(a.dots[0].(*canvas.Circle), off)
	a.scaleDot(a.dots[1].(*canvas.Circle), off1)
	a.scaleDot(a.dots[2].(*canvas.Circle), off2)
}

func (a *activityRenderer) scaleDot(dot *canvas.Circle, off float32) {
	rad := a.maxRad - a.maxRad*off/1.2
	mid := fyne.NewPos(a.bound.Width/2, a.bound.Height/2)

	dot.Move(mid.Subtract(fyne.NewSquareOffsetPos(rad)))
	dot.Resize(fyne.NewSquareSize(rad * 2))

	alpha := uint8(0 + int(float32(a.maxCol.A)*off))
	dot.FillColor = color.NRGBA{R: a.maxCol.R, G: a.maxCol.G, B: a.maxCol.B, A: alpha}
	dot.Refresh()
}

func (a *activityRenderer) start() {
	a.wasStarted = true
	a.anim.Start()
}

func (a *activityRenderer) stop() {
	a.wasStarted = false
	a.anim.Stop()
}

func (a *activityRenderer) updateColor() {
	v := fyne.CurrentApp().Settings().ThemeVariant()
	rr, gg, bb, aa := a.parent.Theme().Color(theme.ColorNameForeground, v).RGBA()
	a.maxCol = color.NRGBA{R: uint8(rr >> 8), G: uint8(gg >> 8), B: uint8(bb >> 8), A: uint8(aa >> 8)}
}
