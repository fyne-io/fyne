// Package truetype provides support for OpenType and TrueType font formats, used in PDF.
//
// It is largely influenced by github.com/ConradIrwin/font and golang.org/x/image/font/sfnt,
// and FreeType2.
package truetype

import (
	"errors"

	"github.com/benoitkugler/textlayout/fonts"
	type1c "github.com/benoitkugler/textlayout/fonts/type1C"
)

var _ fonts.Face = (*Font)(nil)

type fixed struct {
	Major int16
	Minor uint16
}

type longdatetime struct {
	SecondsSince1904 uint64
}

var (
	// errUnsupportedFormat is returned from Parse if parsing failed
	errUnsupportedFormat = errors.New("unsupported font format")

	// errMissingTable is returned from *Table if the table does not exist in the font.
	errMissingTable = errors.New("missing table")

	errUnsupportedTableOffsetLength = errors.New("unsupported table offset or length")
	errInvalidDfont                 = errors.New("invalid dfont")
)

type gid = uint16

// Font represents a SFNT font, which is the underlying representation found
// in .otf and .ttf files.
// SFNT is a container format, which contains a number of tables identified by
// Tags. Depending on the type of glyphs embedded in the file which tables will
// exist. In particular, there's a big different between TrueType glyphs (usually .ttf)
// and CFF/PostScript Type 2 glyphs (usually .otf)
type Font struct {
	cmap         Cmap
	cmapVar      unicodeVariations
	cmapEncoding fonts.CmapEncoding

	Names TableName

	hhea, vhea *TableHVhea
	vorg       *tableVorg // optional
	cff        *type1c.Font
	post       TablePost // optional
	svg        tableSVG  // optional

	// Optionnal, only present in variable fonts

	varCoords  []float32   // coordinates in usage, may be nil
	hvar, vvar *tableHVvar // optional
	avar       tableAvar
	mvar       TableMvar
	gvar       tableGvar
	fvar       TableFvar

	Glyf       TableGlyf
	vmtx, Hmtx TableHVmtx
	bitmap     bitmapTable // CBDT or EBLC or BLOC
	sbix       tableSbix

	OS2 *TableOS2 // optional

	// graphite font, optionnal
	Graphite *GraphiteTables

	// Advanced layout tables.
	layoutTables LayoutTables

	fontSummary fontSummary

	Head TableHead

	// NumGlyphs exposes the number of glyph indexes present in the font,
	// as exposed in the 'maxp' table.
	NumGlyphs int // TODO: check usage

	// Type represents the kind of glyphs in this font.
	// It is one of TypeTrueType, TypeTrueTypeApple, TypePostScript1, TypeOpenType
	Type Tag

	upem uint16 // cached value

	// HasHint is true if the font has a prep table.
	HasHint bool
}

// LayoutTables exposes advanced layout tables.
// All the fields are optionnals.
type LayoutTables struct {
	GDEF TableGDEF // An absent table has a nil Class
	Trak TableTrak
	Ankr TableAnkr
	Feat TableFeat
	Morx TableMorx
	Kern TableKernx
	Kerx TableKernx
	GSUB TableGSUB // An absent table has a nil slice of lookups
	GPOS TableGPOS // An absent table has a nil slice of lookups
}

// LayoutTables returns the valid advanced layout tables.
// When parsing yields an error, it is ignored and an empty table is returned.
// See the individual methods for more control over error handling.
func (font *Font) LayoutTables() LayoutTables { return font.layoutTables }
