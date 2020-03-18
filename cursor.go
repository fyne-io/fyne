package fyne

// Cursor represents a standard fyne cursor
type Cursor string

const (
	// DefaultCursor is the default cursor typically an arrow
	DefaultCursor Cursor = "Arrow"
	// TextCursor is the cursor often used to indicate text selection
	TextCursor Cursor = "Text"
	// CrosshairCursor is the cursor often used to indicate bitmaps
	CrosshairCursor Cursor = "Crosshair"
	// PointerCursor is the cursor often used to indicate a link
	PointerCursor Cursor = "Pointer"
	// HResizeCursor is the cursor often used for horizontal resize
	HResizeCursor Cursor = "HResize"
	// VResizeCursor is the cursor often used for vertical resize
	VResizeCursor Cursor = "VResize"
)
