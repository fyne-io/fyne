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

// TypedShortcut handle the registered shortcut
func (sh *ShortcutHandler) TypedShortcut(shortcut Shortcut) {
	if _, ok := sh.entry[shortcut.ShortcutName()]; !ok {
		return
	}

	sh.entry[shortcut.ShortcutName()](shortcut)
}

// AddShortcut register an handler to be executed when the shortcut action is triggered
func (sh *ShortcutHandler) AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut)) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.entry == nil {
		sh.entry = make(map[string]func(Shortcut))
	}
	sh.entry[shortcut.ShortcutName()] = handler
}

// RemoveShortcut removes a registered shortcut
func (sh *ShortcutHandler) RemoveShortcut(shortcut Shortcut) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.entry == nil {
		return
	}

	delete(sh.entry, shortcut.ShortcutName())
}

// Shortcut is the interface used to describe a shortcut action
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

// ShortcutSelectAll describes a shortcut selectAll action.
type ShortcutSelectAll struct{}

// ShortcutName returns the shortcut name
func (se *ShortcutSelectAll) ShortcutName() string {
	return "SelectAll"
}
