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

	focusIndicatorSize := theme.IconInlineSize() + 2*theme.Padding()
	heightCenterOffset := (100 - focusIndicatorSize) / 2
	assert.Equal(t, fyne.NewPos(theme.Padding()/2, heightCenterOffset), render.focusIndicator.Position1)
}

func TestHorizontalRadioGroupWidths(t *testing.T) {
	radio := &RadioGroup{Options: []string{"A", "Extra Long", "B"}, Horizontal: true}
	rend := radio.CreateRenderer().(*radioGroupRenderer)

	maxMinItemWidth := float32(0)
	for _, i := range rend.items {
		minWidth := i.MinSize().Width
		maxMinItemWidth = fyne.Max(maxMinItemWidth, minWidth)
	}

	radioMinWidth := rend.MinSize().Width

	assert.GreaterOrEqual(t, radioMinWidth, float32(len(radio.Options))*maxMinItemWidth)
}
