package metadata

import (
	"bytes"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
)

func TestScaleIcon(t *testing.T) {
	data, err := fyne.LoadResourceFromPath("./testdata/fyne.png")
	require.NoError(t, err)

	assert.Equal(t, "fyne.png", data.Name())
	img, err := png.Decode(bytes.NewReader(data.Content()))
	require.NoError(t, err)
	assert.Equal(t, 512, img.Bounds().Dx())

	smallData := ScaleIcon(data, 256)

	assert.Equal(t, "fyne-256.png", smallData.Name())
	img, err = png.Decode(bytes.NewReader(smallData.Content()))
	require.NoError(t, err)
	assert.Equal(t, 256, img.Bounds().Dx())
}
