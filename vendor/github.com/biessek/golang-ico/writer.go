package ico

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"io"
	"image/draw"
	"bufio"
)

func Encode(w io.Writer, im image.Image) error {
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	header := head{
		0,
		1,
		1,
	}
	entry := direntry{
		Plane:  1,
		Bits:   32,
		Offset: 22,
	}

	pngbuffer := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbuffer)
	err := png.Encode(pngwriter, m)
	if err != nil {
		return err
	}
	err = pngwriter.Flush()
	if err != nil {
		return err
	}
	entry.Size = uint32(len(pngbuffer.Bytes()))

	bounds := m.Bounds()
	entry.Width = uint8(bounds.Dx())
	entry.Height = uint8(bounds.Dy())
	bb := new(bytes.Buffer)

	var e error
	e = binary.Write(bb, binary.LittleEndian, header)
	if e != nil {
		return e
	}
	e = binary.Write(bb, binary.LittleEndian, entry)
	if e != nil {
		return e
	}

	_, e = w.Write(bb.Bytes())
	if e != nil {
		return e
	}
	_, e = w.Write(pngbuffer.Bytes())
	if e != nil {
		return e
	}

	return e
}
