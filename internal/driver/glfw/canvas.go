package glfw

import (
	"image"
	"math"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Declare conformity with Canvas interface
var _ fyne.Canvas = (*glCanvas)(nil)

type glCanvas struct {
	sync.RWMutex

	content         fyne.CanvasObject
	contentFocusMgr *app.FocusManager
	menu            fyne.CanvasObject
	menuFocusMgr    *app.FocusManager
	overlays        *overlayStack
	padded          bool
	size            fyne.Size

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)
	onKeyDown   func(*fyne.KeyEvent)
	onKeyUp     func(*fyne.KeyEvent)
	shortcut    fyne.ShortcutHandler

	scale, detectedScale, texScale float32
	painter                        gl.Painter

	dirty                 bool
	dirtyMutex            *sync.Mutex
	refreshQueue          chan fyne.CanvasObject
	contentTree, menuTree *renderCacheTree
	context               driver.WithContext
}

func (c *glCanvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
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

func (c *glCanvas) Focus(obj fyne.Focusable) {
	c.focusManager().Focus(obj)
}

func (c *glCanvas) Focused() fyne.Focusable {
	return c.focusManager().Focused()
}

func (c *glCanvas) FocusGained() {
	c.focusManager().FocusGained()
}

func (c *glCanvas) FocusLost() {
	c.focusManager().FocusLost()
}

func (c *glCanvas) FocusNext() {
	c.focusManager().FocusNext()
}

func (c *glCanvas) FocusPrevious() {
	c.focusManager().FocusPrevious()
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

// Deprecated: Use Overlays() instead.
func (c *glCanvas) Overlay() fyne.CanvasObject {
	return c.Overlays().Top()
}

func (c *glCanvas) Overlays() fyne.OverlayStack {
	c.RLock()
	defer c.RUnlock()
	return c.overlays
}

func (c *glCanvas) Padded() bool {
	return c.padded
}

func (c *glCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	texScale := c.texScale
	multiple := float64(c.Scale() * texScale)
	scaleInt := func(x int) int {
		return int(math.Round(float64(x) * multiple))
	}

	return scaleInt(pos.X), scaleInt(pos.Y)
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

func (c *glCanvas) Resize(size fyne.Size) {
	c.Lock()
	c.size = size
	c.Unlock()

	for _, overlay := range c.overlays.List() {
		if p, ok := overlay.(*widget.PopUp); ok {
			// TODO: remove this when #707 is being addressed.
			// “Notifies” the PopUp of the canvas size change.
			p.Resize(p.Content.Size().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		} else {
			overlay.Resize(size)
		}
	}

	c.RLock()
	c.content.Resize(c.contentSize(size))
	c.content.Move(c.contentPos())

	if c.menu != nil {
		c.menu.Refresh()
		c.menu.Resize(fyne.NewSize(size.Width, c.menu.MinSize().Height))
	}
	c.RUnlock()
}

func (c *glCanvas) Scale() float32 {
	c.RLock()
	defer c.RUnlock()
	return c.scale
}

func (c *glCanvas) SetContent(content fyne.CanvasObject) {
	c.Lock()
	c.setContent(content)

	c.content.Resize(c.content.MinSize()) // give it the space it wants then calculate the real min
	// the pass above makes some layouts wide enough to wrap, so we ask again what the true min is.
	newSize := c.size.Max(c.canvasSize(c.content.MinSize()))
	c.Unlock()

	c.Resize(newSize)
	c.setDirty(true)
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

// Deprecated: Use Overlays() instead.
func (c *glCanvas) SetOverlay(overlay fyne.CanvasObject) {
	c.RLock()
	o := c.overlays
	c.RUnlock()
	o.setOverlay(overlay)
}

func (c *glCanvas) SetPadded(padded bool) {
	c.Lock()
	defer c.Unlock()
	c.padded = padded

	c.content.Move(c.contentPos())
}

// SetScale sets the render scale for this specific canvas
//
// Deprecated: Settings are now calculated solely on the user configuration and system setup.
func (c *glCanvas) SetScale(_ float32) {
	if !c.context.(*window).visible {
		return
	}

	c.Lock()
	c.scale = c.context.(*window).calculatedScale()
	c.Unlock()
	c.setDirty(true)

	c.context.RescaleContext()
}

func (c *glCanvas) Size() fyne.Size {
	c.RLock()
	defer c.RUnlock()
	return c.size
}

func (c *glCanvas) Unfocus() {
	c.focusManager().Focus(nil)
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
		c.setMenuOverlay(buildMenuOverlay(m, c))
	}
}

func (c *glCanvas) RemoveShortcut(shortcut fyne.Shortcut) {
	c.shortcut.RemoveShortcut(shortcut)
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

func (c *glCanvas) contentPos() fyne.Position {
	contentPos := fyne.NewPos(0, c.menuHeight())
	if c.Padded() {
		contentPos = contentPos.Add(fyne.NewPos(theme.Padding(), theme.Padding()))
	}
	return contentPos
}

func (c *glCanvas) contentSize(canvasSize fyne.Size) fyne.Size {
	contentSize := fyne.NewSize(canvasSize.Width, canvasSize.Height-c.menuHeight())
	if c.Padded() {
		pad := theme.Padding() * 2
		contentSize = contentSize.Subtract(fyne.NewSize(pad, pad))
	}
	return contentSize
}

func (c *glCanvas) ensureMinSize() bool {
	if c.Content() == nil {
		return false
	}
	var lastParent fyne.CanvasObject

	windowNeedsMinSizeUpdate := false
	ensureMinSize := func(node *renderCacheNode) {
		obj := node.obj
		canvasMutex.Lock()
		canvases[obj] = c
		canvasMutex.Unlock()

		if !obj.Visible() {
			return
		}
		minSize := obj.MinSize()
		minSizeChanged := node.minSize != minSize
		if minSizeChanged {
			objToLayout := obj
			node.minSize = minSize
			if node.parent != nil {
				objToLayout = node.parent.obj
			} else {
				windowNeedsMinSizeUpdate = true
				size := obj.Size()
				expectedSize := minSize.Max(size)
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
	}
	c.walkTrees(nil, ensureMinSize)

	min := c.MinSize()
	c.RLock()
	shouldResize := windowNeedsMinSizeUpdate && (c.size.Width < min.Width || c.size.Height < min.Height)
	c.RUnlock()
	if shouldResize {
		c.Resize(c.Size().Max(c.MinSize()))
	}

	if lastParent != nil {
		c.RLock()
		updateLayout(lastParent)
		c.RUnlock()
	}
	return windowNeedsMinSizeUpdate
}

func (c *glCanvas) focusManager() *app.FocusManager {
	c.RLock()
	defer c.RUnlock()
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	if c.isMenuActive() {
		return c.menuFocusMgr
	}
	return c.contentFocusMgr
}

func (c *glCanvas) isDirty() bool {
	c.dirtyMutex.Lock()
	defer c.dirtyMutex.Unlock()

	return c.dirty
}

func (c *glCanvas) isMenuActive() bool {
	return c.menu != nil && c.menu.(*MenuBar).IsActive()
}

func (c *glCanvas) menuHeight() int {
	switch c.menu {
	case nil:
		// no menu or native menu -> does not consume space on the canvas
		return 0
	default:
		return c.menu.MinSize().Height
	}
}

func (c *glCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+2)
	trees = append(trees, c.content)
	if c.menu != nil {
		trees = append(trees, c.menu)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

func (c *glCanvas) overlayChanged() {
	c.Lock()
	defer c.Unlock()
	c.dirty = true
}

func (c *glCanvas) paint(size fyne.Size) {
	if c.Content() == nil {
		return
	}
	c.setDirty(false)
	c.painter.Clear()

	paint := func(node *renderCacheNode, pos fyne.Position) {
		obj := node.obj
		// TODO should this be somehow not scroll container specific?
		if _, ok := obj.(*widget.ScrollContainer); ok {
			c.painter.StartClipping(
				fyne.NewPos(pos.X, c.Size().Height-pos.Y-obj.Size().Height),
				obj.Size(),
			)
		}
		c.painter.Paint(obj, pos, size)
	}
	afterPaint := func(node *renderCacheNode) {
		if _, ok := node.obj.(*widget.ScrollContainer); ok {
			c.painter.StopClipping()
		}
	}

	c.walkTrees(paint, afterPaint)
}

func (c *glCanvas) setContent(content fyne.CanvasObject) {
	c.content = content
	c.contentTree = &renderCacheTree{root: &renderCacheNode{obj: c.content}}
	c.contentFocusMgr = app.NewFocusManager(c.content)
}

func (c *glCanvas) setDirty(dirty bool) {
	c.dirtyMutex.Lock()
	defer c.dirtyMutex.Unlock()

	c.dirty = dirty
}

func (c *glCanvas) setMenuOverlay(b fyne.CanvasObject) {
	c.menu = b
	c.menuTree = &renderCacheTree{root: &renderCacheNode{obj: c.menu}}
	c.menuFocusMgr = app.NewFocusManager(c.menu)
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

func (c *glCanvas) walkTree(
	tree *renderCacheTree,
	beforeChildren func(*renderCacheNode, fyne.Position),
	afterChildren func(*renderCacheNode),
) {
	tree.Lock()
	defer tree.Unlock()
	var node, parent, prev *renderCacheNode
	node = tree.root

	bc := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if node != nil && node.obj != obj {
			if parent.firstChild == node {
				parent.firstChild = nil
			}
			node = nil
		}
		if node == nil {
			node = &renderCacheNode{parent: parent, obj: obj}
			if parent.firstChild == nil {
				parent.firstChild = node
			} else {
				prev.nextSibling = node
			}
		}
		if prev != nil && prev.parent != parent {
			prev = nil
		}

		if beforeChildren != nil {
			beforeChildren(node, pos)
		}

		parent = node
		node = parent.firstChild
		return false
	}
	ac := func(obj fyne.CanvasObject, _ fyne.CanvasObject) {
		node = parent
		parent = node.parent
		if prev != nil && prev.parent != parent {
			prev.nextSibling = nil
		}

		if afterChildren != nil {
			afterChildren(node)
		}

		prev = node
		node = node.nextSibling
	}
	driver.WalkVisibleObjectTree(tree.root.obj, bc, ac)
}

func (c *glCanvas) walkTrees(
	beforeChildren func(*renderCacheNode, fyne.Position),
	afterChildren func(*renderCacheNode),
) {
	c.walkTree(c.contentTree, beforeChildren, afterChildren)
	if c.menu != nil {
		c.walkTree(c.menuTree, beforeChildren, afterChildren)
	}
	for _, tree := range c.overlays.renderCaches {
		if tree != nil {
			c.walkTree(tree, beforeChildren, afterChildren)
		}
	}
}

type overlayStack struct {
	internal.OverlayStack

	propertyLock sync.RWMutex
	renderCaches []*renderCacheTree
}

func (o *overlayStack) Add(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()
	o.add(overlay)
}

func (o *overlayStack) Remove(overlay fyne.CanvasObject) {
	if overlay == nil || len(o.List()) == 0 {
		return
	}
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()
	o.remove(overlay)
}

func (o *overlayStack) add(overlay fyne.CanvasObject) {
	o.renderCaches = append(o.renderCaches, &renderCacheTree{root: &renderCacheNode{obj: overlay}})
	o.OverlayStack.Add(overlay)
}

func (o *overlayStack) remove(overlay fyne.CanvasObject) {
	o.OverlayStack.Remove(overlay)
	overlayCount := len(o.List())
	o.renderCaches = o.renderCaches[:overlayCount]
}

// concurrency safe implementation of deprecated c.SetOverlay
func (o *overlayStack) setOverlay(overlay fyne.CanvasObject) {
	o.propertyLock.Lock()
	defer o.propertyLock.Unlock()

	if len(o.List()) > 0 {
		o.remove(o.List()[0])
	}
	if overlay != nil {
		o.add(overlay)
	}
}

type renderCacheNode struct {
	// structural data
	firstChild  *renderCacheNode
	nextSibling *renderCacheNode
	obj         fyne.CanvasObject
	parent      *renderCacheNode
	// cache data
	minSize fyne.Size
	// painterData is some data from the painter associated with the drawed node
	// it may for instance point to a GL texture
	// it should free all associated resources when released
	// i.e. it should not simply be a texture reference integer
	painterData interface{}
}

type renderCacheTree struct {
	sync.RWMutex
	root *renderCacheNode
}

func newCanvas() *glCanvas {
	c := &glCanvas{scale: 1.0, texScale: 1.0}
	c.setContent(&canvas.Rectangle{FillColor: theme.BackgroundColor()})
	c.padded = true

	c.overlays = &overlayStack{
		OverlayStack: internal.OverlayStack{
			OnChange: c.overlayChanged,
			Canvas:   c,
		},
	}

	c.refreshQueue = make(chan fyne.CanvasObject, 4096)
	c.dirtyMutex = &sync.Mutex{}

	c.setupThemeListener()

	return c
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
