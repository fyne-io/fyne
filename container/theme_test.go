package container

import (
	"image"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestThemeOverride_Icons(t *testing.T) {
	b := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {})
	o := NewThemeOverride(b, test.Theme())
	w := test.NewWindow(o)
	plain := w.Canvas().Capture().(*image.NRGBA)

	o.Theme = test.NewTheme()
	o.Refresh()
	changed := w.Canvas().Capture().(*image.NRGBA)

	assert.NotEqual(t, plain.Pix, changed.Pix)
}

func TestThemeOverride_Refresh(t *testing.T) {
	b := widget.NewButton("Test", func() {})
	o := NewThemeOverride(b, test.Theme())
	w := test.NewWindow(o)
	plain := w.Canvas().Capture().(*image.NRGBA)

	o.Theme = test.NewTheme()
	o.Refresh()
	changed := w.Canvas().Capture().(*image.NRGBA)

	assert.NotEqual(t, plain.Pix, changed.Pix)
}
