package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
)

func dragSelect(sel *focusSelectable, from, to fyne.Position) {
	sel.MouseDown(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: from},
	})
	sel.Dragged(&fyne.DragEvent{
		Dragged:    fyne.Delta{DX: to.X - from.X, DY: to.Y - from.Y},
		PointEvent: fyne.PointEvent{Position: to},
	})
	sel.DragEnd()
	sel.MouseUp(&desktop.MouseEvent{
		Button:     desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: to},
	})
}

func TestRichText_Select(t *testing.T) {
	rich := NewRichTextWithText("Hello")
	assert.Len(t, test.WidgetRenderer(rich).Objects(), 1)

	rich.Selectable = true
	rich.Refresh()
	rich.Resize(rich.MinSize())
	assert.Empty(t, rich.SelectedText())
	assert.Len(t, test.WidgetRenderer(rich).Objects(), 2)

	sel := test.WidgetRenderer(rich).Objects()[0].(*focusSelectable)
	assert.Equal(t, rich.Size(), sel.Size())
	dragSelect(sel, fyne.NewPos(15, 10), fyne.NewPos(30, 10))
	assert.Equal(t, "el", rich.SelectedText())

	sel.TypedShortcut(&fyne.ShortcutCopy{})
	assert.Equal(t, "el", fyne.CurrentApp().Clipboard().Content())

	rich.ClearSelection()
	assert.Empty(t, rich.SelectedText())

	dragSelect(sel, fyne.NewPos(15, 10), fyne.NewPos(30, 10))
	rich.Selectable = false
	rich.Refresh()
	assert.Len(t, test.WidgetRenderer(rich).Objects(), 1)
	assert.Empty(t, rich.SelectedText())
}

func TestRichText_Select_AcrossSegments(t *testing.T) {
	rich := NewRichText(
		&TextSegment{Text: "Hello ", Style: RichTextStyleInline},
		&TextSegment{Text: "bold", Style: RichTextStyleStrong},
		&TextSegment{Text: " first", Style: RichTextStyleParagraph},
		&TextSegment{Text: "Second", Style: RichTextStyleParagraph},
	)
	rich.Selectable = true
	rich.Resize(rich.MinSize())

	sel := test.WidgetRenderer(rich).Objects()[0].(*focusSelectable)
	sel.DoubleTapped(&fyne.PointEvent{Position: fyne.NewPos(15, 10)})
	assert.Equal(t, "Hello", rich.SelectedText())

	y2, _ := rich.rowGeometry(1)
	dragSelect(sel, fyne.NewPos(0, 10), fyne.NewPos(rich.Size().Width, y2+10))
	assert.Equal(t, "Hello bold first\nSecond", rich.SelectedText())
}

func TestRichText_Select_Blocks(t *testing.T) {
	rich := NewRichTextFromMarkdown("Intro\n\n* one\n* two\n\n```\ncode\nmore\n```\n")
	rich.Selectable = true
	rich.Resize(rich.MinSize())

	sel := test.WidgetRenderer(rich).Objects()[0].(*focusSelectable)
	dragSelect(sel, fyne.NewPos(0, 10), fyne.NewPos(rich.Size().Width, rich.Size().Height))
	assert.Equal(t, "Intro\n• one\n• two\ncode\nmore", rich.SelectedText())
}

func TestRichText_Select_Wrapped(t *testing.T) {
	rich := NewRichTextWithText("Hello World")
	rich.Wrapping = fyne.TextWrapWord
	rich.Selectable = true
	rich.Resize(fyne.NewSize(60, 100))
	assert.Equal(t, 2, rich.rows())

	sel := test.WidgetRenderer(rich).Objects()[0].(*focusSelectable)
	dragSelect(sel, fyne.NewPos(0, 10), fyne.NewPos(60, 90))
	assert.Equal(t, "Hello World", rich.SelectedText())
}

func TestRichText_Select_ContentChanged(t *testing.T) {
	rich := NewRichTextWithText("Hello World")
	rich.Selectable = true
	rich.Resize(rich.MinSize())

	sel := test.WidgetRenderer(rich).Objects()[0].(*focusSelectable)
	dragSelect(sel, fyne.NewPos(0, 10), fyne.NewPos(rich.Size().Width, 10))
	assert.Equal(t, "Hello World", rich.SelectedText())

	rich.Segments[0].(*TextSegment).Text = "Hi"
	rich.Refresh()
	assert.Empty(t, rich.SelectedText())
}

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
