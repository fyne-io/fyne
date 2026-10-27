package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
)

func TestRichText_Select_Panel(t *testing.T) {
	rich := NewRichTextFromMarkdown("Intro\n\n```\ncode here\n```\n")
	rich.Selectable = true
	rich.Resize(rich.MinSize())

	r := test.WidgetRenderer(rich)
	before := len(r.Objects())
	sel := r.Objects()[0].(*focusSelectable)

	row := rich.rows() - 1
	y, h := rich.rowGeometry(row)
	sel.DoubleTapped(&fyne.PointEvent{Position: fyne.NewPos(rich.Size().Width/4, y+h/2+rich.Theme().Size("innerPadding"))})
	assert.NotEmpty(t, rich.SelectedText())

	objs := r.Objects()
	assert.Len(t, objs, before+1)
	// the highlight sits above the selection widget and panel, below the text
	_, isRect := objs[2].(*canvas.Rectangle)
	assert.True(t, isRect)
}

func TestRichText_Select_Scroll(t *testing.T) {
	rich := NewRichTextFromMarkdown("Intro\n\n```\ncode here\n```\n\nOne\n\nTwo\n\nThree")
	rich.Scroll = widget.ScrollBoth
	rich.Selectable = true
	w := test.NewTempWindow(t, rich)
	w.Resize(fyne.NewSize(100, 60))

	r := test.WidgetRenderer(rich)
	assert.Len(t, r.Objects(), 1) // just the scroller
	inner := scrollInnerContainer(rich.scr)
	sel := inner.Objects[0].(*focusSelectable)
	assert.Equal(t, rich.scr.Content.Size(), sel.Size())
	before := len(inner.Objects)

	sel.DoubleTapped(&fyne.PointEvent{Position: fyne.NewPos(15, 10)})
	assert.Equal(t, "Intro", rich.SelectedText())
	assert.Len(t, scrollInnerContainer(rich.scr).Objects, before+1)

	rich.ClearSelection()
	assert.Len(t, scrollInnerContainer(rich.scr).Objects, before)
}

func TestRichText_Select_TappedSecondaryPosition(t *testing.T) {
	rich := NewRichTextWithText("Hello")
	rich.Selectable = true
	rich.Resize(rich.MinSize())
	w := test.NewTempWindow(t, rich)
	w.Resize(fyne.NewSize(200, 200))

	sel := test.WidgetRenderer(rich).Objects()[0].(*focusSelectable)

	rel := fyne.NewPos(12, 10)
	// On mobile the tap's AbsolutePosition is relative to the whole canvas while the
	// pop-up must be placed relative to the interactive area. Simulate that mismatch
	// so this test fails if the absolute position is used directly.
	absolute := fyne.NewPos(500, 500)
	sel.TappedSecondary(&fyne.PointEvent{Position: rel, AbsolutePosition: absolute})

	top := w.Canvas().Overlays().Top()
	require.NotNil(t, top)
	pop := top.(*widget.OverlayContainer).Content.(*PopUpMenu)

	expected := fyne.CurrentApp().Driver().AbsolutePositionForObject(sel).Add(rel)
	assert.Equal(t, expected, pop.Position())
}
