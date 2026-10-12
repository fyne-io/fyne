//go:build mobile

package widget

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
)

// On mobile a drag over text must scroll rather than select, so the text
// widgets are not Draggable and the gesture reaches the enclosing scroller.
func TestEntry_Drag_NotDraggable_Mobile(t *testing.T) {
	_, ok := any(&entryContent{}).(fyne.Draggable)
	assert.False(t, ok)
	_, ok = any(&focusSelectable{}).(fyne.Draggable)
	assert.False(t, ok)
}

// Entry keeps its public Dragged method everywhere, but on mobile it leaves the selection alone.
func TestEntry_Dragged_DoesNotSelect_Mobile(t *testing.T) {
	e := NewEntry()
	e.SetText("Hello World")
	e.Resize(e.MinSize())

	e.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)}, Dragged: fyne.NewDelta(20, 0)})
	e.DragEnd()
	assert.Empty(t, e.SelectedText())
}

func TestEntry_Drag_ScrollsOwnContent_Mobile(t *testing.T) {
	e := NewMultiLineEntry()
	e.SetText("one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten")
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(120, 60))

	test.Drag(w.Canvas(), fyne.NewPos(20, 20), 0, -20)
	assert.Empty(t, e.SelectedText())
	assert.Greater(t, e.scroller.Offset.Y, float32(0))
}

// An entry with no scroller of its own stays Draggable for API compatibility, so it
// receives the drag itself. It must still leave the selection alone.
func TestEntry_Drag_NoScroller_DoesNotSelect_Mobile(t *testing.T) {
	e := NewEntry()
	e.Wrapping = fyne.TextWrapOff
	e.Scroll = fyne.ScrollNone
	e.SetText("Hello World")
	filler := canvas.NewRectangle(color.Black)
	filler.SetMinSize(fyne.NewSize(10, 300))
	scroll := intWidget.NewVScroll(&fyne.Container{Layout: layout.NewVBoxLayout(), Objects: []fyne.CanvasObject{e, filler}})
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(120, 60))

	test.Drag(w.Canvas(), fyne.NewPos(20, 15), 0, -20)
	assert.Empty(t, e.SelectedText())
}

func TestLabel_Drag_ScrollsOuterContainer_Mobile(t *testing.T) {
	l := NewLabel("Hello World")
	l.Selectable = true
	filler := canvas.NewRectangle(color.Black)
	filler.SetMinSize(fyne.NewSize(10, 300))
	scroll := intWidget.NewVScroll(&fyne.Container{Layout: layout.NewVBoxLayout(), Objects: []fyne.CanvasObject{l, filler}})
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(120, 60))

	test.Drag(w.Canvas(), fyne.NewPos(20, 15), 0, -20)
	assert.Empty(t, l.SelectedText())
	assert.Greater(t, scroll.Offset.Y, float32(0))
}
