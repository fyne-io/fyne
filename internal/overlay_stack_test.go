package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
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
	assert.Nil(t, s.TopFocusManager())
	assert.Empty(t, s.ListFocusManagers())

	// add one & remove
	s.Add(o1)
	assert.Equal(t, []fyne.CanvasObject{o1}, s.List())
	assert.Equal(t, o1, s.Top())
	fm := s.TopFocusManager()
	assert.NotNil(t, fm)
	assert.Equal(t, []*app.FocusManager{fm}, s.ListFocusManagers())
	// remove other does nothing
	s.Remove(o2)
	assert.Equal(t, []fyne.CanvasObject{o1}, s.List())
	assert.Equal(t, o1, s.Top())
	assert.Equal(t, fm, s.TopFocusManager())
	assert.Equal(t, []*app.FocusManager{fm}, s.ListFocusManagers())
	// remove the correct one
	s.Remove(o1)
	assert.Empty(t, s.List())
	assert.Nil(t, s.Top())
	assert.Nil(t, s.TopFocusManager())
	assert.Empty(t, s.ListFocusManagers())

	// add multiple & remove
	s.Add(o1)
	fm1 := s.TopFocusManager()
	assert.NotNil(t, fm1)
	assert.Equal(t, []*app.FocusManager{fm1}, s.ListFocusManagers())
	s.Add(o2)
	fm2 := s.TopFocusManager()
	assert.NotNil(t, fm2)
	assert.NotEqual(t, fm1, fm2)
	assert.Equal(t, []*app.FocusManager{fm1, fm2}, s.ListFocusManagers())
	s.Add(o3)
	fm3 := s.TopFocusManager()
	assert.NotNil(t, fm3)
	assert.NotEqual(t, fm2, fm3)
	assert.Equal(t, []*app.FocusManager{fm1, fm2, fm3}, s.ListFocusManagers())
	s.Add(o4)
	fm4 := s.TopFocusManager()
	assert.NotNil(t, fm4)
	assert.NotEqual(t, fm3, fm4)
	assert.Equal(t, []*app.FocusManager{fm1, fm2, fm3, fm4}, s.ListFocusManagers())
	s.Add(o5)
	assert.Equal(t, []fyne.CanvasObject{o1, o2, o3, o4, o5}, s.List())
	assert.Equal(t, o5, s.Top())
	fm5 := s.TopFocusManager()
	assert.NotNil(t, fm5)
	assert.NotEqual(t, fm4, fm5)
	assert.Equal(t, []*app.FocusManager{fm1, fm2, fm3, fm4, fm5}, s.ListFocusManagers())
	s.Remove(o5)
	assert.Equal(t, []fyne.CanvasObject{o1, o2, o3, o4}, s.List())
	assert.Equal(t, o4, s.Top())
	assert.Equal(t, fm4, s.TopFocusManager())
	assert.Equal(t, []*app.FocusManager{fm1, fm2, fm3, fm4}, s.ListFocusManagers())
	// remove cuts the stack
	s.Remove(o2)
	assert.Equal(t, []fyne.CanvasObject{o1}, s.List())
	assert.Equal(t, o1, s.Top())
	assert.Equal(t, fm1, s.TopFocusManager())
	assert.Equal(t, []*app.FocusManager{fm1}, s.ListFocusManagers())
	s.Remove(o1)
	assert.Empty(t, s.List())
	assert.Nil(t, s.Top())
	assert.Nil(t, s.TopFocusManager())
	assert.Empty(t, s.ListFocusManagers())
}
