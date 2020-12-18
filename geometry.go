package fyne

// Vec2 marks geometry types that can operate as a coordinate vector.
type Vec2 interface {
	Components() (float32, float32)
	IsZero() bool
}

// Vector is a generic X, Y coordinate or size representation.
type Vector struct {
	X, Y float32
}

// NewVector returns a newly allocated Vector representing a generic X, Y coordinate.
func NewVector(x float32, y float32) Vector {
	return Vector{x, y}
}

// Components returns the X and Y elements of this Vector.
func (v Vector) Components() (float32, float32) {
	return v.X, v.Y
}

// IsZero returns whether the Position is at the zero-point.
func (v Vector) IsZero() bool {
	return v.X == 0.0 && v.Y == 0.0
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
func (p Position) Add(v Vec2) Position {
	x, y := v.Components()
	return Position{p.X + x, p.Y + y}
}

// Components returns the X and Y elements of this Position
func (p Position) Components() (float32, float32) {
	return p.X, p.Y
}

// IsZero returns whether the Position is at the zero-point.
func (p Position) IsZero() bool {
	return p.X == 0.0 && p.Y == 0.0
}

// Subtract returns a new Position that is the result of offsetting the current
// position by p2 -X and -Y.
func (p Position) Subtract(v Vec2) Position {
	x, y := v.Components()
	return Position{p.X - x, p.Y - y}
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
func (s Size) Add(v Vec2) Size {
	w, h := v.Components()
	return Size{s.Width + w, s.Height + h}
}

// IsZero returns whether the Size has zero width and zero height.
func (s Size) IsZero() bool {
	return s.Width == 0.0 && s.Height == 0.0
}

// Max returns a new Size that is the maximum of the current Size and s2.
func (s Size) Max(v Vec2) Size {
	x, y := v.Components()

	maxW := Max(s.Width, x)
	maxH := Max(s.Height, y)

	return NewSize(maxW, maxH)
}

// Min returns a new Size that is the minimum of the current Size and s2.
func (s Size) Min(v Vec2) Size {
	x, y := v.Components()

	minW := Min(s.Width, x)
	minH := Min(s.Height, y)

	return NewSize(minW, minH)
}

// Components returns the Width and Height elements of this Size
func (s Size) Components() (float32, float32) {
	return s.Width, s.Height
}

// Subtract returns a new Size that is the result of decreasing the current size
// by s2 Width and Height.
func (s Size) Subtract(v Vec2) Size {
	w, h := v.Components()
	return Size{s.Width - w, s.Height - h}
}
