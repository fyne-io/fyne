package fyne

import (
	"sync"
)

// ShortcutHandler is a default implementation of the shortcut handler
// for the canvasObject
type ShortcutHandler struct {
	mu    sync.RWMutex
	entry map[string]func(Shortcut)
}

// HandleShortcut handle the registered shortcut
func (sh *ShortcutHandler) HandleShortcut(shortcut Shortcut) bool {
	if shortcut == nil {
		return false
	}
	if sc, ok := sh.entry[shortcut.ShortcutName()]; ok {
		sc(shortcut)
		return true
	}
	return false
}

// AddShortcut register an handler to be executed when ShortcutName command is triggered
func (sh *ShortcutHandler) AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut)) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.entry == nil {
		sh.entry = make(map[string]func(Shortcut))
	}
	sh.entry[shortcut.ShortcutName()] = handler
}

// Shortcut is the interface implemented by values with a custom formatter.
// The implementation of Format may call Sprint(f) or Fprint(f) etc.
// to generate its output.
type Shortcut interface {
	ShortcutName() string
}

// ShortcutPaste describes a shortcut paste action.
type ShortcutPaste struct {
	Clipboard Clipboard
}

// ShortcutName returns the shortcut name
func (se *ShortcutPaste) ShortcutName() string {
	return "Paste"
}

// ShortcutCopy describes a shortcut copy action.
type ShortcutCopy struct {
	Clipboard Clipboard
}

// ShortcutName returns the shortcut name
func (se *ShortcutCopy) ShortcutName() string {
	return "Copy"
}

// ShortcutCut describes a shortcut cut action.
type ShortcutCut struct {
	Clipboard Clipboard
}

// ShortcutName returns the shortcut name
func (se *ShortcutCut) ShortcutName() string {
	return "Cut"
}
