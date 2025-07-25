package noos

import (
	"image"
	"time"

	"fyne.io/fyne/v2"
	noos2 "fyne.io/fyne/v2/driver/noos"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	intdriver "fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter"
)

type noosDriver struct {
	events chan noos2.Event
	queue  chan funcData
	render func(image.Image)
	done   bool

	wins    []fyne.Window
	current int
	device  noosDevice
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
			case *noos2.CharacterEvent:
				if focused := w.c.Focused(); focused != nil {
					focused.TypedRune(t.Rune)
				} else if tr := w.c.OnTypedRune(); tr != nil {
					tr(t.Rune)
				}

				n.renderWindow(n.wins[n.current])
			case *noos2.KeyEvent:
				keyEvent := &fyne.KeyEvent{Name: t.Name}

				if t.Direction == noos2.KeyReleased {
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
			}
		}
	}
}

func (n *noosDriver) Quit() {
	n.done = true

	go func() {
		n.queue <- funcData{f: func() {}}
	}()
}

func (n *noosDriver) StartAnimation(*fyne.Animation) {
	// no animations on noos
}

func (n *noosDriver) StopAnimation(*fyne.Animation) {
	// no animations on noos
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
