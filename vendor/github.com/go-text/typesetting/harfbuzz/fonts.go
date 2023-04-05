package harfbuzz

import (
	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from src/hb-font.hh, src/hb-font.cc  Copyright © 2009  Red Hat, Inc., 2012  Google, Inc.  Behdad Esfahbod

type Face = *font.Face

// Font is used internally as a light wrapper around the provided Face.
//
// Font are constructed with `NewFont` and adjusted by accessing the fields
// Ptem, XScale, YScale.
type Font struct {
	face Face

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
}

// NewFont constructs a new font object from the specified face.
//
// The scale is set to the face Upem, meaning that by default
// the output results will be expressed in font units.
//
// The `face` object should not be modified after this call.
func NewFont(face Face) *Font {
	var font Font

	font.face = face
	font.faceUpem = Position(font.face.Upem())
	font.XScale = font.faceUpem
	font.YScale = font.faceUpem

	// accelerators
	font.gsubAccels = make([]otLayoutLookupAccelerator, len(face.GSUB.Lookups))
	for i, l := range face.GSUB.Lookups {
		font.gsubAccels[i].init(lookupGSUB(l))
	}
	font.gposAccels = make([]otLayoutLookupAccelerator, len(face.GPOS.Lookups))
	for i, l := range face.GPOS.Lookups {
		font.gposAccels[i].init(lookupGPOS(l))
	}

	return &font
}

// SetVarCoordsDesign applies a list of variation coordinates, in design-space units,
// to the font.
func (f *Font) SetVarCoordsDesign(coords []float32) {
	f.face.Coords = f.face.NormalizeVariations(coords)
}

// Face returns the underlying face.
// Note that field is readonly, since some caching may happen
// in the `NewFont` constructor.
func (f *Font) Face() Face { return f.face }

func (f *Font) nominalGlyph(r rune, notFound GID) (GID, bool) {
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
func (f *Font) GlyphExtents(glyph GID) (out GlyphExtents, ok bool) {
	ext, ok := f.face.GlyphExtents(glyph)
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
func (f *Font) GlyphAdvanceForDirection(glyph GID, dir Direction) (x, y Position) {
	if dir.isHorizontal() {
		return f.GlyphHAdvance(glyph), 0
	}
	return 0, f.getGlyphVAdvance(glyph)
}

// GlyphHAdvance fetches the advance for a glyph ID in the font,
// for horizontal text segments.
func (f *Font) GlyphHAdvance(glyph GID) Position {
	adv := f.face.HorizontalAdvance(glyph)
	return f.emScalefX(adv)
}

// Fetches the advance for a glyph ID in the font,
// for vertical text segments.
func (f *Font) getGlyphVAdvance(glyph GID) Position {
	adv := f.face.VerticalAdvance(glyph)
	return f.emScalefY(adv)
}

// Subtracts the origin coordinates from an (X,Y) point coordinate,
// in the specified glyph ID in the specified font.
//
// Calls the appropriate direction-specific variant (horizontal
// or vertical) depending on the value of @direction.
func (f *Font) subtractGlyphOriginForDirection(glyph GID, direction Direction,
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
func (f *Font) getGlyphOriginForDirection(glyph GID, direction Direction) (x, y Position) {
	if direction.isHorizontal() {
		return f.getGlyphHOriginWithFallback(glyph)
	}
	return f.getGlyphVOriginWithFallback(glyph)
}

func (f *Font) getGlyphHOriginWithFallback(glyph GID) (Position, Position) {
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

func (f *Font) getGlyphVOriginWithFallback(glyph GID) (Position, Position) {
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

func (f *Font) guessVOriginMinusHOrigin(glyph GID) (x, y Position) {
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

func (f *Font) subtractGlyphHOrigin(glyph GID, x, y Position) (Position, Position) {
	originX, originY := f.getGlyphHOriginWithFallback(glyph)
	return x - originX, y - originY
}

func (f *Font) subtractGlyphVOrigin(glyph GID, x, y Position) (Position, Position) {
	originX, originY := f.getGlyphVOriginWithFallback(glyph)
	return x - originX, y - originY
}

func (f *Font) addGlyphHOrigin(glyph GID, x, y Position) (Position, Position) {
	originX, originY := f.getGlyphHOriginWithFallback(glyph)
	return x + originX, y + originY
}

func (f *Font) getGlyphContourPointForOrigin(glyph GID, pointIndex uint16, direction Direction) (x, y Position, ok bool) {
	x, y, ok = f.face.GetGlyphContourPoint(glyph, pointIndex)
	if ok {
		x, y = f.subtractGlyphOriginForDirection(glyph, direction, x, y)
	}

	return x, y, ok
}

// ExtentsForDirection fetches the extents for a font in a text segment of the
// specified direction, applying the scaling.
//
// Calls the appropriate direction-specific variant (horizontal
// or vertical) depending on the value of `direction`.
func (f *Font) ExtentsForDirection(direction Direction) api.FontExtents {
	var (
		extents api.FontExtents
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

func (font *Font) varCoords() []float32 { return font.face.Coords }

func (font *Font) getXDelta(varStore tables.ItemVarStore, device tables.DeviceTable) Position {
	switch device := device.(type) {
	case tables.DeviceHinting:
		return device.GetDelta(font.face.XPpem, font.XScale)
	case tables.DeviceVariation:
		return font.emScalefX(varStore.GetDelta(tables.VariationStoreIndex(device), font.varCoords()))
	default:
		return 0
	}
}

func (font *Font) getYDelta(varStore tables.ItemVarStore, device tables.DeviceTable) Position {
	switch device := device.(type) {
	case tables.DeviceHinting:
		return device.GetDelta(font.face.YPpem, font.YScale)
	case tables.DeviceVariation:
		return font.emScalefY(varStore.GetDelta(tables.VariationStoreIndex(device), font.varCoords()))
	default:
		return 0
	}
}

// GetOTLigatureCarets fetches a list of the caret positions defined for a ligature glyph in the GDEF
// table of the font (or nil if not found).
func (f *Font) GetOTLigatureCarets(direction Direction, glyph GID) []Position {
	varStore := f.face.GDEF.ItemVarStore

	list := f.face.GDEF.LigCaretList
	if list.Coverage == nil {
		return nil
	}

	index, ok := list.Coverage.Index(gID(glyph))
	if !ok {
		return nil
	}

	glyphCarets := list.LigGlyphs[index].CaretValues
	out := make([]Position, len(glyphCarets))
	for i, c := range glyphCarets {
		out[i] = f.getCaretValue(c, direction, glyph, varStore)
	}
	return out
}

// interpreted the CaretValue according to its format
func (f *Font) getCaretValue(caret tables.CaretValue, direction Direction, glyph GID, varStore tables.ItemVarStore) Position {
	switch caret := caret.(type) {
	case tables.CaretValue1:
		if direction.isHorizontal() {
			return f.emScaleX(int16(caret.Coordinate))
		} else {
			return f.emScaleY(int16(caret.Coordinate))
		}
	case tables.CaretValue2:
		x, y, _ := f.getGlyphContourPointForOrigin(glyph, uint16(caret.CaretValuePointIndex), direction)
		if direction.isHorizontal() {
			return x
		} else {
			return y
		}
	case tables.CaretValue3:
		if direction.isHorizontal() {
			return f.emScaleX(caret.Coordinate) + f.getXDelta(varStore, caret.Device)
		} else {
			return f.emScaleY(caret.Coordinate) + f.getYDelta(varStore, caret.Device)
		}
	default:
		return 0
	}
}
