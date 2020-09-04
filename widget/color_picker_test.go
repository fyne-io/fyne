package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestColorGreyscalePicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	color := widget.NewColorGreyscalePicker(nil)

	window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
	window.Resize(color.MinSize().Max(fyne.NewSize(100, 50)))

	test.AssertImageMatches(t, "color/greyscale_picker_layout.png", window.Canvas().Capture())

	window.Close()
}

func TestColorBasicPicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	color := widget.NewColorBasicPicker(nil)

	window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
	window.Resize(color.MinSize().Max(fyne.NewSize(100, 50)))

	test.AssertImageMatches(t, "color/basic_picker_layout.png", window.Canvas().Capture())

	window.Close()
}

func TestColorRecentPicker_Layout(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#0000FF,#008000,#FF0000")

	color := widget.NewColorRecentPicker(nil)

	window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
	window.Resize(color.MinSize().Max(fyne.NewSize(100, 50)))

	test.AssertImageMatches(t, "color/recent_picker_layout.png", window.Canvas().Capture())

	window.Close()
}

func TestColorAdvancedPicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		color color.Color
		model string
	}{
		"primary_rgb": {
			color: theme.PrimaryColor(),
			model: "rgb",
		},
		"primary_hsl": {
			color: theme.PrimaryColor(),
			model: "hsl",
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := widget.NewColorAdvancedPicker(tt.color, nil)

			if tt.model != "" {
				color.ColorModel = tt.model
			}

			color.Refresh()

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
			window.Resize(color.MinSize().Max(fyne.NewSize(200, 200)))

			test.AssertImageMatches(t, "color/advanced_picker_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
