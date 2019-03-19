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

	assert.True(t, handle.TypedShortcut(&ShortcutCut{}))
	assert.True(t, cutCalled)
	assert.True(t, handle.TypedShortcut(&ShortcutCopy{}))
	assert.True(t, copyCalled)
	assert.True(t, handle.TypedShortcut(&ShortcutPaste{}))
	assert.True(t, pasteCalled)
}

func TestShortcutHandler_HandleShortcut_Failures(t *testing.T) {
	handle := &ShortcutHandler{}
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.False(t, handle.TypedShortcut(nil))
	assert.False(t, handle.TypedShortcut(&ShortcutCut{}))

}
