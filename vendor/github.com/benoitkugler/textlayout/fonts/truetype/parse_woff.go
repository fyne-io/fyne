package truetype

import (
	"encoding/binary"
	"io"

	"github.com/benoitkugler/textlayout/fonts"
)

type woffHeader struct {
	Signature      Tag
	Flavor         Tag
	Length         uint32
	NumTables      uint16
	Reserved       uint16
	TotalSfntSize  uint32
	Version        fixed
	MetaOffset     uint32
	MetaLength     uint32
	MetaOrigLength uint32
	PrivOffset     uint32
	PrivLength     uint32
}

type woffEntry struct {
	Tag          Tag
	Offset       uint32
	CompLength   uint32
	OrigLength   uint32
	OrigChecksum uint32
}

const (
	woffHeaderSize = 44
	woffEntrySize  = 20
)

func readWOFFHeader(r io.Reader) (woffHeader, error) {
	var (
		buf    [woffHeaderSize]byte
		header woffHeader
	)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return header, err
	}

	header.Signature = newTag(buf[0:4])
	header.Flavor = newTag(buf[4:8])
	header.Length = binary.BigEndian.Uint32(buf[8:12])
	header.NumTables = binary.BigEndian.Uint16(buf[12:14])
	header.Reserved = binary.BigEndian.Uint16(buf[14:16])
	header.TotalSfntSize = binary.BigEndian.Uint32(buf[16:20])
	header.Version.Major = int16(binary.BigEndian.Uint16(buf[20:22]))
	header.Version.Minor = binary.BigEndian.Uint16(buf[22:24])
	header.MetaOffset = binary.BigEndian.Uint32(buf[24:28])
	header.MetaLength = binary.BigEndian.Uint32(buf[28:32])
	header.MetaOrigLength = binary.BigEndian.Uint32(buf[32:36])
	header.PrivOffset = binary.BigEndian.Uint32(buf[36:40])
	header.PrivLength = binary.BigEndian.Uint32(buf[40:44])
	return header, nil
}

func readWOFFEntry(r io.Reader) (woffEntry, error) {
	var (
		buf   [woffEntrySize]byte
		entry woffEntry
	)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return entry, err
	}
	entry.Tag = newTag(buf[0:4])
	entry.Offset = binary.BigEndian.Uint32(buf[4:8])
	entry.CompLength = binary.BigEndian.Uint32(buf[8:12])
	entry.OrigLength = binary.BigEndian.Uint32(buf[12:16])
	entry.OrigChecksum = binary.BigEndian.Uint32(buf[16:20])
	return entry, nil
}

// `offset` is the beginning of the ressource in the file (non zero for collections)
// `relativeOffset` is true when the table offset are expresed relatively ot the ressource
// (that is, `offset`) rather than to the file
func parseWOFF(file fonts.Resource, offset uint32, relativeOffset bool) (*FontParser, error) {
	_, err := file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	header, err := readWOFFHeader(file)
	if err != nil {
		return nil, err
	}

	fontParser := &FontParser{
		file:   file,
		tables: make(map[Tag]tableSection, header.NumTables),
		Type:   header.Flavor,
	}
	for i := 0; i < int(header.NumTables); i++ {
		entry, err := readWOFFEntry(file)
		if err != nil {
			return nil, err
		}

		// TODO Check the checksum.

		if _, found := fontParser.tables[entry.Tag]; found {
			// ignore duplicate tables – the first one wins
			continue
		}

		sec := tableSection{
			offset:  entry.Offset,
			length:  entry.CompLength,
			zLength: entry.OrigLength,
		}
		// adapt the relative offsets
		if relativeOffset {
			sec.offset += offset
			if sec.offset < offset { // check for overflow
				return nil, errUnsupportedTableOffsetLength
			}
		}

		fontParser.tables[entry.Tag] = sec
	}

	return fontParser, nil
}
