// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import "encoding/binary"

// Ankr is the anchor point table
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6ankr.html
type Ankr struct {
	version uint16 // Version number (set to zero)
	flags   uint16 // Flags (currently unused; set to zero)
	// Offset to the table's lookup table; currently this is always 0x0000000C
	// The lookup table returns uint16 offset from the beginning of the glyph data table, not indices.
	lookupTable AATLookup `offsetSize:"Offset32"`
	// Offset to the glyph data table
	glyphDataTable []byte `offsetSize:"Offset32" arrayCount:"ToEnd"`
}

// GetAnchor return the i-th anchor for `glyph`, or {0,0} if not found.
func (ank Ankr) GetAnchor(glyph GlyphID, index int) (anchor AnkrAnchor) {
	offset, ok := ank.lookupTable.Class(glyph)
	if !ok || int(offset)+4 >= len(ank.glyphDataTable) {
		return anchor
	}

	count := int(binary.BigEndian.Uint32(ank.glyphDataTable[offset:]))
	if index >= count {
		return anchor // invalid index
	}

	indexStart := int(offset) + 4 + 4*index
	if len(ank.glyphDataTable) < indexStart+4 {
		return anchor // invalid table
	}
	anchor.X = int16(binary.BigEndian.Uint16(ank.glyphDataTable[indexStart:]))
	anchor.Y = int16(binary.BigEndian.Uint16(ank.glyphDataTable[indexStart+2:]))
	return anchor
}

// AnkrAnchor is a point within the coordinate space of a given glyph
// independent of the control points used to render the glyph
type AnkrAnchor struct {
	X, Y int16
}
