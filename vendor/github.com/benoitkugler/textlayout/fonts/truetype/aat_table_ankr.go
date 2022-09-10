package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type TableAnkr struct {
	class   Class
	anchors []byte
	flags   uint16
}

// AATAnchor is a point within the coordinate space of a given glyph
// independent of the control points used to render the glyph
type AATAnchor struct {
	X, Y int16
}

// GetAnchor return the i-th anchor for `glyph`, or {0,0} if not found.
func (t TableAnkr) GetAnchor(glyph GID, index int) (anchor AATAnchor) {
	offset, ok := t.class.ClassID(glyph)
	if !ok {
		return anchor
	}

	count := int(binary.BigEndian.Uint32(t.anchors[offset:])) // access sanitized during parsing
	if index >= count {
		return anchor // invalid index
	}

	indexStart := int(offset) + 4 + 4*index
	if len(t.anchors) < indexStart+4 {
		return anchor // invalid table
	}
	anchor.X = int16(binary.BigEndian.Uint16(t.anchors[indexStart:]))
	anchor.Y = int16(binary.BigEndian.Uint16(t.anchors[indexStart+2:]))
	return anchor
}

func parseTableAnkr(data []byte, numGlyphs int) (out TableAnkr, err error) {
	if len(data) < 12 {
		return out, errors.New("invalid 'ankr' table (EOF)")
	}
	out.flags = binary.BigEndian.Uint16(data[2:])
	lookupOffset := binary.BigEndian.Uint32(data[4:])
	glyphOffset := binary.BigEndian.Uint32(data[8:])
	if lookupOffset > glyphOffset || len(data) < int(glyphOffset) {
		return out, errors.New("invalid 'ankr' table")
	}
	out.class, err = parseAATLookupTable(data, lookupOffset, numGlyphs, false)
	if err != nil {
		return out, fmt.Errorf("invalid 'ankr' table: %s", err)
	}
	out.anchors = data[glyphOffset:]
	if e := out.class.Extent(); e+4 > len(out.anchors) {
		return out, errors.New("invalid 'ankr' table (EOF)")
	}
	return out, nil
}
