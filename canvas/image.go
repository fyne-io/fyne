package canvas

import "image"
import "image/color"

import "github.com/fyne-io/fyne"

// Image describes a raster image area that can render in a Fyne canvas
type Image struct {
	Size     fyne.Size     // The current size of the Image
	Position fyne.Position // The current position of the Image
	Options  Options       // Options to pass to the renderer

	// one of the following sources will provide our image data
	PixelColor func(x, y, w, h int) color.RGBA // Render the image from code
	File       string                          // Load the image froma file
}

// CurrentSize returns the current size of this image object
func (r *Image) CurrentSize() fyne.Size {
	return r.Size
}

// Resize sets a new size for the image object
func (r *Image) Resize(size fyne.Size) {
	r.Size = size
}

// CurrentPosition gets the current position of this image object, relative to it's parent / canvas
func (r *Image) CurrentPosition() fyne.Position {
	return r.Position
}

// Move the image object to a new position, relative to it's parent / canvas
func (r *Image) Move(pos fyne.Position) {
	r.Position = pos
}

// MinSize for a Image simply returns Size{1, 1} as there is no
// explicit content
func (r *Image) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

// NewRaster returns a new Image instance that is rendered dynamically using
// the specified pixelColor function.
// Images returned from this method should draw dynamically to fill the width
// and height parameters passed to pixelColor.
func NewRaster(pixelColor func(x, y, w, h int) color.RGBA) *Image {
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
		PixelColor: func(x, y, w, h int) color.RGBA {
			r, g, b, a := img.At(x, y).RGBA()
			return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
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
func NewImageFromResource(res *fyne.Resource) *Image {
	return NewImageFromFile(res.CachePath())
}
