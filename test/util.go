package test

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertImageMatches asserts that the given image is the same as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
func AssertImageMatches(t *testing.T, masterFilename string, img image.Image) bool {
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
	expected := image.NewRGBA(raw.Bounds())
	draw.Draw(expected, expected.Bounds(), raw, image.Pt(0, 0), draw.Src)

	if !assert.Equal(t, expected, img, "Image did not match master. Actual image written to %s.", failedPath) {
		require.NoError(t, writeImage(failedPath, img))
		return false
	}
	return true
}

// Tap simulates a left mouse click on the specified object.
func Tap(obj fyne.Tappable) {
	TapAt(obj, fyne.NewPos(1, 1))
}

// TapAt simulates a left mouse click on the passed object at a specified place within it.
func TapAt(obj fyne.Tappable, pos fyne.Position) {
	if focus, ok := obj.(fyne.Focusable); ok {
		if focus != Canvas().Focused() {
			Canvas().Focus(focus)
		}
	}

	ev := &fyne.PointEvent{Position: pos}
	obj.Tapped(ev)
}

// TapSecondary simulates a right mouse click on the specified object.
func TapSecondary(obj fyne.SecondaryTappable) {
	TapSecondaryAt(obj, fyne.NewPos(1, 1))
}

// TapSecondaryAt simulates a right mouse click on the passed object at a specified place within it.
func TapSecondaryAt(obj fyne.SecondaryTappable, pos fyne.Position) {
	if focus, ok := obj.(fyne.Focusable); ok {
		if focus != Canvas().Focused() {
			Canvas().Focus(focus)
		}
	}

	ev := &fyne.PointEvent{Position: pos}
	obj.TappedSecondary(ev)
}

// Type performs a series of key events to simulate typing of a value into the specified object.
// The focusable object will be focused before typing begins.
// The chars parameter will be input one rune at a time to the focused object.
func Type(obj fyne.Focusable, chars string) {
	obj.FocusGained()

	typeChars([]rune(chars), obj.TypedRune)
}

// TypeOnCanvas is like the Type function but it passes the key events to the canvas object
// rather than a focusable widget.
func TypeOnCanvas(c fyne.Canvas, chars string) {
	typeChars([]rune(chars), c.OnTypedRune())
}

// WaitForThemeToBeApplied waits for the current theme to be applied to the current app.
func WaitForThemeToBeApplied(t *testing.T) {
	require.IsType(t, &testApp{}, fyne.CurrentApp())
	a := fyne.CurrentApp().(*testApp)
	for a.appliedTheme != a.Settings().Theme() {
		time.Sleep(1 * time.Millisecond)
	}
}

// WidgetRenderer allows test scripts to gain access to the current renderer for a widget.
// This can be used for verifying correctness of rendered components for a widget in unit tests.
func WidgetRenderer(wid fyne.Widget) fyne.WidgetRenderer {
	return cache.Renderer(wid)
}

// WithTestTheme runs a function with the testTheme temporarily set.
func WithTestTheme(t *testing.T, f func()) {
	settings := fyne.CurrentApp().Settings()
	current := settings.Theme()
	settings.SetTheme(&testTheme{})
	WaitForThemeToBeApplied(t)
	defer func() {
		settings.SetTheme(current)
		WaitForThemeToBeApplied(t)
	}()
	f()
}

func typeChars(chars []rune, keyDown func(rune)) {
	for _, char := range chars {
		keyDown(char)
	}
}

func writeImage(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
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
