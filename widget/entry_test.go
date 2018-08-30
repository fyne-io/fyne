package widget

import "testing"

import "github.com/stretchr/testify/assert"

import "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/theme"

func TestEntrySize(t *testing.T) {
	entry := NewEntry()
	min := entry.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestEntryAppend(t *testing.T) {
	entry := NewEntry()

	key := new(fyne.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)

	assert.Equal(t, entry.Text, "Hi")
}

func TestEntryBackspace(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")

	key := new(fyne.KeyEvent)
	key.Name = "BackSpace"
	entry.OnKeyDown(key)

	assert.Equal(t, entry.Text, "H")
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
	canvas := test.GetTestCanvas()

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
