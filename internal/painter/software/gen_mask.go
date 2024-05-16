//go:build !software

package software

import (
	"image"

	"fyne.io/fyne/v2"
)

func genMask(c fyne.Canvas, obj fyne.CanvasObject, objBounds, bounds image.Rectangle) *image.Alpha {
	mask := image.NewAlpha(bounds)

	for i := 0; i < len(mask.Pix); i++ {
		mask.Pix[i] = 0xff
	}

	return mask
}
