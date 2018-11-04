package canvas

import "image"
import "image/color"

import "github.com/fyne-io/fyne"

// Image describes a raster image area that can render in a Fyne canvas
type Image struct {
	baseObject

	// one of the following sources will provide our image data
	PixelColor func(x, y, w, h int) color.Color // Render the image from code
	File       string                           // Load the image froma file

	Translucency float64 // Set a translucency value > 0.0 to fade the image
}

// Alpha is a convenience function that returns the alpha value for an image
// based on it's Translucency value. The result is 1.0 - Translucency.
func (i *Image) Alpha() float64 {
	return 1.0 - i.Translucency
}

// NewRaster returns a new Image instance that is rendered dynamically using
// the specified pixelColor function.
// Images returned from this method should draw dynamically to fill the width
// and height parameters passed to pixelColor.
func NewRaster(pixelColor func(x, y, w, h int) color.Color) *Image {
	return &Image{
		PixelColor: pixelColor,
	}
}

// NewImageFromImage returns a new Image instance that is rendered from the Go
// image.Image passed in.
// Images returned from this method will map pixel for pixel to the screen
// starting img.Bounds().Min pixels from the top left of the canvas object.
func NewImageFromImage(img image.Image) *Image {
	return &Image{
		PixelColor: func(x, y, w, h int) color.Color {
			return img.At(x, y)
		},
	}
}

// NewImageFromFile creates a new image from a local file.
// Images returned from this method will scale to fit the canvas object.
func NewImageFromFile(file string) *Image {
	return &Image{
		File: file,
	}
}

// NewImageFromResource creates a new image by loading the specified resource.
// Images returned from this method will scale to fit the canvas object.
func NewImageFromResource(res fyne.Resource) *Image {
	return NewImageFromFile(res.CachePath())
}
