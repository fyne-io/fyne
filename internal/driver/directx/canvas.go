//go:build windows && directx

package directx

import (
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter/dx"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with the Canvas interfaces.
var (
	_ fyne.Canvas           = (*dxCanvas)(nil)
	_ common.SizeableCanvas = (*dxCanvas)(nil)
	_ dx.Repainter          = (*dxCanvas)(nil)
)

type dxCanvas struct {
	common.Canvas

	content fyne.CanvasObject
	padded  bool
	size    fyne.Size

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
	onKeyDown   func(*fyne.KeyEvent)
	onKeyUp     func(*fyne.KeyEvent)

	scale, texScale float32

	win *window
}

func newCanvas() *dxCanvas {
	c := &dxCanvas{scale: 1.0, texScale: 1.0, padded: true}
	c.Initialize(c, c.overlayChanged)
	c.setContent(&canvas.Rectangle{FillColor: theme.Color(theme.ColorNameBackground)})
	return c
}

func (c *dxCanvas) Capture() image.Image {
	return c.Painter().Capture(c)
}

func (c *dxCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *dxCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

func (c *dxCanvas) MinSize() fyne.Size {
	return c.canvasSize(c.content.MinSize())
}

func (c *dxCanvas) OnKeyDown() func(*fyne.KeyEvent) { return c.onKeyDown }
func (c *dxCanvas) OnKeyUp() func(*fyne.KeyEvent)   { return c.onKeyUp }
func (c *dxCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}
func (c *dxCanvas) OnTypedRune() func(rune) { return c.onTypedRune }

func (c *dxCanvas) SetOnKeyDown(f func(*fyne.KeyEvent))  { c.onKeyDown = f }
func (c *dxCanvas) SetOnKeyUp(f func(*fyne.KeyEvent))    { c.onKeyUp = f }
func (c *dxCanvas) SetOnTypedKey(f func(*fyne.KeyEvent)) { c.onTypedKey = f }
func (c *dxCanvas) SetOnTypedRune(f func(rune))          { c.onTypedRune = f }

func (c *dxCanvas) Padded() bool { return c.padded }

func (c *dxCanvas) SetPadded(padded bool) {
	c.padded = padded
	c.content.Move(c.contentPos())
}

func (c *dxCanvas) PixelCoordinateForPosition(pos fyne.Position) (x, y int) {
	multiple := c.scale * c.texScale
	scaleInt := func(v float32) int {
		return int(math.Round(float64(v * multiple)))
	}
	return scaleInt(pos.X), scaleInt(pos.Y)
}

func (c *dxCanvas) Resize(size fyne.Size) {
	// At scale 1 the requested size may be fractional; round up so the canvas
	// covers the window fully and the first frame stays sharp.
	nearestSize := size
	if c.scale == 1 {
		nearestSize = fyne.NewSize(float32(math.Ceil(float64(size.Width))),
			float32(math.Ceil(float64(size.Height))))
	}
	c.size = nearestSize

	for _, overlay := range c.Overlays().List() {
		overlay.Resize(nearestSize)
	}

	c.content.Resize(c.contentSize(nearestSize))
	c.content.Move(c.contentPos())
}

func (c *dxCanvas) Scale() float32 { return c.scale }

func (c *dxCanvas) SetContent(content fyne.CanvasObject) {
	newSize := c.size.Max(c.canvasSize(content.MinSize()))
	c.setContent(content)
	c.Resize(newSize)
	c.SetDirty()
}

func (c *dxCanvas) Size() fyne.Size { return c.size }

func (c *dxCanvas) setContent(content fyne.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

// canvasSize computes the canvas size needed to hold the given content.
func (c *dxCanvas) canvasSize(contentSize fyne.Size) fyne.Size {
	if c.Padded() {
		return contentSize.Add(fyne.NewSquareSize(theme.Padding() * 2))
	}
	return contentSize
}

func (c *dxCanvas) contentPos() fyne.Position {
	if c.Padded() {
		return fyne.NewSquareOffsetPos(theme.Padding())
	}
	return fyne.Position{}
}

func (c *dxCanvas) contentSize(canvasSize fyne.Size) fyne.Size {
	if c.Padded() {
		return canvasSize.Subtract(fyne.NewSquareSize(theme.Padding() * 2))
	}
	return canvasSize
}

func (c *dxCanvas) overlayChanged() {
	c.SetDirty()
}

func (c *dxCanvas) reloadScale() {
	if c.win == nil || !c.win.visible {
		return
	}
	c.scale = c.win.calculatedScale()
	c.SetDirty()
	c.win.rescale()
}

func (c *dxCanvas) applyThemeOutOfTreeObjects() {
	c.SetPadded(c.padded)
}

// Repaint satisfies dx.Repainter: the painter calls it from Capture, because the
// swap chain discards the back buffer on Present and there is nothing valid left
// to read back.
func (c *dxCanvas) Repaint(size fyne.Size) {
	c.paint(size)
}

// paint walks the render tree and draws every object, pushing and popping
// scissor rectangles as clipping containers are entered and left.
func (c *dxCanvas) paint(size fyne.Size) {
	if c.Content() == nil {
		return
	}
	clips := &internal.ClipStack{}
	c.Painter().Clear()

	paint := func(node *common.RenderCacheNode, pos fyne.Position) {
		obj := node.Obj()
		if driver.IsClip(obj) {
			inner := clips.Push(pos, obj.Size())
			c.Painter().StartClipping(inner.Rect())
		}
		if size.Width <= 0 || size.Height <= 0 { // minimising can report a zero size
			return
		}
		c.Painter().Paint(obj, pos, size, clips.Top())
	}
	afterPaint := func(node *common.RenderCacheNode, pos fyne.Position) {
		if driver.IsClip(node.Obj()) {
			clips.Pop()
			if top := clips.Top(); top != nil {
				c.Painter().StartClipping(top.Rect())
			} else {
				c.Painter().StopClipping()
			}
		}
		if build.Mode == fyne.BuildDebug {
			c.DrawDebugOverlay(node.Obj(), pos, size, clips.Top())
		}
	}
	c.WalkTrees(paint, afterPaint)
	// The painter batches glyph quads, so the tail of the tree can still be
	// sitting in a queue when the walk ends.
	if p, ok := c.Painter().(*dx.Painter); ok {
		p.Flush()
	}
}

// applyTheme re-applies the current theme to objects that live outside the
// content tree.
func (c *dxCanvas) applyTheme() {
	app.ApplyThemeTo(c.content, c)
}
