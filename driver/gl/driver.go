// Package gl provides a full Fyne render implementation using system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
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
	canvasMutex.RLock()
	defer canvasMutex.RUnlock()
	return canvases[obj]
}

var fontFaces = make(map[fyne.TextStyle]*truetype.Font)

func loadFont(data fyne.Resource) *truetype.Font {
	loaded, err := truetype.Parse(data.Content())
	if err != nil {
		log.Println("Font load error", err)
	}

	return loaded
}

func fontCache(style fyne.TextStyle) *truetype.Font {
	if fontFaces[style] != nil {
		return fontFaces[style]
	}

	var loaded *truetype.Font
	if style.Monospace {
		loaded = loadFont(theme.TextMonospaceFont())
	} else if style.Bold {
		if style.Italic {
			loaded = loadFont(theme.TextBoldItalicFont())
		} else {
			loaded = loadFont(theme.TextBoldFont())
		}
	} else if style.Italic {
		loaded = loadFont(theme.TextItalicFont())
	} else {
		loaded = loadFont(theme.TextFont())
	}

	fontFaces[style] = loaded
	return loaded
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	var opts truetype.Options
	opts.Size = float64(size)
	opts.DPI = textDPI

	face := truetype.NewFace(fontCache(style), &opts)
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
