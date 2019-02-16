package gl

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type glCanvas struct {
	window  *window
	content fyne.CanvasObject
	focused fyne.Focusable

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	program uint32
	scale   float32

	dirty1, dirty2 bool
	refreshQueue   chan fyne.CanvasObject
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
	if !c.window.Padded() {
		pad = 0
	}
	width := unscaleInt(c, int(w)) - pad*2
	height := unscaleInt(c, int(h)) - pad*2

	c.content.Resize(fyne.NewSize(width, height))
	c.content.Move(fyne.NewPos(pad, pad))
	c.setDirty()
}

func (c *glCanvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
	c.setDirty()
}

func (c *glCanvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		c.focused.(fyne.Focusable).FocusLost()
	}

	c.focused = obj
	obj.FocusGained()
}

func (c *glCanvas) Focused() fyne.Focusable {
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
	c.setDirty()
}

func (c *glCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *glCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *glCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *glCanvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *glCanvas) paint(size fyne.Size) {
	if c.dirty1 {
		c.dirty1 = false
	} else {
		if c.dirty2 {
			c.dirty2 = false
		}
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)

	if c.content == nil {
		return
	}

	paintObj := func(obj fyne.CanvasObject, pos fyne.Position) {
		c.drawObject(obj, pos, size)
	}
	c.walkObjects(c.content, fyne.NewPos(0, 0), paintObj)
}

func (c *glCanvas) setDirty() {
	// we must set twice as it's double buffered
	c.dirty1 = true
	c.dirty2 = true
}

func (c *glCanvas) isDirty() bool {
	return c.dirty1 || c.dirty2
}

func newCanvas(win *window) *glCanvas {
	c := &glCanvas{window: win, scale: 1.0}
	c.refreshQueue = make(chan fyne.CanvasObject, 1024)

	c.initOpenGL()

	return c
}
