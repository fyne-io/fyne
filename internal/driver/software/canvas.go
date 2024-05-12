package software

import (
	"context"
	"image"
	"image/draw"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

var (
	dummyCanvas WindowlessCanvas
)

// WindowlessCanvas provides functionality for a canvas to operate without a window
type WindowlessCanvas interface {
	common.SizeableCanvas

	Padded() bool
	SetPadded(bool)
	SetScale(float32)

	// Initialize(common.SizeableCanvas, func())
	Clear()
}

type SoftwareCanvas struct {
	common.Canvas

	size  fyne.Size
	scale float32

	content     fyne.CanvasObject
	focusMgr    *app.FocusManager
	hovered     desktop.Hoverable
	padded      bool
	transparent bool

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

	touched map[int]mobile.Touchable

	lastTapDown           map[int]time.Time
	lastTapDownPos        map[int]fyne.Position
	dragging              fyne.Draggable
	dragStart, dragOffset fyne.Position

	touchTapCount   int
	touchCancelFunc context.CancelFunc
	touchLastTapped fyne.CanvasObject
}

func (c *SoftwareCanvas) MinSize() fyne.Size {
	// TODO implement me
	panic("implement me")
}

// Canvas returns a reusable in-memory canvas used for testing
func Canvas() WindowlessCanvas {
	if dummyCanvas == nil {
		dummyCanvas = NewCanvasWithPainter(software.NewPainter())
		dummyCanvas.(*SoftwareCanvas).Initialize(dummyCanvas, nil)
	}

	return dummyCanvas
}

// NewCanvas returns a single use in-memory canvas used for testing.
// This canvas has no painter so calls to Capture() will return a blank image.
func NewCanvas() WindowlessCanvas {
	c := &SoftwareCanvas{
		focusMgr: app.NewFocusManager(nil),
		padded:   true,
		scale:    1.0,
		size:     fyne.NewSize(10, 10),

		touched:        make(map[int]mobile.Touchable),
		lastTapDown:    make(map[int]time.Time),
		lastTapDownPos: make(map[int]fyne.Position),
	}
	c.Initialize(c, nil)
	return c
}

// NewCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewCanvasWithPainter(painter painter.Painter) WindowlessCanvas {
	canvas := NewCanvas().(*SoftwareCanvas)
	canvas.SetPainter(painter)

	return canvas
}

// NewTransparentCanvasWithPainter allows creation of an in-memory canvas with a specific painter without a background color.
// The painter will be used to render in the Capture() call.
//
// Since: 2.2
func NewTransparentCanvasWithPainter(painter painter.Painter) WindowlessCanvas {
	canvas := NewCanvasWithPainter(painter).(*SoftwareCanvas)
	canvas.transparent = true

	return canvas
}

func (c *SoftwareCanvas) Capture() image.Image {
	cache.Clean(true)
	bounds := image.Rect(0, 0, scale.ToScreenCoordinate(c, c.Size().Width), scale.ToScreenCoordinate(c, c.Size().Height))
	c.FreeDirtyTextures()
	var img *image.NRGBA
	if !c.transparent {
		// TODO: this makes the cache useless, we need to find a better way to let the painter know we want it to paint everything
		c.Clear()
		img = image.NewNRGBA(bounds)
		// TODO: this is slow, and is slower if the bg color is not color.NRGBA
		draw.Draw(img, bounds, image.NewUniform(theme.BackgroundColor()), image.Point{}, draw.Src)
	}

	if c.Painter() != nil {
		x := c.Painter().Capture(c)
		if c.transparent {
			return x
		}
		// TODO: it's slow (and somewhat useless) to draw the image twice (once here, twice the fb).
		//  Not sure what's a good solution here cause we do care about the bg (or do we?)
		draw.Draw(img, bounds, x, image.Point{}, draw.Over)
	}

	return img
}

func (c *SoftwareCanvas) Clear() {
	driver.WalkCompleteObjectTree(c.content, func(obj fyne.CanvasObject, _, _ fyne.Position, _ fyne.Size) bool {
		cache.DeleteTexture(obj)
		return false
	}, nil)
	if c.Painter() != nil {
		c.Painter().Clear()
	}
}

func (c *SoftwareCanvas) Content() fyne.CanvasObject {
	c.RLock()
	defer c.RUnlock()

	return c.content
}

func (c *SoftwareCanvas) Hovered() desktop.Hoverable {
	c.RLock()
	defer c.RUnlock()

	return c.hovered
}

func (c *SoftwareCanvas) IsTransparent() bool {
	c.RLock()
	defer c.RUnlock()

	return c.transparent
}

func (c *SoftwareCanvas) SetTransparent(transparent bool) {
	c.Lock()
	c.transparent = transparent
	c.Unlock()
}

func (c *SoftwareCanvas) SetHovered(hovered desktop.Hoverable) {
	c.Lock()
	defer c.Unlock()

	c.hovered = hovered
}

func (c *SoftwareCanvas) findObjectAtPositionMatching(pos fyne.Position, test func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	return driver.FindObjectAtPositionMatching(pos, test, c.Overlays().Top(), c.content)
}

func (c *SoftwareCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

func (c *SoftwareCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.RLock()
	defer c.RUnlock()

	return c.onTypedKey
}

func (c *SoftwareCanvas) OnTypedRune() func(rune) {
	c.RLock()
	defer c.RUnlock()

	return c.onTypedRune
}

func (c *SoftwareCanvas) Padded() bool {
	c.RLock()
	defer c.RUnlock()

	return c.padded
}

func (c *SoftwareCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *SoftwareCanvas) Resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.Lock()
	content := c.content
	overlays := c.Overlays()
	padded := c.padded
	c.size = size
	c.Unlock()

	if content == nil {
		return
	}

	// Ensure SoftwareCanvas mimics real canvas.Resize behavior
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
			// overlay.Move(fyne.NewSquareOffsetPos(theme.Padding()))
		}
	}

	if padded {
		content.Resize(size.Subtract(fyne.NewSquareSize(theme.Padding() * 2)))
		content.Move(fyne.NewSquareOffsetPos(theme.Padding()))
	} else {
		content.Resize(size)
		content.Move(fyne.NewPos(0, 0))
	}
}

func (c *SoftwareCanvas) Scale() float32 {
	c.RLock()
	defer c.RUnlock()

	return c.scale
}

func (c *SoftwareCanvas) SetContent(content fyne.CanvasObject) {
	c.Lock()
	// newSize := c.size.Max(c.canvasSize(content.MinSize()))
	c.setContent(content)
	c.Unlock()

	if content == nil {
		return
	}

	padding := fyne.NewSize(0, 0)
	if c.padded {
		padding = fyne.NewSquareSize(theme.Padding() * 2)
	}
	c.Resize(content.MinSize().Add(padding))
	// c.Resize(newSize)
	c.SetDirty()
}

func (c *SoftwareCanvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.Lock()
	defer c.Unlock()

	c.onTypedKey = handler
}

func (c *SoftwareCanvas) SetOnTypedRune(handler func(rune)) {
	c.Lock()
	defer c.Unlock()

	c.onTypedRune = handler
}

func (c *SoftwareCanvas) SetPadded(padded bool) {
	c.Lock()
	c.padded = padded
	c.Unlock()

	c.Resize(c.Size())
}

func (c *SoftwareCanvas) SetScale(scale float32) {
	c.Lock()
	defer c.Unlock()

	c.scale = scale
}

func (c *SoftwareCanvas) Size() fyne.Size {
	c.RLock()
	defer c.RUnlock()

	return c.size
}

// canvasSize computes the needed canvas size for the given content size
// func (c *SoftwareCanvas) canvasSize(contentSize fyne.Size) fyne.Size {
// 	canvasSize := contentSize.Add(fyne.NewSize(0, 0))
// 	if c.padded {
// 		return canvasSize.Add(fyne.NewSquareSize(theme.Padding() * 2))
// 	}
// 	return canvasSize
// }

func (c *SoftwareCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+1)
	if c.content != nil {
		trees = append(trees, c.content)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

func (c *SoftwareCanvas) setContent(content fyne.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

func (c *SoftwareCanvas) tapDown(pos fyne.Position, tapID int) {
	c.lastTapDown[tapID] = time.Now()
	c.lastTapDownPos[tapID] = pos
	c.dragging = nil

	co, objPos, layer := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case mobile.Touchable, fyne.Focusable:
			return true
		}

		return false
	})

	if wid, ok := co.(mobile.Touchable); ok {
		touchEv := &mobile.TouchEvent{}
		touchEv.Position = objPos
		touchEv.AbsolutePosition = pos
		wid.TouchDown(touchEv)
		c.touched[tapID] = wid
	}

	if layer != 1 { // 0 - overlay, 1 - window head / menu, 2 - content
		if wid, ok := co.(fyne.Focusable); !ok || wid != c.Focused() {
			c.Unfocus()
		}
	}
}

const tapMoveThreshold = 4.0

func (c *SoftwareCanvas) tapMove(pos fyne.Position, tapID int,
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	previousPos := c.lastTapDownPos[tapID]
	deltaX := pos.X - previousPos.X
	deltaY := pos.Y - previousPos.Y

	if c.dragging == nil && (math.Abs(float64(deltaX)) < tapMoveThreshold && math.Abs(float64(deltaY)) < tapMoveThreshold) {
		return
	}
	c.lastTapDownPos[tapID] = pos

	co, objPos, _ := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Draggable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		}

		return false
	})

	if c.touched[tapID] != nil {
		if touch, ok := co.(mobile.Touchable); !ok || c.touched[tapID] != touch {
			touchEv := &mobile.TouchEvent{}
			touchEv.Position = objPos
			touchEv.AbsolutePosition = pos
			c.touched[tapID].TouchCancel(touchEv)
			c.touched[tapID] = nil
		}
	}

	if c.dragging == nil {
		if drag, ok := co.(fyne.Draggable); ok {
			c.dragging = drag
			c.dragOffset = previousPos.Subtract(objPos)
			c.dragStart = co.Position()
		} else {
			return
		}
	}

	ev := &fyne.DragEvent{}
	draggedObjDelta := c.dragStart.Subtract(c.dragging.(fyne.CanvasObject).Position())
	ev.Position = pos.Subtract(c.dragOffset).Add(draggedObjDelta)
	ev.Dragged = fyne.Delta{DX: deltaX, DY: deltaY}

	dragCallback(c.dragging, ev)
}

const tapSecondaryDelay = 300 * time.Millisecond

func (c *SoftwareCanvas) tapUp(pos fyne.Position, tapID int,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.SecondaryTappable, *fyne.PointEvent),
	doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable)) {

	if c.dragging != nil {
		dragCallback(c.dragging)

		c.dragging = nil
		return
	}

	duration := time.Since(c.lastTapDown[tapID])

	// if c.menu != nil && c.Overlays().Top() == nil && pos.X > c.menu.Size().Width {
	// 	c.menu.Hide()
	// 	c.menu.Refresh()
	// 	c.setMenu(nil)
	// 	return
	// }

	co, objPos, _ := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.SecondaryTappable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		} else if _, ok := object.(fyne.DoubleTappable); ok {
			return true
		}

		return false
	})

	if wid, ok := co.(mobile.Touchable); ok {
		touchEv := &mobile.TouchEvent{}
		touchEv.Position = objPos
		touchEv.AbsolutePosition = pos
		wid.TouchUp(touchEv)
		c.touched[tapID] = nil
	}

	ev := &fyne.PointEvent{
		Position:         objPos,
		AbsolutePosition: pos,
	}

	if duration < tapSecondaryDelay {
		_, doubleTap := co.(fyne.DoubleTappable)
		if doubleTap {
			c.touchTapCount++
			c.touchLastTapped = co
			if c.touchCancelFunc != nil {
				c.touchCancelFunc()
				return
			}
			go c.waitForDoubleTap(co, ev, tapCallback, doubleTapCallback)
		} else {
			if wid, ok := co.(fyne.Tappable); ok {
				tapCallback(wid, ev)
			}
		}
	} else {
		if wid, ok := co.(fyne.SecondaryTappable); ok {
			tapAltCallback(wid, ev)
		}
	}
}

const tapDoubleDelay = 500 * time.Millisecond

func (c *SoftwareCanvas) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent, tapCallback func(fyne.Tappable, *fyne.PointEvent), doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent)) {
	var ctx context.Context
	ctx, c.touchCancelFunc = context.WithDeadline(context.TODO(), time.Now().Add(tapDoubleDelay))
	defer c.touchCancelFunc()
	<-ctx.Done()
	if c.touchTapCount == 2 && c.touchLastTapped == co {
		if wid, ok := co.(fyne.DoubleTappable); ok {
			doubleTapCallback(wid, ev)
		}
	} else {
		if wid, ok := co.(fyne.Tappable); ok {
			tapCallback(wid, ev)
		}
	}
	c.touchTapCount = 0
	c.touchCancelFunc = nil
	c.touchLastTapped = nil
}
