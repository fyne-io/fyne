package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestOverlayStack(t *testing.T) {
	s := &internal.OverlayStack{Canvas: test.NewCanvas()}
	o1 := widget.NewLabel("A")
	o2 := widget.NewLabel("B")
	o3 := widget.NewLabel("C")
	o4 := widget.NewLabel("D")
	o5 := widget.NewLabel("E")

	// initial empty
	assert.Empty(t, s.List())
	assert.Nil(t, s.Top())

	// add one & remove
	s.Add(o1)
	assert.Equal(t, []fyne.CanvasObject{o1}, s.List())
	assert.Equal(t, o1, s.Top())
	// remove other does nothing
	s.Remove(o2)
	assert.Equal(t, []fyne.CanvasObject{o1}, s.List())
	assert.Equal(t, o1, s.Top())
	// remove the correct one
	s.Remove(o1)
	assert.Empty(t, s.List())
	assert.Nil(t, s.Top())

	// add multiple & remove
	s.Add(o1)
	s.Add(o2)
	s.Add(o3)
	s.Add(o4)
	s.Add(o5)
	assert.Equal(t, []fyne.CanvasObject{o1, o2, o3, o4, o5}, s.List())
	assert.Equal(t, o5, s.Top())
	s.Remove(o5)
	assert.Equal(t, []fyne.CanvasObject{o1, o2, o3, o4}, s.List())
	assert.Equal(t, o4, s.Top())
	// remove cuts the stack
	s.Remove(o2)
	assert.Equal(t, []fyne.CanvasObject{o1}, s.List())
	assert.Equal(t, o1, s.Top())
	s.Remove(o1)
	assert.Empty(t, s.List())
	assert.Nil(t, s.Top())
}
