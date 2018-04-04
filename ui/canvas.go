package ui

// Canvas defines a graphical canvas to which a CanvasObject or Container can be added.
// Each canvas has a scale which is automatically applied during the render process.
type Canvas interface {
	SetContent(CanvasObject)
	Refresh(CanvasObject)

	Size() Size
	Scale() float32
	SetScale(float32)
}
