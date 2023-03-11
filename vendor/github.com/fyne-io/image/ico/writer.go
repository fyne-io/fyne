package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"image/png"
	"io"
)

type icondir struct {
	reserved  uint16 //lint:ignore U1000 in spec but not used
	imageType uint16
	numImages uint16
}

type icondirentry struct {
	imageWidth   uint8
	imageHeight  uint8
	numColors    uint8 //lint:ignore U1000 in spec but not used
	reserved     uint8 //lint:ignore U1000 in spec but not used
	colorPlanes  uint16
	bitsPerPixel uint16
	sizeInBytes  uint32
	offset       uint32
}

func newIcondir() icondir {
	var id icondir
	id.imageType = 1
	id.numImages = 1
	return id
}

func newIcondirentry() icondirentry {
	var ide icondirentry
	ide.colorPlanes = 1   // windows is supposed to not mind 0 or 1, but other icon files seem to have 1 here
	ide.bitsPerPixel = 32 // can be 24 for bitmap or 24/32 for png. Set to 32 for now
	ide.offset = 22       //6 icondir + 16 icondirentry, next image will be this image size + 16 icondirentry, etc
	return ide
}

func Encode(w io.Writer, im image.Image) error {
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	id := newIcondir()
	ide := newIcondirentry()

	pngbb := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbb)
	png.Encode(pngwriter, m)
	pngwriter.Flush()
	ide.sizeInBytes = uint32(len(pngbb.Bytes()))

	bounds := m.Bounds()
	ide.imageWidth = uint8(bounds.Dx())
	ide.imageHeight = uint8(bounds.Dy())
	bb := new(bytes.Buffer)

	var e error
	binary.Write(bb, binary.LittleEndian, id)
	binary.Write(bb, binary.LittleEndian, ide)

	w.Write(bb.Bytes())
	w.Write(pngbb.Bytes())

	return e
}
