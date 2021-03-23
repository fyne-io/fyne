package fyne

import "image"

// Canvas defines a graphical canvas to which a CanvasObject or Container can be added.
// Each canvas has a scale which is automatically applied during the render process.
type Canvas interface {
	Content() CanvasObject
	SetContent(CanvasObject)

	Refresh(CanvasObject)

	// Focus makes the provided item focused.
	// The item has to be added to the contents of the canvas before calling this.
	Focus(Focusable)
	// FocusNext focuses the next focusable item.
	// If no item is currently focused, the first focusable item is focused.
	// If the last focusable item is currently focused, the first focusable item is focused.
	//
	// Since: 2.0
	FocusNext()
	// FocusPrevious focuses the previous focusable item.
	// If no item is currently focused, the last focusable item is focused.
	// If the first focusable item is currently focused, the last focusable item is focused.
	//
	// Since: 2.0
	FocusPrevious()
	Unfocus()
	Focused() Focusable

	// Size returns the current size of this canvas
	Size() Size
	// Scale returns the current scale (multiplication factor) this canvas uses to render
	// The pixel size of a CanvasObject can be found by multiplying by this value.
	Scale() float32

	// Overlays returns the overlay stack.
	Overlays() OverlayStack

	OnTypedRune() func(rune)
	SetOnTypedRune(func(rune))
	OnTypedKey() func(*KeyEvent)
	SetOnTypedKey(func(*KeyEvent))
	AddShortcut(shortcut Shortcut, handler func(shortcut Shortcut))
	RemoveShortcut(shortcut Shortcut)

	Capture() image.Image

	// PixelCoordinateForPosition returns the x and y pixel coordinate for a given position on this canvas.
	// This can be used to find absolute pixel positions or pixel offsets relative to an object top left.
	PixelCoordinateForPosition(Position) (int, int)

	// InteractiveArea returns the position and size of the central interactive area.
	// Operating system elements may overlap the portions outside this area and widgets should avoid being outside.
	//
	// Since: 1.4
	InteractiveArea() (Position, Size)
}
