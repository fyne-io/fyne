package canvas

import (
	"image"
	"image/color"
	"image/draw"
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

// GradientColor ...
type GradientColor struct {
	Start color.Color
	End   color.Color
}

func (gc *GradientColor) linearGradient(w, h, x, y int) *color.RGBA {
	//d := float64(x) / float64(w) // horizontal
	d := float64(x) / float64(w)
	//d := float64(y) / float64(h) // top down

	aR, aG, aB, aA := gc.Start.RGBA()
	bR, bG, bB, bA := gc.End.RGBA()

	dR := d * (float64(bR) - float64(aR))
	dG := d * (float64(bG) - float64(aG))
	dB := d * (float64(bB) - float64(aB))
	dA := d * (float64(bA) - float64(aA))

	//fmt.Printf("[%d][%d] || dR: %f | dG: %f | dB: %f | dA: %f\n", x, y, dR, dG, dB, dA)

	pixel := &color.RGBA{
		R: uint8(float64(aR) + d*dR),
		G: uint8(float64(aG) + d*dG),
		B: uint8(float64(aB) + d*dB),
		A: uint8(float64(aA) + d*dA),
	}

	return pixel

}

type pixelGradient struct {
	g   *Gradient
	img draw.Image
}

// NewRectangleGradient returns a new Image instance that dynamically
// renders a rectangular gradient
func NewRectangleGradient(start color.Color, end color.Color) *Gradient {

	gc := &GradientColor{
		Start: start,
		End:   end,
	}
	pix := &pixelGradient{}
	pix.g = &Gradient{
		Generator: func(w, h int) image.Image {

			if pix.img == nil || pix.img.Bounds().Size().X != w || pix.img.Bounds().Size().Y != h {
				rect := image.Rect(0, 0, w, h)
				pix.img = image.NewRGBA(rect)
			}

			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					pix.img.Set(x, y, gc.linearGradient(w, h, x, y))
				}
			}
			return pix.img
		},
	}
	return pix.g
}
