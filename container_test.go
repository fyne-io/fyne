package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, pos, box.Position())

	pos = NewPos(10, 10)
	container.Move(pos)
	assert.Equal(t, pos, container.Position())
	assert.Equal(t, NewPos(0, 0), box.Position())

	box.Move(pos)
	assert.Equal(t, pos, box.Position())
}

func TestNilLayout(t *testing.T) {
	box := new(dummyObject)
	boxSize := box.size
	container := NewContainer(box)

	size := NewSize(100, 100)
	container.Resize(size)
	assert.Equal(t, size, container.Size())
	assert.Equal(t, boxSize, box.Size())

	container.AddObject(box)
	assert.Equal(t, boxSize, box.Size())
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
	assert.Equal(t, size, container.MinSize())
	assert.Equal(t, size, container.Size())
	assert.Equal(t, size, box.Size())

	container.AddObject(box)
	assert.Equal(t, size, box.Size())
}

func TestContainer_Hide(t *testing.T) {
	box := new(dummyObject)
	container := NewContainer(box)

	assert.True(t, container.Visible())
	assert.True(t, box.Visible())
	container.Hide()
	assert.False(t, container.Visible())
	assert.True(t, box.Visible())
}

func TestContainer_Show(t *testing.T) {
	box := new(dummyObject)
	container := NewContainer(box)

	container.Hide()
	assert.True(t, box.Visible())
	assert.False(t, container.Visible())

	container.Show()
	assert.True(t, box.Visible())
	assert.True(t, container.Visible())
}

type dummyObject struct {
	size   Size
	pos    Position
	hidden bool
}

func (d *dummyObject) Size() Size {
	return d.size
}

func (d *dummyObject) Resize(size Size) {
	d.size = size
}

func (d *dummyObject) Position() Position {
	return d.pos
}

func (d *dummyObject) Move(pos Position) {
	d.pos = pos
}

func (d *dummyObject) MinSize() Size {
	return NewSize(5, 5)
}

func (d *dummyObject) Visible() bool {
	return !d.hidden
}

func (d *dummyObject) Show() {
	d.hidden = false
}

func (d *dummyObject) Hide() {
	d.hidden = true
}
