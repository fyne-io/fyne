package desktop

// Cursor represents a standard fyne cursor
type Cursor int

const (
	// DefaultCursor is the default cursor typically an arrow
	DefaultCursor Cursor = iota
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
)

// Cursorable describes any CanvasObject that needs a cursor change
type Cursorable interface {
	Cursor() Cursor
}
