package fyne

import (
	"sync"
)

// ShortcutHandler is a default implementation of the shortcut handler
// for the canvasObject
type ShortcutHandler struct {
	mu    sync.RWMutex
	entry map[string]func(Shortcuter)
}

// HandleShortcut handle the registered shortcut
func (sh *ShortcutHandler) HandleShortcut(shortcut Shortcuter) {
	if shortcut == nil {
		return
	}
	if sc, ok := sh.entry[shortcut.Shortcut()]; ok {
		sc(shortcut)
	}
}

// AddShortcut register an handler to be executed when ShortcutName command is triggered
func (sh *ShortcutHandler) AddShortcut(shortcut Shortcuter, handler func(shortcut Shortcuter)) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.entry == nil {
		sh.entry = make(map[string]func(Shortcuter))
	}
	sh.entry[shortcut.Shortcut()] = handler
}

// Shortcuter is the interface implemented by values with a custom formatter.
// The implementation of Format may call Sprint(f) or Fprint(f) etc.
// to generate its output.
type Shortcuter interface {
	Shortcut() string
}

// ShortcutPaste describes a shortcut paste action.
type ShortcutPaste struct {
	Clipboard
}

// Shortcut returns the shortcut type
func (se *ShortcutPaste) Shortcut() string {
	return "Paste"
}

// ShortcutCopy describes a shortcut copy action.
type ShortcutCopy struct {
	Clipboard
}

// Shortcut returns the shortcut type
func (se *ShortcutCopy) Shortcut() string {
	return "Copy"
}

// ShortcutCut describes a shortcut cut action.
type ShortcutCut struct {
	Clipboard
}

// Shortcut returns the shortcut type
func (se *ShortcutCut) Shortcut() string {
	return "Cut"
}
