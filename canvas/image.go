package canvas

import "image"
import "image/color"

import "github.com/fyne-io/fyne"

// Image describes a raster image area that can render in a Fyne canvas
type Image struct {
	baseObject

	// one of the following sources will provide our image data
	PixelColor func(x, y, w, h int) color.RGBA // Render the image from code
	File       string                          // Load the image froma file
	Alpha      float64                         // Set an alpha value < 1.0 to fade the image
}

// NewRaster returns a new Image instance that is rendered dynamically using
// the specified pixelColor function.
// Images returned from this method should draw dynamically to fill the width
// and height parameters passed to pixelColor.
func NewRaster(pixelColor func(x, y, w, h int) color.RGBA) *Image {
	return &Image{
		PixelColor: pixelColor,
		Alpha:      1.0,
	}
}

// NewImageFromImage returns a new Image instance that is rendered from the Go
// image.Image passed in.
// Images returned from this method will map pixel for pixel to the screen
// starting img.Bounds().Min pixels from the top left of the canvas object.
func NewImageFromImage(img image.Image) *Image {
	return &Image{
		PixelColor: func(x, y, w, h int) color.RGBA {
			r, g, b, a := img.At(x, y).RGBA()
			return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		},
		Alpha: 1.0,
	}
}

// NewImageFromFile creates a new image from a local file.
// Images returned from this method will scale to fit the canvas object.
func NewImageFromFile(file string) *Image {
	return &Image{
		File:  file,
		Alpha: 1.0,
	}
}

// NewImageFromResource creates a new image by loading the specified resource.
// Images returned from this method will scale to fit the canvas object.
func NewImageFromResource(res fyne.Resource) *Image {
	return NewImageFromFile(res.CachePath())
}
