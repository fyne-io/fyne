package widget

import (
	"testing"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewIcon(t *testing.T) {
	icon := NewIcon(theme.ConfirmIcon())
	render := Renderer(icon)

	assert.Equal(t, 1, len(render.Objects()))
	obj := render.Objects()[0]
	img, ok := obj.(*canvas.Image)
	if !ok {
		t.Fail()
	}
	assert.Equal(t, theme.ConfirmIcon(), img.Resource)
}

func TestIcon_Nil(t *testing.T) {
	icon := NewIcon(nil)
	render := Renderer(icon)

	assert.Equal(t, 0, len(render.Objects()))
}

func TestIcon_MinSize(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	min := icon.MinSize()

	assert.Equal(t, theme.IconInlineSize(), min.Width)
	assert.Equal(t, theme.IconInlineSize(), min.Height)
}

func TestIconRenderer_ApplyTheme(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	icon.Hide()
	render := Renderer(icon).(*iconRenderer)

	assert.False(t, render.objects[0].Visible())

	render.ApplyTheme()
	assert.False(t, render.objects[0].Visible())
}
