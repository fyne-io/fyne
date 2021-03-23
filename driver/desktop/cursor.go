package desktop

import "image"

// Cursor interface is used for objects that desire a specific cursor.
//
// Since: 2.0
type Cursor interface {
	// Image returns the image for the given cursor, or nil if none should be shown.
	// It also returns the x and y pixels that should act as the hot-spot (measured from top left corner).
	Image() (image.Image, int, int)
}

// StandardCursor represents a standard Fyne cursor.
// These values were previously of type `fyne.Cursor`.
//
// Since: 2.0
type StandardCursor int

// Image is not used for any of the StandardCursor types.
//
// Since: 2.0
func (d StandardCursor) Image() (image.Image, int, int) {
	return nil, 0, 0
}

const (
	// DefaultCursor is the default cursor typically an arrow
	DefaultCursor StandardCursor = iota
	// TextCursor is the cursor often used to indicate text selection
	TextCursor
	// CrosshairCursor is the cursor often used to indicate bitmaps
	CrosshairCursor
	// PointerCursor is the cursor often used to indicate a link
	PointerCursor
	// HResizeCursor is the cursor often used to indicate horizontal resize
	HResizeCursor
	// VResizeCursor is the cursor often used to indicate vertical resize
	VResizeCursor
	// HiddenCursor will cause the cursor to not be shown
	HiddenCursor
)

// Cursorable describes any CanvasObject that needs a cursor change
type Cursorable interface {
	Cursor() Cursor
}
