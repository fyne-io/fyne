package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestMultipleWindows_Add(t *testing.T) {
	m := NewMultipleWindows()
	assert.Zero(t, len(m.Windows))

	m.Add(NewInnerWindow("1", widget.NewLabel("Inside")))
	assert.Equal(t, 1, len(m.Windows))
}

func TestMultipleWindows_Drag(t *testing.T) {
	w := NewInnerWindow("1", widget.NewLabel("Inside"))
	m := NewMultipleWindows(w)
	m.Resize(fyne.NewSize(30, 30))
	_ = test.TempWidgetRenderer(t, m) // initialise display
	assert.Equal(t, 1, len(m.Windows))

	assert.True(t, w.Position().IsZero())
	w.OnDragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 10, DY: 5}})
	assert.Equal(t, float32(10), w.Position().X)
	assert.Equal(t, float32(5), w.Position().Y)

	drag := &fyne.DragEvent{Dragged: fyne.Delta{DX: -10, DY: -5}}
	drag.Position = fyne.NewPos(5, 5)
	drag.AbsolutePosition = fyne.NewPos(5, 5)
	w.OnDragged(drag)
	assert.Equal(t, float32(0), w.Position().X)
	assert.Equal(t, float32(0), w.Position().Y)
}

func TestMultipleWindows_RaiseToTop(t *testing.T) {
	w1 := NewInnerWindow("1", widget.NewLabel("Content"))
	m := NewMultipleWindows(w1)
	assert.Equal(t, w1, m.Top())

	w2 := NewInnerWindow("2", widget.NewLabel("Content"))
	m.Add(w2)
	assert.Equal(t, w2, m.Top())

	m.RaiseToTop(w1)
	assert.Equal(t, w1, m.Top())
}

func TestMultipleWindows_Top(t *testing.T) {
	m := NewMultipleWindows()
	assert.Nil(t, m.Top())

	w1 := NewInnerWindow("1", widget.NewLabel("Content"))
	m.Add(w1)
	assert.Equal(t, w1, m.Top())
}

func TestMultipleWindows_OnTappedBarPreserved(t *testing.T) {
	w1 := NewInnerWindow("1", widget.NewLabel("Content"))
	called := 0
	w1.OnTappedBar = func() { called++ }

	w2 := NewInnerWindow("2", widget.NewLabel("Content"))
	m := NewMultipleWindows(w1, w2)
	_ = test.TempWidgetRenderer(t, m) // triggers the first refreshChildren/setupChild

	assert.Equal(t, w2, m.Top())

	w1.OnTappedBar()
	assert.Equal(t, 1, called, "user's OnTappedBar set before adding the window should still be called")
	assert.Equal(t, w1, m.Top(), "the built-in raise-to-top behaviour must still run")
}

func TestMultipleWindows_OnTappedBarNotDoubled(t *testing.T) {
	w1 := NewInnerWindow("1", widget.NewLabel("Content"))
	called := 0
	w1.OnTappedBar = func() { called++ }

	m := NewMultipleWindows(w1)
	_ = test.TempWidgetRenderer(t, m)

	m.Refresh()
	m.Refresh()
	m.Add(NewInnerWindow("2", widget.NewLabel("Content")))

	w1.OnTappedBar()
	assert.Equal(t, 1, called, "repeated Refresh()/Add() calls must not wrap the callback more than once")
}
