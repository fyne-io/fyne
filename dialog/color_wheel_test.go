package dialog

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func Test_colorWheel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	wheel := newColorWheel(nil)
	wheel.SetHSLA(180, 100, 50, 255)
	window := test.NewWindow(wheel)
	window.Resize(wheel.MinSize().Max(fyne.NewSize(100, 100)))

	test.AssertRendersToImage(t, "color/wheel_layout.png", window.Canvas())
	test.AssertRendersToMarkup(t, "color/wheel_layout.xml", window.Canvas())

	window.Close()
}
