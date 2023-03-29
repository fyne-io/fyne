package truetype

import (
	"math"

	"github.com/benoitkugler/textlayout/fonts"
)

var _ fonts.FaceMetrics = (*Font)(nil)

// Returns true if the font has Graphite capabilities,
// but does not check if the tables are actually valid.
func (font *Font) IsGraphite() (*Font, bool) {
	return font, font.Graphite != nil
}

func (f *Font) GetGlyphContourPoint(glyph fonts.GID, pointIndex uint16) (x, y int32, ok bool) {
	// harfbuzz seems not to implement this feature
	return 0, 0, false
}

func (f *Font) GlyphName(glyph GID) string {
	if postNames := f.post.Names; postNames != nil {
		if name := postNames.GlyphName(glyph); name != "" {
			return name
		}
	}
	if f.cff != nil {
		return f.cff.GlyphName(glyph)
	}
	return ""
}

func (f *Font) Upem() uint16 { return f.upem }

var (
	metricsTagHorizontalAscender  = MustNewTag("hasc")
	metricsTagHorizontalDescender = MustNewTag("hdsc")
	metricsTagHorizontalLineGap   = MustNewTag("hlgp")
	metricsTagVerticalAscender    = MustNewTag("vasc")
	metricsTagVerticalDescender   = MustNewTag("vdsc")
	metricsTagVerticalLineGap     = MustNewTag("vlgp")
)

func fixAscenderDescender(value float32, metricsTag Tag) float32 {
	if metricsTag == metricsTagHorizontalAscender || metricsTag == metricsTagVerticalAscender {
		return float32(math.Abs(float64(value)))
	}
	if metricsTag == metricsTagHorizontalDescender || metricsTag == metricsTagVerticalDescender {
		return float32(-math.Abs(float64(value)))
	}
	return value
}

func (f *Font) getPositionCommon(metricTag Tag) (float32, bool) {
	deltaVar := f.mvar.getVar(metricTag, f.varCoords)
	switch metricTag {
	case metricsTagHorizontalAscender:
		if f.OS2.useTypoMetrics() && f.OS2.hasData() {
			return fixAscenderDescender(float32(f.OS2.STypoAscender)+deltaVar, metricTag), true
		} else if f.hhea != nil {
			return fixAscenderDescender(float32(f.hhea.Ascent)+deltaVar, metricTag), true
		}

	case metricsTagHorizontalDescender:
		if f.OS2.useTypoMetrics() && f.OS2.hasData() {
			return fixAscenderDescender(float32(f.OS2.STypoDescender)+deltaVar, metricTag), true
		} else if f.hhea != nil {
			return fixAscenderDescender(float32(f.hhea.Descent)+deltaVar, metricTag), true
		}
	case metricsTagHorizontalLineGap:
		if f.OS2.useTypoMetrics() && f.OS2.hasData() {
			return fixAscenderDescender(float32(f.OS2.STypoLineGap)+deltaVar, metricTag), true
		} else if f.hhea != nil {
			return fixAscenderDescender(float32(f.hhea.LineGap)+deltaVar, metricTag), true
		}
	case metricsTagVerticalAscender:
		if f.vhea != nil {
			return fixAscenderDescender(float32(f.vhea.Ascent)+deltaVar, metricTag), true
		}
	case metricsTagVerticalDescender:
		if f.vhea != nil {
			return fixAscenderDescender(float32(f.vhea.Descent)+deltaVar, metricTag), true
		}
	case metricsTagVerticalLineGap:
		if f.vhea != nil {
			return fixAscenderDescender(float32(f.vhea.LineGap)+deltaVar, metricTag), true
		}
	}
	return 0, false
}

func (f *Font) FontHExtents() (fonts.FontExtents, bool) {
	var (
		out           fonts.FontExtents
		ok1, ok2, ok3 bool
	)
	out.Ascender, ok1 = f.getPositionCommon(metricsTagHorizontalAscender)
	out.Descender, ok2 = f.getPositionCommon(metricsTagHorizontalDescender)
	out.LineGap, ok3 = f.getPositionCommon(metricsTagHorizontalLineGap)
	return out, ok1 && ok2 && ok3
}

func (f *Font) FontVExtents() (fonts.FontExtents, bool) {
	var (
		out           fonts.FontExtents
		ok1, ok2, ok3 bool
	)
	out.Ascender, ok1 = f.getPositionCommon(metricsTagVerticalAscender)
	out.Descender, ok2 = f.getPositionCommon(metricsTagVerticalDescender)
	out.LineGap, ok3 = f.getPositionCommon(metricsTagVerticalLineGap)
	return out, ok1 && ok2 && ok3
}

var (
	tagStrikeoutSize      = MustNewTag("strs")
	tagStrikeoutOffset    = MustNewTag("stro")
	tagUnderlineSize      = MustNewTag("unds")
	tagUnderlineOffset    = MustNewTag("undo")
	tagSuperscriptYSize   = MustNewTag("spys")
	tagSuperscriptXOffset = MustNewTag("spxo")
	tagSubscriptYSize     = MustNewTag("sbys")
	tagSubscriptYOffset   = MustNewTag("sbyo")
	tagSubscriptXOffset   = MustNewTag("sbxo")
	tagXHeight            = MustNewTag("xhgt")
	tagCapHeight          = MustNewTag("cpht")
)

func (f *Font) LineMetric(metric fonts.LineMetric) (float32, bool) {
	switch metric {
	case fonts.UnderlinePosition:
		return float32(f.post.UnderlinePosition) + f.mvar.getVar(tagUnderlineOffset, f.varCoords), true
	case fonts.UnderlineThickness:
		return float32(f.post.UnderlineThickness) + f.mvar.getVar(tagUnderlineSize, f.varCoords), true
	case fonts.StrikethroughPosition:
		return float32(f.OS2.YStrikeoutPosition) + f.mvar.getVar(tagStrikeoutOffset, f.varCoords), true
	case fonts.StrikethroughThickness:
		return float32(f.OS2.YStrikeoutSize) + f.mvar.getVar(tagStrikeoutSize, f.varCoords), true
	case fonts.SuperscriptEmYSize:
		return float32(f.OS2.YSuperscriptYSize) + f.mvar.getVar(tagSuperscriptYSize, f.varCoords), true
	case fonts.SuperscriptEmXOffset:
		return float32(f.OS2.YSuperscriptXOffset) + f.mvar.getVar(tagSuperscriptXOffset, f.varCoords), true
	case fonts.SubscriptEmYSize:
		return float32(f.OS2.YSubscriptYSize) + f.mvar.getVar(tagSubscriptYSize, f.varCoords), true
	case fonts.SubscriptEmYOffset:
		return float32(f.OS2.YSubscriptYOffset) + f.mvar.getVar(tagSubscriptYOffset, f.varCoords), true
	case fonts.SubscriptEmXOffset:
		return float32(f.OS2.YSubscriptXOffset) + f.mvar.getVar(tagSubscriptXOffset, f.varCoords), true
	case fonts.XHeight:
		return float32(f.OS2.SxHeigh) + f.mvar.getVar(tagXHeight, f.varCoords), true
	case fonts.CapHeight:
		return float32(f.OS2.SCapHeight) + f.mvar.getVar(tagCapHeight, f.varCoords), true
	}
	return 0, false
}

func (f *Font) NominalGlyph(ch rune) (GID, bool) {
	return f.cmap.Lookup(ch)
}

func (f *Font) VariationGlyph(ch, varSelector rune) (GID, bool) {
	gid, kind := f.cmapVar.getGlyphVariant(ch, varSelector)
	switch kind {
	case variantNotFound:
		return 0, false
	case variantFound:
		return gid, true
	default: // variantUseDefault
		return f.NominalGlyph(ch)
	}
}

// do not take into account variations
func (f *Font) getBaseAdvance(gid GID, table TableHVmtx) int16 {
	if int(gid) >= len(table) {
		/* If `table` is empty, it means we don't have the metrics table
		 * for this direction: return default advance.  Otherwise, it means that the
		 * glyph index is out of bound: return zero. */
		if len(table) == 0 {
			return int16(f.upem)
		}
		return 0
	}
	return table[gid].Advance
}

const (
	phantomLeft = iota
	phantomRight
	phantomTop
	phantomBottom
	phantomCount
)

// use the `glyf` table to fetch the contour points,
// applying variation if needed.
// for composite, recursively calls itself; allPoints includes phantom points and will be at least of length 4
func (f *Font) getPointsForGlyph(gid GID, currentDepth int, allPoints *[]contourPoint /* OUT */) {
	// adapted from harfbuzz/src/hb-ot-glyf-table.hh

	if currentDepth > maxCompositeNesting || int(gid) >= len(f.Glyf) {
		return
	}
	g := f.Glyf[gid]

	var points []contourPoint
	if data, ok := g.data.(simpleGlyphData); ok {
		points = data.getContourPoints() // fetch the "real" points
	} else { // zeros values are enough
		points = make([]contourPoint, g.pointNumbersCount())
	}

	// init phantom point
	points = append(points, make([]contourPoint, phantomCount)...)
	phantoms := points[len(points)-phantomCount:]

	hDelta := float32(g.Xmin - f.Hmtx.getSideBearing(gid))
	vOrig := float32(g.Ymax + f.vmtx.getSideBearing(gid))
	hAdv := float32(f.getBaseAdvance(gid, f.Hmtx))
	vAdv := float32(f.getBaseAdvance(gid, f.vmtx))
	phantoms[phantomLeft].X = hDelta
	phantoms[phantomRight].X = hAdv + hDelta
	phantoms[phantomTop].Y = vOrig
	phantoms[phantomBottom].Y = vOrig - vAdv

	if f.isVar() {
		f.gvar.applyDeltasToPoints(gid, f.varCoords, points)
	}

	switch data := g.data.(type) {
	case simpleGlyphData:
		*allPoints = append(*allPoints, points...)
	case compositeGlyphData:
		for compIndex, item := range data.glyphs {
			// recurse on component
			var compPoints []contourPoint

			f.getPointsForGlyph(item.glyphIndex, currentDepth+1, &compPoints)

			LC := len(compPoints)
			if LC < phantomCount { // in case of max depth reached
				return
			}

			/* Copy phantom points from component if USE_MY_METRICS flag set */
			if item.hasUseMyMetrics() {
				copy(phantoms, compPoints[LC-phantomCount:])
			}

			/* Apply component transformation & translation */
			item.transformPoints(compPoints)

			/* Apply translation from gvar */
			tx, ty := points[compIndex].X, points[compIndex].Y
			for i := range compPoints {
				compPoints[i].translate(tx, ty)
			}

			if item.isAnchored() {
				p1, p2 := item.argsAsIndices()
				if p1 < len(*allPoints) && p2 < LC {
					tx, ty := (*allPoints)[p1].X-compPoints[p2].X, (*allPoints)[p1].Y-compPoints[p2].Y
					for i := range compPoints {
						compPoints[i].translate(tx, ty)
					}
				}
			}

			*allPoints = append(*allPoints, compPoints[0:LC-phantomCount]...)
		}

		*allPoints = append(*allPoints, phantoms...)
	default: // no data for the glyph
		*allPoints = append(*allPoints, phantoms...)
	}

	// apply at top level
	if currentDepth == 0 {
		/* Undocumented rasterizer behavior:
		 * Shift points horizontally by the updated left side bearing */
		tx := -phantoms[phantomLeft].X
		for i := range *allPoints {
			(*allPoints)[i].translate(tx, 0)
		}
	}
}

func extentsFromPoints(allPoints []contourPoint) (ext fonts.GlyphExtents) {
	truePoints := allPoints[:len(allPoints)-phantomCount]
	if len(truePoints) == 0 {
		// zero extent for the empty glyph
		return ext
	}
	minX, minY := truePoints[0].X, truePoints[0].Y
	maxX, maxY := minX, minY
	for _, p := range truePoints {
		minX = minF(minX, p.X)
		minY = minF(minY, p.Y)
		maxX = maxF(maxX, p.X)
		maxY = maxF(maxY, p.Y)
	}
	ext.XBearing = minX
	ext.YBearing = maxY
	ext.Width = maxX - minX
	ext.Height = minY - maxY
	return ext
}

// walk through the contour points of the given glyph to compute its extends and its phantom points
// As an optimization, if `computeExtents` is false, the extents computation is skipped (a zero value is returned).
func (f *Font) getGlyfPoints(gid GID, computeExtents bool) (ext fonts.GlyphExtents, ph [phantomCount]contourPoint) {
	if int(gid) >= len(f.Glyf) {
		return
	}
	var allPoints []contourPoint
	f.getPointsForGlyph(gid, 0, &allPoints)

	copy(ph[:], allPoints[len(allPoints)-phantomCount:])

	if computeExtents {
		ext = extentsFromPoints(allPoints)
	}

	return ext, ph
}

func clamp(v float32) float32 {
	if v < 0 {
		v = 0
	}
	return v
}

func ceil(v float32) int16 {
	return int16(math.Ceil(float64(v)))
}

func (f *Font) getGlyphAdvanceVar(gid GID, isVertical bool) float32 {
	_, phantoms := f.getGlyfPoints(gid, false)
	if isVertical {
		return clamp(phantoms[phantomTop].Y - phantoms[phantomBottom].Y)
	}
	return clamp(phantoms[phantomRight].X - phantoms[phantomLeft].X)
}

func (f *Font) HorizontalAdvance(gid GID) float32 {
	advance := f.getBaseAdvance(gid, f.Hmtx)
	if !f.isVar() {
		return float32(advance)
	}
	if f.hvar != nil {
		return float32(advance) + f.hvar.getAdvanceVar(gid, f.varCoords)
	}
	return f.getGlyphAdvanceVar(gid, false)
}

// return `true` is the font is variable and `varCoords` is valid
func (f *Font) isVar() bool {
	return len(f.varCoords) != 0 && len(f.varCoords) == len(f.fvar.Axis)
}

func (f *Font) VerticalAdvance(gid GID) float32 {
	// return the opposite of the advance from the font
	advance := f.getBaseAdvance(gid, f.vmtx)
	if !f.isVar() {
		return -float32(advance)
	}
	if f.vvar != nil {
		return -float32(advance) - f.vvar.getAdvanceVar(gid, f.varCoords)
	}
	return -f.getGlyphAdvanceVar(gid, true)
}

func (f *Font) getGlyphSideBearingVar(gid GID, isVertical bool) int16 {
	extents, phantoms := f.getGlyfPoints(gid, true)
	if isVertical {
		return ceil(phantoms[phantomTop].Y - extents.YBearing)
	}
	return int16(phantoms[phantomLeft].X)
}

// take variations into account
func (f *Font) getVerticalSideBearing(glyph GID) int16 {
	// base side bearing
	sideBearing := f.vmtx.getSideBearing(glyph)
	if !f.isVar() {
		return sideBearing
	}
	if f.vvar != nil {
		return sideBearing + int16(f.vvar.getSideBearingVar(glyph, f.varCoords))
	}
	return f.getGlyphSideBearingVar(glyph, true)
}

func (f *Font) GlyphHOrigin(GID) (x, y int32, found bool) {
	// zero is the right value here
	return 0, 0, true
}

func (f *Font) GlyphVOrigin(glyph GID) (x, y int32, found bool) {
	x = int32(f.HorizontalAdvance(glyph) / 2)

	if f.vorg != nil {
		y = int32(f.vorg.getYOrigin(glyph))
		return x, y, true
	}

	if extents, ok := f.getExtentsFromGlyf(glyph); ok {
		tsb := f.getVerticalSideBearing(glyph)
		y = int32(extents.YBearing) + int32(tsb)
		return x, y, true
	}

	fontExtents, ok := f.FontHExtents()
	y = int32(fontExtents.Ascender)

	return x, y, ok
}

func (f *Font) getExtentsFromGlyf(glyph GID) (fonts.GlyphExtents, bool) {
	if int(glyph) >= len(f.Glyf) {
		return fonts.GlyphExtents{}, false
	}
	g := f.Glyf[glyph]
	if f.isVar() { // we have to compute the outline points and apply variations
		extents, _ := f.getGlyfPoints(glyph, true)
		return extents, true
	}
	return g.getExtents(f.Hmtx, glyph), true
}

func (f *Font) getExtentsFromCBDT(glyph GID, xPpem, yPpem uint16) (fonts.GlyphExtents, bool) {
	strike := f.bitmap.chooseStrike(xPpem, yPpem)
	if strike == nil || strike.ppemX == 0 || strike.ppemY == 0 {
		return fonts.GlyphExtents{}, false
	}
	subtable := strike.findTable(glyph)
	if subtable == nil {
		return fonts.GlyphExtents{}, false
	}
	image := subtable.getImage(glyph)
	if image == nil {
		return fonts.GlyphExtents{}, false
	}
	extents := image.metrics.glyphExtents()

	/* convert to font units. */
	xScale := float32(f.upem) / float32(strike.ppemX)
	yScale := float32(f.upem) / float32(strike.ppemY)
	extents.XBearing *= xScale
	extents.YBearing *= yScale
	extents.Width *= xScale
	extents.Height *= yScale
	return extents, true
}

func (f *Font) getExtentsFromSbix(glyph GID, xPpem, yPpem uint16) (fonts.GlyphExtents, bool) {
	strike := f.sbix.chooseStrike(xPpem, yPpem)
	if strike == nil || strike.ppem == 0 {
		return fonts.GlyphExtents{}, false
	}
	data := strike.getGlyph(glyph, 0)
	if data.isNil() {
		return fonts.GlyphExtents{}, false
	}
	extents, ok := data.glyphExtents()

	/* convert to font units. */
	scale := float32(f.upem) / float32(strike.ppem)
	extents.XBearing *= scale
	extents.YBearing *= scale
	extents.Width *= scale
	extents.Height *= scale
	return extents, ok
}

func (f *Font) getExtentsFromCff1(glyph GID) (fonts.GlyphExtents, bool) {
	if f.cff == nil {
		return fonts.GlyphExtents{}, false
	}
	_, bounds, err := f.cff.LoadGlyph(glyph)
	if err != nil {
		return fonts.GlyphExtents{}, false
	}
	return bounds.ToExtents(), true
}

// func (f *fontMetrics) getExtentsFromCff2(glyph , coords []float32) (fonts.GlyphExtents, bool) {
// }

func (f *Font) GlyphExtents(glyph GID, xPpem, yPpem uint16) (fonts.GlyphExtents, bool) {
	out, ok := f.getExtentsFromSbix(glyph, xPpem, yPpem)
	if ok {
		return out, ok
	}
	out, ok = f.getExtentsFromGlyf(glyph)
	if ok {
		return out, ok
	}
	out, ok = f.getExtentsFromCff1(glyph)
	if ok {
		return out, ok
	}
	out, ok = f.getExtentsFromCBDT(glyph, xPpem, yPpem)
	return out, ok
}
