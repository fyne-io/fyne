//go:build !mobile

package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestLabel_Select(t *testing.T) {
	l := NewLabel("Hello")
	l.Selectable = true

	assert.Empty(t, l.SelectedText())
	assert.Equal(t, 2, len(test.WidgetRenderer(l).Objects()))

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	sel.MouseDown(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 10)},
	})
	sel.Dragged(&fyne.DragEvent{
		Dragged:    fyne.Delta{DX: 15, DY: 0},
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)},
	})
	sel.DragEnd()
	sel.MouseUp(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)},
	})
	assert.Equal(t, "el", l.SelectedText())

	sel.TypedShortcut(&fyne.ShortcutCopy{})
	assert.Equal(t, "el", fyne.CurrentApp().Clipboard().Content())

	l.SetText("World")
	assert.Equal(t, "", l.SelectedText())

	l.Selectable = false
	l.Refresh()
	assert.Equal(t, 1, len(test.WidgetRenderer(l).Objects()))
	assert.Empty(t, l.SelectedText())
}

func TestLabel_ClearSelection(t *testing.T) {
	l := NewLabel("Hello")
	l.Selectable = true

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	sel.MouseDown(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 10)},
	})
	sel.Dragged(&fyne.DragEvent{
		Dragged:    fyne.Delta{DX: 15, DY: 0},
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)},
	})
	sel.DragEnd()
	sel.MouseUp(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)},
	})
	assert.Equal(t, "el", l.SelectedText())

	l.ClearSelection()
	assert.Equal(t, "", l.SelectedText())

	// calling again with nothing selected is a no-op
	l.ClearSelection()
	assert.Equal(t, "", l.SelectedText())
}

func TestLabelSizeNameWithSelection(t *testing.T) {
	l := NewLabel("Hello")
	l.Selectable = true

	w := test.NewWindow(l)
	defer w.Close()

	assert.Empty(t, l.SelectedText())
	assert.Equal(t, 2, len(test.WidgetRenderer(l).Objects()))

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	sel.MouseDown(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 10)},
	})
	sel.Dragged(&fyne.DragEvent{
		Dragged:    fyne.Delta{DX: 15, DY: 0},
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)},
	})
	sel.DragEnd()
	sel.MouseUp(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)},
	})
	assert.Equal(t, "el", l.SelectedText())

	test.AssertRendersToImage(t, "label/label_selection_defaultsize.png", w.Canvas())

	l.SizeName = theme.SizeNameHeadingText
	l.Refresh()
	w.SetContent(l) // updates window size

	test.AssertRendersToImage(t, "label/label_selection_headersize.png", w.Canvas())
}
