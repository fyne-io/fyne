package painter

import (
	"image"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images

	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// PaintImage renders a given fyne Image to a Go standard image
// If a fyne.Canvas is given and the image’s fill mode is “fill original” the image’s min size has
// to fit its original size. If it doesn’t, PaintImage does not paint the image but adjusts its min size.
// The image will then be painted on the next frame because of the min size change.
func PaintImage(img *canvas.Image, c fyne.Canvas, width, height int) image.Image {
	if img.Size().IsZero() && c == nil { // an image without size or canvas won't get rendered unless we setup
		img.Resize(fyne.NewSize(float32(width), float32(height)))
	}
	dst, err := paintImage(img, width, height)
	if err != nil {
		fyne.LogError("failed to paint image", err)
	}

	return dst
}

func paintImage(img *canvas.Image, width, height int) (dst image.Image, err error) {
	if width <= 0 || height <= 0 {
		return
	}

	dst = img.Image
	if dst == nil {
		dst = image.NewNRGBA(image.Rect(0, 0, width, height))
	}

	size := dst.Bounds().Size()
	if width != size.X || height != size.Y {
		dst = scaleImage(dst, width, height, img.ScaleMode)
	}
	return
}

func scaleImage(pixels image.Image, scaledW, scaledH int, scale canvas.ImageScale) image.Image {
	if scale == canvas.ImageScaleFastest || scale == canvas.ImageScalePixels {
		// do not perform software scaling
		return pixels
	}

	bounds := pixels.Bounds()
	pixW := int(fyne.Min(float32(scaledW), float32(bounds.Dx()))) // don't push more pixels than we have to
	pixH := int(fyne.Min(float32(scaledH), float32(bounds.Dy()))) // the GL calls will scale this up on GPU.
	scaledBounds := image.Rect(0, 0, pixW, pixH)
	tex := image.NewNRGBA(scaledBounds)
	switch scale {
	case canvas.ImageScalePixels:
		draw.NearestNeighbor.Scale(tex, scaledBounds, pixels, bounds, draw.Over, nil)
	default:
		if scale != canvas.ImageScaleSmooth {
			fyne.LogError("Invalid canvas.ImageScale value, using canvas.ImageScaleSmooth", nil)
		}
		draw.CatmullRom.Scale(tex, scaledBounds, pixels, bounds, draw.Over, nil)
	}
	return tex
}
