package fyne

import "testing"

import "github.com/stretchr/testify/assert"

func TestMinSize(t *testing.T) {
	box := new(dummyObject)
	minSize := box.MinSize()

	container := NewContainer(box)
	assert.Equal(t, minSize, container.MinSize())

	container.AddObject(box)
	assert.Equal(t, minSize, container.MinSize())
}

func TestMove(t *testing.T) {
	box := new(dummyObject)
	container := NewContainer(box)

	size := NewSize(100, 100)
	pos := NewPos(0, 0)
	container.Resize(size)
	assert.Equal(t, pos, box.CurrentPosition())

	pos = NewPos(10, 10)
	container.Move(pos)
	assert.Equal(t, pos, container.CurrentPosition())
	assert.Equal(t, pos, box.CurrentPosition())
}

func TestDefaultLayout(t *testing.T) {
	box := new(dummyObject)
	container := NewContainer(box)

	size := NewSize(100, 100)
	container.Resize(size)
	assert.Equal(t, size, container.CurrentSize())
	assert.Equal(t, size, box.CurrentSize())

	container.AddObject(box)
	assert.Equal(t, size, box.CurrentSize())
}

type customLayout struct {
}

func (c *customLayout) Layout(objs []CanvasObject, size Size) {
	for _, child := range objs {
		child.Resize(size)
	}
}

func (c *customLayout) MinSize(objects []CanvasObject) Size {
	return NewSize(10, 10)
}

func TestCustomLayout(t *testing.T) {
	box := new(dummyObject)
	layout := new(customLayout)
	container := NewContainerWithLayout(layout, box)

	size := layout.MinSize(container.Objects)
	assert.Equal(t, size, container.CurrentSize())
	assert.Equal(t, size, box.CurrentSize())

	container.AddObject(box)
	assert.Equal(t, size, box.CurrentSize())
}

type dummyObject struct {
	size Size
	pos  Position
}

func (d *dummyObject) CurrentSize() Size {
	return d.size
}

func (d *dummyObject) Resize(size Size) {
	d.size = size
}

func (d *dummyObject) CurrentPosition() Position {
	return d.pos
}

func (d *dummyObject) Move(pos Position) {
	d.pos = pos
}

func (d *dummyObject) MinSize() Size {
	return NewSize(5, 5)
}
