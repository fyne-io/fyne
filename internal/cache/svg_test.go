package cache

import (
	"image"
	"image/color"
	"sync"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestColorizedSvgCache(t *testing.T) {
	testClearAll()
	tm := &timeMock{}
	tm.setTime(10, 10)
	src := []byte("<svg xmlns=\"http://www.w3.org/2000/svg\"/>")
	red := color.NRGBA{R: 0xff, A: 0xff}
	SetColorizedSvg(src, red, []byte("red"))

	got, ok := GetColorizedSvg(src, red)
	assert.True(t, ok)
	assert.Equal(t, []byte("red"), got)
	_, ok = GetColorizedSvg(src, color.NRGBA{B: 0xff, A: 0xff})
	assert.False(t, ok)
	_, ok = GetColorizedSvg([]byte("<svg/>"), red)
	assert.False(t, ok)

	lastClean = tm.createTime(10, 10)
	tm.setTime(10, 50)
	Clean(false)
	_, ok = GetColorizedSvg(src, red)
	assert.True(t, ok)

	tm.setTime(11, 30) // only alive because of the last Get
	Clean(false)
	_, ok = GetColorizedSvg(src, red)
	assert.True(t, ok)

	tm.setTime(12, 40)
	Clean(false)
	_, ok = GetColorizedSvg(src, red)
	assert.False(t, ok)
}

func TestColorizedSvgCache_Concurrent(t *testing.T) {
	testClearAll()
	src := []byte("<svg/>")
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := 0; a < 256; a++ {
				clr := color.NRGBA{A: uint8(a)}
				if _, ok := GetColorizedSvg(src, clr); !ok {
					SetColorizedSvg(src, clr, src)
				}
			}
		}()
	}
	wg.Wait()
	assert.Len(t, colorizedSvgs, 256)
}

func TestSvgCacheGet(t *testing.T) {
	ResetThemeCaches()
	img := addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	assert.Equal(t, 1, svgs.Len())

	newImg := GetSvg("empty.svg", nil, 25, 25)
	assert.Equal(t, img, newImg)

	miss := GetSvg("missing.svg", nil, 25, 25)
	assert.Nil(t, miss)
	miss = GetSvg("empty.svg", nil, 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheGet_File(t *testing.T) {
	ResetThemeCaches()
	img := addFileToCache("testdata/stroke.svg", 25, 25)
	assert.Equal(t, 1, svgs.Len())

	newImg := GetSvg("testdata/stroke.svg", nil, 25, 25)
	assert.Equal(t, img, newImg)

	miss := GetSvg("missing.svg", nil, 25, 25)
	assert.Nil(t, miss)
	miss = GetSvg("testdata/stroke.svg", nil, 30, 30)
	assert.Nil(t, miss)
}

func TestSvgCacheReset(t *testing.T) {
	ResetThemeCaches()
	_ = addToCache("empty.svg", "<svg xmlns=\"http://www.w3.org/2000/svg\"/>", 25, 25)
	SetColorizedSvg([]byte("<svg/>"), color.NRGBA{A: 0xff}, []byte("<svg/>"))
	assert.Equal(t, 1, svgs.Len())
	assert.Len(t, colorizedSvgs, 1)

	ResetThemeCaches()
	assert.Equal(t, 0, svgs.Len())
	assert.Empty(t, colorizedSvgs)
}

func addFileToCache(path string, w, h int) image.Image {
	tex := image.NewNRGBA(image.Rect(0, 0, w, h))
	SetSvg(path, nil, tex, w, h)
	return tex
}

func addToCache(name, content string, w, h int) image.Image {
	resource := fyne.NewStaticResource(name, []byte(content))
	tex := image.NewNRGBA(image.Rect(0, 0, w, h))
	SetSvg(resource.Name(), nil, tex, w, h)
	return tex
}
