package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

var globalRadioRenderer fyne.WidgetRenderer

func BenchmarkRadioCreateRenderer(b *testing.B) {
	var renderer fyne.WidgetRenderer
	widget := &radioItem{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		renderer = widget.CreateRenderer()
	}

	// Avoid having the value optimized out by the compiler.
	globalRadioRenderer = renderer
}

func TestRadioItem_FocusIndicator_Centered_Vertically(t *testing.T) {
	item := newRadioItem("Hello", nil)
	render := test.TempWidgetRenderer(t, item).(*radioItemRenderer)
	render.Layout(fyne.NewSize(200, 100))

	focusIndicatorSize := theme.IconInlineSize() + 2*theme.Padding()
	heightCenterOffset := (100 - focusIndicatorSize) / 2
	assert.Equal(t, fyne.NewPos(theme.Padding()/2, heightCenterOffset), render.focusIndicator.Position1)
}
