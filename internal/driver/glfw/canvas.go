package glfw

import (
	"fmt"
	"image"
	"math"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Declare conformity with Canvas interface
var _ fyne.Canvas = (*glCanvas)(nil)

type glCanvas struct {
	sync.RWMutex
	content, overlay fyne.CanvasObject
	menu             *widget.Toolbar
	padded           bool
	size             fyne.Size
	focused          fyne.Focusable
	focusMgr         *app.FocusManager

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
	onKeyDown   func(*fyne.KeyEvent)
	onKeyUp     func(*fyne.KeyEvent)
	shortcut    fyne.ShortcutHandler

	scale, detectedScale float32
	painter              gl.Painter

	dirty        bool
	dirtyMutex   *sync.Mutex
	minSizes     map[fyne.CanvasObject]fyne.Size
	refreshQueue chan fyne.CanvasObject
	context      driver.WithContext
}

func scaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(math.Round(float64(v) * float64(c.Scale())))
	}
}

func unscaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}

func (c *glCanvas) Capture() image.Image {
	var img image.Image
	runOnMain(func() {
		img = c.painter.Capture(c)
	})
	return img
}

func (c *glCanvas) Content() fyne.CanvasObject {
	c.RLock()
	retval := c.content
	c.RUnlock()
	return retval
}

func (c *glCanvas) SetContent(content fyne.CanvasObject) {
	c.Lock()
	c.content = content
	c.minSizes = map[fyne.CanvasObject]fyne.Size{}
	c.Unlock()

	newSize := c.size.Union(c.canvasSize(c.content.MinSize()))
	c.Resize(newSize)

	c.setDirty(true)
}

func (c *glCanvas) Overlay() fyne.CanvasObject {
	c.RLock()
	retval := c.overlay
	c.RUnlock()
	return retval
}

func (c *glCanvas) SetOverlay(overlay fyne.CanvasObject) {
	c.Lock()
	c.overlay = overlay
	c.Unlock()

	if overlay != nil {
		c.overlay.Resize(c.Size())
	}
	c.setDirty(true)
}

func (c *glCanvas) Padded() bool {
	return c.padded
}

func (c *glCanvas) SetPadded(padded bool) {
	c.padded = padded

	c.content.Move(c.contentPos())
}

func (c *glCanvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
	c.setDirty(true)
}

func (c *glCanvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		c.focused.(fyne.Focusable).FocusLost()
	}

	c.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (c *glCanvas) Unfocus() {
	if c.focused != nil {
		c.focused.(fyne.Focusable).FocusLost()
	}
	c.focused = nil
}

func (c *glCanvas) Focused() fyne.Focusable {
	return c.focused
}

func (c *glCanvas) Resize(size fyne.Size) {
	c.size = size
	c.content.Resize(c.contentSize(size))
	c.content.Move(c.contentPos())

	if c.overlay != nil {
		c.overlay.Resize(size)
	}
	if c.menu != nil {
		c.menu.Resize(fyne.NewSize(size.Width, c.menu.MinSize().Height))
	}
	c.Refresh(c.content)
}

func (c *glCanvas) Size() fyne.Size {
	return c.size
}

func (c *glCanvas) MinSize() fyne.Size {
	return c.canvasSize(c.content.MinSize())
}

func (c *glCanvas) Scale() float32 {
	return c.scale
}

func (c *glCanvas) SetScale(scale float32) {
	if scale == fyne.SettingsScaleAuto {
		c.scale = c.detectedScale
	} else {
		c.scale = scale
	}
	c.setDirty(true)

	c.context.RescaleContext()
}

func (c *glCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *glCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *glCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *glCanvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}

func (c *glCanvas) OnKeyDown() func(*fyne.KeyEvent) {
	return c.onKeyDown
}

func (c *glCanvas) SetOnKeyDown(typed func(*fyne.KeyEvent)) {
	c.onKeyDown = typed
}

func (c *glCanvas) OnKeyUp() func(*fyne.KeyEvent) {
	return c.onKeyUp
}

func (c *glCanvas) SetOnKeyUp(typed func(*fyne.KeyEvent)) {
	c.onKeyUp = typed
}

func (c *glCanvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

func (c *glCanvas) ensureMinSize() bool {
	if c.Content() == nil {
		return false
	}
	var objToLayout fyne.CanvasObject
	windowNeedsMinSizeUpdate := false
	ensureMinSize := func(obj fyne.CanvasObject, parent fyne.CanvasObject) {
		canvasMutex.Lock()
		canvases[obj] = c
		canvasMutex.Unlock()

		if !obj.Visible() {
			return
		}
		minSize := obj.MinSize()
		minSizeChanged := c.minSizes[obj] != minSize
		if minSizeChanged {
			c.minSizes[obj] = minSize
			if parent != nil {
				objToLayout = parent
			} else {
				windowNeedsMinSizeUpdate = true
				size := obj.Size()
				expectedSize := minSize.Union(size)
				if expectedSize != size {
					objToLayout = nil
					obj.Resize(expectedSize)
				}
			}
		}
		if obj == objToLayout {
			switch cont := obj.(type) {
			case *fyne.Container:
				if cont.Layout != nil {
					cont.Layout.Layout(cont.Objects, cont.Size())
				}
			case fyne.Widget:
				widget.Renderer(cont).Layout(cont.Size())
			default:
				fmt.Printf("implementation error - unknown container type: %T\n", cont)
			}
		}
	}
	c.walkTree(nil, ensureMinSize)
	if windowNeedsMinSizeUpdate && (c.size.Width < c.MinSize().Width || c.size.Height < c.MinSize().Height) {
		c.Resize(c.Size().Union(c.MinSize()))
	}
	return windowNeedsMinSizeUpdate
}

func (c *glCanvas) paint(size fyne.Size) {
	if c.Content() == nil {
		return
	}
	c.setDirty(false)
	c.painter.Clear()

	paint := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		// TODO should this be somehow not scroll container specific?
		if _, ok := obj.(*widget.ScrollContainer); ok {
			c.painter.StartClipping(
				fyne.NewPos(pos.X, c.Size().Height-pos.Y-obj.Size().Height),
				obj.Size(),
			)
		}
		c.painter.Paint(obj, pos, size)
		return false
	}
	afterPaint := func(obj, _ fyne.CanvasObject) {
		if _, ok := obj.(*widget.ScrollContainer); ok {
			c.painter.StopClipping()
		}
	}

	c.walkTree(paint, afterPaint)
}

func (c *glCanvas) walkTree(
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) {
	driver.WalkVisibleObjectTree(c.content, beforeChildren, afterChildren)
	if c.menu != nil {
		driver.WalkVisibleObjectTree(c.menu, beforeChildren, afterChildren)
	}
	if c.overlay != nil {
		driver.WalkVisibleObjectTree(c.overlay, beforeChildren, afterChildren)
	}
}

func (c *glCanvas) setDirty(dirty bool) {
	c.dirtyMutex.Lock()
	defer c.dirtyMutex.Unlock()

	c.dirty = dirty
}

func (c *glCanvas) isDirty() bool {
	c.dirtyMutex.Lock()
	defer c.dirtyMutex.Unlock()

	return c.dirty
}

func (c *glCanvas) setupThemeListener() {
	listener := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			if c.menu != nil {
				app.ApplyThemeTo(c.menu, c) // Ensure our menu gets the theme change message as it's out-of-tree
			}

			c.SetPadded(c.padded) // refresh the padding for potential theme differences
		}
	}()
}

func (c *glCanvas) buildMenuBar(m *fyne.MainMenu) {
	c.setMenuBar(nil)
	if m == nil {
		return
	}
	if hasNativeMenu() {
		setupNativeMenu(m)
	} else {
		c.setMenuBar(buildMenuBar(m, c))
	}
}

func (c *glCanvas) setMenuBar(b *widget.Toolbar) {
	c.Lock()
	c.menu = b
	c.Unlock()
}

func (c *glCanvas) menuBar() *widget.Toolbar {
	c.RLock()
	defer c.RUnlock()
	return c.menu
}

func (c *glCanvas) menuHeight() int {
	switch c.menuBar() {
	case nil:
		// no menu or native menu -> does not consume space on the canvas
		return 0
	default:
		return c.menuBar().MinSize().Height
	}
}

// canvasSize computes the needed canvas size for the given content size
func (c *glCanvas) canvasSize(contentSize fyne.Size) fyne.Size {
	canvasSize := contentSize.Add(fyne.NewSize(0, c.menuHeight()))
	if c.Padded() {
		pad := theme.Padding() * 2
		canvasSize = canvasSize.Add(fyne.NewSize(pad, pad))
	}
	return canvasSize
}

func (c *glCanvas) contentSize(canvasSize fyne.Size) fyne.Size {
	contentSize := fyne.NewSize(canvasSize.Width, canvasSize.Height-c.menuHeight())
	if c.Padded() {
		pad := theme.Padding() * 2
		contentSize = contentSize.Subtract(fyne.NewSize(pad, pad))
	}
	return contentSize
}

func (c *glCanvas) contentPos() fyne.Position {
	contentPos := fyne.NewPos(0, c.menuHeight())
	if c.Padded() {
		contentPos = contentPos.Add(fyne.NewPos(theme.Padding(), theme.Padding()))
	}
	return contentPos
}

func newCanvas() *glCanvas {
	c := &glCanvas{scale: 1.0}
	c.content = &canvas.Rectangle{FillColor: theme.BackgroundColor()}
	c.minSizes = map[fyne.CanvasObject]fyne.Size{}
	c.padded = true

	c.focusMgr = app.NewFocusManager(c)
	c.refreshQueue = make(chan fyne.CanvasObject, 1024)
	c.dirtyMutex = &sync.Mutex{}

	c.setupThemeListener()

	return c
}
