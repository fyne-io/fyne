package fyne

import "sync"

// ShortcutHandler is a default implementation of the shortcut handler
// for [CanvasObject].
type ShortcutHandler struct {
	entry sync.Map // map[string]func(Shortcut)
}

// TypedShortcut handle the registered shortcut
func (sh *ShortcutHandler) TypedShortcut(shortcut Shortcut) {
	val, ok := sh.entry.Load(shortcut.ShortcutName())
	if !ok {
		return
	}

	f := val.(func(Shortcut))
	f(shortcut)
}

// AddShortcut register a handler to be executed when the shortcut action is triggered
func (sh *ShortcutHandler) AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut)) {
	sh.entry.Store(shortcut.ShortcutName(), handler)
}

// RemoveShortcut removes a registered shortcut
func (sh *ShortcutHandler) RemoveShortcut(shortcut Shortcut) {
	sh.entry.Delete(shortcut.ShortcutName())
}

// Shortcut is the interface used to describe a shortcut action
type Shortcut interface {
	ShortcutName() string
}

// KeyboardShortcut describes a shortcut meant to be triggered by a keyboard action.
type KeyboardShortcut interface {
	Shortcut
	Key() KeyName
	Mod() KeyModifier
}

var _ KeyboardShortcut = (*ShortcutPaste)(nil)

// ShortcutPaste describes a shortcut paste action.
type ShortcutPaste struct {
	Clipboard Clipboard
}

// Key returns the [KeyName] for this shortcut.
func (se *ShortcutPaste) Key() KeyName {
	return KeyV
}

// Mod returns the [KeyModifier] for this shortcut.
func (se *ShortcutPaste) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutPaste) ShortcutName() string {
	return "Paste"
}

var _ KeyboardShortcut = (*ShortcutCopy)(nil)

// ShortcutCopy describes a shortcut copy action.
type ShortcutCopy struct {
	Clipboard Clipboard
}

// Key returns the [KeyName] for this shortcut.
func (se *ShortcutCopy) Key() KeyName {
	return KeyC
}

// Mod returns the [KeyModifier] for this shortcut.
func (se *ShortcutCopy) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutCopy) ShortcutName() string {
	return "Copy"
}

var _ KeyboardShortcut = (*ShortcutCut)(nil)

// ShortcutCut describes a shortcut cut action.
type ShortcutCut struct {
	Clipboard Clipboard
}

// Key returns the [KeyName] for this shortcut.
func (se *ShortcutCut) Key() KeyName {
	return KeyX
}

// Mod returns the [KeyModifier] for this shortcut.
func (se *ShortcutCut) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutCut) ShortcutName() string {
	return "Cut"
}

var _ KeyboardShortcut = (*ShortcutSelectAll)(nil)

// ShortcutSelectAll describes a shortcut selectAll action.
type ShortcutSelectAll struct{}

// Key returns the [KeyName] for this shortcut.
func (se *ShortcutSelectAll) Key() KeyName {
	return KeyA
}

// Mod returns the [KeyModifier] for this shortcut.
func (se *ShortcutSelectAll) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutSelectAll) ShortcutName() string {
	return "SelectAll"
}

var _ KeyboardShortcut = (*ShortcutUndo)(nil)

// ShortcutUndo describes a shortcut undo action.
//
// Since: 2.5
type ShortcutUndo struct{}

// Key returns the [KeyName] for this shortcut.
func (se *ShortcutUndo) Key() KeyName {
	return KeyZ
}

// Mod returns the [KeyModifier] for this shortcut.
func (se *ShortcutUndo) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutUndo) ShortcutName() string {
	return "Undo"
}

var _ KeyboardShortcut = (*ShortcutRedo)(nil)

// ShortcutRedo describes a shortcut redo action.
//
// Since: 2.5
type ShortcutRedo struct{}

// Key returns the [KeyName] for this shortcut.
func (se *ShortcutRedo) Key() KeyName {
	return KeyY
}

// Mod returns the [KeyModifier] for this shortcut.
func (se *ShortcutRedo) Mod() KeyModifier {
	return KeyModifierShortcutDefault
}

// ShortcutName returns the shortcut name
func (se *ShortcutRedo) ShortcutName() string {
	return "Redo"
}
