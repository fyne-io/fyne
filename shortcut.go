package fyne

import (
	"sync"
)

// ShortcutName reprents a particular combination of actions that result in a shortcut.(i.e. copy, paste...)
type ShortcutName string

const (
	// ShortcutNone represents a combination of actions that do not have shortcut associated to
	ShortcutNone ShortcutName = "None"
	// ShortcutUnknown represents a combination of actions that do not have fyne shortcut associated to
	ShortcutUnknown ShortcutName = "Unknown"
	// ShortcutCopy represents a combination of actions used to copy text from the clipboard
	ShortcutCopy ShortcutName = "Copy"
	// ShortcutPaste represents a combination of actions used to paste text to the clipboard
	ShortcutPaste ShortcutName = "Paste"
	// ShortcutCut represents a combination of actions used to cut text to the clipboard
	ShortcutCut ShortcutName = "Cut"
)

// ShortcutDefaultHandler is a default implementation of the shortcut handler
// for the canvasObject
type ShortcutDefaultHandler struct {
	mu    sync.RWMutex
	entry map[ShortcutName]func(ShortcutEventer)
}

// TriggerShortcutHandler handle the shortcut registered as ShortcutName passing the ShortcutEvent
func (sh *ShortcutDefaultHandler) TriggerShortcutHandler(event ShortcutEventer) {
	if event == nil {
		return
	}
	if sc, ok := sh.entry[event.Shortcut()]; ok {
		sc(event)
	}
}

// RegisterShortcutHandler register an handler to be executed when ShortcutName command is triggered
func (sh *ShortcutDefaultHandler) RegisterShortcutHandler(shortcutName ShortcutName, handler func(ShortcutEventer)) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.entry == nil {
		sh.entry = make(map[ShortcutName]func(ShortcutEventer))
	}
	sh.entry[shortcutName] = handler
}
