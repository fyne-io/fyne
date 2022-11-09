// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lzw implements the Lempel-Ziv-Welch compressed data format,
// described in T. A. Welch, “A Technique for High-Performance Data
// Compression”, Computer, 17(6) (June 1984), pp 8-19.
//
// In particular, it implements LZW as used by the TIFF file format, including
// an "off by one" algorithmic difference when compared to standard LZW.
package lzw // import "golang.org/x/image/tiff/lzw"

/*
This file was branched from src/pkg/compress/lzw/reader.go in the
standard library. Differences from the original are marked with "NOTE".

The tif_lzw.c file in the libtiff C library has this comment:

----
The 5.0 spec describes a different algorithm than Aldus
implements.  Specifically, Aldus does code length transitions
one code earlier than should be done (for real LZW).
Earlier versions of this library implemented the correct
LZW algorithm, but emitted codes in a bit order opposite
to the TIFF spec.  Thus, to maintain compatibility w/ Aldus
we interpret MSB-LSB ordered codes to be images written w/
old versions of this library, but otherwise adhere to the
Aldus "off by one" algorithm.
----

The Go code doesn't read (invalid) TIFF files written by old versions of
libtiff, but the LZW algorithm in this package still differs from the one in
Go's standard package library to accommodate this "off by one" in valid TIFFs.
*/

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// Order specifies the bit ordering in an LZW data stream.
type Order int

const (
	// LSB means Least Significant Bits first, as used in the GIF file format.
	LSB Order = iota
	// MSB means Most Significant Bits first, as used in the TIFF and PDF
	// file formats.
	MSB
)

const (
	maxWidth           = 12
	decoderInvalidCode = 0xffff
	flushBuffer        = 1 << maxWidth
)

// decoder is the state from which the readXxx method converts a byte
// stream into a code stream.
type decoder struct {
	r        io.ByteReader
	bits     uint32
	nBits    uint
	width    uint
	read     func(*decoder) (uint16, error) // readLSB or readMSB
	litWidth int                            // width in bits of literal codes
	err      error

	// The first 1<<litWidth codes are literal codes.
	// The next two codes mean clear and EOF.
	// Other valid codes are in the range [lo, hi] where lo := clear + 2,
	// with the upper bound incrementing on each code seen.
	// overflow is the code at which hi overflows the code width. NOTE: TIFF's LZW is "off by one".
	// last is the most recently seen code, or decoderInvalidCode.
	clear, eof, hi, overflow, last uint16

	// Each code c in [lo, hi] expands to two or more bytes. For c != hi:
	//   suffix[c] is the last of these bytes.
	//   prefix[c] is the code for all but the last byte.
	//   This code can either be a literal code or another code in [lo, c).
	// The c == hi case is a special case.
	suffix [1 << maxWidth]uint8
	prefix [1 << maxWidth]uint16

	// output is the temporary output buffer.
	// Literal codes are accumulated from the start of the buffer.
	// Non-literal codes decode to a sequence of suffixes that are first
	// written right-to-left from the end of the buffer before being copied
	// to the start of the buffer.
	// It is flushed when it contains >= 1<<maxWidth bytes,
	// so that there is always room to decode an entire code.
	output [2 * 1 << maxWidth]byte
	o      int    // write index into output
	toRead []byte // bytes to return from Read
}

// readLSB returns the next code for "Least Significant Bits first" data.
func (d *decoder) readLSB() (uint16, error) {
	for d.nBits < d.width {
		x, err := d.r.ReadByte()
		if err != nil {
			return 0, err
		}
		d.bits |= uint32(x) << d.nBits
		d.nBits += 8
	}
	code := uint16(d.bits & (1<<d.width - 1))
	d.bits >>= d.width
	d.nBits -= d.width
	return code, nil
}

// readMSB returns the next code for "Most Significant Bits first" data.
func (d *decoder) readMSB() (uint16, error) {
	for d.nBits < d.width {
		x, err := d.r.ReadByte()
		if err != nil {
			return 0, err
		}
		d.bits |= uint32(x) << (24 - d.nBits)
		d.nBits += 8
	}
	code := uint16(d.bits >> (32 - d.width))
	d.bits <<= d.width
	d.nBits -= d.width
	return code, nil
}

func (d *decoder) Read(b []byte) (int, error) {
	for {
		if len(d.toRead) > 0 {
			n := copy(b, d.toRead)
			d.toRead = d.toRead[n:]
			return n, nil
		}
		if d.err != nil {
			return 0, d.err
		}
		d.decode()
	}
}

// decode decompresses bytes from r and leaves them in d.toRead.
// read specifies how to decode bytes into codes.
// litWidth is the width in bits of literal codes.
func (d *decoder) decode() {
	// Loop over the code stream, converting codes into decompressed bytes.
loop:
	for {
		code, err := d.read(d)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			d.err = err
			break
		}
		switch {
		case code < d.clear:
			// We have a literal code.
			d.output[d.o] = uint8(code)
			d.o++
			if d.last != decoderInvalidCode {
				// Save what the hi code expands to.
				d.suffix[d.hi] = uint8(code)
				d.prefix[d.hi] = d.last
			}
		case code == d.clear:
			d.width = 1 + uint(d.litWidth)
			d.hi = d.eof
			d.overflow = 1 << d.width
			d.last = decoderInvalidCode
			continue
		case code == d.eof:
			d.err = io.EOF
			break loop
		case code <= d.hi:
			c, i := code, len(d.output)-1
			if code == d.hi && d.last != decoderInvalidCode {
				// code == hi is a special case which expands to the last expansion
				// followed by the head of the last expansion. To find the head, we walk
				// the prefix chain until we find a literal code.
				c = d.last
				for c >= d.clear {
					c = d.prefix[c]
				}
				d.output[i] = uint8(c)
				i--
				c = d.last
			}
			// Copy the suffix chain into output and then write that to w.
			for c >= d.clear {
				d.output[i] = d.suffix[c]
				i--
				c = d.prefix[c]
			}
			d.output[i] = uint8(c)
			d.o += copy(d.output[d.o:], d.output[i:])
			if d.last != decoderInvalidCode {
				// Save what the hi code expands to.
				d.suffix[d.hi] = uint8(c)
				d.prefix[d.hi] = d.last
			}
		default:
			d.err = errors.New("lzw: invalid code")
			break loop
		}
		d.last, d.hi = code, d.hi+1
		if d.hi+1 >= d.overflow { // NOTE: the "+1" is where TIFF's LZW differs from the standard algorithm.
			if d.width == maxWidth {
				d.last = decoderInvalidCode
			} else {
				d.width++
				d.overflow <<= 1
			}
		}
		if d.o >= flushBuffer {
			break
		}
	}
	// Flush pending output.
	d.toRead = d.output[:d.o]
	d.o = 0
}

var errClosed = errors.New("lzw: reader/writer is closed")

func (d *decoder) Close() error {
	d.err = errClosed // in case any Reads come along
	return nil
}

// NewReader creates a new io.ReadCloser.
// Reads from the returned io.ReadCloser read and decompress data from r.
// If r does not also implement io.ByteReader,
// the decompressor may read more data than necessary from r.
// It is the caller's responsibility to call Close on the ReadCloser when
// finished reading.
// The number of bits to use for literal codes, litWidth, must be in the
// range [2,8] and is typically 8. It must equal the litWidth
// used during compression.
func NewReader(r io.Reader, order Order, litWidth int) io.ReadCloser {
	d := new(decoder)
	switch order {
	case LSB:
		d.read = (*decoder).readLSB
	case MSB:
		d.read = (*decoder).readMSB
	default:
		d.err = errors.New("lzw: unknown order")
		return d
	}
	if litWidth < 2 || 8 < litWidth {
		d.err = fmt.Errorf("lzw: litWidth %d out of range", litWidth)
		return d
	}
	if br, ok := r.(io.ByteReader); ok {
		d.r = br
	} else {
		d.r = bufio.NewReader(r)
	}
	d.litWidth = litWidth
	d.width = 1 + uint(litWidth)
	d.clear = uint16(1) << uint(litWidth)
	d.eof, d.hi = d.clear+1, d.clear+1
	d.overflow = uint16(1) << d.width
	d.last = decoderInvalidCode

	return d
}
