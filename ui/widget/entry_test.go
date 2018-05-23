package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne/ui"

func TestAppend(t *testing.T) {
	entry := NewEntry()

	key := new(ui.KeyEvent)
	key.String = "H"
	entry.OnKeyDown(key)
	key.String = "i"
	entry.OnKeyDown(key)

	assert.Equal(t, entry.Text(), "Hi")
}

func TestBackspace(t *testing.T) {
	entry := NewEntry()
	entry.SetText("Hi")

	key := new(ui.KeyEvent)
	key.Name = "BackSpace"
	entry.OnKeyDown(key)

	assert.Equal(t, entry.Text(), "H")
}
