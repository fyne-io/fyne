//go:build mobile && (!windows || !ci)

package mobile

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func Test_canvas_Dragged(t *testing.T) {
	dragged := false
	var draggedObj fyne.Draggable
	scroll := container.NewScroll(widget.NewLabel("Hi\nHi\nHi"))
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(scroll)
	c.Resize(fyne.NewSize(40, 24))
	assert.Equal(t, float32(0), scroll.Offset.Y)

	c.tapDown(fyne.NewPos(32, 3), 0)
	c.tapMove(fyne.NewPos(32, 10), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
		draggedObj = wid
	})

	assert.True(t, dragged)
	assert.Equal(t, scroll, draggedObj)
	dragged = false
	c.tapMove(fyne.NewPos(32, 5), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Equal(t, fyne.NewPos(0, 5), scroll.Offset)
}

func Test_canvas_Dragged_nested(t *testing.T) {
	dragged := false
	var draggedObj fyne.Draggable
	childScroll := container.NewHScroll(widget.NewLabel("Child *Scroll"))
	parentScroll := container.NewScroll(container.NewVBox(childScroll, widget.NewLabel("Hi\nHi\nHi")))
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(parentScroll)
	c.Resize(fyne.NewSize(40, 24))
	assert.Equal(t, fyne.NewPos(0, 0), parentScroll.Offset)
	assert.Equal(t, fyne.NewPos(0, 0), childScroll.Offset)

	c.tapDown(fyne.NewPos(10, 10), 0)
	c.tapMove(fyne.NewPos(15, 15), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
		draggedObj = wid
	})
	assert.True(t, dragged)
	assert.Equal(t, parentScroll, draggedObj)
	assert.Equal(t, fyne.NewPos(0, 0), parentScroll.Offset)
	assert.Equal(t, fyne.NewPos(0, 0), childScroll.Offset)

	dragged = false
	c.tapMove(fyne.NewPos(9, 7), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Equal(t, fyne.NewPos(0, 8), parentScroll.Offset)
	assert.Equal(t, fyne.NewPos(6, 0), childScroll.Offset)
}
