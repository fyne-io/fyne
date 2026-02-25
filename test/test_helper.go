//go:build !tamago && !noos

package test

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/internal/test"
)

// AssertCanvasTappableAt asserts that the canvas is tappable at the given position.
func AssertCanvasTappableAt(t *testing.T, c fyne.Canvas, pos fyne.Position) bool {
	if o, _ := findTappable(c, pos); o == nil {
		t.Errorf("No tappable found at %#v", pos)
		return false
	}
	return true
}

// AssertObjectRendersToImage asserts that the given `CanvasObject` renders the same image as the one stored in the master file.
// The theme used is the standard test theme which may look different to how it shows on your device.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Since 2.3
func AssertObjectRendersToImage(t *testing.T, masterFilename string, o fyne.CanvasObject, msgAndArgs ...any) bool {
	c := NewCanvasWithPainter(software.NewPainter())
	c.SetPadded(false)
	size := o.MinSize().Max(o.Size())
	c.SetContent(o)
	c.Resize(size) // ensure we are large enough for current size

	return AssertRendersToImage(t, masterFilename, c, msgAndArgs...)
}

// AssertObjectRendersToMarkup asserts that the given `CanvasObject` renders the same markup as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the rendered markup is not equal to the loaded master markup.
// In this case the rendered markup is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Be aware, that the indentation has to use tab characters ('\t') instead of spaces.
// Every element starts on a new line indented one more than its parent.
// Closing elements stand on their own line, too, using the same indentation as the opening element.
// The only exception to this are text elements which do not contain line breaks unless the text includes them.
//
// Since 2.3
func AssertObjectRendersToMarkup(t *testing.T, masterFilename string, o fyne.CanvasObject, msgAndArgs ...any) bool {
	c := NewCanvas()
	c.SetPadded(false)
	size := o.MinSize().Max(o.Size())
	c.SetContent(o)
	c.Resize(size) // ensure we are large enough for current size

	return AssertRendersToMarkup(t, masterFilename, c, msgAndArgs...)
}

// AssertImageMatches asserts that the given image is the same as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
func AssertImageMatches(t *testing.T, masterFilename string, img image.Image, msgAndArgs ...any) bool {
	return test.AssertImageMatches(t, masterFilename, img, msgAndArgs...)
}

// AssertRendersToImage asserts that the given canvas renders the same image as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Since 2.3
func AssertRendersToImage(t *testing.T, masterFilename string, c fyne.Canvas, msgAndArgs ...any) bool {
	return AssertImageMatches(t, masterFilename, c.Capture(), msgAndArgs...)
}

// AssertRendersToMarkup asserts that the given canvas renders the same markup as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the rendered markup is not equal to the loaded master markup.
// In this case the rendered markup is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Be aware, that the indentation has to use tab characters ('\t') instead of spaces.
// Every element starts on a new line indented one more than its parent.
// Closing elements stand on their own line, too, using the same indentation as the opening element.
// The only exception to this are text elements which do not contain line breaks unless the text includes them.
//
// Since: 2.0
func AssertRendersToMarkup(t *testing.T, masterFilename string, c fyne.Canvas, msgAndArgs ...any) bool {
	wd, err := os.Getwd()
	require.NoError(t, err)

	got := snapshot(c)
	masterPath := filepath.Join(wd, "testdata", masterFilename)
	failedPath := filepath.Join(wd, "testdata/failed", masterFilename)
	_, err = os.Stat(masterPath)
	if os.IsNotExist(err) {
		require.NoError(t, writeMarkup(failedPath, got))
		t.Errorf("Master not found at %s. Markup written to %s might be used as master.", masterPath, failedPath)
		return false
	}

	raw, err := os.ReadFile(masterPath)
	require.NoError(t, err)
	master := strings.ReplaceAll(string(raw), "\r", "")

	var msg string
	if len(msgAndArgs) > 0 {
		msg = fmt.Sprintf(msgAndArgs[0].(string)+"\n", msgAndArgs[1:]...)
	}
	if !assert.Equal(t, master, got, "%sMarkup did not match master. Actual markup written to file://%s.", msg, failedPath) {
		require.NoError(t, writeMarkup(failedPath, got))
		return false
	}
	return true
}

// ApplyTheme sets the given theme and waits for it to be applied to the current app.
func ApplyTheme(t *testing.T, theme fyne.Theme) {
	require.IsType(t, &app{}, fyne.CurrentApp())
	a := fyne.CurrentApp().(*app)
	a.Settings().SetTheme(theme)
	for a.lastAppliedTheme() != theme {
		time.Sleep(5 * time.Millisecond)
	}
}

// TempWidgetRenderer allows test scripts to gain access to the current renderer for a widget.
// This can be used for verifying correctness of rendered components for a widget in unit tests.
// The widget renderer is automatically destroyed when the test ends.
//
// Since: 2.5
func TempWidgetRenderer(t *testing.T, wid fyne.Widget) fyne.WidgetRenderer {
	t.Cleanup(func() { cache.DestroyRenderer(wid) })
	return cache.Renderer(wid)
}

// WithTestTheme runs a function with the testTheme temporarily set.
func WithTestTheme(t *testing.T, f func()) {
	settings := fyne.CurrentApp().Settings()
	current := settings.Theme()
	ApplyTheme(t, NewTheme())
	defer ApplyTheme(t, current)
	f()
}
