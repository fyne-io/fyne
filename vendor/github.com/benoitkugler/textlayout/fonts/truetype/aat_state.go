package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// state tables: https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html#StateTables

// AATStateTable is an extended state table.
type AATStateTable struct {
	class    Class
	entries  []AATStateEntry
	states   [][]uint16 // _ rows, nClasses columns
	nClasses uint32     //  for some reasons, this may differ from Class.Extent
}

const (
	aatStateHeaderSize    = 8
	aatExtStateHeaderSize = 16
)

// extended is true for morx/kerx, false for kern
// every `Data` field of the entries will be of length `entryDataSize`
// meaning they can later be safely interpreted
func parseStateTable(data []byte, entryDataSize int, extended bool, numGlyphs int) (out AATStateTable, err error) {
	headerSize := aatStateHeaderSize
	if extended {
		headerSize = aatExtStateHeaderSize
	}
	if len(data) < headerSize {
		return out, errors.New("invalid AAT state table (EOF)")
	}
	var stateOffset, entryOffset uint32
	if extended {
		out.nClasses = binary.BigEndian.Uint32(data)
		classOffset := binary.BigEndian.Uint32(data[4:])
		stateOffset = binary.BigEndian.Uint32(data[8:])
		entryOffset = binary.BigEndian.Uint32(data[12:])

		out.class, err = parseAATLookupTable(data, classOffset, numGlyphs, false)
	} else {
		out.nClasses = uint32(binary.BigEndian.Uint16(data))
		classOffset := binary.BigEndian.Uint16(data[2:])
		stateOffset = uint32(binary.BigEndian.Uint16(data[4:]))
		entryOffset = uint32(binary.BigEndian.Uint16(data[6:]))

		if len(data) < int(classOffset) {
			return out, errors.New("invalid AAT state table (EOF)")
		}
		// no class format here
		out.class, err = parseClassFormat1(data[classOffset:], 1)
	}
	if err != nil {
		return out, fmt.Errorf("invalid AAT state table: %s", err)
	}
	nC := int(out.nClasses)
	// Ensure pre-defined classes fit.
	if nC < 4 {
		return out, fmt.Errorf("invalid number of classes in state table: %d", nC)
	}

	if stateOffset > entryOffset || len(data) < int(entryOffset) {
		return out, errors.New("invalid AAT state table (EOF)")
	}

	var states []uint16
	if extended {
		states, err = parseUint16s(data[stateOffset:entryOffset], int(entryOffset-stateOffset)/2)
		if err != nil {
			return out, err
		}
	} else {
		states = make([]uint16, entryOffset-stateOffset)
		for i, b := range data[stateOffset:entryOffset] {
			states[i] = uint16(b)
		}
	}

	out.states = make([][]uint16, len(states)/nC)
	for i := range out.states {
		out.states[i] = states[i*nC : (i+1)*nC]
	}

	// find max index
	var maxi uint16
	for _, stateIndex := range states {
		if stateIndex > maxi {
			maxi = stateIndex
		}
	}

	out.entries, err = parseStateEntries(data[entryOffset:], int(maxi)+1, entryDataSize, stateOffset, out.nClasses, extended)

	return out, err
}

// GetClass return the class for the given glyph, with the correct default value.
func (t *AATStateTable) GetClass(glyph GID) uint16 {
	if glyph == 0xFFFF { // deleted glyph
		return 2 // class deleted
	}
	c, ok := t.class.ClassID(glyph)
	if !ok {
		return 1 // class out of bounds
	}
	return uint16(c) // class for a state table can't be uint32
}

// GetEntry return the entry for the given state and class,
// and handle invalid values (by returning an empty entry).
func (t AATStateTable) GetEntry(state, class uint16) AATStateEntry {
	if uint32(class) >= t.nClasses {
		class = 1 // class out of bounds
	}
	if int(state) >= len(t.states) {
		return AATStateEntry{}
	}
	entry := t.states[state][class]
	return t.entries[entry]
}

// AATStateEntry is the shared type for entries
// in a state table. See the various AsXXX methods
// to exploit its content.
type AATStateEntry struct {
	NewState uint16
	Flags    uint16 // Table specific.
	// Remaining of the entry, context specific
	data [4]byte
}

// data is at the start of the entries array
// assume extraDataSize <= 4
func parseStateEntries(data []byte, count, extraDataSize int, stateTableOffset, nClasses uint32, extended bool) ([]AATStateEntry, error) {
	entrySize := 4 + extraDataSize
	if len(data) < count*entrySize {
		return nil, errors.New("invalid AAT state entry array (EOF)")
	}
	out := make([]AATStateEntry, count)
	for i := range out {
		newState := binary.BigEndian.Uint16(data[i*entrySize:])
		if extended { // newState is directly the index
			out[i].NewState = newState
		} else { // newState is an offset: convert back to index
			out[i].NewState = uint16((int(newState) - int(stateTableOffset)) / int(nClasses))
		}
		out[i].Flags = binary.BigEndian.Uint16(data[i*entrySize+2:])
		copy(out[i].data[:], data[i*entrySize+4:(i+1)*entrySize])
	}
	return out, nil
}

// AsMorxContextual reads the internal data for entries in morx contextual subtable.
// The returned indexes use 0xFFFF as empty value.
func (e AATStateEntry) AsMorxContextual() (markIndex, currentIndex uint16) {
	markIndex = binary.BigEndian.Uint16(e.data[:])
	currentIndex = binary.BigEndian.Uint16(e.data[2:])
	return
}

// AsMorxInsertion reads the internal data for entries in morx insertion subtable.
// The returned indexes use 0xFFFF as empty value.
func (e AATStateEntry) AsMorxInsertion() (currentIndex, markedIndex uint16) {
	currentIndex = binary.BigEndian.Uint16(e.data[:])
	markedIndex = binary.BigEndian.Uint16(e.data[2:])
	return
}

// AsMorxLigature reads the internal data for entries in morx ligature subtable.
func (e AATStateEntry) AsMorxLigature() (ligActionIndex uint16) {
	return binary.BigEndian.Uint16(e.data[:])
}

// AsKernxIndex reads the internal data for entries in 'kern/x' subtable format 1.
// An entry with no index returns 0xFFFF
func (e AATStateEntry) AsKernxIndex() uint16 {
	// for kern table, during parsing, we store the resolved index
	// at the same place as kerx tables
	return binary.BigEndian.Uint16(e.data[:])
}

// AAT lookup implementing Class

type lookupFormat0 []uint32

func (l lookupFormat0) ClassID(gid GID) (uint32, bool) {
	if int(gid) >= len(l) {
		return 0, false
	}
	return l[gid], true
}

func (l lookupFormat0) GlyphSize() int { return len(l) }

func (l lookupFormat0) Extent() int {
	max := uint32(0)
	for _, r := range l {
		if r >= max {
			max = r
		}
	}
	return int(max) + 1
}

func parseAATLookupFormat0(data []byte, numGlyphs int, isLong bool) (lookupFormat0, error) {
	if isLong {
		if len(data) < 2+numGlyphs*4 {
			return nil, errors.New("invalid AAT lookup format 0 (EOF)")
		}
		return parseUint32s(data[2:], numGlyphs), nil
	}
	uint16s, err := parseUint16s(data[2:], numGlyphs)
	if err != nil {
		return nil, err
	}
	out := make(lookupFormat0, numGlyphs)
	for i, v := range uint16s {
		out[i] = uint32(v)
	}
	return out, nil
}

// lookupFormat2 is the same as classFormat2, but with start and end are reversed in the binary
func parseAATLookupFormat2(data []byte, isLong bool) (classFormat2, error) {
	const headerSize = 12 // including classFormat
	if len(data) < headerSize {
		return nil, errors.New("invalid AAT lookup format 2 (EOF)")
	}

	unitSize := binary.BigEndian.Uint16(data[2:])
	num := int(binary.BigEndian.Uint16(data[4:]))
	// 3 other field ignored
	if unitSize != 6 {
		return nil, fmt.Errorf("unexpected AAT lookup segment size: %d", unitSize)
	}

	recordSize := 6
	if isLong {
		recordSize = 8
	}
	if len(data) < headerSize+num*recordSize {
		return nil, errors.New("invalid AAT lookup format 2 (EOF)")
	}

	out := make(classFormat2, num)
	for i := range out {
		out[i].end = gid(binary.BigEndian.Uint16(data[headerSize+i*recordSize:]))
		out[i].start = gid(binary.BigEndian.Uint16(data[headerSize+i*recordSize+2:]))
		if isLong {
			out[i].targetClassID = binary.BigEndian.Uint32(data[headerSize+i*recordSize+4:])
		} else {
			out[i].targetClassID = uint32(binary.BigEndian.Uint16(data[headerSize+i*recordSize+4:]))
		}
	}
	return out, nil
}

// sorted accordins to `last`
type lookupFormat4 []struct {
	values      []uint32 // length last - first + 1
	first, last GID
}

func (l lookupFormat4) ClassID(gid GID) (uint32, bool) {
	// binary search
	for i, j := 0, len(l); i < j; {
		h := i + (j-i)/2
		entry := l[h]
		if gid < entry.first {
			j = h
		} else if entry.last < gid {
			i = h + 1
		} else {
			return entry.values[gid-entry.first], true
		}
	}
	return 0, false
}

func (l lookupFormat4) GlyphSize() int {
	var out int
	for _, rec := range l {
		out += len(rec.values)
	}
	return out
}

func (l lookupFormat4) Extent() int {
	max := uint32(0)
	for _, rec := range l {
		for _, val := range rec.values {
			if val >= max {
				max = val
			}
		}
	}
	return int(max) + 1
}

func parseAATLookupFormat4(data []byte, isLong bool) (lookupFormat4, error) {
	const headerSize = 12 // including classFormat
	if len(data) < headerSize {
		return nil, errors.New("invalid AAT lookup format 4 (EOF)")
	}
	unitSize := int(binary.BigEndian.Uint16(data[2:]))
	num := int(binary.BigEndian.Uint16(data[4:]))
	// 3 other field ignored
	if unitSize < 6 {
		return nil, fmt.Errorf("unexpected AAT lookup segment size: %d", unitSize)
	}

	if len(data) < headerSize+num*unitSize {
		return nil, errors.New("invalid AAT lookup format 4 (EOF)")
	}

	valSize := 2
	if isLong {
		valSize = 4
	}

	// we do not include the termination segment
	out := make(lookupFormat4, num-1)
	for i := range out {
		out[i].last = GID(binary.BigEndian.Uint16(data[headerSize+i*unitSize:]))
		out[i].first = GID(binary.BigEndian.Uint16(data[headerSize+i*unitSize+2:]))
		// if out[i].last == 0xffff {
		// 	continue
		// }
		offset := int(binary.BigEndian.Uint16(data[headerSize+i*unitSize+4:]))
		if out[i].last < out[i].first {
			return nil, fmt.Errorf("invalid AAT lookup format 4 (first, last : %d, %d)", out[i].first, out[i].last)
		}
		count := int(out[i].last) - int(out[i].first) + 1

		if len(data) < offset+count*valSize {
			return nil, fmt.Errorf("invalid AAT lookup format 4 (offset %d and count %d for length %d)", offset, count, len(data))
		}

		if isLong {
			out[i].values = parseUint32s(data[offset:], count)
		} else {
			tmp, _ := parseUint16s(data[offset:], count) // length already checked
			out[i].values = make([]uint32, count)
			for j, v := range tmp {
				out[i].values[j] = uint32(v)
			}
		}
	}
	return out, nil
}

// sorted pairs of GlyphIndex, value
type lookupFormat6 []struct {
	gid   GID
	value uint32
}

func (l lookupFormat6) ClassID(gid GID) (uint32, bool) {
	// binary search
	for i, j := 0, len(l); i < j; {
		h := i + (j-i)/2
		entry := l[h]
		if gid < entry.gid {
			j = h
		} else if entry.gid < gid {
			i = h + 1
		} else {
			return entry.value, true
		}
	}
	return 0, false
}

func (l lookupFormat6) GlyphSize() int { return len(l) }

func (l lookupFormat6) Extent() int {
	max := uint32(0)
	for _, r := range l {
		if r.value >= max {
			max = r.value
		}
	}
	return int(max) + 1
}

func parseAATLookupFormat6(data []byte, isLong bool) (lookupFormat6, error) {
	const headerSize = 12 // including classFormat
	if len(data) < headerSize {
		return nil, errors.New("invalid AAT lookup format 6 (EOF)")
	}

	unitSize := int(binary.BigEndian.Uint16(data[2:]))
	num := int(binary.BigEndian.Uint16(data[4:]))
	// 3 other field ignored
	if isLong && unitSize < 6 || unitSize < 4 {
		return nil, fmt.Errorf("unexpected AAT lookup segment size: %d", unitSize)
	}

	if len(data) < headerSize+num*unitSize {
		return nil, errors.New("invalid AAT lookup format 6 (EOF)")
	}

	out := make(lookupFormat6, num)
	for i := range out {
		out[i].gid = GID(binary.BigEndian.Uint16(data[headerSize+i*unitSize:]))
		if isLong {
			out[i].value = binary.BigEndian.Uint32(data[headerSize+i*unitSize+2:])
		} else {
			out[i].value = uint32(binary.BigEndian.Uint16(data[headerSize+i*unitSize+2:]))
		}
	}
	return out, nil
}

// lookupFormat8 is the same as ClassFormat1
func parseAATLookupFormat8(data []byte) (classFormat1, error) {
	return parseClassFormat1(data[2:], 2)
}

// lookupFormat10 is the same as ClassFormat1, with 2 or 4 byte size
func parseAATLookupFormat10(data []byte, isLong bool) (classFormat1, error) {
	byteSize := 2
	if isLong {
		byteSize = 4
	}
	// skip unit size
	if len(data) < 4 {
		return classFormat1{}, errors.New("invalid AAT lookup format 10 (EOF)")
	}
	return parseClassFormat1(data[4:], byteSize)
}

// numGlyphs is used for unbounded lookups
// if isLong is true, the class values are uint32
func parseAATLookupTable(data []byte, offset uint32, numGlyphs int, isLong bool) (Class, error) {
	if len(data) < int(offset)+2 {
		return nil, errors.New("invalid AAT lookup table (EOF)")
	}
	data = data[offset:]
	switch format := binary.BigEndian.Uint16(data); format {
	case 0:
		return parseAATLookupFormat0(data, numGlyphs, isLong)
	case 2:
		return parseAATLookupFormat2(data, isLong)
	case 4:
		return parseAATLookupFormat4(data, isLong)
	case 6:
		return parseAATLookupFormat6(data, isLong)
	case 8:
		return parseAATLookupFormat8(data)
	case 10:
		return parseAATLookupFormat10(data, isLong)
	default:
		return nil, fmt.Errorf("invalid AAT lookup table kind : %d", format)
	}
}
