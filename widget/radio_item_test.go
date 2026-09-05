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

func TestRadioItem_AccessibilityRoleAndLabel(t *testing.T) {
	item := newRadioItem("Hello", nil)

	assert.Equal(t, fyne.AccessibleRoleRadio, item.AccessibilityRole())
	assert.Equal(t, "Hello", item.AccessibilityLabel())
	assert.Equal(t,
		[]fyne.AccessibleAction{fyne.AccessibleActionPress, fyne.AccessibleActionSelect},
		item.AccessibilityActions())
}

func TestRadioItem_AccessibilityStates(t *testing.T) {
	item := newRadioItem("Hello", nil)
	assert.Empty(t, item.AccessibilityStates())

	item.Selected = true
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateSelected}, item.AccessibilityStates())

	item.Disable()
	assert.Equal(t,
		[]fyne.AccessibleState{fyne.AccessibleStateSelected, fyne.AccessibleStateDisabled},
		item.AccessibilityStates())
}

func TestRadioItem_AccessibilityPerformAction_TriggersTap(t *testing.T) {
	tapped := 0
	item := newRadioItem("Hello", func(*radioItem) { tapped++ })

	assert.True(t, item.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.True(t, item.AccessibilityPerformAction(fyne.AccessibleActionSelect))
	assert.Equal(t, 2, tapped)

	assert.False(t, item.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
	assert.Equal(t, 2, tapped)

	item.Disable()
	assert.False(t, item.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 2, tapped)
}
