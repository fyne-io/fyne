package painter

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

const quarterCircleControl = 1 - 0.55228

// DrawArc rasterizes the given arc object into an image.
// The scale function is used to understand how many pixels are required per unit of size.
// The arc is drawn from StartAngle to EndAngle (in degrees).
// 0°/360 is right, 90° is top, 180° is left, 270° is bottom
// 0°/-360 is right, -90° is bottom, -180° is left, -270° is top
func DrawArc(arc *canvas.Arc, vectorPad float32, scale func(float32) float32) *image.RGBA {
	size := arc.Size()
	radius := fyne.Min(size.Width, size.Height) / 2

	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))
	stroke := scale(arc.StrokeWidth)
	// cornerRadius := float64(scale(arc.CornerRadius)) // TODO rounded corners have not been implemented yet

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	outerRadius := float64(scale(radius))
	innerRadius := float64(scale(arc.InnerRadius))

	// Convert to radians
	startRad := float64(arc.StartAngle) * (math.Pi / 180.0)
	endRad := float64(arc.EndAngle) * (math.Pi / 180.0)

	if arc.EndAngle < arc.StartAngle {
		// Ensure always draw counter-clockwise
		startRad, endRad = endRad, startRad
	}

	angleDiff := endRad - startRad

	// Normalize angleDiff to [-2π, 2π]
	angleDiff = math.Mod(angleDiff+2*math.Pi, 2*math.Pi)
	if angleDiff == 0 && arc.StartAngle != arc.EndAngle {
		angleDiff = 2 * math.Pi // full circle
	}

	if math.Abs(angleDiff) < 1e-6 {
		return raw // empty
	}

	// Avoid full circle becoming zero-length
	if math.Abs(angleDiff) >= 2*math.Pi {
		angleDiff = 2*math.Pi - 1e-6
	}

	endRad = startRad + angleDiff

	if arc.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(arc.FillColor)

		point := func(r, angle float64) (x, y float64) {
			x = centerX + r*math.Cos(angle)
			y = centerY - r*math.Sin(angle)
			return
		}

		startX, startY := point(innerRadius, startRad)
		filler.Start(rasterx.ToFixedP(startX, startY))

		outerStartX, outerStartY := point(outerRadius, startRad)
		filler.Line(rasterx.ToFixedP(outerStartX, outerStartY))
		outerEndX, outerEndY := point(outerRadius, endRad)
		outerArc := []float64{
			outerRadius, outerRadius, 0, 0, 0,
			outerEndX, outerEndY,
		}
		rasterx.AddArc(outerArc, centerX, centerY, outerStartX, outerStartY, filler)

		innerEndX, innerEndY := point(innerRadius, endRad)
		filler.Line(rasterx.ToFixedP(innerEndX, innerEndY))
		innerArc := []float64{
			innerRadius, innerRadius, 0, 1, 1,
			startX, startY,
		}
		rasterx.AddArc(innerArc, centerX, centerY, innerEndX, innerEndY, filler)

		filler.Stop(true)
		filler.Draw()
	}

	if arc.StrokeColor != nil && arc.StrokeWidth > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(arc.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)

		point := func(r, angle float64) (x, y float64) {
			x = centerX + r*math.Cos(angle)
			y = centerY - r*math.Sin(angle)
			return
		}

		// Outer arc
		outerStartX, outerStartY := point(outerRadius, startRad)
		outerEndX, outerEndY := point(outerRadius, endRad)
		dasher.Start(rasterx.ToFixedP(outerStartX, outerStartY))
		outerArc := []float64{
			outerRadius, outerRadius, 0, 0, 0,
			outerEndX, outerEndY,
		}
		rasterx.AddArc(outerArc, centerX, centerY, outerStartX, outerStartY, dasher)

		// If it's a ring, draw the inner arc and the connecting lines
		fullCircle := math.Abs(angleDiff-2*math.Pi) < 1e-4

		innerEndX, innerEndY := point(innerRadius, endRad)
		innerStartX, innerStartY := point(innerRadius, startRad)

		// Only draw connecting lines if not full circle
		if !fullCircle {
			dasher.Line(rasterx.ToFixedP(innerEndX, innerEndY))
		} else {
			dasher.Start(rasterx.ToFixedP(innerEndX, innerEndY))
		}

		// Inner arc (reverse direction)
		innerArc := []float64{
			innerRadius, innerRadius, 0, 1, 1,
			innerStartX, innerStartY,
		}
		rasterx.AddArc(innerArc, centerX, centerY, innerEndX, innerEndY, dasher)

		dasher.Stop(true)
		dasher.Draw()
	}

	return raw
}

// DrawCircle rasterizes the given circle object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawCircle(circle *canvas.Circle, vectorPad float32, scale func(float32) float32) *image.RGBA {
	size := circle.Size()
	radius := fyne.Min(size.Width, size.Height) / 2

	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))
	stroke := scale(circle.StrokeWidth)

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())

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
	size := line.Size()
	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))
	stroke := scale(line.StrokeWidth)
	if stroke < 1 { // software painter doesn't fade lines to compensate
		stroke = 1
	}

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(col)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	position := line.Position()
	p1x, p1y := scale(line.Position1.X-position.X+vectorPad), scale(line.Position1.Y-position.Y+vectorPad)
	p2x, p2y := scale(line.Position2.X-position.X+vectorPad), scale(line.Position2.Y-position.Y+vectorPad)

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
func DrawRectangle(rect *canvas.Rectangle, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	return drawOblong(rect, rect.FillColor, rect.StrokeColor, rect.StrokeWidth, rect.CornerRadius, rect.Aspect, rWidth, rHeight, vectorPad, scale)
}

func DrawSquare(sq *canvas.Square, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	return drawOblong(sq, sq.FillColor, sq.StrokeColor, sq.StrokeWidth, sq.CornerRadius, 1.0, rWidth, rHeight, vectorPad, scale)
}

func drawOblong(obj fyne.CanvasObject, fill, strokeCol color.Color, strokeWidth, radius, aspect float32, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	width := int(scale(rWidth + vectorPad*2))
	height := int(scale(rHeight + vectorPad*2))
	stroke := scale(strokeWidth)

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(rWidth), int(rHeight), raw, raw.Bounds())

	scaledPad := scale(vectorPad)
	p1x, p1y := scaledPad, scaledPad
	p2x, p2y := scale(rWidth)+scaledPad, scaledPad
	p3x, p3y := scale(rWidth)+scaledPad, scale(rHeight)+scaledPad
	p4x, p4y := scaledPad, scale(rHeight)+scaledPad

	if fill != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(fill)
		if radius == 0 {
			rasterx.AddRect(float64(p1x), float64(p1y), float64(p3x), float64(p3y), 0, filler)
		} else {
			r := float64(scale(radius))
			rasterx.AddRoundRect(float64(p1x), float64(p1y), float64(p3x), float64(p3y), r, r, 0, rasterx.RoundGap, filler)
		}
		filler.Draw()
	}

	if strokeCol != nil && strokeWidth > 0 {
		r := scale(radius)
		c := quarterCircleControl * r
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(strokeCol)
		dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
		if c != 0 {
			dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y+r)))
			dasher.CubeBezier(rasterx.ToFixedP(float64(p1x), float64(p1y+c)), rasterx.ToFixedP(float64(p1x+c), float64(p1y)), rasterx.ToFixedP(float64(p1x+r), float64(p2y)))
		} else {
			dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
		}
		dasher.Line(rasterx.ToFixedP(float64(p2x-r), float64(p2y)))
		if c != 0 {
			dasher.CubeBezier(rasterx.ToFixedP(float64(p2x-c), float64(p2y)), rasterx.ToFixedP(float64(p2x), float64(p2y+c)), rasterx.ToFixedP(float64(p2x), float64(p2y+r)))
		}
		dasher.Line(rasterx.ToFixedP(float64(p3x), float64(p3y-r)))
		if c != 0 {
			dasher.CubeBezier(rasterx.ToFixedP(float64(p3x), float64(p3y-c)), rasterx.ToFixedP(float64(p3x-c), float64(p3y)), rasterx.ToFixedP(float64(p3x-r), float64(p3y)))
		}
		dasher.Line(rasterx.ToFixedP(float64(p4x+r), float64(p4y)))
		if c != 0 {
			dasher.CubeBezier(rasterx.ToFixedP(float64(p4x+c), float64(p4y)), rasterx.ToFixedP(float64(p4x), float64(p4y-c)), rasterx.ToFixedP(float64(p4x), float64(p4y-r)))
		}
		dasher.Stop(true)
		dasher.Draw()
	}

	return raw
}
