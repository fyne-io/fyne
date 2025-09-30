//go:build !tamago && !noos

package test

import (
	"testing"

	"fyne.io/fyne/v2"
)

// NewTempWindow creates and registers a new window for test purposes.
// This window will get removed automatically once the running test ends.
//
// Since: 2.5
func NewTempWindow(t testing.TB, content fyne.CanvasObject) fyne.Window {
	window := NewWindow(content)
	t.Cleanup(window.Close)
	return window
}
