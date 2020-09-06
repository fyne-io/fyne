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

func TestShortcutHandler_RemoveShortut(t *testing.T) {
	handle := &ShortcutHandler{}
	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, len(handle.entry))

	handle.RemoveShortcut(&ShortcutCopy{})

	assert.Equal(t, 1, len(handle.entry))

	handle.RemoveShortcut(&ShortcutPaste{})

	assert.Equal(t, 0, len(handle.entry))
}

func TestShortcutHandler_HandleShortcut(t *testing.T) {
	handle := &ShortcutHandler{}
	cutCalled, copyCalled, pasteCalled := false, false, false

	handle.AddShortcut(&ShortcutCut{}, func(shortcut Shortcut) {
		cutCalled = true
	})
	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {
		copyCalled = true
	})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {
		pasteCalled = true
	})

	handle.TypedShortcut(&ShortcutCut{})
	assert.True(t, cutCalled)
	handle.TypedShortcut(&ShortcutCopy{})
	assert.True(t, copyCalled)
	handle.TypedShortcut(&ShortcutPaste{})
	assert.True(t, pasteCalled)
}
