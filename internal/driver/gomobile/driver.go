package gomobile

import (
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"

	"fyne.io/fyne"
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
	canvas := &canvas{scale: 2} // TODO detect scale
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
}

func (d *driver) Run() {
	app.Main(func(a app.App) {
		d.app = a

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
			case size.Event:
				sz = e
			case paint.Event:
				if d.glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				d.onPaint(sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				// TODO handle input
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
	}
}

// NewGoMobileDriver sets up a new Driver instance implemented using the Go
// Mobile extension and OpenGL bindings.
func NewGoMobileDriver() fyne.Driver {
	driver := new(driver)

	return driver
}
