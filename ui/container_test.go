package ui

import "testing"

import "github.com/stretchr/testify/assert"

func TestDefaultMinSize(t *testing.T) {
	text := new(dummyObject)
	minSize := text.MinSize()

	container := NewContainer(text)
	container.AddObject(text)
	layoutMin := container.MinSize()

	assert.Equal(t, minSize, layoutMin)
}

type dummyObject struct {
}

func (d *dummyObject) CurrentSize() Size {
	return NewSize(10, 10)
}

func (d *dummyObject) Resize(Size) {
}

func (d *dummyObject) CurrentPosition() Position {
	return NewPos(10, 10)
}

func (d *dummyObject) Move(Position) {
}

func (d *dummyObject) MinSize() Size {
	return NewSize(5, 5)
}
