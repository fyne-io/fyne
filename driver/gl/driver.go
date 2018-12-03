// +build !ci,gl

package gl

import (
	"log"
	"sync"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
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
	return canvases[obj]
}

// TODO for styles...
var fontRegular *truetype.Font

func fontCache() *truetype.Font {
	if fontRegular != nil {
		return fontRegular
	}

	loaded, err := truetype.Parse(theme.TextFont().Content())
	if err != nil {
		log.Println("Font error", err)
	}
	fontRegular = loaded
	return fontRegular
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	var opts truetype.Options
	opts.Size = float64(size)
	opts.DPI = textDPI

	face := truetype.NewFace(fontCache(), &opts)
	advance := font.MeasureString(face, text)

	return fyne.NewSize(advance.Ceil(), face.Metrics().Height.Ceil())
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
