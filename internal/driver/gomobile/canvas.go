package gomobile

import (
	"image"
	"math"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type mobileCanvas struct {
	content, overlay fyne.CanvasObject
	windowHead, menu fyne.CanvasObject
	painter          gl.Painter
	scale            float32
	size             fyne.Size

	focused fyne.Focusable
	touched map[int]mobile.Touchable
	padded  bool

	onTypedRune func(rune)
	onTypedKey  func(event *fyne.KeyEvent)
	shortcut    fyne.ShortcutHandler

	inited         bool
	lastTapDown    map[int]time.Time
	lastTapDownPos map[int]fyne.Position
	dragging       fyne.Draggable
	refreshQueue   chan fyne.CanvasObject
}

func (c *mobileCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *mobileCanvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	c.sizeContent(c.Size().Union(content.MinSize()))
}

func (c *mobileCanvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
}

func (c *mobileCanvas) sizeContent(size fyne.Size) {
	offset := fyne.NewPos(0, 0)
	devicePadTopLeft, devicePadBottomRight := devicePadding()

	if c.windowHead != nil {
		topHeight := c.windowHead.MinSize().Height

		if len(c.windowHead.(*widget.Box).Children) > 1 {
			c.windowHead.Resize(fyne.NewSize(size.Width-devicePadTopLeft.Width-devicePadBottomRight.Width, topHeight))
			offset = fyne.NewPos(0, topHeight)
		} else {
			c.windowHead.Resize(c.windowHead.MinSize())
		}
		c.windowHead.Move(fyne.NewPos(devicePadTopLeft.Width, devicePadTopLeft.Height))
	}

	innerSize := size.Subtract(devicePadTopLeft).Subtract(devicePadBottomRight)
	topLeft := offset.Add(fyne.NewPos(devicePadTopLeft.Width, devicePadTopLeft.Height))

	c.size = size
	if c.padded {
		c.content.Resize(innerSize.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		c.content.Move(topLeft.Add(fyne.NewPos(theme.Padding(), theme.Padding())))
	} else {
		c.content.Resize(innerSize)
		c.content.Move(topLeft)
	}
}

func (c *mobileCanvas) resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.sizeContent(size)
}

func (c *mobileCanvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		c.focused.FocusLost()
	}

	c.focused = obj
	if obj != nil {
		obj.FocusGained()

		if _, ok := obj.(*widget.Entry); ok {
			showVirtualKeyboard()
		}
	}
}

func (c *mobileCanvas) Unfocus() {
	if c.focused != nil {
		c.focused.FocusLost()
		hideVirtualKeyboard()
	}
	c.focused = nil
}

func (c *mobileCanvas) Focused() fyne.Focusable {
	return c.focused
}

func (c *mobileCanvas) Size() fyne.Size {
	return c.size
}

func (c *mobileCanvas) Scale() float32 {
	return c.scale
}

// Deprecated: Settings are now calculated solely on the user configuration and system setup.
func (c *mobileCanvas) SetScale(_ float32) {
	c.scale = fyne.CurrentDevice().SystemScale()
}

func (c *mobileCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
}

func (c *mobileCanvas) Overlay() fyne.CanvasObject {
	return c.overlay
}

func (c *mobileCanvas) SetOverlay(overlay fyne.CanvasObject) {
	c.overlay = overlay
}

func (c *mobileCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *mobileCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *mobileCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *mobileCanvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *mobileCanvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

func (c *mobileCanvas) Capture() image.Image {
	return c.painter.Capture(c)
}

func (c *mobileCanvas) walkTree(
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) {
	driver.WalkVisibleObjectTree(c.content, beforeChildren, afterChildren)
	if c.windowHead != nil {
		driver.WalkVisibleObjectTree(c.windowHead, beforeChildren, afterChildren)
	}
	if c.menu != nil {
		driver.WalkVisibleObjectTree(c.menu, beforeChildren, afterChildren)
	}
	if c.overlay != nil {
		driver.WalkVisibleObjectTree(c.overlay, beforeChildren, afterChildren)
	}
}

func (c *mobileCanvas) findObjectAtPositionMatching(pos fyne.Position, test func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position) {
	if c.menu != nil && c.overlay == nil {
		return driver.FindObjectAtPositionMatching(pos, test, c.menu)
	}

	return driver.FindObjectAtPositionMatching(pos, test, c.overlay, c.windowHead, c.content)
}

func (c *mobileCanvas) tapDown(pos fyne.Position, tapID int) {
	c.lastTapDown[tapID] = time.Now()
	c.lastTapDownPos[tapID] = pos
	c.dragging = nil

	co, objPos := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
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

	needsFocus := true
	wid := c.Focused()
	if wid != nil {
		if wid.(fyne.CanvasObject) != co {
			c.Unfocus()
		} else {
			needsFocus = false
		}
	}
	if wid, ok := co.(fyne.Focusable); ok && needsFocus {
		if dis, ok := wid.(fyne.Disableable); !ok || !dis.Disabled() {
			c.Focus(wid)
		}
	}
}

func (c *mobileCanvas) tapMove(pos fyne.Position, tapID int,
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	deltaX := pos.X - c.lastTapDownPos[tapID].X
	deltaY := pos.Y - c.lastTapDownPos[tapID].Y

	if math.Abs(float64(deltaX)) < 3 && math.Abs(float64(deltaY)) < 3 {
		return
	}
	c.lastTapDownPos[tapID] = pos

	co, objPos := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
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
		} else {
			return
		}
	}
	objPos = pos.Subtract(c.dragging.(fyne.CanvasObject).Position())

	ev := new(fyne.DragEvent)
	ev.Position = objPos
	ev.DraggedX = deltaX
	ev.DraggedY = deltaY

	dragCallback(c.dragging, ev)
}

func (c *mobileCanvas) tapUp(pos fyne.Position, tapID int,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.Tappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	if c.dragging != nil {
		c.dragging.DragEnd()

		c.dragging = nil
	}

	duration := time.Since(c.lastTapDown[tapID])

	if c.menu != nil && c.overlay == nil && pos.X > c.menu.Size().Width {
		c.menu.Hide()
		c.menu = nil
		return
	}

	co, objPos := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
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

	ev := new(fyne.PointEvent)
	ev.Position = objPos
	ev.AbsolutePosition = pos

	if wid, ok := co.(fyne.Tappable); ok {
		// TODO move event queue to common code w.queueEvent(func() { wid.Tapped(ev) })
		if duration < tapSecondaryDelay {
			tapCallback(wid, ev)
		} else {
			tapAltCallback(wid, ev)
		}
	}
}

func (c *mobileCanvas) setupThemeListener() {
	listener := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			if c.menu != nil {
				app.ApplyThemeTo(c.menu, c) // Ensure our menu gets the theme change message as it's out-of-tree
			}
			if c.windowHead != nil {
				app.ApplyThemeTo(c.windowHead, c) // Ensure our child windows get the theme change message as it's out-of-tree
			}
		}
	}()
}

// NewCanvas creates a new gomobile mobileCanvas. This is a mobileCanvas that will render on a mobile device using OpenGL.
func NewCanvas() fyne.Canvas {
	ret := &mobileCanvas{padded: true}
	ret.scale = fyne.CurrentDevice().SystemScale()
	ret.refreshQueue = make(chan fyne.CanvasObject, 1024)
	ret.touched = make(map[int]mobile.Touchable)
	ret.lastTapDownPos = make(map[int]fyne.Position)
	ret.lastTapDown = make(map[int]time.Time)

	ret.setupThemeListener()

	return ret
}
