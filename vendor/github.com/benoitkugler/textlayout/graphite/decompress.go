package graphite

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const uncompressedLimit = 10000000

// Some fonts contains tables compressed with the lz4 format
// this function expect the table to begin with a version (uint32)
// and check if it then has a valid encoding scheme
func decompressTable(data []byte) ([]byte, uint16, error) {
	if len(data) < 4 {
		return nil, 0, errors.New("invalid table (EOF)")
	}
	version := uint16(binary.BigEndian.Uint32(data) >> 16) // major

	if version >= 3 && len(data) >= 8 {
		compression := binary.BigEndian.Uint32(data[4:8])
		scheme := compression >> 27
		switch scheme {
		case 0: // no compression, just go on

		case 1: // lz4 compression
			size := compression & 0x07ffffff
			if size > uncompressedLimit {
				return nil, 0, fmt.Errorf("unsupported size for uncompressed Glat table: %d", size)
			}
			uncompressed := make([]byte, size)
			err := decodeLz4Block(uncompressed, data[8:])
			if err != nil {
				return nil, 0, fmt.Errorf("invalid lz4 compressed table: %s", err)
			}
			return uncompressed, version, nil
		default:
			return nil, 0, fmt.Errorf("unsupported compression scheme: %d", scheme)
		}
	}
	return data, version, nil
}

func decodeLz4Block(dst, src []byte) (err error) {
	if len(src) == 0 {
		return nil
	}

	defer func() {
		if rc := recover(); rc != nil {
			err = fmt.Errorf("panic: %s", rc)
		}
	}()

	err = decodeBlockPanic(dst, src)
	return err
}

// copied from github.com/pierrec/lz4
func decodeBlockPanic(dst, src []byte) error {
	const minMatch = 4 // the minimum size of the match sequence size (4 bytes)

	var si, di int
	for {
		// Literals and match lengths (token).
		b := int(src[si])
		si++

		// Literals.
		if lLen := b >> 4; lLen > 0 {
			switch {
			case lLen < 0xF && si+16 < len(src):
				// Shortcut 1
				// if we have enough room in src and dst, and the literals length
				// is small enough (0..14) then copy all 16 bytes, even if not all
				// are part of the literals.
				copy(dst[di:], src[si:si+16])
				si += lLen
				di += lLen
				if mLen := b & 0xF; mLen < 0xF {
					// Shortcut 2
					// if the match length (4..18) fits within the literals, then copy
					// all 18 bytes, even if not all are part of the literals.
					mLen += 4
					if offset := int(src[si]) | int(src[si+1])<<8; mLen <= offset {
						i := di - offset
						end := i + 18
						if end > len(dst) {
							// The remaining buffer may not hold 18 bytes.
							// See https://github.com/pierrec/lz4/issues/51.
							end = len(dst)
						}
						copy(dst[di:], dst[i:end])
						si += 2
						di += mLen
						continue
					}
				}
			case lLen == 0xF:
				for src[si] == 0xFF {
					lLen += 0xFF
					si++
				}
				lLen += int(src[si])
				si++
				fallthrough
			default:
				copy(dst[di:di+lLen], src[si:si+lLen])
				si += lLen
				di += lLen
			}
		}
		if si >= len(src) {
			return nil
		}

		offset := int(src[si]) | int(src[si+1])<<8
		if offset == 0 {
			return errors.New("uncompress lz4: invalid offset")
		}
		si += 2

		// Match.
		mLen := b & 0xF
		if mLen == 0xF {
			for src[si] == 0xFF {
				mLen += 0xFF
				si++
			}
			mLen += int(src[si])
			si++
		}
		mLen += minMatch

		// Copy the match.
		expanded := dst[di-offset:]
		if mLen > offset {
			// Efficiently copy the match dst[di-offset:di] into the dst slice.
			bytesToCopy := offset * (mLen / offset)
			for n := offset; n <= bytesToCopy+offset; n *= 2 {
				copy(expanded[n:], expanded[:n])
			}
			di += bytesToCopy
			mLen -= bytesToCopy
		}
		di += copy(dst[di:di+mLen], expanded[:mLen])
	}
}
