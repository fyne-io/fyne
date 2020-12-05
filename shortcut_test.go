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

func TestShortcutHandler_RemoveShortcut(t *testing.T) {
	handler := &ShortcutHandler{}
	handler.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handler.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, len(handler.entry))

	handler.RemoveShortcut(&ShortcutCopy{})

	assert.Equal(t, 1, len(handler.entry))

	handler.RemoveShortcut(&ShortcutPaste{})

	assert.Equal(t, 0, len(handler.entry))
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
