package gif

import (
	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsFileGIF(t *testing.T) {
	assert.True(t, IsFileGIF("test.gif"))
	assert.True(t, IsFileGIF("test.GIF"))
	assert.False(t, IsFileGIF("testgif"))
	assert.False(t, IsFileGIF("test.png"))
}

func TestIsResourceGIF(t *testing.T) {
	res, err := fyne.LoadResourceFromPath("../../cmd/fyne/internal/templates/data/spinner_light.gif")
	assert.Nil(t, err)
	assert.True(t, IsResourceGIF(res))

	res.(*fyne.StaticResource).StaticName = "stroke"
	assert.True(t, IsResourceGIF(res))
}
