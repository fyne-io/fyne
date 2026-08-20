package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestNewIcon(t *testing.T) {
	icon := NewIcon(theme.ConfirmIcon())
	render := test.TempWidgetRenderer(t, icon)

	assert.Len(t, render.Objects(), 1)
	obj := render.Objects()[0]
	img, ok := obj.(*canvas.Image)
	if !ok {
		t.Fail()
	}
	assert.Equal(t, theme.ConfirmIcon(), img.Resource)
}

func TestIcon_Nil(t *testing.T) {
	icon := NewIcon(nil)
	render := test.TempWidgetRenderer(t, icon)

	assert.Len(t, render.Objects(), 1)
	assert.Nil(t, render.Objects()[0].(*canvas.Image).Resource)
}

func TestIcon_MinSize(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	minSize := icon.MinSize()

	assert.Equal(t, theme.IconInlineSize(), minSize.Width)
	assert.Equal(t, theme.IconInlineSize(), minSize.Height)
}

func TestIconRenderer_ApplyTheme(t *testing.T) {
	icon := NewIcon(theme.CancelIcon())
	render := test.TempWidgetRenderer(t, icon).(*iconRenderer)
	visible := render.Objects()[0].Visible()

	render.Refresh()
	assert.Equal(t, visible, render.Objects()[0].Visible())
}
