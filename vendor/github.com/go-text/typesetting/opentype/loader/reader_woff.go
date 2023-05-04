// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package loader

import (
	"encoding/binary"
	"errors"
	"io"
)

type woffEntry struct {
	Tag          Tag
	Offset       uint32
	CompLength   uint32
	OrigLength   uint32
	OrigChecksum uint32
}

const (
	woffHeaderSize = 44 // for the full header, but we only read Flavor and NumTables
	woffEntrySize  = 20
)

func readWOFFHeader(r io.Reader) (flavor Tag, numTables uint16, err error) {
	var buf [woffHeaderSize]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, 0, err
	}

	return NewTag(buf[4], buf[5], buf[6], buf[7]), binary.BigEndian.Uint16(buf[12:14]), nil
}

func readWOFFEntry(r io.Reader) (woffEntry, error) {
	var (
		buf   [woffEntrySize]byte
		entry woffEntry
	)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return entry, err
	}
	entry.Tag = NewTag(buf[0], buf[1], buf[2], buf[3])
	entry.Offset = binary.BigEndian.Uint32(buf[4:8])
	entry.CompLength = binary.BigEndian.Uint32(buf[8:12])
	entry.OrigLength = binary.BigEndian.Uint32(buf[12:16])
	entry.OrigChecksum = binary.BigEndian.Uint32(buf[16:20])
	return entry, nil
}

// `offset` is the beginning of the ressource in the file (non zero for collections)
// `relativeOffset` is true when the table offset are expresed relatively ot the ressource
// (that is, `offset`) rather than to the file
func parseWOFF(file Resource, offset uint32, relativeOffset bool) (*Loader, error) {
	_, err := file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	flavor, numTables, err := readWOFFHeader(file)
	if err != nil {
		return nil, err
	}

	fontParser := &Loader{
		file:   file,
		tables: make(map[Tag]tableSection, numTables),
		Type:   flavor,
	}
	for i := 0; i < int(numTables); i++ {
		entry, err := readWOFFEntry(file)
		if err != nil {
			return nil, err
		}

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
				return nil, errors.New("unsupported table offset or length")
			}
		}

		fontParser.tables[entry.Tag] = sec
	}

	return fontParser, nil
}
