package test_test

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestAssertCanvasTappableAt(t *testing.T) {
	c := test.NewCanvas()
	b := widget.NewButton("foo", nil)
	c.SetContent(b)
	c.Resize(fyne.NewSize(300, 300))
	b.Resize(fyne.NewSize(100, 100))
	b.Move(fyne.NewPos(100, 100))

	tt := &testing.T{}
	assert.True(t, test.AssertCanvasTappableAt(tt, c, fyne.NewPos(101, 101)), "tappable found")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertCanvasTappableAt(tt, c, fyne.NewPos(99, 99)), "tappable not found")
	assert.True(t, tt.Failed(), "test failed")
}

func TestAssertImageMatches(t *testing.T) {
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

func readImage(t *testing.T, path string) image.Image {
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()
	raw, _, err := image.Decode(file)
	require.NoError(t, err)
	img := image.NewRGBA(raw.Bounds())
	draw.Draw(img, img.Bounds(), raw, image.Pt(0, 0), draw.Src)
	return img
}
