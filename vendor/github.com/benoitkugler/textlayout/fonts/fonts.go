// Package fonts provides supports for parsing
// several font formats (postscript, bitmap and truetype)
// and provides a common API, inspired by freetype.
//
// It does not currently support CIDType1 fonts.
package fonts

import "math"

// Resource is a combination of io.Reader, io.Seeker and io.ReaderAt.
// This interface is satisfied by most things that you'd want
// to parse, for example *os.File, io.SectionReader or *bytes.Reader.
type Resource interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	Seek(int64, int) (int64, error)
}

// PSInfo exposes global properties of a postscript font.
type PSInfo struct {
	FontName    string // Postscript font name.
	FullName    string // full name of the font.
	FamilyName  string // family name of the font.
	Version     string // font program version identifier (optional)
	Notice      string // font name trademark or copyright notice (optional)
	Weight      string // Weight of the font: normal, bold, etc.
	ItalicAngle int    // italic angle of the font, usually 0. or negative.

	IsFixedPitch bool // true if all the characters have the same width.

	UnderlinePosition  int
	UnderlineThickness int
}

// FontSummary stores basic informations about the
// style of the font.
type FontSummary struct {
	Familly          string
	Style            string
	IsItalic, IsBold bool

	HasScalableGlyphs, HasBitmapGlyphs, HasColorGlyphs bool
}

// FaceMetadata exposes some summary information about the font.
type FaceMetadata interface {
	// Cmap returns the mapping between input character codes
	// and glyph ids. The returned encoding identifies which system is used
	// to describe characters.
	Cmap() (Cmap, CmapEncoding)

	// PostscriptInfo returns the Postscript information,
	// or false is not available.
	PostscriptInfo() (PSInfo, bool)

	// PoscriptName returns the PoscriptName of the font,
	// or an empty string.
	PoscriptName() string

	// LoadSummary fetchs global details about the font.
	// Conceptually, this method just returns it receiver, but this enables lazy loading.
	LoadSummary() (FontSummary, error)

	// LoadBitmaps returns the available bitmap sizes, or an empty
	// slice for outline fonts.
	LoadBitmaps() []BitmapSize
}

// Face provides a unified access to various font formats.
// It describes the content of one font from a font file.
// Implementation must be pointer to simplify caching and hashing.
type Face interface {
	FaceMetadata
	FaceMetrics
	FaceRenderer
}

// Faces is the parsed content of a font ressource.
// Note that variable fonts are not repeated in this slice,
// since instances are accessed on each font.
type Faces = []Face

// FontLoader implements the general parsing
// of a font file. Some font format support to store several
// fonts inside one file. For the other formats, the returned slice will
// have length 1.
type FontLoader = func(file Resource) (Faces, error)

// GID is used to identify glyphs in a font.
// It is mostly internal to the font and should not be confused with
// Unicode code points.
type GID uint32

// EmptyGlyph represents an invisible glyph, which should not be drawn,
// but whose advance and offsets should still be accounted for when rendering.
const EmptyGlyph GID = math.MaxUint32

type CmapEncoding uint8

const (
	EncOther CmapEncoding = iota
	EncUnicode
	EncSymbol
)

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
	// The returned interface is garanted not to be nil
	Iter() CmapIter

	// Lookup avoid the construction of a map and provides
	// an alternative when only few runes need to be fetched.
	// It returns a default value and false when no glyph is provided.
	Lookup(rune) (GID, bool)
}

var (
	_ Cmap     = CmapSimple(nil)
	_ CmapIter = (*cmap0Iter)(nil)
)

// CmapSimple is a map based Cmap implementation.
type CmapSimple map[rune]GID

type cmap0Iter struct {
	data CmapSimple
	keys []rune
	pos  int
}

func (it *cmap0Iter) Next() bool {
	return it.pos < len(it.keys)
}

func (it *cmap0Iter) Char() (rune, GID) {
	r := it.keys[it.pos]
	it.pos++
	return r, it.data[r]
}

func (s CmapSimple) Iter() CmapIter {
	keys := make([]rune, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return &cmap0Iter{data: s, keys: keys}
}

func (s CmapSimple) Lookup(r rune) (GID, bool) {
	v, ok := s[r] // will be 0 if r is not in s
	return v, ok
}

// FontExtents exposes font-wide extent values, measured in font units.
// Note that typically ascender is positive and descender negative in coordinate systems that grow up.
type FontExtents struct {
	Ascender  float32 // Typographic ascender.
	Descender float32 // Typographic descender.
	LineGap   float32 // Suggested line spacing gap.
}

// LineMetric identifies one metric about the font.
// Some formats only support a subset of the metrics defined by the constants.
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

// FaceMetrics exposes details of the font content.
// Implementation must be valid map keys to simplify caching.
type FaceMetrics interface {
	// Upem returns the units per em of the font file.
	// If not found, should return 1000 as fallback value.
	// This value is only relevant for scalable fonts.
	Upem() uint16

	// GlyphName returns the name of the given glyph, or an empty
	// string if the glyph is invalid or has no name.
	GlyphName(gid GID) string

	// LineMetric returns the metric identified by `metric` (in fonts units), or false
	// if the font does not provide such information.
	LineMetric(metric LineMetric) (float32, bool)

	// FontHExtents returns the extents of the font for horizontal text, or false
	// it not available, in font units.
	// `varCoords` (in normalized coordinates) is only useful for variable fonts.
	FontHExtents() (FontExtents, bool)

	// FontVExtents is the same as `FontHExtents`, but for vertical text.
	FontVExtents() (FontExtents, bool)

	// NominalGlyph returns the glyph used to represent the given rune,
	// or false if not found.
	NominalGlyph(ch rune) (GID, bool)

	// HorizontalAdvance returns the horizontal advance in font units.
	// When no data is available but the glyph index is valid, this method
	// should return a default value (the upem number for example).
	// If the glyph is invalid it should return 0.
	// `coords` is used by variable fonts, and is specified in normalized coordinates.
	HorizontalAdvance(gid GID) float32

	// VerticalAdvance is the same as `HorizontalAdvance`, but for vertical advance.
	VerticalAdvance(gid GID) float32

	// GlyphHOrigin fetches the (X,Y) coordinates of the origin (in font units) for a glyph ID,
	// for horizontal text segments.
	// Returns `false` if not available.
	GlyphHOrigin(GID) (x, y int32, found bool)

	// GlyphVOrigin is the same as `GlyphHOrigin`, but for vertical text segments.
	GlyphVOrigin(GID) (x, y int32, found bool)

	// GlyphExtents retrieve the extents for a specified glyph, of false, if not available.
	// `coords` is used by variable fonts, and is specified in normalized coordinates.
	// For bitmap glyphs, the closest resolution to `xPpem` and `yPpem` is selected.
	GlyphExtents(glyph GID, xPpem, yPpem uint16) (GlyphExtents, bool)
}

// FaceRenderer exposes access to glyph contents
type FaceRenderer interface {
	// GlyphData loads the glyph content, or return nil
	// if 'gid' is not supported.
	// For bitmap glyphs, the closest resolution to `xPpem` and `yPpem` is selected.
	GlyphData(gid GID, xPpem, yPpem uint16) GlyphData
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

// FaceID represents an identifier of a face (possibly in a collection),
// and an optional variable instance.
type FaceID struct {
	File string // The filename or identifier of the font file.

	// The index of the face in a collection. It is always 0 for
	// single font files.
	Index uint16

	// For variable fonts, stores 1 + the instance index.
	// (0 to ignore variations).
	Instance uint16
}

// FaceDescription is a summary of a font file.
type FaceDescription struct {
	Family string

	// TODO: more fields will be added
}

// Style (also called slant) allows italic or oblique faces to be selected.
type Style uint8

const (
	// A face that is neither italic not obliqued.
	StyleNormal Style = iota + 1
	// A form that is generally cursive in nature.
	StyleItalic
	// A typically-sloped version of the regular face.
	StyleOblique
)

// Weight is the degree of blackness or stroke thickness of a font.
// This value ranges from 100.0 to 900.0, with 400.0 as normal.
type Weight float32

const (
	// Thin weight (100), the thinnest value.
	WeightThin Weight = 100
	// Extra light weight (200).
	WeightExtraLight Weight = 200
	// Light weight (300).
	WeightLight Weight = 300
	// Normal (400).
	WeightNormal Weight = 400
	// Medium weight (500, higher than normal).
	WeightMedium Weight = 500
	// Semibold weight (600).
	WeightSemibold Weight = 600
	// Bold weight (700).
	WeightBold Weight = 700
	// Extra-bold weight (800).
	WeightExtraBold Weight = 800
	// Black weight (900), the thickest value.
	WeightBlack Weight = 900
)

// Stretch is the width of a font as an approximate fraction of the normal width.
// Widths range from 0.5 to 2.0 inclusive, with 1.0 as the normal width.
type Stretch float32

const (
	// Ultra-condensed width (50%), the narrowest possible.
	StretchUltraCondensed Stretch = 0.5
	// Extra-condensed width (62.5%).
	StretchExtraCondensed Stretch = 0.625
	// Condensed width (75%).
	StretchCondensed Stretch = 0.75
	// Semi-condensed width (87.5%).
	StretchSemiCondensed Stretch = 0.875
	// Normal width (100%).
	StretchNormal Stretch = 1.0
	// Semi-expanded width (112.5%).
	StretchSemiExpanded Stretch = 1.125
	// Expanded width (125%).
	StretchExpanded Stretch = 1.25
	// Extra-expanded width (150%).
	StretchExtraExpanded Stretch = 1.5
	// Ultra-expanded width (200%), the widest possible.
	StretchUltraExpanded Stretch = 2.0
)

// FontDescriptor is a handle on a font, able to efficiently query
// some global information.
type FontDescriptor interface {
	// Family queries the font family name.
	Family() string

	// Aspect queries the visual properties of the font.
	// If not found, zero values should be returned.
	Aspect() (Style, Weight, Stretch)

	// AdditionalStyle returns a description of the style of the font,
	// including information not found by Aspect()
	AdditionalStyle() string

	// Cmap returns the Unicode to Glyph mapping
	LoadCmap() (Cmap, error)
}
