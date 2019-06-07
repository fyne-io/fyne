// Package gl provides a full Fyne render implementation using system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package gl

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
)

var canvasMutex sync.RWMutex
var canvases = make(map[fyne.CanvasObject]fyne.Canvas)

const textDPI = 78

// Declare conformity with Driver
var _ fyne.Driver = (*gLDriver)(nil)

type gLDriver struct {
	windows []fyne.Window
	done    chan interface{}
}

func (d *gLDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	canvasMutex.RLock()
	defer canvasMutex.RUnlock()
	return canvases[obj]
}

func (d *gLDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	var pos fyne.Position
	c := fyne.CurrentApp().Driver().CanvasForObject(co).(*glCanvas)

	driver.WalkObjectTree(c.content, func(o fyne.CanvasObject, p fyne.Position) bool {
		if o == co {
			pos = p
			return true
		}
		return false
	}, nil)

	return pos
}

func loadFont(data fyne.Resource) *truetype.Font {
	loaded, err := truetype.Parse(data.Content())
	if err != nil {
		fyne.LogError("font load error", err)
	}

	return loaded
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	var opts truetype.Options
	opts.Size = float64(size)
	opts.DPI = textDPI

	face := cachedFontFace(style, &opts)
	advance := font.MeasureString(face, text)

	return fyne.NewSize(advance.Ceil(), face.Metrics().Height.Ceil())
}

func (d *gLDriver) Quit() {
	defer func() {
		recover() // we could be called twice - no safe way to check if d.done is closed
	}()
	close(d.done)
}

func (d *gLDriver) Run() {
	go svgCacheMonitorTheme()
	d.runGL()
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() fyne.Driver {
	driver := new(gLDriver)
	driver.done = make(chan interface{})

	return driver
}
