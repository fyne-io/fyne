//go:build !wasm && !test_web_driver && !mobile

package glfw

import (
	"bytes"
	"log"
	"testing"

	"fyne.io/fyne/v2"

	"github.com/stretchr/testify/assert"
)

func TestWindow_SizeLimitWarning(t *testing.T) {
	var output bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(previous) })

	for _, size := range []fyne.Size{fyne.NewSize(100_000, 100), fyne.NewSize(100, 100_000)} {
		w := &window{canvas: &glCanvas{scale: 1, texScale: 1}}
		output.Reset()
		w.screenSize(fyne.NewSize(16384, 16384))
		assert.Empty(t, output.String())
		w.screenSize(size)
		assert.Contains(t, output.String(), "wrapping, truncation, or scrolling")
		output.Reset()
		w.screenSize(size)
		assert.Empty(t, output.String())
	}
}
