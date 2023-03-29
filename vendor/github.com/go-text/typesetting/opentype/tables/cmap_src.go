// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// Cmap is the Character to Glyph Index Mapping table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap
type Cmap struct {
	version   uint16           // Table version number (0).
	numTables uint16           // Number of encoding tables that follow.
	Records   []EncodingRecord `arrayCount:"ComputedField-numTables"`
}

type EncodingRecord struct {
	PlatformID PlatformID   // Platform ID.
	EncodingID EncodingID   // Platform-specific encoding ID.
	Subtable   CmapSubtable `offsetSize:"Offset32" offsetRelativeTo:"Parent"` // Byte offset from beginning of table to the subtable for this encoding.
}

// CmapSubtable is the union type for the various cmap formats
type CmapSubtable interface {
	isCmapSubtable()
}

func (CmapSubtable0) isCmapSubtable()  {}
func (CmapSubtable2) isCmapSubtable()  {}
func (CmapSubtable4) isCmapSubtable()  {}
func (CmapSubtable6) isCmapSubtable()  {}
func (CmapSubtable10) isCmapSubtable() {}
func (CmapSubtable12) isCmapSubtable() {}
func (CmapSubtable13) isCmapSubtable() {}
func (CmapSubtable14) isCmapSubtable() {}

type CmapSubtable0 struct {
	format       uint16 `unionTag:"0"` // Format number is set to 0.
	length       uint16 // This is the length in bytes of the subtable.
	language     uint16
	GlyphIdArray [256]uint8 //	An array that maps character codes to glyph index values.
}

type CmapSubtable2 struct {
	format  uint16 `unionTag:"2"` // Format number is set to 2.
	rawData []byte `arrayCount:"ToEnd"`
}

type CmapSubtable4 struct {
	format         uint16 `unionTag:"4"` // Format number is set to 4.
	length         uint16 // This is the length in bytes of the subtable.
	language       uint16
	segCountX2     uint16   // 2 × segCount.
	searchRange    uint16   // Maximum power of 2 less than or equal to segCount, times 2 ((2**floor(log2(segCount))) * 2, where “**” is an exponentiation operator)
	entrySelector  uint16   // Log2 of the maximum power of 2 less than or equal to numTables (log2(searchRange/2), which is equal to floor(log2(segCount)))
	rangeShift     uint16   // segCount times 2, minus searchRange ((segCount * 2) - searchRange)
	EndCode        []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]uint16 End characterCode for each segment, last=0xFFFF.
	reservedPad    uint16   // Set to 0.
	StartCode      []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]uint16 Start character code for each segment.
	IdDelta        []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]int16 Delta for all character codes in segment.
	IdRangeOffsets []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]uint16 Offsets into glyphIdArray or 0
	GlyphIDArray   []byte   `arrayCount:"ToEnd"`                        // glyphIdArray : uint16[] glyph index array (arbitrary length)
}

type CmapSubtable6 struct {
	format       uint16 `unionTag:"6"` // Format number is set to 6.
	length       uint16 // This is the length in bytes of the subtable.
	language     uint16
	FirstCode    uint16    // First character code of subrange.
	GlyphIdArray []GlyphID `arrayCount:"FirstUint16"` // Array of glyph index values for character codes in the range.
}

type CmapSubtable10 struct {
	format        uint16 `unionTag:"10"` //	Subtable format; set to 10.
	reserved      uint16 //	Reserved; set to 0
	length        uint32 //	Byte length of this subtable (including the header)
	language      uint32
	StartCharCode uint32    // First character code covered
	GlyphIdArray  []GlyphID `arrayCount:"FirstUint32"` // Array of glyph indices for the character codes covered
}

type CmapSubtable12 struct {
	format   uint16               `unionTag:"12"` //	Subtable format; set to 12.
	reserved uint16               //	Reserved; set to 0
	length   uint32               //	Byte length of this subtable (including the header)
	language uint32               //	For requirements on use of the language field, see “Use of the language field in 'cmap' subtables” in this document.
	Groups   []SequentialMapGroup `arrayCount:"FirstUint32"` // Array of SequentialMapGroup records.
}

type SequentialMapGroup struct {
	StartCharCode uint32 //	First character code in this group
	EndCharCode   uint32 //	Last character code in this group
	StartGlyphID  uint32 //	Glyph index corresponding to the starting character code
}

type CmapSubtable13 struct {
	format   uint16               `unionTag:"13"` //	Subtable format; set to 13.
	reserved uint16               //	Reserved; set to 0
	length   uint32               //	Byte length of this subtable (including the header)
	language uint32               //	For requirements on use of the language field, see “Use of the language field in 'cmap' subtables” in this document.
	Groups   []SequentialMapGroup `arrayCount:"FirstUint32"` // Array of SequentialMapGroup records.
}

type CmapSubtable14 struct {
	format       uint16              `unionTag:"14"` // Subtable format. Set to 14.
	length       uint32              // Byte length of this subtable (including this header)
	VarSelectors []VariationSelector `arrayCount:"FirstUint32"` // [numVarSelectorRecords]	Array of VariationSelector records.
}

type VariationSelector struct {
	VarSelector   [3]byte         // uint24 Variation selector
	DefaultUVS    DefaultUVSTable `offsetSize:"Offset32"  offsetRelativeTo:"Parent"` // Offset from the start of the format 14 subtable to Default UVS Table. May be 0.
	NonDefaultUVS UVSMappingTable `offsetSize:"Offset32"  offsetRelativeTo:"Parent"` // Offset from the start of the format 14 subtable to Non-Default UVS Table. May be 0.
}

// DefaultUVSTable is used in Cmap format 14
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap#default-uvs-table
type DefaultUVSTable struct {
	Ranges []UnicodeRange `arrayCount:"FirstUint32"`
}

type UnicodeRange struct {
	StartUnicodeValue [3]byte // uint24 First value in this range
	AdditionalCount   uint8   // Number of additional values in this range
}

// UVSMappingTable is used in Cmap format 14
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap#non-default-uvs-table
type UVSMappingTable struct {
	Ranges []UvsMappingRecord `arrayCount:"FirstUint32"`
}

type UvsMappingRecord struct {
	UnicodeValue [3]byte // uint24 Base Unicode value of the UVS
	GlyphID      GlyphID //	Glyph ID of the UVS
}
