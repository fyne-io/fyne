package fyne

import (
	"sync"
)

// Declare conformity to CanvasObject
var _ CanvasObject = (*Container)(nil)

// Container is a CanvasObject that contains a collection of child objects.
// The layout of the children is set by the specified Layout.
type Container struct {
	size     Size     // The current size of the Container
	position Position // The current position of the Container
	Hidden   bool     // Is this Container hidden

	Layout  Layout // The Layout algorithm for arranging child CanvasObjects
	lock    sync.Mutex
	Objects []CanvasObject // The set of CanvasObjects this container holds
}

// NewContainer returns a new Container instance holding the specified CanvasObjects.
//
// Deprecated: Use container.NewWithoutLayout() to create a container that uses manual layout.
func NewContainer(objects ...CanvasObject) *Container {
	return NewContainerWithoutLayout(objects...)
}

// NewContainerWithoutLayout returns a new Container instance holding the specified
// CanvasObjects that are manually arranged.
//
// Deprecated: Use container.NewWithoutLayout() instead
func NewContainerWithoutLayout(objects ...CanvasObject) *Container {
	ret := &Container{
		Objects: objects,
	}

	ret.size = ret.MinSize()
	return ret
}

// NewContainerWithLayout returns a new Container instance holding the specified
// CanvasObjects which will be laid out according to the specified Layout.
//
// Deprecated: Use container.New() instead
func NewContainerWithLayout(layout Layout, objects ...CanvasObject) *Container {
	ret := &Container{
		Objects: objects,
		Layout:  layout,
	}

	ret.size = layout.MinSize(objects)
	if layout != nil {
		layout.Layout(objects, ret.size)
	}
	return ret
}

// Add appends the specified object to the items this container manages.
//
// Since: 1.4
func (c *Container) Add(add CanvasObject) {
	if add == nil {
		return
	}

	c.lock.Lock()
	c.Objects = append(c.Objects, add)
	layout := c.Layout
	size := c.size
	var obj []CanvasObject
	if layout != nil {
		copyPtr := containerSlicePool.CopyOf(c.Objects)
		defer containerSlicePool.Put(copyPtr)
		obj = *copyPtr
	}
	c.lock.Unlock()

	if layout != nil {
		layout.Layout(obj, size)
	}
}

// AddObject adds another CanvasObject to the set this Container holds.
//
// Deprecated: Use replacement Add() function
func (c *Container) AddObject(o CanvasObject) {
	c.Add(o)
}

// Hide sets this container, and all its children, to be not visible.
func (c *Container) Hide() {
	c.lock.Lock()
	if c.Hidden {
		return
	}

	c.Hidden = true
	c.lock.Unlock()
	repaint(c)
}

// MinSize calculates the minimum size of a Container.
// This is delegated to the Layout, if specified, otherwise it will mimic MaxLayout.
func (c *Container) MinSize() Size {
	c.lock.Lock()
	copyPtr := containerSlicePool.CopyOf(c.Objects)
	defer containerSlicePool.Put(copyPtr)
	obj := *copyPtr
	layout := c.Layout
	c.lock.Unlock()

	if layout != nil {
		return c.Layout.MinSize(obj)
	}

	minSize := NewSize(1, 1)
	for _, child := range obj {
		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

// Move the container (and all its children) to a new position, relative to its parent.
func (c *Container) Move(pos Position) {
	c.lock.Lock()
	c.position = pos
	c.lock.Unlock()
	repaint(c)
}

// Position gets the current position of this Container, relative to its parent.
func (c *Container) Position() Position {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.position
}

// Refresh causes this object to be redrawn in it's current state
func (c *Container) Refresh() {
	c.lock.Lock()
	copyPtr := containerSlicePool.CopyOf(c.Objects)
	defer containerSlicePool.Put(copyPtr)
	obj := *copyPtr
	layout := c.Layout
	size := c.size
	c.lock.Unlock()

	if layout != nil {
		layout.Layout(obj, size)
	}

	for _, child := range obj {
		child.Refresh()
	}

	// this is basically just canvas.Refresh(c) without the package loop
	o := CurrentApp().Driver().CanvasForObject(c)
	if o == nil {
		return
	}
	o.Refresh(c)
}

// Remove updates the contents of this container to no longer include the specified object.
// This method is not intended to be used inside a loop, to remove all the elements.
// It is much more efficient to call RemoveAll() instead.
func (c *Container) Remove(rem CanvasObject) {
	c.lock.Lock()

	if len(c.Objects) == 0 {
		c.lock.Unlock()
		return
	}

	for i, o := range c.Objects {
		if o != rem {
			continue
		}

		removed := make([]CanvasObject, len(c.Objects)-1)
		copy(removed, c.Objects[:i])
		copy(removed[i:], c.Objects[i+1:])

		c.Objects = removed
		layout := c.Layout
		size := c.size
		var obj []CanvasObject
		if layout != nil {
			copyPtr := containerSlicePool.CopyOf(c.Objects)
			defer containerSlicePool.Put(copyPtr)
			obj = *copyPtr
		}
		c.lock.Unlock()

		if layout != nil {
			layout.Layout(obj, size)
		}
		return
	}
}

// RemoveAll updates the contents of this container to no longer include any objects.
//
// Since: 2.2
func (c *Container) RemoveAll() {
	c.lock.Lock()
	c.Objects = nil
	layout := c.Layout
	size := c.size
	c.lock.Unlock()

	if layout != nil {
		layout.Layout(nil, size)
	}
}

// Resize sets a new size for the Container.
func (c *Container) Resize(size Size) {
	c.lock.Lock()

	if c.size == size {
		c.lock.Unlock()
		return
	}

	var obj []CanvasObject
	c.size = size
	layout := c.Layout
	if layout != nil {
		copyPtr := containerSlicePool.CopyOf(c.Objects)
		defer containerSlicePool.Put(copyPtr)
		obj = *copyPtr
	}
	c.lock.Unlock()

	if layout != nil {
		layout.Layout(obj, size)
	}
}

// Show sets this container, and all its children, to be visible.
func (c *Container) Show() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.Hidden {
		return
	}

	c.Hidden = false
}

// Size returns the current size of this container.
func (c *Container) Size() Size {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.size
}

// Visible returns true if the container is currently visible, false otherwise.
func (c *Container) Visible() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return !c.Hidden
}

// repaint instructs the containing canvas to redraw, even if nothing changed.
// This method is a duplicate of what is in `canvas/canvas.go` to avoid a dependency loop or public API.
func repaint(obj *Container) {
	app := CurrentApp()
	if app == nil || app.Driver() == nil {
		return
	}

	c := app.Driver().CanvasForObject(obj)
	if c != nil {
		if paint, ok := c.(interface{ SetDirty() }); ok {
			paint.SetDirty()
		}
	}
}

var containerSlicePool canvasObjectSlicePool

// A tiered pool of []CanvasObject slices
// see https://github.com/golang/go/issues/27735 for why a single sync.Pool
// is a potential memory leak, as well as discussions for other alternatives
type canvasObjectSlicePool struct {
	xsPool sync.Pool // holds slices with length [xsPooledSliceMinSize, sPooledSliceMinSize)
	sPool  sync.Pool // etc
	mPool  sync.Pool
	lPool  sync.Pool
	xlPool sync.Pool // holds slices with length [xlPooledSliceMinSize, xlPooledSliceMaxSize]
}

const (
	xsPooledSliceMinSize = 8
	sPooledSliceMinSize  = 16
	mPooledSliceMinSize  = 32
	lPooledSliceMinSize  = 64
	xlPooledSliceMinSize = 256
	xlPooledSliceMaxSize = 2048
	// do not pool slices of more than xlPooledSliceMaxSize
)

// Get gets a pointer to a slice of at least the given capacity from the pool,
// allocating a new slice if needed. The slice may have any length, but it is
// always safe to re-slice to a new length that is <= capacity.
func (s *canvasObjectSlicePool) Get(capacity int) *[]CanvasObject {
	var obj any
	switch {
	case capacity <= xsPooledSliceMinSize:
		obj = s.xsPool.Get()
	case capacity <= sPooledSliceMinSize:
		obj = s.sPool.Get()
	case capacity <= mPooledSliceMinSize:
		obj = s.mPool.Get()
	case capacity <= lPooledSliceMinSize:
		obj = s.lPool.Get()
	default:
		obj = s.xlPool.Get()
	}

	if obj == nil {
		if capacity < xsPooledSliceMinSize {
			capacity = xsPooledSliceMinSize
		}
		tmp := make([]CanvasObject, 0, capacity)
		obj = &tmp
	}
	slice := obj.(*[]CanvasObject)
	if cap(*slice) < capacity {
		// can only happen in xl pool - reallocate a longer slice
		tmp := make([]CanvasObject, 0, capacity)
		slice = &tmp
	}
	return slice
}

// Put releases the given slice pointer into the pool.
func (s *canvasObjectSlicePool) Put(slice *[]CanvasObject) {
	*slice = (*slice)[:0]
	switch c := cap(*slice); {
	case c < xsPooledSliceMinSize:
		return // don't pool a slice too small
	case c < sPooledSliceMinSize:
		s.xsPool.Put(slice)
	case c < mPooledSliceMinSize:
		s.sPool.Put(slice)
	case c < lPooledSliceMinSize:
		s.mPool.Put(slice)
	case c < xlPooledSliceMinSize:
		s.lPool.Put(slice)
	case c <= xlPooledSliceMaxSize:
		s.xlPool.Put(slice)
	default:
		return // don't pool a slice too large
	}
}

// CopyOf is a convenience method that gets a slice from the pool
// of capacity at least len(slice), and copies the contents of slice
// into the new slice, before returning its pointer.
func (s *canvasObjectSlicePool) CopyOf(slice []CanvasObject) *[]CanvasObject {
	l := len(slice)
	ptr := s.Get(l)
	cp := *ptr
	cp = cp[:l]
	copy(cp, slice)
	*ptr = cp
	return ptr
}
