package dialog

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func Test_colorGreyscalePicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	color := newColorGreyscalePicker(nil)

	window := test.NewWindow(container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(360, 60)))

	test.AssertRendersToImage(t, "color/picker_layout_greyscale.png", window.Canvas())

	window.Close()
}

func Test_colorBasicPicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	color := newColorBasicPicker(nil)

	window := test.NewWindow(container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(360, 60)))

	test.AssertRendersToImage(t, "color/picker_layout_basic.png", window.Canvas())

	window.Close()
}

func Test_colorRecentPicker_Layout(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#0000FF,#008000,#FF0000")

	color := newColorRecentPicker(nil)

	window := test.NewWindow(container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(360, 60)))

	test.AssertRendersToImage(t, "color/picker_layout_recent.png", window.Canvas())

	window.Close()
}

func Test_colorAdvancedPicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	color := newColorAdvancedPicker(theme.PrimaryColor(), nil)

	color.Refresh()

	window := test.NewWindow(container.NewCenter(color))
	window.Resize(color.MinSize().Max(fyne.NewSize(200, 200)))

	test.AssertRendersToImage(t, "color/picker_layout_advanced.png", window.Canvas())

	window.Close()
}
