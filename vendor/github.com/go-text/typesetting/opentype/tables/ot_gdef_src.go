// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"fmt"
)

type GDEF struct {
	majorVersion    uint16       // Major version of the GDEF table, = 1
	minorVersion    uint16       // Minor version of the GDEF table, = 0, 2, 3
	GlyphClassDef   ClassDef     `offsetSize:"Offset16"` // Offset to class definition table for glyph type, from beginning of GDEF header (may be NULL)
	AttachList      AttachList   `offsetSize:"Offset16"` // Offset to attachment point list table, from beginning of GDEF header (may be NULL)
	LigCaretList    LigCaretList `offsetSize:"Offset16"` // Offset to ligature caret list table, from beginning of GDEF header (may be NULL)
	MarkAttachClass ClassDef     `offsetSize:"Offset16"` // Offset to class definition table for mark attachment type, from beginning of GDEF header (may be NULL)

	MarkGlyphSetsDef MarkGlyphSets `isOpaque:""` // Offset to the table of mark glyph set definitions, from beginning of GDEF header (may be NULL)
	ItemVarStore     ItemVarStore  `isOpaque:""` // Offset to the Item Variation Store table, from beginning of GDEF header (may be NULL)
}

func (gdef *GDEF) parseMarkGlyphSetsDef(src []byte) error {
	const headerSize = 12
	if gdef.minorVersion < 2 {
		return nil
	}
	if L := len(src); L < headerSize+2 {
		return fmt.Errorf("EOF: expected length: %d, got %d", headerSize+2, L)
	}
	offset := binary.BigEndian.Uint16(src[headerSize:])
	if offset != 0 {
		var err error
		gdef.MarkGlyphSetsDef, _, err = ParseMarkGlyphSets(src[offset:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (gdef *GDEF) parseItemVarStore(src []byte) (int, error) {
	const headerSize = 12 + 2
	if gdef.minorVersion < 3 {
		return 0, nil
	}
	if L := len(src); L < headerSize+4 {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", headerSize+4, L)
	}
	offset := binary.BigEndian.Uint32(src[headerSize:])
	if offset != 0 {
		var err error
		gdef.ItemVarStore, _, err = ParseItemVarStore(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return headerSize + 4, nil
}

type AttachList struct {
	Coverage     Coverage      `offsetSize:"Offset16"`                            // Offset to Coverage table - from beginning of AttachList table
	AttachPoints []AttachPoint `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // [glyphCount] Array of offsets to AttachPoint tables-from beginning of AttachList table-in Coverage Index order
}

type AttachPoint struct {
	PointIndices []uint16 `arrayCount:"FirstUint16"` // [pointCount]	Array of contour point indices -in increasing numerical order
}

type LigCaretList struct {
	Coverage  Coverage   `offsetSize:"Offset16"`                            // Offset to Coverage table - from beginning of LigCaretList table
	LigGlyphs []LigGlyph `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // [ligGlyphCount]	Array of offsets to LigGlyph tables, from beginning of LigCaretList table —in Coverage Index order
}

type LigGlyph struct {
	CaretValues []CaretValue `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // [caretCount] Array of offsets to CaretValue tables, from beginning of LigGlyph table — in increasing coordinate order
}

type CaretValue interface {
	isCaretValue()
}

func (CaretValue1) isCaretValue() {}
func (CaretValue2) isCaretValue() {}
func (CaretValue3) isCaretValue() {}

type CaretValue1 struct {
	caretValueFormat uint16 `unionTag:"1"` // Format identifier: format = 1
	Coordinate       int16  //	X or Y value, in design units
}

type CaretValue2 struct {
	caretValueFormat     uint16 `unionTag:"2"` // Format identifier: format = 2
	CaretValuePointIndex uint16 // Contour point index on glyph
}

type CaretValue3 struct {
	caretValueFormat uint16      `unionTag:"3"` // Format identifier: format = 3
	Coordinate       int16       // X or Y value, in design units
	deviceOffset     Offset16    // Offset to Device table (non-variable font) / Variation Index table (variable font) for X or Y value-from beginning of CaretValue table
	Device           DeviceTable `isOpaque:""`
}

func (cv *CaretValue3) parseDevice(src []byte) (err error) {
	cv.Device, err = parseDeviceTable(src, uint16(cv.deviceOffset))
	return err
}

type MarkGlyphSets struct {
	format    uint16     // Format identifier == 1
	Coverages []Coverage `arrayCount:"FirstUint16" offsetsArray:"Offset32"` // [markGlyphSetCount] Array of offsets to mark glyph set coverage tables, from the start of the MarkGlyphSets table.
}
