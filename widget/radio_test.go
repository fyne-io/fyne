package widget

import (
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestRadioSize(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	min := radio.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	radio2 := NewRadio([]string{"Hi", "Hi"}, nil)
	min2 := radio2.MinSize()

	assert.Equal(t, min.Width, min2.Width)
	assert.True(t, min2.Height > min.Height)
}

func TestRadioSelected(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.OnMouseDown(&fyne.PointerEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.True(t, selected == "Hi")
}

func TestRadioUnselected(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radio.OnMouseDown(&fyne.PointerEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.True(t, selected == "")
}

func TestRadioSelectedOther(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi", "Hi2"}, func(sel string) {
		selected = sel
	})
	radio.OnMouseDown(&fyne.PointerEvent{Position: fyne.NewPos(theme.Padding(), radio.MinSize().Height-theme.Padding())})

	assert.True(t, selected == "Hi2")
}
