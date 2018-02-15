package ui

type Canvas interface {
	SetContent(CanvasObject)

	Scale() float32
	SetScale(float32)
}
