package graphite

import (
	"errors"
	"fmt"
	"math/bits"

	"github.com/benoitkugler/textlayout/fonts/binaryreader"
)

type subboxMetrics struct {
	Left       uint8 // Left of subbox
	Right      uint8 // Right of subbox
	Bottom     uint8 // Bottom of subbox
	Top        uint8 // Top of subbox
	DiagNegMin uint8 // Defines minimum negatively-sloped diagonal
	DiagNegMax uint8 // Defines maximum negatively -sloped diagonal
	DiagPosMin uint8 // Defines minimum positively-sloped diagonal
	DiagPosMax uint8 // Defines maximum positively-sloped diagonal
}

type octaboxMetrics struct {
	subBbox    []subboxMetrics
	bitmap     uint16
	diagNegMin uint8 // Defines minimum negatively-sloped diagonal
	diagNegMax uint8 // Defines maximum negatively -sloped diagonal
	diagPosMin uint8 // Defines minimum positively-sloped diagonal
	diagPosMax uint8 // Defines maximum positively-sloped diagonal
}

func scaleBox(b rect, zxmin, zymin, zxmax, zymax uint8) rect {
	scaleTo := func(t uint8, zmin, zmax float32) float32 {
		return (zmin + float32(t)*(zmax-zmin)/255)
	}
	return rect{
		Position{scaleTo(zxmin, b.bl.X, b.tr.X), scaleTo(zymin, b.bl.Y, b.tr.Y)},
		Position{scaleTo(zxmax, b.bl.X, b.tr.X), scaleTo(zymax, b.bl.Y, b.tr.Y)},
	}
}

func (oc octaboxMetrics) computeBoxes(bbox rect) glyphBoxes {
	diamax := rect{
		Position{bbox.bl.X + bbox.bl.Y, bbox.bl.X - bbox.tr.Y},
		Position{bbox.tr.X + bbox.tr.Y, bbox.tr.X - bbox.bl.Y},
	}
	diabound := scaleBox(diamax, oc.diagNegMin, oc.diagPosMin, oc.diagNegMax, oc.diagPosMax)

	out := glyphBoxes{bitmap: oc.bitmap, slant: diabound}

	out.subBboxes = make([]rect, len(oc.subBbox))
	out.slantSubBboxes = make([]rect, len(oc.subBbox))
	for i, subbox := range oc.subBbox {
		out.subBboxes[i] = scaleBox(bbox, subbox.Left, subbox.Bottom, subbox.Right, subbox.Top)
		out.slantSubBboxes[i] = scaleBox(diamax, subbox.DiagNegMin, subbox.DiagPosMin, subbox.DiagNegMax, subbox.DiagPosMax)
	}

	return out
}

type glyphAttributes struct {
	octaboxMetrics *octaboxMetrics // may be nil
	attributes     attributeSet
}

type attributSetEntry struct {
	attributes []int16
	firstKey   uint16
}

// sorted by `firstKey`; firstKey + len(attributes) <= nextFirstKey
type attributeSet []attributSetEntry

func (as attributeSet) get(key uint16) int16 {
	// binary search
	for i, j := 0, len(as); i < j; {
		h := i + (j-i)/2
		entry := as[h]
		if key < entry.firstKey {
			j = h
		} else if entry.firstKey+uint16(len(entry.attributes))-1 < key {
			i = h + 1
		} else {
			return entry.attributes[key-entry.firstKey]
		}
	}
	return 0
}

type tableGlat []glyphAttributes // with len >= numGlyphs

func parseTableGloc(data []byte, numGlyphs int) ([]uint32, uint16, error) {
	r := binaryreader.NewReader(data)
	if len(data) < 8 {
		return nil, 0, errors.New("invalid Gloc table (EOF)")
	}
	_, _ = r.Uint32()
	flags, _ := r.Uint16()
	numAttributes, _ := r.Uint16()
	isLong := flags&1 != 0

	// the number of locations may be greater than numGlyphs,
	// since there may be pseudo-glyphs
	// compute it from the end of the table:
	byteLength := len(data) - (int(numAttributes) * int(flags&2)) - 8
	numLocations := byteLength / 2
	if isLong {
		numLocations = byteLength / 4
	}

	if numLocations < numGlyphs+1 {
		return nil, 0, fmt.Errorf("invalid Gloc table: %d locations for %d glyphs ", numLocations, numGlyphs)
	}

	var locations []uint32
	if isLong {
		var err error
		if err != nil {
			return nil, 0, fmt.Errorf("invalid Gloc table: %s", err)
		}
		locations, err = r.Uint32s(numLocations)
		if err != nil {
			return nil, 0, err
		}
	} else {
		tmp, err := r.Uint16s(numLocations)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid Gloc table: %s", err)
		}
		locations = make([]uint32, len(tmp))
		for i, o := range tmp {
			locations[i] = uint32(o)
		}
	}

	return locations, numAttributes, nil
}

// locations has length numAttrs + 1
func parseTableGlat(data []byte, locations []uint32) (tableGlat, error) {
	data, version, err := decompressTable(data)
	if err != nil {
		return nil, fmt.Errorf("invalid table Glat: %s", err)
	}

	out := make(tableGlat, len(locations)-1)
	for i := range out {
		start, end := locations[i], locations[i+1]
		if start >= end {
			continue
		}
		if len(data) < int(end) {
			return nil, fmt.Errorf("invalid offset for table Glat: %d < %d", len(data), end)
		}
		glyphData := data[start:end]
		out[i], err = parseOneGlyphAttr(glyphData, version)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func parseOneGlyphAttr(data []byte, version uint16) (out glyphAttributes, err error) {
	r := binaryreader.NewReader(data)
	if version >= 3 { // read the octabox metrics
		out.octaboxMetrics = new(octaboxMetrics)
		out.octaboxMetrics.bitmap, err = r.Uint16()
		if err != nil {
			return out, fmt.Errorf("invalid Glat entry: %s", err)
		}
		nbSubboxes := bits.OnesCount16(out.octaboxMetrics.bitmap)
		metricsLength := 4 + 8*nbSubboxes
		if len(r.Data()) < metricsLength {
			return out, errors.New("invalid Glat entry (EOF)")
		}
		// NOTE: we follow the ordering of the fields from fonttools, not
		// from the graphite spec
		out.octaboxMetrics.diagNegMin, _ = r.Byte()
		out.octaboxMetrics.diagNegMax, _ = r.Byte()
		out.octaboxMetrics.diagPosMin, _ = r.Byte()
		out.octaboxMetrics.diagPosMax, _ = r.Byte()
		out.octaboxMetrics.subBbox = make([]subboxMetrics, nbSubboxes)
		_ = r.ReadStruct(out.octaboxMetrics.subBbox)
	}
	lastEndKey := -1
	if version < 2 { // one byte
		for len(r.Data()) >= 2 {
			attNum, _ := r.Byte()
			num, _ := r.Byte()
			attributes, err := r.Int16s(int(num))
			if err != nil {
				return out, fmt.Errorf("invalid Glat entry attributes: %s", err)
			}

			if int(attNum) < lastEndKey {
				return out, fmt.Errorf("invalid Glat entry attribute key: %d", attNum)
			}
			lastEndKey = int(attNum) + len(attributes)

			out.attributes = append(out.attributes, attributSetEntry{
				firstKey:   uint16(attNum),
				attributes: attributes,
			})
		}
	} else { // same with two bytes header fields
		for len(r.Data()) >= 4 {
			attNum, _ := r.Uint16()
			num, _ := r.Uint16()
			attributes, err := r.Int16s(int(num))
			if err != nil {
				return out, fmt.Errorf("invalid Glat entry attributes: %s", err)
			}

			if int(attNum) < lastEndKey {
				return out, fmt.Errorf("invalid Glat entry attribute key: %d < %d", attNum, lastEndKey)
			}
			lastEndKey = int(attNum) + len(attributes)

			out.attributes = append(out.attributes, attributSetEntry{
				firstKey:   attNum,
				attributes: attributes,
			})
		}
	}

	return out, nil
}
