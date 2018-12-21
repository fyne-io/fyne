package canvas

import (
	"image"
	"image/color"

	"github.com/fyne-io/fyne"
)

// ImageFill defines the different type of ways an image can stretch to fill it's space.
type ImageFill int

const (
	// ImageFillStretch will scale the image to match the Size() values.
	// This is the default and does not maintain aspect ratio.
	ImageFillStretch ImageFill = iota
	// ImageFillContain makes the image fit within the object Size(),
	// centrally and maintaining aspect ratio.
	// There may be transparent sections top and bottom or left and right.
	ImageFillContain //(Fit)
)

// Image describes a raster image area that can render in a Fyne canvas
type Image struct {
	baseObject

	// one of the following sources will provide our image data
	File        string                           // Load the image from a file
	PixelColor  func(x, y, w, h int) color.Color // Render the image from code
	PixelAspect float32                          // Set an aspect ratio for pixel based images

	Translucency float64   // Set a translucency value > 0.0 to fade the image
	FillMode     ImageFill // Specify how the image should scale to fill or fit
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

// NewRasterFromImage returns a new Image instance that is rendered from the Go
// image.Image passed in.
// Images returned from this method will map pixel for pixel to the screen
// starting img.Bounds().Min pixels from the top left of the canvas object.
func NewRasterFromImage(img image.Image) *Image {
	return &Image{
		PixelColor: func(x, y, w, h int) color.Color {
			return img.At(x, y)
		},
	}
}

// NewImageFromFile creates a new image from a local file.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromFile(file string) *Image {
	return &Image{
		File: file,
	}
}

// NewImageFromResource creates a new image by loading the specified resource.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromResource(res fyne.Resource) *Image {
	return NewImageFromFile(res.CachePath())
}

// NewImageFromImage returns a new Image instance that is rendered from the Go
// image.Image passed in.
// Images returned from this method will scale to fit the canvas object.
// The method for scaling can be set using the Fill field.
func NewImageFromImage(img image.Image) *Image {
	ret := &Image{}
	ret.PixelColor = func(x, y, w, h int) color.Color {
		size := img.Bounds().Size()
		if size != image.ZP {
			ret.PixelAspect = float32(size.X) / float32(size.Y)
		}

		hScale := float32(size.X) / float32(w)
		vScale := float32(size.Y) / float32(h)

		return img.At(int(float32(x)*hScale), int(float32(y)*vScale))
	}

	return ret
}
