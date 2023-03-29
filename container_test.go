package fyne

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer_Add(t *testing.T) {
	box := new(dummyObject)
	container := NewContainerWithLayout(new(customLayout))
	assert.Equal(t, 0, len(container.Objects))
	assert.Equal(t, float32(10), container.MinSize().Width)
	assert.Equal(t, float32(0), container.MinSize().Height)

	container.Add(box)
	assert.Equal(t, 1, len(container.Objects))
	assert.Equal(t, float32(10), container.MinSize().Width)
	assert.Equal(t, float32(10), container.MinSize().Height)

	oldLength := len(container.Objects)
	container.Add(nil)
	assert.Equal(t, oldLength, len(container.Objects))

	box2 := new(dummyObject)
	container.Add(box2)
	assert.Equal(t, 2, len(container.Objects))
	assert.Equal(t, float32(10), container.MinSize().Width)
	assert.Equal(t, float32(20), container.MinSize().Height)
	assert.Equal(t, float32(0), box2.Position().X)
	assert.Equal(t, float32(10), box2.Position().Y)
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
	box1 := new(dummyObject)
	box2 := new(dummyObject)
	container := NewContainerWithLayout(new(customLayout), box1, box2)
	assert.Equal(t, 2, len(container.Objects))
	assert.Equal(t, float32(10), container.MinSize().Width)
	assert.Equal(t, float32(20), container.MinSize().Height)

	container.Remove(box1)
	assert.Equal(t, 1, len(container.Objects))
	assert.Equal(t, float32(10), container.MinSize().Width)
	assert.Equal(t, float32(10), container.MinSize().Height)
	assert.Equal(t, float32(0), box2.Position().X)
	assert.Equal(t, float32(0), box2.Position().Y)
}

func TestContainer_Remove_Race(t *testing.T) {
	var objs []CanvasObject
	for i := 0; i < 100; i++ {
		objs = append(objs, new(dummyObject))
	}

	container := NewContainerWithLayout(new(customLayout), objs...)

	wg := &sync.WaitGroup{}
	wg.Add(100)
	for _, o := range objs {
		rmo := o
		go func() {
			container.Remove(rmo)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, 0, len(container.Objects))
}

func TestContainer_Add_Race(t *testing.T) {
	container := NewContainerWithLayout(new(customLayout))
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			container.Add(new(dummyObject))
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, 100, len(container.Objects))
}

func TestContainer_RemoveAll(t *testing.T) {
	box1 := new(dummyObject)
	box2 := new(dummyObject)
	container := NewContainerWithLayout(new(customLayout), box1, box2)
	assert.Equal(t, 2, len(container.Objects))

	container.RemoveAll()
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
	y := float32(0)
	for _, child := range objs {
		child.Resize(size)
		child.Move(NewPos(0, y))
		y += 10
	}
}

func (c *customLayout) MinSize(objs []CanvasObject) Size {
	return NewSize(10, float32(10*len(objs)))
}

type dummyObject struct {
	size   Size
	pos    Position
	hidden bool
}

func (d *dummyObject) Hide() {
	d.hidden = true
}

func (d *dummyObject) MinSize() Size {
	return NewSize(5, 5)
}

func (d *dummyObject) Move(pos Position) {
	d.pos = pos
}

func (d *dummyObject) Position() Position {
	return d.pos
}

func (d *dummyObject) Refresh() {
}

func (d *dummyObject) Resize(size Size) {
	d.size = size
}

func (d *dummyObject) Show() {
	d.hidden = false
}

func (d *dummyObject) Size() Size {
	return d.size
}

func (d *dummyObject) Visible() bool {
	return !d.hidden
}
