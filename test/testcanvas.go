package test

import (
	"image"
	"image/draw"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var (
	dummyCanvas fyne.Canvas
)

type testCanvas struct {
	size  fyne.Size
	scale float32

	content, overlay fyne.CanvasObject
	focused          fyne.Focusable

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	fyne.ShortcutHandler
}

func (c *testCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *testCanvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	if content == nil {
		return
	}

	theme := fyne.CurrentApp().Settings().Theme()
	padding := fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
	c.Resize(content.MinSize().Add(padding))
}

func (c *testCanvas) Overlay() fyne.CanvasObject {
	return c.overlay
}

func (c *testCanvas) SetOverlay(overlay fyne.CanvasObject) {
	c.overlay = overlay
	if overlay != nil {
		overlay.Resize(c.Size())
	}
}

func (c *testCanvas) Refresh(fyne.CanvasObject) {
}

func (c *testCanvas) Focus(obj fyne.Focusable) {
	if obj == c.focused {
		return
	}

	if c.focused != nil {
		c.focused.FocusLost()
	}

	c.focused = obj

	if obj != nil {
		obj.FocusGained()
	}
}

func (c *testCanvas) Unfocus() {
	if c.focused != nil {
		c.focused.FocusLost()
	}
	c.focused = nil
}

func (c *testCanvas) Focused() fyne.Focusable {
	return c.focused
}

func (c *testCanvas) Size() fyne.Size {
	return c.size
}

func (c *testCanvas) Resize(size fyne.Size) {
	c.size = size

	c.content.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	c.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
}

func (c *testCanvas) Scale() float32 {
	return c.scale
}

func (c *testCanvas) SetScale(scale float32) {
	c.scale = scale
}

func (c *testCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *testCanvas) SetOnTypedRune(handler func(rune)) {
	c.onTypedRune = handler
}

func (c *testCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *testCanvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.onTypedKey = handler
}

func (c *testCanvas) Capture() image.Image {
	theme := fyne.CurrentApp().Settings().Theme()
	// TODO actually implement rendering

	bounds := image.Rect(0, 0, int(float32(c.Size().Width)*c.Scale()), int(float32(c.Size().Height)*c.Scale()))
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(theme.BackgroundColor()), image.ZP, draw.Src)

	return img
}

// NewCanvas returns a single use in-memory canvas used for testing
func NewCanvas() fyne.Canvas {
	theme := fyne.CurrentApp().Settings().Theme()
	padding := fyne.NewSize(theme.Padding(), theme.Padding())
	return &testCanvas{size: padding, scale: 1.0}
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() fyne.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvas()
	}

	return dummyCanvas
}
