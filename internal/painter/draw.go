package painter

import (
	"image"
	"math"

	"fyne.io/fyne/canvas"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

func DrawCircle(circle *canvas.Circle, vectorPad int, scale func(float32) int) *image.RGBA {
	radius := float32(math.Min(float64(circle.Size().Width), float64(circle.Size().Height)) / 2)

	width := scale(float32(circle.Size().Width + vectorPad*2))
	height := scale(float32(circle.Size().Height + vectorPad*2))
	stroke := scale(circle.StrokeWidth)

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(circle.Size().Width, circle.Size().Height, raw, raw.Bounds())

	if circle.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(circle.FillColor)
		rasterx.AddCircle(float64(width/2), float64(height/2), float64(scale(radius)), filler)
		filler.Draw()
	}

	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(circle.StrokeColor)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	rasterx.AddCircle(float64(width/2), float64(height/2), float64(scale(radius)), dasher)
	dasher.Draw()

	return raw
}
