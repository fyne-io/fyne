package test_test

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"testing"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font"

	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestAssertImageEqualsMaster(t *testing.T) {
	bounds := image.Rect(0, 0, 100, 50)
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(color.White), image.ZP, draw.Src)

	txtImg := image.NewRGBA(bounds)
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
	draw.Draw(img, bounds, txtImg, image.ZP, draw.Over)

	tt := &testing.T{}
	assert.False(t, test.AssertImageEqualsMaster(tt, "non_existing_master.png", img), "non existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test failed")
	assert.Equal(t, img, readImage(t, "failed/non_existing_master.png"), "image was written to disk")

	tt = &testing.T{}
	assert.True(t, test.AssertImageEqualsMaster(tt, "master.png", img), "existing master is equal a given image")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertImageEqualsMaster(tt, "diffing_master.png", img), "existing master is not equal a given image")
	assert.True(t, tt.Failed(), "test did not fail")
	assert.Equal(t, img, readImage(t, "failed/diffing_master.png"), "image was written to disk")
}

func readImage(t *testing.T, name string) image.Image {
	file, err := os.Open(filepath.Join("testdata", name))
	require.NoError(t, err)
	defer file.Close()
	raw, _, err := image.Decode(file)
	require.NoError(t, err)
	img := image.NewRGBA(raw.Bounds())
	draw.Draw(img, img.Bounds(), raw, image.Pt(0, 0), draw.Src)
	return img
}
