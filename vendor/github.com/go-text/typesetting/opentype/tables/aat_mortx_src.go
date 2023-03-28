// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Morx is the extended glyph metamorphosis table
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6morx.html
type Morx struct {
	version uint16      // Version number of the extended glyph metamorphosis table (either 2 or 3)
	unused  uint16      // Set to 0
	nChains uint32      // Number of metamorphosis chains contained in this table.
	Chains  []MorxChain `arrayCount:"ComputedField-nChains"`
}

// MorxChain is a set of subtables
type MorxChain struct {
	Flags           uint32              // The default specification for subtables.
	chainLength     uint32              // Total byte count, including this header; must be a multiple of 4.
	nFeatureEntries uint32              // Number of feature subtable entries.
	nSubtable       uint32              // The number of subtables in the chain.
	Features        []AATFeature        `arrayCount:"ComputedField-nFeatureEntries"`
	Subtables       []MorxChainSubtable `arrayCount:"ComputedField-nSubtable"`
}

type AATFeature struct {
	FeatureType    uint16
	FeatureSetting uint16
	EnableFlags    uint32 // Flags for the settings that this feature and setting enables.
	DisableFlags   uint32 // Complement of flags for the settings that this feature and setting disable.
}

type MorxChainSubtable struct {
	length uint32 // Total subtable length, including this header.

	// Coverage flags and subtable type.
	Coverage byte
	ignored  [2]byte
	version  MorxSubtableVersion

	SubFeatureFlags uint32 // The 32-bit mask identifying which subtable this is (the subtable being executed if the AND of this value and the processed defaultFlags is nonzero)

	Data MorxSubtable `unionField:"version"`
}

// check and return the subtable length
func (mc *MorxChainSubtable) parseEnd(src []byte, _ int) (int, error) {
	if L := len(src); L < int(mc.length) {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", mc.length, L)
	}
	return int(mc.length), nil
}

// MorxSubtableVersion indicates the kind of 'morx' subtable.
// See the constants.
type MorxSubtableVersion uint8

const (
	MorxSubtableVersionRearrangement MorxSubtableVersion = iota
	MorxSubtableVersionContextual
	MorxSubtableVersionLigature
	_ // reserved
	MorxSubtableVersionNonContextual
	MorxSubtableVersionInsertion
)

type MorxSubtable interface {
	isMorxSubtable()
}

func (MorxSubtableRearrangement) isMorxSubtable() {}
func (MorxSubtableContextual) isMorxSubtable()    {}
func (MorxSubtableLigature) isMorxSubtable()      {}
func (MorxSubtableNonContextual) isMorxSubtable() {}
func (MorxSubtableInsertion) isMorxSubtable()     {}

// binarygen: argument=valuesCount int
type MorxSubtableRearrangement struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount,entryDataSize=0"`
}

// binarygen: argument=valuesCount int
type MorxSubtableContextual struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount,entryDataSize=4"`
	// Byte offset from the beginning of the state subtable to the beginning of the substitution tables :
	// each value of the array is itself an offet to a aatLookupTable, and the number of
	// items is computed from the header
	Substitutions SubstitutionsTable `offsetSize:"Offset32" arguments:"substitutionsCount=.nSubs(), valuesCount=valuesCount"`
}

type SubstitutionsTable struct {
	Substitutions []AATLookup `offsetsArray:"Offset32"`
}

func (ct *MorxSubtableContextual) nSubs() int {
	// find the maximum index need in the substitution array
	var maxi uint16
	for _, entry := range ct.Entries {
		markIndex, currentIndex := entry.AsMorxContextual()
		if markIndex != 0xFFFF && markIndex > maxi {
			maxi = markIndex
		}
		if currentIndex != 0xFFFF && currentIndex > maxi {
			maxi = currentIndex
		}
	}
	return int(maxi) + 1
}

// binarygen: argument=valuesCount int
type MorxSubtableLigature struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount, entryDataSize=2"`
	ligActionOffset  Offset32  // Byte offset from stateHeader to the start of the ligature action table.
	componentOffset  Offset32  // Byte offset from stateHeader to the start of the component table.
	ligatureOffset   Offset32  // Byte offset from stateHeader to the start of the actual ligature lists.
	LigActions       []uint32  `isOpaque:""`
	Components       []uint16  `isOpaque:""`
	Ligatures        []GlyphID `isOpaque:""`
}

// MorxLigatureSubtable flags
const (
	// Push this glyph onto the component stack for
	// eventual processing.
	MLSetComponent = 0x8000
	// Leave the glyph pointer at this glyph for the
	// next iteration.
	MLDontAdvance = 0x4000
	// Use the ligActionIndex to process a ligature group.
	MLPerformAction = 0x2000
	// Byte offset from beginning of subtable to the
	// ligature action list. This value must be a
	// multiple of 4.
	MLOffset = 0x3FFF

	// This is the last action in the list. This also
	// implies storage.
	MLActionLast = 1 << 31
	// Store the ligature at the current cumulated index
	// in the ligature table in place of the marked
	// (i.e. currently-popped) glyph.
	MLActionStore = 1 << 30
	// A 30-bit value which is sign-extended to 32-bits
	// and added to the glyph ID, resulting in an index
	// into the component table.
	MLActionOffset = 0x3FFFFFFF
)

// the LigActions length is not specified. Instead, we have to parse uint32 one by one
// until we reach last action or reach EOF
func (lig *MorxSubtableLigature) parseLigActions(src []byte, _ int) error {
	// fetch the maximum start index
	maxIndex := -1
	for _, entry := range lig.Entries {
		if entry.Flags&MLPerformAction == 0 {
			continue
		}
		if index := int(entry.AsMorxLigature()); index > maxIndex {
			maxIndex = index
		}
	}

	if L := len(src); L < int(lig.ligActionOffset)+4*int(maxIndex+1) {
		return fmt.Errorf("EOF: expected length: %d, got %d", lig.ligActionOffset, L)
	}

	// fetch the action table, up to the last entry
	src = src[lig.ligActionOffset:]
	for len(src) >= 4 { // stop gracefully if the last action was not found
		action := binary.BigEndian.Uint32(src)
		lig.LigActions = append(lig.LigActions, action)
		src = src[4:]
		// dont break before maxIndex
		if len(lig.LigActions) > maxIndex && action&MLActionLast != 0 {
			break
		}
	}
	return nil
}

func (lig *MorxSubtableLigature) parseComponents(src []byte, _ int) error {
	// we rely on offset being sorted, which seems to be the case in practice
	if lig.componentOffset > lig.ligatureOffset {
		return errors.New("unsupported non sorted offsets")
	}
	if L := len(src); L < int(lig.componentOffset) {
		return fmt.Errorf("EOF: expected length: %d, got %d", lig.componentOffset, L)
	}
	src = src[lig.componentOffset:]
	componentCount := (lig.ligatureOffset - lig.componentOffset) / 2
	lig.Components = make([]uint16, componentCount)
	for i := range lig.Components {
		lig.Components[i] = binary.BigEndian.Uint16(src[2*i:])
	}
	return nil
}

func (lig *MorxSubtableLigature) parseLigatures(src []byte, _ int) error {
	if L := len(src); L < int(lig.ligatureOffset) {
		return fmt.Errorf("EOF: expected length: %d, got %d", lig.ligatureOffset, L)
	}
	src = src[lig.ligatureOffset:]
	ligatureCount := len(src) / 2
	lig.Ligatures = make([]GlyphID, ligatureCount)
	for i := range lig.Ligatures {
		lig.Ligatures[i] = GlyphID(binary.BigEndian.Uint16(src[2*i:]))
	}
	return nil
}

type MorxSubtableNonContextual struct {
	// The lookup value is interpreted as a GlyphIndex
	Class AATLookup
}

// binarygen: argument=valuesCount int
type MorxSubtableInsertion struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount,entryDataSize=4"`
	Insertions       []GlyphID `offsetSize:"Offset32" arrayCount:"ComputedField-nInsertions()"` // Byte offset from stateHeader to the start of the insertion glyph table.
}

// MorxInsertionSubtable flags
const (
	// If set, mark the current glyph.
	MISetMark = 0x8000
	// If set, don't advance to the next glyph before
	// going to the new state.  This does not mean
	// that the glyph pointed to is the same one as
	// before. If you've made insertions immediately
	// downstream of the current glyph, the next glyph
	// processed would in fact be the first one
	// inserted.
	MIDontAdvance = 0x4000
	// If set, and the currentInsertList is nonzero,
	// then the specified glyph list will be inserted
	// as a kashida-like insertion, either before or
	// after the current glyph (depending on the state
	// of the currentInsertBefore flag). If clear, and
	// the currentInsertList is nonzero, then the
	// specified glyph list will be inserted as a
	// split-vowel-like insertion, either before or
	// after the current glyph (depending on the state
	// of the currentInsertBefore flag).
	MICurrentIsKashidaLike = 0x2000
	// If set, and the markedInsertList is nonzero,
	// then the specified glyph list will be inserted
	// as a kashida-like insertion, either before or
	// after the marked glyph (depending on the state
	// of the markedInsertBefore flag). If clear, and
	// the markedInsertList is nonzero, then the
	// specified glyph list will be inserted as a
	// split-vowel-like insertion, either before or
	// after the marked glyph (depending on the state
	// of the markedInsertBefore flag).
	MIMarkedIsKashidaLike = 0x1000
	// If set, specifies that insertions are to be made
	// to the left of the current glyph. If clear,
	// they're made to the right of the current glyph.
	MICurrentInsertBefore = 0x0800
	// If set, specifies that insertions are to be
	// made to the left of the marked glyph. If clear,
	// they're made to the right of the marked glyph.
	MIMarkedInsertBefore = 0x0400
	// This 5-bit field is treated as a count of the
	// number of glyphs to insert at the current
	// position. Since zero means no insertions, the
	// largest number of insertions at any given
	// current location is 31 glyphs.
	MICurrentInsertCount = 0x3E0
	// This 5-bit field is treated as a count of the
	// number of glyphs to insert at the marked
	// position. Since zero means no insertions, the
	// largest number of insertions at any given
	// marked location is 31 glyphs.
	MIMarkedInsertCount = 0x001F
)

func (msi *MorxSubtableInsertion) nInsertions() int {
	// find the maximum index needed in the insertions array,
	// taking into account the number of insertions
	var maxi uint16
	for _, entry := range msi.Entries {
		currentIndex, markedIndex := entry.AsMorxInsertion()
		if currentIndex != 0xFFFF {
			indexEnd := currentIndex + (entry.Flags&MICurrentInsertCount)>>5
			if indexEnd > maxi {
				maxi = indexEnd
			}
		}
		if markedIndex != 0xFFFF {
			indexEnd := markedIndex + entry.Flags&MIMarkedInsertCount
			if indexEnd > maxi {
				maxi = indexEnd
			}
		}
	}
	return int(maxi)
}
