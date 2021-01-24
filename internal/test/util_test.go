package test_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font"

	"fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestAssertImageMatches(t *testing.T) {
	bounds := image.Rect(0, 0, 100, 50)
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(color.White), image.Point{}, draw.Src)

	txtImg := image.NewNRGBA(bounds)
	opts := truetype.Options{Size: 20, DPI: 96}
	f, _ := truetype.Parse(theme.TextFont().Content())
	face := truetype.NewFace(f, &opts)
	d := font.Drawer{
		Dst:  txtImg,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  freetype.Pt(0, 50-face.Metrics().Descent.Ceil()),
	}
	d.DrawString("Hello!")
	draw.Draw(img, bounds, txtImg, image.Point{}, draw.Over)

	tt := &testing.T{}
	assert.False(t, test.AssertImageMatches(tt, "non_existing_master.png", img), "non existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test failed")
	assert.Equal(t, img, readImage(t, "testdata/failed/non_existing_master.png"), "image was written to disk")

	tt = &testing.T{}
	assert.True(t, test.AssertImageMatches(tt, "master.png", img), "existing master is equal a given image")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertImageMatches(tt, "diffing_master.png", img), "existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test did not fail")
	assert.Equal(t, img, readImage(t, "testdata/failed/diffing_master.png"), "image was written to disk")

	if !t.Failed() {
		os.RemoveAll("testdata/failed")
	}
}

func TestNewCheckedImage(t *testing.T) {
	img := test.NewCheckedImage(10, 10, 5, 2)
	expectedColorValues := [][]uint8{
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
		{0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0},
	}
	for y, xv := range expectedColorValues {
		for x, v := range xv {
			assert.Equal(t, color.NRGBA{R: v, G: v, B: v, A: 0xff}, img.At(x, y), fmt.Sprintf("color value at %d,%d", x, y))
		}
	}
}

func readImage(t *testing.T, path string) image.Image {
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()
	raw, _, err := image.Decode(file)
	require.NoError(t, err)
	img := image.NewNRGBA(raw.Bounds())
	draw.Draw(img, img.Bounds(), raw, image.Pt(0, 0), draw.Src)
	return img
}
