package icns

import (
	"errors"
	"image"
	"io"
	"sync"

	"github.com/nfnt/resize"
)

// Encoder encodes ICNS files from a source image.
type Encoder struct {
	Wr        io.Writer
	Algorithm InterpolationFunction
}

// NewEncoder initialises an encoder.
func NewEncoder(wr io.Writer) *Encoder {
	return &Encoder{
		Wr: wr,
	}
}

// WithAlgorithm applies the interpolation function used to resize the image.
func (enc *Encoder) WithAlgorithm(a InterpolationFunction) *Encoder {
	enc.Algorithm = a
	return enc
}

// Encode icns with the given configuration.
func (enc *Encoder) Encode(img image.Image) error {
	if enc.Wr == nil {
		return errors.New("cannot write to nil writer")
	}
	if img == nil {
		return errors.New("cannot encode nil image")
	}
	iconset, err := NewIconSet(img, enc.Algorithm)
	if err != nil {
		return err
	}
	if _, err := iconset.WriteTo(enc.Wr); err != nil {
		return err
	}
	return nil
}

// Encode writes img to wr in ICNS format.
// img is assumed to be a rectangle; non-square dimensions will be squared
// without preserving the aspect ratio.
// Uses nearest neighbor as interpolation algorithm.
func Encode(wr io.Writer, img image.Image) error {
	return NewEncoder(wr).Encode(img)
}

// NewIconSet uses the source image to create an IconSet.
// If width != height, the image will be resized using the largest side without
// preserving the aspect ratio.
func NewIconSet(img image.Image, interp InterpolationFunction) (*IconSet, error) {
	biggest := findNearestSize(img)
	if biggest == 0 {
		return nil, ErrImageTooSmall{image: img, need: 16}
	}
	icons := make([]*Icon, len(osTypes))
	work := sync.WaitGroup{}
	for ii, size := range sizesFrom(biggest) {
		t, ok := getTypeFromSize(size)
		if !ok {
			continue
		}
		ii := ii
		size := size
		work.Add(1)
		go func() {
			iconImg := resize.Resize(size, size, img, interp)
			icons[ii] = &Icon{
				Type:  t,
				Image: iconImg,
			}
			work.Done()
		}()
	}
	work.Wait()
	iconset := &IconSet{
		Icons: icons,
	}
	return iconset, nil
}

// Big-endian.
// https://golang.org/src/image/png/writer.go
func writeUint32(b []uint8, u uint32) {
	b[0] = uint8(u >> 24)
	b[1] = uint8(u >> 16)
	b[2] = uint8(u >> 8)
	b[3] = uint8(u >> 0)
}

var sizes = []uint{
	1024,
	512,
	256,
	128,
	64,
	32,
}

// findNearestSize finds the biggest icon size we can use for this image.
func findNearestSize(img image.Image) uint {
	size := biggestSide(img)
	for _, s := range sizes {
		if size >= s {
			return s
		}
	}
	return 0
}

func biggestSide(img image.Image) uint {
	var size uint
	b := img.Bounds()
	w, h := uint(b.Max.X), uint(b.Max.Y)
	size = w
	if h > size {
		size = h
	}
	return size
}

// sizesFrom returns a slice containing the sizes less than and including max.
func sizesFrom(max uint) []uint {
	for ii, s := range sizes {
		if s <= max {
			return sizes[ii:len(sizes)]
		}
	}
	return []uint{}
}

// OsType is a 4 character identifier used to differentiate icon types.
type OsType struct {
	ID   string
	Size uint
}

var osTypes = []OsType{
	{ID: "ic10", Size: uint(1024)},
	{ID: "ic14", Size: uint(512)},
	{ID: "ic13", Size: uint(256)},
	{ID: "ic07", Size: uint(128)},
	{ID: "ic12", Size: uint(64)},
	{ID: "ic11", Size: uint(32)},
}

// getTypeFromSize returns the type for the given icon size (in px).
// The boolean indicates whether the type exists.
func getTypeFromSize(size uint) (OsType, bool) {
	for _, t := range osTypes {
		if t.Size == size {
			return t, true
		}
	}
	return OsType{}, false
}

func getTypeFromID(ID string) (OsType, bool) {
	for _, t := range osTypes {
		if t.ID == ID {
			return t, true
		}
	}
	return OsType{}, false
}

func osTypeFromID(ID string) OsType {
	t, _ := getTypeFromID(ID)
	return t
}
