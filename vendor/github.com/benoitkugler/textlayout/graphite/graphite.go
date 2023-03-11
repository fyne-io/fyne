// Package graphite Graphite implements a "smart font" system developed
// specifically to handle the complexities of lesser-known languages of the world.
package graphite

import (
	"github.com/benoitkugler/textlayout/fonts"
	"github.com/benoitkugler/textlayout/fonts/truetype"
)

const debugMode = 0

type (
	GID = fonts.GID
	Tag = truetype.Tag
	gid = uint16
)

type Position struct {
	X, Y float32
}

func (p Position) add(other Position) Position {
	return Position{p.X + other.X, p.Y + other.Y}
}

// returns p - other
func (p Position) sub(other Position) Position {
	return Position{p.X - other.X, p.Y - other.Y}
}

func (p Position) scale(s float32) Position {
	return Position{p.X * s, p.Y * s}
}

type rect struct {
	bl, tr Position
}

func (r rect) width() float32  { return r.tr.X - r.bl.X }
func (r rect) height() float32 { return r.tr.Y - r.bl.Y }

func (r rect) scale(s float32) rect {
	return rect{r.bl.scale(s), r.tr.scale(s)}
}

func (r rect) addPosition(pos Position) rect {
	return rect{r.bl.add(pos), r.tr.add(pos)}
}

func (r rect) widen(other rect) rect {
	out := r
	if r.bl.X > other.bl.X {
		out.bl.X = other.bl.X
	}
	if r.bl.Y > other.bl.Y {
		out.bl.Y = other.bl.Y
	}
	if r.tr.X < other.tr.X {
		out.tr.X = other.tr.X
	}
	if r.tr.Y < other.tr.Y {
		out.tr.Y = other.tr.Y
	}
	return out
}

func min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

const iSQRT2 = 0.707106781

const (
	kgmetLsb = iota
	kgmetRsb
	kgmetBbTop
	kgmetBbBottom
	kgmetBbLeft
	kgmetBbRight
	kgmetBbHeight
	kgmetBbWidth
	kgmetAdvWidth
	kgmetAdvHeight
	kgmetAscent
	kgmetDescent
)

type glyphBoxes struct {
	subBboxes      []rect // same length
	slantSubBboxes []rect // same length
	slant          rect
	bitmap         uint16
}

type glyph struct {
	attrs   attributeSet
	boxes   glyphBoxes
	advance struct{ x, y int16 }
	bbox    rect
}

func (g glyph) getMetric(metric uint8) int32 {
	switch metric {
	case kgmetLsb:
		return int32(g.bbox.bl.X)
	case kgmetRsb:
		return int32(g.advance.x) - int32(g.bbox.tr.X)
	case kgmetBbTop:
		return int32(g.bbox.tr.Y)
	case kgmetBbBottom:
		return int32(g.bbox.bl.Y)
	case kgmetBbLeft:
		return int32(g.bbox.bl.X)
	case kgmetBbRight:
		return int32(g.bbox.tr.X)
	case kgmetBbHeight:
		return int32(g.bbox.tr.Y - g.bbox.bl.Y)
	case kgmetBbWidth:
		return int32(g.bbox.tr.X - g.bbox.bl.X)
	case kgmetAdvWidth:
		return int32(g.advance.x)
	case kgmetAdvHeight:
		return int32(g.advance.y)
	default:
		return 0
	}
}

// FontOptions allows to specify a scale to get position
// in user units rather than in font units.
type FontOptions struct {
	scale float32 // scales from design units to ppm
	// isHinted bool
}

// NewFontOptions builds options from the given points per em.
func NewFontOptions(ppem uint16, face *GraphiteFace) *FontOptions {
	return &FontOptions{scale: float32(ppem) / float32(face.Upem())}
}

var _ fonts.FaceMetrics = GraphiteFace{}

// GraphiteFace contains the specific OpenType tables
// used by the Graphite engine.
// It also wraps the common OpenType metrics record.
type GraphiteFace struct {
	fonts.FaceMetrics

	names truetype.TableName
	cmap  truetype.Cmap
	silf  tableSilf
	sill  tableSill
	feat  tableFeat

	// aggregation of metrics and attributes
	glyphs []glyph

	numAttributes uint16 //  number of glyph attributes per glyph

	ascent, descent int32
}

// LoadGraphite read graphite tables from the given OpenType font.
// It returns an error if the tables are invalid or if the font is
// not a graphite font.
func LoadGraphite(font *truetype.Font) (*GraphiteFace, error) {
	var (
		out GraphiteFace
		err error
	)

	out.FaceMetrics = font
	out.cmap, _ = font.Cmap()
	out.names = font.Names

	htmx, glyphs := font.Hmtx, font.Glyf
	tables := font.Graphite

	out.sill, err = parseTableSill(tables.Sill)
	if err != nil {
		return nil, err
	}

	out.feat, err = parseTableFeat(tables.Feat)
	if err != nil {
		return nil, err
	}

	locations, numAttributes, err := parseTableGloc(tables.Gloc, font.NumGlyphs)
	if err != nil {
		return nil, err
	}

	out.numAttributes = numAttributes
	attrs, err := parseTableGlat(tables.Glat, locations)
	if err != nil {
		return nil, err
	}

	out.silf, err = parseTableSilf(tables.Silf, numAttributes, uint16(len(out.feat)))
	if err != nil {
		return nil, err
	}

	out.preprocessGlyphsAttributes(glyphs, htmx, attrs)

	return &out, nil
}

// process the 'glyf', 'htmx' and 'glat' tables to extract relevant info.
func (f *GraphiteFace) preprocessGlyphsAttributes(glyphs truetype.TableGlyf, htmx truetype.TableHVmtx,
	attrs tableGlat,
) {
	// take into account pseudo glyphs (len(glyphs) <= len(attrs))
	L := len(glyphs)

	f.glyphs = make([]glyph, len(attrs))

	for gid, attr := range attrs {
		dst := &f.glyphs[gid]
		if gid < L {
			dst.advance.x = htmx[gid].Advance
			data := glyphs[gid]
			dst.bbox = rect{
				bl: Position{float32(data.Xmin), float32(data.Ymin)},
				tr: Position{float32(data.Xmax), float32(data.Ymax)},
			}
		}
		dst.attrs = attr.attributes
		if attr.octaboxMetrics != nil {
			dst.boxes = attr.octaboxMetrics.computeBoxes(dst.bbox)
		}
	}
}

// FeaturesForLang selects the features and values for the given language, or
// the default ones if the language is not found.
func (f *GraphiteFace) FeaturesForLang(lang Tag) FeaturesValue {
	return f.sill.getFeatures(lang, f.feat)
}

// getGlyph return nil for invalid gid
func (f *GraphiteFace) getGlyph(gid GID) *glyph {
	if int(gid) < len(f.glyphs) {
		return &f.glyphs[gid]
	}
	return nil
}

func (f *GraphiteFace) getGlyphAttr(gid GID, attr uint16) int16 {
	if glyph := f.getGlyph(gid); glyph != nil {
		return glyph.attrs.get(attr)
	}
	return 0
}

func (f *GraphiteFace) getGlyphMetric(gid GID, metric uint8) int32 {
	switch metric {
	case kgmetAscent:
		return f.ascent
	case kgmetDescent:
		return f.descent
	}
	if glyph := f.getGlyph(gid); glyph != nil {
		return glyph.getMetric(metric)
	}
	return 0
}

func (f *GraphiteFace) runGraphite(seg *Segment, silf *passes) {
	if seg.dir&3 == 3 && silf.indexBidiPass == 0xFF {
		seg.doMirror(silf.attrMirroring)
	}
	res := silf.runGraphite(seg, 0, silf.indexPosPass, true)
	if res {
		seg.associateChars(0, len(seg.charinfo))
		if silf.hasCollision {
			ok := seg.initCollisions()
			res = ok
		}
		if res {
			res = silf.runGraphite(seg, silf.indexPosPass, uint8(len(silf.passes)), false)
		}
	}

	if debugMode >= 2 {
		seg.positionSlots(nil, nil, nil, seg.currdir(), true)
		tr.finaliseOutput(seg)
	}
}

// Shape process the given `text` and applies the graphite tables
// found in the font, returning a shaped segment of text.
// `font` is optional: if given, the positions are scaled; otherwise they are
// expressed in font units.
// If `features` is nil, the default features from the `Sill` table are used.
// Note that this not the same as passing an empty slice, which would desactivate any feature.
// `script` is optional and may help to select the correct `Silf` subtable.
// `dir` sets the direction of the text.
func (face *GraphiteFace) Shape(font *FontOptions, text []rune, script Tag, features FeaturesValue, dir int8) *Segment {
	var seg Segment

	seg.face = face

	// allocate memory
	seg.charinfo = make([]charInfo, len(text))

	// choose silf: for now script is unused
	// script = spaceToZero(script) // adapt convention
	if len(face.silf) != 0 {
		seg.silf = &face.silf[0]
	} else {
		seg.silf = &passes{}
	}

	seg.dir = dir
	if seg.silf.hasCollision {
		seg.flags = 1 << 1
	}
	if seg.silf.attrSkipPasses != 0 {
		seg.passBits = ^uint32(0)
	}

	if features == nil {
		features = face.FeaturesForLang(0)
	}
	seg.feats = features

	seg.processRunes(text)

	face.runGraphite(&seg, seg.silf)

	seg.finalise(font, true)
	return &seg
}
