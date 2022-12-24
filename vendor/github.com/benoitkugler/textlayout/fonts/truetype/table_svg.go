package truetype

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
	"sort"

	"github.com/benoitkugler/textlayout/fonts"
	"github.com/benoitkugler/textlayout/fonts/binaryreader"
)

var tagSVG = NewTag('S', 'V', 'G', ' ')

// sorted by startGlyphID
type tableSVG []svgDocumentIndexEntry

func (s tableSVG) glyphData(gid GID) (fonts.GlyphSVG, bool) {
	data, ok := s.rawGlyphData(gid)
	if !ok {
		return fonts.GlyphSVG{}, false
	}

	// un-compress if needed
	if r, err := gzip.NewReader(bytes.NewReader(data)); err == nil {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err == nil {
			data = buf.Bytes()
		}
	}

	return fonts.GlyphSVG{Source: data}, true
}

func (s tableSVG) rawGlyphData(gid GID) ([]byte, bool) {
	// binary search
	for i, j := 0, len(s); i < j; {
		h := i + (j-i)/2
		entry := s[h]
		if gid < GID(entry.first) {
			j = h
		} else if GID(entry.last) < gid {
			i = h + 1
		} else {
			return entry.svg, true
		}
	}
	return nil, false
}

type svgDocumentIndexEntry struct {
	// svg document
	// each glyph description must be written
	// in an element with id=glyphXXX
	svg   []byte
	first gid // The first glyph ID in the range described by this index entry.
	last  gid // The last glyph ID in the range described by this index entry. Must be >= startGlyphID.
}

func parseTableSVG(buf []byte) (tableSVG, error) {
	if len(buf) < 6 {
		return nil, errors.New("invalid SVG table (EOF)")
	}

	// version := binary.BigEndian.Uint16(buf)
	offset := binary.BigEndian.Uint32(buf[2:])
	if len(buf) < int(offset) {
		return nil, errors.New("invalid SVG table (invalid offset)")
	}
	buf = buf[offset:]
	r := binaryreader.NewReader(buf)
	numEntries, err := r.Uint16()
	if err != nil {
		return nil, err
	}
	type docHeader struct {
		StartGlyphID, EndGlyphID gid
		Offset                   uint32
		Length                   uint32
	}

	headers := make([]docHeader, numEntries)
	err = r.ReadStruct(&headers)
	if err != nil {
		return nil, err
	}
	out := make(tableSVG, numEntries)
	for i, header := range headers {
		endOffset := header.Offset + header.Length
		if int(endOffset) > len(buf) {
			return nil, errors.New("invalid SVG table (invalid offset)")
		}
		out[i].first = header.StartGlyphID
		out[i].last = header.EndGlyphID
		out[i].svg = buf[header.Offset:endOffset]
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].first < out[j].first
	})

	return out, nil
}
