package dialog

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestcolorGreyscalePicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	color := newColorGreyscalePicker(nil)

	window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
	window.Resize(color.MinSize().Max(fyne.NewSize(100, 50)))

	test.AssertImageMatches(t, "color/greyscale_picker_layout.png", window.Canvas().Capture())

	window.Close()
}

func TestcolorBasicPicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	color := newColorBasicPicker(nil)

	window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
	window.Resize(color.MinSize().Max(fyne.NewSize(100, 50)))

	test.AssertImageMatches(t, "color/basic_picker_layout.png", window.Canvas().Capture())

	window.Close()
}

func TestcolorRecentPicker_Layout(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	// Inject recent preferences
	a.Preferences().SetString("color_recents", "#0000FF,#008000,#FF0000")

	color := newColorRecentPicker(nil)

	window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
	window.Resize(color.MinSize().Max(fyne.NewSize(100, 50)))

	test.AssertImageMatches(t, "color/recent_picker_layout.png", window.Canvas().Capture())

	window.Close()
}

func TestcolorAdvancedPicker_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		color color.Color
		model string
	}{
		"primary_rgb": {
			color: theme.PrimaryColor(),
			model: "RGB",
		},
		"primary_hsl": {
			color: theme.PrimaryColor(),
			model: "HSL",
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := newColorAdvancedPicker(tt.color, nil)

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
