package fyne

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	String string

	Name      KeyName
	Modifiers Modifier
}

// Shortcut returns the shortcut associated to the KeyEvent, if any
func (ev *KeyEvent) Shortcut() KeyShortcut {
	if ev.Modifiers == 0 {
		return ShortcutNone
	}

	switch ev.Name {
	case KeyV:
		if ev.Modifiers&ControlModifier != 0 {
			return ShortcutPaste
		}
	case KeyC:
		if ev.Modifiers&ControlModifier != 0 {
			return ShortcutCopy
		}
	case KeyX:
		if ev.Modifiers&ControlModifier != 0 {
			return ShortcutCut
		}
	}

	return ShortcutNone
}

// MouseEvent describes a pointer input event. The position is relative to the top-left
// of the CanvasObject this is triggered on.
type MouseEvent struct {
	Position Position    // The position of the event
	Button   MouseButton // The mouse button which caused the event
}
