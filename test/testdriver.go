package test

import (
	"image"
	"log"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/internal/painter/software"

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

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver() fyne.Driver {
	drv := new(testDriver)
	drv.windowsMutex = sync.RWMutex{}

	// make a single dummy window for rendering tests
	drv.CreateWindow("")

	return drv
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter SoftwarePainter) fyne.Driver {
	drv := new(testDriver)
	drv.painter = painter
	drv.windowsMutex = sync.RWMutex{}

	return drv
}

func (d *testDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	tc := c.(*testCanvas)
	return driver.AbsolutePositionForObject(co, tc.objectTrees())
}

func (d *testDriver) AllWindows() []fyne.Window {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	return d.windows
}

func (d *testDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	// cheating: probably the last created window is meant
	return d.windows[len(d.windows)-1].Canvas()
}

func (d *testDriver) CreateWindow(string) fyne.Window {
	canvas := NewCanvas().(*testCanvas)
	if d.painter != nil {
		canvas.painter = d.painter
	} else {
		canvas.painter = software.NewPainter()
	}

	window := &testWindow{canvas: canvas, driver: d}
	window.clipboard = &testClipboard{}

	d.windowsMutex.Lock()
	d.windows = append(d.windows, window)
	d.windowsMutex.Unlock()
	return window
}

func (d *testDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}
	return d.device
}

// RenderedTextSize looks up how bit a string would be if drawn on screen
func (d *testDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	var opts truetype.Options
	opts.Size = float64(size)
	opts.DPI = painter.TextDPI

	face := painter.CachedFontFace(style, &opts)
	advance := font.MeasureString(face, text)

	sws := fyne.NewSize(advance.Ceil(), face.Metrics().Height.Ceil())
	gls := painter.RenderedTextSize(text, size, style)
	if sws != gls {
		log.Println("SoftwareTextSize:", sws)
		log.Println("GLTextSize:", gls)
	}
	return sws
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
