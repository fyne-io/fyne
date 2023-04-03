// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import (
	"fmt"

	"github.com/go-text/typesetting/opentype/tables"
)

type svg []svgDocument

func newSvg(table tables.SVG) (svg, error) {
	rawData := table.SVGDocumentList.SVGRawData
	out := make(svg, len(table.SVGDocumentList.DocumentRecords))
	for i, rec := range table.SVGDocumentList.DocumentRecords {
		start, end := rec.SvgDocOffset, rec.SvgDocOffset+tables.Offset32(rec.SvgDocLength)
		if len(rawData) < int(end) {
			return nil, fmt.Errorf("invalid svg table (EOF: expected %d, got %d)", end, len(rawData))
		}
		out[i] = svgDocument{
			first: rec.StartGlyphID,
			last:  rec.EndGlyphID,
			svg:   rawData[start:end],
		}
	}
	return out, nil
}

type svgDocument struct {
	// svg document
	// each glyph description must be written
	// in an element with id=glyphXXX
	svg   []byte
	first gID // The first glyph ID in the range described by this index entry.
	last  gID // The last glyph ID in the range described by this index entry. Must be >= startGlyphID.
}

// rawGlyphData returns the SVG document for [gid], or false.
func (s svg) rawGlyphData(gid gID) ([]byte, bool) {
	// binary search
	for i, j := 0, len(s); i < j; {
		h := i + (j-i)/2
		entry := s[h]
		if gid < entry.first {
			j = h
		} else if entry.last < gid {
			i = h + 1
		} else {
			return entry.svg, true
		}
	}
	return nil, false
}
