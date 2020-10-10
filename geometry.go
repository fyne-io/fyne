package fyne

// Position describes a generic X, Y coordinate relative to a parent Canvas
// or CanvasObject.
type Position struct {
	X int // The position from the parent's left edge
	Y int // The position from the parent's top edge
}

// NewPos returns a newly allocated Position representing the specified coordinates.
func NewPos(x int, y int) Position {
	return Position{x, y}
}

// Add returns a new Position that is the result of offsetting the current
// position by p2 X and Y.
func (p Position) Add(p2 Position) Position {
	return Position{p.X + p2.X, p.Y + p2.Y}
}

// IsZero returns whether the Position is at the zero-point.
func (p Position) IsZero() bool {
	return p.X == 0 && p.Y == 0
}

// Subtract returns a new Position that is the result of offsetting the current
// position by p2 -X and -Y.
func (p Position) Subtract(p2 Position) Position {
	return Position{p.X - p2.X, p.Y - p2.Y}
}

// Size describes something with width and height.
type Size struct {
	Width  int // The number of units along the X axis.
	Height int // The number of units along the Y axis.
}

// NewSize returns a newly allocated Size of the specified dimensions.
func NewSize(w int, h int) Size {
	return Size{w, h}
}

// Add returns a new Size that is the result of increasing the current size by
// s2 Width and Height.
func (s Size) Add(s2 Size) Size {
	return Size{s.Width + s2.Width, s.Height + s2.Height}
}

// IsZero returns whether the Size has zero width and zero height.
func (s Size) IsZero() bool {
	return s.Width == 0 && s.Height == 0
}

// Max returns a new Size that is the maximum of the current Size and s2.
func (s Size) Max(s2 Size) Size {
	maxW := Max(s.Width, s2.Width)
	maxH := Max(s.Height, s2.Height)

	return NewSize(maxW, maxH)
}

// Min returns a new Size that is the minimum of the current Size and s2.
func (s Size) Min(s2 Size) Size {
	minW := Min(s.Width, s2.Width)
	minH := Min(s.Height, s2.Height)

	return NewSize(minW, minH)
}

// Subtract returns a new Size that is the result of decreasing the current size
// by s2 Width and Height.
func (s Size) Subtract(s2 Size) Size {
	return Size{s.Width - s2.Width, s.Height - s2.Height}
}

// Union returns a new Size that is the maximum of the current Size and s2.
//
// Deprecated: use Max() instead
func (s Size) Union(s2 Size) Size {
	return s.Max(s2)
}
