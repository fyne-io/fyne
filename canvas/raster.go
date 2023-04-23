package canvas

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Raster)(nil)

// Raster describes a raster image area that can render in a Fyne canvas
type Raster struct {
	baseObject

	// Render the raster image from code
	Generator func(w, h int) image.Image

	// Set a translucency value > 0.0 to fade the raster
	Translucency float64
	// Specify the type of scaling interpolation applied to the raster if it is not full-size
	// Since: 1.4.1
	ScaleMode ImageScale
}

// Alpha is a convenience function that returns the alpha value for a raster
// based on its Translucency value. The result is 1.0 - Translucency.
func (r *Raster) Alpha() float64 {
	return 1.0 - r.Translucency
}

// Hide will set this raster to not be visible
func (r *Raster) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the raster to a new position, relative to its parent / canvas
func (r *Raster) Move(pos fyne.Position) {
	r.baseObject.Move(pos)

	repaint(r)
}

// Resize on a raster image causes the new size to be set and then calls Refresh.
// This causes the underlying data to be recalculated and a new output to be drawn.
func (r *Raster) Resize(s fyne.Size) {
	if s == r.Size() {
		return
	}

	r.baseObject.Resize(s)
	Refresh(r)
}

// Refresh causes this raster to be redrawn with its configured state.
func (r *Raster) Refresh() {
	Refresh(r)
}

// NewRaster returns a new Image instance that is rendered dynamically using
// the specified generate function.
// Images returned from this method should draw dynamically to fill the width
// and height parameters passed to pixelColor.
func NewRaster(generate func(w, h int) image.Image) *Raster {
	return &Raster{Generator: generate}
}

type pixelRaster struct {
	r *Raster

	img draw.Image
}

// NewRasterWithPixels returns a new Image instance that is rendered dynamically
// by iterating over the specified pixelColor function for each x, y pixel.
// Images returned from this method should draw dynamically to fill the width
// and height parameters passed to pixelColor.
func NewRasterWithPixels(pixelColor func(x, y, w, h int) color.Color) *Raster {
	pix := &pixelRaster{}
	pix.r = &Raster{
		Generator: func(w, h int) image.Image {
			if pix.img == nil || pix.img.Bounds().Size().X != w || pix.img.Bounds().Size().Y != h {
				// raster first pixel, figure out color type
				var dst draw.Image
				rect := image.Rect(0, 0, w, h)
				switch pixelColor(0, 0, w, h).(type) {
				case color.Alpha:
					dst = image.NewAlpha(rect)
				case color.Alpha16:
					dst = image.NewAlpha16(rect)
				case color.CMYK:
					dst = image.NewCMYK(rect)
				case color.Gray:
					dst = image.NewGray(rect)
				case color.Gray16:
					dst = image.NewGray16(rect)
				case color.NRGBA:
					dst = image.NewNRGBA(rect)
				case color.NRGBA64:
					dst = image.NewNRGBA64(rect)
				case color.RGBA:
					dst = image.NewRGBA(rect)
				case color.RGBA64:
					dst = image.NewRGBA64(rect)
				default:
					dst = image.NewRGBA(rect)
				}
				pix.img = dst
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					pix.img.Set(x, y, pixelColor(x, y, w, h))
				}
			}

			return pix.img
		},
	}
	return pix.r
}

type subImg interface {
	SubImage(r image.Rectangle) image.Image
}

// NewRasterFromImage returns a new Raster instance that is rendered from the Go
// image.Image passed in.
// Rasters returned from this method will map pixel for pixel to the screen
// starting img.Bounds().Min pixels from the top left of the canvas object.
// Truncates rather than scales the image.
// If smaller than the target space, the image will be padded with zero-pixels to the target size.
func NewRasterFromImage(img image.Image) *Raster {
	return &Raster{
		Generator: func(w int, h int) image.Image {
			bounds := img.Bounds()

			rect := image.Rect(0, 0, w, h)

			switch {
			case w == bounds.Max.X && h == bounds.Max.Y:
				return img
			case w >= bounds.Max.X && h >= bounds.Max.Y:
				// try quickly truncating
				if sub, ok := img.(subImg); ok {
					return sub.SubImage(image.Rectangle{
						Min: bounds.Min,
						Max: image.Point{
							X: bounds.Min.X + w,
							Y: bounds.Min.Y + h,
						},
					})
				}
			default:
				if !rect.Overlaps(bounds) {
					return image.NewUniform(color.RGBA{})
				}
				bounds = bounds.Intersect(rect)
			}

			// respect the user's pixel format (if possible)
			var dst draw.Image
			switch i := img.(type) {
			case *image.Alpha:
				dst = image.NewAlpha(rect)
			case *image.Alpha16:
				dst = image.NewAlpha16(rect)
			case *image.CMYK:
				dst = image.NewCMYK(rect)
			case *image.Gray:
				dst = image.NewGray(rect)
			case *image.Gray16:
				dst = image.NewGray16(rect)
			case *image.NRGBA:
				dst = image.NewNRGBA(rect)
			case *image.NRGBA64:
				dst = image.NewNRGBA64(rect)
			case *image.Paletted:
				dst = image.NewPaletted(rect, i.Palette)
			case *image.RGBA:
				dst = image.NewRGBA(rect)
			case *image.RGBA64:
				dst = image.NewRGBA64(rect)
			default:
				dst = image.NewRGBA(rect)
			}

			draw.Draw(dst, bounds, img, bounds.Min, draw.Over)
			return dst
		},
	}
}
