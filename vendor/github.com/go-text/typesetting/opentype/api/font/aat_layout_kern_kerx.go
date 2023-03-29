// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import (
	"encoding/binary"

	"github.com/go-text/typesetting/opentype/tables"
)

// Kernx represents a 'kern' or 'kerx' kerning table.
// It supports both Microsoft and Apple formats.
type Kernx []KernSubtable

func newKernxFromKerx(kerx tables.Kerx) Kernx {
	if len(kerx.Tables) == 0 {
		return nil
	}
	out := make(Kernx, len(kerx.Tables))
	for i, ta := range kerx.Tables {
		out[i] = newKerxSubtable(ta)
	}
	return out
}

func newKernxFromKern(kern tables.Kern) Kernx {
	if len(kern.Tables) == 0 {
		return nil
	}
	out := make(Kernx, len(kern.Tables))
	for i, ta := range kern.Tables {
		out[i] = newKernSubtable(ta)
	}
	return out
}

// KernSubtable represents a 'kern' or 'kerx' subtable.
type KernSubtable struct {
	Data interface{ isKernSubtable() }

	// high bit of the Coverage field, following 'kerx' conventions
	coverage byte

	// IsExtended [true] for AAT `kerx` subtables, false for 'kern' subtables
	IsExtended bool

	// 0 for scalar values
	TupleCount int
}

func newKernSubtable(table tables.KernSubtable) (out KernSubtable) {
	out.IsExtended = false
	switch table := table.(type) {
	case tables.OTKernSubtableHeader:
		// synthesize a coverage flag following kerx conventions
		const (
			Horizontal  = 0x01
			CrossStream = 0x04
		)
		if table.Coverage&Horizontal == 0 { // vertical
			out.coverage |= kerxVertical
		}
		if table.Coverage&CrossStream != 0 {
			out.coverage |= kerxCrossStream
		}
	case tables.AATKernSubtableHeader:
		out.coverage = table.Coverage
		out.TupleCount = int(table.TupleCount)
	}
	switch data := table.Data().(type) {
	case tables.KernData0:
		out.Data = newKern0(data)
	case tables.KernData1:
		out.Data = newKern1(data)
	case tables.KernData2:
		out.Data = newKern2(data)
	case tables.KernData3:
		out.Data = Kern3(data)
	}
	return out
}

func newKerxSubtable(table tables.KerxSubtable) (out KernSubtable) {
	out.IsExtended = true
	out.TupleCount = int(table.TupleCount)
	out.coverage = byte(table.Coverage >> 8) // high bit only

	switch data := table.Data.(type) {
	case tables.KerxData0:
		out.Data = newKern0x(data)
	case tables.KerxData1:
		out.Data = newKern1x(data)
	case tables.KerxData2:
		out.Data = Kern2(data)
	case tables.KerxData4:
		out.Data = newKern4(data)
	case tables.KerxData6:
		out.Data = Kern6(data)
	}
	return out
}

func (Kern0) isKernSubtable() {}
func (Kern1) isKernSubtable() {}
func (Kern2) isKernSubtable() {}
func (Kern3) isKernSubtable() {}
func (Kern4) isKernSubtable() {}
func (Kern6) isKernSubtable() {}

var (
	_ SimpleKerns = Kern0(nil)
	_ SimpleKerns = (*Kern2)(nil)
	_ SimpleKerns = (*Kern3)(nil)
	_ SimpleKerns = (*Kern6)(nil)
)

// SimpleKerns store a compact form of the kerning values,
// which is restricted to (one direction) kerning pairs.
// It is only implemented by [Kern0], [Kern2], [Kern3] and [Kern6],
// where [Kern1] and [Kern4] requires a state machine to be interpreted.
type SimpleKerns interface {
	// KernPair return the kern value for the given pair, or zero.
	// The value is expressed in glyph units and
	// is negative when glyphs should be closer.
	KernPair(left, right GID) int16
}

// kernx coverage flags
const (
	kerxBackwards   = 1 << (12 - 8)
	kerxVariation   = 1 << (13 - 8)
	kerxCrossStream = 1 << (14 - 8)
	kerxVertical    = 1 << (15 - 8)
)

// IsHorizontal returns true if the subtable has horizontal kerning values.
func (k KernSubtable) IsHorizontal() bool { return k.coverage&kerxVertical == 0 }

// IsBackwards returns true if state-table based should process the glyphs backwards.
func (k KernSubtable) IsBackwards() bool { return k.coverage&kerxBackwards != 0 }

// IsCrossStream returns true if the subtable has cross-stream kerning values.
func (k KernSubtable) IsCrossStream() bool { return k.coverage&kerxCrossStream != 0 }

// IsVariation returns true if the subtable has variation kerning values.
func (k KernSubtable) IsVariation() bool { return k.coverage&kerxVariation != 0 }

type Kern0 []tables.Kernx0Record

func newKern0(k tables.KernData0) Kern0  { return k.Pairs }
func newKern0x(k tables.KerxData0) Kern0 { return k.Pairs }

func kernPair(records []tables.Kernx0Record, left, right GID) int16 {
	key := uint32(left)<<16 | uint32(right)
	low, high := 0, len(records)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := recordKey(records[mid])
		if key < p {
			high = mid
		} else if key > p {
			low = mid + 1
		} else {
			return records[mid].Value
		}
	}
	return 0
}

func recordKey(kp tables.Kernx0Record) uint32 { return uint32(kp.Left)<<16 | uint32(kp.Right) }

func (kd Kern0) KernPair(left, right GID) int16 { return kernPair(kd, left, right) }

type Kern1 struct {
	Values  []int16 // After successful parsing, may be safely indexed by AATStateEntry.AsKernxIndex() from `Machine`
	Machine AATStateTable
}

// convert from non extended to extended
func newKern1(k tables.KernData1) Kern1 {
	class := tables.AATLoopkup8{
		AATLoopkup8Data: tables.AATLoopkup8Data{
			FirstGlyph: k.ClassTable.StartGlyph,
			Values:     make([]uint16, len(k.ClassTable.Values)),
		},
	}
	for i, b := range k.ClassTable.Values {
		class.Values[i] = uint16(b)
	}
	states := make([][]uint16, len(k.States))
	for i, row := range k.States {
		v := make([]uint16, len(row))
		for j, b := range row {
			v[j] = uint16(b)
		}
		states[i] = v
	}
	return Kern1{
		Values: k.Values,
		Machine: AATStateTable{
			nClass:  uint32(k.StateSize),
			class:   class,
			states:  states,
			entries: k.Entries,
		},
	}
}

func newKern1x(k tables.KerxData1) Kern1 {
	return Kern1{Values: k.Values, Machine: newAATStableTable(k.AATStateTableExt)}
}

type Kern2 tables.KerxData2

// convert from non extended to extended
func newKern2(k tables.KernData2) Kern2 {
	return Kern2{
		Left:         tables.AATLoopkup8{AATLoopkup8Data: k.Left},
		Right:        tables.AATLoopkup8{AATLoopkup8Data: k.Right},
		KerningStart: tables.Offset32(k.KerningStart),
		KerningData:  k.KerningData,
	}
}

func (kd Kern2) KernPair(left, right GID) int16 {
	l, _ := kd.Left.Class(tables.GlyphID(left))
	r, _ := kd.Right.Class(tables.GlyphID(right))
	index := int(l) + int(r)
	if len(kd.KerningData) < index+2 || index < int(kd.KerningStart) {
		return 0
	}
	kernVal := binary.BigEndian.Uint16(kd.KerningData[index:])
	return int16(kernVal)
}

type Kern3 tables.KernData3

func (kd Kern3) KernPair(left, right GID) int16 {
	if int(left) >= len(kd.LeftClass) || int(right) >= len(kd.RightClass) { // should not happend
		return 0
	}

	lc, rc := int(kd.LeftClass[left]), int(kd.RightClass[right])
	index := kd.KernIndex[lc*int(kd.RightClassCount)+rc] // sanitized during parsing
	return kd.Kernings[index]                            // sanitized during parsing
}

type Kern4 struct {
	Anchors tables.KerxAnchors
	Machine AATStateTable
	flags   uint32
}

func newKern4(k tables.KerxData4) Kern4 {
	return Kern4{
		Machine: newAATStableTable(k.AATStateTableExt),
		Anchors: k.Anchors,
		flags:   k.Flags,
	}
}

// ActionType returns 0, 1 or 2 .
func (k Kern4) ActionType() uint8 {
	const ActionType = 0xC0000000 // A two-bit field containing the action type.
	return uint8(k.flags & ActionType >> 30)
}

type Kern6 tables.KerxData6

func (kd Kern6) KernPair(left, right GID) int16 {
	l := kd.Row.ClassUint32(tables.GlyphID(left))
	r := kd.Column.ClassUint32(tables.GlyphID(right))
	index := int(l) + int(r)
	if len(kd.Kernings) <= index {
		return 0
	}
	return kd.Kernings[index]
}

// --------------------------------------- state machine ---------------------------------------

// AATStateTable supports both regular and extended AAT state machines
type AATStateTable struct {
	nClass  uint32
	class   tables.AATLookup
	states  [][]uint16             // each sub array has length stateSize
	entries []tables.AATStateEntry // length is the maximum state + 1
}

func newAATStableTable(k tables.AATStateTableExt) AATStateTable {
	return AATStateTable{
		nClass:  k.StateSize,
		class:   k.Class,
		states:  k.States,
		entries: k.Entries,
	}
}

// GetClass return the class for the given glyph, with the correct default value.
func (st *AATStateTable) GetClass(glyph GID) uint16 {
	if glyph == 0xFFFF { // deleted glyph
		return 2 // class deleted
	}
	c, ok := st.class.Class(tables.GlyphID(glyph))
	if !ok {
		return 1 // class out of bounds
	}
	return c // class for a state table can't be uint32
}

// GetEntry return the entry for the given state and class,
// and handle invalid values (by returning an empty entry).
func (st *AATStateTable) GetEntry(state, class uint16) tables.AATStateEntry {
	if uint32(class) >= st.nClass {
		class = 1 // class out of bounds
	}
	if int(state) >= len(st.states) {
		return tables.AATStateEntry{}
	}
	entry := st.states[state][class] // access check when parsing
	return st.entries[entry]         // access check when parsing
}
