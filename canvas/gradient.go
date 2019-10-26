package canvas

import (
	"image"
	"image/color"
	"math"
)

// LinearGradient defines a Gradient travelling straight at a given angle.
// The only supported values for the angle are `0.0` (vertical) and `90.0` (horizontal), currently.
type LinearGradient struct {
	baseObject

	StartColor color.Color // The beginning RGBA color of the gradient
	EndColor   color.Color // The end RGBA color of the gradient
	Angle      float64     // The angle of the gradient (0/180 for vertical; 90/270 for horizontal)
}

// Generate calculates an image of the gradient with the specified width and height.
func (g *LinearGradient) Generate(w, h int) image.Image {
	var generator func(x, y, w, h float64) float64
	if g.Angle == 90 {
		// horizontal flipped
		generator = func(x, _, w, _ float64) float64 {
			return ((w - 1) - x) / (w - 1)
		}
	} else if g.Angle == 270 {
		// horizontal
		generator = func(x, _, w, _ float64) float64 {
			return x / (w - 1)
		}
	} else if g.Angle == 45 {
		// diagonal negative flipped
		generator = func(x, y, w, h float64) float64 {
			return math.Abs(((w-1)+(h-1))-(x+((h-1)-y))) / math.Abs((w-1)+(h-1))
		}
	} else if g.Angle == 225 {
		// diagonal negative
		generator = func(x, y, w, h float64) float64 {
			return math.Abs(x+((h-1)-y)) / math.Abs((w-1)+(h-1))
		}
	} else if g.Angle == 135 {
		// diagonal positive flipped
		generator = func(x, y, w, h float64) float64 {
			return math.Abs(((w-1)+(h-1))-(x+y)) / math.Abs((w-1)+(h-1))
		}
	} else if g.Angle == 315 {
		// diagonal positive
		generator = func(x, y, w, h float64) float64 {
			return math.Abs(x+y) / math.Abs((w-1)+(h-1))
		}
	} else if g.Angle == 180 {
		// vertical flipped
		generator = func(_, y, _, h float64) float64 {
			return ((h - 1) - y) / (h - 1)
		}
	} else {
		// vertical
		generator = func(_, y, _, h float64) float64 {
			return y / (h - 1)
		}
	}
	return computeGradient(generator, w, h, g.StartColor, g.EndColor)
}

// RadialGradient defines a Gradient travelling radially from a center point outward.
type RadialGradient struct {
	baseObject

	StartColor color.Color // The beginning RGBA color of the gradient
	EndColor   color.Color // The end RGBA color of the gradient
	// The offset of the center for generation of the gradient.
	// This is not a DP measure but relates to the width/height.
	// A value of 0.5 would move the center by the half width/height.
	CenterOffsetX, CenterOffsetY float64
}

// Generate calculates an image of the gradient with the specified width and height.
func (g *RadialGradient) Generate(w, h int) image.Image {
	generator := func(x, y, w, h float64) float64 {
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
	return computeGradient(generator, w, h, g.StartColor, g.EndColor)
}

func calculatePixel(d float64, startColor, endColor color.Color) *color.RGBA64 {
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

func computeGradient(generator func(x, y, w, h float64) float64, w, h int, startColor, endColor color.Color) image.Image {
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
			distance := generator(float64(x), float64(y), float64(w), float64(h))
			img.Set(x, y, calculatePixel(distance, startColor, endColor))
		}
	}
	return img
}

// NewHorizontalGradient creates a new horizontally travelling linear gradient.
// The start color will be at the left of the gradient and the end color will be at the right.
func NewHorizontalGradient(start, end color.Color) *LinearGradient {
	g := &LinearGradient{StartColor: start, EndColor: end}
	g.Angle = 270
	return g
}

// NewLinearGradient creates a linear gradient at a the specified angle.
// The angle parameter is the degree angle along which the gradient is calculated.
// A NewHorizontalGradient uses 270 degrees and NewVerticalGradient is 0 degrees.
func NewLinearGradient(start, end color.Color, angle float64) *LinearGradient {
	g := &LinearGradient{StartColor: start, EndColor: end}
	g.Angle = angle
	return g
}

// NewRadialGradient creates a new radial gradient.
func NewRadialGradient(start, end color.Color) *RadialGradient {
	return &RadialGradient{StartColor: start, EndColor: end}
}

// NewVerticalGradient creates a new vertically travelling linear gradient.
// The start color will be at the top of the gradient and the end color will be at the bottom.
func NewVerticalGradient(start color.Color, end color.Color) *LinearGradient {
	return &LinearGradient{StartColor: start, EndColor: end}
}
