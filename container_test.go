package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer_Add(t *testing.T) {
	box := new(dummyObject)
	container := NewContainerWithoutLayout()
	assert.Equal(t, 0, len(container.Objects))

	container.Add(box)
	assert.Equal(t, 1, len(container.Objects))
}

func TestContainer_CustomLayout(t *testing.T) {
	box := new(dummyObject)
	layout := new(customLayout)
	container := NewContainerWithLayout(layout, box)

	size := layout.MinSize(container.Objects)
	assert.Equal(t, size, container.MinSize())
	assert.Equal(t, size, container.Size())
	assert.Equal(t, size, box.Size())
}

func TestContainer_Hide(t *testing.T) {
	box := new(dummyObject)
	container := NewContainerWithoutLayout(box)

	assert.True(t, container.Visible())
	assert.True(t, box.Visible())
	container.Hide()
	assert.False(t, container.Visible())
	assert.True(t, box.Visible())
}

func TestContainer_MinSize(t *testing.T) {
	box := new(dummyObject)
	minSize := box.MinSize()

	container := NewContainerWithoutLayout(box)
	assert.Equal(t, minSize, container.MinSize())
}

func TestContainer_Move(t *testing.T) {
	box := new(dummyObject)
	container := NewContainerWithoutLayout(box)

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

func TestContainer_NilLayout(t *testing.T) {
	box := new(dummyObject)
	boxSize := box.size
	container := NewContainerWithoutLayout(box)

	size := NewSize(100, 100)
	container.Resize(size)
	assert.Equal(t, size, container.Size())
	assert.Equal(t, boxSize, box.Size())
}

func TestContainer_Remove(t *testing.T) {
	box := new(dummyObject)
	container := NewContainerWithoutLayout(box)
	assert.Equal(t, 1, len(container.Objects))

	container.Remove(box)
	assert.Equal(t, 0, len(container.Objects))
}

func TestContainer_Show(t *testing.T) {
	box := new(dummyObject)
	container := NewContainerWithoutLayout(box)

	container.Hide()
	assert.True(t, box.Visible())
	assert.False(t, container.Visible())

	container.Show()
	assert.True(t, box.Visible())
	assert.True(t, container.Visible())
}

type customLayout struct {
}

func (c *customLayout) Layout(objs []CanvasObject, size Size) {
	for _, child := range objs {
		child.Resize(size)
	}
}

func (c *customLayout) MinSize(_ []CanvasObject) Size {
	return NewSize(10, 10)
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

func (d *dummyObject) Refresh() {
}
