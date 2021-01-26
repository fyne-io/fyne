package painter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestIsFileSVG(t *testing.T) {
	assert.True(t, isFileSVG("test.svg"))
	assert.True(t, isFileSVG("test.SVG"))
	assert.False(t, isFileSVG("testsvg"))
	assert.False(t, isFileSVG("test.png"))
}

func TestIsResourceSVG(t *testing.T) {
	res, err := fyne.LoadResourceFromPath("./testdata/stroke.svg")
	assert.Nil(t, err)
	assert.True(t, isResourceSVG(res))

	res.(*fyne.StaticResource).StaticName = "stroke"
	assert.True(t, isResourceSVG(res))

	svgString := "<svg version=\"1.1\"></svg>"
	res.(*fyne.StaticResource).StaticContent = []byte(svgString)
	assert.True(t, isResourceSVG(res))

	res.(*fyne.StaticResource).StaticContent = []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" + svgString)
	assert.True(t, isResourceSVG(res))
}
