package svg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestIsFileSVG(t *testing.T) {
	assert.True(t, IsFileSVG("test.svg"))
	assert.True(t, IsFileSVG("test.SVG"))
	assert.False(t, IsFileSVG("testsvg"))
	assert.False(t, IsFileSVG("test.png"))
}

func TestIsResourceSVG(t *testing.T) {
	res, err := fyne.LoadResourceFromPath("./testdata/circles.svg")
	assert.Nil(t, err)
	assert.True(t, IsResourceSVG(res))

	res.(*fyne.StaticResource).StaticName = "stroke"
	assert.True(t, IsResourceSVG(res))

	svgString := "<svg version=\"1.1\"></svg>"
	res.(*fyne.StaticResource).StaticContent = []byte(svgString)
	assert.True(t, IsResourceSVG(res))

	res.(*fyne.StaticResource).StaticContent = []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" + svgString)
	assert.True(t, IsResourceSVG(res))
}
