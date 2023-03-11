package fyne

// CanvasObject describes any graphical object that can be added to a canvas.
// Objects have a size and position that can be controlled through this API.
// MinSize is used to determine the minimum size which this object should be displayed.
// An object will be visible by default but can be hidden with Hide() and re-shown with Show().
//
// Note: If this object is controlled as part of a Layout you should not call
// Resize(Size) or Move(Position).
type CanvasObject interface {
	// geometry

	// MinSize returns the minimum size this object needs to be drawn.
	MinSize() Size
	// Move moves this object to the given position relative to its parent.
	// This should only be called if your object is not in a container with a layout manager.
	Move(Position)
	// Position returns the current position of the object relative to its parent.
	Position() Position
	// Resize resizes this object to the given size.
	// This should only be called if your object is not in a container with a layout manager.
	Resize(Size)
	// Size returns the current size of this object.
	Size() Size

	// visibility

	// Hide hides this object.
	Hide()
	// Visible returns whether this object is visible or not.
	Visible() bool
	// Show shows this object.
	Show()

	// Refresh must be called if this object should be redrawn because its inner state changed.
	Refresh()
}

// Disableable describes any CanvasObject that can be disabled.
// This is primarily used with objects that also implement the Tappable interface.
type Disableable interface {
	Enable()
	Disable()
	Disabled() bool
}

// DoubleTappable describes any CanvasObject that can also be double tapped.
type DoubleTappable interface {
	DoubleTapped(*PointEvent)
}

// Draggable indicates that a CanvasObject can be dragged.
// This is used for any item that the user has indicated should be moved across the screen.
type Draggable interface {
	Dragged(*DragEvent)
	DragEnd()
}

// Focusable describes any CanvasObject that can respond to being focused.
// It will receive the FocusGained and FocusLost events appropriately.
// When focused it will also have TypedRune called as text is input and
// TypedKey called when other keys are pressed.
//
// Note: You must not change canvas state (including overlays or focus) in FocusGained or FocusLost
// or you would end up with a dead-lock.
type Focusable interface {
	// FocusGained is a hook called by the focus handling logic after this object gained the focus.
	FocusGained()
	// FocusLost is a hook called by the focus handling logic after this object lost the focus.
	FocusLost()

	// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
	TypedRune(rune)
	// TypedKey is a hook called by the input handling logic on key events if this object is focused.
	TypedKey(*KeyEvent)
}

// Scrollable describes any CanvasObject that can also be scrolled.
// This is mostly used to implement the widget.ScrollContainer.
type Scrollable interface {
	Scrolled(*ScrollEvent)
}

// SecondaryTappable describes a CanvasObject that can be right-clicked or long-tapped.
type SecondaryTappable interface {
	TappedSecondary(*PointEvent)
}

// Shortcutable describes any CanvasObject that can respond to shortcut commands (quit, cut, copy, and paste).
type Shortcutable interface {
	TypedShortcut(Shortcut)
}

// Tabbable describes any object that needs to accept the Tab key presses.
//
// Since: 2.1
type Tabbable interface {
	// AcceptsTab() is a hook called by the key press handling logic.
	// If it returns true then the Tab key events will be sent using TypedKey.
	AcceptsTab() bool
}

// Tappable describes any CanvasObject that can also be tapped.
// This should be implemented by buttons etc that wish to handle pointer interactions.
type Tappable interface {
	Tapped(*PointEvent)
}
