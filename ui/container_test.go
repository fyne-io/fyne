package ui

import "testing"

import "reflect"

func TestDefaultMinSize(t *testing.T) {
	text := new(dummyObject)
	minSize := text.MinSize()

	container := NewContainer(text)
	layoutMin := container.MinSize()

	if !reflect.DeepEqual(minSize, layoutMin) {
		t.Fatal("Expected", minSize, "but got", layoutMin)
	}
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
