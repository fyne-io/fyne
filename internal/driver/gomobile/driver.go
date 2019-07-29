package gomobile

import (
	"log"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"

	"fyne.io/fyne"
	util "fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
	pgl "fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
)

type driver struct {
	app   app.App
	glctx gl.Context

	windows []fyne.Window
}

func (d *driver) CreateWindow(title string) fyne.Window {
	canvas := NewCanvas().(*canvas) // silence lint
	ret := &window{title: title, canvas: canvas, padded: true}
	canvas.painter = pgl.NewPainter(canvas, ret)

	d.windows = append(d.windows, ret)
	return ret
}

func (d *driver) AllWindows() []fyne.Window {
	return d.windows
}

func (d *driver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	return painter.RenderedTextSize(text, size, style)
}

func (d *driver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	if len(d.windows) == 0 {
		return nil
	}

	// TODO figure out how we handle multiple windows...
	return d.windows[0].Canvas()
}

func (d *driver) AbsolutePositionForObject(fyne.CanvasObject) fyne.Position {
	log.Println("TODO - absolute position!")
	return fyne.NewPos(0, 0)
}

func (d *driver) Quit() {
	// TODO? often mobile apps should not allow this...
	d.app.Send(lifecycle.Event{From: lifecycle.StageVisible, To: lifecycle.StageDead, DrawContext: nil})
}

func (d *driver) scheduleFrames(a app.App) {
	fps := time.NewTicker(time.Second / 60)
	go func() {
		for {
			select {
			case <-fps.C:
				a.Send(paint.Event{})
			}
		}
	}()
}

func (d *driver) Run() {
	app.Main(func(a app.App) {
		d.app = a
		quit := false
		d.scheduleFrames(a)

		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					d.glctx, _ = e.DrawContext.(gl.Context)
					d.onStart()
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					d.onStop()
					d.glctx = nil
				}
				if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOff {
					quit = true
				}
			case size.Event:
				sz = e
			case paint.Event:
				if len(d.AllWindows()) == 0 {
					break
				}
				canvas := d.AllWindows()[0].Canvas().(*canvas)

				if canvas.dirty && d.glctx != nil {
					d.freeDirtyTextures(canvas)

					d.onPaint(sz)
					a.Publish()
					canvas.dirty = false
				}
			case touch.Event:
				switch e.Type {
				case touch.TypeBegin:
				case touch.TypeEnd:
					d.onTapEnd(e.X, e.Y)
				}
			}

			if quit {
				break
			}
		}
	})
}

func (d *driver) onStart() {
	for _, win := range d.AllWindows() {
		win.Canvas().(*canvas).painter.Init() // we cannot init until the context is set above
	}
}

func (d *driver) onStop() {
}

func (d *driver) onPaint(sz size.Event) {
	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	d.glctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	d.glctx.Clear(gl.COLOR_BUFFER_BIT)

	if len(d.AllWindows()) > 0 {
		canvas := d.AllWindows()[0].Canvas().(*canvas)
		newSize := fyne.NewSize(int(float32(sz.WidthPx)/canvas.scale), int(float32(sz.HeightPx)/canvas.scale))
		canvas.Resize(newSize)
		canvas.painter.Paint(canvas.content, canvas, canvas.Size())
		if canvas.overlay != nil {
			canvas.painter.Paint(canvas.overlay, canvas, canvas.Size())
		}
	}
}

func (d *driver) onTapEnd(x, y float32) {
	if len(d.AllWindows()) == 0 {
		return
	}

	canvas := d.AllWindows()[0].Canvas().(*canvas)
	tapX := util.UnscaleInt(canvas, int(x))
	tapY := util.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY)

	co, objX, objY := util.FindObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
			return true
		}

		return false
	}, canvas.overlay, canvas.content)

	ev := new(fyne.PointEvent)
	ev.Position = fyne.NewPos(tapX-objX, tapY-objY)

	if wid, ok := co.(fyne.Tappable); ok {
		// TODO move event queue to common code w.queueEvent(func() { wid.Tapped(ev) })
		go wid.Tapped(ev)
	}
}

func (d *driver) freeDirtyTextures(canvas *canvas) {
	for {
		select {
		case object := <-canvas.refreshQueue:
			freeWalked := func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
				canvas.painter.Free(obj)
				log.Println("Free", obj)
				return false
			}
			util.WalkObjectTree(object, freeWalked, nil)
		default:
			return
		}
	}
}

// NewGoMobileDriver sets up a new Driver instance implemented using the Go
// Mobile extension and OpenGL bindings.
func NewGoMobileDriver() fyne.Driver {
	driver := new(driver)

	return driver
}
