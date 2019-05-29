package canvas

import (
	"errors"
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

	Generator func(w, h int) image.Image

	direction GradientDirection // The direction of the gradient.
	start     color.Color       // The beginning RGBA color of the gradient
	end       color.Color       // The end RGBA color of the gradient
	offset    fyne.Position     // The offest for generation of the gradient

	gradientFnc func(x, y, w, h int, ox, oy float64) float64 // The function for calculating the gradient

	img draw.Image // internal cache for pixel generator - may be superfluous
}

// Offset returns the offset of the  gradient
func (g *Gradient) Offset() fyne.Position {
	return g.offset
}

// SetOffset sets the gradient offset.
// Should be set before generation
func (g *Gradient) SetOffset(pos fyne.Position) {
	g.offset = pos
}

// calculatePixel uses the gradientFnc to caculate the pixel
// at x, y as a gradient between start and end
// using w and h to determine rate of graduation
// returns a color.RGBA64
func (g *Gradient) calculatePixel(w, h, x, y int) *color.RGBA64 {
	ox, oy := float64(g.offset.X), float64(g.offset.Y)
	d := g.gradientFnc(x, y, w, h, ox, oy)

	// fetch RGBA values
	aR, aG, aB, aA := g.start.RGBA()
	bR, bG, bB, bA := g.end.RGBA()

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

// pixelGradient is used during construction
type pixelGradient struct {
	g   *Gradient
	img draw.Image // image cache during generator run
}

// NewLinearGradient returns a new Image instance that dynamically
// renders a gradient of type GradientDirection
func NewLinearGradient(start color.Color, end color.Color, direction GradientDirection) (*Gradient, error) {
	if start == nil || end == nil {
		return nil, errors.New("start and end colors are required for a gradient")
	}

	pix := &pixelGradient{}

	pix.g = &Gradient{
		start:  start,
		end:    end,
		offset: fyne.NewPos(0, 0),
	}

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

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				pix.img.Set(x, y, pix.g.calculatePixel(w, h, x, y))

			}
		}
		return pix.img
	}
	return pix.g, nil
}

/*
	Gradient multiplier calculations
*/

// Linear horizontal gradiant
func linearHorizontal(x, y, w, h int, ox, oy float64) float64 {
	return float64(x) / float64(w)
}

// Linear vertical gradiant
func linearVertical(x, y, w, h int, ox, oy float64) float64 {
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
