package canvas_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestBlur(t *testing.T) {
	test.NewTempApp(t)
	bg := canvas.NewImageFromFile("testdata/Utah_teapot.png")
	b1 := canvas.NewBlur(35)
	w := test.NewTempWindow(t, container.NewWithoutLayout(bg, b1))
	w.SetPadded(false)
	size := fyne.NewSize(300, 243)
	bg.Resize(size)
	w.Resize(size)

	b1.Move(fyne.NewPos(120, 110))
	b1.Resize(fyne.NewSize(140, 100))

	test.AssertRendersToImage(t, "blur.png", w.Canvas())
}
