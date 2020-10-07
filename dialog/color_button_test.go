package dialog

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func Test_colorButton_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		color   color.Color
		hovered bool
	}{
		"primary": {
			color: theme.PrimaryColor(),
		},
		"primary_hovered": {
			color:   theme.PrimaryColor(),
			hovered: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			color := newColorButton(tt.color, nil)

			if tt.hovered {
				color.MouseIn(nil)
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), color))
			window.Resize(color.MinSize().Max(fyne.NewSize(50, 50)))

			test.AssertImageMatches(t, "color/button_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
