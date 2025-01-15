package test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
)

func Test_driver_AbsolutePositionForObject(t *testing.T) {
	d := &driver{}
	w := d.CreateWindow("Test Window")
	o := fynecanvas.NewRectangle(color.Black)
	w.SetContent(o)
	w.Resize(fyne.NewSize(320, 200))

	t.Run("for padded window", func(t *testing.T) {
		w.SetPadded(true)
		assert.Equal(t, fyne.NewPos(2, 1), d.AbsolutePositionForObject(o), "safe area offset (2,3) is subtracted")
	})

	t.Run("for non-padded window", func(t *testing.T) {
		w.SetPadded(false)
		assert.Equal(t, fyne.NewPos(-2, -3), d.AbsolutePositionForObject(o), "safe area offset (2,3) is subtracted")
	})
}

func TestDriver_CreateWindow(t *testing.T) {
	d := &driver{}
	w := d.CreateWindow("Test Window")

	assert.Equal(t, "Test Window", w.Title())
}
