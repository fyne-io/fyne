// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"fmt"
)

// AAT layout

// State table header, without the actual data
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html
type AATStateTable struct {
	StateSize  uint16          // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	ClassTable ClassTable      `offsetSize:"Offset16"` // Byte offset from the beginning of the state table to the class subtable.
	stateArray Offset16        // Byte offset from the beginning of the state table to the state array.
	entryTable Offset16        // Byte offset from the beginning of the state table to the entry subtable.
	States     [][]uint8       `isOpaque:""`
	Entries    []AATStateEntry `isOpaque:""` // entry data are empty
}

func (state *AATStateTable) parseStates(src []byte) error {
	if state.stateArray > state.entryTable {
		return fmt.Errorf("invalid AAT state offsets (%d > %d)", state.stateArray, state.entryTable)
	}
	if L := len(src); L < int(state.entryTable) {
		return fmt.Errorf("EOF: expected length: %d, got %d", state.entryTable, L)
	}
	states := src[state.stateArray:state.entryTable]

	nC := int(state.StateSize)
	// Ensure pre-defined classes fit.
	if nC < 4 {
		return fmt.Errorf("invalid number of classes in AAT state table: %d", nC)
	}
	state.States = make([][]uint8, len(states)/nC)
	for i := range state.States {
		state.States[i] = states[i*nC : (i+1)*nC]
	}

	return nil
}

func (state *AATStateTable) parseEntries(src []byte) (int, error) {
	// find max index
	var maxi uint8
	for _, l := range state.States {
		for _, stateIndex := range l {
			if stateIndex > maxi {
				maxi = stateIndex
			}
		}
	}

	src = src[state.entryTable:] // checked in parseStates
	count := int(maxi) + 1
	var err error
	state.Entries, err = parseAATStateEntries(src, count, 0)
	if err != nil {
		return 0, err
	}

	// newState is an offset: convert back to index
	for i, entry := range state.Entries {
		state.Entries[i].NewState = uint16((int(entry.NewState) - int(state.stateArray)) / int(state.StateSize))
	}

	// the own header data stop at the entryTable offset
	return 8, err
}

// src starts at the entryTable
func parseAATStateEntries(src []byte, count, entryDataSize int) ([]AATStateEntry, error) {
	entrySize := 4 + entryDataSize
	if L := len(src); L < count*entrySize {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", count*entrySize, L)
	}
	out := make([]AATStateEntry, count)
	for i := range out {
		out[i].NewState = binary.BigEndian.Uint16(src[i*entrySize:])
		out[i].Flags = binary.BigEndian.Uint16(src[i*entrySize+2:])
		copy(out[i].data[:], src[i*entrySize+4:(i+1)*entrySize])
	}

	return out, nil
}

// ClassTable is the same as AATLookup8, but with no format and with bytes instead of uint16s
type ClassTable struct {
	StartGlyph GlyphID
	Values     []byte `arrayCount:"FirstUint16"`
}

// Extended state table, including the data
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html - State tables
// binarygen: argument=entryDataSize int
type AATStateTableExt struct {
	StateSize  uint32          // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	Class      AATLookup       `offsetSize:"Offset32"` // Byte offset from the beginning of the state table to the class subtable.
	stateArray Offset32        // Byte offset from the beginning of the state table to the state array.
	entryTable Offset32        // Byte offset from the beginning of the state table to the entry subtable.
	States     [][]uint16      `isOpaque:""` // each sub array has length stateSize
	Entries    []AATStateEntry `isOpaque:""` // length is the maximum state + 1
}

func (state *AATStateTableExt) parseStates(src []byte, _, _ int) error {
	if state.stateArray > state.entryTable {
		return fmt.Errorf("invalid AAT state offsets (%d > %d)", state.stateArray, state.entryTable)
	}
	if L := len(src); L < int(state.entryTable) {
		return fmt.Errorf("EOF: expected length: %d, got %d", state.entryTable, L)
	}

	statesArray := src[state.stateArray:state.entryTable]
	states, err := ParseUint16s(statesArray, len(statesArray)/2)
	if err != nil {
		return err
	}

	nC := int(state.StateSize)
	// Ensure pre-defined classes fit.
	if nC < 4 {
		return fmt.Errorf("invalid number of classes in AAT state table: %d", nC)
	}
	state.States = make([][]uint16, len(states)/nC)
	for i := range state.States {
		state.States[i] = states[i*nC : (i+1)*nC]
	}
	return nil
}

func (state *AATStateTableExt) parseEntries(src []byte, _, entryDataSize int) (int, error) {
	// find max index
	var maxi uint16
	for _, l := range state.States {
		for _, stateIndex := range l {
			if stateIndex > maxi {
				maxi = stateIndex
			}
		}
	}

	src = src[state.entryTable:] // checked in parseStates
	count := int(maxi) + 1
	var err error
	state.Entries, err = parseAATStateEntries(src, count, entryDataSize)

	// the own header data stop at the entryTable offset
	return 16, err
}

// AATStateEntry is shared between old and extended state tables,
// and between the different kind of entries.
// See the various AsXXX() methods.
type AATStateEntry struct {
	NewState uint16
	Flags    uint16  // Table specific.
	data     [4]byte // Table specific.
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

// AsKernxIndex reads the internal data for entries in 'kern/x' subtable format 1 or 4.
// An entry with no index returns 0xFFFF
func (e AATStateEntry) AsKernxIndex() uint16 {
	// for kern table, during parsing, we store the resolved index
	// at the same place as kerx tables
	return binary.BigEndian.Uint16(e.data[:])
}

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 // The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 // The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 // The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

// AATLookup is conceptually a map[GlyphID]uint16, but it may
// be implemented more efficiently.
type AATLookup interface {
	AatLookupMixed

	isAATLookup()

	// Class returns the class ID for the provided glyph, or (0, false)
	// for glyphs not covered by this class.
	Class(g GlyphID) (uint16, bool)
}

func (AATLoopkup0) isAATLookup()  {}
func (AATLoopkup2) isAATLookup()  {}
func (AATLoopkup4) isAATLookup()  {}
func (AATLoopkup6) isAATLookup()  {}
func (AATLoopkup8) isAATLookup()  {}
func (AATLoopkup10) isAATLookup() {}

type AATLoopkup0 struct {
	version uint16   `unionTag:"0"`
	Values  []uint16 `arrayCount:""`
}

type AATLoopkup2 struct {
	version uint16 `unionTag:"2"`
	binSearchHeader
	Records []LookupRecord2 `arrayCount:"ComputedField-nUnits"`
}

type LookupRecord2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint16
}

type AATLoopkup4 struct {
	version uint16 `unionTag:"4"`
	binSearchHeader
	// Do not include the termination segment
	Records []AATLookupRecord4 `arrayCount:"ComputedField-nUnits-1"`
}

type AATLookupRecord4 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	Values []uint16 `offsetSize:"Offset16" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nValues()"`
}

func (lk AATLookupRecord4) nValues() int { return int(lk.LastGlyph) - int(lk.FirstGlyph) + 1 }

type AATLoopkup6 struct {
	version uint16 `unionTag:"6"`
	binSearchHeader
	Records []loopkupRecord6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecord6 struct {
	Glyph GlyphID
	Value uint16
}

type AATLoopkup8 struct {
	version uint16 `unionTag:"8"`
	AATLoopkup8Data
}

type AATLoopkup8Data struct {
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

type AATLoopkup10 struct {
	version    uint16 `unionTag:"10"`
	unitSize   uint16
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

// extended versions

// AATLookupExt is the same as AATLookup, but class values are uint32
type AATLookupExt interface {
	AatLookupMixed

	isAATLookupExt()

	// Class returns the class ID for the provided glyph, or (0, false)
	// for glyphs not covered by this class.
	Class(g GlyphID) (uint32, bool)
}

func (AATLoopkupExt0) isAATLookupExt()  {}
func (AATLoopkupExt2) isAATLookupExt()  {}
func (AATLoopkupExt4) isAATLookupExt()  {}
func (AATLoopkupExt6) isAATLookupExt()  {}
func (AATLoopkupExt8) isAATLookupExt()  {}
func (AATLoopkupExt10) isAATLookupExt() {}

type AATLoopkupExt0 struct {
	version uint16   `unionTag:"0"`
	Values  []uint32 `arrayCount:""`
}

type AATLoopkupExt2 struct {
	version uint16 `unionTag:"2"`
	binSearchHeader
	Records []lookupRecordExt2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecordExt2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint32
}

type AATLoopkupExt4 struct {
	version uint16 `unionTag:"4"`
	binSearchHeader
	// the values pointed by the record are uint32
	Records []loopkupRecordExt4 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecordExt4 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	Values []uint32 `offsetSize:"Offset16" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nValues()"`
}

func (lk loopkupRecordExt4) nValues() int { return int(lk.LastGlyph) - int(lk.FirstGlyph) + 1 }

type AATLoopkupExt6 struct {
	version uint16 `unionTag:"6"`
	binSearchHeader
	Records []loopkupRecordExt6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecordExt6 struct {
	Glyph GlyphID
	Value uint32
}

type AATLoopkupExt8 AATLoopkup8

type AATLoopkupExt10 struct {
	version    uint16 `unionTag:"10"`
	unitSize   uint16
	FirstGlyph GlyphID
	Values     []uint32 `arrayCount:"FirstUint16"`
}
