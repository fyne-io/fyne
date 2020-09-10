package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type dummyHolder struct {
	r, g, b, a, h, s, l float64
}

func (d *dummyHolder) RGBA() (float64, float64, float64, float64) {
	return d.r, d.g, d.b, d.a
}

func (d *dummyHolder) SetRGBA(float64, float64, float64, float64) {
	// Do nothing
}

func (d *dummyHolder) HSLA() (float64, float64, float64, float64) {
	return d.h, d.s, d.l, d.a
}

func (d *dummyHolder) SetHSLA(float64, float64, float64, float64) {
	// Do nothing
}

func TestColorArea_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	holder := &dummyHolder{
		r: 0.5,
		g: 0.5,
		b: 0.5,
		a: 0.5,
		h: 0.5,
		s: 0.5,
		l: 0.5,
	}

	for name, tt := range map[string]struct {
		area *widget.ColorArea
	}{
		"rgb": {
			area: widget.NewRGBAColorArea(holder),
		},
		"hsl": {
			area: widget.NewHSLAColorArea(holder),
		},
	} {
		t.Run(name, func(t *testing.T) {
			window := test.NewWindow(tt.area)
			window.Resize(tt.area.MinSize().Max(fyne.NewSize(100, 100)))

			test.AssertImageMatches(t, "color/area_layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
