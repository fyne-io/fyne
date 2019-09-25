package test

import (
	"image"
	"sync"

	"fyne.io/fyne"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
)

// SoftwarePainter describes a simple type that can render canvases
type SoftwarePainter interface {
	Paint(fyne.Canvas) image.Image
}

type testDriver struct {
	device       *device
	painter      SoftwarePainter
	windows      []fyne.Window
	windowsMutex sync.RWMutex
}

// Declare conformity with Driver
var _ fyne.Driver = (*testDriver)(nil)

func (d *testDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	// cheating as we only have a single test window
	return d.windows[0].Canvas()
}

func (d *testDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	return co.Position() // TODO get the real answer
}

func (d *testDriver) CreateWindow(string) fyne.Window {
	canvas := NewCanvas().(*testCanvas)
	canvas.painter = d.painter

	window := &testWindow{canvas: canvas, driver: d}
	window.clipboard = &testClipboard{}

	d.windowsMutex.Lock()
	d.windows = append(d.windows, window)
	d.windowsMutex.Unlock()
	return window
}

func (d *testDriver) AllWindows() []fyne.Window {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	return d.windows
}

// RenderedTextSize looks up how bit a string would be if drawn on screen
func (d *testDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	var opts truetype.Options
	opts.Size = float64(size)
	opts.DPI = 78 // TODO move this?

	theme := fyne.CurrentApp().Settings().Theme()
	// TODO check style
	f, err := truetype.Parse(theme.TextFont().Content())
	if err != nil {
		fyne.LogError("Unable to load theme font", err)
	}
	face := truetype.NewFace(f, &opts)
	advance := font.MeasureString(face, text)

	return fyne.NewSize(advance.Ceil(), face.Metrics().Height.Ceil())
}

func (d *testDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}
	return d.device
}

func (d *testDriver) Run() {
	// no-op
}

func (d *testDriver) Quit() {
	// no-op
}

func (d *testDriver) removeWindow(w *testWindow) {
	d.windowsMutex.Lock()
	i := 0
	for _, window := range d.windows {
		if window == w {
			break
		}
		i++
	}

	d.windows = append(d.windows[:i], d.windows[i+1:]...)
	d.windowsMutex.Unlock()
}

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver() fyne.Driver {
	driver := new(testDriver)
	// make a single dummy window for rendering tests
	driver.CreateWindow("")

	return driver
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter SoftwarePainter) fyne.Driver {
	driver := new(testDriver)
	driver.painter = painter

	driver.windows = make([]fyne.Window, 0)
	driver.windowsMutex = sync.RWMutex{}

	return driver
}
