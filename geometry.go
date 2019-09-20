package fyne

// Size describes something with width and height.
type Size struct {
	Width  int // The number of units along the X axis.
	Height int // The number of units along the Y axis.
}

// Add returns a new Size that is the result of increasing the current size by
// s2 Width and Height.
func (s1 Size) Add(s2 Size) Size {
	return Size{s1.Width + s2.Width, s1.Height + s2.Height}
}

// Subtract returns a new Size that is the result of decreasing the current size
// by s2 Width and Height.
func (s1 Size) Subtract(s2 Size) Size {
	return Size{s1.Width - s2.Width, s1.Height - s2.Height}
}

// Union returns a new Size that is the maximum of the current Size and s2.
func (s1 Size) Union(s2 Size) Size {
	maxW := Max(s1.Width, s2.Width)
	maxH := Max(s1.Height, s2.Height)

	return NewSize(maxW, maxH)
}

// NewSize returns a newly allocated Size of the specified dimensions.
func NewSize(w int, h int) Size {
	return Size{w, h}
}

// Position describes a generic X, Y coordinate relative to a parent Canvas
// or CanvasObject.
type Position struct {
	X int // The position from the parent's left edge
	Y int // The position from the parent's top edge
}

// Add returns a new Position that is the result of offsetting the current
// position by p2 X and Y.
func (p1 Position) Add(p2 Position) Position {
	return Position{p1.X + p2.X, p1.Y + p2.Y}
}

// Subtract returns a new Position that is the result of offsetting the current
// position by p2 -X and -Y.
func (p1 Position) Subtract(p2 Position) Position {
	return Position{p1.X - p2.X, p1.Y - p2.Y}
}

// NewPos returns a newly allocated Position representing the specified coordinates.
func NewPos(x int, y int) Position {
	return Position{x, y}
}
