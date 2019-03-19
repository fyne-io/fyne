package icns

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

// TestDecode relies on Encode being correct.
// We are testing that an ICNS with a series of icons will only yield the
// largest icon in the series.
func TestDecode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc  string
		input image.Image
		want  image.Image
	}{
		{
			"valid square icon, exact size",
			rect(0, 0, 256, 256),
			rect(0, 0, 256, 256),
		},
		{
			"non exact size",
			rect(0, 0, 50, 50),
			rect(0, 0, 32, 32),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(st *testing.T) {
			buf := bytes.NewBuffer(nil)
			if err := Encode(buf, tt.input); err != nil {
				st.Fatalf("unexpected error while encoding: %v", err)
			}
			img, err := Decode(buf)
			if err != nil {
				st.Fatalf("unexpected error: %v", err)
			}
			if tt.want != nil && !imageCompare(img, tt.want) {
				st.Fatalf("decoded image is incorrect")
			}
		})
	}
}

func imageCompare(left, right image.Image) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil && right != nil {
		return false
	}
	if left != nil && right == nil {
		return false
	}
	lb := left.Bounds()
	for ii := lb.Min.X; ii <= lb.Max.X; ii++ {
		for kk := lb.Min.Y; kk <= lb.Max.Y; kk++ {
			lr, lg, lb, la := left.At(ii, kk).RGBA()
			rr, rg, rb, ra := right.At(ii, kk).RGBA()
			if lr != rr || lg != rg || lb != rb || la != ra {
				return false
			}
		}
	}
	return true
}

// TestEncode tests for input validation, sanity checks and errors.
// The validity of the encoding is not tested here.
// Super large images are not tested because the resizing takes too
// long for unit testing.
func TestEncode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc string
		wr   io.Writer
		img  image.Image

		wantErr bool
	}{
		{
			"nil image",
			ioutil.Discard,
			nil,
			true,
		},
		{
			"nil writer",
			nil,
			rect(0, 0, 50, 50),
			true,
		},
		{
			"valid sqaure",
			ioutil.Discard,
			rect(0, 0, 50, 50),
			false,
		},
		{
			"valid non-square",
			ioutil.Discard,
			rect(0, 0, 10, 50),
			false,
		},
		{
			"valid non-square, weird dimensions",
			ioutil.Discard,
			rect(0, 0, 17, 77),
			false,
		},
		{
			"invalid zero img",
			ioutil.Discard,
			rect(0, 0, 0, 0),
			true,
		},
		{
			"invalid small img",
			ioutil.Discard,
			rect(0, 0, 1, 1),
			true,
		},
		{
			"valid square not at origin point",
			ioutil.Discard,
			rect(10, 10, 50, 50),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(st *testing.T) {
			err := Encode(tt.wr, tt.img)
			if !tt.wantErr && err != nil {
				st.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSizesFromMax(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc string
		from uint
		want []uint
	}{
		{
			"small",
			100,
			[]uint{64, 32},
		},
		{
			"large",
			99999,
			[]uint{1024, 512, 256, 64, 32},
		},
		{
			"smallest",
			0,
			[]uint{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(st *testing.T) {
			got := sizesFrom(tt.from)
			if !reflect.DeepEqual(got, tt.want) {
				st.Errorf("want=%d, got=%d", tt.want, got)
			}
		})
	}
}

func TestBiggestSide(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc string
		img  image.Image
		want uint
	}{
		{
			"equal",
			rect(0, 0, 100, 100),
			100,
		},
		{
			"right larger",
			rect(0, 0, 50, 100),
			100,
		},
		{
			"left larger",
			rect(0, 0, 100, 50),
			100,
		},
		{
			"off by one",
			rect(0, 0, 100, 99),
			100,
		},
		{
			"empty",
			rect(0, 0, 0, 0),
			0,
		},
		{
			"left empty",
			rect(0, 0, 0, 10),
			10,
		},
		{
			"right empty",
			rect(0, 0, 10, 0),
			10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(st *testing.T) {
			got := biggestSide(tt.img)
			if got != tt.want {
				st.Errorf("want=%d, got=%d", tt.want, got)
			}
		})
	}
}

func TestFindNearestSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc string
		img  image.Image
		want uint
	}{
		{
			"small",
			rect(0, 0, 100, 100),
			64,
		},
		{
			"very large",
			rect(0, 0, 123456789, 123456789),
			1024,
		},
		{
			"too small",
			rect(0, 0, 16, 16),
			0,
		},
		{
			"off by one",
			rect(0, 0, 33, 33),
			32,
		},
		{
			"exact",
			rect(0, 0, 256, 256),
			256,
		},
		{
			"exact",
			rect(0, 0, 1024, 1024),
			1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(st *testing.T) {
			got := findNearestSize(tt.img)
			if tt.want != got {
				st.Errorf("want=%d, got=%d", tt.want, got)
			}
		})
	}
}

func TestEncodeImage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc string

		img    image.Image
		format string

		want string
	}{
		{
			"png - png",
			_decode(_png(rect(0, 0, 50, 50))),
			"png",
			"png",
		},
		{
			"default png - png",
			_decode(_png(rect(0, 0, 50, 50))),
			"",
			"png",
		},
		{
			"jpg - jpg",
			_decode(_jpg(rect(0, 0, 50, 50))),
			"jpeg",
			"png",
		},
		{
			"default jpg - png",
			_decode(_jpg(rect(0, 0, 50, 50))),
			"",
			"png",
		},
		{
			"invalid format identifier",
			_decode(_jpg(rect(0, 0, 50, 50))),
			"asdf",
			"png",
		},
		{
			"not actually a jpeg",
			_decode(_png(rect(0, 0, 50, 50))),
			"jpeg",
			"png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(st *testing.T) {
			data, err := encodeImage(tt.img)
			if err != nil {
				st.Fatalf("encoding image: %v", err)
			}
			_, f, err := image.Decode(bytes.NewBuffer(data))
			if err != nil {
				st.Fatalf("decoding iamge: %v", err)
			}
			if f != tt.want {
				st.Fatalf("formats: want=%s, got=%s", tt.want, f)
			}
		})
	}
}

func rect(x0, y0, x1, y1 int) image.Image {
	return image.Rect(x0, y0, x1, y1)
}

func _png(img image.Image) io.Reader {
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		panic(errors.Wrapf(err, "encoding png"))
	}
	return buf
}

func _jpg(img image.Image) io.Reader {
	buf := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		panic(errors.Wrapf(err, "encoding jpeg"))
	}
	return buf
}

func _decode(r io.Reader) image.Image {
	m, _, err := image.Decode(r)
	if err != nil {
		panic(errors.Wrapf(err, "decoding image"))
	}
	return m
}
