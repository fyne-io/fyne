package ui

// Container is a CanvasObject that contains a collection of child objects.
// The layout of the children is set by the specified Layout.
type Container struct {
	Size     Size     // The current size of the Container
	Position Position // The current position of the Container

	Layout  Layout         // The Layout algorithm for arranging child CanvasObjects
	Objects []CanvasObject // The set of CanvasObjects this container holds
}

// CurrentSize returns the current size of this container
func (c *Container) CurrentSize() Size {
	return c.Size
}

// Resize sets a new size for the Container
func (c *Container) Resize(size Size) {
	c.Size = size
}

// CurrentPosition gets the current position of this Container, relative to it's parent
func (c *Container) CurrentPosition() Position {
	return c.Position
}

// Move the container (and all it's children) to a new position, relative to it's parent
func (c *Container) Move(pos Position) {
	c.Position = pos
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

// AddObject adds another CanvasObject to the set this Container holds
func (c *Container) AddObject(o CanvasObject) {
	c.Objects = append(c.Objects, o)
}

// NewContainer returns a new Container instance holding the specified CanvasObjects
func NewContainer(objects ...CanvasObject) *Container {
	return &Container{
		Objects: objects,
	}
}

// NewContainerWithLayout returns a new Container instance holding the specified
// CanvasObjects which will be laid out according to the specified Layout
func NewContainerWithLayout(layout Layout, objects ...CanvasObject) *Container {
	return &Container{
		Objects: objects,
		Layout:  layout,
	}
}
