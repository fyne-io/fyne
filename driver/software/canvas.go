package software

import (
	"image"
	"image/draw"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

// WindowlessCanvas provides functionality for a canvas to operate without a window
//
// Since: 2.9
type WindowlessCanvas interface {
	fyne.Canvas

	Padded() bool
	Resize(fyne.Size)
	SetPadded(bool)
	SetScale(float32)
}

// NewCanvas creates a new canvas in memory that can render without hardware support.
func NewCanvas() WindowlessCanvas {
	return newCanvas(software.NewPainter(), false)
}

// NewCanvasWithPainter creates a new canvas in memory that can render without hardware support
// which uses the given driver.Painter for #Capture().
//
// Since: 2.8
func NewCanvasWithPainter(painter driver.Painter) WindowlessCanvas {
	return newCanvas(painter, false)
}

// NewTransparentCanvas creates a new canvas in memory that can render without hardware support without a background color.
//
// Since: 2.2
func NewTransparentCanvas() WindowlessCanvas {
	return newCanvas(software.NewPainter(), true)
}

// NewTransparentCanvasWithPainter creates a new canvas in memory that can render without hardware support
// which uses the given driver.Painter for #Capture() without a background color.
//
// Since: 2.8
func NewTransparentCanvasWithPainter(painter driver.Painter) WindowlessCanvas {
	return newCanvas(painter, true)
}

func newCanvas(painter driver.Painter, transparent bool) WindowlessCanvas {
	c := &canvas{
		focusMgr:    app.NewFocusManager(nil),
		padded:      true,
		painter:     painter,
		scale:       1.0,
		size:        fyne.NewSize(100, 100),
		transparent: transparent,
	}
	c.overlays.Canvas = c
	return c
}

type canvas struct {
	size    fyne.Size
	resized bool
	scale   float32

	content     fyne.CanvasObject
	overlays    internal.OverlayStack
	focusMgr    *app.FocusManager
	padded      bool
	transparent bool

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	fyne.ShortcutHandler
	painter      driver.Painter
	propertyLock sync.RWMutex
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
	return fyne.NewPos(0, 0), c.Size()
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

	return &c.overlays
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
	c.resized = true
	c.propertyLock.Unlock()

	c.doResize(size)
}

func (c *canvas) Scale() float32 {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.scale
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.propertyLock.Lock()
	c.content = content
	c.focusMgr = app.NewFocusManager(c.content)
	resized := c.resized
	c.propertyLock.Unlock()

	if content == nil {
		return
	}

	minSize := content.MinSize()
	if c.padded {
		minSize = minSize.Add(fyne.NewSquareSize(theme.Padding() * 2))
	}

	if resized {
		c.doResize(c.Size().Max(minSize))
	} else {
		c.doResize(minSize)
	}
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

	c.doResize(c.Size())
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

func (c *canvas) doResize(size fyne.Size) {
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
	fullPos, fullSize := c.InteractiveArea()
	for _, overlay := range overlays.List() {
		overlay.Move(fullPos)
		overlay.Resize(fullSize)
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

func (c *canvas) focusManager() *app.FocusManager {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	return c.focusMgr
}
