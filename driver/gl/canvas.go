// +build !ci,gl

package gl

import (
	"math"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type glCanvas struct {
	window  *window
	content fyne.CanvasObject
	focused fyne.FocusableObject

	onKeyDown func(*fyne.KeyEvent)

	program uint32
	scale   float32
}

func scaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(math.Round(float64(v) * float64(c.Scale())))
	}
}

func unscaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}

func (c *glCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *glCanvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	var w, h = c.window.viewport.GetSize()

	pad := theme.Padding()
	if c.window.FullScreen() {
		pad = 0
	}
	width := unscaleInt(c, int(w)) - pad*2
	height := unscaleInt(c, int(h)) - pad*2

	c.content.Resize(fyne.NewSize(width, height))
	c.content.Move(fyne.NewPos(pad, pad))
}

func (c *glCanvas) Refresh(fyne.CanvasObject) {
	// no-op
}

func (c *glCanvas) Contains(fyne.CanvasObject) bool {
	return true
}

func (c *glCanvas) Focus(obj fyne.FocusableObject) {
	if c.focused != nil {
		c.focused.(fyne.FocusableObject).OnFocusLost()
	}

	c.focused = obj
	obj.OnFocusGained()
}

func (c *glCanvas) Focused() fyne.FocusableObject {
	return c.focused
}

func (c *glCanvas) Size() fyne.Size {
	var w, h = c.window.viewport.GetSize()

	width := unscaleInt(c, int(w))
	height := unscaleInt(c, int(h))

	return fyne.NewSize(width, height)
}

func (c *glCanvas) Scale() float32 {
	return c.scale
}

func (c *glCanvas) SetScale(scale float32) {
	c.scale = scale
}

func (c *glCanvas) SetOnKeyDown(keyDown func(*fyne.KeyEvent)) {
	c.onKeyDown = keyDown
}

func (c *glCanvas) refresh() {
	// TODO check if it's really needed...

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r, g, b, a := theme.BackgroundColor().RGBA()
	gl.ClearColor(float32(uint8(r))/255, float32(uint8(g))/255, float32(uint8(b))/255, float32(uint8(a))/255)

	if c.content == nil {
		return
	}

	walkObjects(c.content, fyne.NewPos(0, 0), c.drawObject)
}

func newCanvas(win *window, program uint32) *glCanvas {
	return &glCanvas{window: win, scale: 1.0, program: program}
}
