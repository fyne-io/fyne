// Copyright 2017 The oksvg Authors. All rights reserved.
//
// created: 2/12/2017 by S.R.Wiley
// The oksvg package provides a partial implementation of the SVG 2.0 standard.
// It can perform all SVG2.0 path commands, including arc and miterclip. It also
// has some additional capabilities like arc-clip. Svgdraw does
// not implement all SVG features such as animation or markers, but it can draw
// the many of open source SVG icons correctly. See Readme for
// a list of features.

package oksvg

import (
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"

	"encoding/xml"
	"errors"
	"image/color"
	"log"
	"math"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/colornames"
	"golang.org/x/image/math/fixed"
)

type (
	// PathStyle holds the state of the SVG style
	PathStyle struct {
		FillOpacity, LineOpacity          float64
		LineWidth, DashOffset, MiterLimit float64
		Dash                              []float64
		UseNonZeroWinding                 bool
		fillerColor, linerColor           interface{} // either color.Color or *rasterx.Gradient
		LineGap                           rasterx.GapFunc
		LeadLineCap                       rasterx.CapFunc // This is used if different than LineCap
		LineCap                           rasterx.CapFunc
		LineJoin                          rasterx.JoinMode
		mAdder                            rasterx.MatrixAdder // current transform
	}

	// SvgPath binds a style to a path
	SvgPath struct {
		PathStyle
		Path rasterx.Path
	}

	// SvgIcon holds data from parsed SVGs
	SvgIcon struct {
		ViewBox      struct{ X, Y, W, H float64 }
		Titles       []string // Title elements collect here
		Descriptions []string // Description elements collect here
		Ids          map[string]interface{}
		SVGPaths     []SvgPath
		Transform    rasterx.Matrix2D
	}

	// IconCursor is used while parsing SVG files
	IconCursor struct {
		PathCursor
		icon                                   *SvgIcon
		StyleStack                             []PathStyle
		grad                                   *rasterx.Gradient
		inTitleText, inDescText, inGrad, inDef bool
	}
)

// DefaultStyle sets the default PathStyle to fill black, winding rule,
// full opacity, no stroke, ButtCap line end and Bevel line connect.
var DefaultStyle = PathStyle{1.0, 1.0, 2.0, 0.0, 4.0, nil, true,
	color.NRGBA{0x00, 0x00, 0x00, 0xff}, nil,
	nil, nil, rasterx.ButtCap, rasterx.Bevel, rasterx.MatrixAdder{M: rasterx.Identity}}

// Draw the compiled SVG icon into the GraphicContext.
// All elements should be contained by the Bounds rectangle of the SvgIcon.
func (s *SvgIcon) Draw(r *rasterx.Dasher, opacity float64) {
	for _, svgp := range s.SVGPaths {
		svgp.DrawTransformed(r, opacity, s.Transform)
	}
}

// SetTarget sets the Transform matrix to draw within the bounds of the rectangle arguments
func (s *SvgIcon) SetTarget(x, y, w, h float64) {
	scaleW := w / s.ViewBox.W
	scaleH := h / s.ViewBox.H
	s.Transform = rasterx.Identity.Translate(x-s.ViewBox.X, y-s.ViewBox.Y).Scale(scaleW, scaleH)

}

// Draw the compiled SvgPath into the Dasher.
func (svgp *SvgPath) Draw(r *rasterx.Dasher, opacity float64) {
	svgp.DrawTransformed(r, opacity, rasterx.Identity)
}

// DrawTransformed draws the compiled SvgPath into the Dasher while applying transform t.
func (svgp *SvgPath) DrawTransformed(r *rasterx.Dasher, opacity float64, t rasterx.Matrix2D) {
	m := svgp.mAdder.M
	svgp.mAdder.M = t.Mult(m)
	defer func() { svgp.mAdder.M = m }() // Restore untransformed matrix
	if svgp.fillerColor != nil {
		r.Clear()
		rf := &r.Filler
		rf.SetWinding(svgp.UseNonZeroWinding)
		svgp.mAdder.Adder = rf // This allows transformations to be applied
		svgp.Path.AddTo(&svgp.mAdder)

		switch fillerColor := svgp.fillerColor.(type) {
		case color.Color:
			rf.SetColor(rasterx.ApplyOpacity(fillerColor, svgp.FillOpacity*opacity))
		case *rasterx.Gradient:
			if fillerColor.Units == rasterx.ObjectBoundingBox {
				fRect := rf.Scanner.GetPathExtent()
				mnx, mny := float64(fRect.Min.X)/64, float64(fRect.Min.Y)/64
				mxx, mxy := float64(fRect.Max.X)/64, float64(fRect.Max.Y)/64
				fillerColor.Bounds.X, fillerColor.Bounds.Y = mnx, mny
				fillerColor.Bounds.W, fillerColor.Bounds.H = mxx-mnx, mxy-mny
			}
			rf.SetColor(fillerColor.GetColorFunction(svgp.FillOpacity * opacity))
		}
		rf.Draw()
		// default is true
		rf.SetWinding(true)
	}
	if svgp.linerColor != nil {
		r.Clear()
		svgp.mAdder.Adder = r
		lineGap := svgp.LineGap
		if lineGap == nil {
			lineGap = DefaultStyle.LineGap
		}
		lineCap := svgp.LineCap
		if lineCap == nil {
			lineCap = DefaultStyle.LineCap
		}
		leadLineCap := lineCap
		if svgp.LeadLineCap != nil {
			leadLineCap = svgp.LeadLineCap
		}
		r.SetStroke(fixed.Int26_6(svgp.LineWidth*64),
			fixed.Int26_6(svgp.MiterLimit*64), leadLineCap, lineCap,
			lineGap, svgp.LineJoin, svgp.Dash, svgp.DashOffset)
		svgp.Path.AddTo(&svgp.mAdder)
		switch linerColor := svgp.linerColor.(type) {
		case color.Color:
			r.SetColor(rasterx.ApplyOpacity(linerColor, svgp.LineOpacity*opacity))
		case *rasterx.Gradient:
			if linerColor.Units == rasterx.ObjectBoundingBox {
				fRect := r.Scanner.GetPathExtent()
				mnx, mny := float64(fRect.Min.X)/64, float64(fRect.Min.Y)/64
				mxx, mxy := float64(fRect.Max.X)/64, float64(fRect.Max.Y)/64
				linerColor.Bounds.X, linerColor.Bounds.Y = mnx, mny
				linerColor.Bounds.W, linerColor.Bounds.H = mxx-mnx, mxy-mny
			}
			r.SetColor(linerColor.GetColorFunction(svgp.LineOpacity * opacity))
		}
		r.Draw()
	}
}

// ParseSVGColorNum reads the SFG color string e.g. #FBD9BD
func ParseSVGColorNum(colorStr string) (r, g, b uint8, err error) {
	colorStr = strings.TrimPrefix(colorStr, "#")
	var t uint64
	if len(colorStr) != 6 {
		// SVG specs say duplicate characters in case of 3 digit hex number
		colorStr = string([]byte{colorStr[0], colorStr[0],
			colorStr[1], colorStr[1], colorStr[2], colorStr[2]})
	}
	for _, v := range []struct {
		c *uint8
		s string
	}{
		{&r, colorStr[0:2]},
		{&g, colorStr[2:4]},
		{&b, colorStr[4:6]}} {
		t, err = strconv.ParseUint(v.s, 16, 8)
		if err != nil {
			return
		}
		*v.c = uint8(t)
	}
	return
}

// ParseSVGColor parses an SVG color string in all forms
// including all SVG1.1 names, obtained from the colornames package
func ParseSVGColor(colorStr string) (color.Color, error) {
	//_, _, _, a := curColor.RGBA()
	v := strings.ToLower(colorStr)
	if strings.HasPrefix(v, "url") { // We are not handling urls
		// and gradients and stuff at this point
		return color.NRGBA{0, 0, 0, 255}, nil
	}
	switch v {
	case "none":
		// nil signals that the function (fill or stroke) is off;
		// not the same as black
		return nil, nil
	default:
		cn, ok := colornames.Map[v]
		if ok {
			r, g, b, a := cn.RGBA()
			return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
		}
	}
	cStr := strings.TrimPrefix(colorStr, "rgb(")
	if cStr != colorStr {
		cStr := strings.TrimSuffix(cStr, ")")
		vals := strings.Split(cStr, ",")
		if len(vals) != 3 {
			return color.NRGBA{}, errParamMismatch
		}
		var cvals [3]uint8
		var err error
		for i := range cvals {
			cvals[i], err = parseColorValue(vals[i])
			if err != nil {
				return nil, err
			}
		}
		return color.NRGBA{cvals[0], cvals[1], cvals[2], 0xFF}, nil
	}
	if colorStr[0] == '#' {
		r, g, b, err := ParseSVGColorNum(colorStr)
		if err != nil {
			return nil, err
		}
		return color.NRGBA{r, g, b, 0xFF}, nil
	}
	return nil, errParamMismatch
}

func parseColorValue(v string) (uint8, error) {
	if v[len(v)-1] == '%' {
		n, err := strconv.Atoi(strings.TrimSpace(v[:len(v)-1]))
		if err != nil {
			return 0, err
		}
		return uint8(n * 0xFF / 100), nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if n > 255 {
		n = 255
	}
	return uint8(n), err
}

func (c *IconCursor) readTransformAttr(m1 rasterx.Matrix2D, k string) (rasterx.Matrix2D, error) {
	ln := len(c.points)
	switch k {
	case "rotate":
		if ln == 1 {
			m1 = m1.Rotate(c.points[0] * math.Pi / 180)
		} else if ln == 3 {
			m1 = m1.Translate(c.points[1], c.points[2]).
				Rotate(c.points[0]*math.Pi/180).
				Translate(-c.points[1], -c.points[2])
		} else {
			return m1, errParamMismatch
		}
	case "translate":
		if ln == 1 {
			m1 = m1.Translate(c.points[0], 0)
		} else if ln == 2 {
			m1 = m1.Translate(c.points[0], c.points[1])
		} else {
			return m1, errParamMismatch
		}
	case "skewx":
		if ln == 1 {
			m1 = m1.SkewX(c.points[0] * math.Pi / 180)
		} else {
			return m1, errParamMismatch
		}
	case "skewy":
		if ln == 1 {
			m1 = m1.SkewY(c.points[0] * math.Pi / 180)
		} else {
			return m1, errParamMismatch
		}
	case "scale":
		if ln == 1 {
			m1 = m1.Scale(c.points[0], 0)
		} else if ln == 2 {
			m1 = m1.Scale(c.points[0], c.points[1])
		} else {
			return m1, errParamMismatch
		}
	case "matrix":
		if ln == 6 {
			m1 = m1.Mult(rasterx.Matrix2D{
				A: c.points[0],
				B: c.points[1],
				C: c.points[2],
				D: c.points[3],
				E: c.points[4],
				F: c.points[5]})
		} else {
			return m1, errParamMismatch
		}
	default:
		return m1, errParamMismatch
	}
	return m1, nil
}

func (c *IconCursor) parseTransform(v string) (rasterx.Matrix2D, error) {
	ts := strings.Split(v, ")")
	m1 := c.StyleStack[len(c.StyleStack)-1].mAdder.M
	for _, t := range ts {
		t = strings.TrimSpace(t)
		if len(t) == 0 {
			continue
		}
		d := strings.Split(t, "(")
		if len(d) != 2 || len(d[1]) < 1 {
			return m1, errParamMismatch // badly formed transformation
		}
		err := c.GetPoints(d[1])
		if err != nil {
			return m1, err
		}
		m1, err = c.readTransformAttr(m1, strings.ToLower(strings.TrimSpace(d[0])))
		if err != nil {
			return m1, err
		}
	}
	return m1, nil
}

func (c *IconCursor) readStyleAttr(curStyle *PathStyle, k, v string) error {
	switch k {
	case "fill":
		gradient, err := c.ReadGradURL(v)
		if err != nil {
			return err
		}
		if gradient != nil {
			curStyle.fillerColor = gradient
			break
		}
		curStyle.fillerColor, err = ParseSVGColor(v)
		return err
	case "stroke":
		gradient, err := c.ReadGradURL(v)
		if gradient != nil {
			curStyle.linerColor = gradient
			break
		}
		if err != nil {
			return err
		}
		col, errc := ParseSVGColor(v)
		if errc != nil {
			return errc
		}
		if col != nil {
			curStyle.linerColor = col.(color.NRGBA)
		} else {
			curStyle.linerColor = nil
		}
	case "stroke-linegap":
		switch v {
		case "flat":
			curStyle.LineGap = rasterx.FlatGap
		case "round":
			curStyle.LineGap = rasterx.RoundGap
		case "cubic":
			curStyle.LineGap = rasterx.CubicGap
		case "quadratic":
			curStyle.LineGap = rasterx.QuadraticGap
		}
	case "stroke-leadlinecap":
		switch v {
		case "butt":
			curStyle.LeadLineCap = rasterx.ButtCap
		case "round":
			curStyle.LeadLineCap = rasterx.RoundCap
		case "square":
			curStyle.LeadLineCap = rasterx.SquareCap
		case "cubic":
			curStyle.LeadLineCap = rasterx.CubicCap
		case "quadratic":
			curStyle.LeadLineCap = rasterx.QuadraticCap
		}
	case "stroke-linecap":
		switch v {
		case "butt":
			curStyle.LineCap = rasterx.ButtCap
		case "round":
			curStyle.LineCap = rasterx.RoundCap
		case "square":
			curStyle.LineCap = rasterx.SquareCap
		case "cubic":
			curStyle.LineCap = rasterx.CubicCap
		case "quadratic":
			curStyle.LineCap = rasterx.QuadraticCap
		}
	case "stroke-linejoin":
		switch v {
		case "miter":
			curStyle.LineJoin = rasterx.Miter
		case "miter-clip":
			curStyle.LineJoin = rasterx.MiterClip
		case "arc-clip":
			curStyle.LineJoin = rasterx.ArcClip
		case "round":
			curStyle.LineJoin = rasterx.Round
		case "arc":
			curStyle.LineJoin = rasterx.Arc
		case "bevel":
			curStyle.LineJoin = rasterx.Bevel
		}
	case "stroke-miterlimit":
		mLimit, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		curStyle.MiterLimit = mLimit
	case "stroke-width":
		v = strings.TrimSuffix(v, "px")
		width, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		curStyle.LineWidth = width
	case "stroke-dashoffset":
		dashOffset, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		curStyle.DashOffset = dashOffset
	case "stroke-dasharray":
		if v != "none" {
			dashes := strings.Split(v, ",")
			dList := make([]float64, len(dashes))
			for i, dstr := range dashes {
				d, err := strconv.ParseFloat(strings.TrimSpace(dstr), 64)
				if err != nil {
					return err
				}
				dList[i] = d
			}
			curStyle.Dash = dList
			break
		}
	case "opacity", "stroke-opacity", "fill-opacity":
		op, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		if k != "stroke-opacity" {
			curStyle.FillOpacity *= op
		}
		if k != "fill-opacity" {
			curStyle.LineOpacity *= op
		}
	case "transform":
		m, err := c.parseTransform(v)
		if err != nil {
			return err
		}
		curStyle.mAdder.M = m
	}
	return nil
}

// PushStyle parses the style element, and push it on the style stack. Only color and opacity are supported
// for fill. Note that this parses both the contents of a style attribute plus
// direct fill and opacity attributes.
func (c *IconCursor) PushStyle(se xml.StartElement) error {
	var pairs []string
	for _, attr := range se.Attr {
		switch strings.ToLower(attr.Name.Local) {
		case "style":
			pairs = append(pairs, strings.Split(attr.Value, ";")...)
		default:
			pairs = append(pairs, attr.Name.Local+":"+attr.Value)
		}
	}
	// Make a copy of the top style
	curStyle := c.StyleStack[len(c.StyleStack)-1]
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) >= 2 {
			k := strings.ToLower(kv[0])
			k = strings.TrimSpace(k)
			v := strings.TrimSpace(kv[1])
			err := c.readStyleAttr(&curStyle, k, v)
			if err != nil {
				return err
			}
		}
	}
	c.StyleStack = append(c.StyleStack, curStyle) // Push style onto stack
	return nil
}

// unitSuffixes are suffixes sometimes applied to the width and height attributes
// of the svg element.
var unitSuffixes = [3]string{"cm", "mm", "px"}

func trimSuffixes(a string) (b string) {
	b = a
	for _, v := range unitSuffixes {
		b = strings.TrimSuffix(b, v)
	}
	return
}

func (c *IconCursor) readStartElement(se xml.StartElement) (err error) {
	icon := c.icon
	switch se.Name.Local {
	case "svg":
		icon.ViewBox.X = 0
		icon.ViewBox.Y = 0
		icon.ViewBox.W = 0
		icon.ViewBox.H = 0
		var width, height float64
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "viewBox":
				err = c.GetPoints(attr.Value)
				if len(c.points) != 4 {
					return errParamMismatch
				}
				icon.ViewBox.X = c.points[0]
				icon.ViewBox.Y = c.points[1]
				icon.ViewBox.W = c.points[2]
				icon.ViewBox.H = c.points[3]
			case "width":
				wn := trimSuffixes(attr.Value)
				width, err = strconv.ParseFloat(wn, 64)
			case "height":
				hn := trimSuffixes(attr.Value)
				height, err = strconv.ParseFloat(hn, 64)
			}
			if err != nil {
				return
			}
		}
		if icon.ViewBox.W == 0 {
			icon.ViewBox.W = width
		}
		if icon.ViewBox.H == 0 {
			icon.ViewBox.H = height
		}
	case "g": // G does nothing but push the style
	case "rect":
		var x, y, w, h, rx, ry float64
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "x":
				x, err = strconv.ParseFloat(attr.Value, 64)
			case "y":
				y, err = strconv.ParseFloat(attr.Value, 64)
			case "width":
				w, err = strconv.ParseFloat(attr.Value, 64)
			case "height":
				h, err = strconv.ParseFloat(attr.Value, 64)
			case "rx":
				rx, err = strconv.ParseFloat(attr.Value, 64)
			case "ry":
				ry, err = strconv.ParseFloat(attr.Value, 64)
			}
			if err != nil {
				return
			}
		}
		if w == 0 || h == 0 {
			break
		}
		rasterx.AddRoundRect(x, y, w+x, h+y, rx, ry, 0, rasterx.RoundGap, &c.Path)
	case "circle", "ellipse":
		var cx, cy, rx, ry float64
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "cx":
				cx, err = strconv.ParseFloat(attr.Value, 64)
			case "cy":
				cy, err = strconv.ParseFloat(attr.Value, 64)
			case "r":
				rx, err = strconv.ParseFloat(attr.Value, 64)
				ry = rx
			case "rx":
				rx, err = strconv.ParseFloat(attr.Value, 64)
			case "ry":
				ry, err = strconv.ParseFloat(attr.Value, 64)
			}
			if err != nil {
				return
			}
		}
		if rx == 0 || ry == 0 { // not drawn, but not an error
			break
		}
		c.EllipseAt(cx, cy, rx, ry)
	case "line":
		var x1, x2, y1, y2 float64
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "x1":
				x1, err = strconv.ParseFloat(attr.Value, 64)
			case "x2":
				x2, err = strconv.ParseFloat(attr.Value, 64)
			case "y1":
				y1, err = strconv.ParseFloat(attr.Value, 64)
			case "y2":
				y2, err = strconv.ParseFloat(attr.Value, 64)
			}
			if err != nil {
				return
			}
		}
		c.Path.Start(fixed.Point26_6{
			X: fixed.Int26_6(x1 * 64),
			Y: fixed.Int26_6(y1 * 64)})
		c.Path.Line(fixed.Point26_6{
			X: fixed.Int26_6(x2 * 64),
			Y: fixed.Int26_6(y2 * 64)})
	case "polygon", "polyline":
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "points":
				err = c.GetPoints(attr.Value)
				if len(c.points)%2 != 0 {
					return errors.New("polygon has odd number of points")
				}
			}
			if err != nil {
				return
			}
		}
		if len(c.points) > 4 {
			c.Path.Start(fixed.Point26_6{
				X: fixed.Int26_6(c.points[0] * 64),
				Y: fixed.Int26_6(c.points[1] * 64)})
			for i := 2; i < len(c.points)-1; i += 2 {
				c.Path.Line(fixed.Point26_6{
					X: fixed.Int26_6(c.points[i] * 64),
					Y: fixed.Int26_6(c.points[i+1] * 64)})
			}
			if se.Name.Local == "polygon" { // SVG spec sez polylines dont have close
				c.Path.Stop(true)
			}
		}
	case "path":
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "d":
				err = c.CompilePath(attr.Value)
			}
			if err != nil {
				return err
			}
		}
	case "desc":
		c.inDescText = true
		icon.Descriptions = append(icon.Descriptions, "")
	case "title":
		c.inTitleText = true
		icon.Titles = append(icon.Titles, "")
	case "def":
		c.inDef = true
	case "linearGradient":
		c.inGrad = true
		c.grad = &rasterx.Gradient{Points: [5]float64{0, 0, 1, 0, 0},
			IsRadial: false, Bounds: icon.ViewBox, Matrix: rasterx.Identity}
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "id":
				id := attr.Value
				if len(id) >= 0 {
					icon.Ids[id] = c.grad
				} else {
					return errZeroLengthID
				}
			case "x1":
				c.grad.Points[0], err = readFraction(attr.Value)
			case "y1":
				c.grad.Points[1], err = readFraction(attr.Value)
			case "x2":
				c.grad.Points[2], err = readFraction(attr.Value)
			case "y2":
				c.grad.Points[3], err = readFraction(attr.Value)
			default:
				err = c.ReadGradAttr(attr)
			}
			if err != nil {
				return err
			}
		}
	case "radialGradient":
		c.inGrad = true
		c.grad = &rasterx.Gradient{Points: [5]float64{0.5, 0.5, 0.5, 0.5, 0.5},
			IsRadial: true, Bounds: icon.ViewBox, Matrix: rasterx.Identity}
		var setFx, setFy bool
		for _, attr := range se.Attr {
			switch attr.Name.Local {
			case "id":
				id := attr.Value
				if len(id) >= 0 {
					icon.Ids[id] = c.grad
				} else {
					return errZeroLengthID
				}
			case "r":
				c.grad.Points[4], err = readFraction(attr.Value)
			case "cx":
				c.grad.Points[0], err = readFraction(attr.Value)
			case "cy":
				c.grad.Points[1], err = readFraction(attr.Value)
			case "fx":
				setFx = true
				c.grad.Points[2], err = readFraction(attr.Value)
			case "fy":
				setFy = true
				c.grad.Points[3], err = readFraction(attr.Value)
			default:
				err = c.ReadGradAttr(attr)
			}
			if err != nil {
				return err
			}
		}
		if setFx == false { // set fx to cx by default
			c.grad.Points[2] = c.grad.Points[0]
		}
		if setFy == false { // set fy to cy by default
			c.grad.Points[3] = c.grad.Points[1]
		}
	case "stop":
		if c.inGrad {
			stop := rasterx.GradStop{Opacity: 1.0}
			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "offset":
					stop.Offset, err = readFraction(attr.Value)
				case "stop-color":
					//todo: add current color inherit
					stop.StopColor, err = ParseSVGColor(attr.Value)
				case "stop-opacity":
					stop.Opacity, err = strconv.ParseFloat(attr.Value, 64)
				}
				if err != nil {
					return err
				}
			}
			c.grad.Stops = append(c.grad.Stops, stop)
		}

	default:
		errStr := "Cannot process svg element " + se.Name.Local
		if c.ErrorMode == StrictErrorMode {
			return errors.New(errStr)
		} else if c.ErrorMode == WarnErrorMode {
			log.Println(errStr)
		}
	}
	if len(c.Path) > 0 {
		//The cursor parsed a path from the xml element
		pathCopy := make(rasterx.Path, len(c.Path))
		copy(pathCopy, c.Path)
		icon.SVGPaths = append(icon.SVGPaths,
			SvgPath{c.StyleStack[len(c.StyleStack)-1], pathCopy})
		c.Path = c.Path[:0]
	}
	return
}

// ReadIconStream reads the Icon from the given io.Reader
// This only supports a sub-set of SVG, but
// is enough to draw many icons. If errMode is provided,
// the first value determines if the icon ignores, errors out, or logs a warning
// if it does not handle an element found in the icon file. Ignore warnings is
// the default if no ErrorMode value is provided.
func ReadIconStream(stream io.Reader, errMode ...ErrorMode) (*SvgIcon, error) {
	icon := &SvgIcon{Ids: make(map[string]interface{}), Transform: rasterx.Identity}
	cursor := &IconCursor{StyleStack: []PathStyle{DefaultStyle}, icon: icon}
	if len(errMode) > 0 {
		cursor.ErrorMode = errMode[0]
	}
	decoder := xml.NewDecoder(stream)
	decoder.CharsetReader = charset.NewReaderLabel
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return icon, err
		}
		// Inspect the type of the XML token
		switch se := t.(type) {
		case xml.StartElement:
			// Reads all recognized style attributes from the start element
			// and places it on top of the styleStack
			err = cursor.PushStyle(se)
			if err != nil {
				return icon, err
			}
			err = cursor.readStartElement(se)
			if err != nil {
				return icon, err
			}
		case xml.EndElement:
			// pop style
			cursor.StyleStack = cursor.StyleStack[:len(cursor.StyleStack)-1]
			switch se.Name.Local {
			case "title":
				cursor.inTitleText = false
			case "desc":
				cursor.inDescText = false
			case "def":
				cursor.inDef = false
			case "radialGradient", "linearGradient":
				cursor.inGrad = false
			}
		case xml.CharData:
			if cursor.inTitleText == true {
				icon.Titles[len(icon.Titles)-1] += string(se)
			}
			if cursor.inDescText == true {
				icon.Descriptions[len(icon.Descriptions)-1] += string(se)
			}
		}
	}
	return icon, nil
}

// ReadIcon reads the Icon from the named file
// This only supports a sub-set of SVG, but
// is enough to draw many icons. If errMode is provided,
// the first value determines if the icon ignores, errors out, or logs a warning
// if it does not handle an element found in the icon file. Ignore warnings is
// the default if no ErrorMode value is provided.
func ReadIcon(iconFile string, errMode ...ErrorMode) (*SvgIcon, error) {
	fin, errf := os.Open(iconFile)
	if errf != nil {
		return nil, errf
	}
	defer fin.Close()

	return ReadIconStream(fin, errMode...)
}

func readFraction(v string) (f float64, err error) {
	v = strings.TrimSpace(v)
	d := 1.0
	if strings.HasSuffix(v, "%") {
		d = 100
		v = strings.TrimSuffix(v, "%")
	}
	f, err = strconv.ParseFloat(v, 64)
	f /= d
	if f > 1 {
		f = 1
	} else if f < 0 {
		f = 0
	}
	return
}

// ReadGradURL reads an SVG format gradient url
func (c *IconCursor) ReadGradURL(v string) (grad *rasterx.Gradient, err error) {
	if strings.HasPrefix(v, "url(") && strings.HasSuffix(v, ")") {
		urlStr := strings.TrimSpace(v[4 : len(v)-1])
		if strings.HasPrefix(urlStr, "#") {
			switch grad := c.icon.Ids[urlStr[1:]].(type) {
			case *rasterx.Gradient:
				return grad, nil
			default:
				return nil, nil //missingIdError
			}

		}
	}
	return nil, nil // not a gradient url, and not an error
}

// ReadGradAttr reads an SVG gradient attribute
func (c *IconCursor) ReadGradAttr(attr xml.Attr) (err error) {
	switch attr.Name.Local {
	case "gradientTransform":
		c.grad.Matrix, err = c.parseTransform(attr.Value)
	case "gradientUnits":
		switch strings.TrimSpace(attr.Value) {
		case "userSpaceOnUse":
			c.grad.Units = rasterx.UserSpaceOnUse
		case "objectBoundingBox":
			c.grad.Units = rasterx.ObjectBoundingBox
		}
	case "spreadMethod":
		switch strings.TrimSpace(attr.Value) {
		case "pad":
			c.grad.Spread = rasterx.PadSpread
		case "reflect":
			c.grad.Spread = rasterx.ReflectSpread
		case "repeat":
			c.grad.Spread = rasterx.RepeatSpread
		}
	}
	return
}
