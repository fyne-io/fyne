package playground

import (
	"encoding/base64"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/internal/test"
	_ "fyne.io/fyne/v2/test" // for test app initialization
	"fyne.io/fyne/v2/theme"
)

func TestRender(t *testing.T) {
	obj := canvas.NewRectangle(color.Black)
	obj.SetMinSize(fyne.NewSquareSize(10))

	img := software.Render(obj, test.DarkTheme(theme.DefaultTheme()))
	assert.NotNil(t, img)

	enc, err := encodeImage(img)
	assert.NoError(t, err)
	assert.Equal(t, "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAIAAAACUFjqAAAAFElEQVR4nGJiwAtGpbECQAAAAP//DogAFaNSFa8AAAAASUVORK5CYII=", enc)

	bytes, err := base64.StdEncoding.DecodeString(enc)
	assert.NoError(t, err)
	assert.Equal(t, "PNG", string(bytes)[1:4])
}
