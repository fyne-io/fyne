package cache

import (
	"image"
	"sync"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/stretchr/testify/assert"
)

func syncMapLen(m *sync.Map) (n int) {
	m.Range(func(_, _ interface{}) bool {
		n++
		return true
	})
	return
}

func TestSvgCacheGet(t *testing.T) {
	ResetThemeCaches()
	img := addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, syncMapLen(svgs))

	newImg := GetSvg("empty.svg", 25, 25)
	assert.Equal(t, img, newImg)

	miss := GetSvg("missing.svg", 25, 25)
	assert.Nil(t, miss)
	miss = GetSvg("empty.svg", 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheGet_File(t *testing.T) {
	ResetThemeCaches()
	img := addFileToCache("testdata/stroke.svg", 25, 25)
	assert.Equal(t, 1, syncMapLen(svgs))

	newImg := GetSvg("testdata/stroke.svg", 25, 25)
	assert.Equal(t, img, newImg)

	miss := GetSvg("missing.svg", 25, 25)
	assert.Nil(t, miss)
	miss = GetSvg("testdata/stroke.svg", 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheReset(t *testing.T) {
	ResetThemeCaches()
	_ = addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, syncMapLen(svgs))

	ResetThemeCaches()
	assert.Equal(t, 0, syncMapLen(svgs))
}

func addFileToCache(path string, w, h int) image.Image {
	img := canvas.NewImageFromFile(path)
	tex := image.NewNRGBA(image.Rect(0, 0, w, h))
	SetSvg(img.File, tex, w, h)
	return tex
}

func addToCache(name, content string, w, h int) image.Image {
	img := canvas.NewImageFromResource(fyne.NewStaticResource(name, []byte(content)))
	tex := image.NewNRGBA(image.Rect(0, 0, w, h))
	SetSvg(img.Resource.Name(), tex, w, h)
	return tex
}
