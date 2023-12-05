package glfw

import (
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Canvas interface
var _ fyne.Canvas = (*glCanvas)(nil)

type glCanvas struct {
	common.Canvas

	content       fyne.CanvasObject
	menu          fyne.CanvasObject
	padded, debug bool
	size          fyne.Size

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
	onKeyDown   func(*fyne.KeyEvent)
	onKeyUp     func(*fyne.KeyEvent)
	// shortcut    fyne.ShortcutHandler

	scale, detectedScale, texScale float32

	context driver.WithContext
}

func (c *glCanvas) Capture() image.Image {
	var img image.Image
	runOnDraw(c.context.(*window), func() {
		img = c.Painter().Capture(c)
	})
	return img
}

func (c *glCanvas) Content() fyne.CanvasObject {
	c.RLock()
	retval := c.content
	c.RUnlock()
	return retval
}

func (c *glCanvas) DismissMenu() bool {
	c.RLock()
	menu := c.menu
	c.RUnlock()
	if menu != nil && menu.(*MenuBar).IsActive() {
		menu.(*MenuBar).Toggle()
		return true
	}
	return false
}

func (c *glCanvas) InteractiveArea() (fyne.Position, fyne.Size) {
	return fyne.Position{}, c.Size()
}

func (c *glCanvas) MinSize() fyne.Size {
	c.RLock()
	defer c.RUnlock()
	return c.canvasSize(c.content.MinSize())
}

func (c *glCanvas) OnKeyDown() func(*fyne.KeyEvent) {
	return c.onKeyDown
}

func (c *glCanvas) OnKeyUp() func(*fyne.KeyEvent) {
	return c.onKeyUp
}

func (c *glCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *glCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *glCanvas) Padded() bool {
	return c.padded
}

func (c *glCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	c.RLock()
	texScale := c.texScale
	scale := c.scale
	c.RUnlock()
	multiple := scale * texScale
	scaleInt := func(x float32) int {
		return int(math.Round(float64(x * multiple)))
	}

	return scaleInt(pos.X), scaleInt(pos.Y)
}

func (c *glCanvas) Resize(size fyne.Size) {
	// This might not be the ideal solution, but it effectively avoid the first frame to be blurry due to the
	// rounding of the size to the loower integer when scale == 1. It does not affect the other cases as far as we tested.
	// This can easily be seen with fyne/cmd/hello and a scale == 1 as the text will happear blurry without the following line.
	nearestSize := fyne.NewSize(float32(math.Ceil(float64(size.Width))), float32(math.Ceil(float64(size.Height))))

	c.Lock()
	c.size = nearestSize
	c.Unlock()

	for _, overlay := range c.Overlays().List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Refresh()
		} else {
			overlay.Resize(nearestSize)
		}
	}

	c.RLock()
	content := c.content
	contentSize := c.contentSize(nearestSize)
	contentPos := c.contentPos()
	menu := c.menu
	menuHeight := c.menuHeight()
	c.RUnlock()

	content.Resize(contentSize)
	content.Move(contentPos)

	if menu != nil {
		menu.Refresh()
		menu.Resize(fyne.NewSize(nearestSize.Width, menuHeight))
	}
}

func (c *glCanvas) Scale() float32 {
	c.RLock()
	defer c.RUnlock()
	return c.scale
}

func (c *glCanvas) SetContent(content fyne.CanvasObject) {
	content.Resize(content.MinSize()) // give it the space it wants then calculate the real min

	c.Lock()
	// the pass above makes some layouts wide enough to wrap, so we ask again what the true min is.
	newSize := c.size.Max(c.canvasSize(content.MinSize()))

	c.setContent(content)
	c.Unlock()

	c.Resize(newSize)
	c.SetDirty()
}

func (c *glCanvas) SetOnKeyDown(typed func(*fyne.KeyEvent)) {
	c.onKeyDown = typed
}

func (c *glCanvas) SetOnKeyUp(typed func(*fyne.KeyEvent)) {
	c.onKeyUp = typed
}

func (c *glCanvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *glCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *glCanvas) SetPadded(padded bool) {
	c.Lock()
	content := c.content
	c.padded = padded
	pos := c.contentPos()
	c.Unlock()

	content.Move(pos)
}

func (c *glCanvas) reloadScale() {
	w := c.context.(*window)
	w.viewLock.RLock()
	windowVisible := w.visible
	w.viewLock.RUnlock()
	if !windowVisible {
		return
	}

	c.Lock()
	c.scale = w.calculatedScale()
	c.Unlock()
	c.SetDirty()

	c.context.RescaleContext()
}

func (c *glCanvas) Size() fyne.Size {
	c.RLock()
	defer c.RUnlock()
	return c.size
}

func (c *glCanvas) ToggleMenu() {
	c.RLock()
	menu := c.menu
	c.RUnlock()
	if menu != nil {
		menu.(*MenuBar).Toggle()
	}
}

func (c *glCanvas) buildMenu(w *window, m *fyne.MainMenu) {
	c.Lock()
	defer c.Unlock()
	c.setMenuOverlay(nil)
	if m == nil {
		return
	}
	if hasNativeMenu() {
		setupNativeMenu(w, m)
	} else {
		c.setMenuOverlay(buildMenuOverlay(m, w))
	}
}

// canvasSize computes the needed canvas size for the given content size
func (c *glCanvas) canvasSize(contentSize fyne.Size) fyne.Size {
	canvasSize := contentSize.Add(fyne.NewSize(0, c.menuHeight()))
	if c.Padded() {
		return canvasSize.Add(fyne.NewSquareSize(theme.Padding() * 2))
	}
	return canvasSize
}

func (c *glCanvas) contentPos() fyne.Position {
	contentPos := fyne.NewPos(0, c.menuHeight())
	if c.Padded() {
		return contentPos.Add(fyne.NewSquareOffsetPos(theme.Padding()))
	}
	return contentPos
}

func (c *glCanvas) contentSize(canvasSize fyne.Size) fyne.Size {
	contentSize := fyne.NewSize(canvasSize.Width, canvasSize.Height-c.menuHeight())
	if c.Padded() {
		return contentSize.Subtract(fyne.NewSquareSize(theme.Padding() * 2))
	}
	return contentSize
}

func (c *glCanvas) menuHeight() float32 {
	if c.menu == nil {
		return 0 // no menu or native menu -> does not consume space on the canvas
	}

	return c.menu.MinSize().Height
}

func (c *glCanvas) overlayChanged() {
	c.SetDirty()
}

func (c *glCanvas) paint(size fyne.Size) {
	clips := &internal.ClipStack{}
	if c.Content() == nil {
		return
	}
	c.Painter().Clear()

	paint := func(node *common.RenderCacheNode, pos fyne.Position) {
		obj := node.Obj()
		if _, ok := obj.(fyne.Scrollable); ok {
			inner := clips.Push(pos, obj.Size())
			c.Painter().StartClipping(inner.Rect())
		}
		if size.Width <= 0 || size.Height <= 0 { // iconifying on Windows can do bad things
			return
		}
		c.Painter().Paint(obj, pos, size)
	}
	afterPaint := func(node *common.RenderCacheNode, pos fyne.Position) {
		if _, ok := node.Obj().(fyne.Scrollable); ok {
			clips.Pop()
			if top := clips.Top(); top != nil {
				c.Painter().StartClipping(top.Rect())
			} else {
				c.Painter().StopClipping()
			}
		}

		if c.debug {
			c.DrawDebugOverlay(node.Obj(), pos, size)
		}
	}
	c.WalkTrees(paint, afterPaint)
}

func (c *glCanvas) setContent(content fyne.CanvasObject) {
	c.content = content
	c.SetContentTreeAndFocusMgr(content)
}

func (c *glCanvas) setMenuOverlay(b fyne.CanvasObject) {
	c.menu = b
	c.SetMenuTreeAndFocusMgr(b)

	if c.menu != nil && !c.size.IsZero() {
		c.content.Resize(c.contentSize(c.size))
		c.content.Move(c.contentPos())

		c.menu.Refresh()
		c.menu.Resize(fyne.NewSize(c.size.Width, c.menu.MinSize().Height))
	}
}

func (c *glCanvas) applyThemeOutOfTreeObjects() {
	c.RLock()
	menu := c.menu
	padded := c.padded
	c.RUnlock()
	if menu != nil {
		app.ApplyThemeTo(menu, c) // Ensure our menu gets the theme change message as it's out-of-tree
	}

	c.SetPadded(padded) // refresh the padding for potential theme differences
}

func newCanvas() *glCanvas {
	c := &glCanvas{scale: 1.0, texScale: 1.0, padded: true}
	c.Initialize(c, c.overlayChanged)
	c.setContent(&canvas.Rectangle{FillColor: theme.BackgroundColor()})
	c.debug = fyne.CurrentApp().Settings().BuildType() == fyne.BuildDebug
	return c
}
