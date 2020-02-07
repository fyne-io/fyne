package test

import (
	"image"
	"image/draw"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
)

var (
	dummyCanvas fyne.Canvas
)

// WindowlessCanvas provides functionality for a canvas to operate without a window
type WindowlessCanvas interface {
	fyne.Canvas

	Resize(fyne.Size)

	Padded() bool
	SetPadded(bool)
}

type testCanvas struct {
	size  fyne.Size
	scale float32

	content, overlay fyne.CanvasObject
	focused          fyne.Focusable
	padded           bool

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	fyne.ShortcutHandler
	painter SoftwarePainter
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
	if c.content == nil {
		return
	}

	if c.padded {
		theme := fyne.CurrentApp().Settings().Theme()
		c.content.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		c.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	} else {
		c.content.Resize(size)
	}
}

func (c *testCanvas) Scale() float32 {
	return c.scale
}

func (c *testCanvas) SetScale(scale float32) {
	c.scale = scale
}

func (c *testCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
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

func (c *testCanvas) Padded() bool {
	return c.padded
}

func (c *testCanvas) SetPadded(padded bool) {
	c.padded = padded
	c.Resize(c.Size())
}

func (c *testCanvas) Capture() image.Image {
	if c.painter != nil {
		return c.painter.Paint(c)
	}
	theme := fyne.CurrentApp().Settings().Theme()

	bounds := image.Rect(0, 0, internal.ScaleInt(c, c.Size().Width), internal.ScaleInt(c, c.Size().Height))
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(theme.BackgroundColor()), image.ZP, draw.Src)

	return img
}

// NewCanvas returns a single use in-memory canvas used for testing
func NewCanvas() WindowlessCanvas {
	padding := fyne.NewSize(10, 10)
	return &testCanvas{size: padding, padded: true, scale: 1.0}
}

// NewCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	canvas := NewCanvas().(*testCanvas)
	canvas.painter = painter

	return canvas
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() fyne.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvas()
	}

	return dummyCanvas
}
