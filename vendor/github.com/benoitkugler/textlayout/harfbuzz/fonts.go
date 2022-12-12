package harfbuzz

import (
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
	"github.com/benoitkugler/textlayout/graphite"
)

// ported from src/hb-font.hh, src/hb-font.cc  Copyright © 2009  Red Hat, Inc., 2012  Google, Inc.  Behdad Esfahbod

// Face is the interface providing font metrics and layout information.
// Harfbuzz is mostly useful when used with fonts providing advanced layout capabilities :
// see the extension interface `FaceOpentype`.
type Face = fonts.Face

var _ FaceOpentype = (*tt.Font)(nil)

// FaceOpentype adds support for advanced layout features
// found in Opentype/Truetype font files.
// See the package fonts/truetype for more details.
type FaceOpentype interface {
	Face
	tt.FaceVariable

	// Returns true if the font has Graphite capabilities.
	// Note that tables validity will still be checked in `NewFont`,
	// using the table from the returned `tt.Font`.
	// Overide this method to disable Graphite functionalities.
	IsGraphite() (*tt.Font, bool)

	// LayoutTables fetchs the Opentype layout tables of the font.
	LayoutTables() tt.LayoutTables

	// GetGlyphContourPoint retrieves the (X,Y) coordinates (in font units) for a
	// specified contour point in a glyph, or false if not found.
	GetGlyphContourPoint(glyph fonts.GID, pointIndex uint16) (x, y int32, ok bool)

	// VariationGlyph retrieves the glyph ID for a specified Unicode code point
	// followed by a specified Variation Selector code point, or false if not found
	VariationGlyph(ch, varSelector rune) (fonts.GID, bool)
}

// Font is used internally as a light wrapper around the provided Face.
//
// While a font face is generally the in-memory representation of a static font file,
// `Font` handles dynamic attributes like size, width and
// other parameters (pixels-per-em, points-per-em, variation
// settings).
//
// Font are constructed with `NewFont` and adjusted by accessing the fields
// XPpem, YPpem, Ptem,XScale, YScale and with the method `SetVarCoordsDesign` for
// variable fonts.
type Font struct {
	face Face

	// only non nil for valid graphite fonts
	gr *graphite.GraphiteFace

	// opentype fields, initialized from a FaceOpentype
	otTables               *tt.LayoutTables
	gsubAccels, gposAccels []otLayoutLookupAccelerator // accelators for lookup
	faceUpem               int32                       // cached value of Face.Upem()

	// Point size of the font. Set to zero to unset.
	// This is used in AAT layout, when applying 'trak' table.
	Ptem float32

	// Horizontal and vertical scale of the font.
	// The resulting positions are computed with: fontUnit * Scale / faceUpem,
	// where faceUpem is given by the face.
	//
	// Given a device resolution (in dpi) and a point size, the scale to
	// get result in pixels is given by : pointSize * dpi / 72
	XScale, YScale int32

	// Horizontal and vertical pixels-per-em (ppem) of the font.
	// Is is used to select bitmap sizes and to perform some Opentype
	// positionning.
	XPpem, YPpem uint16
}

// NewFont constructs a new font object from the specified face.
//
// The scale is set to the face Upem, meaning that by default
// the output results will be expressed in font units.
//
// When appropriate, it will load the additional information
// required for Opentype and Graphite layout, which will influence
// the shaping plan used in `Buffer.Shape`.
//
// The `face` object should not be modified after this call.
func NewFont(face Face) *Font {
	var font Font

	font.face = face
	font.faceUpem = Position(font.face.Upem())
	font.XScale = font.faceUpem
	font.YScale = font.faceUpem

	if opentypeFace, ok := face.(FaceOpentype); ok {
		lt := opentypeFace.LayoutTables()
		font.otTables = &lt

		// accelerators
		font.gsubAccels = make([]otLayoutLookupAccelerator, len(lt.GSUB.Lookups))
		for i, l := range lt.GSUB.Lookups {
			font.gsubAccels[i].init(lookupGSUB(l))
		}
		font.gposAccels = make([]otLayoutLookupAccelerator, len(lt.GPOS.Lookups))
		for i, l := range lt.GPOS.Lookups {
			font.gposAccels[i].init(lookupGPOS(l))
		}

		if tables, is := opentypeFace.IsGraphite(); is {
			font.gr, _ = graphite.LoadGraphite(tables)
		}
	}

	return &font
}

// SetVarCoordsDesign applies a list of variation coordinates, in design-space units,
// to the font.
func (f *Font) SetVarCoordsDesign(coords []float32) {
	if varFace, ok := f.face.(FaceOpentype); ok {
		varFace.SetVarCoordinates(varFace.NormalizeVariations(coords))
	}
}

// Face returns the underlying face.
// Note that field is readonly, since some caching may happen
// in the `NewFont` constructor.
func (f *Font) Face() fonts.Face { return f.face }

func (f *Font) nominalGlyph(r rune, notFound fonts.GID) (fonts.GID, bool) {
	g, ok := f.face.NominalGlyph(r)
	if !ok {
		g = notFound
	}
	return g, ok
}

// ---- Convert from font-space to user-space ----

func (f *Font) emScaleX(v int16) Position    { return Position(v) * f.XScale / f.faceUpem }
func (f *Font) emScaleY(v int16) Position    { return Position(v) * f.YScale / f.faceUpem }
func (f *Font) emScalefX(v float32) Position { return emScalef(v, f.XScale, f.faceUpem) }
func (f *Font) emScalefY(v float32) Position { return emScalef(v, f.YScale, f.faceUpem) }
func (f *Font) emFscaleX(v int16) float32    { return emFscale(v, f.XScale, f.faceUpem) }
func (f *Font) emFscaleY(v int16) float32    { return emFscale(v, f.YScale, f.faceUpem) }

func emScalef(v float32, scale, faceUpem int32) Position {
	return roundf(v * float32(scale) / float32(faceUpem))
}

func emFscale(v int16, scale, faceUpem int32) float32 {
	return float32(v) * float32(scale) / float32(faceUpem)
}

// GlyphExtents is the same as fonts.GlyphExtents but with int type
type GlyphExtents struct {
	XBearing int32
	YBearing int32
	Width    int32
	Height   int32
}

// GlyphExtents fetches the GlyphExtents data for a glyph ID
// in the specified font, or false if not found
func (f *Font) GlyphExtents(glyph fonts.GID) (out GlyphExtents, ok bool) {
	ext, ok := f.face.GlyphExtents(glyph, f.XPpem, f.YPpem)
	if !ok {
		return out, false
	}
	out.XBearing = f.emScalefX(ext.XBearing)
	out.Width = f.emScalefX(ext.Width)
	out.YBearing = f.emScalefY(ext.YBearing)
	out.Height = f.emScalefY(ext.Height)
	return out, true
}

// GlyphAdvanceForDirection fetches the advance for a glyph ID from the specified font,
// in a text segment of the specified direction.
//
// Calls the appropriate direction-specific variant (horizontal
// or vertical) depending on the value of `dir`.
func (f *Font) GlyphAdvanceForDirection(glyph fonts.GID, dir Direction) (x, y Position) {
	if dir.isHorizontal() {
		return f.GlyphHAdvance(glyph), 0
	}
	return 0, f.getGlyphVAdvance(glyph)
}

// GlyphHAdvance fetches the advance for a glyph ID in the font,
// for horizontal text segments.
func (f *Font) GlyphHAdvance(glyph fonts.GID) Position {
	adv := f.face.HorizontalAdvance(glyph)
	return f.emScalefX(adv)
}

// Fetches the advance for a glyph ID in the font,
// for vertical text segments.
func (f *Font) getGlyphVAdvance(glyph fonts.GID) Position {
	adv := f.face.VerticalAdvance(glyph)
	return f.emScalefY(adv)
}

// Subtracts the origin coordinates from an (X,Y) point coordinate,
// in the specified glyph ID in the specified font.
//
// Calls the appropriate direction-specific variant (horizontal
// or vertical) depending on the value of @direction.
func (f *Font) subtractGlyphOriginForDirection(glyph fonts.GID, direction Direction,
	x, y Position,
) (Position, Position) {
	originX, originY := f.getGlyphOriginForDirection(glyph, direction)

	return x - originX, y - originY
}

// Fetches the (X,Y) coordinates of the origin for a glyph in
// the specified font.
//
// Calls the appropriate direction-specific variant (horizontal
// or vertical) depending on the value of @direction.
func (f *Font) getGlyphOriginForDirection(glyph fonts.GID, direction Direction) (x, y Position) {
	if direction.isHorizontal() {
		return f.getGlyphHOriginWithFallback(glyph)
	}
	return f.getGlyphVOriginWithFallback(glyph)
}

func (f *Font) getGlyphHOriginWithFallback(glyph fonts.GID) (Position, Position) {
	x, y, ok := f.face.GlyphHOrigin(glyph)
	if !ok {
		x, y, ok = f.face.GlyphVOrigin(glyph)
		if ok {
			dx, dy := f.guessVOriginMinusHOrigin(glyph)
			return x - dx, y - dy
		}
	}
	return x, y
}

func (f *Font) getGlyphVOriginWithFallback(glyph fonts.GID) (Position, Position) {
	x, y, ok := f.face.GlyphVOrigin(glyph)
	if !ok {
		x, y, ok = f.face.GlyphHOrigin(glyph)
		if ok {
			dx, dy := f.guessVOriginMinusHOrigin(glyph)
			return x + dx, y + dy
		}
	}
	return x, y
}

func (f *Font) guessVOriginMinusHOrigin(glyph fonts.GID) (x, y Position) {
	x = f.GlyphHAdvance(glyph) / 2
	y = f.getHExtendsAscender()
	return x, y
}

func (f *Font) getHExtendsAscender() Position {
	extents, ok := f.face.FontHExtents()
	if !ok {
		return f.YScale * 4 / 5
	}
	return f.emScalefY(extents.Ascender)
}

func (f *Font) hasGlyph(ch rune) bool {
	_, ok := f.face.NominalGlyph(ch)
	return ok
}

func (f *Font) subtractGlyphHOrigin(glyph fonts.GID, x, y Position) (Position, Position) {
	originX, originY := f.getGlyphHOriginWithFallback(glyph)
	return x - originX, y - originY
}

func (f *Font) subtractGlyphVOrigin(glyph fonts.GID, x, y Position) (Position, Position) {
	originX, originY := f.getGlyphVOriginWithFallback(glyph)
	return x - originX, y - originY
}

func (f *Font) addGlyphHOrigin(glyph fonts.GID, x, y Position) (Position, Position) {
	originX, originY := f.getGlyphHOriginWithFallback(glyph)
	return x + originX, y + originY
}

func (f *Font) getGlyphContourPointForOrigin(glyph fonts.GID, pointIndex uint16, direction Direction) (x, y Position, ok bool) {
	met, ok := f.face.(FaceOpentype)
	if !ok {
		return
	}

	x, y, ok = met.GetGlyphContourPoint(glyph, pointIndex)
	if ok {
		x, y = f.subtractGlyphOriginForDirection(glyph, direction, x, y)
	}

	return x, y, ok
}

// Generates gidDDD if glyph has no name.
func (f *Font) glyphToString(glyph fonts.GID) string {
	if name := f.face.GlyphName(glyph); name != "" {
		return name
	}

	return fmt.Sprintf("gid%d", glyph)
}

// ExtentsForDirection fetches the extents for a font in a text segment of the
// specified direction, applying the scaling.
//
// Calls the appropriate direction-specific variant (horizontal
// or vertical) depending on the value of `direction`.
func (f *Font) ExtentsForDirection(direction Direction) fonts.FontExtents {
	var (
		extents fonts.FontExtents
		ok      bool
	)
	if direction.isHorizontal() {
		extents, ok = f.face.FontHExtents()
		extents.Ascender = float32(f.emScalefY(extents.Ascender))
		extents.Descender = float32(f.emScalefY(extents.Descender))
		extents.LineGap = float32(f.emScalefY(extents.LineGap))
		if !ok {
			extents.Ascender = float32(f.YScale) * 0.8
			extents.Descender = extents.Ascender - float32(f.YScale)
			extents.LineGap = 0
		}
	} else {
		extents, ok = f.face.FontVExtents()
		extents.Ascender = float32(f.emScalefX(extents.Ascender))
		extents.Descender = float32(f.emScalefX(extents.Descender))
		extents.LineGap = float32(f.emScalefX(extents.LineGap))
		if !ok {
			extents.Ascender = float32(f.XScale) * 0.5
			extents.Descender = extents.Ascender - float32(f.XScale)
			extents.LineGap = 0
		}
	}
	return extents
}

// LineMetric fetches the given metric, applying potential variations
// and scaling.
func (f *Font) LineMetric(metric fonts.LineMetric) (int32, bool) {
	m, ok := f.face.LineMetric(metric)
	return f.emScalefY(m), ok
}

func (font *Font) varCoords() []float32 {
	if ot, ok := font.face.(FaceOpentype); ok {
		return ot.VarCoordinates()
	}
	return nil
}

func (font *Font) getXDelta(varStore tt.VariationStore, device tt.DeviceTable) Position {
	switch device := device.(type) {
	case tt.DeviceHinting:
		return device.GetDelta(font.XPpem, font.XScale)
	case tt.DeviceVariation:
		return font.emScalefX(varStore.GetDelta(tt.VariationStoreIndex(device), font.varCoords()))
	default:
		return 0
	}
}

func (font *Font) getYDelta(varStore tt.VariationStore, device tt.DeviceTable) Position {
	switch device := device.(type) {
	case tt.DeviceHinting:
		return device.GetDelta(font.YPpem, font.YScale)
	case tt.DeviceVariation:
		return font.emScalefY(varStore.GetDelta(tt.VariationStoreIndex(device), font.varCoords()))
	default:
		return 0
	}
}

// GetOTLayoutTables returns the OpenType layout tables, or nil
// if the underlying face is not a FaceOpentype.
// The returned tables should not be modified.
func (f *Font) GetOTLayoutTables() *tt.LayoutTables { return f.otTables }

// GetOTGlyphClass fetches the GDEF class of the requested glyph in the specified face,
// or 0 if not found.
func (f *Font) GetOTGlyphClass(glyph fonts.GID) uint32 {
	if f.otTables == nil {
		return 0
	}

	if cl := f.otTables.GDEF.Class; cl != nil {
		class, _ := cl.ClassID(glyph)
		return class
	}
	return 0
}

// GetOTLigatureCarets fetches a list of the caret positions defined for a ligature glyph in the GDEF
// table of the font (or nil if not found).
func (f *Font) GetOTLigatureCarets(direction Direction, glyph fonts.GID) []Position {
	if f.otTables == nil {
		return nil
	}

	varStore := f.otTables.GDEF.VariationStore

	list := f.otTables.GDEF.LigatureCaretList
	if list.Coverage == nil {
		return nil
	}

	index, ok := list.Coverage.Index(glyph)
	if !ok {
		return nil
	}

	glyphCarets := list.LigCarets[index]
	out := make([]Position, len(glyphCarets))
	for i, c := range glyphCarets {
		out[i] = f.getCaretValue(c, direction, glyph, varStore)
	}
	return out
}

// interpreted the CaretValue according to its format
func (f *Font) getCaretValue(caret tt.CaretValue, direction Direction, glyph fonts.GID, varStore tt.VariationStore) Position {
	switch caret := caret.(type) {
	case tt.CaretValueFormat1:
		if direction.isHorizontal() {
			return f.emScaleX(int16(caret))
		} else {
			return f.emScaleY(int16(caret))
		}
	case tt.CaretValueFormat2:
		x, y, _ := f.getGlyphContourPointForOrigin(glyph, uint16(caret), direction)
		if direction.isHorizontal() {
			return x
		} else {
			return y
		}
	case tt.CaretValueFormat3:
		if direction.isHorizontal() {
			return f.emScaleX(caret.Coordinate) + f.getXDelta(varStore, caret.Device)
		} else {
			return f.emScaleY(caret.Coordinate) + f.getYDelta(varStore, caret.Device)
		}
	default:
		return 0
	}
}
