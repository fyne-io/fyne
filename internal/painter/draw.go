package painter

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

// DrawCircle rasterizes the given circle object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawCircle(circle *canvas.Circle, vectorPad float32, scale func(float32) float32) *image.RGBA {
	radius := fyne.Min(circle.Size().Width, circle.Size().Height) / 2

	width := int(scale(circle.Size().Width + vectorPad*2))
	height := int(scale(circle.Size().Height + vectorPad*2))
	stroke := scale(circle.StrokeWidth)

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(circle.Size().Width), int(circle.Size().Height), raw, raw.Bounds())

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
func DrawLine(line *canvas.Line, vectorPad float32, scale func(float32) float32) *image.RGBA {
	col := line.StrokeColor
	width := int(scale(line.Size().Width + vectorPad*2))
	height := int(scale(line.Size().Height + vectorPad*2))
	stroke := scale(line.StrokeWidth)
	if stroke < 1 { // software painter doesn't fade lines to compensate
		stroke = 1
	}

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(line.Size().Width), int(line.Size().Height), raw, raw.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(col)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	p1x, p1y := scale(line.Position1.X-line.Position().X+vectorPad), scale(line.Position1.Y-line.Position().Y+vectorPad)
	p2x, p2y := scale(line.Position2.X-line.Position().X+vectorPad), scale(line.Position2.Y-line.Position().Y+vectorPad)

	if stroke <= 1.5 { // adjust to support 1px
		if p1x == p2x {
			p1x -= 0.5
			p2x -= 0.5
		}
		if p1y == p2y {
			p1y -= 0.5
			p2y -= 0.5
		}
	}

	dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
	dasher.Line(rasterx.ToFixedP(float64(p2x), float64(p2y)))
	dasher.Stop(true)
	dasher.Draw()

	return raw
}

// DrawRectangle rasterizes the given rectangle object with stroke border into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawRectangle(rect *canvas.Rectangle, vectorPad float32, scale func(float32) float32) *image.RGBA {
	width := int(scale(rect.Size().Width + vectorPad*2))
	height := int(scale(rect.Size().Height + vectorPad*2))
	stroke := scale(rect.StrokeWidth)

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(rect.Size().Width), int(rect.Size().Height), raw, raw.Bounds())

	scaledPad := scale(vectorPad)
	p1x, p1y := scaledPad, scaledPad
	p2x, p2y := scale(rect.Size().Width)+scaledPad, scaledPad
	p3x, p3y := scale(rect.Size().Width)+scaledPad, scale(rect.Size().Height)+scaledPad
	p4x, p4y := scaledPad, scale(rect.Size().Height)+scaledPad

	if rect.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(rect.FillColor)
		rasterx.AddRect(float64(p1x), float64(p1y), float64(p3x), float64(p3y), 0, filler)
		filler.Draw()
	}

	if rect.StrokeColor != nil && rect.StrokeWidth > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(rect.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
		dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
		dasher.Line(rasterx.ToFixedP(float64(p2x), float64(p2y)))
		dasher.Line(rasterx.ToFixedP(float64(p3x), float64(p3y)))
		dasher.Line(rasterx.ToFixedP(float64(p4x), float64(p4y)))
		dasher.Stop(true)
		dasher.Draw()
	}

	return raw
}
