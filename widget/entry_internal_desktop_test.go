//go:build !mobile

package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

func TestEntry_DragSelect(t *testing.T) {
	entry := NewEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = intWidget.ScrollNone
	entry.SetText("The quick brown fox jumped\nover the lazy dog\nThe quick\nbrown fox\njumped over the lazy dog\n")
	entry.Resize(entry.MinSize())

	// get position after the letter 'e' on the second row
	ev1 := getClickPosition("ove", 1)
	// get position after the letter 'z' on the second row
	ev2 := getClickPosition("over the laz", 1)

	// mouse down and drag from 'r' to 'z'
	me := &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.MouseButtonPrimary}
	entry.MouseDown(me)
	for ; ev1.Position.X < ev2.Position.X; ev1.Position.X++ {
		de := &fyne.DragEvent{PointEvent: *ev1, Dragged: fyne.NewDelta(1, 0)}
		entry.Dragged(de)
	}
	me = &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.MouseButtonPrimary}
	entry.MouseUp(me)

	assert.Equal(t, "r the laz", entry.SelectedText())
}

func TestEntry_DragSelectLargeStep(t *testing.T) {
	entry := NewEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = intWidget.ScrollNone
	entry.SetText("The quick brown fox jumped\nover the lazy dog\nThe quick\nbrown fox\njumped over the lazy dog\n")
	entry.Resize(entry.MinSize())

	// get position after the letter 'e' on the second row
	ev1 := getClickPosition("ove", 1)
	// get position after the letter 'z' on the second row
	ev2 := getClickPosition("over the laz", 1)

	// mouse down and drag from 'r' to 'z'
	me := &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.MouseButtonPrimary}
	entry.MouseDown(me)

	delta := ev2.Position.Subtract(ev1.Position)
	de := &fyne.DragEvent{PointEvent: *ev2, Dragged: fyne.NewDelta(delta.X, delta.Y)}
	entry.Dragged(de)

	me = &desktop.MouseEvent{PointEvent: *ev2, Button: desktop.MouseButtonPrimary}
	entry.MouseUp(me)

	assert.Equal(t, "r the laz", entry.SelectedText())
}

func TestEntry_DragSelectEmpty(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Testing")

	ev1 := getClickPosition("T", 0)
	ev2 := getClickPosition("Testing", 0)

	// Test empty selection - drag from 'e' to 'e' (empty)
	de := &fyne.DragEvent{PointEvent: *ev1, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)
	de = &fyne.DragEvent{PointEvent: *ev1, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)

	assert.True(t, entry.sel.selecting)

	entry.DragEnd()
	assert.Equal(t, "", entry.SelectedText())
	assert.False(t, entry.sel.selecting)

	// Test non-empty selection - drag from 'T' to 'g' (empty)
	ev1 = getClickPosition("", 0)
	de = &fyne.DragEvent{PointEvent: *ev1, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)
	de = &fyne.DragEvent{PointEvent: *ev2, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)

	assert.True(t, entry.sel.selecting)

	entry.DragEnd()
	assert.Equal(t, "Testing", entry.SelectedText())
	assert.True(t, entry.sel.selecting)
}

func TestEntry_DragSelectWithScroll(t *testing.T) {
	entry := NewEntry()
	entry.SetText("The quick brown fox jumped over and over the lazy dog.")

	// get position after the letter 'a'.
	ev1 := getClickPosition("The quick brown fox jumped over and over the la", 0)
	// get position after the letter 'u'
	ev2 := getClickPosition("The qu", 0)

	// mouse down and drag from 'a' to 'i'
	me := &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.MouseButtonPrimary}
	entry.MouseDown(me)
	de := &fyne.DragEvent{PointEvent: *ev1, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)
	de = &fyne.DragEvent{PointEvent: *ev2, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)
	me = &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.MouseButtonPrimary}
	entry.MouseUp(me)

	assert.Equal(t, "ick brown fox jumped over and over the la", entry.SelectedText())
}

func TestEntry_MouseClickAndDragOutsideText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("A\nB\n")

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev := &fyne.PointEvent{Position: pos}

	me := &desktop.MouseEvent{PointEvent: *ev, Button: desktop.MouseButtonPrimary}
	entry.MouseDown(me)
	de := &fyne.DragEvent{PointEvent: *ev, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)
	entry.MouseUp(me)
	assert.False(t, entry.sel.selecting)
}
