package test

import (
	"image"
	"image/draw"

	"fyne.io/fyne"
)

var (
	dummyCanvas fyne.Canvas
)

type testCanvas struct {
	size fyne.Size

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
}

func (c *testCanvas) Scale() float32 {
	return 1.0
}

func (c *testCanvas) SetScale(float32) {
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
	return &testCanvas{size: fyne.NewSize(10, 10)}
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() fyne.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvas()
	}

	return dummyCanvas
}
