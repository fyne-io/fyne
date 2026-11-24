package mobile

import (
	"io"
	"log"
	"os"
	"testing"

	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func Test_device(t *testing.T) {
	dev := &device{}

	assert.Equal(t, true, dev.IsMobile())
	assert.Equal(t, false, dev.IsBrowser())
	assert.Equal(t, false, dev.HasKeyboard())
}

func Test_keyboardText(t *testing.T) {
	entry := widget.NewEntry()
	entry.SetText("a😀b")
	entry.CursorColumn = 2

	text, cursor, ok := keyboardText(entry)
	assert.True(t, ok)
	assert.Equal(t, "a😀b", text)
	assert.Equal(t, 2, cursor) // runes, not bytes or UTF-16 units

	_, _, ok = keyboardText(widget.NewButton("Tap", nil))
	assert.False(t, ok)
	_, _, ok = keyboardText(nil)
	assert.False(t, ok)
}

func Test_device_HideVirtualKeyboard_BeforeRun(t *testing.T) {
	// Discarding log output for tests
	// The following method logs an error:
	// (&device{}).HideVirtualKeyboard() // should not crash!
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	(&device{}).HideVirtualKeyboard() // should not crash!
}
