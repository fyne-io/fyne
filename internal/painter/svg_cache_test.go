package painter

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

func TestSvgCacheGet(t *testing.T) {
	SvgCacheReset()
	img := addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, len(rasters))

	newImg := svgCacheGet("empty.svg", 25, 25)
	assert.Equal(t, img, newImg)

	miss := svgCacheGet("missing.svg", 25, 25)
	assert.Nil(t, miss)
	miss = svgCacheGet("empty.svg", 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheGet_File(t *testing.T) {
	SvgCacheReset()
	img := addFileToCache("testdata/stroke.svg", 25, 25)
	assert.Equal(t, 1, len(rasters))

	newImg := svgCacheGet("testdata/stroke.svg", 25, 25)
	assert.Equal(t, img, newImg)

	miss := svgCacheGet("missing.svg", 25, 25)
	assert.Nil(t, miss)
	miss = svgCacheGet("testdata/stroke.svg", 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheReset(t *testing.T) {
	SvgCacheReset()
	_ = addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, len(rasters))

	SvgCacheReset()
	assert.Equal(t, 0, len(rasters))
}

func addFileToCache(path string, w, h int) image.Image {
	img := canvas.NewImageFromFile(path)
	return PaintImage(img, nil, w, h)
}

func addToCache(name, content string, w, h int) image.Image {
	img := canvas.NewImageFromResource(fyne.NewStaticResource(name, []byte(content)))
	return PaintImage(img, nil, w, h)
}
