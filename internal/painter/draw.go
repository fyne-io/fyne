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

	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())

	centerX := float64(width) / 2
	centerY := float64(height) / 2

	outerRadius := float64(scale(fyne.Min(size.Width, size.Height) / 2.0))
	innerRadius := float64(scale(arc.InnerRadius))
	if innerRadius < 0 {
		innerRadius = 0
	}
	if innerRadius > outerRadius {
		innerRadius = outerRadius
	}

	// convert to radians
	startRad := float64(arc.StartAngle * math.Pi / 180.0)
	endRad := float64(arc.EndAngle * math.Pi / 180.0)
	sweep := endRad - startRad
	if sweep == 0 {
		// nothing to draw
		return raw
	}

	if sweep > 2*math.Pi {
		sweep = 2 * math.Pi
	} else if sweep < -2*math.Pi {
		sweep = -2 * math.Pi
	}

	cornerRadius := float64(scale(arc.CornerRadius))
	if arc.CornerRadius == canvas.RadiusMaximum {
		// height (thickness), width (length)
		thickness := outerRadius - innerRadius
		span := math.Sin(0.5 * math.Min(math.Abs(sweep), math.Pi)) // span in (0,1)
		length := 1.5 * outerRadius * span / (1 + span)            // no division-by-zero risk

		cornerRadius = float64(GetMaximumRadius(fyne.NewSize(
			float32(thickness), float32(length),
		)))
	}

	if arc.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(arc.FillColor)
		// rasterx.AddArc is not used because it does not support rounded corners
		drawRoundArc(filler, centerX, centerY, outerRadius, innerRadius, startRad, sweep, cornerRadius)
		filler.Draw()
	}

	stroke := float64(scale(arc.StrokeWidth))
	if arc.StrokeColor != nil && stroke > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(arc.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(stroke*64), 0, nil, nil, nil, 0, nil, 0)
		// rasterx.AddArc is not used because it does not support rounded corners
		drawRoundArc(dasher, centerX, centerY, outerRadius, innerRadius, startRad, sweep, cornerRadius)
		dasher.Draw()
	}

	return raw
}

// DrawCircle rasterizes the given circle object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawCircle(circle *canvas.Circle, vectorPad float32, scale func(float32) float32) *image.RGBA {
	size := circle.Size()
	radius := GetMaximumRadius(size)

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
	topRightRadius := GetCornerRadius(rect.TopRightCornerRadius, rect.CornerRadius)
	topLeftRadius := GetCornerRadius(rect.TopLeftCornerRadius, rect.CornerRadius)
	bottomRightRadius := GetCornerRadius(rect.BottomRightCornerRadius, rect.CornerRadius)
	bottomLeftRadius := GetCornerRadius(rect.BottomLeftCornerRadius, rect.CornerRadius)
	return drawOblong(rect.FillColor, rect.StrokeColor, rect.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rWidth, rHeight, vectorPad, scale)
}

func DrawSquare(sq *canvas.Square, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	topRightRadius := GetCornerRadius(sq.TopRightCornerRadius, sq.CornerRadius)
	topLeftRadius := GetCornerRadius(sq.TopLeftCornerRadius, sq.CornerRadius)
	bottomRightRadius := GetCornerRadius(sq.BottomRightCornerRadius, sq.CornerRadius)
	bottomLeftRadius := GetCornerRadius(sq.BottomLeftCornerRadius, sq.CornerRadius)
	return drawOblong(sq.FillColor, sq.StrokeColor, sq.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rWidth, rHeight, vectorPad, scale)
}

func drawOblong(fill, strokeCol color.Color, strokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	// The maximum possible corner radius for a circular shape
	maxCornerRadius := GetMaximumRadius(fyne.NewSize(rWidth, rHeight))

	if topRightRadius == canvas.RadiusMaximum {
		topRightRadius = maxCornerRadius
	}

	if topLeftRadius == canvas.RadiusMaximum {
		topLeftRadius = maxCornerRadius
	}

	if bottomRightRadius == canvas.RadiusMaximum {
		bottomRightRadius = maxCornerRadius
	}

	if bottomLeftRadius == canvas.RadiusMaximum {
		bottomLeftRadius = maxCornerRadius
	}

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
		if topRightRadius == topLeftRadius && bottomRightRadius == bottomLeftRadius && topRightRadius == bottomRightRadius {
			// If all corners are the same, we can draw a simple rectangle
			radius := topRightRadius
			if radius == 0 {
				rasterx.AddRect(float64(p1x), float64(p1y), float64(p3x), float64(p3y), 0, filler)
			} else {
				r := float64(scale(radius))
				rasterx.AddRoundRect(float64(p1x), float64(p1y), float64(p3x), float64(p3y), r, r, 0, rasterx.RoundGap, filler)
			}
		} else {
			rTL, rTR, rBR, rBL := scale(topLeftRadius), scale(topRightRadius), scale(bottomRightRadius), scale(bottomLeftRadius)
			// Top-left corner
			c := quarterCircleControl * rTL
			if c != 0 {
				filler.Start(rasterx.ToFixedP(float64(p1x), float64(p1y+rTL)))
				filler.CubeBezier(rasterx.ToFixedP(float64(p1x), float64(p1y+c)), rasterx.ToFixedP(float64(p1x+c), float64(p1y)), rasterx.ToFixedP(float64(p1x+rTL), float64(p1y)))
			} else {
				filler.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
			}
			// Top edge to top-right
			c = quarterCircleControl * rTR
			filler.Line(rasterx.ToFixedP(float64(p2x-rTR), float64(p2y)))
			if c != 0 {
				filler.CubeBezier(rasterx.ToFixedP(float64(p2x-c), float64(p2y)), rasterx.ToFixedP(float64(p2x), float64(p2y+c)), rasterx.ToFixedP(float64(p2x), float64(p2y+rTR)))
			}
			// Right edge to bottom-right
			c = quarterCircleControl * rBR
			filler.Line(rasterx.ToFixedP(float64(p3x), float64(p3y-rBR)))
			if c != 0 {
				filler.CubeBezier(rasterx.ToFixedP(float64(p3x), float64(p3y-c)), rasterx.ToFixedP(float64(p3x-c), float64(p3y)), rasterx.ToFixedP(float64(p3x-rBR), float64(p3y)))
			}
			// Bottom edge to bottom-left
			c = quarterCircleControl * rBL
			filler.Line(rasterx.ToFixedP(float64(p4x+rBL), float64(p4y)))
			if c != 0 {
				filler.CubeBezier(rasterx.ToFixedP(float64(p4x+c), float64(p4y)), rasterx.ToFixedP(float64(p4x), float64(p4y-c)), rasterx.ToFixedP(float64(p4x), float64(p4y-rBL)))
			}
			// Left edge to top-left
			filler.Line(rasterx.ToFixedP(float64(p1x), float64(p1y+rTL)))
			filler.Stop(true)
		}
		filler.Draw()
	}

	if strokeCol != nil && strokeWidth > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(strokeCol)
		dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
		rTL, rTR, rBR, rBL := scale(topLeftRadius), scale(topRightRadius), scale(bottomRightRadius), scale(bottomLeftRadius)
		c := quarterCircleControl * rTL
		if c != 0 {
			dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y+rTL)))
			dasher.CubeBezier(rasterx.ToFixedP(float64(p1x), float64(p1y+c)), rasterx.ToFixedP(float64(p1x+c), float64(p1y)), rasterx.ToFixedP(float64(p1x+rTL), float64(p2y)))
		} else {
			dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
		}
		c = quarterCircleControl * rTR
		dasher.Line(rasterx.ToFixedP(float64(p2x-rTR), float64(p2y)))
		if c != 0 {
			dasher.CubeBezier(rasterx.ToFixedP(float64(p2x-c), float64(p2y)), rasterx.ToFixedP(float64(p2x), float64(p2y+c)), rasterx.ToFixedP(float64(p2x), float64(p2y+rTR)))
		}
		c = quarterCircleControl * rBR
		dasher.Line(rasterx.ToFixedP(float64(p3x), float64(p3y-rBR)))
		if c != 0 {
			dasher.CubeBezier(rasterx.ToFixedP(float64(p3x), float64(p3y-c)), rasterx.ToFixedP(float64(p3x-c), float64(p3y)), rasterx.ToFixedP(float64(p3x-rBR), float64(p3y)))
		}
		c = quarterCircleControl * rBL
		dasher.Line(rasterx.ToFixedP(float64(p4x+rBL), float64(p4y)))
		if c != 0 {
			dasher.CubeBezier(rasterx.ToFixedP(float64(p4x+c), float64(p4y)), rasterx.ToFixedP(float64(p4x), float64(p4y-c)), rasterx.ToFixedP(float64(p4x), float64(p4y-rBL)))
		}
		dasher.Stop(true)
		dasher.Draw()
	}

	return raw
}

// drawRoundArc constructs a rounded pie slice or annular sector
func drawRoundArc(adder rasterx.Adder, cx, cy, outer, inner, start, sweep, cr float64) {
	if sweep == 0 {
		return
	}

	cosSinPoint := func(cx, cy, r, ang float64) (x, y float64) {
		return cx + r*math.Cos(ang), cy - r*math.Sin(ang)
	}

	// addCircularArc appends a circular arc to the current path using cubic Bezier approximation.
	// 'adder' must already be positioned at the arc start point.
	// sweep is signed (positive = CCW, negative = CW).
	addCircularArc := func(adder rasterx.Adder, cx, cy, r, start, sweep float64) {
		if sweep == 0 || r == 0 {
			return
		}
		segCount := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2.0)))
		da := sweep / float64(segCount)

		for i := 0; i < segCount; i++ {
			a1 := start + float64(i)*da
			a2 := a1 + da

			x1, y1 := cosSinPoint(cx, cy, r, a1)
			x2, y2 := cosSinPoint(cx, cy, r, a2)

			k := 4.0 / 3.0 * math.Tan((a2-a1)/4.0)
			// tangent unit vectors on our param (x = cx+rcos, y = cy-rsin)
			c1x := x1 + k*r*(-math.Sin(a1))
			c1y := y1 + k*r*(-math.Cos(a1))
			c2x := x2 - k*r*(-math.Sin(a2))
			c2y := y2 - k*r*(-math.Cos(a2))

			adder.CubeBezier(
				rasterx.ToFixedP(c1x, c1y),
				rasterx.ToFixedP(c2x, c2y),
				rasterx.ToFixedP(x2, y2),
			)
		}
	}

	// full-circle/donut paths (two closed subpaths: outer CCW, inner CW if inner > 0)
	if math.Abs(sweep) >= 2*math.Pi {
		// outer loop (CCW)
		ox, oy := cosSinPoint(cx, cy, outer, 0)
		adder.Start(rasterx.ToFixedP(ox, oy))
		addCircularArc(adder, cx, cy, outer, 0, 2*math.Pi)
		adder.Stop(true)
		// inner loop reversed (CW) to create a hole
		if inner > 0 {
			ix, iy := cosSinPoint(cx, cy, inner, 0)
			adder.Start(rasterx.ToFixedP(ix, iy))
			addCircularArc(adder, cx, cy, inner, 0, -2*math.Pi)
			adder.Stop(true)
		}
		return
	}

	if cr <= 0 {
		// sharp-corner fallback
		if inner <= 0 {
			// pie slice
			ox, oy := cosSinPoint(cx, cy, outer, start)
			adder.Start(rasterx.ToFixedP(cx, cy))
			adder.Line(rasterx.ToFixedP(ox, oy))
			addCircularArc(adder, cx, cy, outer, start, sweep)
			adder.Line(rasterx.ToFixedP(cx, cy))
			adder.Stop(true)
			return
		}
		// annular sector
		outerStartX, outerStartY := cosSinPoint(cx, cy, outer, start)
		adder.Start(rasterx.ToFixedP(outerStartX, outerStartY))
		addCircularArc(adder, cx, cy, outer, start, sweep)
		innerEndX, innerEndY := cosSinPoint(cx, cy, inner, start+sweep)
		adder.Line(rasterx.ToFixedP(innerEndX, innerEndY))
		addCircularArc(adder, cx, cy, inner, start+sweep, -sweep)
		adder.Stop(true)
		return
	}

	// rounded corners
	sgn := 1.0
	if sweep < 0 {
		sgn = -1.0
	}
	absSweep := math.Abs(sweep)

	// clamp the corner radius if the value is too large
	cr = math.Min(cr, outer/2)

	// trim angles due to rounds
	sOut := math.Sqrt(math.Max(0, outer*(outer-2*cr)))
	thetaOut := math.Atan2(cr, sOut) // positive

	crIn := math.Min(cr, 0.5*math.Min(outer-inner, math.Abs(sweep)*inner))
	var sIn, thetaIn float64
	if inner > 0 {
		sIn = math.Sqrt(math.Max(0, inner*(inner+2*crIn)))
		thetaIn = math.Atan2(crIn, sIn)
	}

	// ensure the trim does not exceed half the sweep
	thetaOut = math.Min(thetaOut, absSweep/2.0-1e-6)
	if thetaOut < 0 {
		thetaOut = 0
	}
	if inner > 0 {
		thetaIn = math.Min(thetaIn, absSweep/2.0-1e-6)
		if thetaIn < 0 {
			thetaIn = 0
		}
	}

	// trimmed arc angles
	startOuter := start + sgn*thetaOut
	endOuter := start + sweep - sgn*thetaOut

	startInner := 0.0
	endInner := 0.0
	if inner > 0 {
		startInner = start + sgn*thetaIn
		endInner = start + sweep - sgn*thetaIn
	}

	// direction frames at start/end radial lines
	// start side
	vSx, vSy := math.Cos(start), -math.Sin(start)
	tSx, tSy := -math.Sin(start), -math.Cos(start)
	nSx, nSy := sgn*tSx, sgn*tSy // interior side normal at start

	// end side
	endRad := start + sweep
	vEx, vEy := math.Cos(endRad), -math.Sin(endRad)
	tEx, tEy := -math.Sin(endRad), -math.Cos(endRad)
	nEx, nEy := -sgn*tEx, -sgn*tEy // interior side normal at end

	// key points on arcs
	pOutStartX, pOutStartY := cosSinPoint(cx, cy, outer, startOuter)
	pOutEndX, pOutEndY := cosSinPoint(cx, cy, outer, endOuter)

	var pInStartX, pInStartY, pInEndX, pInEndY float64
	if inner > 0 {
		pInStartX, pInStartY = cosSinPoint(cx, cy, inner, startInner)
		pInEndX, pInEndY = cosSinPoint(cx, cy, inner, endInner)
	}

	angleAt := func(cx, cy, x, y float64) float64 {
		return math.Atan2(cy-y, x-cx)
	}

	// round geometry at start/end
	// outer rounds
	aOutSx, aOutSy := cx+sOut*vSx, cy+sOut*vSy                      // radial tangent (start)
	fOutSx, fOutSy := aOutSx+cr*nSx, aOutSy+cr*nSy                  // round center (start)
	aOutEx, aOutEy := cx+sOut*vEx, cy+sOut*vEy                      // radial tangent (end)
	fOutEx, fOutEy := aOutEx+cr*nEx, aOutEy+cr*nEy                  // round center (end)
	phiOutEndB := angleAt(fOutEx, fOutEy, pOutEndX, pOutEndY)       // outer end trimmed point
	phiOutEndA := angleAt(fOutEx, fOutEy, aOutEx, aOutEy)           // end radial tangent
	phiOutStartA := angleAt(fOutSx, fOutSy, aOutSx, aOutSy)         // start radial tangent
	phiOutStartB := angleAt(fOutSx, fOutSy, pOutStartX, pOutStartY) // outer start trimmed point

	// inner rounds
	var aInSx, aInSy, fInSx, fInSy, aInEx, aInEy, fInEx, fInEy float64
	var phiInEndA, phiInEndB, phiInStartA, phiInStartB float64
	if inner > 0 {
		aInSx, aInSy = cx+sIn*vSx, cy+sIn*vSy
		fInSx, fInSy = aInSx+crIn*nSx, aInSy+crIn*nSy
		aInEx, aInEy = cx+sIn*vEx, cy+sIn*vEy
		fInEx, fInEy = aInEx+crIn*nEx, aInEy+crIn*nEy

		phiInEndA = angleAt(fInEx, fInEy, aInEx, aInEy)           // end radial tangent
		phiInEndB = angleAt(fInEx, fInEy, pInEndX, pInEndY)       // inner end trimmed point
		phiInStartB = angleAt(fInSx, fInSy, pInStartX, pInStartY) // inner start trimmed point
		phiInStartA = angleAt(fInSx, fInSy, aInSx, aInSy)         // start radial tangent
	}

	angleDiff := func(delta float64) float64 {
		return math.Atan2(math.Sin(delta), math.Cos(delta))
	}

	adder.Start(rasterx.ToFixedP(pOutStartX, pOutStartY))                                   // start at trimmed outer start
	addCircularArc(adder, cx, cy, outer, startOuter, endOuter-startOuter)                   // outer arc (trimmed)
	addCircularArc(adder, fOutEx, fOutEy, cr, phiOutEndB, angleDiff(phiOutEndA-phiOutEndB)) // end side: outer round to radial

	if inner > 0 {
		adder.Line(rasterx.ToFixedP(aInEx, aInEy))                                                 // end side: radial line to inner
		addCircularArc(adder, fInEx, fInEy, crIn, phiInEndA, angleDiff(phiInEndB-phiInEndA))       // end side: inner round to inner arc
		addCircularArc(adder, cx, cy, inner, endInner, startInner-endInner)                        // inner arc (reverse, trimmed)
		addCircularArc(adder, fInSx, fInSy, crIn, phiInStartB, angleDiff(phiInStartA-phiInStartB)) // start side: inner round to radial
		adder.Line(rasterx.ToFixedP(aOutSx, aOutSy))                                               // start side: radial line to outer
	} else {
		// pie slice: close via center with radial lines
		adder.Line(rasterx.ToFixedP(cx, cy))         // to center from end side
		adder.Line(rasterx.ToFixedP(aOutSx, aOutSy)) // to start-side radial tangent
	}

	// start side: outer round from radial to outer start
	addCircularArc(adder, fOutSx, fOutSy, cr, phiOutStartA, angleDiff(phiOutStartB-phiOutStartA))
	adder.Stop(true)
}

// GetCornerRadius returns the effective corner radius for a rectangle or square corner.
// If the specific corner radius (perCornerRadius) is zero, it falls back to the baseCornerRadius.
// Otherwise, it uses the specific corner radius provided.
//
// This allows for per-corner customization while maintaining a default overall radius.
func GetCornerRadius(perCornerRadius, baseCornerRadius float32) float32 {
	if perCornerRadius == 0.0 {
		return baseCornerRadius
	}
	return perCornerRadius
}

// GetMaximumRadius returns the maximum possible radius that fits within the given size.
// It calculates half of the smaller dimension (width or height) of the provided fyne.Size.
// This is typically used for drawing circular corners in rectangles, circles or squares.
func GetMaximumRadius(size fyne.Size) float32 {
	return fyne.Min(size.Height, size.Width) / 2
}
