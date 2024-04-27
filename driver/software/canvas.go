package software

import (
	"image"
	"image/draw"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter/software"
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

type softwareCanvas struct {
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

	refreshQueue *async.CanvasObjectQueue
	dirty        atomic.Bool
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
	c := &softwareCanvas{
		focusMgr:     app.NewFocusManager(nil),
		padded:       true,
		scale:        1.0,
		size:         fyne.NewSize(10, 10),
		refreshQueue: async.NewCanvasObjectQueue(),
	}
	c.overlays = &internal.OverlayStack{Canvas: c}
	return c
}

// NewCanvasWithPainter allows creation of an in-memory canvas with a specific painter.
// The painter will be used to render in the Capture() call.
func NewCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	canvas := NewCanvas().(*softwareCanvas)
	canvas.painter = painter

	return canvas
}

// NewTransparentCanvasWithPainter allows creation of an in-memory canvas with a specific painter without a background color.
// The painter will be used to render in the Capture() call.
//
// Since: 2.2
func NewTransparentCanvasWithPainter(painter SoftwarePainter) WindowlessCanvas {
	canvas := NewCanvasWithPainter(painter).(*softwareCanvas)
	canvas.transparent = true

	return canvas
}

func (c *softwareCanvas) Capture() image.Image {
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

func (c *softwareCanvas) Content() fyne.CanvasObject {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.content
}

func (c *softwareCanvas) findObjectAtPositionMatching(pos fyne.Position, test func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	return driver.FindObjectAtPositionMatching(pos, test, c.Overlays().Top(), c.content)
}

func (c *softwareCanvas) Focus(obj fyne.Focusable) {
	c.focusManager().Focus(obj)
}

func (c *softwareCanvas) FocusNext() {
	c.focusManager().FocusNext()
}

func (c *softwareCanvas) FocusPrevious() {
	c.focusManager().FocusPrevious()
}

func (c *softwareCanvas) Focused() fyne.Focusable {
	return c.focusManager().Focused()
}

func (c *softwareCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

func (c *softwareCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedKey
}

func (c *softwareCanvas) OnTypedRune() func(rune) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.onTypedRune
}

func (c *softwareCanvas) Overlays() fyne.OverlayStack {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.overlays
}

func (c *softwareCanvas) Padded() bool {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.padded
}

func (c *softwareCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *softwareCanvas) Refresh(obj fyne.CanvasObject) {
	walkNeeded := false
	switch obj.(type) {
	case *fyne.Container:
		walkNeeded = true
	case fyne.Widget:
		walkNeeded = true
	}

	if walkNeeded {
		driver.WalkCompleteObjectTree(obj, func(co fyne.CanvasObject, p1, p2 fyne.Position, s fyne.Size) bool {
			if i, ok := co.(*canvas.Image); ok {
				i.Refresh()
			}
			return false
		}, nil)
	}

	c.refreshQueue.In(obj)
	c.SetDirty()
}

func (c *softwareCanvas) Resize(size fyne.Size) {
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

func (c *softwareCanvas) Scale() float32 {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.scale
}

func (c *softwareCanvas) SetContent(content fyne.CanvasObject) {
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

// CheckDirtyAndClear returns true if the canvas is dirty and
// clears the dirty state atomically.
func (c *softwareCanvas) CheckDirtyAndClear() bool {
	return c.dirty.Swap(false)
}

// SetDirty sets canvas dirty flag atomically.
func (c *softwareCanvas) SetDirty() {
	c.dirty.Store(true)
}

func (c *softwareCanvas) SetOnTypedKey(handler func(*fyne.KeyEvent)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedKey = handler
}

func (c *softwareCanvas) SetOnTypedRune(handler func(rune)) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.onTypedRune = handler
}

func (c *softwareCanvas) SetPadded(padded bool) {
	c.propertyLock.Lock()
	c.padded = padded
	c.propertyLock.Unlock()

	c.Resize(c.Size())
}

func (c *softwareCanvas) SetScale(scale float32) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.scale = scale
}

func (c *softwareCanvas) Size() fyne.Size {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	return c.size
}

func (c *softwareCanvas) Unfocus() {
	c.focusManager().Focus(nil)
}

func (c *softwareCanvas) focusManager() *app.FocusManager {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	return c.focusMgr
}

func (c *softwareCanvas) objectTrees() []fyne.CanvasObject {
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

func (c *softwareCanvas) tapDown(pos fyne.Position, tapID int) {
	// c.lastTapDown[tapID] = time.Now()
	// c.lastTapDownPos[tapID] = pos
	// c.dragging = nil

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
		// c.touched[tapID] = wid
	}

	if layer != 1 { // 0 - overlay, 1 - window head / menu, 2 - content
		if wid, ok := co.(fyne.Focusable); !ok || wid != c.Focused() {
			c.Unfocus()
		}
	}
}

// func (c *softwareCanvas) tapMove(pos fyne.Position, tapID int,
//
//		dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
//		previousPos := c.lastTapDownPos[tapID]
//		deltaX := pos.X - previousPos.X
//		deltaY := pos.Y - previousPos.Y
//
//		if c.dragging == nil && (math.Abs(float64(deltaX)) < tapMoveThreshold && math.Abs(float64(deltaY)) < tapMoveThreshold) {
//			return
//		}
//		c.lastTapDownPos[tapID] = pos
//
//		co, objPos, _ := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
//			if _, ok := object.(fyne.Draggable); ok {
//				return true
//			} else if _, ok := object.(mobile.Touchable); ok {
//				return true
//			}
//
//			return false
//		})
//
//		if c.touched[tapID] != nil {
//			if touch, ok := co.(mobile.Touchable); !ok || c.touched[tapID] != touch {
//				touchEv := &mobile.TouchEvent{}
//				touchEv.Position = objPos
//				touchEv.AbsolutePosition = pos
//				c.touched[tapID].TouchCancel(touchEv)
//				c.touched[tapID] = nil
//			}
//		}
//
//		if c.dragging == nil {
//			if drag, ok := co.(fyne.Draggable); ok {
//				c.dragging = drag
//				c.dragOffset = previousPos.Subtract(objPos)
//				c.dragStart = co.Position()
//			} else {
//				return
//			}
//		}
//
//		ev := &fyne.DragEvent{}
//		draggedObjDelta := c.dragStart.Subtract(c.dragging.(fyne.CanvasObject).Position())
//		ev.Position = pos.Subtract(c.dragOffset).Add(draggedObjDelta)
//		ev.Dragged = fyne.Delta{DX: deltaX, DY: deltaY}
//
//		dragCallback(c.dragging, ev)
//	}
func (c *softwareCanvas) tapUp(pos fyne.Position, tapID int,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.SecondaryTappable, *fyne.PointEvent),
	doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable)) {

	// if c.dragging != nil {
	// 	dragCallback(c.dragging)
	//
	// 	c.dragging = nil
	// 	return
	// }

	// duration := time.Since(c.lastTapDown[tapID])

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
		// c.touched[tapID] = nil
	}

	ev := &fyne.PointEvent{
		Position:         objPos,
		AbsolutePosition: pos,
	}

	// if duration < tapSecondaryDelay {
	// 	_, doubleTap := co.(fyne.DoubleTappable)
	// 	if doubleTap {
	// 		c.touchTapCount++
	// 		c.touchLastTapped = co
	// 		if c.touchCancelFunc != nil {
	// 			c.touchCancelFunc()
	// 			return
	// 		}
	// 		go c.waitForDoubleTap(co, ev, tapCallback, doubleTapCallback)
	// 	} else {
	if wid, ok := co.(fyne.Tappable); ok {
		tapCallback(wid, ev)
	}
	// 	}
	// } else {
	// 	if wid, ok := co.(fyne.SecondaryTappable); ok {
	// 		tapAltCallback(wid, ev)
	// 	}
	// }
}

// NewTransparentCanvas creates a new canvas in memory that can render without hardware support without a background color.
//
// Since: 2.2
func NewTransparentCanvas() WindowlessCanvas {
	return NewTransparentCanvasWithPainter(software.NewPainter())
}
