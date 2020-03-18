package fyne

// Cursor represents a standard fyne cursor
type Cursor string

const (
	// DefaultCursor is the default cursor typically an arrow
	DefaultCursor Cursor = "DefaultCursor"
	// TextCursor is the cursor often used to indicate text selection
	TextCursor Cursor = "TextCursor"
	// CrosshairCursor is the cursor often used to indicate bitmaps
	CrosshairCursor Cursor = "CrosshairCursor"
	// PointerCursor is the cursor often used to indicate a link
	PointerCursor Cursor = "PointerCursor"
	// HResizeCursor is the cursor often used for horizontal resize
	HResizeCursor Cursor = "HResizeCursor"
	// VResizeCursor is the cursor often used for vertical resize
	VResizeCursor Cursor = "VResizeCursor"
)
