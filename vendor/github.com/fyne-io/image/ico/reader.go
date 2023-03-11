package ico

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	bmp "github.com/jsummers/gobmp"
)

func init() {
	image.RegisterFormat("ico", "\x00\x00\x01\x00?????\x00", Decode, DecodeConfig)
}

// ---- public ----

func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	if err := d.decode(r); err != nil {
		return nil, err
	}
	return d.images[0], nil
}

func DecodeAll(r io.Reader) ([]image.Image, error) {
	var d decoder
	if err := d.decode(r); err != nil {
		return nil, err
	}
	return d.images, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var (
		d   decoder
		cfg image.Config
		err error
	)
	if err = d.decodeHeader(r); err != nil {
		return cfg, err
	}
	if err = d.decodeEntries(r); err != nil {
		return cfg, err
	}
	e := d.entries[0]
	buf := make([]byte, e.Size+14)
	n, err := io.ReadFull(r, buf[14:])
	if err != nil && err != io.ErrUnexpectedEOF {
		return cfg, err
	}
	buf = buf[:14+n]
	if n > 8 && bytes.Equal(buf[14:22], pngHeader) {
		return png.DecodeConfig(bytes.NewReader(buf[14:]))
	}

	d.forgeBMPHead(buf, &e)
	return bmp.DecodeConfig(bytes.NewReader(buf))
}

// ---- private ----

type direntry struct {
	Width   byte
	Height  byte
	Palette byte
	_       byte
	Plane   uint16
	Bits    uint16
	Size    uint32
	Offset  uint32
}

type head struct {
	Zero   uint16
	Type   uint16
	Number uint16
}

type decoder struct {
	head    head
	entries []direntry
	images  []image.Image
}

func (d *decoder) decode(r io.Reader) (err error) {
	if err = d.decodeHeader(r); err != nil {
		return err
	}
	if err = d.decodeEntries(r); err != nil {
		return err
	}
	d.images = make([]image.Image, d.head.Number)
	for i := range d.entries {
		e := &(d.entries[i])
		data := make([]byte, e.Size+14)
		n, err := io.ReadFull(r, data[14:])
		if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}
		data = data[:14+n]
		if n > 8 && bytes.Equal(data[14:22], pngHeader) { // decode as PNG
			if d.images[i], err = png.Decode(bytes.NewReader(data[14:])); err != nil {
				return err
			}
		} else { // decode as BMP
			maskData := d.forgeBMPHead(data, e)
			if maskData != nil {
				data = data[:n+14-len(maskData)]
			}
			if d.images[i], err = bmp.Decode(bytes.NewReader(data)); err != nil {
				return err
			}
			bounds := d.images[i].Bounds()
			mask := image.NewAlpha(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			masked := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			for row := 0; row < int(e.Height); row++ {
				for col := 0; col < int(e.Width); col++ {
					if maskData != nil {
						rowSize := (int(e.Width) + 31) / 32 * 4
						if (maskData[row*rowSize+col/8]>>(7-uint(col)%8))&0x01 != 1 {
							mask.SetAlpha(col, int(e.Height)-row-1, color.Alpha{255})
						}
					} else { // 32-Bit
						rowSize := (int(e.Width)*32 + 31) / 32 * 4
						offset := int(binary.LittleEndian.Uint32(data[10:14]))
						mask.SetAlpha(col, int(e.Height)-row-1, color.Alpha{data[offset+row*rowSize+col*4+3]})
					}
				}
			}
			draw.DrawMask(masked, masked.Bounds(), d.images[i], bounds.Min, mask, bounds.Min, draw.Src)
			d.images[i] = masked
		}
	}
	return nil
}

func (d *decoder) decodeHeader(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &(d.head))
	if d.head.Zero != 0 || d.head.Type != 1 {
		return fmt.Errorf("corrupted head: [%x,%x]", d.head.Zero, d.head.Type)
	}
	return nil
}

func (d *decoder) decodeEntries(r io.Reader) error {
	n := int(d.head.Number)
	d.entries = make([]direntry, n)
	for i := 0; i < n; i++ {
		if err := binary.Read(r, binary.LittleEndian, &(d.entries[i])); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) forgeBMPHead(buf []byte, e *direntry) (mask []byte) {
	// See en.wikipedia.org/wiki/BMP_file_format
	data := buf[14:]
	imageSize := len(data)
	if e.Bits != 32 {
		maskSize := (int(e.Width) + 31) / 32 * 4 * int(e.Height)
		imageSize -= maskSize
		if imageSize <= 0 {
			return
		}
		mask = data[imageSize:]
	}

	copy(buf[0:2], "\x42\x4D") // Magic number
	dibSize := binary.LittleEndian.Uint32(data[:4])
	w := binary.LittleEndian.Uint32(data[4:8])
	h := binary.LittleEndian.Uint32(data[8:12])
	if h > w {
		binary.LittleEndian.PutUint32(data[8:12], h/2)
	}

	binary.LittleEndian.PutUint32(buf[2:6], uint32(imageSize)) // File size

	// Calculate offset into image data
	numColors := binary.LittleEndian.Uint32(data[32:36])
	bits := binary.LittleEndian.Uint16(data[14:16])

	switch bits {
	case 1, 2, 4, 8:
		x := uint32(1 << bits)
		if numColors == 0 || numColors > x {
			numColors = x
		}
	default:
		numColors = 0
	}

	var numColorsSize uint32
	switch dibSize {
	case 12, 64:
		numColorsSize = numColors * 3
	default:
		numColorsSize = numColors * 4
	}
	offset := 14 + dibSize + numColorsSize
	if dibSize > 40 && int(dibSize-4) <= len(data) {
		offset += binary.LittleEndian.Uint32(data[dibSize-8 : dibSize-4])
	}
	binary.LittleEndian.PutUint32(buf[10:14], offset)
	return
}

var pngHeader = []byte{'\x89', 'P', 'N', 'G', '\r', '\n', '\x1a', '\n'}
