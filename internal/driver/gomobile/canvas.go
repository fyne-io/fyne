package gomobile

import (
	"image"
	"math"
	"time"

	"fyne.io/fyne"
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
	padded  bool

	onTypedRune func(rune)
	onTypedKey  func(event *fyne.KeyEvent)
	shortcut    fyne.ShortcutHandler

	inited         bool
	lastTapDown    time.Time
	lastTapDownPos fyne.Position
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
	}
}

func (c *mobileCanvas) Unfocus() {
	if c.focused != nil {
		c.focused.FocusLost()
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

func (c *mobileCanvas) SetScale(scale float32) {
	if scale == fyne.SettingsScaleAuto {
		c.scale = deviceScale()
	} else if scale == 0 { // not set in the config
		return
	} else {
		c.scale = scale
	}
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

func (c *mobileCanvas) findObjectMatching(test func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position) {
	if c.menu != nil && c.overlay == nil {
		return driver.FindObjectAtPositionMatching(c.lastTapDownPos, test, c.menu)
	}

	return driver.FindObjectAtPositionMatching(c.lastTapDownPos, test, c.overlay, c.windowHead, c.content)
}

func (c *mobileCanvas) tapDown(pos fyne.Position) {
	c.lastTapDown = time.Now()
	c.lastTapDownPos = pos
	c.dragging = nil

	co, _ := c.findObjectMatching(func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
			return true
		}

		return false
	})

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

func (c *mobileCanvas) tapMove(pos fyne.Position,
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	deltaX := pos.X - c.lastTapDownPos.X
	deltaY := pos.Y - c.lastTapDownPos.Y

	if math.Abs(float64(deltaX)) < 3 && math.Abs(float64(deltaY)) < 3 {
		return
	}

	if c.dragging == nil {
		co, _ := c.findObjectMatching(func(object fyne.CanvasObject) bool {
			if _, ok := object.(fyne.Draggable); ok {
				return true
			}

			return false
		})

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

func (c *mobileCanvas) tapUp(pos fyne.Position,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.Tappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	if c.dragging != nil {
		c.dragging.DragEnd()

		c.dragging = nil
	}

	duration := time.Since(c.lastTapDown)

	if c.menu != nil && c.overlay == nil && pos.X > c.menu.Size().Width {
		c.menu.Hide()
		c.menu = nil
		return
	}

	co, objPos := c.findObjectMatching(func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
			return true
		}

		return false
	})

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
	ret.scale = deviceScale()
	ret.refreshQueue = make(chan fyne.CanvasObject, 1024)

	ret.setupThemeListener()

	return ret
}
