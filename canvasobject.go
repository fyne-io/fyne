package fyne

// CanvasObject describes any graphical object that can be added to a canvas.
// Objects have a size and position that can be controlled through this API.
// MinSize is used to determine the minimum size which this object shuold be displayed.
// Note: If this object is controlled as part of a Layout you should not call
// Resize(Size) or Move(Position).
type CanvasObject interface {
	CurrentSize() Size
	Resize(Size)
	CurrentPosition() Position
	Move(Position)

	MinSize() Size
}

// ThemedObject indicates that the associated CanvasObject reponse to theme
// changes. When the settings detect a theme change the object will be informed
// through the invocation of ApplyTheme().
type ThemedObject interface {
	ApplyTheme()
}

// ClickableObject describes any CanvasObject that can also be clicked
// (i.e. has mouse handlers). This should be implemented by buttons etc that
// wish to handle pointer interactions.
type ClickableObject interface {
	OnMouseDown(*MouseEvent)
}

// FocusableObject describes any CanvasObject that can respond to being focused.
// It will receive the OnFocusGained and OnFocusLost events appropriately and,
// when focussed, it will also have OnKeyDown called as keys are pressed.
type FocusableObject interface {
	OnFocusGained()
	OnFocusLost()
	Focused() bool

	OnKeyDown(*KeyEvent)
}
