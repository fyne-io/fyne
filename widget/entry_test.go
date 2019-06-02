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

func TestPasswordEntry_Placeholder(t *testing.T) {
	entry := NewPasswordEntry()
	entry.SetPlaceHolder("Password")

	assert.Equal(t, "Password", entryRenderPlaceholderTexts(entry)[0].Text)
	assert.False(t, entry.placeholderProvider().presenter.password())
}

var tap = func(e *Entry, k *fyne.KeyEvent) {
	e.KeyDown(k)
	e.TypedKey(k)
	e.KeyUp(k)
}

var right = func(e *Entry) {
	tap(e, &fyne.KeyEvent{Name: fyne.KeyRight})
}
var left = func(e *Entry) {
	tap(e, &fyne.KeyEvent{Name: fyne.KeyLeft})
}
var up = func(e *Entry) {
	tap(e, &fyne.KeyEvent{Name: fyne.KeyUp})
}
var down = func(e *Entry) {
	tap(e, &fyne.KeyEvent{Name: fyne.KeyDown})
}

var shiftDown = func(e *Entry) {
	k := &fyne.KeyEvent{Name: desktop.KeyShiftLeft}
	e.KeyDown(k)
	e.TypedKey(k)
}

var shiftUp = func(e *Entry) {
	k := &fyne.KeyEvent{Name: desktop.KeyShiftLeft}
	e.KeyUp(k)
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
	right(r)
	right(r)
	right(r)
	shiftDown(r)
	right(r)
	right(r)
	assert.Equal(t, 3, r.SelectionStart())
	assert.Equal(t, 5, r.SelectionEnd())

	e := NewEntry()
	e.SetText("Testing")

	// move right, press & hold shift and move right
	right(e)
	shiftDown(e)
	right(e)
	right(e)
	assert.Equal(t, 1, e.SelectionStart())
	assert.Equal(t, 3, e.SelectionEnd())

	// release shift
	shiftUp(e)
	assert.Equal(t, 1, e.SelectionStart())
	assert.Equal(t, 3, e.SelectionEnd())

	// press shift and move
	shiftDown(e)
	right(e)
	assert.Equal(t, 1, e.SelectionStart())
	assert.Equal(t, 4, e.SelectionEnd())

	// release shift and move right
	shiftUp(e)
	right(e)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())

	// press shift and move left
	e.CursorColumn = 4 // we should be here already thanks to snapping
	shiftDown(e)
	left(e)
	left(e)
	assert.Equal(t, 2, e.SelectionStart())
	assert.Equal(t, 4, e.SelectionEnd())
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
	shiftDown(e)
	right(e)
	right(e)
	right(e)
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
	shiftDown(e)
	left(e)
	left(e)
	left(e)
	return e
}

func TestEntry_SelectHomeEnd(t *testing.T) {
	home := &fyne.KeyEvent{Name: fyne.KeyHome}
	end := &fyne.KeyEvent{Name: fyne.KeyEnd}

	// T e[s t i] n g -> end -> // T e[s t i n g]
	e := setup()
	tap(e, end)
	assert.Equal(t, 10, e.SelectionStart())
	assert.Equal(t, 15, e.SelectionEnd())

	// T e s[t i n g] -> home -> ]T e[s t i n g
	tap(e, home)
	assert.Equal(t, 8, e.SelectionStart())
	assert.Equal(t, 10, e.SelectionEnd())

	// home after releasing shift
	e = setup()
	shiftUp(e)
	tap(e, home)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 0, e.CursorColumn)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())

	// end after releasing shift
	e = setup()
	shiftUp(e)
	tap(e, end)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 7, e.CursorColumn)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())
}

func TestEntry_MultilineSelect(t *testing.T) {
	e := setup()

	// Extend the selection down one row
	assert.Equal(t, 1, e.CursorRow)
	down(e)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	assert.Equal(t, 10, e.SelectionStart())
	assert.Equal(t, 21, e.SelectionEnd())

	up(e)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	assert.Equal(t, 10, e.SelectionStart())
	assert.Equal(t, 13, e.SelectionEnd())

	up(e)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	assert.Equal(t, 5, e.SelectionStart())
	assert.Equal(t, 10, e.SelectionEnd())
}

func TestEntry_SelectSnapping(t *testing.T) {

	e := setup()
	shiftUp(e)

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	right(e)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())

	e = setup()
	shiftUp(e)
	left(e)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())

	// up and down snap to start/end respectively, but they also move
	e = setup()
	shiftUp(e)
	down(e)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())

	e = setup()
	shiftUp(e)
	up(e)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())
}

func TestEntry_SelectDelete(t *testing.T) {
	del := &fyne.KeyEvent{Name: fyne.KeyDelete}

	e := setup()
	tap(e, del)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	assert.Equal(t, 20, len(e.Text))
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())

	e = setup()
	down(e)
	tap(e, del)
	assert.Equal(t, "Testing\nTeng", e.Text)
	assert.Equal(t, 12, len(e.Text))

	e = setupReverse()
	down(e)
	tap(e, del)
	assert.Equal(t, "Testing\nTestisting", e.Text)
	assert.Equal(t, 18, len(e.Text))

	{
		// After pressing delete we should be able to press down to get a new selection
		// as we're still holding delete
		e = setup()
		tap(e, del)
		down(e)
		// T e s t i n g
		// T e[n g
		// T e]s t i n g
		assert.Equal(t, 10, e.SelectionStart())
		assert.Equal(t, 15, e.SelectionEnd())

		e = setupReverse()
		tap(e, del)
		down(e)
		assert.Equal(t, 10, e.SelectionStart())
		assert.Equal(t, 15, e.SelectionEnd())
	}

	{
		// Pressing up after delete should
		//  a) delete the selection
		//  b) move the selection start point
		e = setup()
		tap(e, del)
		up(e)
		// T e[s t i n g
		// T e]n g
		// T e s t i n g
		assert.Equal(t, 2, e.SelectionStart())
		assert.Equal(t, 10, e.SelectionEnd())

		e = setupReverse()
		tap(e, del)
		up(e)
		assert.Equal(t, 2, e.SelectionStart())
		assert.Equal(t, 10, e.SelectionEnd())
	}
}

func TestEntry_SelectBackspace(t *testing.T) {

	// AFAIK the backspace on selection behaviour should be identical to delete
	bs := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	e := setup()
	tap(e, bs)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTeng\nTesting", e.Text)
	assert.Equal(t, 20, len(e.Text))
	assert.Equal(t, -1, e.SelectionStart())
	assert.Equal(t, -1, e.SelectionEnd())
}

func TestEntry_SelectEnter(t *testing.T) {

	// Erase the selection and add a newline at selection start
	bs := &fyne.KeyEvent{Name: fyne.KeyEnter}
	e := setup()
	tap(e, bs)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTe\nng\nTesting", e.Text)
	assert.Equal(t, 21, len(e.Text))
	assert.Equal(t, 10, e.SelectionStart()) // Hmm, maybe these should be -1 and -1
	assert.Equal(t, 11, e.SelectionEnd())

	e = setupReverse()
	tap(e, bs)
	// "Testing\nTeng\nTesting"
	assert.Equal(t, "Testing\nTe\nng\nTesting", e.Text)
	assert.Equal(t, 21, len(e.Text))
	assert.Equal(t, 10, e.SelectionStart()) // Hmm, maybe these should be -1 and -1
	assert.Equal(t, 11, e.SelectionEnd())
}
