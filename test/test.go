package test

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/internal/test"

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

// AssertObjectRendersToImage asserts that the given `CanvasObject` renders the same image as the one stored in the master file.
// The theme used is the standard test theme which may look different to how it shows on your device.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Since 2.3
func AssertObjectRendersToImage(t *testing.T, masterFilename string, o fyne.CanvasObject, msgAndArgs ...interface{}) bool {
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
func AssertObjectRendersToMarkup(t *testing.T, masterFilename string, o fyne.CanvasObject, msgAndArgs ...interface{}) bool {
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
func AssertImageMatches(t *testing.T, masterFilename string, img image.Image, msgAndArgs ...interface{}) bool {
	return test.AssertImageMatches(t, masterFilename, img, msgAndArgs...)
}

// AssertRendersToImage asserts that the given canvas renders the same image as the one stored in the master file.
// The master filename is relative to the `testdata` directory which is relative to the test.
// The test `t` fails if the given image is not equal to the loaded master image.
// In this case the given image is written into a file in `testdata/failed/<masterFilename>` (relative to the test).
// This path is also reported, thus the file can be used as new master.
//
// Since 2.3
func AssertRendersToImage(t *testing.T, masterFilename string, c fyne.Canvas, msgAndArgs ...interface{}) bool {
	return test.AssertImageMatches(t, masterFilename, c.Capture(), msgAndArgs...)
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
func AssertRendersToMarkup(t *testing.T, masterFilename string, c fyne.Canvas, msgAndArgs ...interface{}) bool {
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

// Drag drags at an absolute position on the canvas.
// deltaX/Y is the dragging distance: <0 for dragging up/left, >0 for dragging down/right.
func Drag(c fyne.Canvas, pos fyne.Position, deltaX, deltaY float32) {
	matches := func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Draggable); ok {
			return true
		}
		return false
	}
	o, p, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o == nil {
		return
	}
	e := &fyne.DragEvent{
		PointEvent: fyne.PointEvent{Position: p},
		Dragged:    fyne.Delta{DX: deltaX, DY: deltaY},
	}
	o.(fyne.Draggable).Dragged(e)
	o.(fyne.Draggable).DragEnd()
}

// FocusNext focuses the next focusable on the canvas.
func FocusNext(c fyne.Canvas) {
	if tc, ok := c.(*testCanvas); ok {
		tc.focusManager().FocusNext()
	} else {
		fyne.LogError("FocusNext can only be called with a test canvas", nil)
	}
}

// FocusPrevious focuses the previous focusable on the canvas.
func FocusPrevious(c fyne.Canvas) {
	if tc, ok := c.(*testCanvas); ok {
		tc.focusManager().FocusPrevious()
	} else {
		fyne.LogError("FocusPrevious can only be called with a test canvas", nil)
	}
}

// LaidOutObjects returns all fyne.CanvasObject starting at the given fyne.CanvasObject which is laid out previously.
func LaidOutObjects(o fyne.CanvasObject) (objects []fyne.CanvasObject) {
	if o != nil {
		objects = layoutAndCollect(objects, o, o.MinSize().Max(o.Size()))
	}
	return objects
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
	o, p, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o != nil {
		hovered = o.(desktop.Hoverable)
		me := &desktop.MouseEvent{
			PointEvent: fyne.PointEvent{
				AbsolutePosition: pos,
				Position:         p,
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

// Scroll scrolls at an absolute position on the canvas.
// deltaX/Y is the scrolling distance: <0 for scrolling up/left, >0 for scrolling down/right.
func Scroll(c fyne.Canvas, pos fyne.Position, deltaX, deltaY float32) {
	matches := func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Scrollable); ok {
			return true
		}
		return false
	}
	o, _, _ := driver.FindObjectAtPositionMatching(pos, matches, c.Overlays().Top(), c.Content())
	if o == nil {
		return
	}

	e := &fyne.ScrollEvent{Scrolled: fyne.Delta{DX: deltaX, DY: deltaY}}
	o.(fyne.Scrollable).Scrolled(e)
}

// DoubleTap simulates a double left mouse click on the specified object.
func DoubleTap(obj fyne.DoubleTappable) {
	ev, c := prepareTap(obj, fyne.NewPos(1, 1))
	handleFocusOnTap(c, obj)
	obj.DoubleTapped(ev)
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
	if o, p := findTappable(c, pos); o != nil {
		tap(c, o.(fyne.Tappable), &fyne.PointEvent{AbsolutePosition: pos, Position: p})
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
		_, ok := object.(fyne.Tappable)
		return ok
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
				unfocus = true
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

func writeMarkup(path string, markup string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(markup), 0644)
}
