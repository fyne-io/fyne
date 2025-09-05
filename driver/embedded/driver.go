package embedded

import (
	"image"

	"fyne.io/fyne/v2"
)

type Driver interface {
	Render(image.Image)
	Run(func())

	ScreenSize() fyne.Size
	Queue() chan Event
}
