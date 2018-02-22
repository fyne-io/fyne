package ui

type Container struct {
	Layout   Layout
	Size     Size
	Position Position

	Objects []CanvasObject
}

func (c *Container) CurrentSize() Size {
	return c.Size
}

func (c *Container) Resize(size Size) {
	c.Size = size
}

func (c *Container) CurrentPosition() Position {
	return c.Position
}

func (c *Container) Move(pos Position) {
	c.Position = pos
}

func (c *Container) MinSize() Size {
	return NewSize(1, 1)
}

func (c *Container) AddObject(o CanvasObject) {
	c.Objects = append(c.Objects, o)
}

func NewContainer(objects ...CanvasObject) *Container {
	return &Container{
		Objects: objects,
	}
}
