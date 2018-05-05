package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

// Image describes a raster image area that can render in a Fyne canvas
type Image struct {
	Size     ui.Size     // The current size of the Image
	Position ui.Position // The current position of the Image

	PixelColor func(x, y, w, h int) color.RGBA
}

// CurrentSize returns the current size of this image object
func (r *Image) CurrentSize() ui.Size {
	return r.Size
}

// Resize sets a new size for the image object
func (r *Image) Resize(size ui.Size) {
	r.Size = size
}

// CurrentPosition gets the current position of this image object, relative to it's parent / canvas
func (r *Image) CurrentPosition() ui.Position {
	return r.Position
}

// Move the image object to a new position, relative to it's parent / canvas
func (r *Image) Move(pos ui.Position) {
	r.Position = pos
}

// MinSize for a Image simply returns Size{1, 1} as there is no
// explicit content
func (r *Image) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

// NewImage returns a new Image instance
func NewRaster(pixelColor func(x, y, w, h int) color.RGBA) *Image {
	return &Image{
		PixelColor: pixelColor,
	}
}
