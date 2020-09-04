package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestColorChannel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		name  string
		value float64
	}{
		"foobar_0.0": {
			name:  "foobar",
			value: 0.0,
		},
		"foobar_0.5": {
			name:  "foobar",
			value: 0.5,
		},
		"foobar_1.0": {
			name:  "foobar",
			value: 1.0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := widget.NewColorChannel(tt.name, tt.value, nil)

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
			window.Resize(color.MinSize().Max(fyne.NewSize(100, 100)))

			test.AssertImageMatches(t, "color/channel_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
