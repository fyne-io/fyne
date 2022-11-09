package truetype

import (
	"encoding/binary"
	"errors"
)

var (
	errInvalidPostTable     = errors.New("invalid post table")
	errUnsupportedPostTable = errors.New("unsupported post table")
)

// TablePost represents an information stored in the PostScript font section.
type TablePost struct {
	// Names stores the glyph names. It may be nil.
	Names GlyphNames
	// ItalicAngle in counter-clockwise degrees from the vertical. Zero for
	// upright text, negative for text that leans to the right (forward).
	ItalicAngle float64
	// Version of the version tag of the "post" table.
	Version uint32
	// UnderlinePosition is the suggested distance of the top of the
	// underline from the baseline (negative values indicate below baseline).
	UnderlinePosition int16
	// Suggested values for the underline thickness.
	UnderlineThickness int16
	// IsFixedPitch indicates that the font is not proportionally spaced
	// (i.e. monospaced).
	IsFixedPitch bool
}

func parseTablePost(buf []byte, numGlyphs uint16) (TablePost, error) {
	// https://www.microsoft.com/typography/otspec/post.htm

	const headerSize = 32
	if len(buf) < headerSize {
		return TablePost{}, errInvalidPostTable
	}
	var (
		names GlyphNames
		err   error
	)
	u := binary.BigEndian.Uint32(buf)
	switch u {
	case 0x10000:
		names = postNamesFormat10{}
	case 0x30000:
		// No-op.
	case 0x20000:
		if len(buf) < headerSize+2+2*int(numGlyphs) {
			return TablePost{}, errInvalidPostTable
		}
		names, err = parseNameFormat20(buf, numGlyphs)
		if err != nil {
			return TablePost{}, err
		}
	default:
		return TablePost{}, errUnsupportedPostTable
	}

	ang := binary.BigEndian.Uint32(buf[4:])
	up := binary.BigEndian.Uint16(buf[8:])
	ut := binary.BigEndian.Uint16(buf[10:])
	fp := binary.BigEndian.Uint32(buf[12:])
	return TablePost{
		Version:            u,
		ItalicAngle:        float64(int32(ang)) / 0x10000,
		UnderlinePosition:  int16(up),
		UnderlineThickness: int16(ut),
		IsFixedPitch:       fp != 0,
		Names:              names,
	}, nil
}

// GlyphNames stores the names of a 'post' table.
type GlyphNames interface {
	// GlyphName return the postscript name of a
	// glyph, or an empty string if it not found
	GlyphName(x GID) string
}

type postNamesFormat10 struct{}

func (p postNamesFormat10) GlyphName(x GID) string {
	if int(x) >= numBuiltInPostNames {
		return ""
	}
	// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6post.html
	return builtInPostNames[x]
}

type postNamesFormat20 struct {
	glyphNameIndexes []uint16 // size numGlyph
	names            []string
}

func (p postNamesFormat20) GlyphName(x GID) string {
	if int(x) >= len(p.glyphNameIndexes) {
		return ""
	}
	u := int(p.glyphNameIndexes[x])
	if u < numBuiltInPostNames {
		return builtInPostNames[u]
	}
	u -= numBuiltInPostNames
	return p.names[u]
}

func parseNameFormat20(buf []byte, numGlyphs uint16) (postNamesFormat20, error) {
	// The wire format for a Version 2 post table is documented at:
	// https://www.microsoft.com/typography/otspec/post.htm
	const glyphNameIndexOffset = 34
	if len(buf) < glyphNameIndexOffset+2*int(numGlyphs) {
		return postNamesFormat20{}, errInvalidPostTable
	}
	buf = buf[glyphNameIndexOffset:]

	// we check at parse time that all the indexes are valid:
	// we find the maximum
	var maxIndex int
	glyphNameIndexes := make([]uint16, numGlyphs)
	for x := range glyphNameIndexes {
		u := binary.BigEndian.Uint16(buf[2*x:])
		// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6post.html
		// says that "32768 through 65535 are reserved for future use".
		if u > 32767 {
			return postNamesFormat20{}, errUnsupportedPostTable
		}
		if int(u) > maxIndex {
			maxIndex = int(u)
		}
		glyphNameIndexes[x] = u
	}

	// read all the string data until the end of the table
	var names []string
	for i := 2 * int(numGlyphs); i < len(buf); {
		length := int(buf[i])
		if len(buf) < i+1+length {
			return postNamesFormat20{}, errInvalidPostTable
		}
		names = append(names, string(buf[i+1:i+1+length]))
		i += int(length) + 1
	}
	if maxIndex >= numBuiltInPostNames && len(names) < (maxIndex-numBuiltInPostNames) {
		return postNamesFormat20{}, errInvalidPostTable
	}
	return postNamesFormat20{glyphNameIndexes: glyphNameIndexes, names: names}, nil
}
