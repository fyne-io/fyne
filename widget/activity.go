package widget

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Activity)(nil)

// Activity is used to indicate that something is happening that should be waited for,
// or is in the background (depending on usage).
//
// Since: 2.5
type Activity struct {
	BaseWidget

	started bool
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

	return fyne.NewSquareSize(theme.IconInlineSize())
}

// Start the activity indicator animation
func (a *Activity) Start() {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	if a.started {
		return
	}

	a.started = true
	if r, ok := cache.Renderer(a.super()).(*activityRenderer); ok {
		r.start()
	}
}

// Stop the activity indicator animation
func (a *Activity) Stop() {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	if !a.started {
		return
	}

	a.started = false
	if r, ok := cache.Renderer(a.super()).(*activityRenderer); ok {
		r.stop()
	}
}

func (a *Activity) CreateRenderer() fyne.WidgetRenderer {
	dots := make([]fyne.CanvasObject, 3)
	for i := range dots {
		dots[i] = canvas.NewCircle(theme.ForegroundColor())
	}
	r := &activityRenderer{dots: dots}
	r.anim = &fyne.Animation{
		Duration:    time.Second * 2,
		RepeatCount: fyne.AnimationRepeatForever,
		Tick:        r.animate}
	r.updateColor()

	a.propertyLock.RLock()
	if a.started {
		r.start()
	}
	a.propertyLock.RUnlock()

	return r
}

var _ fyne.WidgetRenderer = (*activityRenderer)(nil)

type activityRenderer struct {
	anim *fyne.Animation
	dots []fyne.CanvasObject

	bound  fyne.Size
	maxCol color.NRGBA
	maxRad float32
}

func (a *activityRenderer) Destroy() {
	a.stop()
}

func (a *activityRenderer) Layout(size fyne.Size) {
	a.maxRad = fyne.Min(size.Width, size.Height) / 2
	a.bound = size
}

func (a *activityRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(theme.IconInlineSize())
}

func (a *activityRenderer) Objects() []fyne.CanvasObject {
	return a.dots
}

func (a *activityRenderer) Refresh() {
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
	a.anim.Start()
}

func (a *activityRenderer) stop() {
	a.anim.Stop()
}

func (a *activityRenderer) updateColor() {
	rr, gg, bb, aa := theme.ForegroundColor().RGBA()
	a.maxCol = color.NRGBA{R: uint8(rr >> 8), G: uint8(gg >> 8), B: uint8(bb >> 8), A: uint8(aa >> 8)}
}
