package commands

import "testing"

func TestNewGetter(t *testing.T) {
	g := NewGetter()
	g.SetAppID("io.fyne.text") // would crash if not set up internally correctly
}
