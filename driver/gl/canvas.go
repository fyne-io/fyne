package gl

import (
	"math"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type glCanvas struct {
	sync.RWMutex
	window                 *window
	content, overlay, menu fyne.CanvasObject
	focused                fyne.Focusable

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
	onKeyDown   func(*fyne.KeyEvent)
	onKeyUp     func(*fyne.KeyEvent)
	shortcut    fyne.ShortcutHandler

	program  uint32
	scale    float32
	texScale float32

	dirty        bool
	dirtyMutex   *sync.Mutex
	refreshQueue chan fyne.CanvasObject
}

func scaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(math.Round(float64(v) * float64(c.Scale())))
	}
}

func textureScaleInt(c *glCanvas, v int) int {
	if c.scale == 1.0 && c.texScale == 1.0 {
		return v
	}
	return int(math.Round(float64(v) * float64(c.scale*c.texScale)))
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
	c.RLock()
	retval := c.content
	c.RUnlock()
	return retval
}

func (c *glCanvas) SetContent(content fyne.CanvasObject) {
	c.Lock()
	c.content = content
	c.Unlock()

	var w, h = c.window.viewport.GetSize()

	xpad := theme.Padding()
	ypad := theme.Padding()
	if !c.window.Padded() {
		xpad = 0
		ypad = 0
	}
	menu := c.menuHeight()
	width := unscaleInt(c, int(w)) - xpad*2
	height := unscaleInt(c, int(h+menu)) - ypad*2

	c.content.Resize(fyne.NewSize(width, height))
	c.content.Move(fyne.NewPos(xpad, ypad+menu))
	c.setDirty(true)
}

func (c *glCanvas) Overlay() fyne.CanvasObject {
	c.RLock()
	retval := c.overlay
	c.RUnlock()
	return retval
}

func (c *glCanvas) SetOverlay(overlay fyne.CanvasObject) {
	c.Lock()
	c.overlay = overlay
	c.Unlock()

	if overlay != nil {
		c.overlay.Resize(c.Size())
	}
	c.setDirty(true)
}

func (c *glCanvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
	c.setDirty(true)
}

func (c *glCanvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		c.focused.(fyne.Focusable).FocusLost()
	}

	c.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (c *glCanvas) Unfocus() {
	if c.focused != nil {
		c.focused.(fyne.Focusable).FocusLost()
	}
	c.focused = nil
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
	c.setDirty(true)
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

func (c *glCanvas) OnKeyDown() func(*fyne.KeyEvent) {
	return c.onKeyDown
}

func (c *glCanvas) SetOnKeyDown(typed func(*fyne.KeyEvent)) {
	c.onKeyDown = typed
}

func (c *glCanvas) OnKeyUp() func(*fyne.KeyEvent) {
	return c.onKeyUp
}

func (c *glCanvas) SetOnKeyUp(typed func(*fyne.KeyEvent)) {
	c.onKeyUp = typed
}

func (c *glCanvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

func (c *glCanvas) paint(size fyne.Size) {
	if c.Content() == nil {
		return
	}
	c.setDirty(false)

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	paint := func(obj fyne.CanvasObject, pos fyne.Position) bool {
		// TODO should this be somehow not scroll container specific?
		if _, ok := obj.(*widget.ScrollContainer); ok {
			scrollX := textureScaleInt(c, pos.X)
			scrollY := textureScaleInt(c, pos.Y)
			scrollWidth := textureScaleInt(c, obj.Size().Width)
			scrollHeight := textureScaleInt(c, obj.Size().Height)
			_, pixHeight := c.window.viewport.GetFramebufferSize()
			gl.Scissor(int32(scrollX), int32(pixHeight-scrollY-scrollHeight), int32(scrollWidth), int32(scrollHeight))
			gl.Enable(gl.SCISSOR_TEST)
		}
		if obj.Visible() {
			c.drawObject(obj, pos, size)
		}
		return false
	}
	afterPaint := func(obj fyne.CanvasObject, pos fyne.Position, _ bool) {
		if _, ok := obj.(*widget.ScrollContainer); ok {
			gl.Disable(gl.SCISSOR_TEST)
		}
	}

	driver.WalkObjectTree(c.content, fyne.NewPos(0, 0), paint, afterPaint)
	if c.menu != nil {
		driver.WalkObjectTree(c.menu, fyne.NewPos(0, 0), paint, afterPaint)
	}
	if c.overlay != nil {
		driver.WalkObjectTree(c.overlay, fyne.NewPos(0, 0), paint, afterPaint)
	}
}

func (c *glCanvas) setDirty(dirty bool) {
	c.dirtyMutex.Lock()
	defer c.dirtyMutex.Unlock()

	c.dirty = dirty
}

func (c *glCanvas) isDirty() bool {
	c.dirtyMutex.Lock()
	defer c.dirtyMutex.Unlock()

	return c.dirty
}

func newCanvas(win *window) *glCanvas {
	c := &glCanvas{window: win, scale: 1.0}
	c.content = &canvas.Rectangle{FillColor: theme.BackgroundColor()}
	c.refreshQueue = make(chan fyne.CanvasObject, 1024)
	c.dirtyMutex = &sync.Mutex{}

	c.initOpenGL()

	return c
}
