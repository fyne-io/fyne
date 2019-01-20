// Package gl provides a full Fyne render implementation using system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package gl

import (
	"log"
	"sync"

	"fyne.io/fyne"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var canvasMutex sync.RWMutex
var canvases = make(map[fyne.CanvasObject]fyne.Canvas)

const textDPI = 78

type gLDriver struct {
	windows []fyne.Window
	done    chan interface{}
}

func (d *gLDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	canvasMutex.RLock()
	defer canvasMutex.RUnlock()
	return canvases[obj]
}

func loadFont(data fyne.Resource) *truetype.Font {
	loaded, err := truetype.Parse(data.Content())
	if err != nil {
		log.Println("Font load error", err)
	}

	return loaded
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	var opts truetype.Options
	opts.Size = float64(size)
	opts.DPI = textDPI

	face := cachedFontFace(style, &opts)
	advance := font.MeasureString(face, text)

	return fyne.NewSize(advance.Ceil(), face.Metrics().Height.Ceil()+face.Metrics().Descent.Ceil())
}

func (d *gLDriver) Quit() {
	close(d.done)
}

func (d *gLDriver) Run() {
	d.runGL()
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() fyne.Driver {
	driver := new(gLDriver)
	driver.done = make(chan interface{})

	return driver
}
