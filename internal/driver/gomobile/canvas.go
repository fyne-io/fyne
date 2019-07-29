package gomobile

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
)

type canvas struct {
	content, overlay fyne.CanvasObject
	painter          gl.Painter
	scale            float32
	size             fyne.Size

	dirty        bool
	refreshQueue chan fyne.CanvasObject
}

func (c *canvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	content.Resize(c.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
}

func (c *canvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
	c.dirty = true
}

func (c *canvas) Resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.size = size
	c.content.Resize(c.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	c.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	if c.overlay != nil {
		c.overlay.Resize(size)
	}
}

func (c *canvas) Focus(fyne.Focusable) {
	panic("implement me")
}

func (c *canvas) Unfocus() {
	panic("implement me")
}

func (c *canvas) Focused() fyne.Focusable {
	panic("implement me")
}

func (c *canvas) Size() fyne.Size {
	return c.size
}

func (c *canvas) Scale() float32 {
	return c.scale
}

func (c *canvas) SetScale(scale float32) {
	if scale == fyne.SettingsScaleAuto {
		scale = c.detectScale()
	}
	c.scale = scale
}

func (c *canvas) detectScale() float32 {
	return 2 // TODO real detection
}

func (c *canvas) Overlay() fyne.CanvasObject {
	return c.overlay
}

func (c *canvas) SetOverlay(overlay fyne.CanvasObject) {
	c.overlay = overlay

	if c.overlay != nil {
		c.overlay.Resize(c.size)
	}
}

func (c *canvas) OnTypedRune() func(rune) {
	panic("implement me")
}

func (c *canvas) SetOnTypedRune(func(rune)) {
	//	panic("implement me")
}

func (c *canvas) OnTypedKey() func(*fyne.KeyEvent) {
	panic("implement me")
}

func (c *canvas) SetOnTypedKey(func(*fyne.KeyEvent)) {
	//	panic("implement me")
}

func (c *canvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	panic("implement me")
}

func (c *canvas) Capture() image.Image {
	return c.painter.Capture(c)
}

// NewCanvas creates a new gomobile canvas. This is a canvas that will render on a mobile device using OpenGL.
func NewCanvas() fyne.Canvas {
	ret := &canvas{}
	ret.scale = ret.detectScale()
	ret.refreshQueue = make(chan fyne.CanvasObject, 1024)

	return ret
}
