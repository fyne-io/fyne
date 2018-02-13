package ui

type Canvas interface {
	AddObject(CanvasObject)

	Scale() float32
	SetScale(float32)
}
