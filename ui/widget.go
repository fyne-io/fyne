package ui

type Widget interface {
	CurrentSize() Size
	Resize(Size)
	CurrentPosition() Position
	Move(Position)

	MinSize() Size

	// TODO should this move to a widget impl?... (private)
	Layout() []CanvasObject
}
