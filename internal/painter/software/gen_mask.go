//go:build !software

package software

import (
	"image"

	"fyne.io/fyne/v2"
)

func genMask(obj fyne.CanvasObject, bounds, clip image.Rectangle) *image.RGBA {
	mask := image.NewAlpha(bounds)

	return mask
}
