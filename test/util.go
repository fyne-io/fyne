package test

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertImagesEqual asserts that the master image at given filename (below testdata) and the given image are equal.
// It fails the test if not. In this case the actual image is written to disk and reported.
func AssertImagesEqual(t *testing.T, masterFilename string, got image.Image) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	masterPath := filepath.Join("testdata", masterFilename)
	var expected image.Image
	_, err = os.Stat(masterPath)
	if os.IsNotExist(err) {
		fmt.Printf("Master image does not exist at %s.\nAssume initial run and use empty image for comparison.\n", filepath.Join(wd, masterPath))
		expected = image.NewRGBA(image.Rectangle{})
	} else {
		file, err := os.Open(masterPath)
		require.NoError(t, err)
		defer file.Close()
		raw, _, err := image.Decode(file)
		require.NoError(t, err)
		normalized := image.NewRGBA(raw.Bounds())
		draw.Draw(normalized, normalized.Bounds(), raw, image.Pt(0, 0), draw.Src)
		expected = normalized
	}
	if !assert.Equal(t, expected, got) {
		path := filepath.Join("testdata/failed", masterFilename)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
		file, err := os.Create(path)
		require.NoError(t, err)
		defer file.Close()
		require.NoError(t, png.Encode(file, got))
		fmt.Println("Images were not equal. Actual image written to:", filepath.Join(wd, path))
	}
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

// WidgetRenderer allows test scripts to gain access to the current renderer for a widget.
// This can be used for verifying correctness of rendered components for a widget in unit tests.
func WidgetRenderer(wid fyne.Widget) fyne.WidgetRenderer {
	return cache.Renderer(wid)
}

func typeChars(chars []rune, keyDown func(rune)) {
	for _, char := range chars {
		keyDown(char)
	}
}
