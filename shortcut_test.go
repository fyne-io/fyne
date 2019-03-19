package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortcutHandler_AddShortcut(t *testing.T) {
	handle := &ShortcutHandler{}

	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, len(handle.entry))
}

func TestShortcutHandler_HandleShortcut(t *testing.T) {
	handle := &ShortcutHandler{}
	cut, copy, paste := false, false, false

	handle.AddShortcut(&ShortcutCut{}, func(shortcut Shortcut) {
		cut = true
	})
	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {
		copy = true
	})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {
		paste = true
	})

	handle.TypedShortcut(&ShortcutCut{})
	assert.True(t, cut)
	handle.TypedShortcut(&ShortcutCopy{})
	assert.True(t, copy)
	handle.TypedShortcut(&ShortcutPaste{})
	assert.True(t, paste)
}
