package gomobile

import (
	"time"

	"fyne.io/fyne/widget"
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
	device  fyne.Device
}

// Declare conformity with Driver
var _ fyne.Driver = (*driver)(nil)

func (d *driver) CreateWindow(title string) fyne.Window {
	canvas := NewCanvas().(*canvas) // silence lint
	ret := &window{title: title, canvas: canvas}
	canvas.painter = pgl.NewPainter(canvas, ret)

	d.windows = append(d.windows, ret)
	return ret
}

func (d *driver) AllWindows() []fyne.Window {
	return d.windows
}

// currentWindow returns the most recently opened window - we can only show one at a time.
func (d *driver) currentWindow() fyne.Window {
	if len(d.windows) == 0 {
		return nil
	}

	// TODO ensure we remove new windows once closed
	return d.windows[len(d.windows)-1]
}

func (d *driver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	return painter.RenderedTextSize(text, size, style)
}

func (d *driver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	if len(d.windows) == 0 {
		return nil
	}

	// TODO figure out how we handle multiple windows...
	return d.currentWindow().Canvas()
}

func (d *driver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	var pos fyne.Position
	c := fyne.CurrentApp().Driver().CanvasForObject(co).(*canvas)

	util.WalkVisibleObjectTree(c.content, func(o fyne.CanvasObject, p fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if o == co {
			pos = p
			return true
		}
		return false
	}, nil)

	return pos
}

func (d *driver) Quit() {
	if d.app == nil {
		return
	}

	// TODO? often mobile apps should not allow this...
	d.app.Send(lifecycle.Event{From: lifecycle.StageAlive, To: lifecycle.StageDead, DrawContext: nil})
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
				if e.Crosses(lifecycle.StageAlive) == lifecycle.CrossOff {
					quit = true
				}
			case size.Event:
				sz = e
			case paint.Event:
				current := d.currentWindow()
				if current == nil {
					break
				}
				canvas := current.Canvas().(*canvas)
				if !canvas.inited && d.glctx != nil {
					canvas.inited = true
					canvas.painter.Init() // we cannot init until the context is set above
				}

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
	currentOrientation = sz.Orientation

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	d.glctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	d.glctx.Clear(gl.COLOR_BUFFER_BIT)

	current := d.currentWindow()
	if current == nil {
		return
	}

	canvas := current.Canvas().(*canvas)
	newSize := fyne.NewSize(int(float32(sz.WidthPx)/canvas.scale), int(float32(sz.HeightPx)/canvas.scale))
	canvas.Resize(newSize)

	paint := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		// TODO should this be somehow not scroll container specific?
		if _, ok := obj.(*widget.ScrollContainer); ok {
			canvas.painter.StartClipping(
				fyne.NewPos(pos.X, canvas.Size().Height-pos.Y-obj.Size().Height),
				obj.Size(),
			)
		}
		canvas.painter.Paint(obj, pos, newSize)
		return false
	}
	afterPaint := func(obj, _ fyne.CanvasObject) {
		if _, ok := obj.(*widget.ScrollContainer); ok {
			canvas.painter.StopClipping()
		}
	}

	canvas.walkTree(paint, afterPaint)
}

func (d *driver) onTapEnd(x, y float32) {
	current := d.currentWindow()
	if current == nil {
		return
	}

	canvas := current.Canvas().(*canvas)
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
	ev.Position = fyne.NewPos(objX, objY)

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
				return false
			}
			util.WalkCompleteObjectTree(object, freeWalked, nil)
		default:
			return
		}
	}
}

func (d *driver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}

	return d.device
}

// NewGoMobileDriver sets up a new Driver instance implemented using the Go
// Mobile extension and OpenGL bindings.
func NewGoMobileDriver() fyne.Driver {
	driver := new(driver)

	return driver
}
