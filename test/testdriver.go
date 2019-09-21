package test

import (
	"image"

	"fyne.io/fyne"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
)

// SoftwarePainter describes a simple type that can render canvases
type SoftwarePainter interface {
	Paint(fyne.Canvas) image.Image
}

type testDriver struct {
	device  *device
	painter SoftwarePainter
}

// Declare conformity with Driver
var _ fyne.Driver = (*testDriver)(nil)

func (d *testDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	windowsMutex.RLock()
	defer windowsMutex.RUnlock()
	// cheating as we only have a single test window
	return windows[0].Canvas()
}

func (d *testDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	return co.Position() // TODO get the real answer
}

func (d *testDriver) CreateWindow(string) fyne.Window {
	return NewWindow(nil)
}

func (d *testDriver) AllWindows() []fyne.Window {
	windowsMutex.RLock()
	defer windowsMutex.RUnlock()
	return windows
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

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver() fyne.Driver {
	driver := new(testDriver)
	// make a single dummy window for rendering tests
	NewWindow(nil)

	return driver
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter SoftwarePainter) fyne.Driver {
	driver := new(testDriver)
	driver.painter = painter

	return driver
}
