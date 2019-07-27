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
}

func (c *canvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	content.Resize(c.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
}

func (c *canvas) Refresh(fyne.CanvasObject) {
	//	panic("implement me")
}

func (c *canvas) Resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.size = size
	c.content.Resize(c.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	c.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
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
	c.scale = scale
}

func (c *canvas) Overlay() fyne.CanvasObject {
	return c.overlay
}

func (c *canvas) SetOverlay(overlay fyne.CanvasObject) {
	c.overlay = overlay
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
