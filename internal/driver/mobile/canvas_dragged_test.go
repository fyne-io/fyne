//go:build mobile && (!windows || !ci)

package mobile

import (
	"strings"
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

	c.tapDown(fyne.NewPos(20, 5), 0)
	c.tapMove(fyne.NewPos(20, 10), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
		draggedObj = wid
	})

	assert.True(t, dragged)
	assert.Equal(t, scroll, draggedObj)
	dragged = false
	c.tapMove(fyne.NewPos(20, 5), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
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

// Test_canvas_DraggedCompanionUsesPressTarget covers
// https://github.com/fyne-io/fyne/issues/6666
//
// The drag starts on the left column's horizontal scroller, but the first move
// lands on the right column's horizontal scroller. The companion scroller must
// be the pressed column's vertical parent, not the one under the later pointer
// position.
func Test_canvas_DraggedCompanionUsesPressTarget(t *testing.T) {
	leftHeader := container.NewHScroll(widget.NewLabel("Left header"))
	leftParent := container.NewVScroll(container.NewVBox(leftHeader,
		widget.NewLabel(strings.Repeat("left tall content\n", 40))))
	rightHeader := container.NewHScroll(widget.NewLabel("Right header"))
	rightParent := container.NewVScroll(container.NewVBox(rightHeader,
		widget.NewLabel(strings.Repeat("right tall content\n", 40))))

	content := container.NewWithoutLayout(leftParent, rightParent)
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.padded = false
	c.SetContent(content)
	c.Resize(fyne.NewSize(200, 120))
	leftParent.Move(fyne.NewPos(0, 0))
	leftParent.Resize(fyne.NewSize(100, 120))
	rightParent.Move(fyne.NewPos(100, 0))
	rightParent.Resize(fyne.NewSize(100, 120))

	var targets []fyne.Draggable
	c.tapDown(fyne.NewPos(75, 15), 0)
	c.tapMove(fyne.NewPos(125, 5), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		targets = append(targets, wid)
		wid.Dragged(ev)
	})

	assert.Equal(t, []fyne.Draggable{leftHeader, leftParent}, targets)
	assert.Equal(t, fyne.NewPos(0, 10), leftParent.Offset)
	assert.Equal(t, fyne.NewPos(0, 0), rightParent.Offset)
}
