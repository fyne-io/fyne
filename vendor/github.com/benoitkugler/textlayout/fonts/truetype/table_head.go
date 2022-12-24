package truetype

import (
	"encoding/binary"
	"errors"
)

// TableHead contains critical information about the rest of the font.
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6head.html
// https://docs.microsoft.com/fr-fr/typography/opentype/spec/head
type TableHead struct {
	Created            longdatetime
	Updated            longdatetime
	FontRevision       uint32
	checkSumAdjustment uint32
	UnitsPerEm         uint16
	Flags              uint16
	XMin               int16
	YMin               int16
	XMax               int16
	YMax               int16
	MacStyle           uint16
	LowestRecPPEM      uint16
	FontDirection      int16
	indexToLocFormat   int16
	// glyphDataFormat    int16
}

func parseTableHead(data []byte) (out TableHead, err error) {
	const headerSize = 54
	if len(data) < headerSize {
		return TableHead{}, errors.New("invalid 'head' table (EOF)")
	}
	// out.VersionMajor = binary.BigEndian.Uint16(data)
	// out.VersionMinor = binary.BigEndian.Uint16(data[2:])
	out.FontRevision = binary.BigEndian.Uint32(data[4:])
	out.checkSumAdjustment = binary.BigEndian.Uint32(data[8:])
	// out.MagicNumber = binary.BigEndian.Uint32(data[12:])
	out.Flags = binary.BigEndian.Uint16(data[16:])
	out.UnitsPerEm = binary.BigEndian.Uint16(data[18:])
	out.Created = longdatetime{binary.BigEndian.Uint64(data[20:])}
	out.Updated = longdatetime{binary.BigEndian.Uint64(data[28:])}
	out.XMin = int16(binary.BigEndian.Uint16(data[36:]))
	out.YMin = int16(binary.BigEndian.Uint16(data[38:]))
	out.XMax = int16(binary.BigEndian.Uint16(data[40:]))
	out.YMax = int16(binary.BigEndian.Uint16(data[42:]))
	out.MacStyle = binary.BigEndian.Uint16(data[44:])
	out.LowestRecPPEM = binary.BigEndian.Uint16(data[46:])
	out.FontDirection = int16(binary.BigEndian.Uint16(data[48:]))
	out.indexToLocFormat = int16(binary.BigEndian.Uint16(data[50:]))
	// out.glyphDataFormat = int16(binary.BigEndian.Uint16(data[52:]))
	return out, err
}

// ExpectedChecksum is the checksum that the file should have had.
func (table *TableHead) ExpectedChecksum() uint32 {
	return 0xB1B0AFBA - table.checkSumAdjustment
}

func (head *TableHead) Upem() uint16 {
	if head.UnitsPerEm < 16 || head.UnitsPerEm > 16384 {
		return 1000
	}
	return head.UnitsPerEm
}

// type TableHead struct {
// 	VersionNumber      fixed
// 	FontRevision       uint32
// 	CheckSumAdjustment uint32
// 	MagicNumber        uint32
// 	Flags              uint16
// 	UnitsPerEm         uint16
// 	Created            longdatetime
// 	Updated            longdatetime
// 	XMin               int16
// 	YMin               int16
// 	XMax               int16
// 	YMax               int16
// 	MacStyle           uint16
// 	LowestRecPPEM      uint16
// 	FontDirection      int16
// 	IndexToLocFormat   int16
// 	GlyphDataFormat    int16
// }

// func parseTableHead(buf []byte) (TableHead, error) {
// 	var fields TableHead
// 	err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &fields)
// 	return fields, err
// }
