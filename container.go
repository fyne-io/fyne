package fyne

// Container is a CanvasObject that contains a collection of child objects.
// The layout of the children is set by the specified Layout.
type Container struct {
	size     Size     // The current size of the Container
	position Position // The current position of the Container
	Hidden   bool     // Is this Container hidden

	Layout  Layout         // The Layout algorithm for arranging child CanvasObjects
	Objects []CanvasObject // The set of CanvasObjects this container holds
}

func (c *Container) layout() {
	if c.Layout != nil {
		c.Layout.Layout(c.Objects, c.size)
		return
	}

	for _, child := range c.Objects {
		child.Resize(c.size)
		child.Move(c.position)
	}
}

// Size returns the current size of this container.
func (c *Container) Size() Size {
	return c.size
}

// Resize sets a new size for the Container.
func (c *Container) Resize(size Size) {
	c.size = size
	c.layout()
}

// Position gets the current position of this Container, relative to its parent.
func (c *Container) Position() Position {
	return c.position
}

// Move the container (and all its children) to a new position, relative to its parent.
func (c *Container) Move(pos Position) {
	c.position = pos
	c.layout()
}

// MinSize calculates the minimum size of a Container.
// This is delegated to the Layout, if specified, otherwise it will mimic MaxLayout.
func (c *Container) MinSize() Size {
	if c.Layout != nil {
		return c.Layout.MinSize(c.Objects)
	}

	minSize := NewSize(1, 1)
	for _, child := range c.Objects {
		minSize = minSize.Union(child.MinSize())
	}

	return minSize
}

// Visible returns true if the container is currently visible, false otherwise.
func (c *Container) Visible() bool {
	return !c.Hidden
}

// Show sets this container, and all its children, to be visible.
func (c *Container) Show() {
	c.Hidden = false
	for _, child := range c.Objects {
		child.Show()
	}
}

// Hide sets this container, and all its children, to be not visible.
func (c *Container) Hide() {
	c.Hidden = true
	for _, child := range c.Objects {
		child.Hide()
	}
}

// AddObject adds another CanvasObject to the set this Container holds.
func (c *Container) AddObject(o CanvasObject) {
	c.Objects = append(c.Objects, o)
	c.layout()
}

// NewContainer returns a new Container instance holding the specified CanvasObjects.
func NewContainer(objects ...CanvasObject) *Container {
	ret := &Container{
		Objects: objects,
	}

	ret.size = ret.MinSize()
	ret.layout()

	return ret
}

// NewContainerWithLayout returns a new Container instance holding the specified
// CanvasObjects which will be laid out according to the specified Layout.
func NewContainerWithLayout(layout Layout, objects ...CanvasObject) *Container {
	ret := &Container{
		Objects: objects,
		Layout:  layout,
	}

	ret.size = layout.MinSize(objects)
	ret.layout()
	return ret
}
