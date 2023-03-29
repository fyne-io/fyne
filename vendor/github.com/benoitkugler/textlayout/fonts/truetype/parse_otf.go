package truetype

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/benoitkugler/textlayout/fonts"
)

type otfHeader struct {
	ScalerType    Tag
	NumTables     uint16
	SearchRange   uint16
	EntrySelector uint16
	RangeShift    uint16
}

// An Entry in an OpenType table.
type directoryEntry struct {
	Tag      Tag
	CheckSum uint32
	Offset   uint32
	Length   uint32
}

func readOTFHeader(r io.Reader, header *otfHeader) error {
	var buf [12]byte
	if _, err := r.Read(buf[:]); err != nil {
		return fmt.Errorf("invalid OpenType header: %s", err)
	}

	header.ScalerType = newTag(buf[0:4])
	header.NumTables = binary.BigEndian.Uint16(buf[4:6])
	header.SearchRange = binary.BigEndian.Uint16(buf[6:8])
	header.EntrySelector = binary.BigEndian.Uint16(buf[8:10])
	header.RangeShift = binary.BigEndian.Uint16(buf[10:12])

	return nil
}

func readDirectoryEntry(r io.Reader, entry *directoryEntry) error {
	var buf [16]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return fmt.Errorf("invalid directory entry: %s", err)
	}

	entry.Tag = newTag(buf[0:4])
	entry.CheckSum = binary.BigEndian.Uint32(buf[4:8])
	entry.Offset = binary.BigEndian.Uint32(buf[8:12])
	entry.Length = binary.BigEndian.Uint32(buf[12:16])

	return nil
}

// parseOTF reads an OpenTyp (.otf) or TrueType (.ttf) file and returns a Font.
// If parsing fails, then an error is returned and Font will be nil.
// `offset` is the beginning of the ressource in the file (non zero for collections)
// `relativeOffset` is true when the table offset are expresed relatively ot the ressource
// (that is, `offset`) rather than to the file
func parseOTF(file fonts.Resource, offset uint32, relativeOffset bool) (*FontParser, error) {
	_, err := file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("invalid offset: %s", err)
	}

	var header otfHeader
	if err := readOTFHeader(file, &header); err != nil {
		return nil, err
	}

	fontParser := &FontParser{
		file:   file,
		tables: make(map[Tag]tableSection, header.NumTables),
		Type:   header.ScalerType,
	}

	for i := 0; i < int(header.NumTables); i++ {
		var entry directoryEntry
		if err := readDirectoryEntry(file, &entry); err != nil {
			return nil, err
		}

		// TODO Check the checksum.

		if _, found := fontParser.tables[entry.Tag]; found {
			// ignore duplicate tables – the first one wins
			continue
		}

		sec := tableSection{
			offset: entry.Offset,
			length: entry.Length,
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
