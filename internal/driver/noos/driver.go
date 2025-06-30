package noos

import (
	"fyne.io/fyne/v2"
	noos2 "fyne.io/fyne/v2/driver/noos"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/painter"
	"image"
	"time"
)

type noosDriver struct {
	events chan noos2.Event
	queue  chan funcData
	render func(image.Image)
	done   bool

	wins    []fyne.Window
	current int
}

func (n *noosDriver) CreateWindow(_ string) fyne.Window {
	w := newWindow(n)
	n.wins = append(n.wins, w)
	n.current = len(n.wins) - 1

	return w
}

func (n *noosDriver) AllWindows() []fyne.Window {
	return n.wins
}

func (n *noosDriver) RenderedTextSize(text string, fontSize float32, style fyne.TextStyle, source fyne.Resource) (size fyne.Size, baseline float32) {
	return painter.RenderedTextSize(text, fontSize, style, source)
}

func (n *noosDriver) CanvasForObject(_ fyne.CanvasObject) fyne.Canvas {
	//TODO implement me
	return n.AllWindows()[n.current].Canvas()
}

func (n *noosDriver) AbsolutePositionForObject(o fyne.CanvasObject) fyne.Position {
	//TODO implement me
	return fyne.Position{}
}

func (n *noosDriver) Device() fyne.Device {
	//TODO implement me
	panic("implement me")
}

func (n *noosDriver) Run() {
	for _, w := range n.wins {
		n.renderWindow(w)
	}

	for !n.done {
		select {
		case fn := <-n.queue:
			fn.f()
			if fn.done != nil {
				fn.done <- struct{}{}
			}
		case <-n.events:
			// TODO actually process events, shared keyboard event code for a start
			n.renderWindow(n.wins[n.current])
		}
	}
}

func (n *noosDriver) Quit() {
	n.queue <- funcData{
		f: func() {
			n.done = true
		}}
}

func (n *noosDriver) StartAnimation(anim *fyne.Animation) {
	//TODO implement me
}

func (n *noosDriver) StopAnimation(anim *fyne.Animation) {
	//TODO implement me
}

func (n *noosDriver) DoubleTapDelay() time.Duration {
	return time.Duration(150)
}

func (n *noosDriver) SetDisableScreenBlanking(bool) {}

func (n *noosDriver) DoFromGoroutine(fn func(), wait bool) {
	if wait {
		async.EnsureNotMain(func() {
			done := make(chan struct{})

			n.queue <- funcData{f: fn, done: done}
			<-done
		})
	} else {
		n.queue <- funcData{f: fn}
	}
}

// TODO add some caching to stop allocating images...
func (n *noosDriver) renderWindow(w fyne.Window) {
	img := w.Canvas().Capture()

	n.render(img)
}

func NewNoOSDriver(render func(img image.Image), events chan noos2.Event) fyne.Driver {
	return &noosDriver{events: events, queue: make(chan funcData),
		render: render, wins: make([]fyne.Window, 0)}
}

type funcData struct {
	f    func()
	done chan struct{} // Zero allocation signalling channel
}
