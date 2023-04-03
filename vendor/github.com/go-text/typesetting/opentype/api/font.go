// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

// Package font provides an high level API to access
// Opentype font properties.
// See packages [loader] and [tables] for a lower level, more detailled API.
package api

import (
	"math"
)

// GID is used to identify glyphs in a font.
// It is mostly internal to the font and should not be confused with
// Unicode code points.
// Note that, despite Opentype font files using uint16, we choose to use uint32,
// to allow room for future extension.
type GID uint32

// EmptyGlyph represents an invisible glyph, which should not be drawn,
// but whose advance and offsets should still be accounted for when rendering.
const EmptyGlyph GID = math.MaxUint32

// CmapIter is an interator over a Cmap.
type CmapIter interface {
	// Next returns true if the iterator still has data to yield
	Next() bool

	// Char must be called only when `Next` has returned `true`
	Char() (rune, GID)
}

// Cmap stores a compact representation of a cmap,
// offering both on-demand rune lookup and full rune range.
// It is conceptually equivalent to a map[rune]GID, but is often
// implemented more efficiently.
type Cmap interface {
	// Iter returns a new iterator over the cmap
	// Multiple iterators may be used over the same cmap
	// The returned interface is garanted not to be nil.
	Iter() CmapIter

	// Lookup avoid the construction of a map and provides
	// an alternative when only few runes need to be fetched.
	// It returns a default value and false when no glyph is provided.
	Lookup(rune) (GID, bool)
}

// FontExtents exposes font-wide extent values, measured in font units.
// Note that typically ascender is positive and descender negative in coordinate systems that grow up.
type FontExtents struct {
	Ascender  float32 // Typographic ascender.
	Descender float32 // Typographic descender.
	LineGap   float32 // Suggested line spacing gap.
}

// LineMetric identifies one metric about the font.
type LineMetric uint8

const (
	// Distance above the baseline of the top of the underline.
	// Since most fonts have underline positions beneath the baseline, this value is typically negative.
	UnderlinePosition LineMetric = iota

	// Suggested thickness to draw for the underline.
	UnderlineThickness

	// Distance above the baseline of the top of the strikethrough.
	StrikethroughPosition

	// Suggested thickness to draw for the strikethrough.
	StrikethroughThickness

	SuperscriptEmYSize
	SuperscriptEmXOffset

	SubscriptEmYSize
	SubscriptEmYOffset
	SubscriptEmXOffset

	CapHeight
	XHeight
)

// GlyphExtents exposes extent values, measured in font units.
// Note that height is negative in coordinate systems that grow up.
type GlyphExtents struct {
	XBearing float32 // Left side of glyph from origin
	YBearing float32 // Top side of glyph from origin
	Width    float32 // Distance from left to right side
	Height   float32 // Distance from top to bottom side
}

// GlyphData describe how to graw a glyph.
// It is either an GlyphOutline, GlyphSVG or GlyphBitmap.
type GlyphData interface {
	isGlyphData()
}

func (GlyphOutline) isGlyphData() {}
func (GlyphSVG) isGlyphData()     {}
func (GlyphBitmap) isGlyphData()  {}

// GlyphOutline exposes the path to draw for
// vector glyph.
// Coordinates are expressed in fonts units.
type GlyphOutline struct {
	Segments []Segment
}

type SegmentOp uint8

const (
	SegmentOpMoveTo SegmentOp = iota
	SegmentOpLineTo
	SegmentOpQuadTo
	SegmentOpCubeTo
)

type SegmentPoint struct {
	X, Y float32 // expressed in fonts units
}

// Move translates the point.
func (pt *SegmentPoint) Move(dx, dy float32) {
	pt.X += dx
	pt.Y += dy
}

type Segment struct {
	Op SegmentOp
	// Args is up to three (x, y) coordinates, depending on the
	// operation.
	// The Y axis increases up.
	Args [3]SegmentPoint
}

// ArgsSlice returns the effective slice of points
// used (whose length is between 1 and 3).
func (s *Segment) ArgsSlice() []SegmentPoint {
	switch s.Op {
	case SegmentOpMoveTo, SegmentOpLineTo:
		return s.Args[0:1]
	case SegmentOpQuadTo:
		return s.Args[0:2]
	case SegmentOpCubeTo:
		return s.Args[0:3]
	default:
		panic("unreachable")
	}
}

// GlyphSVG is an SVG description for the glyph,
// as found in Opentype SVG table.
type GlyphSVG struct {
	// The SVG image content, decompressed if needed.
	// The actual glyph description is an SVG element
	// with id="glyph<GID>" (as in id="glyph12"),
	// and several glyphs may share the same Source
	Source []byte

	// According to the specification, a fallback outline
	// should be specified for each SVG glyphs
	Outline GlyphOutline
}

type GlyphBitmap struct {
	// The actual image content, whose interpretation depends
	// on the Format field.
	Data          []byte
	Format        BitmapFormat
	Width, Height int // number of columns and rows

	// Outline may be specified to be drawn with bitmap
	Outline *GlyphOutline
}

// BitmapFormat identifies the format on the glyph
// raw data. Across the various font files, many formats
// may be encountered : black and white bitmaps, PNG, TIFF, JPG.
type BitmapFormat uint8

const (
	_ BitmapFormat = iota
	BlackAndWhite
	PNG
	JPG
	TIFF
)

// BitmapSize expose the size of bitmap glyphs.
// One font may contain several sizes.
type BitmapSize struct {
	Height, Width uint16
	XPpem, YPpem  uint16
}

// FontID represents an identifier of a font (possibly in a collection),
// and an optional variable instance.
type FontID struct {
	File string // The filename or identifier of the font file.

	// The index of the face in a collection. It is always 0 for
	// single font files.
	Index uint16

	// For variable fonts, stores 1 + the instance index.
	// (0 to ignore variations).
	Instance uint16
}
