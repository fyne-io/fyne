package mobile

import (
	"context"
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/app"
	intdriver "fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Canvas = (*canvas)(nil)

type canvas struct {
	common.Canvas
	content        fyne.CanvasObject
	device         *device
	initialized    bool
	lastTapDown    map[int]time.Time
	lastTapDownPos map[int]fyne.Position
	lastTapDelta   map[int]fyne.Delta
	menu           fyne.CanvasObject
	padded         bool
	scale          float32
	size           fyne.Size
	touched        map[int]mobile.Touchable
	windowHead     fyne.CanvasObject

	dragOffset fyne.Position
	dragStart  fyne.Position
	dragging   fyne.Draggable

	onTypedKey  func(event *fyne.KeyEvent)
	onTypedRune func(rune)

	touchCancelFunc context.CancelFunc
	touchLastTapped fyne.CanvasObject
	touchTapCount   int
}

func newCanvas(dev fyne.Device) fyne.Canvas {
	d, _ := dev.(*device)
	ret := &canvas{
		Canvas: common.Canvas{
			OnFocus:   handleKeyboard,
			OnUnfocus: hideVirtualKeyboard,
		},
		device:         d,
		lastTapDown:    make(map[int]time.Time),
		lastTapDownPos: make(map[int]fyne.Position),
		lastTapDelta:   make(map[int]fyne.Delta),
		padded:         true,
		scale:          dev.SystemScaleForWindow(nil), // we don't need a window parameter on mobile,
		touched:        make(map[int]mobile.Touchable),
	}
	ret.Initialize(ret, ret.overlayChanged)
	return ret
}

func (c *canvas) Capture() image.Image {
	return c.Painter().Capture(c)
}

func (c *canvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *canvas) InteractiveArea() (fyne.Position, fyne.Size) {
	var pos fyne.Position
	var size fyne.Size
	if c.device == nil {
		// running in test mode
		size = c.Size()
	} else {
		safeLeft := float32(c.device.safeLeft) / c.scale
		safeTop := float32(c.device.safeTop) / c.scale
		safeRight := float32(c.device.safeRight) / c.scale
		safeBottom := float32(c.device.safeBottom) / c.scale
		pos = fyne.NewPos(safeLeft, safeTop)
		size = c.size.SubtractWidthHeight(safeLeft+safeRight, safeTop+safeBottom)
	}
	if c.windowHeadIsDisplacing() {
		offset := c.windowHead.MinSize().Height
		pos = pos.AddXY(0, offset)
		size = size.SubtractWidthHeight(0, offset)
	}
	return pos, size
}

func (c *canvas) MinSize() fyne.Size {
	return c.size // TODO check
}

func (c *canvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *canvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *canvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *canvas) Resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.sizeContent(size)
}

func (c *canvas) Scale() float32 {
	return c.scale
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.setContent(content)
	c.sizeContent(c.Size()) // fixed window size for mobile, cannot stretch to new content
	c.SetDirty()
}

func (c *canvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *canvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *canvas) Size() fyne.Size {
	return c.size
}

func (c *canvas) applyThemeOutOfTreeObjects() {
	if c.menu != nil {
		app.ApplyThemeTo(c.menu, c) // Ensure our menu gets the theme change message as it's out-of-tree
	}
	if c.windowHead != nil {
		app.ApplyThemeTo(c.windowHead, c) // Ensure our child windows get the theme change message as it's out-of-tree
	}
}

func (c *canvas) findObjectAtPositionMatching(pos fyne.Position, test func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	if c.menu != nil {
		return intdriver.FindObjectAtPositionMatching(pos, test, c.Overlays().Top(), c.menu)
	}

	return intdriver.FindObjectAtPositionMatching(pos, test, c.Overlays().Top(), c.windowHead, c.content)
}

func (c *canvas) overlayChanged() {
	handleKeyboard(c.Focused())
	c.SetDirty()
}

func (c *canvas) setContent(content fyne.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

func (c *canvas) setMenu(menu fyne.CanvasObject) {
	c.menu = menu
	c.SetMenuTreeAndFocusMgr(menu)
}

func (c *canvas) setWindowHead(head fyne.CanvasObject) {
	if c.padded {
		head = container.NewPadded(head)
	}
	c.windowHead = head
	c.SetMobileWindowHeadTree(head)
}

func (c *canvas) sizeContent(size fyne.Size) {
	if c.content == nil { // window may not be configured yet
		return
	}

	c.size = size
	areaPos, areaSize := c.InteractiveArea()

	if c.windowHead != nil {
		var headSize fyne.Size
		headPos := areaPos
		if c.windowHeadIsDisplacing() {
			headSize = fyne.NewSize(areaSize.Width, c.windowHead.MinSize().Height)
			headPos = headPos.SubtractXY(0, headSize.Height)
		} else {
			headSize = c.windowHead.MinSize()
		}
		c.windowHead.Resize(headSize)
		c.windowHead.Move(headPos)
	}

	for _, overlay := range c.Overlays().List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(areaSize)
			overlay.Move(areaPos)
		}
	}

	if c.padded {
		c.content.Resize(areaSize.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		c.content.Move(areaPos.Add(fyne.NewPos(theme.Padding(), theme.Padding())))
	} else {
		c.content.Resize(areaSize)
		c.content.Move(areaPos)
	}
}

func (c *canvas) tapDown(pos fyne.Position, tapID int) {
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

func (c *canvas) tapMove(pos fyne.Position, tapID int,
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	previousPos := c.lastTapDownPos[tapID]
	deltaX := pos.X - previousPos.X
	deltaY := pos.Y - previousPos.Y

	if c.dragging == nil && (math.Abs(float64(deltaX)) < tapMoveThreshold && math.Abs(float64(deltaY)) < tapMoveThreshold) {
		return
	}
	c.lastTapDownPos[tapID] = pos
	offset := fyne.Delta{DX: deltaX, DY: deltaY}
	c.lastTapDelta[tapID] = offset

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
	ev.Dragged = offset

	dragCallback(c.dragging, ev)
}

func (c *canvas) tapUp(pos fyne.Position, tapID int,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.SecondaryTappable, *fyne.PointEvent),
	doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {

	if c.dragging != nil {
		previousDelta := c.lastTapDelta[tapID]
		dragCallback(c.dragging, &fyne.DragEvent{Dragged: previousDelta})

		c.dragging = nil
		return
	}

	duration := time.Since(c.lastTapDown[tapID])

	if c.menu != nil && c.Overlays().Top() == nil && pos.X > c.menu.Size().Width {
		c.menu.Hide()
		c.menu.Refresh()
		c.setMenu(nil)
		return
	}

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

func (c *canvas) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent, tapCallback func(fyne.Tappable, *fyne.PointEvent), doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent)) {
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

func (c *canvas) windowHeadIsDisplacing() bool {
	if c.windowHead == nil {
		return false
	}

	chromeBox := c.windowHead.(*fyne.Container)
	if c.padded {
		chromeBox = chromeBox.Objects[0].(*fyne.Container) // the padded container
	}
	return len(chromeBox.Objects) > 1
}
