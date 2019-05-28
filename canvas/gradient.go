package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/davecgh/go-spew/spew"
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

// LinearGradientColor defines start and end color
type LinearGradientColor struct {
	Start       color.Color
	End         color.Color
	GradientFnc func(x, y, w, h int) float64
}

func (gc *LinearGradientColor) gradient(w, h, x, y int) *color.RGBA64 {
	d := gc.GradientFnc(x, y, w, h)
	// fetch RGBA values
	aR, aG, aB, aA := gc.Start.RGBA()
	bR, bG, bB, bA := gc.End.RGBA()

	// Get difference
	dR := (float64(bR) - float64(aR))
	dG := (float64(bG) - float64(aG))
	dB := (float64(bB) - float64(aB))
	dA := (float64(bA) - float64(aA))

	// Return with applied gradiation
	pixel := &color.RGBA64{
		R: uint16(float64(aR) + d*dR),
		B: uint16(float64(aB) + d*dB),
		G: uint16(float64(aG) + d*dG),
		A: uint16(float64(aA) + d*dA),
	}

	// Prevent rendering dead space
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
func NewRectangleGradient(optFnc ...GradientOption) *Gradient {
	options := &GradientOptions{}
	for _, opt := range optFnc {
		opt(options)
	}

	l := LinearGradientColor{
		Start: options.StartColor,
		End:   options.EndColor,
	}

	pix := &pixelGradient{}

	spew.Dump(options)

	switch options.Direction {
	case HORIZONTAL:
		l.GradientFnc = linearHorizontal
	case VERTICAL:
		l.GradientFnc = linearVertical
	case CIRCULAR:
		l.GradientFnc = linearCircular
	}

	pix.g = &Gradient{
		Generator: func(w, h int) image.Image {
			if pix.img == nil || pix.img.Bounds().Size().X != w || pix.img.Bounds().Size().Y != h {
				rect := image.Rect(0, 0, w, h)
				pix.img = image.NewRGBA(rect)
			}
			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					pix.img.Set(x, y, l.gradient(w, h, x, y))
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
	Direction   GradientDirection
	StartColor  color.Color
	EndColor    color.Color
	GradientFnc func(x, y, w, h int) float64
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

// GradientFunction attaches a function to determine gradiation
func GradientFunction(fnc func(x, y, w, h int) float64) GradientOption {
	return func(opt *GradientOptions) {
		opt.GradientFnc = fnc
	}
}

func linearHorizontal(x, y, w, h int) float64 {
	return float64(x) / float64(w)
}

func linearVertical(x, y, w, h int) float64 {
	return float64(y) / float64(h)
}

func linearCircular(x, y, w, h int) float64 {
	centerX := w / 2
	centerY := h / 2
	dx, dy := float64(centerX-x), float64(centerY-y)
	d := math.Sqrt(dx*dx + dy*dy)
	return d / 255

}
