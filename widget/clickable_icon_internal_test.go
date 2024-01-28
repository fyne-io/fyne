package widget

import (
	"testing"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewClickableIcon(t *testing.T) {
	icon := NewClickableIcon(theme.ConfirmIcon(), func() {})
	render := test.WidgetRenderer(icon)

	assert.Equal(t, 1, len(render.Objects()))
	obj := render.Objects()[0]
	img, ok := obj.(*canvas.Image)
	if !ok {
		t.Fail()
	}
	assert.Equal(t, theme.ConfirmIcon(), img.Resource)
}

func TestClickableIcon_Nil(t *testing.T) {
	icon := NewClickableIcon(nil, func() {})
	render := test.WidgetRenderer(icon)

	assert.Equal(t, 1, len(render.Objects()))
	assert.Nil(t, render.Objects()[0].(*canvas.Image).Resource)
}

func TestClickableIcon_MinSize(t *testing.T) {
	icon := NewClickableIcon(theme.CancelIcon(), func() {})
	min := icon.MinSize()

	assert.Equal(t, theme.IconInlineSize(), min.Width)
	assert.Equal(t, theme.IconInlineSize(), min.Height)
}
