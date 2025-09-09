package fyne

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func syncMapLen(m *sync.Map) (n int) {
	m.Range(func(_, _ any) bool {
		n++
		return true
	})
	return n
}

func TestShortcutHandler_AddShortcut(t *testing.T) {
	handle := &ShortcutHandler{}

	handle.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handle.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, syncMapLen(&handle.entry))
}

func TestShortcutHandler_RemoveShortcut(t *testing.T) {
	handler := &ShortcutHandler{}
	handler.AddShortcut(&ShortcutCopy{}, func(shortcut Shortcut) {})
	handler.AddShortcut(&ShortcutPaste{}, func(shortcut Shortcut) {})

	assert.Equal(t, 2, syncMapLen(&handler.entry))

	handler.RemoveShortcut(&ShortcutCopy{})

	assert.Equal(t, 1, syncMapLen(&handler.entry))

	handler.RemoveShortcut(&ShortcutPaste{})

	assert.Equal(t, 0, syncMapLen(&handler.entry))
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
