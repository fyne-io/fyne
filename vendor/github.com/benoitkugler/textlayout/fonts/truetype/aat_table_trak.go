package truetype

import (
	"encoding/binary"
	"errors"
)

type TableTrak struct {
	Horizontal, Vertical TrakData // may be empty
}

// IsEmpty return `true` it the table has no entries.
func (t TableTrak) IsEmpty() bool {
	return len(t.Horizontal.Entries)+len(t.Vertical.Entries) == 0
}

func parseTrakTable(data []byte) (out TableTrak, err error) {
	if len(data) < 12 {
		return out, errors.New("invalid trak table (EOF)")
	}
	// ignoring version and format
	horizOffset := binary.BigEndian.Uint16(data[6:])
	vertOffset := binary.BigEndian.Uint16(data[8:])

	if horizOffset != 0 {
		out.Horizontal, err = parseTrakData(data, int(horizOffset))
		if err != nil {
			return out, err
		}
	}
	if vertOffset != 0 {
		out.Vertical, err = parseTrakData(data, int(vertOffset))
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

type TrackEntry struct {
	PerSizeTracking []int16 // in font units, with length len(Sizes)
	Track           float32
	NameIndex       NameID
}

type TrakData struct {
	Entries []TrackEntry
	Sizes   []float32
}

// idx is assumed to verify idx <= len(Sizes) - 2
func (td TrakData) interpolateAt(idx int, targetSize float32, trackSizes []int16) float32 {
	s0 := td.Sizes[idx]
	s1 := td.Sizes[idx+1]
	var t float32
	if s0 != s1 {
		t = (targetSize - s0) / (s1 - s0)
	}
	return t*float32(trackSizes[idx+1]) + (1.-t)*float32(trackSizes[idx])
}

// GetTracking select the tracking for the given `trackValue` and apply it
// for `ptem`. It returns 0 if not found.
func (td TrakData) GetTracking(ptem float32, trackValue float32) float32 {
	// Choose track.

	var trackTableEntry *TrackEntry
	for i := range td.Entries {
		/* Note: Seems like the track entries are sorted by values.  But the
		 * spec doesn't explicitly say that.  It just mentions it in the example. */

		if td.Entries[i].Track == trackValue {
			trackTableEntry = &td.Entries[i]
			break
		}
	}
	if trackTableEntry == nil {
		return 0.
	}

	// Choose size.

	if len(td.Sizes) == 0 {
		return 0.
	}
	if len(td.Sizes) == 1 {
		return float32(trackTableEntry.PerSizeTracking[0])
	}

	var sizeIndex int
	for sizeIndex = range td.Sizes {
		if td.Sizes[sizeIndex] >= ptem {
			break
		}
	}
	if sizeIndex != 0 {
		sizeIndex = sizeIndex - 1
	}
	return td.interpolateAt(sizeIndex, ptem, trackTableEntry.PerSizeTracking)
}

func parseTrakData(data []byte, offset int) (out TrakData, err error) {
	if len(data) < int(offset)+8 {
		return out, errors.New("invalid trak data table (EOF)")
	}
	nTracks := binary.BigEndian.Uint16(data[offset:])
	nSizes := binary.BigEndian.Uint16(data[offset+2:])
	sizeTableOffset := int(binary.BigEndian.Uint32(data[offset+4:]))

	if len(data) < offset+8+8*int(nTracks) {
		return out, errors.New("invalid trak data table (EOF)")
	}

	out.Entries = make([]TrackEntry, nTracks)
	for i := range out.Entries {
		out.Entries[i].Track = fixed1616ToFloat(binary.BigEndian.Uint32(data[offset+8+i*8:]))
		out.Entries[i].NameIndex = NameID(binary.BigEndian.Uint16(data[offset+8+i*8+4:]))
		sizeOffset := binary.BigEndian.Uint16(data[offset+8+i*8+6:])
		out.Entries[i].PerSizeTracking, err = parseTrackSizes(data, int(sizeOffset), nSizes)
		if err != nil {
			return out, err
		}
	}

	if len(data) < sizeTableOffset+4*int(nSizes) {
		return out, errors.New("invalid trak data table (EOF)")
	}
	out.Sizes = make([]float32, nSizes)
	for i := range out.Sizes {
		out.Sizes[i] = fixed1616ToFloat(binary.BigEndian.Uint32(data[sizeTableOffset+4*i:]))
	}

	return out, nil
}

func parseTrackSizes(data []byte, offset int, count uint16) ([]int16, error) {
	if len(data) < offset+int(count)*2 {
		return nil, errors.New("invalid trak table per-sizes values (EOF)")
	}
	out := make([]int16, count)
	for i := range out {
		out[i] = int16(binary.BigEndian.Uint16(data[offset+2*i:]))
	}
	return out, nil
}
