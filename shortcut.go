package fyne

// Shortcut reprents a particular combination of actions that result in a shortcut.(i.e. copy, paste...)
type Shortcut int

const (
	// ShortcutNone represents a combination of actions that do not have shortcut associated to
	ShortcutNone Shortcut = iota
	// ShortcutCopy represents a combination of actions used to copy text from the clipboard
	ShortcutCopy
	// ShortcutPaste represents a combination of actions used to paste text to the clipboard
	ShortcutPaste
	// ShortcutCut represents a combination of actions used to cut text to the clipboard
	ShortcutCut
)
