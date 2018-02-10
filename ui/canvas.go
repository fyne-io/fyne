package ui

import "image"

type Canvas interface {
	NewRectangle(image.Rectangle) (CanvasObject)
}
