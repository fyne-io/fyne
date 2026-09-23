//go:build !no_glfw && !mobile && !windows && !darwin

package glfw

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fynecolor "fyne.io/fyne/v2/internal/color"
)

func TestToOSIconEncodesJPGAsPNG(t *testing.T) {
	img := testTrayIconImage()

	jpg := bytes.NewBuffer(nil)
	require.NoError(t, jpeg.Encode(jpg, img, &jpeg.Options{Quality: 95}))

	converted, err := toOSIcon(jpg.Bytes())
	require.NoError(t, err)

	_, format, err := image.DecodeConfig(bytes.NewReader(converted))
	require.NoError(t, err)
	assert.Equal(t, "png", format)

	_, err = png.Decode(bytes.NewReader(converted))
	require.NoError(t, err)
}

func TestToOSIconJPGPixelsMatchSystrayPixelConversion(t *testing.T) {
	source := color.NRGBA{R: 230, G: 43, B: 30, A: 255}
	img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetNRGBA(x, y, source)
		}
	}

	jpg := bytes.NewBuffer(nil)
	require.NoError(t, jpeg.Encode(jpg, img, &jpeg.Options{Quality: 95}))

	converted, err := toOSIcon(jpg.Bytes())
	require.NoError(t, err)

	decoded, _, err := image.Decode(bytes.NewReader(converted))
	require.NoError(t, err)

	r, g, b, a := fynecolor.ToNRGBA(decoded.At(0, 0))
	assert.InDelta(t, source.A, a, 0)
	assert.InDelta(t, source.R, r, 3)
	assert.InDelta(t, source.G, g, 3)
	assert.InDelta(t, source.B, b, 3)
}

func TestToOSIconLeavesPNGUnchanged(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 230, G: 43, B: 30, A: 255})

	source := bytes.NewBuffer(nil)
	require.NoError(t, png.Encode(source, img))

	converted, err := toOSIcon(source.Bytes())
	require.NoError(t, err)
	assert.Equal(t, source.Bytes(), converted)
}

func TestToOSIconDarwinLeavesOriginalBytesUnchanged(t *testing.T) {
	source := []byte("darwin accepts icon bytes unchanged")

	converted, err := toOSIconForRuntime(source, "darwin")
	require.NoError(t, err)
	assert.Equal(t, source, converted)
}

func TestToOSIconWindowsEncodesICO(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 230, G: 43, B: 30, A: 255})

	source := bytes.NewBuffer(nil)
	require.NoError(t, png.Encode(source, img))

	converted, err := toOSIconForRuntime(source.Bytes(), "windows")
	require.NoError(t, err)
	assert.Equal(t, []byte{0, 0, 1, 0}, converted[:4])
}

func TestToOSIconRejectsInvalidImage(t *testing.T) {
	_, err := toOSIconForRuntime([]byte("not an image"), "linux")
	require.Error(t, err)
}

func TestUsesUnixSystrayIcon(t *testing.T) {
	assert.True(t, usesUnixSystrayIcon("linux"))
	assert.True(t, usesUnixSystrayIcon("freebsd"))
	assert.True(t, usesUnixSystrayIcon("openbsd"))
	assert.True(t, usesUnixSystrayIcon("netbsd"))
	assert.False(t, usesUnixSystrayIcon("darwin"))
	assert.False(t, usesUnixSystrayIcon("windows"))
}

func testTrayIconImage() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.SetNRGBA(0, 0, color.NRGBA{R: 230, G: 43, B: 30, A: 255})
	img.SetNRGBA(1, 0, color.NRGBA{R: 33, G: 160, B: 79, A: 255})
	img.SetNRGBA(0, 1, color.NRGBA{R: 33, G: 89, B: 210, A: 255})
	img.SetNRGBA(1, 1, color.NRGBA{R: 246, G: 190, B: 34, A: 255})

	return img
}
