package gomobile

import (
	"image"
	"math"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
)

type canvas struct {
	content, overlay fyne.CanvasObject
	windowHead       fyne.CanvasObject
	painter          gl.Painter
	scale            float32
	size             fyne.Size

	focused fyne.Focusable
	padded  bool

	onTypedRune func(rune)
	onTypedKey  func(event *fyne.KeyEvent)
	shortcut    fyne.ShortcutHandler

	inited, dirty  bool
	lastTapDown    int64
	lastTapDownPos fyne.Position
	dragging       fyne.Draggable
	refreshQueue   chan fyne.CanvasObject
}

func (c *canvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	c.sizeContent(c.Size().Union(content.MinSize()))
}

func (c *canvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
	c.dirty = true
}

func (c *canvas) sizeContent(size fyne.Size) {
	offset := fyne.NewPos(0, 0)
	if c.windowHead != nil {
		topHeight := c.windowHead.MinSize().Height
		offset = fyne.NewPos(0, topHeight)
		c.windowHead.Resize(fyne.NewSize(size.Width, topHeight))
	}

	devicePadTopLeft, devicePadBottomRight := devicePadding()
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

func (c *canvas) resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.sizeContent(size)
}

func (c *canvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		c.focused.FocusLost()
	}

	c.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (c *canvas) Unfocus() {
	if c.focused != nil {
		c.focused.FocusLost()
	}
	c.focused = nil
}

func (c *canvas) Focused() fyne.Focusable {
	return c.focused
}

func (c *canvas) Size() fyne.Size {
	return c.size
}

func (c *canvas) Scale() float32 {
	return c.scale
}

func (c *canvas) SetScale(scale float32) {
	if scale == fyne.SettingsScaleAuto {
		c.scale = deviceScale()
	} else if scale == 0 { // not set in the config
		return
	} else {
		c.scale = scale
	}
}

func (c *canvas) Overlay() fyne.CanvasObject {
	return c.overlay
}

func (c *canvas) SetOverlay(overlay fyne.CanvasObject) {
	c.overlay = overlay

	if c.overlay != nil {
		c.overlay.Resize(c.size)
	}
}

func (c *canvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *canvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *canvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *canvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *canvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

func (c *canvas) Capture() image.Image {
	return c.painter.Capture(c)
}

func (c *canvas) walkTree(
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) {
	driver.WalkVisibleObjectTree(c.content, beforeChildren, afterChildren)
	//if c.menu != nil {
	//	driver.WalkVisibleObjectTree(c.menu, beforeChildren, afterChildren)
	//}
	if c.overlay != nil {
		driver.WalkVisibleObjectTree(c.overlay, beforeChildren, afterChildren)
	}
	if c.windowHead != nil {
		driver.WalkVisibleObjectTree(c.windowHead, beforeChildren, afterChildren)
	}
}

func (c *canvas) tapDown(pos fyne.Position) {
	c.lastTapDown = time.Now().UnixNano()
	c.lastTapDownPos = pos
	c.dragging = nil

	co, _ := driver.FindObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Focusable); ok {
			return true
		}

		return false
	}, c.overlay, c.content, c.windowHead)

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
		c.Focus(wid)
	}
}

func (c *canvas) tapMove(pos fyne.Position,
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	deltaX := pos.X - c.lastTapDownPos.X
	deltaY := pos.Y - c.lastTapDownPos.Y

	if math.Abs(float64(deltaX)) < 3 && math.Abs(float64(deltaY)) < 3 {
		return
	}

	if c.dragging == nil {
		co, _ := driver.FindObjectAtPositionMatching(c.lastTapDownPos, func(object fyne.CanvasObject) bool {
			if _, ok := object.(fyne.Draggable); ok {
				return true
			}

			return false
		}, c.overlay, c.content, c.windowHead)

		if drag, ok := co.(fyne.Draggable); ok {
			c.dragging = drag
		} else {
			return
		}
	}
	objPos := pos.Subtract(c.dragging.(fyne.CanvasObject).Position())

	ev := new(fyne.DragEvent)
	ev.Position = objPos
	ev.DraggedX = deltaX
	ev.DraggedY = deltaY

	dragCallback(c.dragging, ev)
	c.lastTapDownPos = pos
}

func (c *canvas) tapUp(pos fyne.Position,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.Tappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	if c.dragging != nil {
		c.dragging.DragEnd()

		c.dragging = nil
	}

	duration := time.Now().UnixNano() - c.lastTapDown

	co, objPos := driver.FindObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
			return true
		}

		return false
	}, c.overlay, c.content, c.windowHead)

	ev := new(fyne.PointEvent)
	ev.Position = objPos

	if wid, ok := co.(fyne.Tappable); ok {
		// TODO move event queue to common code w.queueEvent(func() { wid.Tapped(ev) })
		if duration < tapSecondaryDelay {
			tapCallback(wid, ev)
		} else {
			tapAltCallback(wid, ev)
		}
	}
}

// NewCanvas creates a new gomobile canvas. This is a canvas that will render on a mobile device using OpenGL.
func NewCanvas() fyne.Canvas {
	ret := &canvas{padded: true}
	ret.scale = deviceScale()
	ret.refreshQueue = make(chan fyne.CanvasObject, 1024)

	return ret
}
