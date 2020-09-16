package dialog

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func Test_colorChannel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	min := 0
	max := 100
	size := fyne.NewSize(250, 50)

	for name, tt := range map[string]struct {
		name  string
		value int
	}{
		"foobar_0": {
			name:  "foobar",
			value: 0,
		},
		"foobar_50": {
			name:  "foobar",
			value: 50,
		},
		"foobar_100": {
			name:  "foobar",
			value: 100,
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := newColorChannel(tt.name, min, max, tt.value, nil)
			color.Resize(size)

			window := test.NewWindow(color)

			test.AssertImageMatches(t, "color/channel_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
