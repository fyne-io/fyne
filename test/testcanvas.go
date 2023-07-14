package test

import (
	"image"
	"image/draw"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

var (
	dummyCanvas fyne.Canvas
)

// WindowlessCanvas provides functionality for a canvas to operate without a window
type WindowlessCanvas interface {
	fyne.Canvas

	Padded() bool
	Resize(fyne.Size)
	SetPadded(bool)
	SetScale(float32)
}

type testCanvas struct {
	size  fyne.Size
	scale float32

	content     fyne.CanvasObject
	overlays    *internal.OverlayStack
	focusMgr    *app.FocusManager
	hovered     desktop.Hoverable
	padded      bool
	transparent bool

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	fyne.ShortcutHandler
	painter      SoftwarePainter
	propertyLock sync.RWMutex
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() fyne.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvas()
	}

	return dummyCanvas
}

// NewCanvas returns a single use in-memory canvas used for testing.
// This canvas has no painter so calls to Capture() will return a blank image.
func NewCanvas() WindowlessCanvas {
	c := &testCanvas{
		focusMgr: app.NewFocusManager(nil),
		padded:   true,
		scale:    1.0,
		size:     fyne.NewSize(10, 10),
	}
	c.overlays = &internal.OverlayStack{Canvas: c}
	return c
}

// NewCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	canvas := NewCanvas().(*testCanvas)
	canvas.painter = painter

	return canvas
}

// NewTransparentCanvasWithPainter allows creation of an in-memory canvas with a specific painter without a background color.
// The painter will be used to render in the Capture() call.
//
// Since: 2.2
func NewTransparentCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	canvas := NewCanvasWithPainter(painter).(*testCanvas)
	canvas.transparent = true

	return canvas
}

func (c *testCanvas) Capture() image.Image {
	cache.Clean(true)
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, c.Size().Width), scale.ToScreenCoordinate(c, c.Size().Height))
	img := image.NewNRGBA(bounds)
	if !c.transparent {
		draw.Draw(img, bounds, image.NewUniform(theme.BackgroundColor()), image.Point{}, draw.Src)
	}

	if c.painter != nil {
		draw.Draw(img, bounds, c.painter.Paint(c), image.Point{}, draw.Over)
	}

	return img
}

func (c *testCanvas) Content() fyne.CanvasObject {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.content
}

func (c *testCanvas) Focus(obj fyne.Focusable) {
	c.focusManager().Focus(obj)
}

func (c *testCanvas) FocusNext() {
	c.focusManager().FocusNext()
}

func (c *testCanvas) FocusPrevious() {
	c.focusManager().FocusPrevious()
}

func (c *testCanvas) Focused() fyne.Focusable {
	return c.focusManager().Focused()
}

func (c *testCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

func (c *testCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedKey
}

func (c *testCanvas) OnTypedRune() func(rune) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedRune
}

func (c *testCanvas) Overlays() fyne.OverlayStack {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.overlays
}

func (c *testCanvas) Padded() bool {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.padded
}

func (c *testCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *testCanvas) Refresh(fyne.CanvasObject) {
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

	// Ensure testcanvas mimics real canvas.Resize behavior
	for _, overlay := range overlays.List() {
		type popupWidget interface {
			fyne.CanvasObject
			ShowAtPosition(fyne.Position)
		}
		if p, ok := overlay.(popupWidget); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(size)
		}
	}

	if padded {
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

func (c *testCanvas) SetContent(content fyne.CanvasObject) {
	c.propertyLock.Lock()
	c.content = content
	c.focusMgr = app.NewFocusManager(c.content)
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

func (c *testCanvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedKey = handler
}

func (c *testCanvas) SetOnTypedRune(handler func(rune)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedRune = handler
}

func (c *testCanvas) SetPadded(padded bool) {
	c.propertyLock.Lock()
	c.padded = padded
	c.propertyLock.Unlock()

	c.Resize(c.Size())
}

func (c *testCanvas) SetScale(scale float32) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.scale = scale
}

func (c *testCanvas) Size() fyne.Size {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.size
}

func (c *testCanvas) Unfocus() {
	c.focusManager().Focus(nil)
}

func (c *testCanvas) focusManager() *app.FocusManager {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	return c.focusMgr
}

func (c *testCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+1)
	if c.content != nil {
		trees = append(trees, c.content)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
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
