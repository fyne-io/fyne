package widget

import "testing"

import "github.com/stretchr/testify/assert"

import "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/theme"

func TestEntryPassword_MinSize(t *testing.T) {
	entry := NewEntryPassword()
	min := entry.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestEntryPassword_OnKeyDown(t *testing.T) {
	entry := NewEntryPassword()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)

	assert.Equal(t, "Hi", entry.Text)
}

func TestEntryPassword_OnKeyDown_Insert(t *testing.T) {
	entry := NewEntryPassword()

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

func TestEntryPassword_OnKeyDown_Backspace(t *testing.T) {
	entry := NewEntryPassword()
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

func TestEntryPassword_OnKeyDown_BackspaceBeyondContent(t *testing.T) {
	entry := NewEntryPassword()
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

func TestEntryPassword_OnKeyDown_Delete(t *testing.T) {
	entry := NewEntryPassword()
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

func TestEntryPassword_OnKeyDown_DeleteBeyondContent(t *testing.T) {
	entry := NewEntryPassword()
	entry.SetText("Hi")

	key := new(fyne.KeyEvent)
	key.Name = "Delete"
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)
	entry.OnKeyDown(key)

	assert.Equal(t, "", entry.Text)
}

func TestEntryPassword_Notify(t *testing.T) {
	entry := NewEntryPassword()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntryPassword_Focus(t *testing.T) {
	entry := NewEntryPassword()

	entry.OnFocusGained()
	assert.True(t, entry.Focused())

	entry.OnFocusLost()
	assert.False(t, entry.Focused())
}

func TestEntryPassword_WindowFocus(t *testing.T) {
	entry := NewEntryPassword()
	canvas := test.Canvas()

	canvas.Focus(entry)
	assert.True(t, entry.Focused())
}

func TestEntryPassword_FocusHighlight(t *testing.T) {
	entry := NewEntryPassword()

	entry.OnFocusGained()
	assert.True(t, entry.focused)

	entry.OnFocusLost()
	assert.False(t, entry.focused)
}

func TestEntryPassword_CursorRow(t *testing.T) {
	entry := NewEntryPassword()
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

func TestEntryPassword_CursorColumn(t *testing.T) {
	entry := NewEntryPassword()
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
