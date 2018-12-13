package widget

import "testing"

import "github.com/stretchr/testify/assert"

import "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/theme"

func TestEntry_MinSize(t *testing.T) {
	entry := NewEntry()
	min := entry.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestEntry_OnKeyDown(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Insert(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hi", entry.Text)

	left := &fyne.KeyEvent{Name: "Left"}
	entry.OnKeyDown(left)

	key.String = "o"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hoi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	key := new(fyne.KeyEvent)
	key.Name = "BackSpace"
	entry.OnKeyDown(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_BackspaceBeyondText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)

	key := new(fyne.KeyEvent)
	key.Name = "BackSpace"
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_BackspaceNewline(t *testing.T) {
	entry := NewEntry()
	entry.SetText("H\ni")

	down := &fyne.KeyEvent{Name: "Down"}
	entry.OnKeyDown(down)

	key := new(fyne.KeyEvent)
	key.Name = "BackSpace"
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntry_OnKeyDown_Backspace_Unicode(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "è"
	entry.OnKeyDown(key)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	bs := new(fyne.KeyEvent)
	bs.Name = "BackSpace"
	entry.OnKeyDown(bs)
	assert.Equal(t, "", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_OnKeyDown_Delete(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")
	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	key := new(fyne.KeyEvent)
	key.Name = "Delete"
	entry.OnKeyDown(key)

	assert.Equal(t, "H", entry.Text)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_OnKeyDown_DeleteBeyondText(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")

	key := new(fyne.KeyEvent)
	key.Name = "Delete"
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntry_OnKeyDown_DeleteNewline(t *testing.T) {
	entry := NewEntry()
	entry.SetText("H\ni")

	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)

	key := new(fyne.KeyEvent)
	key.Name = "Delete"
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
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

	entry.OnFocusGained()
	assert.True(t, entry.Focused())

	entry.OnFocusLost()
	assert.False(t, entry.Focused())
}

func TestEntryWindowFocus(t *testing.T) {
	entry := NewEntry()
	canvas := test.Canvas()

	canvas.Focus(entry)
	assert.True(t, entry.Focused())
}

func TestEntryFocusHighlight(t *testing.T) {
	entry := NewEntry()

	entry.OnFocusGained()
	assert.True(t, entry.focused)

	entry.OnFocusLost()
	assert.False(t, entry.focused)
}

func TestEntry_CursorRow(t *testing.T) {
	entry := NewEntry()
	entry.SetText("test")
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: "Down"}
	entry.OnKeyDown(down)
	assert.Equal(t, 0, entry.CursorRow)

	// 2 lines, this should increment
	entry.SetText("test\nrows")
	entry.OnKeyDown(down)
	assert.Equal(t, 1, entry.CursorRow)

	up := &fyne.KeyEvent{Name: "Up"}
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestEntry_CursorColumn(t *testing.T) {
	entry := NewEntry()
	entry.SetText("")
	assert.Equal(t, 0, entry.CursorColumn)

	// only 0 columns, do nothing
	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorColumn)

	// 1, this should increment
	entry.SetText("a")
	entry.OnKeyDown(right)
	assert.Equal(t, 1, entry.CursorColumn)

	left := &fyne.KeyEvent{Name: "Left"}
	entry.OnKeyDown(left)
	assert.Equal(t, 0, entry.CursorColumn)

	// don't go beyond left
	entry.OnKeyDown(left)
	assert.Equal(t, 0, entry.CursorColumn)
}

func TestEntry_CursorColumn_Wrap(t *testing.T) {
	entry := NewEntry()
	entry.SetText("a\nb")
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// go to end of line
	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)

	// wrap to new line
	entry.OnKeyDown(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 0, entry.CursorColumn)

	// and back
	left := &fyne.KeyEvent{Name: "Left"}
	entry.OnKeyDown(left)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestEntry_CursorColumn_Jump(t *testing.T) {
	entry := NewEntry()
	entry.SetText("a\nbc")

	// go to end of text
	right := &fyne.KeyEvent{Name: "Right"}
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	entry.OnKeyDown(right)
	assert.Equal(t, 1, entry.CursorRow)
	assert.Equal(t, 2, entry.CursorColumn)

	// go up, to a shorter line
	up := &fyne.KeyEvent{Name: "Up"}
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)
	assert.Equal(t, 1, entry.CursorColumn)
}

func TestPasswordEntry_NewlineIgnored(t *testing.T) {
	entry := NewPasswordEntry()
	entry.SetText("test")
	assert.Equal(t, 0, entry.CursorRow)

	// only 1 line, do nothing
	down := &fyne.KeyEvent{Name: "Down"}
	entry.OnKeyDown(down)
	assert.Equal(t, 0, entry.CursorRow)

	// return is ignored, do nothing
	ret := &fyne.KeyEvent{Name: "Return"}
	entry.OnKeyDown(ret)
	assert.Equal(t, 0, entry.CursorRow)

	up := &fyne.KeyEvent{Name: "Up"}
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)

	// don't go beyond top
	entry.OnKeyDown(up)
	assert.Equal(t, 0, entry.CursorRow)
}

func TestPasswordEntry_Obfuscation(t *testing.T) {
	entry := NewPasswordEntry()

	key := new(fyne.KeyEvent)
	key.String = "Hié™שרה"
	entry.OnKeyDown(key)
	assert.Equal(t, "Hié™שרה", entry.Text)
	assert.Equal(t, "*******", entry.label().Text)
}
