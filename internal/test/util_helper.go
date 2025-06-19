//go:build !tamago && !noos

package test

import (
	"fmt"
	"image"
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
func AssertImageMatches(t *testing.T, masterFilename string, img image.Image, msgAndArgs ...any) bool {
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
