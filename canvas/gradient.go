package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne"
)

// Gradient describes a gradient between two colors
type Gradient struct {
	baseObject

	Generator func(w, h int) image.Image

	Translucency float64 // Set a translucency value > 0.0 to fade the raster

	img draw.Image // internal cache for pixel generator
}

// Alpha is a convenience function that returns the alpha value for a raster
// based on its Translucency value. The result is 1.0 - Translucency.
func (g *Gradient) Alpha() float64 {
	return 1.0 - g.Translucency
}

// gradient generates a pixel using the defined gradient
// function by using the w, h, and current x,y  of this pixel
func gradient(gradientFnc func(x, y, ox, oy, w, h int) float64, start color.Color, end color.Color, w, h, ox, oy, x, y int) *color.RGBA64 {
	d := gradientFnc(x, y, ox, oy, w, h)
	// fetch RGBA values
	aR, aG, aB, aA := start.RGBA()
	bR, bG, bB, bA := end.RGBA()

	// Get difference
	dR := (float64(bR) - float64(aR))
	dG := (float64(bG) - float64(aG))
	dB := (float64(bB) - float64(aB))
	dA := (float64(bA) - float64(aA))

	// Apply graduation
	pixel := &color.RGBA64{
		R: uint16(float64(aR) + d*dR),
		B: uint16(float64(aB) + d*dB),
		G: uint16(float64(aG) + d*dG),
		A: uint16(float64(aA) + d*dA),
	}

	// Prevent rendering dead space - transparent corners
	if d > 1 {
		pixel.A = uint16(0)
	}

	return pixel
}

type pixelGradient struct {
	g   *Gradient
	img draw.Image
}

// NewRectangleGradient returns a new Image instance that dynamically
// renders a gradient based off of options
// Accepts
func NewRectangleGradient(optFnc ...GradientOption) *Gradient {
	options := &GradientOptions{}
	for _, opt := range optFnc {
		opt(options)
	}

	pix := &pixelGradient{}

	var gradFnc func(x, y, ox, oy, w, h int) float64

	// Select linear function for appropriate type of gradient
	switch options.Direction {
	case HORIZONTAL:
		gradFnc = linearHorizontal
	case VERTICAL:
		gradFnc = linearVertical
	case CIRCULAR:
		gradFnc = linearCircular
	default:
		gradFnc = linearHorizontal
	}

	pix.g = &Gradient{
		Generator: func(w, h int) image.Image {
			if pix.img == nil || pix.img.Bounds().Size().X != w || pix.img.Bounds().Size().Y != h {
				rect := image.Rect(0, 0, w, h)
				pix.img = image.NewRGBA(rect)
			}

			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					pix.img.Set(x, y, gradient(gradFnc, options.StartColor, options.EndColor, w, h, options.Offset.X, options.Offset.Y, x, y))

				}
			}
			return pix.img
		},
	}

	return pix.g
}

// GradientDirection defines a starting or stopping location
type GradientDirection int

// Our cardinal Targets
const (
	VERTICAL GradientDirection = iota
	HORIZONTAL
	CIRCULAR
)

// GradientOptions options for configuring a new LinearGradient
type GradientOptions struct {
	Direction  GradientDirection
	StartColor color.Color
	EndColor   color.Color
	Offset     fyne.Position
}

// GradientOption API for configuring a new gradient
type GradientOption func(*GradientOptions)

// GradientAlignment attachs a start position
func GradientAlignment(dir GradientDirection) GradientOption {
	return func(opt *GradientOptions) {
		opt.Direction = dir
	}
}

// GradientStartColor attaches a starting color
func GradientStartColor(start color.Color) GradientOption {
	return func(opt *GradientOptions) {
		opt.StartColor = start
	}
}

// GradientEndColor attaches an ending color
func GradientEndColor(end color.Color) GradientOption {
	return func(opt *GradientOptions) {
		opt.EndColor = end
	}
}

// GradientOffset adds an offset for gradient start point
func GradientOffset(p fyne.Position) GradientOption {
	return func(opt *GradientOptions) {
		opt.Offset = p
	}
}

// Linear horizontal gradiant
func linearHorizontal(x, y, ox, oy, w, h int) float64 {
	return float64(x) / float64(w)
}

// Linear vertical gradiant
func linearVertical(x, y, ox, oy, w, h int) float64 {
	return float64(y) / float64(h)
}

// Linear circular gradient - function/math thanks to Tilo Prützs
func linearCircular(x, y, ox, oy, w, h int) float64 {
	centerX := float64(w)/2 + float64(ox)
	centerY := float64(h)/2 + float64(oy)
	var rx, ry float64
	if ox < 0 {
		rx = float64(w) - centerX
	} else {
		rx = centerX
	}
	if oy < 0 {
		ry = float64(h) - centerY
	} else {
		ry = centerY
	}
	_ = ry
	dx, dy := centerX-float64(x), centerY-float64(y)
	d := math.Sqrt(dx*dx + dy*dy)
	if d > rx {
		return 1
	}
	return d / rx
}
