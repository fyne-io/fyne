package painter

import (
	"image"
	"math"

	"fyne.io/fyne/canvas"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

// DrawCircle rasterizes the given circle object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
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

// DrawLine rasterizes the given line object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawLine(line *canvas.Line, vectorPad int, scale func(float32) int) *image.RGBA {
	col := line.StrokeColor
	width := scale(float32(line.Size().Width + vectorPad*2))
	height := scale(float32(line.Size().Height + vectorPad*2))
	stroke := scale(line.StrokeWidth)

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(line.Size().Width, line.Size().Height, raw, raw.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(col)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	p1x, p1y := scale(float32(line.Position1.X-line.Position().X+vectorPad)), scale(float32(line.Position1.Y-line.Position().Y+vectorPad))
	p2x, p2y := scale(float32(line.Position2.X-line.Position().X+vectorPad)), scale(float32(line.Position2.Y-line.Position().Y+vectorPad))

	dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
	dasher.Line(rasterx.ToFixedP(float64(p2x), float64(p2y)))
	dasher.Stop(true)
	dasher.Draw()

	return raw
}