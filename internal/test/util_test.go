package test_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/theme"
)

func TestAssertImageMatches(t *testing.T) {
	bounds := image.Rect(0, 0, 100, 34)
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(color.White), image.Point{}, draw.Src)

	txtImg := image.NewNRGBA(bounds)
	face, err := font.ParseTTF(bytes.NewReader(theme.TextFont().Content()))
	assert.Nil(t, err)

	painter.DrawString(txtImg, "Hello!", color.Black, &test.FontMap{face}, 25, 1, fyne.TextStyle{TabWidth: 4})
	draw.Draw(img, bounds, txtImg, image.Point{}, draw.Over)

	tt := &testing.T{}
	assert.False(t, test.AssertImageMatches(tt, "non_existing_master.png", img), "non existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test failed")
	assert.Equal(t, img, readImage(t, "testdata/failed/non_existing_master.png"), "image was written to disk")

	tt = &testing.T{}
	assert.True(t, test.AssertImageMatches(tt, "master.png", img), "existing master should equal given image")
	assert.False(t, tt.Failed(), "failed image match")

	tt = &testing.T{}
	assert.False(t, test.AssertImageMatches(tt, "diffing_master.png", img), "existing master should not equal given image")
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
