// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// SVG is the SVG (Scalable Vector Graphics) table.
// See - https://learn.microsoft.com/fr-fr/typography/opentype/spec/svg
type SVG struct {
	version         uint16          // Table version (starting at 0). Set to 0.
	SVGDocumentList SVGDocumentList `offsetSize:"Offset32"` // Offset to the SVG Document List, from the start of the SVG table. Must be non-zero.
	reserved        uint32          // Set to 0.
}

type SVGDocumentList struct {
	DocumentRecords []SVGDocumentRecord `arrayCount:"FirstUint16"` // [numEntries]	Array of SVG document records.
	SVGRawData      []byte              `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

// Each SVG document record specifies a range of glyph IDs (from startGlyphID to endGlyphID, inclusive), and the location of its associated SVG document in the SVG table.
type SVGDocumentRecord struct {
	StartGlyphID GlyphID  // The first glyph ID for the range covered by this record.
	EndGlyphID   GlyphID  // The last glyph ID for the range covered by this record.
	SvgDocOffset Offset32 // Offset from the beginning of the SVGDocumentList to an SVG document. Must be non-zero.
	SvgDocLength uint32   // Length of the SVG document data. Must be non-zero.
}

// CFF is the Compact Font Format Table.
// Since it used its own format, quite different from the regular Opentype format,
// its interpretation is handled externally (see font/cff).
// See also https://learn.microsoft.com/fr-fr/typography/opentype/spec/cff
type CFF = []byte

// VORG is the Vertical Origin Table
// See - https://learn.microsoft.com/fr-fr/typography/opentype/spec/vorg
type VORG struct {
	majorVersion       uint16              //	Major version (starting at 1). Set to 1.
	minorVersion       uint16              //	Minor version (starting at 0). Set to 0.
	DefaultVertOriginY int16               //	The y coordinate of a glyph’s vertical origin, in the font’s design coordinate system, to be used if no entry is present for the glyph in the vertOriginYMetrics array.
	VertOriginYMetrics []VertOriginYMetric `arrayCount:"FirstUint16"`
}

// YOrigin returns the vertical origin for [glyph].
func (t *VORG) YOrigin(glyph GlyphID) int16 {
	// binary search
	for i, j := 0, len(t.VertOriginYMetrics); i < j; {
		h := i + (j-i)/2
		entry := t.VertOriginYMetrics[h]
		if glyph < entry.GlyphIndex {
			j = h
		} else if entry.GlyphIndex < glyph {
			i = h + 1
		} else {
			return entry.VertOriginY
		}
	}
	return t.DefaultVertOriginY
}

type VertOriginYMetric struct {
	GlyphIndex  GlyphID //  Glyph index.
	VertOriginY int16   //  Y coordinate, in the font’s design coordinate system, of the vertical origin of glyph with index glyphIndex.
}
