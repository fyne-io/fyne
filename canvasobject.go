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

// ThemedObject indicates that the associated CanvasObject responds to theme
// changes. When the settings detect a theme change the object will be informed
// through the invocation of ApplyTheme().
type ThemedObject interface {
	ApplyTheme()
}

// TappableObject describes any CanvasObject that can also be tapped.
// This should be implemented by buttons etc that wish to handle pointer interactions.
type TappableObject interface {
	Tapped(*PointEvent)
	TappedSecondary(*PointEvent)
}

// ScrollableObject describes any CanvasObject that can also be scrolled.
// This is mostly used to implement the widget.ScrollContainer.
type ScrollableObject interface {
	Scrolled(*ScrollEvent)
}

// FocusableObject describes any CanvasObject that can respond to being focused.
// It will receive the OnFocusGained and OnFocusLost events appropriately and,
// when focused, it will also have OnKeyDown called as keys are pressed.
type FocusableObject interface {
	OnFocusGained()
	OnFocusLost()
	Focused() bool

	OnKeyDown(*KeyEvent)
}
