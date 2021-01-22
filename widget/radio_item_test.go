package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestRadioItem_FocusIndicator_Centered_Vertically(t *testing.T) {
	item := newRadioItem("Hello", nil)
	render := test.WidgetRenderer(item).(*radioItemRenderer)
	render.Layout(fyne.NewSize(200, 100))

	focusIndicatorSize := theme.IconInlineSize() + theme.Padding()*2
	heightCenterOffset := (100-focusIndicatorSize)/2 + theme.Padding()
	assert.Equal(t, fyne.NewPos(theme.Padding()/2, heightCenterOffset), render.focusIndicator.Position1)
}
