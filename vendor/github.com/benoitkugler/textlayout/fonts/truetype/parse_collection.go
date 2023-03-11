package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/benoitkugler/textlayout/fonts"
)

const maxNumFonts = 1024 // security implementation limit

// returns the offsets of each font
func parseTTCHeader(r io.Reader) ([]uint32, error) {
	// The https://www.microsoft.com/typography/otspec/otff.htm "Font
	// Collections" section describes the TTC header.
	var buf [12]byte
	if _, err := r.Read(buf[:]); err != nil {
		return nil, err
	}
	// skip versions
	numFonts := binary.BigEndian.Uint32(buf[8:])
	if numFonts == 0 {
		return nil, errors.New("empty font collection")
	}
	if numFonts > maxNumFonts {
		return nil, fmt.Errorf("number of fonts (%d) in collection exceed implementation limit (%d)",
			numFonts, maxNumFonts)
	}

	offsetsBytes := make([]byte, numFonts*4)
	_, err := io.ReadFull(r, offsetsBytes)
	if err != nil {
		return nil, err
	}
	return parseUint32s(offsetsBytes, int(numFonts)), nil
}

// parseDfont parses a dfont resource map, as per
// https://github.com/kreativekorp/ksfl/wiki/Macintosh-Resource-File-Format
//
// That unofficial wiki page lists all of its fields as *signed* integers,
// which looks unusual. The actual file format might use *unsigned* integers in
// various places, but until we have either an official specification or an
// actual dfont file where this matters, we'll use signed integers and treat
// negative values as invalid.
func parseDfont(r fonts.Resource) ([]uint32, error) {
	var buf [16]byte
	if _, err := r.Read(buf[:]); err != nil {
		return nil, err
	}
	resourceMapOffset := binary.BigEndian.Uint32(buf[4:])
	resourceMapLength := binary.BigEndian.Uint32(buf[12:])

	const (
		// (maxTableOffset + maxTableLength) will not overflow an int32.
		maxTableLength = 1 << 29
		maxTableOffset = 1 << 29
	)
	if resourceMapOffset > maxTableOffset || resourceMapLength > maxTableLength {
		return nil, errUnsupportedTableOffsetLength
	}

	const headerSize = 28
	if resourceMapLength < headerSize {
		return nil, errInvalidDfont
	}
	_, err := r.ReadAt(buf[:2], int64(resourceMapOffset+24))
	if err != nil {
		return nil, err
	}
	typeListOffset := int64(int16(binary.BigEndian.Uint16(buf[:])))
	if typeListOffset < headerSize || resourceMapLength < uint32(typeListOffset)+2 {
		return nil, errInvalidDfont
	}
	_, err = r.ReadAt(buf[:2], int64(resourceMapOffset)+typeListOffset)
	if err != nil {
		return nil, err
	}
	typeCount := int(binary.BigEndian.Uint16(buf[:])) // The number of types, minus one.
	if typeCount == 0xFFFF {
		return nil, errInvalidDfont
	}
	typeCount += 1

	const tSize = 8
	if tSize*uint32(typeCount) > resourceMapLength-uint32(typeListOffset)-2 {
		return nil, errInvalidDfont
	}

	typeList := make([]byte, tSize*typeCount)
	_, err = r.ReadAt(typeList, int64(resourceMapOffset)+typeListOffset+2)
	if err != nil {
		return nil, err
	}
	numFonts, resourceListOffset := 0, 0
	for i := 0; i < typeCount; i++ {
		if binary.BigEndian.Uint32(typeList[tSize*i:]) != 0x73666e74 { // "sfnt".
			continue
		}

		numFonts = int(int16(binary.BigEndian.Uint16(typeList[tSize*i+4:])))
		if numFonts < 0 {
			return nil, errInvalidDfont
		}
		// https://github.com/kreativekorp/ksfl/wiki/Macintosh-Resource-File-Format
		// says that the value in the wire format is "the number of
		// resources of this type, minus one."
		numFonts++

		resourceListOffset = int(int16(binary.BigEndian.Uint16((typeList[tSize*i+6:]))))
		if resourceListOffset < 0 {
			return nil, errInvalidDfont
		}
	}
	if numFonts == 0 {
		return nil, errInvalidDfont
	}
	if numFonts > maxNumFonts {
		return nil, fmt.Errorf("number of fonts (%d) in collection exceed implementation limit (%d)",
			numFonts, maxNumFonts)
	}

	const rSize = 12
	o, n := uint32(int(typeListOffset)+resourceListOffset), rSize*uint32(numFonts)
	if o > resourceMapLength || n > resourceMapLength-o {
		return nil, errInvalidDfont
	}

	offsetsBytes := make([]byte, n)
	_, err = r.ReadAt(offsetsBytes, int64(resourceMapOffset+o))
	if err != nil {
		return nil, err
	}

	offsets := make([]uint32, numFonts)
	for i := range offsets {
		o := 0xffffff & binary.BigEndian.Uint32(offsetsBytes[rSize*i+4:])
		// Offsets are relative to the resource data start, not the file start.
		// A particular resource's data also starts with a 4-byte length, which
		// we skip.
		o += dfontResourceDataOffset + 4
		if o > maxTableOffset {
			return nil, errUnsupportedTableOffsetLength
		}
		offsets[i] = o
	}

	return offsets, nil
}
