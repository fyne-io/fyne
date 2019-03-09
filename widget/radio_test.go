package widget

import (
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestRadio_MinSize(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	min := radio.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	radio2 := NewRadio([]string{"Hi", "Hi"}, nil)
	min2 := radio2.MinSize()

	assert.Equal(t, min.Width, min2.Width)
	assert.True(t, min2.Height > min.Height)
}

func TestRadio_Selected(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "Hi", selected)
}

func TestRadio_Unselected(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected)
}

func TestRadio_SelectedOther(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi", "Hi2"}, func(sel string) {
		selected = sel
	})
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), radio.MinSize().Height-theme.Padding())})

	assert.Equal(t, "Hi2", selected)
}

func TestRadio_SelectedNone(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(0, 2)})
	assert.Equal(t, "", selected)

	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(0, 25)})
	assert.Equal(t, "", selected)
}

func TestRadio_Append(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(Renderer(radio).(*radioRenderer).items))

	radio.Options = append(radio.Options, "Another")
	Refresh(radio)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(Renderer(radio).(*radioRenderer).items))
}

func TestRadio_Remove(t *testing.T) {
	radio := NewRadio([]string{"Hi", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(Renderer(radio).(*radioRenderer).items))

	radio.Options = radio.Options[:1]
	Refresh(radio)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(Renderer(radio).(*radioRenderer).items))
}
