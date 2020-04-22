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
	"fyne.io/fyne/internal/driver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SecondaryTappableCanvasObject is an interface used by the secondary tap helper methods.
type SecondaryTappableCanvasObject interface {
	fyne.CanvasObject
	fyne.SecondaryTappable
}

// TappableCanvasObject is an interface used by the tap helper methods.
type TappableCanvasObject interface {
	fyne.CanvasObject
	fyne.Tappable
}

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
func Tap(obj TappableCanvasObject) {
	TapAt(obj, fyne.NewPos(1, 1))
}

// TapAt simulates a left mouse click on the passed object at a specified place within it.
func TapAt(obj TappableCanvasObject, pos fyne.Position) {
	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	absPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
	ev := &fyne.PointEvent{AbsolutePosition: absPos.Add(pos), Position: pos}
	tap(c, obj, ev)
}

// TapCanvas taps at an absolute position on the canvas.
// It fails the test if there is no fyne.Tappable reachable at the position.
func TapCanvas(t *testing.T, c fyne.Canvas, pos fyne.Position) {
	matches := func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		}
		return false
	}
	o, absPos := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	require.NotNil(t, o, "no tappable found at %#v", pos)
	tap(c, o.(TappableCanvasObject), &fyne.PointEvent{AbsolutePosition: pos, Position: pos.Subtract(absPos)})
}

// TapSecondary simulates a right mouse click on the specified object.
func TapSecondary(obj SecondaryTappableCanvasObject) {
	TapSecondaryAt(obj, fyne.NewPos(1, 1))
}

// TapSecondaryAt simulates a right mouse click on the passed object at a specified place within it.
func TapSecondaryAt(obj SecondaryTappableCanvasObject, pos fyne.Position) {
	c := fyne.CurrentApp().Driver().CanvasForObject(obj)
	handleFocusOnTap(c, obj)
	absPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
	ev := &fyne.PointEvent{AbsolutePosition: absPos.Add(pos), Position: pos}
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

// ApplyTheme sets the given theme and waits for it to be applied to the current app.
func ApplyTheme(t *testing.T, theme fyne.Theme) {
	require.IsType(t, &testApp{}, fyne.CurrentApp())
	a := fyne.CurrentApp().(*testApp)
	a.Settings().SetTheme(theme)
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
	ApplyTheme(t, &testTheme{})
	defer ApplyTheme(t, current)
	f()
}

func tap(c fyne.Canvas, obj TappableCanvasObject, ev *fyne.PointEvent) {
	handleFocusOnTap(c, obj)
	obj.Tapped(ev)
}

func handleFocusOnTap(c fyne.Canvas, obj fyne.CanvasObject) {
	unfocus := true
	if focus, ok := obj.(fyne.Focusable); ok {
		if dis, ok := obj.(fyne.Disableable); !ok || !dis.Disabled() {
			unfocus = false
			if focus != c.Focused() {
				c.Focus(focus)
			}
		}
	}
	if unfocus {
		c.Unfocus()
	}
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
