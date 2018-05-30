package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"

func TestEntryTestSize(t *testing.T) {
	entry := NewEntry()
	min := entry.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestEntryTestAppend(t *testing.T) {
	entry := NewEntry()

	key := new(ui.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)

	assert.Equal(t, entry.Text(), "Hi")
}

func TestEntryTestBackspace(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")

	key := new(ui.KeyEvent)
	key.Name = "BackSpace"
	entry.OnKeyDown(key)

	assert.Equal(t, entry.Text(), "H")
}

func TestEntryTestNotify(t *testing.T) {
	entry := NewEntry()
	changed := false

	entry.OnChanged = func(string) {
		changed = true
	}
	entry.SetText("Test")

	assert.True(t, changed)
}

func TestEntryTestFocusHighlight(t *testing.T) {
	entry := NewEntry()
	bg := entry.Layout(entry.MinSize())[0].(*canvas.Rectangle)
	color := bg.FillColor

	entry.OnFocusGained()
	assert.NotEqual(t, bg.FillColor, color)

	entry.OnFocusLost()
	assert.Equal(t, bg.FillColor, color)
}
