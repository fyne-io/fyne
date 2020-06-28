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

	content  fyne.CanvasObject
	menu     fyne.CanvasObject
	overlays *overlayStack
	padded   bool
	size     fyne.Size
	focused  fyne.Focusable
	focusMgr *app.FocusManager

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

type renderCacheTree struct {
	sync.RWMutex
	root *renderCacheNode
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

type overlayStack struct {
	internal.OverlayStack

	onChange     func()
	renderCaches []*renderCacheTree
}

func (o *overlayStack) Add(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	o.OverlayStack.Add(overlay)
	o.renderCaches = append(o.renderCaches, &renderCacheTree{root: &renderCacheNode{obj: overlay}})
	o.onChange()
}

func (o *overlayStack) Remove(overlay fyne.CanvasObject) {
	o.OverlayStack.Remove(overlay)
	o.renderCaches = o.renderCaches[:len(o.List())]
	o.onChange()
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
	c.contentTree = &renderCacheTree{root: &renderCacheNode{obj: c.content}}

	c.content.Resize(c.content.MinSize()) // give it the space it wants then calculate the real min
	// the pass above makes some layouts wide enough to wrap, so we ask again what the true min is.
	newSize := c.size.Union(c.canvasSize(c.content.MinSize()))
	c.Unlock()

	c.Resize(newSize)
	c.setDirty(true)
}

// Deprecated: Use Overlays() instead.
func (c *glCanvas) Overlay() fyne.CanvasObject {
	c.RLock()
	defer c.RUnlock()
	return c.overlays.Top()
}

func (c *glCanvas) Overlays() fyne.OverlayStack {
	c.RLock()
	defer c.RUnlock()
	return c.overlays
}

// Deprecated: Use Overlays() instead.
func (c *glCanvas) SetOverlay(overlay fyne.CanvasObject) {
	c.Lock()
	defer c.Unlock()
	if len(c.overlays.List()) > 0 {
		c.overlays.Remove(c.overlays.List()[0])
	}
	c.overlays.Add(overlay)
	c.setDirty(true)
}

func (c *glCanvas) Padded() bool {
	return c.padded
}

func (c *glCanvas) SetPadded(padded bool) {
	c.Lock()
	defer c.Unlock()
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
	c.RLock()
	focused := c.focused
	c.RUnlock()

	if focused == obj || obj == nil {
		return
	}

	if dis, ok := obj.(fyne.Disableable); ok && dis.Disabled() {
		c.Unfocus()
		return
	}

	if focused != nil {
		focused.FocusLost()
	}

	c.Lock()
	c.focused = obj
	c.Unlock()
	obj.FocusGained()
}

func (c *glCanvas) Unfocus() {
	c.Lock()
	focused := c.focused
	c.focused = nil
	c.Unlock()

	if focused != nil {
		focused.FocusLost()
	}
}

func (c *glCanvas) Focused() fyne.Focusable {
	c.RLock()
	defer c.RUnlock()
	return c.focused
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

func (c *glCanvas) Size() fyne.Size {
	c.RLock()
	defer c.RUnlock()
	return c.size
}

func (c *glCanvas) MinSize() fyne.Size {
	c.RLock()
	defer c.RUnlock()
	return c.canvasSize(c.content.MinSize())
}

func (c *glCanvas) Scale() float32 {
	c.RLock()
	defer c.RUnlock()
	return c.scale
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

func (c *glCanvas) PixelCoordinateForPosition(pos fyne.Position) (int, int) {
	texScale := c.texScale
	multiple := float64(c.Scale() * texScale)
	scaleInt := func(x int) int {
		return int(math.Round(float64(x) * multiple))
	}

	return scaleInt(pos.X), scaleInt(pos.Y)
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

func (c *glCanvas) objectTrees() []fyne.CanvasObject {
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+2)
	trees = append(trees, c.content)
	if c.menu != nil {
		trees = append(trees, c.menu)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
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
	}
	c.walkTrees(nil, ensureMinSize)

	min := c.MinSize()
	c.RLock()
	shouldResize := windowNeedsMinSizeUpdate && (c.size.Width < min.Width || c.size.Height < min.Height)
	c.RUnlock()
	if shouldResize {
		c.Resize(c.Size().Union(c.MinSize()))
	}

	if lastParent != nil {
		c.RLock()
		updateLayout(lastParent)
		c.RUnlock()
	}
	return windowNeedsMinSizeUpdate
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

func (c *glCanvas) buildMenu(w *window, m *fyne.MainMenu) {
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

func (c *glCanvas) setMenuOverlay(b fyne.CanvasObject) {
	c.Lock()
	c.menu = b
	c.menuTree = &renderCacheTree{root: &renderCacheNode{obj: c.menu}}
	c.Unlock()
}

func (c *glCanvas) menuOverlay() fyne.CanvasObject {
	return c.menu
}

func (c *glCanvas) menuHeight() int {
	switch c.menuOverlay() {
	case nil:
		// no menu or native menu -> does not consume space on the canvas
		return 0
	default:
		return c.menuOverlay().MinSize().Height
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
	c := &glCanvas{scale: 1.0, texScale: 1.0}
	c.content = &canvas.Rectangle{FillColor: theme.BackgroundColor()}
	c.contentTree = &renderCacheTree{root: &renderCacheNode{obj: c.content}}
	c.padded = true

	c.overlays = &overlayStack{onChange: func() { c.setDirty(true) }}

	c.focusMgr = app.NewFocusManager(c)
	c.refreshQueue = make(chan fyne.CanvasObject, 4096)
	c.dirtyMutex = &sync.Mutex{}

	c.setupThemeListener()

	return c
}
