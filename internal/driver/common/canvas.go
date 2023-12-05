package common

import (
	"image/color"
	"reflect"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter/gl"
)

// SizeableCanvas defines a canvas with size related functions.
type SizeableCanvas interface {
	fyne.Canvas
	Resize(fyne.Size)
	MinSize() fyne.Size
}

// Canvas defines common canvas implementation.
type Canvas struct {
	sync.RWMutex

	OnFocus   func(obj fyne.Focusable)
	OnUnfocus func()

	impl SizeableCanvas

	contentFocusMgr *app.FocusManager
	menuFocusMgr    *app.FocusManager
	overlays        *overlayStack

	shortcut fyne.ShortcutHandler

	painter gl.Painter

	// Any object that requestes to enter to the refresh queue should
	// not be omitted as it is always a rendering task's decision
	// for skipping frames or drawing calls.
	//
	// If an object failed to ender the refresh queue, the object may
	// disappear or blink from the view at any frames. As of this reason,
	// the refreshQueue is an unbounded queue which is bale to cache
	// arbitrary number of fyne.CanvasObject for the rendering.
	refreshQueue *async.CanvasObjectQueue
	dirty        uint32 // atomic

	mWindowHeadTree, contentTree, menuTree *renderCacheTree
}

// AddShortcut adds a shortcut to the canvas.
func (c *Canvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

func (c *Canvas) DrawDebugOverlay(obj fyne.CanvasObject, pos fyne.Position, size fyne.Size) {
	switch obj.(type) {
	case fyne.Widget:
		r := canvas.NewRectangle(color.Transparent)
		r.StrokeColor = color.NRGBA{R: 0xcc, G: 0x33, B: 0x33, A: 0xff}
		r.StrokeWidth = 1
		r.Resize(obj.Size())
		c.Painter().Paint(r, pos, size)

		t := canvas.NewText(reflect.ValueOf(obj).Elem().Type().Name(), r.StrokeColor)
		t.TextSize = 10
		c.Painter().Paint(t, pos.AddXY(2, 2), size)
	case *fyne.Container:
		r := canvas.NewRectangle(color.Transparent)
		r.StrokeColor = color.NRGBA{R: 0x33, G: 0x33, B: 0xcc, A: 0xff}
		r.StrokeWidth = 1
		r.Resize(obj.Size())
		c.Painter().Paint(r, pos, size)
	}
}

// EnsureMinSize ensure canvas min size.
//
// This function uses lock.
func (c *Canvas) EnsureMinSize() bool {
	if c.impl.Content() == nil {
		return false
	}
	windowNeedsMinSizeUpdate := false
	csize := c.impl.Size()
	min := c.impl.MinSize()

	c.RLock()
	defer c.RUnlock()

	var parentNeedingUpdate *RenderCacheNode

	ensureMinSize := func(node *RenderCacheNode, pos fyne.Position) {
		obj := node.obj
		cache.SetCanvasForObject(obj, c.impl, func() {
			if img, ok := obj.(*canvas.Image); ok {
				c.RUnlock()
				img.Refresh() // this may now have a different texScale
				c.RLock()
			}
		})

		if parentNeedingUpdate == node {
			c.updateLayout(obj)
			parentNeedingUpdate = nil
		}

		c.RUnlock()
		if !obj.Visible() {
			c.RLock()
			return
		}
		minSize := obj.MinSize()
		c.RLock()

		minSizeChanged := node.minSize != minSize
		if minSizeChanged {
			node.minSize = minSize
			if node.parent != nil {
				parentNeedingUpdate = node.parent
			} else {
				windowNeedsMinSizeUpdate = true
				c.RUnlock()
				size := obj.Size()
				c.RLock()
				expectedSize := minSize.Max(size)
				if expectedSize != size && size != csize {
					c.RUnlock()
					obj.Resize(expectedSize)
					c.RLock()
				} else {
					c.updateLayout(obj)
				}
			}
		}
	}
	c.WalkTrees(nil, ensureMinSize)

	shouldResize := windowNeedsMinSizeUpdate && (csize.Width < min.Width || csize.Height < min.Height)
	if shouldResize {
		c.RUnlock()
		c.impl.Resize(csize.Max(min))
		c.RLock()
	}
	return windowNeedsMinSizeUpdate
}

// Focus makes the provided item focused.
func (c *Canvas) Focus(obj fyne.Focusable) {
	focusMgr := c.focusManager()
	if focusMgr != nil && focusMgr.Focus(obj) { // fast path – probably >99.9% of all cases
		if c.OnFocus != nil {
			c.OnFocus(obj)
		}
		return
	}

	c.RLock()
	focusMgrs := append([]*app.FocusManager{c.contentFocusMgr, c.menuFocusMgr}, c.overlays.ListFocusManagers()...)
	c.RUnlock()

	for _, mgr := range focusMgrs {
		if mgr == nil {
			continue
		}
		if focusMgr != mgr {
			if mgr.Focus(obj) {
				if c.OnFocus != nil {
					c.OnFocus(obj)
				}
				return
			}
		}
	}

	fyne.LogError("Failed to focus object which is not part of the canvas’ content, menu or overlays.", nil)
}

// Focused returns the current focused object.
func (c *Canvas) Focused() fyne.Focusable {
	mgr := c.focusManager()
	if mgr == nil {
		return nil
	}
	return mgr.Focused()
}

// FocusGained signals to the manager that its content got focus.
// Valid only on Desktop.
func (c *Canvas) FocusGained() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusGained()
}

// FocusLost signals to the manager that its content lost focus.
// Valid only on Desktop.
func (c *Canvas) FocusLost() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusLost()
}

// FocusNext focuses the next focusable item.
func (c *Canvas) FocusNext() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusNext()
}

// FocusPrevious focuses the previous focusable item.
func (c *Canvas) FocusPrevious() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	mgr.FocusPrevious()
}

// FreeDirtyTextures frees dirty textures and returns the number of freed textures.
func (c *Canvas) FreeDirtyTextures() (freed uint64) {
	freeObject := func(object fyne.CanvasObject) {
		freeWalked := func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
			// No image refresh while recursing to avoid double texture upload.
			if _, ok := obj.(*canvas.Image); ok {
				return false
			}
			if c.painter != nil {
				c.painter.Free(obj)
			}
			return false
		}

		// Image.Refresh will trigger a refresh specific to the object, while recursing on parent widget would just lead to
		// a double texture upload.
		if img, ok := object.(*canvas.Image); ok {
			if c.painter != nil {
				c.painter.Free(img)
			}
		} else {
			driver.WalkCompleteObjectTree(object, freeWalked, nil)
		}
	}

	// Within a frame, refresh tasks are requested from the Refresh method,
	// and we desire to clear out all requested operations within a frame.
	// See https://github.com/fyne-io/fyne/issues/2548.
	tasksToDo := c.refreshQueue.Len()

	shouldFilterDuplicates := (tasksToDo > 200) // filtering has overhead, not worth enabling for few tasks
	var refreshSet map[fyne.CanvasObject]struct{}
	if shouldFilterDuplicates {
		refreshSet = make(map[fyne.CanvasObject]struct{})
	}

	for c.refreshQueue.Len() > 0 {
		object := c.refreshQueue.Out()
		if !shouldFilterDuplicates {
			freed++
			freeObject(object)
		} else {
			refreshSet[object] = struct{}{}
			tasksToDo--
			if tasksToDo == 0 {
				shouldFilterDuplicates = false // stop collecting messages to avoid starvation
				for object := range refreshSet {
					freed++
					freeObject(object)
				}
			}
		}
	}

	if c.painter != nil {
		cache.RangeExpiredTexturesFor(c.impl, c.painter.Free)
	}
	return
}

// Initialize initializes the canvas.
func (c *Canvas) Initialize(impl SizeableCanvas, onOverlayChanged func()) {
	c.impl = impl
	c.refreshQueue = async.NewCanvasObjectQueue()
	c.overlays = &overlayStack{
		OverlayStack: internal.OverlayStack{
			OnChange: onOverlayChanged,
			Canvas:   impl,
		},
	}
}

// ObjectTrees return canvas object trees.
//
// This function uses lock.
func (c *Canvas) ObjectTrees() []fyne.CanvasObject {
	c.RLock()
	var content, menu fyne.CanvasObject
	if c.contentTree != nil && c.contentTree.root != nil {
		content = c.contentTree.root.obj
	}
	if c.menuTree != nil && c.menuTree.root != nil {
		menu = c.menuTree.root.obj
	}
	c.RUnlock()
	trees := make([]fyne.CanvasObject, 0, len(c.Overlays().List())+2)
	trees = append(trees, content)
	if menu != nil {
		trees = append(trees, menu)
	}
	trees = append(trees, c.Overlays().List()...)
	return trees
}

// Overlays returns the overlay stack.
func (c *Canvas) Overlays() fyne.OverlayStack {
	// we don't need to lock here, because overlays never changes
	return c.overlays
}

// Painter returns the canvas painter.
func (c *Canvas) Painter() gl.Painter {
	return c.painter
}

// Refresh refreshes a canvas object.
func (c *Canvas) Refresh(obj fyne.CanvasObject) {
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

// RemoveShortcut removes a shortcut from the canvas.
func (c *Canvas) RemoveShortcut(shortcut fyne.Shortcut) {
	c.shortcut.RemoveShortcut(shortcut)
}

// SetContentTreeAndFocusMgr sets content tree and focus manager.
//
// This function does not use the canvas lock.
func (c *Canvas) SetContentTreeAndFocusMgr(content fyne.CanvasObject) {
	c.contentTree = &renderCacheTree{root: &RenderCacheNode{obj: content}}
	var focused fyne.Focusable
	if c.contentFocusMgr != nil {
		focused = c.contentFocusMgr.Focused() // keep old focus if possible
	}
	c.contentFocusMgr = app.NewFocusManager(content)
	if focused != nil {
		c.contentFocusMgr.Focus(focused)
	}
}

// CheckDirtyAndClear returns true if the canvas is dirty and
// clears the dirty state atomically.
func (c *Canvas) CheckDirtyAndClear() bool {
	return atomic.SwapUint32(&c.dirty, 0) != 0
}

// SetDirty sets canvas dirty flag atomically.
func (c *Canvas) SetDirty() {
	atomic.AddUint32(&c.dirty, 1)
}

// SetMenuTreeAndFocusMgr sets menu tree and focus manager.
//
// This function does not use the canvas lock.
func (c *Canvas) SetMenuTreeAndFocusMgr(menu fyne.CanvasObject) {
	c.menuTree = &renderCacheTree{root: &RenderCacheNode{obj: menu}}
	if menu != nil {
		c.menuFocusMgr = app.NewFocusManager(menu)
	} else {
		c.menuFocusMgr = nil
	}
}

// SetMobileWindowHeadTree sets window head tree.
//
// This function does not use the canvas lock.
func (c *Canvas) SetMobileWindowHeadTree(head fyne.CanvasObject) {
	c.mWindowHeadTree = &renderCacheTree{root: &RenderCacheNode{obj: head}}
}

// SetPainter sets the canvas painter.
func (c *Canvas) SetPainter(p gl.Painter) {
	c.painter = p
}

// TypedShortcut handle the registered shortcut.
func (c *Canvas) TypedShortcut(shortcut fyne.Shortcut) {
	c.shortcut.TypedShortcut(shortcut)
}

// Unfocus unfocuses all the objects in the canvas.
func (c *Canvas) Unfocus() {
	mgr := c.focusManager()
	if mgr == nil {
		return
	}
	if mgr.Focus(nil) && c.OnUnfocus != nil {
		c.OnUnfocus()
	}
}

// WalkTrees walks over the trees.
func (c *Canvas) WalkTrees(
	beforeChildren func(*RenderCacheNode, fyne.Position),
	afterChildren func(*RenderCacheNode, fyne.Position),
) {
	c.walkTree(c.contentTree, beforeChildren, afterChildren)
	if c.mWindowHeadTree != nil && c.mWindowHeadTree.root.obj != nil {
		c.walkTree(c.mWindowHeadTree, beforeChildren, afterChildren)
	}
	if c.menuTree != nil && c.menuTree.root.obj != nil {
		c.walkTree(c.menuTree, beforeChildren, afterChildren)
	}
	for _, tree := range c.overlays.renderCaches {
		if tree != nil {
			c.walkTree(tree, beforeChildren, afterChildren)
		}
	}
}

func (c *Canvas) focusManager() *app.FocusManager {
	if focusMgr := c.overlays.TopFocusManager(); focusMgr != nil {
		return focusMgr
	}
	c.RLock()
	defer c.RUnlock()
	if c.isMenuActive() {
		return c.menuFocusMgr
	}
	return c.contentFocusMgr
}

func (c *Canvas) isMenuActive() bool {
	if c.menuTree == nil || c.menuTree.root == nil || c.menuTree.root.obj == nil {
		return false
	}
	menu := c.menuTree.root.obj
	if am, ok := menu.(activatableMenu); ok {
		return am.IsActive()
	}
	return true
}

func (c *Canvas) walkTree(
	tree *renderCacheTree,
	beforeChildren func(*RenderCacheNode, fyne.Position),
	afterChildren func(*RenderCacheNode, fyne.Position),
) {
	tree.Lock()
	defer tree.Unlock()
	var node, parent, prev *RenderCacheNode
	node = tree.root

	bc := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if node != nil && node.obj != obj {
			if parent.firstChild == node {
				parent.firstChild = nil
			}
			node = nil
		}
		if node == nil {
			node = &RenderCacheNode{parent: parent, obj: obj}
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
	ac := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.CanvasObject) {
		node = parent
		parent = node.parent
		if prev != nil && prev.parent != parent {
			prev.nextSibling = nil
		}

		if afterChildren != nil {
			afterChildren(node, pos)
		}

		prev = node
		node = node.nextSibling
	}
	driver.WalkVisibleObjectTree(tree.root.obj, bc, ac)
}

// RenderCacheNode represents a node in a render cache tree.
type RenderCacheNode struct {
	// structural data
	firstChild  *RenderCacheNode
	nextSibling *RenderCacheNode
	obj         fyne.CanvasObject
	parent      *RenderCacheNode
	// cache data
	minSize fyne.Size
	// painterData is some data from the painter associated with the drawed node
	// it may for instance point to a GL texture
	// it should free all associated resources when released
	// i.e. it should not simply be a texture reference integer
	painterData interface{}
}

// Obj returns the node object.
func (r *RenderCacheNode) Obj() fyne.CanvasObject {
	return r.obj
}

type activatableMenu interface {
	IsActive() bool
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
	o.renderCaches = append(o.renderCaches, &renderCacheTree{root: &RenderCacheNode{obj: overlay}})
	o.OverlayStack.Add(overlay)
}

func (o *overlayStack) remove(overlay fyne.CanvasObject) {
	o.OverlayStack.Remove(overlay)
	overlayCount := len(o.List())
	o.renderCaches[overlayCount] = nil // release memory reference to removed element
	o.renderCaches = o.renderCaches[:overlayCount]
}

type renderCacheTree struct {
	sync.RWMutex
	root *RenderCacheNode
}

func (c *Canvas) updateLayout(objToLayout fyne.CanvasObject) {
	switch cont := objToLayout.(type) {
	case *fyne.Container:
		if cont.Layout != nil {
			layout := cont.Layout
			objects := cont.Objects
			c.RUnlock()
			layout.Layout(objects, cont.Size())
			c.RLock()
		}
	case fyne.Widget:
		renderer := cache.Renderer(cont)
		c.RUnlock()
		renderer.Layout(cont.Size())
		c.RLock()
	}
}
