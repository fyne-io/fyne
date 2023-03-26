// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package loader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// An Entry in an OpenType table.
type otfEntry struct {
	Tag      Tag
	CheckSum uint32
	Offset   uint32
	Length   uint32
}

func readOTFHeader(r io.Reader) (flavor Tag, numTables uint16, err error) {
	var buf [12]byte
	if _, err := r.Read(buf[:]); err != nil {
		return 0, 0, fmt.Errorf("invalid OpenType header: %s", err)
	}

	return NewTag(buf[0], buf[1], buf[2], buf[3]), binary.BigEndian.Uint16(buf[4:6]), nil
}

func readOTFEntry(r io.Reader) (otfEntry, error) {
	var (
		buf   [16]byte
		entry otfEntry
	)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return entry, fmt.Errorf("invalid directory entry: %s", err)
	}

	entry.Tag = Tag(binary.BigEndian.Uint32(buf[0:4]))
	entry.CheckSum = binary.BigEndian.Uint32(buf[4:8])
	entry.Offset = binary.BigEndian.Uint32(buf[8:12])
	entry.Length = binary.BigEndian.Uint32(buf[12:16])

	return entry, nil
}

// parseOTF reads an OpenTyp (.otf) or TrueType (.ttf) file and returns a Font.
// If the parsing fails, then an error is returned and Font will be nil.
// `offset` is the beginning of the ressource in the file (non zero for collections)
// `relativeOffset` is true when the table offset are expresed relatively to the ressource start
// (that is, `offset`) rather than to the file start.
func parseOTF(file Resource, offset uint32, relativeOffset bool) (*Loader, error) {
	_, err := file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("invalid offset: %s", err)
	}

	flavor, numTables, err := readOTFHeader(file)
	if err != nil {
		return nil, err
	}

	pr := &Loader{
		file:   file,
		tables: make(map[Tag]tableSection, numTables),
		Type:   flavor,
	}

	for i := 0; i < int(numTables); i++ {
		entry, err := readOTFEntry(file)
		if err != nil {
			return nil, err
		}

		if _, found := pr.tables[entry.Tag]; found {
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
				return nil, errors.New("unsupported table offset or length")
			}
		}
		pr.tables[entry.Tag] = sec
	}

	return pr, nil
}
