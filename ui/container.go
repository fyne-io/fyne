package ui

type Container struct {
	Layout  Layout
	Objects []CanvasObject
}

func (c *Container) MinSize() (int, int) {
	return 1, 1
}

func (c *Container) AddObject(o CanvasObject) {
	c.Objects = append(c.Objects, o)
}
