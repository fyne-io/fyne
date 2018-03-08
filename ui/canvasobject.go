package ui

// CanvasObject describes any graphical object that can be added to a canvas.
// Objects have a size and position that can be controlled through this API.
// MinSize is used to determine the minimum size which this object shuold be displayed.
// Note: If this object is controlled as part of a Layout you should not call
// Resize(Size) or Move(Position).
type CanvasObject interface {
	CurrentSize() Size
	Resize(Size)
	CurrentPosition() Position
	Move(Position)

	MinSize() Size
}
