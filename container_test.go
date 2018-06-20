package fyne

import "testing"

import "github.com/stretchr/testify/assert"

func TestDefaultMinSize(t *testing.T) {
	box := new(dummyObject)
	minSize := box.MinSize()

	container := NewContainer(box)
	assert.Equal(t, minSize, container.MinSize())

	container.AddObject(box)
	assert.Equal(t, minSize, container.MinSize())
}

func TestDefaultLayout(t *testing.T) {
	box := new(dummyObject)
	container := NewContainer(box)

	size := NewSize(100, 100)
	container.Resize(size)
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
