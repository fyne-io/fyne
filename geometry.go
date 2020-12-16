package fyne

type Vector struct {
	X, Y float32
}

// NewVector returns a newly allocated Vector representing a generic X, Y coordinate.
func NewVector(x float32, y float32) Vector {
	return Vector{x, y}
}

// Position describes a generic X, Y coordinate relative to a parent Canvas
// or CanvasObject.
type Position struct {
	X float32 // The position from the parent's left edge
	Y float32 // The position from the parent's top edge
}

// NewPos returns a newly allocated Position representing the specified coordinates.
func NewPos(x float32, y float32) Position {
	return Position{x, y}
}

// Add returns a new Position that is the result of offsetting the current
// position by p2 X and Y.
func (p Position) Add(p2 Position) Position {
	return Position{p.X + p2.X, p.Y + p2.Y}
}

// IsZero returns whether the Position is at the zero-point.
func (p Position) IsZero() bool {
	return p.X == 0.0 && p.Y == 0.0
}

// Subtract returns a new Position that is the result of offsetting the current
// position by p2 -X and -Y.
func (p Position) Subtract(p2 Position) Position {
	return Position{p.X - p2.X, p.Y - p2.Y}
}

// Size describes something with width and height.
type Size struct {
	Width  float32 // The number of units along the X axis.
	Height float32 // The number of units along the Y axis.
}

// NewSize returns a newly allocated Size of the specified dimensions.
func NewSize(w float32, h float32) Size {
	return Size{w, h}
}

// Add returns a new Size that is the result of increasing the current size by
// s2 Width and Height.
func (s Size) Add(s2 Size) Size {
	return Size{s.Width + s2.Width, s.Height + s2.Height}
}

// IsZero returns whether the Size has zero width and zero height.
func (s Size) IsZero() bool {
	return s.Width == 0.0 && s.Height == 0.0
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
