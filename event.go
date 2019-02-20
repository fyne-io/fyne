package fyne

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	Name KeyName
}

// PointEvent describes a pointer input event. The position is relative to the
// top-left of the CanvasObject this is triggered on.
type PointEvent struct {
	Position Position // The position of the event
}

// ScrollEvent defines the parameters of a pointer or other scroll event.
// The DeltaX and DeltaY represent how large the scroll was in two dimensions.
type ScrollEvent struct {
	PointEvent
	DeltaX, DeltaY int
}

// ShortcutEvent describes a shortcut input event.
type ShortcutEvent struct {
	ShortcutName
}

// Shortcut returns the shortcut name associated to the event
func (se *ShortcutEvent) Shortcut() ShortcutName {
	return se.ShortcutName
}

// ShortcutClipboardEvent describes a shortcut clipoard event.
type ShortcutClipboardEvent struct {
	ShortcutName
	Clipboard
}

// Shortcut returns the shortcut name associated to the event
func (se *ShortcutClipboardEvent) Shortcut() ShortcutName {
	return se.ShortcutName
}

// ShortcutDesktopEvent describes a shortcut clipoard event.
type ShortcutDesktopEvent struct {
	ShortcutName
	KeyName
	Modifier
}

// Shortcut returns the shortcut name associated to the event
func (se *ShortcutDesktopEvent) Shortcut() ShortcutName {
	return se.ShortcutName
}

// ShortcutEventer describes a shortcut event.
type ShortcutEventer interface {
	Shortcut() ShortcutName
}
