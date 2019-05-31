package canvas

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne"
)

const (
	// GradientDirectionVertical defines a vertically traveling gradient
	GradientDirectionVertical GradientDirection = iota
	// GradientDirectionHorizontal defines a horizontally traveling gradient
	GradientDirectionHorizontal
	// GradientDirectionCircular defines a gradient traveling radially from a center point outward
	GradientDirectionCircular
)

// GradientDirection defines a starting or stopping location
type GradientDirection int

// Gradient describes a gradient fade between two colors
type Gradient struct {
	baseObject
	// Generator is a func used by driver to generate and render. Inadvisable to use directly.
	Generator func(w, h int) image.Image

	Direction    GradientDirection // The direction of the gradient.
	StartColor   color.Color       // The beginning RGBA color of the gradient
	EndColor     color.Color       // The end RGBA color of the gradient
	CenterOffset fyne.Position     // The offset for generation of the gradient

	gradientFnc func(x, y, w, h int, ox, oy float64) float64 // The function for calculating the gradient

	img draw.Image // internal cache for pixel generator - may be superfluous
}

// calculatePixel uses the gradientFnc to calculate the pixel
// at x, y as a gradient between start and end
// using w and h to determine rate of gradation
// returns a color.RGBA64
func (g *Gradient) calculatePixel(w, h, x, y int) *color.RGBA64 {
	ox, oy := float64(g.CenterOffset.X), float64(g.CenterOffset.Y)
	d := g.gradientFnc(x, y, w, h, ox, oy)

	// fetch RGBA values
	aR, aG, aB, aA := g.StartColor.RGBA()
	bR, bG, bB, bA := g.EndColor.RGBA()

	// Get difference
	dR := float64(bR) - float64(aR)
	dG := float64(bG) - float64(aG)
	dB := float64(bB) - float64(aB)
	dA := float64(bA) - float64(aA)

	// Apply gradations
	pixel := &color.RGBA64{
		R: uint16(float64(aR) + d*dR),
		B: uint16(float64(aB) + d*dB),
		G: uint16(float64(aG) + d*dG),
		A: uint16(float64(aA) + d*dA),
	}

	return pixel

}

// pixelGradient is used during construction
type pixelGradient struct {
	g   *Gradient
	img draw.Image // image cache during generator run
}

// NewLinearGradient returns a new Image instance that dynamically
// renders a gradient of type GradientDirection
func NewLinearGradient(start color.Color, end color.Color, direction GradientDirection) *Gradient {
	pix := &pixelGradient{}

	pix.g = &Gradient{
		StartColor: start,
		EndColor:   end,
	}

	pix.g.Direction = direction
	switch direction {
	case GradientDirectionHorizontal:
		pix.g.gradientFnc = linearHorizontal
	case GradientDirectionVertical:
		pix.g.gradientFnc = linearVertical
	case GradientDirectionCircular:
		pix.g.gradientFnc = linearCircular
	default:
		pix.g.gradientFnc = linearHorizontal
	}

	pix.g.Generator = func(w, h int) image.Image {
		if pix.img == nil || pix.img.Bounds().Size().X != w || pix.img.Bounds().Size().Y != h {
			rect := image.Rect(0, 0, w, h)
			pix.img = image.NewRGBA(rect)
		}

		if pix.g.StartColor == nil && pix.g.EndColor == nil {
			return pix.img
		} else if pix.g.StartColor == nil {
			pix.g.StartColor = color.Transparent
		} else if pix.g.EndColor == nil {
			pix.g.EndColor = color.Transparent
		}

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				pix.img.Set(x, y, pix.g.calculatePixel(w, h, x, y))

			}
		}
		return pix.img
	}
	return pix.g
}

/*
	Gradient multiplier calculations
*/

// Linear horizontal gradient
func linearHorizontal(x, _, w, _ int, _, _ float64) float64 {
	return float64(x) / float64(w)
}

// Linear vertical gradient
func linearVertical(_, y, _, h int, _, _ float64) float64 {
	return float64(y) / float64(h)
}

// Linear circular gradient - function/math thanks to Tilo Prützs
func linearCircular(x, y, w, h int, ox, oy float64) float64 {
	// define center plus offset
	centerX := float64(w)/2 + ox
	centerY := float64(h)/2 + oy

	// handle negative offsets
	var a, b float64
	if ox < 0 {
		a = float64(w) - centerX
	} else {
		a = centerX
	}
	if oy < 0 {
		b = float64(h) - centerY
	} else {
		b = centerY
	}

	// calculate distance from center for gradient multiplier
	dx, dy := centerX-float64(x), centerY-float64(y)
	da := math.Sqrt(dx*dx + dy*dy*a*a/b/b)
	if da > a {
		return 1
	}
	return da / a
}
