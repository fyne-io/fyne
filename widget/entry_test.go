package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func entryRenderTexts(e *Entry) []*canvas.Text {
	textWid := Renderer(e).(*entryRenderer).text
	return Renderer(textWid).(*textRenderer).texts
}

func entryRenderPlaceholderTexts(e *Entry) []*canvas.Text {
	textWid := Renderer(e).(*entryRenderer).placeholder
	return Renderer(textWid).(*textRenderer).texts
}

func TestEntry_MinSize(t *testing.T) {
	entry := NewEntry()
	min := entry.MinSize()
	entry.SetPlaceHolder("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetText("")
	assert.Equal(t, min, entry.MinSize())
	entry.SetPlaceHolder("Hi")
	assert.True(t, entry.MinSize().Width > min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestMultiLineEntry_MinSize(t *testing.T) {
	entry := NewEntry()
	entry.MinSize()
	singleMin := entry.MinSize()

	multi := NewMultiLineEntry()
	multiMin := multi.MinSize()

	assert.Equal(t, singleMin.Width, multiMin.Width)
	assert.True(t, multiMin.Height > singleMin.Height)

	multi.MultiLine = false
	multiMin = multi.MinSize()
	assert.Equal(t, singleMin.Height, multiMin.Height)
}

func TestEntry_SetPlaceHolder(t *testing.T) {
	entry := NewEntry()

	assert.Equal(t, 0, len(entry.Text))
	assert.Equal(t, 0, entry.textProvider().len())

	entry.SetPlaceHolder("Test")
	assert.Equal(t, 0, len(entry.Text))
	assert.Equal(t, 0, entry.textProvider().len())
	assert.Equal(t, 4, entry.placeholderProvider().len())
	assert.False(t, entry.placeholderProvider().Hidden)

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	assert.True(t, entry.placeholderProvider().Hidden)

	assert.Equal(t, 2, entry.textProvider().len())
}

func TestEntry_SetTextEmptyString(t *testing.T) {
	entry := NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	entry = NewMultiLineEntry()
	test.Type(entry, "test\ntest")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)

	assert.Equal(t, 4, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_OnKeyDown(t *testing.T) {
	entry := NewEntry()

	test.Type(entry, "Hi")

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_SetReadOnly_KeyDown(t *testing.T) {
	entry := NewEntry()

	test.Type(entry, "H")
	entry.SetReadOnly(true)
	test.Type(entry, "i")
	assert.Equal(t, "H", entry.Text)

	entry.SetReadOnly(false)
	test.Type(entry, "i")
	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_SetReadOnly_OnFocus(t *testing.T) {
	entry := NewEntry()
	entry.SetReadOnly(true)

	entry.FocusGained()
	assert.False(t, entry.Focused())

	entry.SetReadOnly(false)
	entry.FocusGained()
	assert.True(t, entry.Focused())
}

func TestEntry_OnKeyDown_Insert(t *testing.T) {
	entry := NewEntry()

	test.Type(entry, "Hi")
	assert.Equal(t, "Hi", entry.Text)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)

	test.Type(entry, "o")
	assert.Equal(t, "Hoi", entry.Text)
}

func TestEntry_OnKeyDown_Newline(t *testing.T) {
	entry := &Entry{MultiLine: true}
	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyReturn}
	entry.TypedKey(key)

	assert.Equal(t, "H\ni", entry.Text)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "o")
	assert.Equal(t, "H\noi", entry.textProvider().String())
	assert.Equal(t, "H", entryRenderTexts(entry)[0].Text)
	assert.Equal(t, "oi", entryRenderTexts(entry)[1].Text)
}

func TestEntry_OnKeyDown_Backspace(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_BackspaceBeyondText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceNewline(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace_Unicode(t *testing.T) {
	entry := NewEntry()

	test.Type(entry, "è")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	bs := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(bs)
	assert.Equal(t, "", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Delete(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_DeleteBeyondText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_DeleteNewline(t *testing.T) {
	entry := NewEntry()
	entry.SetText("H\ni")

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Home_End(t *testing.T) {
	entry := &Entry{}
	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	end := &fyne.KeyEvent{Name: fyne.KeyEnd}
	entry.TypedKey(end)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	home := &fyne.KeyEvent{Name: fyne.KeyHome}
	entry.TypedKey(home)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntryNotify(t *testing.T) {
	entry := NewEntry()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntryFocus(t *testing.T) {
	entry := NewEntry()

	entry.FocusGained()
	assert.True(t, entry.Focused())

	entry.FocusLost()
	assert.False(t, entry.Focused())
}

func TestEntryWindowFocus(t *testing.T) {
	entry := NewEntry()

	test.Canvas().Focus(entry)
	assert.True(t, entry.Focused())
}

func TestEntry_Tapped(t *testing.T) {
	entry := NewEntry()
	entry.SetText("MMM")

	test.Tap(entry)
	assert.True(t, entry.Focused())

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(int(float32(testCharSize)*1.5), testCharSize/2) // tap in the middle of the 2nd "M"
	ev := &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)

	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	pos = fyne.NewPos(int(float32(testCharSize)*2.5), testCharSize/2) // tap in the middle of the 3rd "M"
	ev = &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)

	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)
}

func TestEntry_Tapped_AfterCol(t *testing.T) {
	entry := NewEntry()
	entry.SetText("M")

	test.Tap(entry)
	assert.True(t, entry.Focused())

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize*2, testCharSize/2) // tap after text
	ev := &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)

	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_Tapped_AfterRow(t *testing.T) {
	entry := NewEntry()
	entry.SetText("M\nM\n")

	test.Tap(entry)
	assert.True(t, entry.Focused())

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev := &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)

	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func getClickPosition(e *Entry, str string, row int) *fyne.PointEvent {
	x := textMinSize(str, theme.TextSize(), e.textStyle()).Width + theme.Padding()

	rowHeight := e.textProvider().charMinSize().Height
	y := theme.Padding() + row*rowHeight + rowHeight/2

	pos := fyne.NewPos(x, y)
	return &fyne.PointEvent{Position: pos}
}

func TestEntry_ExpandSelectionForDoubleTap(t *testing.T) {
	str := []rune(" fish 日本語日  \t  test 本日本 moose  \t")

	// select invalid (before start)
	start, end := getTextWhitespaceRegion(str, -1)
	assert.Equal(t, -1, start)
	assert.Equal(t, -1, end)

	// select invalid (after end)
	start, end = getTextWhitespaceRegion(str, len(str))
	assert.Equal(t, -1, start)
	assert.Equal(t, -1, end)

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

func TestEntry_DoubleTapped(t *testing.T) {
	entry := NewEntry()
	entry.SetText("The quick brown fox\njumped    over the lazy dog\n")

	// select the word 'quick'
	ev := getClickPosition(entry, "The qui", 0)
	entry.Tapped(ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "quick", entry.selectedText())

	// select the whitespace after 'quick'
	ev = getClickPosition(entry, "The quick", 0)
	// add half a ' ' character
	ev.Position.X += textMinSize(" ", theme.TextSize(), entry.textStyle()).Width / 2
	entry.Tapped(ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, " ", entry.selectedText())

	// select all whitespace after 'jumped'
	ev = getClickPosition(entry, "jumped  ", 1)
	entry.Tapped(ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "    ", entry.selectedText())
}

func TestEntry_DoubleTapped_AfterCol(t *testing.T) {
	entry := NewEntry()
	entry.SetText("A\nB\n")

	test.Tap(entry)
	assert.True(t, entry.Focused())

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev := &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)
	entry.DoubleTapped(ev)

	assert.Equal(t, "", entry.selectedText())
}

func TestEntry_CursorRow(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("test")
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 0, entry.CursorRow)

	// 2 lines, this should increment
	entry.SetText("test\nrows")
	entry.TypedKey(down)
	assert.Equal(t, 1, entry.CursorRow)

	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_CursorColumn(t *testing.T) {
	entry := NewEntry()
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	// only 0 columns, do nothing
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorColumn)

	// 1, this should increment
	entry.SetText("a")
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorColumn)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorColumn)

	// don't go beyond left
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_CursorColumn_Wrap(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("a\nb")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// go to end of line
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	// wrap to new line
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// and back
	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorColumn_Jump(t *testing.T) {
	entry := NewMultiLineEntry()
	entry.SetText("a\nbc")

	// go to end of text
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)
	entry.TypedKey(right)
	entry.TypedKey(right)
	entry.TypedKey(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	// go up, to a shorter line
	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func checkNewlineIgnored(t *testing.T, entry *Entry) {
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)
	assert.Equal(t, 0, entry.CursorRow)

	// return is ignored, do nothing
	ret := &fyne.KeyEvent{Name: fyne.KeyReturn}
	entry.TypedKey(ret)
	assert.Equal(t, 0, entry.CursorRow)

	up := &fyne.KeyEvent{Name: fyne.KeyUp}
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.TypedKey(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestSingleLineEntry_NewlineIgnored(t *testing.T) {
	entry := &Entry{MultiLine: false}
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_NewlineIgnored(t *testing.T) {
	entry := NewPasswordEntry()
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_Obfuscation(t *testing.T) {
	entry := NewPasswordEntry()

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	assert.Equal(t, "*******", entryRenderTexts(entry)[0].Text)
}

func TestEntry_OnCut(t *testing.T) {
	e := NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCut{Clipboard: clipboard}
	handled := e.TypedShortcut(shortcut)

	assert.True(t, handled)
	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Teng", e.Text)
}

func TestEntry_OnCopy(t *testing.T) {
	e := NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCopy{Clipboard: clipboard}
	handled := e.TypedShortcut(shortcut)

	assert.True(t, handled)
	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnPaste(t *testing.T) {
	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	tests := []struct {
		name             string
		entry            *Entry
		clipboardContent string
		wantText         string
		wantRow, wantCol int
	}{
		{
			name:             "singleline: empty content",
			entry:            NewEntry(),
			clipboardContent: "",
			wantText:         "",
			wantRow:          0,
			wantCol:          0,
		},
		{
			name:             "singleline: simple text",
			entry:            NewEntry(),
			clipboardContent: "clipboard content",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: UTF8 text",
			entry:            NewEntry(),
			clipboardContent: "Hié™שרה",
			wantText:         "Hié™שרה",
			wantRow:          0,
			wantCol:          7,
		},
		{
			name:             "singleline: with new line",
			entry:            NewEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: with tab",
			entry:            NewEntry(),
			clipboardContent: "clipboard\tcontent",
			wantText:         "clipboard\tcontent",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "password: with new line",
			entry:            NewPasswordEntry(),
			clipboardContent: "3SB=y+)z\nkHGK(hx6 -e_\"1TZu q^bF3^$u H[:e\"1O.",
			wantText:         `3SB=y+)z kHGK(hx6 -e_"1TZu q^bF3^$u H[:e"1O.`,
			wantRow:          0,
			wantCol:          44,
		},
		{
			name:             "multiline: with new line",
			entry:            NewMultiLineEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard\ncontent",
			wantRow:          1,
			wantCol:          7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipboard.SetContent(tt.clipboardContent)
			handled := tt.entry.TypedShortcut(shortcut)
			assert.True(t, handled)
			assert.Equal(t, tt.wantText, tt.entry.Text)
			assert.Equal(t, tt.wantRow, tt.entry.CursorRow)
			assert.Equal(t, tt.wantCol, tt.entry.CursorColumn)
		})
	}
}

func TestEntry_PasteOverSelection(t *testing.T) {
	e := NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	clipboard.SetContent("Insert")
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	handled := e.TypedShortcut(shortcut)

	assert.True(t, handled)
	assert.Equal(t, "Insert", clipboard.Content())
	assert.Equal(t, "TeInsertng", e.Text)
}

func TestPasswordEntry_Placeholder(t *testing.T) {
	entry := NewPasswordEntry()
	entry.SetPlaceHolder("Password")

	assert.Equal(t, "Password", entryRenderPlaceholderTexts(entry)[0].Text)
	assert.False(t, entry.placeholderProvider().presenter.password())
}

const (
	keyShiftLeftDown  fyne.KeyName = "LeftShiftDown"
	keyShiftLeftUp    fyne.KeyName = "LeftShiftUp"
	keyShiftRightDown fyne.KeyName = "RightShiftDown"
	keyShiftRightUp   fyne.KeyName = "RightShiftUp"
)

var typeKeys = func(e *Entry, keys ...fyne.KeyName) {
	var keyDown = func(key *fyne.KeyEvent) {
		e.KeyDown(key)
		e.TypedKey(key)
	}

	for _, key := range keys {
		switch key {
		case keyShiftLeftDown:
			keyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
		case keyShiftLeftUp:
			e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
		case keyShiftRightDown:
			keyDown(&fyne.KeyEvent{Name: desktop.KeyShiftRight})
		case keyShiftRightUp:
			e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftRight})
		default:
			keyDown(&fyne.KeyEvent{Name: key})
			e.KeyUp(&fyne.KeyEvent{Name: key})
		}
	}
}

func TestEntry_SweetSweetCoverage(t *testing.T) {
	e := NewEntry()
	row, col := e.rowColFromTextPos(1)
	assert.Equal(t, 0, row)
	assert.Equal(t, 0, col)
}

func TestEntry_BasicSelect(t *testing.T) {

	// SeletionStart/SelectionEnd documentation example
	r := NewEntry()
	r.SetText("Testing")
	typeKeys(r, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight)
	a, b := r.selection()
	assert.Equal(t, 3, a)
	assert.Equal(t, 5, b)
	assert.Equal(t, "ti", r.selectedText())

	e := NewEntry()
	e.SetText("Testing")

	// move right, press & hold shift and move right
	typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight)
	a, b = e.selection()
	assert.Equal(t, 1, a)
	assert.Equal(t, 3, b)
	assert.Equal(t, "es", e.selectedText())

	// release shift
	typeKeys(e, keyShiftLeftUp)
	a, b = e.selection()
	assert.Equal(t, 1, a)
	assert.Equal(t, 3, b)

	// press shift and move
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight)
	a, b = e.selection()
	assert.Equal(t, 1, a)
	assert.Equal(t, 4, b)

	// release shift and move right
	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
	assert.Equal(t, "", e.selectedText())

	// press shift and move left
	e.CursorColumn = 4 // we should be here already thanks to snapping
	typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft)
	a, b = e.selection()
	assert.Equal(t, 2, a)
	assert.Equal(t, 4, b)
	assert.Equal(t, "st", e.selectedText())
}

// Selects "sti" on line 2 of a new multiline
// T e s t i n g
// T e[s t i]n g
// T e s t i n g
var setup = func() *Entry {
	e := NewMultiLineEntry()
	e.SetText("Testing\nTesting\nTesting")
	e.CursorRow = 1
	e.CursorColumn = 2
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)
	return e
}

// Selects "sti" on line 2 of a new multiline (but in reverse)
// T e s t i n g
// T e]s t i[n g
// T e s t i n g
var setupReverse = func() *Entry {
	e := NewMultiLineEntry()
	e.SetText("Testing\nTesting\nTesting")
	e.CursorRow = 1
	e.CursorColumn = 5
	typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft, fyne.KeyLeft)
	return e
}

func TestEntry_SelectHomeEnd(t *testing.T) {

	// T e[s t i] n g -> end -> // T e[s t i n g]
	e := setup()
	typeKeys(e, fyne.KeyEnd)
	a, b := e.selection()
	assert.Equal(t, 10, a)
	assert.Equal(t, 15, b)

	// T e s[t i n g] -> home -> ]T e[s t i n g
	typeKeys(e, fyne.KeyHome)
	a, b = e.selection()
	assert.Equal(t, 8, a)
	assert.Equal(t, 10, b)

	// home after releasing shift
	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyHome)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 0, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	// end after releasing shift
	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyEnd)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 7, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}

func TestEntry_MultilineSelect(t *testing.T) {
	e := setup()

	// Extend the selection down one row
	assert.Equal(t, 1, e.CursorRow)
	typeKeys(e, fyne.KeyDown)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	a, b := e.selection()
	assert.Equal(t, 10, a)
	assert.Equal(t, 21, b)

	typeKeys(e, fyne.KeyUp)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, 10, a)
	assert.Equal(t, 13, b)

	typeKeys(e, fyne.KeyUp)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, 5, a)
	assert.Equal(t, 10, b)
}

func TestEntry_SelectSnapping(t *testing.T) {

	e := setup()
	typeKeys(e, keyShiftLeftUp)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, fyne.KeyRight)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	// up and down snap to start/end respectively, but they also move
	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyDown)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyUp)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}

func TestEntry_SelectDelete(t *testing.T) {

	e := setup()
	typeKeys(e, fyne.KeyDelete)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	assert.Equal(t, 20, len(e.Text))
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	e = setup()
	typeKeys(e, fyne.KeyDown, fyne.KeyDelete)
	assert.Equal(t, "Testing\nTeng", e.Text)
	assert.Equal(t, 12, len(e.Text))

	e = setupReverse()
	typeKeys(e, fyne.KeyDown, fyne.KeyDelete)
	assert.Equal(t, "Testing\nTestisting", e.Text)
	assert.Equal(t, 18, len(e.Text))

	{
		// After pressing delete we should be able to press down to get a new selection
		// as we're still holding delete
		e = setup()
		typeKeys(e, fyne.KeyDelete, fyne.KeyDown)
		// T e s t i n g
		// T e[n g
		// T e]s t i n g
		a, b = e.selection()
		assert.Equal(t, 10, a)
		assert.Equal(t, 15, b)

		e = setupReverse()
		typeKeys(e, fyne.KeyDelete, fyne.KeyDown)
		a, b = e.selection()
		assert.Equal(t, 10, a)
		assert.Equal(t, 15, b)
	}

	{
		// Pressing up after delete should
		//  a) delete the selection
		//  b) move the selection start point
		e = setup()
		typeKeys(e, fyne.KeyDelete, fyne.KeyUp)
		// T e[s t i n g
		// T e]n g
		// T e s t i n g
		a, b = e.selection()
		assert.Equal(t, 2, a)
		assert.Equal(t, 10, b)

		e = setupReverse()
		typeKeys(e, fyne.KeyDelete, fyne.KeyUp)
		a, b = e.selection()
		assert.Equal(t, 2, a)
		assert.Equal(t, 10, b)
	}
}

func TestEntry_SelectBackspace(t *testing.T) {

	// AFAIK the backspace on selection behaviour should be identical to delete
	e := setup()
	typeKeys(e, fyne.KeyBackspace)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	assert.Equal(t, 20, len(e.Text))
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}

func TestEntry_SelectEnter(t *testing.T) {

	// Erase the selection and add a newline at selection start
	e := setup()
	typeKeys(e, fyne.KeyEnter)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTe\nng\nTesting", e.Text)
	assert.Equal(t, 21, len(e.Text))
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	e = setupReverse()
	typeKeys(e, fyne.KeyEnter)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTe\nng\nTesting", e.Text)
	assert.Equal(t, 21, len(e.Text))
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}

func TestEntry_SelectReplace(t *testing.T) {
	e := setup()
	test.Type(e, "hello")
	assert.Equal(t, "Testing\nTehellong\nTesting", e.Text)

	e = setupReverse()
	test.Type(e, "hello")
	assert.Equal(t, "Testing\nTehellong\nTesting", e.Text)
}

func TestEntry_EraseSelection(t *testing.T) {
	e := setup()
	e.eraseSelection()
	e.updateText(e.textProvider().String())
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}

func TestEntry_EraseEmptySelection(t *testing.T) {
	e := setup()
	// clear empty selection
	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft, fyne.KeyDelete)
	assert.Equal(t, "Testing\nTeting\nTesting", e.Text)
	a, b := e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	e = setup()
	// clear empty selection while shift is held
	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyDelete)
	assert.Equal(t, "Testing\nTeting\nTesting", e.Text)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	// ensure that backspace doesn't leave a selection start at the old cursor position
	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyBackspace)
	assert.Equal(t, "Testing\nTsting\nTesting", e.Text)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)

	// clear selection, select a character and while holding shift issue two backspaces
	e = setup()
	typeKeys(e, keyShiftLeftUp, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyBackspace, fyne.KeyBackspace)
	assert.Equal(t, "Testing\nTeing\nTesting", e.Text)
	a, b = e.selection()
	assert.Equal(t, -1, a)
	assert.Equal(t, -1, b)
}
