package test

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/driver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertCanvasTappableAt asserts that the canvas is tappable at the given position.
func AssertCanvasTappableAt(t *testing.T, c fyne.Canvas, pos fyne.Position) bool {
	if o, _ := findTappable(c, pos); o == nil {
		t.Errorf("No tappable found at %#v", pos)
		return false
	}
	return true
}

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

// MoveMouse simulates a mouse movement to the given position.
func MoveMouse(c fyne.Canvas, pos fyne.Position) {
	if fyne.CurrentDevice().IsMobile() {
		return
	}

	tc, _ := c.(*testCanvas)
	var oldHovered, hovered desktop.Hoverable
	if tc != nil {
		oldHovered = tc.hovered
	}
	matches := func(object fyne.CanvasObject) bool {
		if _, ok := object.(desktop.Hoverable); ok {
			return true
		}
		return false
	}
	o, absPos, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o != nil {
		hovered = o.(desktop.Hoverable)
		me := &desktop.MouseEvent{
			PointEvent: fyne.PointEvent{
				AbsolutePosition: pos,
				Position:         pos.Subtract(absPos),
			},
		}
		if hovered == oldHovered {
			hovered.MouseMoved(me)
		} else {
			if oldHovered != nil {
				oldHovered.MouseOut()
			}
			hovered.MouseIn(me)
		}
	} else if oldHovered != nil {
		oldHovered.MouseOut()
	}
	if tc != nil {
		tc.hovered = hovered
	}
}

// Tap simulates a left mouse click on the specified object.
func Tap(obj fyne.Tappable) {
	TapAt(obj, fyne.NewPos(1, 1))
}

// TapAt simulates a left mouse click on the passed object at a specified place within it.
func TapAt(obj fyne.Tappable, pos fyne.Position) {
	ev, c := prepareTap(obj, pos)
	tap(c, obj, ev)
}

// TapCanvas taps at an absolute position on the canvas.
func TapCanvas(c fyne.Canvas, pos fyne.Position) {
	if o, absPos := findTappable(c, pos); o != nil {
		tap(c, o.(fyne.Tappable), &fyne.PointEvent{AbsolutePosition: pos, Position: pos.Subtract(absPos)})
	}
}

// TapSecondary simulates a right mouse click on the specified object.
func TapSecondary(obj fyne.SecondaryTappable) {
	TapSecondaryAt(obj, fyne.NewPos(1, 1))
}

// TapSecondaryAt simulates a right mouse click on the passed object at a specified place within it.
func TapSecondaryAt(obj fyne.SecondaryTappable, pos fyne.Position) {
	ev, c := prepareTap(obj, pos)
	handleFocusOnTap(c, obj)
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
	for a.lastAppliedTheme() != theme {
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
	ApplyTheme(t, NewTheme())
	defer ApplyTheme(t, current)
	f()
}

func findTappable(c fyne.Canvas, pos fyne.Position) (o fyne.CanvasObject, p fyne.Position) {
	matches := func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		}
		return false
	}
	o, p, _ = driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	return
}

func prepareTap(obj interface{}, pos fyne.Position) (*fyne.PointEvent, fyne.Canvas) {
	d := fyne.CurrentApp().Driver()
	ev := &fyne.PointEvent{Position: pos}
	var c fyne.Canvas
	if co, ok := obj.(fyne.CanvasObject); ok {
		c = d.CanvasForObject(co)
		ev.AbsolutePosition = d.AbsolutePositionForObject(co).Add(pos)
	}
	return ev, c
}

func tap(c fyne.Canvas, obj fyne.Tappable, ev *fyne.PointEvent) {
	handleFocusOnTap(c, obj)
	obj.Tapped(ev)
}

func handleFocusOnTap(c fyne.Canvas, obj interface{}) {
	if c == nil {
		return
	}
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
