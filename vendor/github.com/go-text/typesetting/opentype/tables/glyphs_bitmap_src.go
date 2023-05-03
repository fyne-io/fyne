// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import "fmt"

// CBLC is the Color Bitmap Location Table
// See - https://learn.microsoft.com/fr-fr/typography/opentype/spec/cblc
type CBLC struct {
	majorVersion   uint16             //	Major version of the CBLC table, = 3.
	minorVersion   uint16             //	Minor version of the CBLC table, = 0.
	BitmapSizes    []BitmapSize       `arrayCount:"FirstUint32"` // BitmapSize records array.
	IndexSubTables [][]BitmapSubtable `isOpaque:""`              // with same length as [BitmapSizes]
}

func (cb *CBLC) parseIndexSubTables(src []byte) error {
	cb.IndexSubTables = make([][]BitmapSubtable, len(cb.BitmapSizes))
	for i, size := range cb.BitmapSizes {
		start := int(size.indexSubTableArrayOffset)
		if L := len(src); L < start {
			return fmt.Errorf("EOF: expected length: %d, got %d", start, L)
		}
		subtables, _, err := ParseIndexSubTableArray(src[start:], int(size.numberOfIndexSubTables))
		if err != nil {
			return err
		}
		sizeSubtables := make([]BitmapSubtable, len(subtables.Subtables))
		for j, subtable := range subtables.Subtables {
			numGlyphs := int(subtable.LastGlyph) - int(subtable.FirstGlyph) + 1
			subtableStart := start + int(subtable.additionalOffsetToIndexSubtable)

			sizeSubtables[j].FirstGlyph = subtable.FirstGlyph
			sizeSubtables[j].LastGlyph = subtable.LastGlyph
			sizeSubtables[j].IndexSubHeader, _, err = ParseIndexSubHeader(src[subtableStart:], numGlyphs+1)
			if err != nil {
				return err
			}
		}
		cb.IndexSubTables[i] = sizeSubtables
	}
	return nil
}

type BitmapSize struct {
	indexSubTableArrayOffset Offset32        //	Offset to index subtable from beginning of CBLC.
	indexTablesSize          uint32          //	Number of bytes in corresponding index subtables and array.
	numberOfIndexSubTables   uint32          //	There is an index subtable for each range or format change.
	colorRef                 uint32          //	Not used; set to 0.
	Hori                     SbitLineMetrics //	Line metrics for text rendered horizontally.
	Vert                     SbitLineMetrics //	Line metrics for text rendered vertically.
	startGlyphIndex          uint16          //	Lowest glyph index for this size.
	endGlyphIndex            uint16          //	Highest glyph index for this size.
	PpemX                    uint8           //	Horizontal pixels per em.
	PpemY                    uint8           //	Vertical pixels per em.
	bitDepth                 uint8           //	In addtition to already defined bitDepth values 1, 2, 4, and 8 supported by existing implementations, the value of 32 is used to identify color bitmaps with 8 bit per pixel RGBA channels.
	flags                    int8            //	Vertical or horizontal (see the Bitmap Flags section of the EBLC table chapter).
}

type SbitLineMetrics struct {
	Ascender              int8
	Descender             int8
	widthMax              uint8
	caretSlopeNumerator   int8
	caretSlopeDenominator int8
	caretOffset           int8
	minOriginSB           int8
	minAdvanceSB          int8
	MaxBeforeBL           int8
	MinAfterBL            int8
	pad1                  int8
	pad2                  int8
}

type IndexSubTableArray struct {
	Subtables []IndexSubTableHeader
}

type IndexSubTableHeader struct {
	FirstGlyph                      GlyphID  //	First glyph ID of this range.
	LastGlyph                       GlyphID  //	Last glyph ID of this range (inclusive).
	additionalOffsetToIndexSubtable Offset32 //	Add to indexSubTableArrayOffset to get offset from beginning of EBLC.
}

type IndexSubHeader struct {
	indexFormat     indexVersion // Format of this IndexSubTable.
	ImageFormat     uint16       // Format of EBDT image data.
	ImageDataOffset Offset32     // Offset to image data in EBDT table.
	IndexData       IndexData    `unionField:"indexFormat"`
}

type indexVersion uint16

const (
	indexVersion1 indexVersion = iota + 1
	indexVersion2
	indexVersion3
	indexVersion4
	indexVersion5
)

type IndexData interface {
	isIndexData()
}

func (IndexData1) isIndexData() {}
func (IndexData2) isIndexData() {}
func (IndexData3) isIndexData() {}
func (IndexData4) isIndexData() {}
func (IndexData5) isIndexData() {}

type IndexData1 struct {
	// sizeOfArray = (lastGlyph - firstGlyph + 1) + 1 + 1 pad if needed
	// sbitOffsets[glyphIndex] + imageDataOffset = glyphData
	SbitOffsets []Offset32
}

type IndexData2 struct {
	ImageSize  uint32          // All the glyphs are of the same size.
	BigMetrics BigGlyphMetrics // All glyphs have the same metrics; glyph data may be compressed, byte-aligned, or bit-aligned.
}

type IndexData3 struct {
	// sizeOfArray = (lastGlyph - firstGlyph + 1) + 1 + 1 pad if needed
	// sbitOffets[glyphIndex] + imageDataOffset = glyphData
	SbitOffsets []Offset16
}

type IndexData4 struct {
	numGlyphs  uint32              //	Array length.
	GlyphArray []GlyphIdOffsetPair `arrayCount:"ComputedField-numGlyphs+1"` //[numGlyphs + 1]	One per glyph.
}

type GlyphIdOffsetPair struct {
	GlyphID    GlyphID  //	Glyph ID of glyph present.
	SbitOffset Offset16 //	Location in EBDT.
}

type IndexData5 struct {
	ImageSize    uint32          //	All glyphs have the same data size.
	BigMetrics   BigGlyphMetrics //	All glyphs have the same metrics.
	GlyphIdArray []GlyphID       `arrayCount:"FirstUint32"` // [numGlyphs] One per glyph, sorted by glyph ID.
}

// ------------------------- actual data : shared by EBDT / CBDT / BDAT -------------------------
// for now, we simplify the implementation to two cases:
//	- data, metrics (small)
//  - data only

type SmallGlyphMetrics struct {
	Height   uint8 // Number of rows of data.
	Width    uint8 // Number of columns of data.
	BearingX int8  // Distance in pixels from the horizontal origin to the left edge of the bitmap.
	BearingY int8  // Distance in pixels from the horizontal origin to the top edge of the bitmap.
	Advance  uint8 // Horizontal advance width in pixels.
}

type BigGlyphMetrics struct {
	SmallGlyphMetrics
	vertBearingX int8  // Distance in pixels from the vertical origin to the left edge of the bitmap.
	vertBearingY int8  // Distance in pixels from the vertical origin to the top edge of the bitmap.
	vertAdvance  uint8 // Vertical advance width in pixels.
}

// Format 2: small metrics, bit-aligned data
type BitmapData2 struct {
	SmallGlyphMetrics
	Image []byte `arrayCount:"ToEnd"`
}

// Format 5: metrics in CBLC table, bit-aligned image data only
type BitmapData5 struct {
	Image []byte `arrayCount:"ToEnd"`
}

// Format 17: small metrics, PNG image data
type BitmapData17 struct {
	SmallGlyphMetrics
	Image []byte `arrayCount:"FirstUint32"`
}

// Format 18: big metrics, PNG image data
type BitmapData18 struct {
	BigGlyphMetrics
	Image []byte `arrayCount:"FirstUint32"`
}

// Format 19: metrics in CBLC table, PNG image data
type BitmapData19 struct {
	Image []byte `arrayCount:"FirstUint32"`
}
