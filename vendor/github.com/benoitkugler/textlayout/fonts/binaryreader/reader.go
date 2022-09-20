// Package binaryreader provides a convenient
// reader to decode (big endian) binary content.
package binaryreader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Reader holds a byte slice and provide convenient methods
// to interpret its content.
type Reader struct {
	data []byte
	pos  int
}

// NewReader returns a reader filled with `data`.
func NewReader(data []byte) *Reader { return &Reader{data: data} }

// NewReaderAt returns a reader filled with `data[offset:]`.
func NewReaderAt(data []byte, offset uint32) (*Reader, error) {
	if len(data) < int(offset) {
		return nil, fmt.Errorf("invalid offset %d for length %d", offset, len(data))
	}
	return NewReader(data[offset:]), nil
}

// Data returns the remaining (unread) slice.
func (r *Reader) Data() []byte {
	return r.data[r.pos:]
}

// Skip advances from `count` bytes, delaying potential bound checking.
func (r *Reader) Skip(count int) {
	r.pos += count
}

// SetPos seeks from the start, delaying potential bound checking.
func (r *Reader) SetPos(pos int) { r.pos = pos }

func (r *Reader) Read(out []byte) (int, error) {
	if len(r.data) < r.pos+len(out) {
		return 0, io.EOF
	}
	n := copy(out, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// Byte reads one byte and advances. The only error possible is
// reaching the end of the slice.
func (r *Reader) Byte() (byte, error) {
	if len(r.data) <= r.pos {
		return 0, errors.New("no byte (EOF)")
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

// Uint16 reads one uint16 and advances. The only error possible is
// reaching the end of the slice.
func (r *Reader) Uint16() (uint16, error) {
	if len(r.data) <= r.pos+1 {
		return 0, errors.New("no uint16 (EOF)")
	}
	b := binary.BigEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return b, nil
}

// Uint32 reads one uint32 and advances. The only error possible is
// reaching the end of the slice.
func (r *Reader) Uint32() (uint32, error) {
	if len(r.data) <= r.pos+3 {
		return 0, errors.New("no uint32 (EOF)")
	}
	b := binary.BigEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return b, nil
}

// Uint16s reads a slice of uint16 with length `count` and advances.
func (r *Reader) Uint16s(count int) ([]uint16, error) {
	if len(r.data) < r.pos+2*count {
		return nil, fmt.Errorf("EOF for %d uint16s", count)
	}
	out := make([]uint16, count)
	for i := range out {
		out[i] = binary.BigEndian.Uint16(r.data[r.pos+i*2:])
	}
	r.pos += 2 * count
	return out, nil
}

// Int16s reads a slice of int16 with length `count` and advances.
func (r *Reader) Int16s(count int) ([]int16, error) {
	if len(r.data) < r.pos+2*count {
		return nil, fmt.Errorf("EOF for %d int16s (%d bytes left)", count, len(r.data)-r.pos)
	}
	out := make([]int16, count)
	for i := range out {
		out[i] = int16(binary.BigEndian.Uint16(r.data[r.pos+i*2:]))
	}
	r.pos += 2 * count
	return out, nil
}

// Uint32s reads a slice of uint32 with length `count` and advances.
func (r *Reader) Uint32s(count int) ([]uint32, error) {
	if len(r.data) < r.pos+4*count {
		return nil, fmt.Errorf("EOF for %d uint32s", count)
	}
	out := make([]uint32, count)
	for i := range out {
		out[i] = binary.BigEndian.Uint32(r.data[r.pos+i*4:])
	}
	r.pos += 4 * count
	return out, nil
}

// FixedSizes return a slice with length count*size, and advances
func (r *Reader) FixedSizes(count, size int) ([]byte, error) {
	L := count * size
	if len(r.data) < r.pos+L {
		return nil, fmt.Errorf("EOF for %d records with size %d", count, size)
	}
	out := r.data[r.pos : r.pos+L]
	r.pos += L
	return out, nil
}

// ReadStruct calls binary.Read and advances. The only error possible
// is reaching the end of the slice.
func (r *Reader) ReadStruct(out interface{}) error {
	return binary.Read(r, binary.BigEndian, out)
}
