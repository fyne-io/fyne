package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func clickPrimary(e *Entry, ev *fyne.PointEvent) {
	mouseEvent := &desktop.MouseEvent{
		PointEvent: *ev,
		Button:     desktop.MouseButtonPrimary,
	}
	e.MouseDown(mouseEvent)
	e.MouseUp(mouseEvent)
	e.Tapped(ev) // in the glfw driver there is a double click delay before Tapped()
}

func TestEntry_Cursor(t *testing.T) {
	entry := NewEntry()
	assert.Equal(t, desktop.TextCursor, entry.Cursor())
}

func TestEntry_DoubleTapped(t *testing.T) {
	entry := NewEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = intWidget.ScrollNone
	entry.SetText("The quick brown fox\njumped    over the lazy dog\n")
	entry.Resize(entry.MinSize())

	// select the word 'quick'
	ev := getClickPosition("The qui", 0)
	clickPrimary(entry, ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "quick", entry.SelectedText())

	entry.sel.doubleTappedAtUnixMillis = 0 // make sure we don't register a triple tap next

	// select the whitespace after 'quick'
	ev = getClickPosition("The quick", 0)
	clickPrimary(entry, ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, " ", entry.SelectedText())

	entry.sel.doubleTappedAtUnixMillis = 0

	// select all whitespace after 'jumped'
	ev = getClickPosition("jumped  ", 1)
	clickPrimary(entry, ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "    ", entry.SelectedText())
}

func TestEntry_DoubleTapped_AfterCol(t *testing.T) {
	entry := NewEntry()
	entry.SetText("A\nB\n")

	window := test.NewWindow(entry)
	defer window.Close()
	window.SetPadded(false)
	window.Resize(entry.MinSize())
	entry.Resize(entry.MinSize())
	c := window.Canvas()

	ev := getClickPosition("", 0)
	clickPrimary(entry, ev)
	assert.Equal(t, entry, c.Focused())

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev = &fyne.PointEvent{Position: pos}
	clickPrimary(entry, ev)
	entry.DoubleTapped(ev)

	assert.Equal(t, "", entry.SelectedText())
}

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

func TestEntry_ExpandSelectionForDoubleTap(t *testing.T) {
	str := []rune(" fish 日本語日  \t  test 本日本 moose  \t")

	// select invalid (before start)
	start, end := getTextWhitespaceRegion(str, -1, false)
	assert.Equal(t, -1, start)
	assert.Equal(t, -1, end)

	// select whitespace at the end of text
	start, end = getTextWhitespaceRegion(str, len(str), false)
	assert.Equal(t, 29, start)
	assert.Equal(t, 32, end)
	start, end = getTextWhitespaceRegion(str, len(str)+100, false)
	assert.Equal(t, 29, start)
	assert.Equal(t, 32, end)

	// select the whitespace
	start, end = getTextWhitespaceRegion(str, 0, false)
	assert.Equal(t, 0, start)
	assert.Equal(t, 1, end)
	// select the whitespace - grab adjacent words
	start, end = getTextWhitespaceRegion(str, 0, true)
	assert.Equal(t, 0, start)
	assert.Equal(t, 5, end)

	// select "fish"
	start, end = getTextWhitespaceRegion(str, 1, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 5, end)
	start, end = getTextWhitespaceRegion(str, 4, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 5, end)

	// select "日本語日"
	start, end = getTextWhitespaceRegion(str, 7, false)
	assert.Equal(t, 6, start)
	assert.Equal(t, 10, end)
	start, end = getTextWhitespaceRegion(str, 9, false)
	assert.Equal(t, 6, start)
	assert.Equal(t, 10, end)

	// select "  \t  "
	start, end = getTextWhitespaceRegion(str, 10, false)
	assert.Equal(t, 10, start)
	assert.Equal(t, 15, end)

	// select "  \t"
	start, end = getTextWhitespaceRegion(str, 30, false)
	assert.Equal(t, 29, start)
	assert.Equal(t, len(str), end)
}

func TestEntry_ExpandSelectionWithWordSeparators(t *testing.T) {
	// select "is_a"
	str := []rune("This-is_a-test")
	start, end := getTextWhitespaceRegion(str, 6, false)
	assert.Equal(t, 5, start)
	assert.Equal(t, 9, end)
}

func TestEntry_EraseSelection(t *testing.T) {
	// Selects "sti" on border 2 of a new multiline
	// T e s t i n g
	// T e[s t i]n g
	// T e s t i n g
	e := NewMultiLineEntry()
	e.SetText("Testing\nTesting\nTesting")
	e.CursorRow = 1
	e.CursorColumn = 2
	e.sel.cursorRow, e.sel.cursorRow = e.CursorRow, e.CursorColumn
	keyDown := func(key *fyne.KeyEvent) {
		e.KeyDown(key)
		e.TypedKey(key)
	}
	keyPress := func(key *fyne.KeyEvent) {
		keyDown(key)
		e.KeyUp(key)
	}
	keyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	keyPress(&fyne.KeyEvent{Name: fyne.KeyRight})
	keyPress(&fyne.KeyEvent{Name: fyne.KeyRight})
	keyPress(&fyne.KeyEvent{Name: fyne.KeyRight})

	_ = e.Theme()
	e.eraseSelectionAndUpdate()
	e.updateText(e.textProvider().String(), false)
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	a, b := e.sel.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}

func TestEntry_CallbackLocking(t *testing.T) {
	e := &Entry{}
	called := 0
	e.OnChanged = func(_ string) {
		called++ // Just to not have an empty critical section.
	}

	_ = e.Theme()
	test.Type(e, "abc123")
	e.selectAll()
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, 7, called)
}

func TestEntry_IconContentSizeAndPlacement(t *testing.T) {
	entry := NewEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = fyne.ScrollNone
	icon := theme.MailComposeIcon()
	entry.SetIcon(icon)
	entry.SetText("SomeText")
	r := test.TempWidgetRenderer(t, entry)
	r.Layout(entry.MinSize())

	var iconObj *canvas.Image
	for _, obj := range r.Objects() {
		if img, ok := obj.(*canvas.Image); ok && img.Resource == icon {
			iconObj = img
			break
		}
	}

	// Icon should be at the left, with correct size
	assert.NotNil(t, iconObj)
	assert.Equal(t, theme.IconInlineSize(), iconObj.Size().Width)
	assert.Equal(t, theme.IconInlineSize(), iconObj.Size().Height)
	assert.Equal(t, fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()), iconObj.Position())
	// Content should be positioned after the icon, with correct padding
	assert.Equal(t, fyne.NewPos(theme.InnerPadding()+theme.LineSpacing()+theme.IconInlineSize(), theme.InputBorderSize()), entry.content.Position())
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

func TestEntry_MouseDownOnSelect(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Ahnj\nBuki\n")
	entry.TypedShortcut(&fyne.ShortcutSelectAll{})

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev := &fyne.PointEvent{Position: pos}

	me := &desktop.MouseEvent{PointEvent: *ev, Button: desktop.MouseButtonSecondary}
	entry.MouseDown(me)
	entry.MouseUp(me)

	assert.Equal(t, "Ahnj\nBuki\n", entry.SelectedText())

	me = &desktop.MouseEvent{PointEvent: *ev, Button: desktop.MouseButtonPrimary}
	entry.MouseDown(me)
	entry.MouseUp(me)

	assert.Equal(t, "", entry.SelectedText())
}

func TestEntry_PasteFromClipboard(t *testing.T) {
	entry := NewEntry()

	w := test.NewApp().NewWindow("")
	defer w.Close()
	w.SetContent(entry)

	testContent := "test"

	clipboard := fyne.CurrentApp().Clipboard()
	clipboard.SetContent(testContent)

	entry.pasteFromClipboard(clipboard)

	assert.Equal(t, testContent, entry.Text)
}

func TestEntry_PasteFromClipboard_MultilineWrapping(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord

	w := test.NewApp().NewWindow("")
	defer w.Close()
	w.SetContent(entry)
	w.Resize(fyne.NewSize(108, 64))

	test.Type(entry, "T")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	clipboard := fyne.CurrentApp().Clipboard()
	clipboard.SetContent("esting entry")

	entry.pasteFromClipboard(clipboard)

	assert.Equal(t, "Testing entry", entry.Text)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 5, entry.CursorColumn)

	clipboard.SetContent(" paste\ncontent")
	entry.pasteFromClipboard(clipboard)

	assert.Equal(t, "Testing entry paste\ncontent", entry.Text)
	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 7, entry.CursorColumn)
}

func TestEntry_PasteFromClipboardValidation(t *testing.T) {
	entry := NewEntry()
	var triggered int
	entry.Validator = func(s string) error {
		triggered++
		return nil
	}

	testContent := "test"
	clipboard := test.NewTempApp(t).Clipboard()
	clipboard.SetContent(testContent)

	entry.pasteFromClipboard(clipboard)
	assert.Equal(t, 2, triggered)
}

func TestEntry_PlaceholderTextStyle(t *testing.T) {
	e := NewEntry()
	e.TextStyle = fyne.TextStyle{Bold: true, Italic: true}

	w := test.NewTempWindow(t, e)
	assert.Equal(t, e.TextStyle, e.placeholder.Segments[0].(*TextSegment).Style.TextStyle)

	w.Canvas().Focus(e)
	assert.Equal(t, e.TextStyle, e.placeholder.Segments[0].(*TextSegment).Style.TextStyle)
}

func TestEntry_Tab(t *testing.T) {
	e := NewEntry()
	e.TextStyle.Monospace = true
	e.SetText("a\n\tb\nc")

	_ = e.Theme()
	r := cache.Renderer(e.textProvider()).(*textRenderer)
	assert.Len(t, r.Objects(), 3)
	assert.Equal(t, "a", r.Objects()[0].(*canvas.Text).Text)
	assert.Equal(t, "\tb", r.Objects()[1].(*canvas.Text).Text)

	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(86, 86))
	w.Canvas().Focus(e)
	test.AssertImageMatches(t, "entry/tab-content.png", w.Canvas().Capture())
}

func TestEntry_TabSelection(t *testing.T) {
	e := NewEntry()
	e.SetText("a\n\tb\nc")
	e.TextStyle.Monospace = true

	e.CursorRow = 1
	e.sel.cursorRow = 1
	e.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})

	assert.Equal(t, "\tb", e.SelectedText())

	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(86, 86))
	w.Canvas().Focus(e)
	test.AssertImageMatches(t, "entry/tab-select.png", w.Canvas().Capture())
}

func TestEntry_ShiftSelection_ResetOnFocusLost(t *testing.T) {
	e := NewEntry()
	e.SetText("Hello")

	e.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	assert.True(t, e.selectKeyDown)

	e.FocusLost()
	assert.False(t, e.selectKeyDown)
}

func getClickPosition(str string, row int) *fyne.PointEvent {
	x := fyne.MeasureText(str, theme.TextSize(), fyne.TextStyle{}).Width + theme.Padding()

	rowHeight := fyne.MeasureText("M", theme.TextSize(), fyne.TextStyle{}).Height
	y := float32(row)*rowHeight + rowHeight/2

	// add a couple of pixels, this is currently a workaround for weird mouse to column logic on text with kerning
	pos := fyne.NewPos(x+2, y)
	return &fyne.PointEvent{Position: pos}
}
