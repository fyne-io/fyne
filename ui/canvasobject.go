package ui

type CanvasObject interface {
	CurrentSize() Size
	Resize(Size)
	CurrentPosition() Position
	Move(Position)

	MinSize() Size
}
