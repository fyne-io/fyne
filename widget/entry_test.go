package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Cursor(t *testing.T) {
	entry := widget.NewEntry()
	assert.Equal(t, desktop.TextCursor, entry.Cursor())
}

func TestEntry_CursorColumn(t *testing.T) {
	entry := widget.NewEntry()
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

func TestEntry_CursorColumn_Jump(t *testing.T) {
	entry := widget.NewMultiLineEntry()
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

func TestEntry_CursorColumn_Wrap(t *testing.T) {
	entry := widget.NewMultiLineEntry()
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

func TestEntry_CursorPasswordRevealer(t *testing.T) {
	pr := widget.NewPasswordEntry().ActionItem.(desktop.Cursorable)
	assert.Equal(t, desktop.DefaultCursor, pr.Cursor())
}

func TestEntry_CursorRow(t *testing.T) {
	entry := widget.NewMultiLineEntry()
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

func TestEntry_Disableable(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.False(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_enabled_empty.png", c.Capture())

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_disabled_empty.png", c.Capture())

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_enabled_empty.png", c.Capture())

	entry.SetPlaceHolder("Type!")
	assert.False(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_enabled_placeholder.png", c.Capture())

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_disabled_placeholder.png", c.Capture())

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_enabled_placeholder.png", c.Capture())

	entry.SetText("Hello")
	assert.False(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_enabled_custom_value.png", c.Capture())

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_disabled_custom_value.png", c.Capture())

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertImageMatches(t, "entry/disableable_enabled_custom_value.png", c.Capture())
}

func TestEntry_DoubleTapped(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("The quick brown fox\njumped    over the lazy dog\n")

	// select the word 'quick'
	ev := getClickPosition("The qui", 0)
	entry.Tapped(ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "quick", entry.SelectedText())

	// select the whitespace after 'quick'
	ev = getClickPosition("The quick", 0)
	// add half a ' ' character
	ev.Position.X += fyne.MeasureText(" ", theme.TextSize(), fyne.TextStyle{}).Width / 2
	entry.Tapped(ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, " ", entry.SelectedText())

	// select all whitespace after 'jumped'
	ev = getClickPosition("jumped  ", 1)
	entry.Tapped(ev)
	entry.DoubleTapped(ev)
	assert.Equal(t, "    ", entry.SelectedText())
}

func TestEntry_DoubleTapped_AfterCol(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("A\nB\n")
	c := fyne.CurrentApp().Driver().CanvasForObject(entry)

	test.Tap(entry)
	assert.Equal(t, c.Focused(), entry)

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev := &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)
	entry.DoubleTapped(ev)

	assert.Equal(t, "", entry.SelectedText())
}

func TestEntry_DragSelect(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("The quick brown fox jumped\nover the lazy dog\nThe quick\nbrown fox\njumped over the lazy dog\n")

	// get position after the letter 'e' on the second row
	ev1 := getClickPosition("ove", 1)
	// get position after the letter 'z' on the second row
	ev2 := getClickPosition("over the laz", 1)
	// add a couple of pixels, this is currently a workaround for weird mouse to column logic on text with kerning
	ev2.Position.X += 2

	// mouse down and drag from 'r' to 'z'
	me := &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.LeftMouseButton}
	entry.MouseDown(me)
	for ; ev1.Position.X < ev2.Position.X; ev1.Position.X++ {
		de := &fyne.DragEvent{PointEvent: *ev1, DraggedX: 1, DraggedY: 0}
		entry.Dragged(de)
	}
	me = &desktop.MouseEvent{PointEvent: *ev1, Button: desktop.LeftMouseButton}
	entry.MouseUp(me)

	assert.Equal(t, "r the laz", entry.SelectedText())
}

func TestEntry_EmptySelection(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("text")

	// trying to select at the edge
	typeKeys(entry, keyShiftLeftDown, fyne.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())

	typeKeys(entry, fyne.KeyRight)
	assert.Equal(t, 1, entry.CursorColumn)

	// stop selecting at the edge when nothing is selected
	typeKeys(entry, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())
	assert.Equal(t, 0, entry.CursorColumn)

	// check that the selection has been removed
	typeKeys(entry, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftUp)
	assert.Equal(t, "", entry.SelectedText())
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_Focus(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/focus_initial.png", c.Capture())

	entry.FocusGained()
	test.AssertImageMatches(t, "entry/focus_focus_gained.png", c.Capture())

	entry.FocusLost()
	test.AssertImageMatches(t, "entry/focus_focus_lost.png", c.Capture())

	test.Canvas().Focus(entry)
	test.AssertImageMatches(t, "entry/focus_focus_gained.png", c.Capture())
}

func TestEntry_FocusWithPopUp(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.TapSecondaryAt(entry, fyne.NewPos(1, 1))
	test.AssertImageMatches(t, "entry/focus_with_popup_initial.png", c.Capture())

	test.TapCanvas(c, fyne.NewPos(20, 20))
	test.AssertImageMatches(t, "entry/focus_with_popup_entry_selected.png", c.Capture())

	test.TapSecondaryAt(entry, fyne.NewPos(1, 1))
	test.AssertImageMatches(t, "entry/focus_with_popup_initial.png", c.Capture())

	test.TapCanvas(c, fyne.NewPos(5, 5))
	test.AssertImageMatches(t, "entry/focus_with_popup_dismissed.png", c.Capture())
}

func TestEntry_HidePopUpOnEntry(t *testing.T) {
	entry := widget.NewEntry()
	tapPos := fyne.NewPos(1, 1)
	c := fyne.CurrentApp().Driver().CanvasForObject(entry)

	assert.Nil(t, c.Overlays().Top())

	test.TapSecondaryAt(entry, tapPos)
	assert.NotNil(t, c.Overlays().Top())

	test.Type(entry, "KJGFD")
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "KJGFD", entry.Text)
}

func TestEntry_MinSize(t *testing.T) {
	entry := widget.NewEntry()
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

	min = entry.MinSize()
	entry.ActionItem = canvas.NewCircle(color.Black)
	assert.Equal(t, min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0)), entry.MinSize())
}

func TestEntry_MouseDownOnSelect(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Ahnj\nBuki\n")
	entry.TypedShortcut(&fyne.ShortcutSelectAll{})

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev := &fyne.PointEvent{Position: pos}

	me := &desktop.MouseEvent{PointEvent: *ev, Button: desktop.RightMouseButton}
	entry.MouseDown(me)
	entry.MouseUp(me)

	assert.Equal(t, entry.SelectedText(), "Ahnj\nBuki\n")

	me = &desktop.MouseEvent{PointEvent: *ev, Button: desktop.LeftMouseButton}
	entry.MouseDown(me)
	entry.MouseUp(me)

	assert.Equal(t, entry.SelectedText(), "")
}

func TestEntry_MultilineSelect(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	// Extend the selection down one row
	typeKeys(e, fyne.KeyDown)
	test.AssertImageMatches(t, "entry/selection_add_one_row_down.png", c.Capture())

	typeKeys(e, fyne.KeyUp)
	test.AssertImageMatches(t, "entry/selection_remove_one_row_up.png", c.Capture())

	typeKeys(e, fyne.KeyUp)
	test.AssertImageMatches(t, "entry/selection_remove_add_one_row_up.png", c.Capture())
}

func TestEntry_Notify(t *testing.T) {
	entry := widget.NewEntry()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntry_OnCopy(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCopy{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnCopy_Password(t *testing.T) {
	e := widget.NewPasswordEntry()
	e.SetText("Testing")
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCopy{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnCut(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCut{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "sti", clipboard.Content())
	assert.Equal(t, "Teng", e.Text)
}

func TestEntry_OnCut_Password(t *testing.T) {
	e := widget.NewPasswordEntry()
	e.SetText("Testing")
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutCut{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "", clipboard.Content())
	assert.Equal(t, "Testing", e.Text)
}

func TestEntry_OnKeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "Hi")

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace(t *testing.T) {
	entry := widget.NewEntry()
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
	entry := widget.NewEntry()
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

func TestEntry_OnKeyDown_BackspaceBeyondTextAndNewLine(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)
	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, "H", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceNewline(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceUnicode(t *testing.T) {
	entry := widget.NewEntry()

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
	entry := widget.NewEntry()
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
	entry := widget.NewEntry()
	entry.SetText("Hi")

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)
	entry.TypedKey(key)
	entry.TypedKey(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_DeleteNewline(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("H\ni")

	right := &fyne.KeyEvent{Name: fyne.KeyRight}
	entry.TypedKey(right)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_HomeEnd(t *testing.T) {
	entry := &widget.Entry{}
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

func TestEntry_OnKeyDown_Insert(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "Hi")
	assert.Equal(t, "Hi", entry.Text)

	left := &fyne.KeyEvent{Name: fyne.KeyLeft}
	entry.TypedKey(left)

	test.Type(entry, "o")
	assert.Equal(t, "Hoi", entry.Text)
}

func TestEntry_OnKeyDown_Newline(t *testing.T) {
	entry, window := setupImageTest(true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	test.AssertImageMatches(t, "entry/on_key_down_newline_initial.png", c.Capture())

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
	assert.Equal(t, "H\noi", entry.Text)
	test.AssertImageMatches(t, "entry/on_key_down_newline_typed.png", c.Capture())
}

func TestEntry_OnPaste(t *testing.T) {
	clipboard := test.NewClipboard()
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	tests := []struct {
		name             string
		entry            *widget.Entry
		clipboardContent string
		wantText         string
		wantRow, wantCol int
	}{
		{
			name:             "singleline: empty content",
			entry:            widget.NewEntry(),
			clipboardContent: "",
			wantText:         "",
			wantRow:          0,
			wantCol:          0,
		},
		{
			name:             "singleline: simple text",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard content",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: UTF8 text",
			entry:            widget.NewEntry(),
			clipboardContent: "Hié™שרה",
			wantText:         "Hié™שרה",
			wantRow:          0,
			wantCol:          7,
		},
		{
			name:             "singleline: with new line",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard content",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "singleline: with tab",
			entry:            widget.NewEntry(),
			clipboardContent: "clipboard\tcontent",
			wantText:         "clipboard\tcontent",
			wantRow:          0,
			wantCol:          17,
		},
		{
			name:             "password: with new line",
			entry:            widget.NewPasswordEntry(),
			clipboardContent: "3SB=y+)z\nkHGK(hx6 -e_\"1TZu q^bF3^$u H[:e\"1O.",
			wantText:         `3SB=y+)z kHGK(hx6 -e_"1TZu q^bF3^$u H[:e"1O.`,
			wantRow:          0,
			wantCol:          44,
		},
		{
			name:             "multiline: with new line",
			entry:            widget.NewMultiLineEntry(),
			clipboardContent: "clipboard\ncontent",
			wantText:         "clipboard\ncontent",
			wantRow:          1,
			wantCol:          7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipboard.SetContent(tt.clipboardContent)
			tt.entry.TypedShortcut(shortcut)
			assert.Equal(t, tt.wantText, tt.entry.Text)
			assert.Equal(t, tt.wantRow, tt.entry.CursorRow)
			assert.Equal(t, tt.wantCol, tt.entry.CursorColumn)
		})
	}
}

func TestEntry_PageUpDown(t *testing.T) {
	t.Run("single line", func(*testing.T) {
		e, window := setupImageTest(false)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing")
		test.AssertImageMatches(t, "entry/select_initial.png", c.Capture())

		// move right, press & hold shift and pagedown
		typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyPageDown)
		assert.Equal(t, "esting", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertImageMatches(t, "entry/select_single_line_shift_pagedown.png", c.Capture())

		// while shift is held press pageup
		typeKeys(e, fyne.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertImageMatches(t, "entry/select_single_line_shift_pageup.png", c.Capture())

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, fyne.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertImageMatches(t, "entry/select_single_line_pagedown.png", c.Capture())
	})

	t.Run("page down single line", func(*testing.T) {
		e, window := setupImageTest(true)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing\nTesting\nTesting")
		test.AssertImageMatches(t, "entry/select_multi_line_initial.png", c.Capture())

		// move right, press & hold shift and pagedown
		typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyPageDown)
		assert.Equal(t, "esting\nTesting\nTesting", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertImageMatches(t, "entry/select_multi_line_shift_pagedown.png", c.Capture())

		// while shift is held press pageup
		typeKeys(e, fyne.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertImageMatches(t, "entry/select_multi_line_shift_pageup.png", c.Capture())

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, fyne.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertImageMatches(t, "entry/select_multi_line_pagedown.png", c.Capture())
	})
}

func TestEntry_PasteOverSelection(t *testing.T) {
	e := widget.NewEntry()
	e.SetText("Testing")
	typeKeys(e, fyne.KeyRight, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)

	clipboard := test.NewClipboard()
	clipboard.SetContent("Insert")
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "Insert", clipboard.Content())
	assert.Equal(t, "TeInsertng", e.Text)
}

func TestEntry_PasteUnicode(t *testing.T) {
	e := widget.NewMultiLineEntry()
	e.SetText("line")
	e.CursorColumn = 4

	clipboard := test.NewClipboard()
	clipboard.SetContent("thing {\n\titem: 'val测试'\n}")
	shortcut := &fyne.ShortcutPaste{Clipboard: clipboard}
	e.TypedShortcut(shortcut)

	assert.Equal(t, "thing {\n\titem: 'val测试'\n}", clipboard.Content())
	assert.Equal(t, "linething {\n\titem: 'val测试'\n}", e.Text)

	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 1, e.CursorColumn)
}

func TestEntry_Placeholder(t *testing.T) {
	entry := &widget.Entry{}
	entry.Text = "Text"
	entry.PlaceHolder = "Placehold"

	window := test.NewWindow(entry)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, "Text", entry.Text)
	test.AssertImageMatches(t, "entry/placeholder_withtext.png", c.Capture())

	entry.SetText("")
	assert.Equal(t, "", entry.Text)
	test.AssertImageMatches(t, "entry/placeholder_initial.png", c.Capture())
}

func TestEntry_Select(t *testing.T) {
	for name, tt := range map[string]struct {
		keys          []fyne.KeyName
		text          string
		setupReverse  bool
		wantImage     string
		wantSelection string
		wantText      string
	}{
		"delete single-line": {
			keys:      []fyne.KeyName{fyne.KeyDelete},
			wantText:  "Testing\nTeng\nTesting",
			wantImage: "entry/selection_delete_single_line.png",
		},
		"delete multi-line": {
			keys:      []fyne.KeyName{fyne.KeyDown, fyne.KeyDelete},
			wantText:  "Testing\nTeng",
			wantImage: "entry/selection_delete_multi_line.png",
		},
		"delete reverse multi-line": {
			keys:         []fyne.KeyName{fyne.KeyDown, fyne.KeyDelete},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantImage:    "entry/selection_delete_reverse_multi_line.png",
		},
		"delete select down with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyDelete, fyne.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantImage:     "entry/selection_delete_and_add_down.png",
		},
		"delete reverse select down with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyDelete, fyne.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantImage:     "entry/selection_delete_and_add_down.png",
		},
		"delete select up with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyDelete, fyne.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantImage:     "entry/selection_delete_and_add_up.png",
		},
		"delete reverse select up with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyDelete, fyne.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantImage:     "entry/selection_delete_and_add_up.png",
		},
		// The backspace delete behaviour is the same as via delete.
		"backspace single-line": {
			keys:      []fyne.KeyName{fyne.KeyBackspace},
			wantText:  "Testing\nTeng\nTesting",
			wantImage: "entry/selection_delete_single_line.png",
		},
		"backspace multi-line": {
			keys:      []fyne.KeyName{fyne.KeyDown, fyne.KeyBackspace},
			wantText:  "Testing\nTeng",
			wantImage: "entry/selection_delete_multi_line.png",
		},
		"backspace reverse multi-line": {
			keys:         []fyne.KeyName{fyne.KeyDown, fyne.KeyBackspace},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantImage:    "entry/selection_delete_reverse_multi_line.png",
		},
		"backspace select down with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyBackspace, fyne.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantImage:     "entry/selection_delete_and_add_down.png",
		},
		"backspace reverse select down with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyBackspace, fyne.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantImage:     "entry/selection_delete_and_add_down.png",
		},
		"backspace select up with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyBackspace, fyne.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantImage:     "entry/selection_delete_and_add_up.png",
		},
		"backspace reverse select up with Shift still hold": {
			keys:          []fyne.KeyName{fyne.KeyBackspace, fyne.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantImage:     "entry/selection_delete_and_add_up.png",
		},
		// Erase the selection and add a newline at selection start
		"enter": {
			keys:      []fyne.KeyName{fyne.KeyEnter},
			wantText:  "Testing\nTe\nng\nTesting",
			wantImage: "entry/selection_enter.png",
		},
		"enter reverse": {
			keys:         []fyne.KeyName{fyne.KeyEnter},
			setupReverse: true,
			wantText:     "Testing\nTe\nng\nTesting",
			wantImage:    "entry/selection_enter.png",
		},
		"replace": {
			text:      "hello",
			wantText:  "Testing\nTehellong\nTesting",
			wantImage: "entry/selection_replace.png",
		},
		"replace reverse": {
			text:         "hello",
			setupReverse: true,
			wantText:     "Testing\nTehellong\nTesting",
			wantImage:    "entry/selection_replace.png",
		},
		"deselect and delete": {
			keys:      []fyne.KeyName{keyShiftLeftUp, fyne.KeyLeft, fyne.KeyDelete},
			wantText:  "Testing\nTeting\nTesting",
			wantImage: "entry/selection_deselect_delete.png",
		},
		"deselect and delete holding shift": {
			keys:      []fyne.KeyName{keyShiftLeftUp, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyDelete},
			wantText:  "Testing\nTeting\nTesting",
			wantImage: "entry/selection_deselect_delete.png",
		},
		// ensure that backspace doesn't leave a selection start at the old cursor position
		"deselect and backspace holding shift": {
			keys:      []fyne.KeyName{keyShiftLeftUp, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyBackspace},
			wantText:  "Testing\nTsting\nTesting",
			wantImage: "entry/selection_deselect_backspace.png",
		},
		// clear selection, select a character and while holding shift issue two backspaces
		"deselect, select and double backspace": {
			keys:      []fyne.KeyName{keyShiftLeftUp, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyBackspace, fyne.KeyBackspace},
			wantText:  "Testing\nTeing\nTesting",
			wantImage: "entry/selection_deselect_select_backspace.png",
		},
	} {
		t.Run(name, func(t *testing.T) {
			entry, window := setupSelection(tt.setupReverse)
			defer teardownImageTest(window)
			c := window.Canvas()

			if tt.setupReverse {
				test.AssertImageMatches(t, "entry/selection_reverse_initial.png", c.Capture())
			} else {
				test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())
			}

			if tt.text != "" {
				test.Type(entry, tt.text)
			} else {
				typeKeys(entry, tt.keys...)
			}
			assert.Equal(t, tt.wantText, entry.Text)
			assert.Equal(t, tt.wantSelection, entry.SelectedText())
			test.AssertImageMatches(t, tt.wantImage, c.Capture())
		})
	}
}

func TestEntry_SelectAll(t *testing.T) {
	e, window := setupImageTest(true)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("First Row\nSecond Row\nThird Row")
	test.AssertImageMatches(t, "entry/select_all_initial.png", c.Capture())

	shortcut := &fyne.ShortcutSelectAll{}
	e.TypedShortcut(shortcut)
	test.AssertImageMatches(t, "entry/select_all_selected.png", c.Capture())
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 9, e.CursorColumn)
}

func TestEntry_SelectAll_EmptyEntry(t *testing.T) {
	entry := widget.NewEntry()
	entry.TypedShortcut(&fyne.ShortcutSelectAll{})

	assert.Equal(t, "", entry.SelectedText())
}

func TestEntry_SelectEndWithoutShift(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	// end after releasing shift
	typeKeys(e, keyShiftLeftUp, fyne.KeyEnd)
	test.AssertImageMatches(t, "entry/selection_end.png", c.Capture())
}

func TestEntry_SelectHomeEnd(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	// T e[s t i] n g -> end -> // T e[s t i n g]
	typeKeys(e, fyne.KeyEnd)
	test.AssertImageMatches(t, "entry/selection_add_to_end.png", c.Capture())

	// T e s[t i n g] -> home -> ]T e[s t i n g
	typeKeys(e, fyne.KeyHome)
	test.AssertImageMatches(t, "entry/selection_add_to_home.png", c.Capture())
}

func TestEntry_SelectHomeWithoutShift(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	// home after releasing shift
	typeKeys(e, keyShiftLeftUp, fyne.KeyHome)
	test.AssertImageMatches(t, "entry/selection_home.png", c.Capture())
}

func TestEntry_SelectSnapDown(t *testing.T) {
	// down snaps to end, but it also moves
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	typeKeys(e, keyShiftLeftUp, fyne.KeyDown)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_snap_down.png", c.Capture())
}

func TestEntry_SelectSnapLeft(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_snap_left.png", c.Capture())
}

func TestEntry_SelectSnapRight(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_snap_right.png", c.Capture())
}

func TestEntry_SelectSnapUp(t *testing.T) {
	// up snaps to start, but it also moves
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	typeKeys(e, keyShiftLeftUp, fyne.KeyUp)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertImageMatches(t, "entry/selection_snap_up.png", c.Capture())
}

func TestEntry_SelectedText(t *testing.T) {
	e, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("Testing")
	test.AssertImageMatches(t, "entry/select_initial.png", c.Capture())

	// move right, press & hold shift and move right
	typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight)
	assert.Equal(t, "es", e.SelectedText())
	test.AssertImageMatches(t, "entry/select_selected.png", c.Capture())

	// release shift
	typeKeys(e, keyShiftLeftUp)
	// press shift and move
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight)
	assert.Equal(t, "est", e.SelectedText())
	test.AssertImageMatches(t, "entry/select_add_selection.png", c.Capture())

	// release shift and move right
	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	assert.Equal(t, "", e.SelectedText())
	test.AssertImageMatches(t, "entry/select_move_wo_shift.png", c.Capture())

	// press shift and move left
	typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft)
	assert.Equal(t, "st", e.SelectedText())
	test.AssertImageMatches(t, "entry/select_select_left.png", c.Capture())
}

func TestEntry_SelectionHides(t *testing.T) {
	e, window := setupSelection(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/selection_initial.png", c.Capture())

	c.Unfocus()
	test.AssertImageMatches(t, "entry/selection_focus_lost.png", c.Capture())

	c.Focus(e)
	test.AssertImageMatches(t, "entry/selection_focus_gained.png", c.Capture())
}

func TestEntry_SetPlaceHolder(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 0, len(entry.Text))
	test.AssertImageMatches(t, "entry/set_placeholder_initial.png", c.Capture())

	entry.SetPlaceHolder("Test")
	assert.Equal(t, 0, len(entry.Text))
	test.AssertImageMatches(t, "entry/set_placeholder_set.png", c.Capture())

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	test.AssertImageMatches(t, "entry/set_placeholder_replaced.png", c.Capture())
}

func TestEntry_SetReadOnly_KeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "H")
	entry.SetReadOnly(true)
	test.Type(entry, "i")
	assert.Equal(t, "H", entry.Text)

	entry.SetReadOnly(false)
	test.Type(entry, "i")
	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_SetReadOnly_OnFocus(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetReadOnly(true)
	entry.FocusGained()
	test.AssertImageMatches(t, "entry/set_readonly_on_focus_readonly.png", c.Capture())

	entry.SetReadOnly(false)
	entry.FocusGained()
	test.AssertImageMatches(t, "entry/set_readonly_on_focus_writable.png", c.Capture())
}

func TestEntry_SetText_EmptyString(t *testing.T) {
	entry := widget.NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	entry = widget.NewMultiLineEntry()
	test.Type(entry, "test\ntest")

	down := &fyne.KeyEvent{Name: fyne.KeyDown}
	entry.TypedKey(down)

	assert.Equal(t, 4, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_SetText_Manual(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/set_text_initial.png", c.Capture())

	entry.Text = "Test"
	entry.Refresh()
	test.AssertImageMatches(t, "entry/set_text_changed.png", c.Capture())
}

func TestEntry_SetText_Overflow(t *testing.T) {
	entry := widget.NewEntry()

	assert.Equal(t, 0, entry.CursorColumn)

	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)

	entry.SetText("x")
	assert.Equal(t, 1, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyDelete}
	entry.TypedKey(key)

	assert.Equal(t, 1, entry.CursorColumn)
	assert.Equal(t, "x", entry.Text)

	key = &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, "", entry.Text)
}

func TestEntry_SetText_Underflow(t *testing.T) {
	entry := widget.NewEntry()
	test.Type(entry, "test")
	assert.Equal(t, 4, entry.CursorColumn)

	entry.Text = ""
	entry.Refresh()
	assert.Equal(t, 0, entry.CursorColumn)

	key := &fyne.KeyEvent{Name: fyne.KeyBackspace}
	entry.TypedKey(key)

	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, "", entry.Text)
}

func TestEntry_Tapped(t *testing.T) {
	entry, window := setupImageTest(true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("MMM\nWWW\n")
	test.AssertImageMatches(t, "entry/tapped_initial.png", c.Capture())

	test.Tap(entry)
	test.AssertImageMatches(t, "entry/tapped_focused.png", c.Capture())

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(int(float32(testCharSize)*1.5), testCharSize/2) // tap in the middle of the 2nd "M"
	ev := &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)
	test.AssertImageMatches(t, "entry/tapped_tapped_2nd_m.png", c.Capture())
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	pos = fyne.NewPos(int(float32(testCharSize)*2.5), testCharSize/2) // tap in the middle of the 3rd "M"
	ev = &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)
	test.AssertImageMatches(t, "entry/tapped_tapped_3nd_m.png", c.Capture())
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	pos = fyne.NewPos(testCharSize*4, testCharSize/2) // tap after text
	ev = &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)
	test.AssertImageMatches(t, "entry/tapped_tapped_after_last_col.png", c.Capture())
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 3, entry.CursorColumn)

	pos = fyne.NewPos(testCharSize, testCharSize*4) // tap below rows
	ev = &fyne.PointEvent{Position: pos}
	entry.Tapped(ev)
	test.AssertImageMatches(t, "entry/tapped_tapped_after_last_row.png", c.Capture())
	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_TappedSecondary(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "entry/tapped_secondary_initial.png", c.Capture())

	tapPos := fyne.NewPos(20, 10)
	test.TapSecondaryAt(entry, tapPos)
	test.AssertImageMatches(t, "entry/tapped_secondary_full_menu.png", c.Capture())
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Disable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertImageMatches(t, "entry/tapped_secondary_read_menu.png", c.Capture())
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Password = true
	entry.Refresh()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertImageMatches(t, "entry/tapped_secondary_no_password_menu.png", c.Capture())
	assert.Nil(t, c.Overlays().Top(), "No popup for disabled password")

	entry.Enable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertImageMatches(t, "entry/tapped_secondary_password_menu.png", c.Capture())
	assert.Equal(t, 1, len(c.Overlays().List()))
}

func TestEntry_TextWrap(t *testing.T) {
	for name, tt := range map[string]struct {
		multiLine bool
		wantImage string
		wrap      fyne.TextWrap
	}{
		"single line WrapOff": {
			wantImage: "entry/wrap_single_line_off.png",
		},
		// Disallowed - fallback to TextWrapOff
		"single line Truncate": {
			wrap:      fyne.TextTruncate,
			wantImage: "entry/wrap_single_line_off.png",
		},
		// Disallowed - fallback to TextWrapOff
		"single line WrapBreak": {
			wrap:      fyne.TextWrapBreak,
			wantImage: "entry/wrap_single_line_off.png",
		},
		// Disallowed - fallback to TextWrapOff
		"single line WrapWord": {
			wrap:      fyne.TextWrapWord,
			wantImage: "entry/wrap_single_line_off.png",
		},
		"multi line WrapOff": {
			multiLine: true,
			wantImage: "entry/wrap_multi_line_off.png",
		},
		// Disallowed - fallback to TextWrapOff
		"multi line Truncate": {
			multiLine: true,
			wrap:      fyne.TextTruncate,
			wantImage: "entry/wrap_multi_line_off.png",
		},
		"multi line WrapBreak": {
			multiLine: true,
			wrap:      fyne.TextWrapBreak,
			wantImage: "entry/wrap_multi_line_wrap_break.png",
		},
		"multi line WrapWord": {
			multiLine: true,
			wrap:      fyne.TextWrapWord,
			wantImage: "entry/wrap_multi_line_wrap_word.png",
		},
	} {
		t.Run(name, func(t *testing.T) {
			e, window := setupImageTest(tt.multiLine)
			defer teardownImageTest(window)
			c := window.Canvas()

			c.Focus(e)
			e.Wrapping = tt.wrap
			if tt.multiLine {
				e.SetText("A long text on short words w/o NLs or LFs.")
			} else {
				e.SetText("Testing Wrapping")
			}
			test.AssertImageMatches(t, tt.wantImage, c.Capture())
		})
	}
}

func TestMultiLineEntry_MinSize(t *testing.T) {
	entry := widget.NewEntry()
	singleMin := entry.MinSize()

	multi := widget.NewMultiLineEntry()
	multiMin := multi.MinSize()

	assert.Equal(t, singleMin.Width, multiMin.Width)
	assert.True(t, multiMin.Height > singleMin.Height)

	multi.MultiLine = false
	multiMin = multi.MinSize()
	assert.Equal(t, singleMin.Height, multiMin.Height)
}

func TestPasswordEntry_ActionItemSizeAndPlacement(t *testing.T) {
	e := widget.NewEntry()
	b := widget.NewButton("", func() {})
	b.Icon = theme.CancelIcon()
	e.ActionItem = b
	test.WidgetRenderer(e).Layout(e.MinSize())
	assert.Equal(t, fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()), b.Size())
	assert.Equal(t, fyne.NewPos(e.MinSize().Width-2*theme.Padding()-b.Size().Width, 2*theme.Padding()), b.Position())
}

func TestPasswordEntry_NewlineIgnored(t *testing.T) {
	entry := widget.NewPasswordEntry()
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

func TestPasswordEntry_Obfuscation(t *testing.T) {
	entry, window := setupPasswordImageTest()
	defer teardownImageTest(window)
	c := window.Canvas()

	test.AssertImageMatches(t, "password_entry/obfuscation_initial.png", c.Capture())

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertImageMatches(t, "password_entry/obfuscation_typed.png", c.Capture())
}

func TestPasswordEntry_Placeholder(t *testing.T) {
	entry, window := setupPasswordImageTest()
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetPlaceHolder("Password")
	test.AssertImageMatches(t, "password_entry/placeholder_initial.png", c.Capture())

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertImageMatches(t, "password_entry/placeholder_typed.png", c.Capture())
}

func TestPasswordEntry_Reveal(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	t.Run("NewPasswordEntry constructor", func(t *testing.T) {
		entry := widget.NewPasswordEntry()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(fyne.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
		entry.Move(fyne.NewPos(10, 10))
		c := window.Canvas()

		test.AssertImageMatches(t, "password_entry/initial.png", c.Capture())
		c.Focus(entry)

		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/concealed.png", c.Capture())

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/revealed.png", c.Capture())
		assert.Equal(t, entry, c.Focused())

		// update the Password field
		entry.Password = true
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/concealed.png", c.Capture())
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		tapPos := fyne.NewPos(140-theme.Padding()*2-theme.IconInlineSize()/2, 10+entry.Size().Height/2)
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/revealed.png", c.Capture())
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/concealed.png", c.Capture())
		assert.Equal(t, entry, c.Focused())
	})

	// This test cover backward compatibility use case when on an Entry widget
	// the Password field is set to true.
	// In this case the action item will be set when the renderer is created.
	t.Run("Entry with Password field", func(t *testing.T) {
		entry := &widget.Entry{}
		entry.Password = true
		entry.Refresh()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(fyne.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
		entry.Move(fyne.NewPos(10, 10))
		c := window.Canvas()

		test.AssertImageMatches(t, "password_entry/initial.png", c.Capture())
		c.Focus(entry)

		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/concealed.png", c.Capture())

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertImageMatches(t, "password_entry/revealed.png", c.Capture())
		assert.Equal(t, entry, c.Focused())
	})
}

func TestSingleLineEntry_NewlineIgnored(t *testing.T) {
	entry := &widget.Entry{MultiLine: false}
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

const (
	keyShiftLeftDown  fyne.KeyName = "LeftShiftDown"
	keyShiftLeftUp    fyne.KeyName = "LeftShiftUp"
	keyShiftRightDown fyne.KeyName = "RightShiftDown"
	keyShiftRightUp   fyne.KeyName = "RightShiftUp"
)

var typeKeys = func(e *widget.Entry, keys ...fyne.KeyName) {
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

func checkNewlineIgnored(t *testing.T, entry *widget.Entry) {
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

func getClickPosition(str string, row int) *fyne.PointEvent {
	x := fyne.MeasureText(str, theme.TextSize(), fyne.TextStyle{}).Width + theme.Padding()

	rowHeight := fyne.MeasureText("M", theme.TextSize(), fyne.TextStyle{}).Height
	y := theme.Padding() + row*rowHeight + rowHeight/2

	pos := fyne.NewPos(x, y)
	return &fyne.PointEvent{Position: pos}
}

func setupImageTest(multiLine bool) (*widget.Entry, fyne.Window) {
	app := test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	entry := &widget.Entry{MultiLine: multiLine}
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(150, 200))

	if multiLine {
		entry.Resize(fyne.NewSize(120, 100))
	} else {
		entry.Resize(entry.MinSize().Max(fyne.NewSize(120, 0)))
	}
	entry.Move(fyne.NewPos(10, 10))

	return entry, w
}

func setupPasswordImageTest() (*widget.Entry, fyne.Window) {
	app := test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	entry := widget.NewPasswordEntry()
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(150, 100))

	entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
	entry.Move(fyne.NewPos(10, 10))

	return entry, w
}

// Selects "sti" on line 2 of a new multiline
// T e s t i n g
// T e[s t i]n g
// T e s t i n g
func setupSelection(reverse bool) (*widget.Entry, fyne.Window) {
	e, window := setupImageTest(true)
	e.SetText("Testing\nTesting\nTesting")
	c := window.Canvas()
	c.Focus(e)
	if reverse {
		e.CursorRow = 1
		e.CursorColumn = 5
		typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft, fyne.KeyLeft)
	} else {
		e.CursorRow = 1
		e.CursorColumn = 2
		typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)
	}
	return e, window
}

func teardownImageTest(w fyne.Window) {
	w.Close()
	test.NewApp()
}
