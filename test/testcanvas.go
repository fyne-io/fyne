package test

import (
	"image"
	"image/draw"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/theme"
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

	content  fyne.CanvasObject
	overlays *internal.OverlayStack
	focused  fyne.Focusable
	hovered  desktop.Hoverable
	padded   bool

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	fyne.ShortcutHandler
	painter      SoftwarePainter
	propertyLock sync.RWMutex
}

func (c *testCanvas) Content() fyne.CanvasObject {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.content
}

func (c *testCanvas) SetContent(content fyne.CanvasObject) {
	c.propertyLock.Lock()
	c.content = content
	c.propertyLock.Unlock()

	if content == nil {
		return
	}

	padding := fyne.NewSize(0, 0)
	if c.padded {
		padding = fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
	}
	c.Resize(content.MinSize().Add(padding))
}

// Deprecated
func (c *testCanvas) Overlay() fyne.CanvasObject {
	panic("deprecated method should not be used")
}

func (c *testCanvas) Overlays() fyne.OverlayStack {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.overlays
}

// Deprecated
func (c *testCanvas) SetOverlay(_ fyne.CanvasObject) {
	panic("deprecated method should not be used")
}

func (c *testCanvas) Refresh(fyne.CanvasObject) {
}

func (c *testCanvas) Focus(obj fyne.Focusable) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

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
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if c.focused != nil {
		c.focused.FocusLost()
	}
	c.focused = nil
}

func (c *testCanvas) Focused() fyne.Focusable {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.focused
}

func (c *testCanvas) Size() fyne.Size {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.size
}

func (c *testCanvas) Resize(size fyne.Size) {
	c.propertyLock.Lock()
	content := c.content
	overlays := c.overlays
	padded := c.padded
	c.size = size
	c.propertyLock.Unlock()

	if content == nil {
		return
	}

	for _, overlay := range overlays.List() {
		overlay.Resize(size)
	}

	if padded {
		theme := fyne.CurrentApp().Settings().Theme()
		content.Resize(size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	} else {
		content.Resize(size)
		content.Move(fyne.NewPos(0, 0))
	}
}

func (c *testCanvas) Scale() float32 {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.scale
}

func (c *testCanvas) SetScale(scale float32) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.scale = scale
}

func (c *testCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *testCanvas) OnTypedRune() func(rune) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedRune
}

func (c *testCanvas) SetOnTypedRune(handler func(rune)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedRune = handler
}

func (c *testCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedKey
}

func (c *testCanvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedKey = handler
}

func (c *testCanvas) Padded() bool {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.padded
}

func (c *testCanvas) SetPadded(padded bool) {
	c.propertyLock.Lock()
	c.padded = padded
	c.propertyLock.Unlock()

	c.Resize(c.Size())
}

func (c *testCanvas) Capture() image.Image {
	if c.painter != nil {
		return c.painter.Paint(c)
	}
	theme := fyne.CurrentApp().Settings().Theme()

	bounds := image.Rect(0, 0, internal.ScaleInt(c, c.Size().Width), internal.ScaleInt(c, c.Size().Height))
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(theme.BackgroundColor()), image.Point{}, draw.Src)

	return img
}

// LaidOutObjects returns all fyne.CanvasObject starting at the given fyne.CanvasObject which is laid out previously.
func LaidOutObjects(o fyne.CanvasObject) (objects []fyne.CanvasObject) {
	if o != nil {
		objects = layoutAndCollect(objects, o, o.MinSize().Max(o.Size()))
	}
	return objects
}

func layoutAndCollect(objects []fyne.CanvasObject, o fyne.CanvasObject, size fyne.Size) []fyne.CanvasObject {
	objects = append(objects, o)
	switch c := o.(type) {
	case fyne.Widget:
		r := c.CreateRenderer()
		r.Layout(size)
		for _, child := range r.Objects() {
			objects = layoutAndCollect(objects, child, child.Size())
		}
	case *fyne.Container:
		if c.Layout != nil {
			c.Layout.Layout(c.Objects, size)
		}
		for _, child := range c.Objects {
			objects = layoutAndCollect(objects, child, child.Size())
		}
	}
	return objects
}

func (c *testCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+1)
	if c.content != nil {
		trees = append(trees, c.content)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

// NewCanvas returns a single use in-memory canvas used for testing
func NewCanvas() WindowlessCanvas {
	padding := fyne.NewSize(10, 10)
	return &testCanvas{size: padding, padded: true, scale: 1.0, overlays: &internal.OverlayStack{}}
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
