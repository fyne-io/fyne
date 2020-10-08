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
	c.SetContent(fyne.NewContainerWithoutLayout(d))
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

func TestFocusNext(t *testing.T) {
	c := test.NewCanvas()
	f1 := &focusable{}
	f2 := &focusable{}
	f3 := &focusable{}
	c.SetContent(fyne.NewContainerWithoutLayout(f1, f2, f3))

	assert.Nil(t, c.Focused())
	assert.False(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.False(t, f3.Focused())

	test.FocusNext(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.False(t, f3.Focused())

	test.FocusNext(c)
	assert.Equal(t, f2, c.Focused())
	assert.False(t, f1.Focused())
	assert.True(t, f2.Focused())
	assert.False(t, f3.Focused())

	test.FocusNext(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.True(t, f3.Focused())

	test.FocusNext(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.False(t, f3.Focused())
}

func TestFocusPrevious(t *testing.T) {
	c := test.NewCanvas()
	f1 := &focusable{}
	f2 := &focusable{}
	f3 := &focusable{}
	c.SetContent(fyne.NewContainerWithoutLayout(f1, f2, f3))

	assert.Nil(t, c.Focused())
	assert.False(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.False(t, f3.Focused())

	test.FocusPrevious(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.True(t, f3.Focused())

	test.FocusPrevious(c)
	assert.Equal(t, f2, c.Focused())
	assert.False(t, f1.Focused())
	assert.True(t, f2.Focused())
	assert.False(t, f3.Focused())

	test.FocusPrevious(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.False(t, f3.Focused())

	test.FocusPrevious(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.Focused())
	assert.False(t, f2.Focused())
	assert.True(t, f3.Focused())
}

func TestScroll(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	s := &scrollable{}
	c.SetContent(fyne.NewContainerWithoutLayout(s))
	c.Resize(fyne.NewSize(30, 30))
	s.Resize(fyne.NewSize(20, 20))
	s.Move(fyne.NewPos(10, 10))

	test.Scroll(c, fyne.NewPos(5, 5), 10, 10)
	assert.Nil(t, s.event, "nothing happens if no scrollable was found at position")

	test.Scroll(c, fyne.NewPos(15, 15), 17, 42)
	assert.Equal(t, &fyne.ScrollEvent{DeltaX: 17, DeltaY: 42}, s.event)
}

var _ fyne.Draggable = (*draggable)(nil)

type draggable struct {
	widget.BaseWidget
	event      *fyne.DragEvent
	wasDragged bool
}

func (d *draggable) DragEnd() {
	d.wasDragged = true
}

func (d *draggable) Dragged(event *fyne.DragEvent) {
	d.event = event
}

var _ fyne.Focusable = (*focusable)(nil)

type focusable struct {
	widget.BaseWidget
	focused bool
}

func (f *focusable) FocusGained() {
	f.focused = true
}

func (f *focusable) FocusLost() {
	f.focused = false
}

func (f *focusable) Focused() bool {
	return f.focused
}

func (f *focusable) TypedKey(event *fyne.KeyEvent) {
}

func (f *focusable) TypedRune(r rune) {
}

var _ fyne.Scrollable = (*scrollable)(nil)

type scrollable struct {
	widget.BaseWidget
	event *fyne.ScrollEvent
}

func (s *scrollable) Scrolled(event *fyne.ScrollEvent) {
	s.event = event
}
