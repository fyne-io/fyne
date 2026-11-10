package widget

import (
	"strings"
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
	assert.Equal(t, textSpace, entry.SelectedText())

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

func TestEntry_FocusKeepsScrollOffset(t *testing.T) {
	e := NewMultiLineEntry()
	e.Validator = func(string) error { return nil }
	e.SetText(strings.Repeat("line\n", 50))
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(150, 100))

	w.Canvas().Focus(e)
	e.scroller.ScrollToOffset(fyne.NewPos(0, 200))
	assert.Equal(t, float32(200), e.scroller.Offset.Y)

	w.Canvas().Unfocus()
	assert.Equal(t, float32(200), e.scroller.Offset.Y)
	w.Canvas().Focus(e)
	assert.Equal(t, float32(200), e.scroller.Offset.Y)

	e.TypedRune('a')
	assert.Zero(t, e.scroller.Offset.Y)
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

func TestEntry_ContentSizeAndPlacementWithIcon(t *testing.T) {
	entry := NewEntry()
	entry.SetIcon(theme.MailComposeIcon())
	entry.SetText("SomeText")
	renderer := entry.CreateRenderer()
	contentPos := fyne.NewPos(theme.LineSpacing()+theme.IconInlineSize(), theme.InputBorderSize())

	renderer.Layout(entry.MinSize())
	// Scrollable content should be positioned after the icon, with correct padding
	assert.Equal(t, contentPos, entry.scroller.Position())

	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = fyne.ScrollNone
	renderer.Layout(entry.MinSize())
	// Content should be positioned after the icon, with correct padding
	assert.Equal(t, contentPos, entry.content.Position())
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
	entry.Validator = func(string) error {
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
	assert.Equal(t, e.TextStyle, e.placeholderWidget.Segments[0].(*TextSegment).Style.TextStyle)

	w.Canvas().Focus(e)
	assert.Equal(t, e.TextStyle, e.placeholderWidget.Segments[0].(*TextSegment).Style.TextStyle)
}

func TestEntry_ClickAfterScroll(t *testing.T) {
	e := NewMultiLineEntry()
	e.SetText(strings.Repeat("line\n", 50))
	w := test.NewTempWindow(t, e)
	w.Resize(fyne.NewSize(150, 100))

	e.scroller.ScrollToOffset(fyne.NewPos(0, 200))
	clickPrimary(e, getClickPosition("li", 2))

	assert.Equal(t, float32(200), e.scroller.Offset.Y)
	cursorY := e.CursorPosition().Y
	assert.GreaterOrEqual(t, cursorY, e.scroller.Offset.Y)
	assert.Less(t, cursorY, e.scroller.Offset.Y+e.scroller.Size().Height)
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

// A line that ends with a space at the word wrap boundary used to lose that space
// from its rows, so the position after it resolved on to the row below and joining
// the next line with Backspace left the cursor at the end of the line after it.
func TestEntry_Backspace_TrailingSpaceAtWrapBoundary(t *testing.T) {
	line1 := "aa bb "
	text := line1 + "\n\nmore text here"
	lineEnd := len([]rune(line1))

	wrapped := 0
	for width := float32(40); width <= 160; width += 4 {
		e := NewMultiLineEntry()
		e.Wrapping = fyne.TextWrapWord
		w := test.NewTempWindow(t, e)
		w.Resize(fyne.NewSize(width, 150))
		e.SetText(text)

		provider := e.textProvider()
		row, col := e.rowColFromTextPos(lineEnd)
		assert.LessOrEqual(t, col, provider.rowLength(row),
			"width %v: the end of the line is not inside any row", width)
		if provider.rowBoundary(row).docBegin > 0 {
			wrapped++ // the trailing space is on a row that the wrapping produced
		}

		// put the cursor at the start of the empty line and join it to the one above
		e.CursorRow, e.CursorColumn = e.rowColFromTextPos(lineEnd + 1)
		e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})

		assert.Equal(t, line1+"\nmore text here", e.Text, "width %v", width)
		assert.Equal(t, lineEnd, e.CursorTextOffset(), "width %v", width)
	}

	assert.NotZero(t, wrapped, "no width wrapped the first line")
}

// Resizing a wrapped Entry moves the wrap points, so the cursor has to follow the
// position it was on to whichever row now holds it. Comparing text positions could
// not see that for a cursor on the first row, where the position is just the column,
// and left the cursor stranded past the end of that row.
func TestEntry_Resize_CursorFollowsWrap(t *testing.T) {
	text := "one two three four five six " // six words, and a space at the end
	end := len([]rune(text))

	wrapped := 0
	for width := float32(40); width <= 200; width += 4 {
		e := NewMultiLineEntry()
		e.Wrapping = fyne.TextWrapWord
		w := test.NewTempWindow(t, e)
		w.Resize(fyne.NewSize(400, 150))
		e.SetText(text)

		e.CursorRow, e.CursorColumn = e.rowColFromTextPos(end)
		assert.Equal(t, 0, e.CursorRow, "the text did not start out on one row")

		w.Resize(fyne.NewSize(width, 150))

		provider := e.textProvider()
		assert.Equal(t, end, provider.rowBoundary(provider.rows()-1).docEnd,
			"width %v: the trailing space is not in any row", width)

		row, col := e.rowColFromTextPos(end)
		assert.Equal(t, row, e.CursorRow, "width %v", width)
		assert.Equal(t, col, e.CursorColumn, "width %v", width)
		assert.LessOrEqual(t, e.CursorColumn, provider.rowLength(e.CursorRow),
			"width %v: the cursor is past the end of its own row", width)

		if e.CursorRow > 0 {
			wrapped++
		}
	}

	assert.NotZero(t, wrapped, "no width wrapped the text, so the cursor never had to move")
}
