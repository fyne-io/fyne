// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"fmt"
)

// Sbix is the Standard Bitmap Graphics Table
// See - https://learn.microsoft.com/fr-fr/typography/opentype/spec/sbix
type Sbix struct {
	version uint16 //	Table version number — set to 1
	// Bit 0: Set to 1.
	// Bit 1: Draw outlines.
	// Bits 2 to 15: reserved (set to 0).
	Flags   uint16
	Strikes []Strike `arrayCount:"FirstUint32" offsetsArray:"Offset32"` // [numStrikes]	Offsets from the beginning of the 'sbix' table to data for each individual bitmap strike.
}

// Strike stores one size of bitmap glyphs in the 'sbix' table.
// binarygen: argument=numGlyphs int
type Strike struct {
	Ppem       uint16            // The PPEM size for which this strike was designed.
	Ppi        uint16            // The device pixel density (in PPI) for which this strike was designed. (E.g., 96 PPI, 192 PPI.)
	GlyphDatas []BitmapGlyphData `isOpaque:""` //[numGlyphs+1] Offset from the beginning of the strike data header to bitmap data for an individual glyph ID.
}

func (st *Strike) parseGlyphDatas(src []byte, numGlyphs int) error {
	const headerSize = 4
	offsets, err := ParseLoca(src[headerSize:], numGlyphs, true)
	if err != nil {
		return err
	}
	st.GlyphDatas = make([]BitmapGlyphData, numGlyphs)
	for i := range st.GlyphDatas {
		start, end := offsets[i], offsets[i+1]
		if start == end { // no data
			continue
		}

		if start > end {
			return fmt.Errorf("invalid strike offsets %d > %d", start, end)
		}

		if L := len(src); L < int(end) {
			return fmt.Errorf("EOF: expected length: %d, got %d", end, L)
		}

		st.GlyphDatas[i], _, err = ParseBitmapGlyphData(src[start:end])
		if err != nil {
			return err
		}
	}
	return nil
}

type BitmapGlyphData struct {
	OriginOffsetX int16  //	The horizontal (x-axis) position of the left edge of the bitmap graphic in relation to the glyph design space origin.
	OriginOffsetY int16  //	The vertical (y-axis) position of the bottom edge of the bitmap graphic in relation to the glyph design space origin.
	GraphicType   Tag    //	Indicates the format of the embedded graphic data: one of 'jpg ', 'png ' or 'tiff', or the special format 'dupe'.
	Data          []byte `arrayCount:"ToEnd"` //	The actual embedded graphic data. The total length is inferred from sequential entries in the glyphDataOffsets array and the fixed size (8 bytes) of the preceding fields.
}
