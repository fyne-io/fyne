package fyne

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
)

// Declare conformity to CanvasObject
var _ CanvasObject = (*Container)(nil)

// Container is a CanvasObject that contains a collection of child objects.
// The layout of the children is set by the specified Layout.
type Container struct {
	id       string         // The Object Id of the Container,
	parent   LinkableObject // The parent object, which contains this container
	size     Size           // The current size of the Container
	position Position       // The current position of the Container
	Hidden   bool           // Is this Container hidden

	Layout     Layout // The Layout algorithm for arranging child CanvasObjects
	lock       sync.Mutex
	Objects    []CanvasObject            // The set of CanvasObjects this container holds
	objectsMap map[string]NameableObject // The map of CanvasObjects this container holds (their Id is the key)
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
	ret.layout()
	return ret
}

// InitObjectMap initializes the internal object map. Usually it should not be called outside of that library
func (c *Container) InitObjectMap() {
	c.objectsMap = make(map[string]NameableObject)
}

// SetId is used to set the object id, then it can be used to retrieve the object from the parent's object map.
func (c *Container) SetId(id string) error {
	if c.id != "" {
		return errors.New("object ID is already set")
	} else {
		c.id = id
		return nil
	}
}

// ID is used to get the object id, then it can be used to retrieve the object from the parent's object map.
func (c *Container) ID() string {
	if c.id == "" {
		c.id = fmt.Sprintf("cntr-%d", rand.Intn(math.MaxInt32))
	}
	return c.id
}

// SetParent is used to set the parent object pointer. Should be used by the object where this widget is added to.
func (c *Container) SetParent(parent LinkableObject) {
	c.parent = parent
}

// Parent is used to get the parent object pointer. Can be used to access the parent object and its object map.
func (c *Container) Parent() LinkableObject {
	return c.parent
}

// LinkedObject is used to get a linked object, registered by the Add method, called by its ID.
func (c *Container) LinkedObject(id string) CanvasObject {
	var co NameableObject = c.objectsMap[id]
	if co == nil {
		return nil
	} else {
		return co.(CanvasObject)
	}
}

// Add appends the specified object to the items this container manages.
//
// Since: 1.4
func (c *Container) Add(add CanvasObject) {
	if add == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.Objects = append(c.Objects, add)
	if nObj, compliant := add.(NameableObject); compliant {
		c.objectsMap[nObj.ID()] = nObj
	}
	if lObj, compliant := add.(LinkableObject); compliant {
		lObj.SetParent(c)
	}
	c.layout()
}

// AddObject adds another CanvasObject to the set this Container holds.
//
// Deprecated: Use replacement Add() function
func (c *Container) AddObject(o CanvasObject) {
	c.Add(o)
}

// Hide sets this container, and all its children, to be not visible.
func (c *Container) Hide() {
	if c.Hidden {
		return
	}

	c.Hidden = true
	repaint(c)
}

// MinSize calculates the minimum size of a Container.
// This is delegated to the Layout, if specified, otherwise it will mimic MaxLayout.
func (c *Container) MinSize() Size {
	if c.Layout != nil {
		return c.Layout.MinSize(c.Objects)
	}

	minSize := NewSize(1, 1)
	for _, child := range c.Objects {
		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

// Move the container (and all its children) to a new position, relative to its parent.
func (c *Container) Move(pos Position) {
	c.position = pos
	repaint(c)
}

// Position gets the current position of this Container, relative to its parent.
func (c *Container) Position() Position {
	return c.position
}

// Refresh causes this object to be redrawn in its current state
func (c *Container) Refresh() {
	c.layout()

	for _, child := range c.Objects {
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
	defer c.lock.Unlock()
	if len(c.Objects) == 0 {
		return
	}

	for i, o := range c.Objects {
		if o != rem {
			continue
		}

		if wid, compliant := rem.(NameableObject); compliant {
			delete(c.objectsMap, wid.ID())
		}

		if wid, compliant := rem.(LinkableObject); compliant {
			wid.SetParent(nil)
		}

		removed := make([]CanvasObject, len(c.Objects)-1)
		copy(removed, c.Objects[:i])
		copy(removed[i:], c.Objects[i+1:])

		c.Objects = removed
		c.layout()
		return
	}
}

// RemoveAll updates the contents of this container to no longer include any objects.
//
// Since: 2.2
func (c *Container) RemoveAll() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Objects = nil
	clear(c.objectsMap)
	c.layout()
}

// Resize sets a new size for the Container.
func (c *Container) Resize(size Size) {
	if c.size == size {
		return
	}

	c.size = size
	c.layout()
}

// Show sets this container, and all its children, to be visible.
func (c *Container) Show() {
	if !c.Hidden {
		return
	}

	c.Hidden = false
}

// Size returns the current size of this container.
func (c *Container) Size() Size {
	return c.size
}

// Visible returns true if the container is currently visible, false otherwise.
func (c *Container) Visible() bool {
	return !c.Hidden
}

func (c *Container) layout() {
	if c.Layout == nil {
		return
	}

	c.Layout.Layout(c.Objects, c.size)
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
