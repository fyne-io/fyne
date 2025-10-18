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
// 0°/360 is top, 90° is right, 180° is bottom, 270° is left
// 0°/-360 is top, -90° is left, -180° is bottom, -270° is right
func DrawArc(arc *canvas.Arc, vectorPad float32, scale func(float32) float32) *image.RGBA {
	size := arc.Size()

	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())

	centerX := float64(width) / 2
	centerY := float64(height) / 2

	outerRadius := fyne.Min(size.Width, size.Height) / 2
	innerRadius := float32(float64(outerRadius) * math.Min(1.0, math.Max(0.0, float64(arc.CutoutRatio))))
	cornerRadius := fyne.Min(GetMaximumRadiusArc(outerRadius, innerRadius, arc.EndAngle-arc.StartAngle), arc.CornerRadius)
	startAngle, endAngle := NormalizeArcAngles(arc.StartAngle, arc.EndAngle)

	// convert to radians
	startRad := float64(startAngle * math.Pi / 180.0)
	endRad := float64(endAngle * math.Pi / 180.0)
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

	cornerRadius = scale(cornerRadius)
	outerRadius = scale(outerRadius)
	innerRadius = scale(innerRadius)

	if arc.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(arc.FillColor)
		// rasterx.AddArc is not used because it does not support rounded corners
		drawRoundArc(filler, centerX, centerY, float64(outerRadius), float64(innerRadius), startRad, sweep, float64(cornerRadius))
		filler.Draw()
	}

	stroke := float64(scale(arc.StrokeWidth))
	if arc.StrokeColor != nil && stroke > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(arc.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(stroke*64), 0, nil, nil, nil, 0, nil, 0)
		// rasterx.AddArc is not used because it does not support rounded corners
		drawRoundArc(dasher, centerX, centerY, float64(outerRadius), float64(innerRadius), startRad, sweep, float64(cornerRadius))
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

// DrawPolygon rasterizes the given regular polygon object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawPolygon(polygon *canvas.Polygon, vectorPad float32, scale func(float32) float32) *image.RGBA {
	size := polygon.Size()

	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))
	outerRadius := scale(fyne.Min(size.Width, size.Height) / 2)
	cornerRadius := scale(fyne.Min(GetMaximumRadius(size), polygon.CornerRadius))
	sides := int(polygon.Sides)
	angle := polygon.Angle

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())

	if polygon.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(polygon.FillColor)
		drawRegularPolygon(float64(width/2), float64(height/2), float64(outerRadius), float64(cornerRadius), float64(angle), int(sides), filler)
		filler.Draw()
	}

	if polygon.StrokeColor != nil && polygon.StrokeWidth > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(polygon.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(float64(scale(polygon.StrokeWidth))*64), 0, nil, nil, nil, 0, nil, 0)
		drawRegularPolygon(float64(width/2), float64(height/2), float64(outerRadius), float64(cornerRadius), float64(angle), int(sides), dasher)
		dasher.Draw()
	}

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

func drawOblong(fill, strokeCol color.Color, strokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	// The maximum possible corner radii for a circular shape
	size := fyne.NewSize(rWidth, rHeight)
	topRightRadius = GetMaximumCornerRadius(topRightRadius, topLeftRadius, bottomRightRadius, size)
	topLeftRadius = GetMaximumCornerRadius(topLeftRadius, topRightRadius, bottomLeftRadius, size)
	bottomRightRadius = GetMaximumCornerRadius(bottomRightRadius, bottomLeftRadius, topRightRadius, size)
	bottomLeftRadius = GetMaximumCornerRadius(bottomLeftRadius, bottomRightRadius, topLeftRadius, size)

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

// drawRegularPolygon draws a regular n-sides centered at (cx,cy) with
// radius, rounded corners of cornerRadius, rotated by rot degrees.
func drawRegularPolygon(cx, cy, radius, cornerRadius, rot float64, sides int, p rasterx.Adder) {
	if sides < 3 || radius <= 0 {
		return
	}
	gf := rasterx.RoundGap
	angleStep := 2 * math.Pi / float64(sides)
	rotRads := rot*math.Pi/180 - math.Pi/2

	// fully rounded, draw circle
	if math.Min(cornerRadius, radius) == radius {
		rasterx.AddCircle(cx, cy, radius, p)
		return
	}

	// sharp polygon fast path
	if cornerRadius <= 0 {
		x0 := cx + radius*math.Cos(rotRads)
		y0 := cy + radius*math.Sin(rotRads)
		p.Start(rasterx.ToFixedP(x0, y0))
		for i := 1; i < sides; i++ {
			t := rotRads + angleStep*float64(i)
			p.Line(rasterx.ToFixedP(cx+radius*math.Cos(t), cy+radius*math.Sin(t)))
		}
		p.Stop(true)
		return
	}

	norm := func(x, y float64) (nx, ny float64) {
		l := math.Hypot(x, y)
		if l == 0 {
			return 0, 0
		}
		return x / l, y / l
	}

	// regular polygon vertices
	xs := make([]float64, sides)
	ys := make([]float64, sides)
	for i := 0; i < sides; i++ {
		t := rotRads + angleStep*float64(i)
		xs[i] = cx + radius*math.Cos(t)
		ys[i] = cy + radius*math.Sin(t)
	}

	// interior angle and side length
	alpha := math.Pi * (float64(sides) - 2) / float64(sides)
	r := cornerRadius

	// distances for tangency and center placement
	tTrim := r / math.Tan(alpha/2) // along each edge from vertex to tangency
	d := r / math.Sin(alpha/2)     // from vertex to arc center along interior bisector

	// precompute fillet geometry per vertex
	type pt struct{ x, y float64 }
	sPts := make([]pt, sides) // start tangency (on incoming edge)
	vS := make([]pt, sides)   // center->start vector
	vE := make([]pt, sides)   // center->end vector
	cPts := make([]pt, sides) // arc centers

	for i := 0; i < sides; i++ {
		prv := (i - 1 + sides) % sides
		nxt := (i + 1) % sides

		// unit directions
		uInX, uInY := xs[i]-xs[prv], ys[i]-ys[prv]   // prev -> i
		uOutX, uOutY := xs[nxt]-xs[i], ys[nxt]-ys[i] // i -> next
		uInX, uInY = norm(uInX, uInY)
		uOutX, uOutY = norm(uOutX, uOutY)

		// tangency points along edges from the vertex
		sx, sy := xs[i]-uInX*tTrim, ys[i]-uInY*tTrim   // incoming (toward prev)
		ex, ey := xs[i]+uOutX*tTrim, ys[i]+uOutY*tTrim // outgoing (toward next)

		// interior bisector direction and arc center
		bx, by := -uInX+uOutX, -uInY+uOutY
		bx, by = norm(bx, by)
		cxI, cyI := xs[i]+bx*d, ys[i]+by*d

		// center->tangent vectors
		vsx, vsy := sx-cxI, sy-cyI
		velx, vely := ex-cxI, ey-cyI

		sPts[i] = pt{sx, sy}
		vS[i] = pt{vsx, vsy}
		vE[i] = pt{velx, vely}
		cPts[i] = pt{cxI, cyI}
	}

	// start at s0, arc corner 0, then line+arc around, close last edge
	p.Start(rasterx.ToFixedP(sPts[0].x, sPts[0].y))
	gf(p,
		rasterx.ToFixedP(cPts[0].x, cPts[0].y),
		rasterx.ToFixedP(vS[0].x, vS[0].y),
		rasterx.ToFixedP(vE[0].x, vE[0].y),
	)
	for i := 1; i < sides; i++ {
		p.Line(rasterx.ToFixedP(sPts[i].x, sPts[i].y))
		gf(p,
			rasterx.ToFixedP(cPts[i].x, cPts[i].y),
			rasterx.ToFixedP(vS[i].x, vS[i].y),
			rasterx.ToFixedP(vE[i].x, vE[i].y),
		)
	}
	p.Line(rasterx.ToFixedP(sPts[0].x, sPts[0].y))
	p.Stop(true)
}

// drawRoundArc constructs a rounded pie slice or annular sector
// it uses the Unit circle coordinate system
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

// GetMaximumRadius returns the maximum possible corner radius that fits within the given size.
// It calculates half of the smaller dimension (width or height) of the provided fyne.Size.
//
// This is typically used for drawing circular corners in rectangles, circles or squares with the same radius for all corners.
func GetMaximumRadius(size fyne.Size) float32 {
	return fyne.Min(size.Height, size.Width) / 2
}

// GetMaximumCornerRadius returns the maximum possible corner radius for an individual corner,
// considering the specified corner radius, the radii of adjacent corners, and the maximum radii
// allowed for the width and height of the shape. Corner radius may utilize unused capacity from adjacent corners with radius smaller than maximum value
// so this corner can grow up to double the maximum radius of the smaller dimension (width or height) without causing overlaps.
//
// This is typically used for drawing circular corners in rectangles or squares with different corner radii.
func GetMaximumCornerRadius(radius, adjacentWidthRadius, adjacentHeightRadius float32, size fyne.Size) float32 {
	maxWidthRadius := size.Width / 2
	maxHeightRadius := size.Height / 2
	// fast path: corner radius fits within both per-axis maxima
	if radius <= fyne.Min(maxWidthRadius, maxHeightRadius) {
		return radius
	}
	// expand per-axis limits by borrowing any unused capacity from adjacent corners
	expandedMaxWidthRadius := 2*maxWidthRadius - fyne.Min(maxWidthRadius, adjacentWidthRadius)
	expandedMaxHeightRadius := 2*maxHeightRadius - fyne.Min(maxHeightRadius, adjacentHeightRadius)

	// respect the smaller axis and never exceed the requested radius
	expandedMaxRadius := fyne.Min(expandedMaxWidthRadius, expandedMaxHeightRadius)
	return fyne.Min(expandedMaxRadius, radius)
}

// GetMaximumRadiusArc returns the maximum possible corner radius for an arc segment based on the outer radius,
// inner radius, and sweep angle in degrees.
// It calculates half of the smaller dimension (thickness or effective length) of the provided arc parameters
func GetMaximumRadiusArc(outerRadius, innerRadius, sweepAngle float32) float32 {
	// height (thickness), width (length)
	thickness := outerRadius - innerRadius
	// TODO: length formula can be improved to get a fully rounded (pill shape) outer edge for thin (small sweep) arc segments
	span := math.Sin(0.5 * math.Min(math.Abs(float64(sweepAngle))*math.Pi/180.0, math.Pi)) // span in (0,1)
	length := 1.5 * float64(outerRadius) * span / (1 + span)                               // no division-by-zero risk

	return GetMaximumRadius(fyne.NewSize(
		thickness, float32(length),
	))
}

// NormalizeArcAngles adjusts the given start and end angles for arc drawing.
// It converts the angles from the Unit circle coordinate system (where 0 degrees is along the positive X-axis)
// to the coordinate system used by the painter, where 0 degrees is at the top (12 o'clock position).
// The function also reverses the direction: positive is clockwise, negative is counter-clockwise
func NormalizeArcAngles(startAngle, endAngle float32) (float32, float32) {
	return -(startAngle - 90), -(endAngle - 90)
}
