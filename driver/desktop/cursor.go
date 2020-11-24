package desktop

import "image"

// Cursor interface is used for objects usable as cursor
type Cursor interface {
	Image() (image.Image, int, int) // Image returns the image for the given cursor, or nil to hide the cursor
}

// StandardCursor represents a standard fyne cursor
type StandardCursor int

// Image will always return nil for StandardCursor
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
