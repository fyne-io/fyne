package driver

import (
	"image"

	"fyne.io/fyne/v2"
)

// Painter describes a simple type that can render canvases
//
// Since: 2.9
type Painter interface {
	Paint(fyne.Canvas) image.Image
}
