// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"fmt"
)

type SinglePos struct {
	Data SinglePosData
}

type SinglePosData interface {
	isSinglePosData()

	Cov() Coverage
}

func (SinglePosData1) isSinglePosData() {}
func (SinglePosData2) isSinglePosData() {}

type SinglePosData1 struct {
	format      uint16      `unionTag:"1"`
	coverage    Coverage    `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of SinglePos subtable.
	ValueFormat ValueFormat //	Defines the types of data in the ValueRecord.
	ValueRecord ValueRecord `isOpaque:""` //	Defines positioning value(s) — applied to all glyphs in the Coverage table.
}

func (sp *SinglePosData1) parseValueRecord(src []byte) (err error) {
	sp.ValueRecord, _, err = parseValueRecord(sp.ValueFormat, src, 6)
	return err
}

type SinglePosData2 struct {
	format       uint16        `unionTag:"2"`
	coverage     Coverage      `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of SinglePos subtable.
	ValueFormat  ValueFormat   // Defines the types of data in the ValueRecords.
	valueCount   uint16        // Number of ValueRecords — must equal glyphCount in the Coverage table.
	ValueRecords []ValueRecord `isOpaque:""` //[valueCount]	Array of ValueRecords — positioning values applied to glyphs.
}

func (sp *SinglePosData2) parseValueRecords(src []byte) (err error) {
	offset := 8
	sp.ValueRecords = make([]ValueRecord, sp.valueCount)
	for i := range sp.ValueRecords {
		sp.ValueRecords[i], offset, err = parseValueRecord(sp.ValueFormat, src, offset)
		if err != nil {
			return err
		}
	}

	return err
}

type PairPos struct {
	Data PairPosData
}

type PairPosData interface {
	isPairPosData()

	Cov() Coverage
}

func (PairPosData1) isPairPosData() {}
func (PairPosData2) isPairPosData() {}

type PairPosData1 struct {
	format   uint16   `unionTag:"1"`
	coverage Coverage `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of PairPos subtable.

	ValueFormat1 ValueFormat // Defines the types of data in valueRecord1 — for the first glyph in the pair (may be zero).
	ValueFormat2 ValueFormat // Defines the types of data in valueRecord2 — for the second glyph in the pair (may be zero).
	PairSets     []PairSet   `arrayCount:"FirstUint16" offsetsArray:"Offset16" arguments:"valueFormat1=.ValueFormat1, valueFormat2=.ValueFormat2"` //[pairSetCount] Array of offsets to PairSet tables. Offsets are from beginning of PairPos subtable, ordered by Coverage Index.
}

// binarygen: argument=valueFormat1  ValueFormat
// binarygen: argument=valueFormat2  ValueFormat
type PairSet struct {
	pairValueCount   uint16            // Number of PairValueRecords
	PairValueRecords []PairValueRecord `isOpaque:""` // [pairValueCount] Array of PairValueRecords, ordered by glyph ID of the second glyph.
}

func (ps *PairSet) parsePairValueRecords(src []byte, fmt1, fmt2 ValueFormat) error {
	out := make([]PairValueRecord, ps.pairValueCount)
	offsetR := 2
	var err error
	for i := range out {
		if L := len(src); L < 2+offsetR {
			return fmt.Errorf("EOF: expected length: %d, got %d", 2+offsetR, L)
		}
		out[i].SecondGlyph = GlyphID(binary.BigEndian.Uint16(src[offsetR:]))
		out[i].ValueRecord1, offsetR, err = parseValueRecord(fmt1, src, offsetR+2)
		if err != nil {
			return fmt.Errorf("invalid pair set table: %s", err)
		}
		out[i].ValueRecord2, offsetR, err = parseValueRecord(fmt2, src, offsetR)
		if err != nil {
			return fmt.Errorf("invalid pair set table: %s", err)
		}
	}
	ps.PairValueRecords = out
	return nil
}

type PairPosData2 struct {
	format   uint16   `unionTag:"2"`
	coverage Coverage `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of PairPos subtable.

	ValueFormat1 ValueFormat //	Defines the types of data in valueRecord1 — for the first glyph in the pair (may be zero).
	ValueFormat2 ValueFormat //	Defines the types of data in valueRecord2 — for the second glyph in the pair (may be zero).

	ClassDef1     ClassDef       `offsetSize:"Offset16"` // Offset to ClassDef table, from beginning of PairPos subtable — for the first glyph of the pair.
	ClassDef2     ClassDef       `offsetSize:"Offset16"` // Offset to ClassDef table, from beginning of PairPos subtable — for the second glyph of the pair.
	class1Count   uint16         //	Number of classes in classDef1 table — includes Class 0.
	class2Count   uint16         //	Number of classes in classDef2 table — includes Class 0.
	Class1Records []Class1Record `isOpaque:""` //[class1Count]	Array of Class1 records, ordered by classes in classDef1.
}

func (pp *PairPosData2) parseClass1Records(src []byte) error {
	const headerSize = 16 // including posFormat and coverageOffset

	pp.Class1Records = make([]Class1Record, pp.class1Count)

	offset := headerSize
	for i := range pp.Class1Records {
		vi := make(Class1Record, pp.class2Count)
		for j := range vi {
			var err error
			vi[j].ValueRecord1, offset, err = parseValueRecord(pp.ValueFormat1, src, offset)
			if err != nil {
				return err
			}
			vi[j].ValueRecord2, offset, err = parseValueRecord(pp.ValueFormat2, src, offset)
			if err != nil {
				return err
			}
		}
		pp.Class1Records[i] = vi
	}
	return nil
}

// DeviceTableHeader is the common header for DeviceTable
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/chapter2#device-and-variationindex-tables
type DeviceTableHeader struct {
	first       uint16
	second      uint16
	deltaFormat uint16 // Format of deltaValue array data
}

type Anchor interface {
	isAnchor()
}

type EntryExit struct {
	EntryAnchor Anchor
	ExitAnchor  Anchor
}

type CursivePos struct {
	posFormat        uint16            //	Format identifier: format = 1
	coverage         Coverage          `offsetSize:"Offset16"`    //	Offset to Coverage table, from beginning of CursivePos subtable.
	entryExitRecords []entryExitRecord `arrayCount:"FirstUint16"` //[entryExitCount]	Array of EntryExit records, in Coverage index order.
	EntryExits       []EntryExit       `isOpaque:""`
}

type entryExitRecord struct {
	entryAnchorOffset Offset16 // Offset to entryAnchor table, from beginning of CursivePos subtable (may be NULL).
	exitAnchorOffset  Offset16 // Offset to exitAnchor table, from beginning of CursivePos subtable (may be NULL).
}

func (cp *CursivePos) parseEntryExits(src []byte) error {
	cp.EntryExits = make([]EntryExit, len(cp.entryExitRecords))
	var err error
	for i, rec := range cp.entryExitRecords {
		if rec.entryAnchorOffset != 0 {
			if L := len(src); L < int(rec.entryAnchorOffset) {
				return fmt.Errorf("EOF: expected length: %d, got %d", rec.entryAnchorOffset, L)
			}
			cp.EntryExits[i].EntryAnchor, _, err = ParseAnchor(src[rec.entryAnchorOffset:])
			if err != nil {
				return err
			}
		}
		if rec.exitAnchorOffset != 0 {
			if L := len(src); L < int(rec.exitAnchorOffset) {
				return fmt.Errorf("EOF: expected length: %d, got %d", rec.exitAnchorOffset, L)
			}
			cp.EntryExits[i].ExitAnchor, _, err = ParseAnchor(src[rec.exitAnchorOffset:])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type MarkBasePos struct {
	posFormat      uint16    // Format identifier: format = 1
	markCoverage   Coverage  `offsetSize:"Offset16"` // Offset to markCoverage table, from beginning of MarkBasePos subtable.
	BaseCoverage   Coverage  `offsetSize:"Offset16"` // Offset to baseCoverage table, from beginning of MarkBasePos subtable.
	markClassCount uint16    // Number of classes defined for marks
	MarkArray      MarkArray `offsetSize:"Offset16"`                                          // Offset to MarkArray table, from beginning of MarkBasePos subtable.
	BaseArray      BaseArray `offsetSize:"Offset16" arguments:"offsetsCount=.markClassCount"` // Offset to BaseArray table, from beginning of MarkBasePos subtable.
}

type BaseArray struct {
	baseRecords []anchorOffsets `arrayCount:"FirstUint16"` //  [markClassCount] Array of offsets (one per mark class) to Anchor tables. Offsets are from beginning of BaseArray table, ordered by class (offsets may be NULL).
	BaseAnchors [][]Anchor      `isOpaque:""`
}

func (ba *BaseArray) parseBaseAnchors(src []byte, _ int) (err error) {
	ba.BaseAnchors, err = resolveAnchorOffsets(ba.baseRecords, src)
	return err
}

type anchorOffsets struct {
	offsets []Offset16 // Array of offsets to Anchor tables, with external length
}

// resolveAnchorOffsets resolve the offsset using the given input slice
func resolveAnchorOffsets(offsets []anchorOffsets, src []byte) ([][]Anchor, error) {
	out := make([][]Anchor, len(offsets))
	var err error
	for i, list := range offsets {
		bi := make([]Anchor, len(list.offsets))
		for j, offset := range list.offsets {
			if offset == 0 {
				continue
			}

			if L := len(src); L < int(offset) {
				return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
			}
			bi[j], _, err = ParseAnchor(src[offset:])
			if err != nil {
				return nil, err
			}
		}
		out[i] = bi
	}
	return out, nil
}

type MarkLigPos struct {
	posFormat        uint16        // Format identifier: format = 1
	MarkCoverage     Coverage      `offsetSize:"Offset16"` // Offset to markCoverage table, from beginning of MarkLigPos subtable.
	LigatureCoverage Coverage      `offsetSize:"Offset16"` // Offset to ligatureCoverage table, from beginning of MarkLigPos subtable.
	MarkClassCount   uint16        // Number of defined mark classes
	MarkArray        MarkArray     `offsetSize:"Offset16"`                                          // Offset to MarkArray table, from beginning of MarkLigPos subtable.
	LigatureArray    LigatureArray `offsetSize:"Offset16" arguments:"offsetsCount=.MarkClassCount"` // Offset to LigatureArray table, from beginning of MarkLigPos subtable.
}

type LigatureArray struct {
	LigatureAttachs []LigatureAttach `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // [ligatureCount]	Array of offsets to LigatureAttach tables. Offsets are from beginning of LigatureArray table, ordered by ligatureCoverage index.
}

type LigatureAttach struct {
	// [componentCount]	Array of Component records, ordered in writing direction.
	// Each element is an array of offsets (one per class, length = [markClassCount]) to Anchor tables. Offsets are from beginning of LigatureAttach table, ordered by class (offsets may be NULL).
	componentRecords []anchorOffsets `arrayCount:"FirstUint16"`
	ComponentAnchors [][]Anchor      `isOpaque:""`
}

func (la *LigatureAttach) parseComponentAnchors(src []byte, _ int) (err error) {
	la.ComponentAnchors, err = resolveAnchorOffsets(la.componentRecords, src)
	return err
}

type MarkMarkPos struct {
	PosFormat      uint16     //	Format identifier: format = 1
	Mark1Coverage  Coverage   `offsetSize:"Offset16"` // Offset to Combining Mark Coverage table, from beginning of MarkMarkPos subtable.
	Mark2Coverage  Coverage   `offsetSize:"Offset16"` // Offset to Base Mark Coverage table, from beginning of MarkMarkPos subtable.
	MarkClassCount uint16     //	Number of Combining Mark classes defined
	Mark1Array     MarkArray  `offsetSize:"Offset16"`                                          //	Offset to MarkArray table for mark1, from beginning of MarkMarkPos subtable.
	Mark2Array     Mark2Array `offsetSize:"Offset16" arguments:"offsetsCount=.MarkClassCount"` //	Offset to Mark2Array table for mark2, from beginning of MarkMarkPos subtable.
}

type Mark2Array struct {
	// [mark2Count]	Array of Mark2Records, in Coverage order.
	// Each element if an array of offsets (one per class, length = [markClassCount]) to Anchor tables. Offsets are from beginning of Mark2Array table, in class order (offsets may be NULL).
	mark2Records []anchorOffsets `arrayCount:"FirstUint16"`
	Mark2Anchors [][]Anchor      `isOpaque:""`
}

func (ma *Mark2Array) parseMark2Anchors(src []byte, _ int) (err error) {
	ma.Mark2Anchors, err = resolveAnchorOffsets(ma.mark2Records, src)
	return err
}

type ContextualPos struct {
	Data ContextualPosITF
}

type ContextualPosITF interface {
	isContextualPosITF()

	Cov() Coverage
}

type (
	ContextualPos1 SequenceContextFormat1
	ContextualPos2 SequenceContextFormat2
	ContextualPos3 SequenceContextFormat3
)

func (ContextualPos1) isContextualPosITF() {}
func (ContextualPos2) isContextualPosITF() {}
func (ContextualPos3) isContextualPosITF() {}

type ChainedContextualPos struct {
	Data ChainedContextualPosITF
}

type ChainedContextualPosITF interface {
	isChainedContextualPosITF()

	Cov() Coverage
}

type (
	ChainedContextualPos1 ChainedSequenceContextFormat1
	ChainedContextualPos2 ChainedSequenceContextFormat2
	ChainedContextualPos3 ChainedSequenceContextFormat3
)

func (ChainedContextualPos1) isChainedContextualPosITF() {}
func (ChainedContextualPos2) isChainedContextualPosITF() {}
func (ChainedContextualPos3) isChainedContextualPosITF() {}

type ExtensionPos Extension
