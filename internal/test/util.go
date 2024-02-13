package test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertImageMatches asserts that the given image is the same as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
func AssertImageMatches(t *testing.T, masterFilename string, img image.Image, msgAndArgs ...interface{}) bool {
	wd, err := os.Getwd()
	require.NoError(t, err)
	masterPath := filepath.Join(wd, "testdata", masterFilename)
	failedPath := filepath.Join(wd, "testdata/failed", masterFilename)
	_, err = os.Stat(masterPath)
	if os.IsNotExist(err) {
		require.NoError(t, writeImage(failedPath, img))
		t.Errorf("Master not found at %s. Image written to %s might be used as master.", masterPath, failedPath)
		return false
	}

	file, err := os.Open(masterPath)
	require.NoError(t, err)
	defer file.Close()
	raw, _, err := image.Decode(file)
	require.NoError(t, err)

	masterPix := pixelsForImage(t, raw) // let's just compare the pixels directly
	capturePix := pixelsForImage(t, img)

	// On darwin/arm64, there are slight differences in the rendering.
	// Use a slower, more lenient comparison. If that fails,
	// fall back to the strict comparison for a more detailed error message.
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" && pixCloseEnough(masterPix, capturePix) {
		return true
	}

	var msg string
	if len(msgAndArgs) > 0 {
		msg = fmt.Sprintf(msgAndArgs[0].(string)+"\n", msgAndArgs[1:]...)
	}
	if !assert.Equal(t, masterPix, capturePix, "%sImage did not match master. Actual image written to file://%s.", msg, failedPath) {
		require.NoError(t, writeImage(failedPath, img))
		return false
	}
	return true
}

// pixCloseEnough reports whether a and b are mostly the same.
func pixCloseEnough(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	mismatches := 0

	for i, v := range a {
		w := b[i]
		if v == w {
			continue
		}
		// Allow a small delta for rendering variation.
		delta := int(v) - int(w) // use int to avoid overflow
		if delta < 0 {
			delta *= -1
		}
		if delta > 4 {
			return false
		}
		mismatches++
	}

	// Allow up to 1% of pixels to mismatch.
	return mismatches < len(a)/100
}

// NewCheckedImage returns a new black/white checked image with the specified size
// and the specified amount of horizontal and vertical tiles.
func NewCheckedImage(w, h, hTiles, vTiles int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	colors := []color.Color{color.White, color.Black}
	tileWidth := float64(w) / float64(hTiles)
	tileHeight := float64(h) / float64(vTiles)
	for y := 0; y < h; y++ {
		yTile := int(math.Floor(float64(y) / tileHeight))
		for x := 0; x < w; x++ {
			xTile := int(math.Floor(float64(x) / tileWidth))
			img.Set(x, y, colors[(xTile+yTile)%2])
		}
	}
	return img
}

func pixelsForImage(t *testing.T, img image.Image) []uint8 {
	var pix []uint8
	if data, ok := img.(*image.RGBA); ok {
		pix = data.Pix
	} else if data, ok := img.(*image.NRGBA); ok {
		pix = data.Pix
	}
	if pix == nil {
		t.Error("Master image is unsupported type")
	}

	return pix
}

func writeImage(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err = png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
