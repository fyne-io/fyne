package test

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne"
)

var (
	dummyCanvas fyne.Canvas

	testBackground = color.RGBA{R: 255, G: 128, B: 255, A: 255}
)

type testCanvas struct {
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
	c.focused = obj
	obj.FocusGained()
}

func (c *testCanvas) Unfocus() {
	if c.focused != nil {
		c.focused.(fyne.Focusable).FocusLost()
	}
	c.focused = nil
}

func (c *testCanvas) Focused() fyne.Focusable {
	return c.focused
}

func (c *testCanvas) Size() fyne.Size {
	return fyne.NewSize(10, 10)
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
	// TODO actually implement rendering

	bounds := image.Rect(0, 0, int(float32(c.Size().Width)*c.Scale()), int(float32(c.Size().Height)*c.Scale()))
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(testBackground), image.ZP, draw.Src)

	return img
}

// NewCanvas returns a single use in-memory canvas used for testing
func NewCanvas() fyne.Canvas {
	return &testCanvas{}
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() fyne.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvas()
	}

	return dummyCanvas
}
