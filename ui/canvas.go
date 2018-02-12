package ui

import "image"

type Canvas interface {
	NewRectangle(image.Rectangle) CanvasObject
	NewText(string) CanvasTextObject

	Scale() float32
	SetScale(float32)
}
