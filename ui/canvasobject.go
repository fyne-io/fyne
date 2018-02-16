package ui

type CanvasObject interface {
	CurrentSize() Size
	Resize(Size)
	MinSize() Size
}
