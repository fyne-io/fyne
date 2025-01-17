package test

import (
	"image"
	"image/draw"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

var dummyCanvas WindowlessCanvas

// WindowlessCanvas provides functionality for a canvas to operate without a window
type WindowlessCanvas interface {
	fyne.Canvas

	Padded() bool
	Resize(fyne.Size)
	SetPadded(bool)
	SetScale(float32)
}

type canvas struct {
	size  fyne.Size
	scale float32

	content     fyne.CanvasObject
	overlays    *internal.OverlayStack
	focusMgr    *intapp.FocusManager
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
	c := &canvas{
		focusMgr: intapp.NewFocusManager(nil),
		padded:   true,
		scale:    1.0,
		size:     fyne.NewSize(100, 100),
	}
	c.overlays = &internal.OverlayStack{Canvas: c}
	return c
}

// NewCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	c := NewCanvas().(*canvas)
	c.painter = painter

	return c
}

// NewTransparentCanvasWithPainter allows creation of an in-memory canvas with a specific painter without a background color.
// The painter will be used to render in the Capture() call.
//
// Since: 2.2
func NewTransparentCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	c := NewCanvasWithPainter(painter).(*canvas)
	c.transparent = true

	return c
}

func (c *canvas) Capture() image.Image {
	cache.Clean(true)
	size := c.Size()
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, size.Width), scale.ToScreenCoordinate(c, size.Height))
	img := image.NewNRGBA(bounds)
	if !c.transparent {
		draw.Draw(img, bounds, image.NewUniform(theme.Color(theme.ColorNameBackground)), image.Point{}, draw.Src)
	}

	if c.painter != nil {
		draw.Draw(img, bounds, c.painter.Paint(c), image.Point{}, draw.Over)
	}

	return img
}

func (c *canvas) Content() fyne.CanvasObject {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.content
}

func (c *canvas) Focus(obj fyne.Focusable) {
	c.focusManager().Focus(obj)
}

func (c *canvas) FocusNext() {
	c.focusManager().FocusNext()
}

func (c *canvas) FocusPrevious() {
	c.focusManager().FocusPrevious()
}

func (c *canvas) Focused() fyne.Focusable {
	return c.focusManager().Focused()
}

func (c *canvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.NewPos(2, 3), c.Size().SubtractWidthHeight(4, 5)
}

func (c *canvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedKey
}

func (c *canvas) OnTypedRune() func(rune) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedRune
}

func (c *canvas) Overlays() fyne.OverlayStack {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.overlays
}

func (c *canvas) Padded() bool {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.padded
}

func (c *canvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(pos.X * c.scale), int(pos.Y * c.scale)
}

func (c *canvas) Refresh(fyne.CanvasObject) {
}

func (c *canvas) Resize(size fyne.Size) {
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
		padding := theme.Padding()
		content.Resize(size.Subtract(fyne.NewSquareSize(padding * 2)))
		content.Move(fyne.NewSquareOffsetPos(padding))
	} else {
		content.Resize(size)
		content.Move(fyne.NewPos(0, 0))
	}
}

func (c *canvas) Scale() float32 {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.scale
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.propertyLock.Lock()
	c.content = content
	c.focusMgr = intapp.NewFocusManager(c.content)
	c.propertyLock.Unlock()

	if content == nil {
		return
	}

	minSize := content.MinSize()
	if c.padded {
		minSize = minSize.Add(fyne.NewSquareSize(theme.Padding() * 2))
	}
	c.Resize(minSize)
}

func (c *canvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedKey = handler
}

func (c *canvas) SetOnTypedRune(handler func(rune)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedRune = handler
}

func (c *canvas) SetPadded(padded bool) {
	c.propertyLock.Lock()
	c.padded = padded
	c.propertyLock.Unlock()

	c.Resize(c.Size())
}

func (c *canvas) SetScale(scale float32) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.scale = scale
}

func (c *canvas) Size() fyne.Size {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.size
}

func (c *canvas) Unfocus() {
	c.focusManager().Focus(nil)
}

func (c *canvas) focusManager() *intapp.FocusManager {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	return c.focusMgr
}

func (c *canvas) objectTrees() []fyne.CanvasObject {
	overlays := c.Overlays().List()
	trees := make([]fyne.CanvasObject, 0, len(overlays)+1)
	if c.content != nil {
		trees = append(trees, c.content)
	}
	trees = append(trees, overlays...)
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
