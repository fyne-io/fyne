package dialog

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func Test_colorGreyscalePicker_Layout(t *testing.T) {
	test.NewTempApp(t)

	color := newColorGreyscalePicker(nil)

	window := test.NewTempWindow(t, container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(360, 60)))

	test.AssertRendersToImage(t, "color/picker_layout_greyscale.png", window.Canvas())
}

func Test_colorBasicPicker_Layout(t *testing.T) {
	test.NewTempApp(t)

	color := newColorBasicPicker(nil)

	window := test.NewTempWindow(t, container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(360, 60)))

	test.AssertRendersToImage(t, "color/picker_layout_basic.png", window.Canvas())
}

func Test_colorRecentPicker_Layout(t *testing.T) {
	a := test.NewTempApp(t)

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#0000FF,#008000,#FF0000")

	color := newColorRecentPicker(nil)

	window := test.NewTempWindow(t, container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(360, 60)))

	test.AssertRendersToImage(t, "color/picker_layout_recent.png", window.Canvas())
}

func Test_colorAdvancedPicker_Layout(t *testing.T) {
	test.NewTempApp(t)

	color := newColorAdvancedPicker(theme.Color(theme.ColorNamePrimary), nil)

	color.Refresh()

	window := test.NewTempWindow(t, container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(200, 200)))

	test.AssertRendersToImage(t, "color/picker_layout_advanced.png", window.Canvas())
}
