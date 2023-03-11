package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
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
	entry.SetText("The quick brown fox\njumped    over the lazy dog\n")
	entry.Resize(entry.MinSize())

	// select the word 'quick'
	ev := getClickPosition("The qui", 0)
	clickPrimary(entry, ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "quick", entry.SelectedText())

	// select the whitespace after 'quick'
	ev = getClickPosition("The quick", 0)
	clickPrimary(entry, ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, " ", entry.SelectedText())

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

	entry.propertyLock.RLock()
	assert.True(t, entry.selecting)
	entry.propertyLock.RUnlock()

	entry.DragEnd()
	assert.Equal(t, "", entry.SelectedText())
	entry.propertyLock.RLock()
	assert.False(t, entry.selecting)
	entry.propertyLock.RUnlock()

	// Test non-empty selection - drag from 'T' to 'g' (empty)
	ev1 = getClickPosition("", 0)
	de = &fyne.DragEvent{PointEvent: *ev1, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)
	de = &fyne.DragEvent{PointEvent: *ev2, Dragged: fyne.NewDelta(1, 0)}
	entry.Dragged(de)

	entry.propertyLock.RLock()
	assert.True(t, entry.selecting)
	entry.propertyLock.RUnlock()

	entry.DragEnd()
	assert.Equal(t, "Testing", entry.SelectedText())
	entry.propertyLock.RLock()
	assert.True(t, entry.selecting)
	entry.propertyLock.RUnlock()
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
	start, end := getTextWhitespaceRegion(str, -1)
	assert.Equal(t, -1, start)
	assert.Equal(t, -1, end)

	// select whitespace at the end of text
	start, end = getTextWhitespaceRegion(str, len(str))
	assert.Equal(t, 29, start)
	assert.Equal(t, 32, end)
	start, end = getTextWhitespaceRegion(str, len(str)+100)
	assert.Equal(t, 29, start)
	assert.Equal(t, 32, end)

	// select the whitespace
	start, end = getTextWhitespaceRegion(str, 0)
	assert.Equal(t, 0, start)
	assert.Equal(t, 1, end)

	// select "fish"
	start, end = getTextWhitespaceRegion(str, 1)
	assert.Equal(t, 1, start)
	assert.Equal(t, 5, end)
	start, end = getTextWhitespaceRegion(str, 4)
	assert.Equal(t, 1, start)
	assert.Equal(t, 5, end)

	// select "日本語日"
	start, end = getTextWhitespaceRegion(str, 6)
	assert.Equal(t, 6, start)
	assert.Equal(t, 10, end)
	start, end = getTextWhitespaceRegion(str, 9)
	assert.Equal(t, 6, start)
	assert.Equal(t, 10, end)

	// select "  \t  "
	start, end = getTextWhitespaceRegion(str, 10)
	assert.Equal(t, 10, start)
	assert.Equal(t, 15, end)

	// select "  \t"
	start, end = getTextWhitespaceRegion(str, 30)
	assert.Equal(t, 29, start)
	assert.Equal(t, len(str), end)
}

func TestEntry_ExpandSelectionWithWordSeparators(t *testing.T) {
	// select "is_a"
	str := []rune("This-is_a-test")
	start, end := getTextWhitespaceRegion(str, 6)
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
	var keyDown = func(key *fyne.KeyEvent) {
		e.KeyDown(key)
		e.TypedKey(key)
	}
	var keyPress = func(key *fyne.KeyEvent) {
		keyDown(key)
		e.KeyUp(key)
	}
	keyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	keyPress(&fyne.KeyEvent{Name: fyne.KeyRight})
	keyPress(&fyne.KeyEvent{Name: fyne.KeyRight})
	keyPress(&fyne.KeyEvent{Name: fyne.KeyRight})

	e.eraseSelection()
	e.updateText(e.textProvider().String())
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
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
	assert.False(t, entry.selecting)
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
	w.SetContent(entry)

	testContent := "test"

	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	clipboard.SetContent(testContent)

	entry.pasteFromClipboard(clipboard)

	assert.Equal(t, testContent, entry.Text)
}

func TestEntry_PasteFromClipboard_MultilineWrapping(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord

	w := test.NewApp().NewWindow("")
	w.SetContent(entry)
	w.Resize(fyne.NewSize(108, 64))

	test.Type(entry, "T")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	clipboard.SetContent("esting entry")

	entry.pasteFromClipboard(clipboard)

	assert.Equal(t, entry.Text, "Testing entry")
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 5, entry.CursorColumn)

	clipboard.SetContent(" paste\ncontent")
	entry.pasteFromClipboard(clipboard)

	assert.Equal(t, "Testing entry paste\ncontent", entry.Text)
	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 7, entry.CursorColumn)
}

func TestEntry_Tab(t *testing.T) {
	e := NewEntry()
	e.TextStyle.Monospace = true
	e.SetText("a\n\tb\nc")

	r := cache.Renderer(e.textProvider()).(*textRenderer)
	assert.Equal(t, 3, len(r.Objects()))
	assert.Equal(t, "a", r.Objects()[0].(*canvas.Text).Text)
	assert.Equal(t, "\tb", r.Objects()[1].(*canvas.Text).Text)

	w := test.NewWindow(e)
	w.Resize(fyne.NewSize(86, 86))
	w.Canvas().Focus(e)
	test.AssertImageMatches(t, "entry/tab-content.png", w.Canvas().Capture())
}

func TestEntry_TabSelection(t *testing.T) {
	e := NewEntry()
	e.SetText("a\n\tb\nc")
	e.TextStyle.Monospace = true

	e.CursorRow = 1
	e.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})

	assert.Equal(t, "\tb", e.SelectedText())

	w := test.NewWindow(e)
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
