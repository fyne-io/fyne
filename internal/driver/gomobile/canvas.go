package gomobile

import (
	"context"
	"image"
	"math"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type mobileCanvas struct {
	content          fyne.CanvasObject
	windowHead, menu fyne.CanvasObject
	overlays         *internal.OverlayStack
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
	minSizeCache   map[fyne.CanvasObject]fyne.Size

	touchTapCount   int
	touchCancelFunc context.CancelFunc
	touchLastTapped fyne.CanvasObject
}

const (
	doubleClickDelay = 500 // ms (maximum interval between clicks for double click detection)
)

func (c *mobileCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *mobileCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	scale := fyne.CurrentDevice().SystemScaleForWindow(nil) // we don't need a window parameter on mobile

	dev, ok := fyne.CurrentDevice().(*device)
	if !ok || dev.safeWidth == 0 || dev.safeHeight == 0 {
		return fyne.NewPos(0, 0), c.Size() // running in test mode
	}

	return fyne.NewPos(int(float32(dev.safeLeft)/scale), int(float32(dev.safeTop)/scale)),
		fyne.NewSize(int(float32(dev.safeWidth)/scale), int(float32(dev.safeHeight)/scale))
}

func (c *mobileCanvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	c.sizeContent(c.Size().Max(content.MinSize()))
}

// Deprecated: Use Overlays() instead.
func (c *mobileCanvas) Overlay() fyne.CanvasObject {
	return c.overlays.Top()
}

func (c *mobileCanvas) Overlays() fyne.OverlayStack {
	return c.overlays
}

// Deprecated: Use Overlays() instead.
func (c *mobileCanvas) SetOverlay(overlay fyne.CanvasObject) {
	if len(c.overlays.List()) > 0 {
		c.overlays.Remove(c.overlays.List()[0])
	}
	c.overlays.Add(overlay)
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
	if c.content == nil { // window may not be configured yet
		return
	}
	c.size = size

	offset := fyne.NewPos(0, 0)
	areaPos, areaSize := c.InteractiveArea()

	if c.windowHead != nil {
		topHeight := c.windowHead.MinSize().Height

		if len(c.windowHead.(*widget.Box).Children) > 1 {
			c.windowHead.Resize(fyne.NewSize(areaSize.Width, topHeight))
			offset = fyne.NewPos(0, topHeight)
		} else {
			c.windowHead.Resize(c.windowHead.MinSize())
		}
		c.windowHead.Move(areaPos)
	}

	topLeft := areaPos.Add(offset)
	for _, overlay := range c.overlays.List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			size := p.Content.Size().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)).Min(areaSize)
			p.Resize(size)
		} else {
			overlay.Resize(areaSize)
			overlay.Move(topLeft)
		}
	}

	if c.padded {
		c.content.Resize(areaSize.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		c.content.Move(topLeft.Add(fyne.NewPos(theme.Padding(), theme.Padding())))
	} else {
		c.content.Resize(areaSize)
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
	if c.focused == obj || obj == nil {
		return
	}

	if dis, ok := obj.(fyne.Disableable); ok && dis.Disabled() {
		c.Unfocus()
		return
	}

	if c.focused != nil {
		c.focused.FocusLost()
	}

	c.focused = obj
	obj.FocusGained()

	if keyb, ok := obj.(mobile.Keyboardable); ok {
		showVirtualKeyboard(keyb.Keyboard())
	} else {
		hideVirtualKeyboard()
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
	c.scale = fyne.CurrentDevice().SystemScaleForWindow(nil) // we don't need a window parameter on mobile
}

func (c *mobileCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	return int(float32(pos.X) * c.scale), int(float32(pos.Y) * c.scale)
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

func (c *mobileCanvas) RemoveShortcut(shortcut fyne.Shortcut) {
	c.shortcut.RemoveShortcut(shortcut)
}

func (c *mobileCanvas) Capture() image.Image {
	return c.painter.Capture(c)
}

func (c *mobileCanvas) minSizeChanged() bool {
	if c.Content() == nil {
		return false
	}
	minSizeChange := false

	ensureMinSize := func(obj, parent fyne.CanvasObject) {
		if !obj.Visible() {
			return
		}
		minSize := obj.MinSize()

		if minSize != c.minSizeCache[obj] {
			minSizeChange = true

			c.minSizeCache[obj] = minSize
		}
	}
	c.walkTree(nil, ensureMinSize)

	return minSizeChange
}

func (c *mobileCanvas) ensureMinSize() {
	if c.Content() == nil {
		return
	}
	var lastParent fyne.CanvasObject

	ensureMinSize := func(obj, parent fyne.CanvasObject) {
		if !obj.Visible() {
			return
		}
		minSize := obj.MinSize()

		objToLayout := obj
		if parent != nil {
			objToLayout = parent
		} else {
			size := obj.Size()
			expectedSize := minSize.Union(size)
			if expectedSize != size && size != c.size {
				objToLayout = nil
				obj.Resize(expectedSize)
			}
		}

		if objToLayout != lastParent {
			updateLayout(lastParent)
			lastParent = objToLayout
		}
	}
	c.walkTree(nil, ensureMinSize)

	if lastParent != nil {
		updateLayout(lastParent)
	}
}

func updateLayout(objToLayout fyne.CanvasObject) {
	switch cont := objToLayout.(type) {
	case *fyne.Container:
		if cont.Layout != nil {
			cont.Layout.Layout(cont.Objects, cont.Size())
		}
	case fyne.Widget:
		cache.Renderer(cont).Layout(cont.Size())
	}
}

func (c *mobileCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+2)
	trees = append(trees, c.content)
	if c.menu != nil {
		trees = append(trees, c.menu)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
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
	for _, overlay := range c.overlays.List() {
		if overlay != nil {
			driver.WalkVisibleObjectTree(overlay, beforeChildren, afterChildren)
		}
	}
}

func (c *mobileCanvas) findObjectAtPositionMatching(pos fyne.Position, test func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	if c.menu != nil {
		return driver.FindObjectAtPositionMatching(pos, test, c.overlays.Top(), c.menu)
	}

	return driver.FindObjectAtPositionMatching(pos, test, c.overlays.Top(), c.windowHead, c.content)
}

func (c *mobileCanvas) tapDown(pos fyne.Position, tapID int) {
	c.lastTapDown[tapID] = time.Now()
	c.lastTapDownPos[tapID] = pos
	c.dragging = nil

	co, objPos, layer := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case fyne.Tappable, mobile.Touchable, fyne.Focusable:
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
		if wid, ok := co.(fyne.Focusable); ok {
			c.Focus(wid)
		} else {
			c.Unfocus()
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
		} else {
			return
		}
	}

	ev := new(fyne.DragEvent)
	ev.Position = objPos
	ev.DraggedX = deltaX
	ev.DraggedY = deltaY

	dragCallback(c.dragging, ev)
}

func (c *mobileCanvas) tapUp(pos fyne.Position, tapID int,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.SecondaryTappable, *fyne.PointEvent),
	doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable, *fyne.DragEvent)) {
	if c.dragging != nil {
		c.dragging.DragEnd()

		c.dragging = nil
		return
	}

	duration := time.Since(c.lastTapDown[tapID])

	if c.menu != nil && c.overlays.Top() == nil && pos.X > c.menu.Size().Width {
		c.menu.Hide()
		c.menu.Refresh()
		c.menu = nil
		return
	}

	co, objPos, _ := c.findObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.SecondaryTappable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
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

	ev := new(fyne.PointEvent)
	ev.Position = objPos
	ev.AbsolutePosition = pos

	// TODO move event queue to common code w.queueEvent(func() { wid.Tapped(ev) })
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

func (c *mobileCanvas) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent, tapCallback func(fyne.Tappable, *fyne.PointEvent), doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent)) {
	var ctx context.Context
	ctx, c.touchCancelFunc = context.WithDeadline(context.TODO(), time.Now().Add(time.Millisecond*doubleClickDelay))
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
	ret.scale = fyne.CurrentDevice().SystemScaleForWindow(nil) // we don't need a window parameter on mobile
	ret.refreshQueue = make(chan fyne.CanvasObject, 1024)
	ret.touched = make(map[int]mobile.Touchable)
	ret.lastTapDownPos = make(map[int]fyne.Position)
	ret.lastTapDown = make(map[int]time.Time)
	ret.minSizeCache = make(map[fyne.CanvasObject]fyne.Size)
	ret.overlays = &internal.OverlayStack{Canvas: ret}

	ret.setupThemeListener()

	return ret
}
