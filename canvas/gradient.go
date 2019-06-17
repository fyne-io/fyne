package canvas

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne"
)

// Gradient describes a gradient fade between two colors.
type Gradient interface {
	fyne.CanvasObject
	calculateGradient(x, y, w, h float64) float64
	Compute(w, h int) image.Image
}

type baseGradient struct {
	baseObject

	StartColor color.Color // The beginning RGBA color of the gradient
	EndColor   color.Color // The end RGBA color of the gradient
}

var _ Gradient = (*LinearGradient)(nil)
var _ Gradient = (*RadialGradient)(nil)

// LinearGradient defines a Gradient travelling straight at a given angle.
// The only supported values for the angle are `0.0` (vertical) and `90.0` (horizontal), currently.
type LinearGradient struct {
	baseGradient

	Angle float64 // The angle of the gradient (0/180 for vertical; 90/270 for horizontal)
}

func (g *LinearGradient) calculateGradient(x, y, w, h float64) float64 {
	if g.Angle == 90 {
		// horizontal
		return x / w
	}
	// vertical
	return y / h
}

// Compute calculates an image of the gradient with the specified width and height.
func (g *LinearGradient) Compute(w, h int) image.Image {
	return computeGradient(g, w, h, g.StartColor, g.EndColor)
}

// RadialGradient defines a Gradient travelling radially from a center point outward.
type RadialGradient struct {
	baseGradient

	// The offset of the center for generation of the gradient.
	// This is not a DP measure but relates to the width/height.
	// A value of 0.5 would move the center by the half width/height.
	CenterOffsetX, CenterOffsetY float64
}

func (g *RadialGradient) calculateGradient(x, y, w, h float64) float64 {
	// define center plus offset
	centerX := w/2 + w*g.CenterOffsetX
	centerY := h/2 + h*g.CenterOffsetY

	// handle negative offsets
	var a, b float64
	if g.CenterOffsetX < 0 {
		a = w - centerX
	} else {
		a = centerX
	}
	if g.CenterOffsetY < 0 {
		b = h - centerY
	} else {
		b = centerY
	}

	// calculate distance from center for gradient multiplier
	dx, dy := centerX-x, centerY-y
	da := math.Sqrt(dx*dx + dy*dy*a*a/b/b)
	if da > a {
		return 1
	}
	return da / a
}

// Compute calculates an image of the gradient with the specified width and height.
func (g *RadialGradient) Compute(w, h int) image.Image {
	return computeGradient(g, w, h, g.StartColor, g.EndColor)
}

// calculatePixel calculates the pixel at x, y as a gradient between start and end
// using w and h to determine rate of gradation
// returns a color.RGBA64
func calculatePixel(g Gradient, w, h, x, y float64, startColor, endColor color.Color) *color.RGBA64 {
	d := g.calculateGradient(x, y, w, h)

	// fetch RGBA values
	aR, aG, aB, aA := startColor.RGBA()
	bR, bG, bB, bA := endColor.RGBA()

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

func computeGradient(g Gradient, w, h int, startColor, endColor color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	if startColor == nil && endColor == nil {
		return img
	} else if startColor == nil {
		startColor = color.Transparent
	} else if endColor == nil {
		endColor = color.Transparent
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, calculatePixel(g, float64(w), float64(h), float64(x), float64(y), startColor, endColor))
		}
	}
	return img
}

// NewHorizontalGradient creates a new horiziontally travelling linear gradient.
func NewHorizontalGradient(start color.Color, end color.Color) *LinearGradient {
	g := &LinearGradient{baseGradient: baseGradient{StartColor: start, EndColor: end}}
	g.Angle = 90
	return g
}

// NewRadialGradient creates a new radial gradient.
func NewRadialGradient(start color.Color, end color.Color) *RadialGradient {
	return &RadialGradient{baseGradient: baseGradient{StartColor: start, EndColor: end}}
}

// NewVerticalGradient creates a new vertically travelling linear gradient.
func NewVerticalGradient(start color.Color, end color.Color) *LinearGradient {
	return &LinearGradient{baseGradient: baseGradient{StartColor: start, EndColor: end}}
}
