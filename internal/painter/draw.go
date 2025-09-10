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

// DrawPolygon rasterizes the given regular polygon object into an image.
// The bounds of the output image will be increased by vectorPad to allow for stroke overflow at the edges.
// The scale function is used to understand how many pixels are required per unit of size.
func DrawPolygon(polygon *canvas.Polygon, vectorPad float32, scale func(float32) float32) *image.RGBA {
	size := polygon.Size()
	cornerRadius := scale(polygon.CornerRadius)

	width := int(scale(size.Width + vectorPad*2))
	height := int(scale(size.Height + vectorPad*2))
	shapeRadius := fyne.Min(size.Width, size.Height) / 2

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(int(size.Width), int(size.Height), raw, raw.Bounds())

	if polygon.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(polygon.FillColor)
		drawRegularPolygon(float64(width/2), float64(height/2), float64(shapeRadius), float64(cornerRadius), float64(polygon.Rotation), int(polygon.Sides), filler)
		filler.Draw()
	}

	if polygon.StrokeColor != nil && polygon.StrokeWidth > 0 {
		dasher := rasterx.NewDasher(width, height, scanner)
		dasher.SetColor(polygon.StrokeColor)
		dasher.SetStroke(fixed.Int26_6(float64(polygon.StrokeWidth)*64), 0, nil, nil, nil, 0, nil, 0)
		drawRegularPolygon(float64(width/2), float64(height/2), float64(shapeRadius), float64(cornerRadius), float64(polygon.Rotation), int(polygon.Sides), dasher)
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

func DrawSquare(sq *canvas.Square, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
	topRightRadius := GetCornerRadius(sq.TopRightCornerRadius, sq.CornerRadius)
	topLeftRadius := GetCornerRadius(sq.TopLeftCornerRadius, sq.CornerRadius)
	bottomRightRadius := GetCornerRadius(sq.BottomRightCornerRadius, sq.CornerRadius)
	bottomLeftRadius := GetCornerRadius(sq.BottomLeftCornerRadius, sq.CornerRadius)
	return drawOblong(sq.FillColor, sq.StrokeColor, sq.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rWidth, rHeight, vectorPad, scale)
}

func drawOblong(fill, strokeCol color.Color, strokeWidth float32, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius float32, rWidth, rHeight, vectorPad float32, scale func(float32) float32) *image.RGBA {
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

	// sharp polygon fast path
	if cornerRadius <= 0 {
		angleStep := 2 * math.Pi / float64(sides)
		rotRads := -rot*math.Pi/180 - math.Pi/2
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
	angleStep := 2 * math.Pi / float64(sides)
	rotRads := -rot*math.Pi/180 - math.Pi/2
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
