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
	Size() Size
	Resize(Size)
	Position() Position
	Move(Position)
	MinSize() Size

	// visibility
	Visible() bool
	Show()
	Hide()
}

// Themeable indicates that the associated CanvasObject responds to theme
// changes. When the settings detect a theme change the object will be informed
// through the invocation of ApplyTheme().
type Themeable interface {
	ApplyTheme()
}

// Tappable describes any CanvasObject that can also be tapped.
// This should be implemented by buttons etc that wish to handle pointer interactions.
type Tappable interface {
	Tapped(*PointEvent)
	TappedSecondary(*PointEvent)
}

// Scrollable describes any CanvasObject that can also be scrolled.
// This is mostly used to implement the widget.ScrollContainer.
type Scrollable interface {
	Scrolled(*ScrollEvent)
}

// Focusable describes any CanvasObject that can respond to being focused.
// It will receive the FocusGained and FocusLost events appropriately.
// When focused it will also have TypedRune called as text is input and
// TypedKey called when other keys are pressed.
type Focusable interface {
	FocusGained()
	FocusLost()
	Focused() bool

	TypedRune(rune)
	TypedKey(*KeyEvent)
}

// Shortcutable describes any CanvasObject that can respond to shortcut commands (quit, cut, copy, and paste).
type Shortcutable interface {
	TypedShortcut(shortcut Shortcut) bool
}
