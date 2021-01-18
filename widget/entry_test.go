package widget_test

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Binding(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("Init")
	assert.Equal(t, "Init", entry.Text)

	str := binding.NewString()
	entry.Bind(str)
	waitForBinding()
	assert.Equal(t, "", entry.Text)

	err := str.Set("Updated")
	assert.Nil(t, err)
	waitForBinding()
	assert.Equal(t, "Updated", entry.Text)

	entry.SetText("Typed")
	v, err := str.Get()
	assert.Nil(t, err)
	assert.Equal(t, "Typed", v)

	entry.Unbind()
	waitForBinding()
	assert.Equal(t, "Typed", entry.Text)
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

func TestEntry_CursorColumn_Wrap2(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.Resize(fyne.NewSize(64, 64))
	entry.SetText("1234")
	entry.CursorColumn = 3
	test.Type(entry, "a")
	test.Type(entry, "b")
	test.Type(entry, "c")
	assert.Equal(t, 0, entry.CursorColumn)
	assert.Equal(t, 1, entry.CursorRow)
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(70, 70))
	test.AssertImageMatches(t, "entry/wrap_multi_line_cursor.png", w.Canvas().Capture())
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
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_empty.xml", c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_disabled_empty.xml", c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_empty.xml", c)

	entry.SetPlaceHolder("Type!")
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_placeholder.xml", c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_disabled_placeholder.xml", c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_placeholder.xml", c)

	entry.SetText("Hello")
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_custom_value.xml", c)

	entry.Disable()
	assert.True(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_disabled_custom_value.xml", c)

	entry.Enable()
	assert.False(t, entry.Disabled())
	test.AssertRendersToMarkup(t, "entry/disableable_enabled_custom_value.xml", c)
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
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/focus_gained.xml", c)

	entry.FocusLost()
	test.AssertRendersToMarkup(t, "entry/focus_lost.xml", c)

	window.Canvas().Focus(entry)
	test.AssertRendersToMarkup(t, "entry/focus_gained.xml", c)
}

func TestEntry_FocusWithPopUp(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.TapSecondaryAt(entry, fyne.NewPos(1, 1))

	test.AssertRendersToMarkup(t, "entry/focus_with_popup_initial.xml", c)

	test.TapCanvas(c, fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, "entry/focus_with_popup_entry_selected.xml", c)

	test.TapSecondaryAt(entry, fyne.NewPos(1, 1))
	test.AssertRendersToMarkup(t, "entry/focus_with_popup_initial.xml", c)

	test.TapCanvas(c, fyne.NewPos(5, 5))
	test.AssertRendersToMarkup(t, "entry/focus_with_popup_dismissed.xml", c)
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
	entry.SetPlaceHolder("Hello")
	assert.Equal(t, entry.MinSize().Width, min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	entry.Wrapping = fyne.TextWrapOff
	entry.Refresh()
	assert.Greater(t, entry.MinSize().Width, min.Width)

	min = entry.MinSize()
	entry.ActionItem = canvas.NewCircle(color.Black)
	assert.Equal(t, min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0)), entry.MinSize())
}

func TestEntryMultiline_MinSize(t *testing.T) {
	entry := widget.NewMultiLineEntry()
	min := entry.MinSize()
	entry.SetText("Hello")
	assert.Equal(t, entry.MinSize().Width, min.Width)
	assert.Equal(t, entry.MinSize().Height, min.Height)

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	entry.Wrapping = fyne.TextWrapOff
	entry.Refresh()
	assert.Greater(t, entry.MinSize().Width, min.Width)

	entry.Wrapping = fyne.TextWrapBreak
	entry.Refresh()
	assert.Equal(t, entry.MinSize().Width, min.Width)

	min = entry.MinSize()
	entry.ActionItem = canvas.NewCircle(color.Black)
	assert.Equal(t, min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0)), entry.MinSize())
}

func TestEntry_MultilineSelect(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// Extend the selection down one row
	typeKeys(e, fyne.KeyDown)
	test.AssertRendersToMarkup(t, "entry/selection_add_one_row_down.xml", c)
	assert.Equal(t, "sting\nTesti", e.SelectedText())

	typeKeys(e, fyne.KeyUp)
	test.AssertRendersToMarkup(t, "entry/selection_remove_one_row_up.xml", c)
	assert.Equal(t, "sti", e.SelectedText())

	typeKeys(e, fyne.KeyUp)
	test.AssertRendersToMarkup(t, "entry/selection_remove_add_one_row_up.xml", c)
	assert.Equal(t, "ng\nTe", e.SelectedText())
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
	entry, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("Hi")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/on_key_down_newline_initial.xml", c)

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
	test.AssertRendersToMarkup(t, "entry/on_key_down_newline_typed.xml", c)
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
		e, window := setupImageTest(t, false)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing")
		test.AssertRendersToMarkup(t, "entry/select_initial.xml", c)

		// move right, press & hold shift and pagedown
		typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyPageDown)
		assert.Equal(t, "esting", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_single_line_shift_pagedown.xml", c)

		// while shift is held press pageup
		typeKeys(e, fyne.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_single_line_shift_pageup.xml", c)

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, fyne.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_single_line_pagedown.xml", c)
	})

	t.Run("page down single line", func(*testing.T) {
		e, window := setupImageTest(t, true)
		defer teardownImageTest(window)
		c := window.Canvas()

		c.Focus(e)
		e.SetText("Testing\nTesting\nTesting")
		test.AssertRendersToMarkup(t, "entry/select_multi_line_initial.xml", c)

		// move right, press & hold shift and pagedown
		typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyPageDown)
		assert.Equal(t, "esting\nTesting\nTesting", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_multi_line_shift_pagedown.xml", c)

		// while shift is held press pageup
		typeKeys(e, fyne.KeyPageUp)
		assert.Equal(t, "T", e.SelectedText())
		assert.Equal(t, 0, e.CursorRow)
		assert.Equal(t, 0, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_multi_line_shift_pageup.xml", c)

		// release shift and press pagedown
		typeKeys(e, keyShiftLeftUp, fyne.KeyPageDown)
		assert.Equal(t, "", e.SelectedText())
		assert.Equal(t, 2, e.CursorRow)
		assert.Equal(t, 7, e.CursorColumn)
		test.AssertRendersToMarkup(t, "entry/select_multi_line_pagedown.xml", c)
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
	test.AssertRendersToMarkup(t, "entry/placeholder_with_text.xml", c)

	entry.SetText("")
	assert.Equal(t, "", entry.Text)
	test.AssertRendersToMarkup(t, "entry/placeholder_without_text.xml", c)
}

func TestEntry_Select(t *testing.T) {
	for name, tt := range map[string]struct {
		keys          []fyne.KeyName
		text          string
		setupReverse  bool
		wantMarkup    string
		wantSelection string
		wantText      string
	}{
		"delete single-line": {
			keys:       []fyne.KeyName{fyne.KeyDelete},
			wantText:   "Testing\nTeng\nTesting",
			wantMarkup: "entry/selection_delete_single_line.xml",
		},
		"delete multi-line": {
			keys:       []fyne.KeyName{fyne.KeyDown, fyne.KeyDelete},
			wantText:   "Testing\nTeng",
			wantMarkup: "entry/selection_delete_multi_line.xml",
		},
		"delete reverse multi-line": {
			keys:         []fyne.KeyName{keyShiftLeftDown, fyne.KeyDown, fyne.KeyDelete},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantMarkup:   "entry/selection_delete_reverse_multi_line.xml",
		},
		"delete select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"delete reverse select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"delete select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		"delete reverse select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyDelete, fyne.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		// The backspace delete behaviour is the same as via delete.
		"backspace single-line": {
			keys:       []fyne.KeyName{fyne.KeyBackspace},
			wantText:   "Testing\nTeng\nTesting",
			wantMarkup: "entry/selection_delete_single_line.xml",
		},
		"backspace multi-line": {
			keys:       []fyne.KeyName{fyne.KeyDown, fyne.KeyBackspace},
			wantText:   "Testing\nTeng",
			wantMarkup: "entry/selection_delete_multi_line.xml",
		},
		"backspace reverse multi-line": {
			keys:         []fyne.KeyName{keyShiftLeftDown, fyne.KeyDown, fyne.KeyBackspace},
			setupReverse: true,
			wantText:     "Testing\nTestisting",
			wantMarkup:   "entry/selection_delete_reverse_multi_line.xml",
		},
		"backspace select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyDown},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"backspace reverse select down with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyDown},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "ng\nTe",
			wantMarkup:    "entry/selection_delete_and_add_down.xml",
		},
		"backspace select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyUp},
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		"backspace reverse select up with Shift held": {
			keys:          []fyne.KeyName{keyShiftLeftDown, fyne.KeyBackspace, fyne.KeyUp},
			setupReverse:  true,
			wantText:      "Testing\nTeng\nTesting",
			wantSelection: "sting\nTe",
			wantMarkup:    "entry/selection_delete_and_add_up.xml",
		},
		// Erase the selection and add a newline at selection start
		"enter": {
			keys:       []fyne.KeyName{fyne.KeyEnter},
			wantText:   "Testing\nTe\nng\nTesting",
			wantMarkup: "entry/selection_enter.xml",
		},
		"enter reverse": {
			keys:         []fyne.KeyName{fyne.KeyEnter},
			setupReverse: true,
			wantText:     "Testing\nTe\nng\nTesting",
			wantMarkup:   "entry/selection_enter.xml",
		},
		"replace": {
			text:       "hello",
			wantText:   "Testing\nTehellong\nTesting",
			wantMarkup: "entry/selection_replace.xml",
		},
		"replace reverse": {
			text:         "hello",
			setupReverse: true,
			wantText:     "Testing\nTehellong\nTesting",
			wantMarkup:   "entry/selection_replace.xml",
		},
		"deselect and delete": {
			keys:       []fyne.KeyName{keyShiftLeftUp, fyne.KeyLeft, fyne.KeyDelete},
			wantText:   "Testing\nTeting\nTesting",
			wantMarkup: "entry/selection_deselect_delete.xml",
		},
		"deselect and delete holding shift": {
			keys:       []fyne.KeyName{keyShiftLeftUp, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyDelete},
			wantText:   "Testing\nTeting\nTesting",
			wantMarkup: "entry/selection_deselect_delete.xml",
		},
		// ensure that backspace doesn't leave a selection start at the old cursor position
		"deselect and backspace holding shift": {
			keys:       []fyne.KeyName{keyShiftLeftUp, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyBackspace},
			wantText:   "Testing\nTsting\nTesting",
			wantMarkup: "entry/selection_deselect_backspace.xml",
		},
		// clear selection, select a character and while holding shift issue two backspaces
		"deselect, select and double backspace": {
			keys:       []fyne.KeyName{keyShiftLeftUp, fyne.KeyRight, fyne.KeyLeft, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyBackspace, fyne.KeyBackspace},
			wantText:   "Testing\nTeing\nTesting",
			wantMarkup: "entry/selection_deselect_select_backspace.xml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			entry, window := setupSelection(t, tt.setupReverse)
			defer teardownImageTest(window)
			c := window.Canvas()

			if tt.text != "" {
				test.Type(entry, tt.text)
			} else {
				typeKeys(entry, tt.keys...)
			}
			assert.Equal(t, tt.wantText, entry.Text)
			assert.Equal(t, tt.wantSelection, entry.SelectedText())
			test.AssertRendersToMarkup(t, tt.wantMarkup, c)
		})
	}
}

func TestEntry_SelectAll(t *testing.T) {
	e, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("First Row\nSecond Row\nThird Row")
	test.AssertRendersToMarkup(t, "entry/select_all_initial.xml", c)

	shortcut := &fyne.ShortcutSelectAll{}
	e.TypedShortcut(shortcut)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 9, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/select_all_selected.xml", c)
}

func TestEntry_SelectAll_EmptyEntry(t *testing.T) {
	entry := widget.NewEntry()
	entry.TypedShortcut(&fyne.ShortcutSelectAll{})

	assert.Equal(t, "", entry.SelectedText())
}

func TestEntry_SelectEndWithoutShift(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// end after releasing shift
	typeKeys(e, keyShiftLeftUp, fyne.KeyEnd)
	test.AssertRendersToMarkup(t, "entry/selection_end.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectHomeEnd(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// Hold shift to continue selection
	typeKeys(e, keyShiftLeftDown)

	// T e[s t i]n g -> end -> // T e[s t i n g]
	typeKeys(e, fyne.KeyEnd)
	test.AssertRendersToMarkup(t, "entry/selection_add_to_end.xml", c)
	assert.Equal(t, "sting", e.SelectedText())

	// T e[s t i n g] -> home -> [T e]s t i n g
	typeKeys(e, fyne.KeyHome)
	test.AssertRendersToMarkup(t, "entry/selection_add_to_home.xml", c)
	assert.Equal(t, "Te", e.SelectedText())
}

func TestEntry_SelectHomeWithoutShift(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	// home after releasing shift
	typeKeys(e, keyShiftLeftUp, fyne.KeyHome)
	test.AssertRendersToMarkup(t, "entry/selection_home.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapDown(t *testing.T) {
	// down snaps to end, but it also moves
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyDown)
	assert.Equal(t, 2, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_down.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapLeft(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyLeft)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 2, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_left.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapRight(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_right.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectSnapUp(t *testing.T) {
	// up snaps to start, but it also moves
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 1, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)

	typeKeys(e, keyShiftLeftUp, fyne.KeyUp)
	assert.Equal(t, 0, e.CursorRow)
	assert.Equal(t, 5, e.CursorColumn)
	test.AssertRendersToMarkup(t, "entry/selection_snap_up.xml", c)
	assert.Equal(t, "", e.SelectedText())
}

func TestEntry_SelectedText(t *testing.T) {
	e, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Focus(e)
	e.SetText("Testing")
	test.AssertRendersToMarkup(t, "entry/select_initial.xml", c)
	assert.Equal(t, "", e.SelectedText())

	// move right, press & hold shift and move right
	typeKeys(e, fyne.KeyRight, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight)
	assert.Equal(t, "es", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_selected.xml", c)

	// release shift
	typeKeys(e, keyShiftLeftUp)
	// press shift and move
	typeKeys(e, keyShiftLeftDown, fyne.KeyRight)
	assert.Equal(t, "est", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_add_selection.xml", c)

	// release shift and move right
	typeKeys(e, keyShiftLeftUp, fyne.KeyRight)
	assert.Equal(t, "", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_move_wo_shift.xml", c)

	// press shift and move left
	typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft)
	assert.Equal(t, "st", e.SelectedText())
	test.AssertRendersToMarkup(t, "entry/select_select_left.xml", c)
}

func TestEntry_SelectionHides(t *testing.T) {
	e, window := setupSelection(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	c.Unfocus()
	test.AssertRendersToMarkup(t, "entry/selection_focus_lost.xml", c)
	assert.Equal(t, "sti", e.SelectedText())

	c.Focus(e)
	test.AssertRendersToMarkup(t, "entry/selection_focus_gained.xml", c)
	assert.Equal(t, "sti", e.SelectedText())
}

func TestEntry_SetPlaceHolder(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	assert.Equal(t, 0, len(entry.Text))

	entry.SetPlaceHolder("Test")
	assert.Equal(t, 0, len(entry.Text))
	test.AssertRendersToMarkup(t, "entry/set_placeholder_set.xml", c)

	entry.SetText("Hi")
	assert.Equal(t, 2, len(entry.Text))
	test.AssertRendersToMarkup(t, "entry/set_placeholder_replaced.xml", c)
}

func TestEntry_Disable_KeyDown(t *testing.T) {
	entry := widget.NewEntry()

	test.Type(entry, "H")
	entry.Disable()
	test.Type(entry, "i")
	assert.Equal(t, "H", entry.Text)

	entry.Enable()
	test.Type(entry, "i")
	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_Disable_OnFocus(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Disable()
	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/focused_disabled.xml", c)

	entry.Enable()
	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/focused_enabled.xml", c)
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
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Text = "Test"
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_changed.xml", c)
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

func TestEntry_SetTextStyle(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Text = "Styled Text"
	entry.TextStyle = fyne.TextStyle{Bold: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_style_bold.xml", c)

	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_style_monospace.xml", c)

	entry.TextStyle = fyne.TextStyle{Italic: true}
	entry.Refresh()
	test.AssertRendersToMarkup(t, "entry/set_text_style_italic.xml", c)
}

func TestEntry_Submit(t *testing.T) {
	t.Run("Callback", func(t *testing.T) {
		var submission string
		entry := &widget.Entry{
			OnSubmitted: func(s string) {
				submission = s
			},
		}
		t.Run("SingleLine_Enter", func(t *testing.T) {
			entry.MultiLine = false
			entry.SetText("a")
			entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
			assert.Equal(t, "a", entry.Text)
			assert.Equal(t, "a", submission)
		})
		t.Run("SingleLine_Return", func(t *testing.T) {
			entry.MultiLine = false
			entry.SetText("b")
			entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
			assert.Equal(t, "b", entry.Text)
			assert.Equal(t, "b", submission)
		})
		t.Run("MultiLine_ShiftEnter", func(t *testing.T) {
			entry.MultiLine = true
			entry.SetText("c")
			typeKeys(entry, keyShiftLeftDown, fyne.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "c", entry.Text)
			assert.Equal(t, "c", submission)
			entry.SetText("d")
			typeKeys(entry, keyShiftRightDown, fyne.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "d", entry.Text)
			assert.Equal(t, "d", submission)
		})
		t.Run("MultiLine_ShiftReturn", func(t *testing.T) {
			entry.MultiLine = true
			entry.SetText("e")
			typeKeys(entry, keyShiftLeftDown, fyne.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "e", entry.Text)
			assert.Equal(t, "e", submission)
			entry.SetText("f")
			typeKeys(entry, keyShiftRightDown, fyne.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "f", entry.Text)
			assert.Equal(t, "f", submission)
		})
	})
	t.Run("NoCallback", func(t *testing.T) {
		entry := &widget.Entry{}
		t.Run("SingleLine_Enter", func(t *testing.T) {
			entry.MultiLine = false
			entry.SetText("a")
			entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
			assert.Equal(t, "a", entry.Text)
		})
		t.Run("SingleLine_Return", func(t *testing.T) {
			entry.MultiLine = false
			entry.SetText("b")
			entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
			assert.Equal(t, "b", entry.Text)
		})
		t.Run("MultiLine_ShiftEnter", func(t *testing.T) {
			entry.MultiLine = true
			entry.SetText("c")
			typeKeys(entry, keyShiftLeftDown, fyne.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "\nc", entry.Text)
			entry.SetText("d")
			typeKeys(entry, keyShiftRightDown, fyne.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "\nd", entry.Text)
		})
		t.Run("MultiLine_ShiftReturn", func(t *testing.T) {
			entry.MultiLine = true
			entry.SetText("e")
			typeKeys(entry, keyShiftLeftDown, fyne.KeyReturn, keyShiftLeftUp)
			assert.Equal(t, "\ne", entry.Text)
			entry.SetText("f")
			typeKeys(entry, keyShiftRightDown, fyne.KeyReturn, keyShiftRightUp)
			assert.Equal(t, "\nf", entry.Text)
		})
	})
}

func TestEntry_Tapped(t *testing.T) {
	entry, window := setupImageTest(t, true)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetText("MMM\nWWW\n")
	test.AssertRendersToMarkup(t, "entry/tapped_initial.xml", c)

	entry.FocusGained()
	test.AssertRendersToMarkup(t, "entry/tapped_focused.xml", c)

	testCharSize := theme.TextSize()
	pos := fyne.NewPos(entryOffset+theme.Padding()+testCharSize*1.5, entryOffset+theme.Padding()+testCharSize/2) // tap in the middle of the 2nd "M"
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_2nd_m.xml", c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	pos = fyne.NewPos(entryOffset+theme.Padding()+testCharSize*2.5, entryOffset+theme.Padding()+testCharSize/2) // tap in the middle of the 3rd "M"
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_3rd_m.xml", c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	pos = fyne.NewPos(entryOffset+theme.Padding()+testCharSize*4, entryOffset+theme.Padding()+testCharSize/2) // tap after text
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_after_last_col.xml", c)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 3, entry.CursorColumn)

	pos = fyne.NewPos(entryOffset+testCharSize, entryOffset+testCharSize*4) // tap below rows
	test.TapCanvas(window.Canvas(), pos)
	test.AssertRendersToMarkup(t, "entry/tapped_tapped_after_last_row.xml", c)
	assert.Equal(t, 2, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_TappedSecondary(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	tapPos := fyne.NewPos(20, 10)
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_full_menu.xml", c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Disable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_read_menu.xml", c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	c.Overlays().Remove(c.Overlays().Top())

	entry.Password = true
	entry.Refresh()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_no_password_menu.xml", c)
	assert.Nil(t, c.Overlays().Top(), "No popup for disabled password")

	entry.Enable()
	test.TapSecondaryAt(entry, tapPos)
	test.AssertRendersToMarkup(t, "entry/tapped_secondary_password_menu.xml", c)
	assert.Equal(t, 1, len(c.Overlays().List()))
}

func TestEntry_TextWrap(t *testing.T) {
	for name, tt := range map[string]struct {
		multiLine bool
		want      string
		wrap      fyne.TextWrap
	}{
		"single line WrapOff": {
			want: "entry/wrap_single_line_off.xml",
		},
		// Disallowed - fallback to TextWrapOff
		"single line Truncate": {
			wrap: fyne.TextTruncate,
			want: "entry/wrap_single_line_off.xml",
		},
		// Disallowed - fallback to TextWrapOff
		"single line WrapBreak": {
			wrap: fyne.TextWrapBreak,
			want: "entry/wrap_single_line_off.xml",
		},
		// Disallowed - fallback to TextWrapOff
		"single line WrapWord": {
			wrap: fyne.TextWrapWord,
			want: "entry/wrap_single_line_off.xml",
		},
		"multi line WrapOff": {
			multiLine: true,
			want:      "entry/wrap_multi_line_off.xml",
		},
		// Disallowed - fallback to TextWrapOff
		"multi line Truncate": {
			multiLine: true,
			wrap:      fyne.TextTruncate,
			want:      "entry/wrap_multi_line_off.xml",
		},
		"multi line WrapBreak": {
			multiLine: true,
			wrap:      fyne.TextWrapBreak,
			want:      "entry/wrap_multi_line_wrap_break.xml",
		},
		"multi line WrapWord": {
			multiLine: true,
			wrap:      fyne.TextWrapWord,
			want:      "entry/wrap_multi_line_wrap_word.xml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			e, window := setupImageTest(t, tt.multiLine)
			defer teardownImageTest(window)
			c := window.Canvas()

			c.Focus(e)
			e.Wrapping = tt.wrap
			if tt.multiLine {
				e.SetText("A long text on short words w/o NLs or LFs.")
			} else {
				e.SetText("Testing Wrapping")
			}
			test.AssertRendersToMarkup(t, tt.want, c)
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

func TestNewEntryWithData(t *testing.T) {
	str := binding.NewString()
	err := str.Set("Init")
	assert.Nil(t, err)

	entry := widget.NewEntryWithData(str)
	waitForBinding()
	assert.Equal(t, "Init", entry.Text)

	entry.SetText("Typed")
	v, err := str.Get()
	assert.Nil(t, err)
	assert.Equal(t, "Typed", v)
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
	entry, window := setupPasswordTest(t)
	defer teardownImageTest(window)
	c := window.Canvas()

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertRendersToMarkup(t, "password_entry/obfuscation_typed.xml", c)
}

func TestPasswordEntry_Placeholder(t *testing.T) {
	entry, window := setupPasswordTest(t)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.SetPlaceHolder("Password")
	test.AssertRendersToMarkup(t, "password_entry/placeholder_initial.xml", c)

	test.Type(entry, "Hié™שרה")
	assert.Equal(t, "Hié™שרה", entry.Text)
	test.AssertRendersToMarkup(t, "password_entry/placeholder_typed.xml", c)
}

func TestPasswordEntry_Reveal(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	t.Run("NewPasswordEntry constructor", func(t *testing.T) {
		entry := widget.NewPasswordEntry()
		window := test.NewWindow(entry)
		defer window.Close()
		window.Resize(fyne.NewSize(150, 100))
		entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
		entry.Move(fyne.NewPos(10, 10))
		c := window.Canvas()

		test.AssertRendersToMarkup(t, "password_entry/initial.xml", c)

		c.Focus(entry)
		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/revealed.xml", c)
		assert.Equal(t, entry, c.Focused())

		// update the Password field
		entry.Password = true
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		tapPos := fyne.NewPos(140-theme.Padding()*2-theme.IconInlineSize()/2, 10+entry.Size().Height/2)
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/revealed.xml", c)
		assert.Equal(t, entry, c.Focused())

		// tap on action icon
		test.TapCanvas(c, tapPos)
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)
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

		test.AssertRendersToMarkup(t, "password_entry/initial.xml", c)

		c.Focus(entry)
		test.Type(entry, "Hié™שרה")
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/concealed.xml", c)

		// update the Password field
		entry.Password = false
		entry.Refresh()
		assert.Equal(t, "Hié™שרה", entry.Text)
		test.AssertRendersToMarkup(t, "password_entry/revealed.xml", c)
		assert.Equal(t, entry, c.Focused())
	})
}

func TestSingleLineEntry_NewlineIgnored(t *testing.T) {
	entry := &widget.Entry{MultiLine: false}
	entry.SetText("test")

	checkNewlineIgnored(t, entry)
}

const (
	entryOffset = 10

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

func setupImageTest(t *testing.T, multiLine bool) (*widget.Entry, fyne.Window) {
	test.NewApp()

	entry := &widget.Entry{MultiLine: multiLine}
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(150, 200))

	if multiLine {
		entry.Resize(fyne.NewSize(120, 100))
	} else {
		entry.Resize(entry.MinSize().Max(fyne.NewSize(120, 0)))
	}
	entry.Move(fyne.NewPos(10, 10))

	if multiLine {
		test.AssertRendersToMarkup(t, "entry/initial_multiline.xml", w.Canvas())
	} else {
		test.AssertRendersToMarkup(t, "entry/initial.xml", w.Canvas())
	}

	return entry, w
}

func setupPasswordTest(t *testing.T) (*widget.Entry, fyne.Window) {
	test.NewApp()

	entry := widget.NewPasswordEntry()
	w := test.NewWindow(entry)
	w.Resize(fyne.NewSize(150, 100))

	entry.Resize(entry.MinSize().Max(fyne.NewSize(130, 0)))
	entry.Move(fyne.NewPos(entryOffset, entryOffset))

	test.AssertRendersToMarkup(t, "password_entry/initial.xml", w.Canvas())

	return entry, w
}

// Selects "sti" on line 2 of a new multiline
// T e s t i n g
// T e[s t i]n g
// T e s t i n g
func setupSelection(t *testing.T, reverse bool) (*widget.Entry, fyne.Window) {
	e, window := setupImageTest(t, true)
	e.SetText("Testing\nTesting\nTesting")
	c := window.Canvas()
	c.Focus(e)
	if reverse {
		e.CursorRow = 1
		e.CursorColumn = 5
		typeKeys(e, keyShiftLeftDown, fyne.KeyLeft, fyne.KeyLeft, fyne.KeyLeft)
		test.AssertRendersToMarkup(t, "entry/selection_initial_reverse.xml", c)
		assert.Equal(t, "sti", e.SelectedText())
	} else {
		e.CursorRow = 1
		e.CursorColumn = 2
		typeKeys(e, keyShiftLeftDown, fyne.KeyRight, fyne.KeyRight, fyne.KeyRight)
		test.AssertRendersToMarkup(t, "entry/selection_initial.xml", c)
		assert.Equal(t, "sti", e.SelectedText())
	}

	return e, window
}

func teardownImageTest(w fyne.Window) {
	w.Close()
	test.NewApp()
}

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}
