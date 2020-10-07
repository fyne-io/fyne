package test_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestAssertCanvasTappableAt(t *testing.T) {
	c := test.NewCanvas()
	b := widget.NewButton("foo", nil)
	c.SetContent(b)
	c.Resize(fyne.NewSize(300, 300))
	b.Resize(fyne.NewSize(100, 100))
	b.Move(fyne.NewPos(100, 100))

	tt := &testing.T{}
	assert.True(t, test.AssertCanvasTappableAt(tt, c, fyne.NewPos(101, 101)), "tappable found")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertCanvasTappableAt(tt, c, fyne.NewPos(99, 99)), "tappable not found")
	assert.True(t, tt.Failed(), "test failed")
}

func TestDrag(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	d := &draggable{}
	c.SetContent(fyne.NewContainer(d))
	c.Resize(fyne.NewSize(30, 30))
	d.Resize(fyne.NewSize(20, 20))
	d.Move(fyne.NewPos(10, 10))

	test.Drag(c, fyne.NewPos(5, 5), 10, 10)
	assert.Nil(t, d.event, "nothing happens if no draggable was found at position")
	assert.False(t, d.wasDragged)

	test.Drag(c, fyne.NewPos(15, 15), 17, 42)
	assert.Equal(t, &fyne.DragEvent{
		PointEvent: fyne.PointEvent{Position: fyne.Position{X: 5, Y: 5}},
		DraggedX:   17,
		DraggedY:   42,
	}, d.event)
	assert.True(t, d.wasDragged)
}

func TestScroll(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	s := &scrollable{}
	c.SetContent(fyne.NewContainer(s))
	c.Resize(fyne.NewSize(30, 30))
	s.Resize(fyne.NewSize(20, 20))
	s.Move(fyne.NewPos(10, 10))

	test.Scroll(c, fyne.NewPos(5, 5), 10, 10)
	assert.Nil(t, s.event, "nothing happens if no scrollable was found at position")

	test.Scroll(c, fyne.NewPos(15, 15), 17, 42)
	assert.Equal(t, &fyne.ScrollEvent{DeltaX: 17, DeltaY: 42}, s.event)
}

type draggable struct {
	widget.BaseWidget
	event      *fyne.DragEvent
	wasDragged bool
}

var _ fyne.Draggable = (*draggable)(nil)

func (d *draggable) Dragged(event *fyne.DragEvent) {
	d.event = event
}

func (d *draggable) DragEnd() {
	d.wasDragged = true
}

type scrollable struct {
	widget.BaseWidget
	event *fyne.ScrollEvent
}

var _ fyne.Scrollable = (*scrollable)(nil)

func (s *scrollable) Scrolled(event *fyne.ScrollEvent) {
	s.event = event
}
