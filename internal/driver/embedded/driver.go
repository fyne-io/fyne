package embedded

import (
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	intdriver "fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter"
)

type noosDriver struct {
	events chan embedded.Event
	queue  chan funcData
	render func(image.Image)
	run    func(func())
	size   func() fyne.Size
	done   bool

	wins    []fyne.Window
	current int
	device  noosDevice
}

func (n *noosDriver) CreateWindow(_ string) fyne.Window {
	w := newWindow(n)
	n.wins = append(n.wins, w)
	n.current = len(n.wins) - 1

	if f := n.size; f != nil {
		w.Resize(f())
	}
	return w
}

func (n *noosDriver) AllWindows() []fyne.Window {
	return n.wins
}

func (n *noosDriver) RenderedTextSize(text string, fontSize float32, style fyne.TextStyle, source fyne.Resource) (size fyne.Size, baseline float32) {
	return painter.RenderedTextSize(text, fontSize, style, source)
}

func (n *noosDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return cache.GetCanvasForObject(obj)
}

func (n *noosDriver) AbsolutePositionForObject(o fyne.CanvasObject) fyne.Position {
	c := n.CanvasForObject(o)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	pos := intdriver.AbsolutePositionForObject(o, []fyne.CanvasObject{c.Content()})
	inset, _ := c.InteractiveArea()
	return pos.Subtract(inset)
}

func (n *noosDriver) Device() fyne.Device {
	return &n.device
}

func (n *noosDriver) Run() {
	n.run(n.doRun)
}

func (n *noosDriver) doRun() {
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
		case e := <-n.events:
			if e == nil {
				// closing
				n.Quit()
				continue
			}

			w := n.wins[n.current].(*noosWindow)

			switch t := e.(type) {
			case *embedded.CharacterEvent:
				if focused := w.c.Focused(); focused != nil {
					focused.TypedRune(t.Rune)
				} else if tr := w.c.OnTypedRune(); tr != nil {
					tr(t.Rune)
				}

				n.renderWindow(n.wins[n.current])
			case *embedded.KeyEvent:
				keyEvent := &fyne.KeyEvent{Name: t.Name}

				if t.Direction == embedded.KeyReleased {
					// No desktop events so key/up down not reported
					continue // ignore key up in other core events
				}

				if t.Name == fyne.KeyTab {
					captures := false

					if ent, ok := w.Canvas().Focused().(fyne.Tabbable); ok {
						captures = ent.AcceptsTab()
					}
					if !captures {
						// TODO handle shift
						w.Canvas().FocusNext()
						n.renderWindow(n.wins[n.current])
						continue
					}
				}

				// No shortcut detected, pass down to TypedKey
				focused := w.c.Focused()
				if focused != nil {
					focused.TypedKey(keyEvent)
				} else if tk := w.c.OnTypedKey(); tk != nil {
					tk(keyEvent)
				}

				n.renderWindow(n.wins[n.current])
			case *embedded.TouchDownEvent:
				n.handleTouchDown(t, n.wins[n.current].(*noosWindow))
			case *embedded.TouchMoveEvent:
				n.handleTouchMove(t, n.wins[n.current].(*noosWindow))
			case *embedded.TouchUpEvent:
				n.handleTouchUp(t, n.wins[n.current].(*noosWindow))
			}
		}
	}
}

func (n *noosDriver) handleTouchDown(ev *embedded.TouchDownEvent, w *noosWindow) {
	w.c.tapDown(ev.Position, ev.ID)
	n.renderWindow(w)
}

func (n *noosDriver) handleTouchMove(ev *embedded.TouchMoveEvent, w *noosWindow) {
	w.c.tapMove(ev.Position, ev.ID, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
	})
	n.renderWindow(w)
}

func (n *noosDriver) handleTouchUp(ev *embedded.TouchUpEvent, w *noosWindow) {
	w.c.tapUp(ev.Position, ev.ID, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		wid.TappedSecondary(ev)
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		if math.Abs(float64(ev.Dragged.DX)) <= tapMoveEndThreshold && math.Abs(float64(ev.Dragged.DY)) <= tapMoveEndThreshold {
			wid.DragEnd()
			return
		}

		go func() {
			for math.Abs(float64(ev.Dragged.DX)) > tapMoveEndThreshold || math.Abs(float64(ev.Dragged.DY)) > tapMoveEndThreshold {
				if math.Abs(float64(ev.Dragged.DX)) > 0 {
					ev.Dragged.DX *= tapMoveDecay
				}
				if math.Abs(float64(ev.Dragged.DY)) > 0 {
					ev.Dragged.DY *= tapMoveDecay
				}

				n.DoFromGoroutine(func() {
					wid.Dragged(ev)
				}, false)
				time.Sleep(time.Millisecond * 16)
			}

			n.DoFromGoroutine(wid.DragEnd, false)
		}()
	})
	n.renderWindow(w)
}

func (n *noosDriver) Quit() {
	n.done = true

	go func() {
		n.queue <- funcData{f: func() {}}
	}()
}

func (n *noosDriver) StartAnimation(*fyne.Animation) {
	// no animations on embedded
}

func (n *noosDriver) StopAnimation(*fyne.Animation) {
	// no animations on embedded
}

func (n *noosDriver) DoubleTapDelay() time.Duration {
	return tapDoubleDelay
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

func NewNoOSDriver(render func(img image.Image), run func(func()), events chan embedded.Event, size func() fyne.Size) fyne.Driver {
	return &noosDriver{
		events: events, queue: make(chan funcData), size: size,
		render: render, run: run, wins: make([]fyne.Window, 0),
	}
}

type funcData struct {
	f    func()
	done chan struct{} // Zero allocation signalling channel
}
