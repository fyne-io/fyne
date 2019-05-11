package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewSelect(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, 2, len(combo.Options))
	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_Invalid(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("3")

	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected_InvalidReplace(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")
	combo.SetSelected("3")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_Callback(t *testing.T) {
	selected := ""
	combo := NewSelect([]string{"1", "2"}, func(s string) {
		selected = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", selected)
}

func TestSelect_Tapped(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(s string) {})
	test.Tap(combo)

	over := fyne.CurrentApp().Driver().CanvasForObject(combo).Overlay()
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(over)
	assert.NotNil(t, over)

	cont := over.(*PopOver).Content
	assert.True(t, cont.Position().X > pos.X)
	assert.True(t, cont.Position().Y > pos.Y)

	items := cont.(*Box).Children
	assert.Equal(t, 2, len(items))
}

func TestSelect_Tapped_Constrained(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(s string) {})
	win := test.NewWindow(combo)
	win.Resize(fyne.NewSize(20, 20))
	test.Tap(combo)

	over := fyne.CurrentApp().Driver().CanvasForObject(combo).Overlay()
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(over)

	cont := over.(*PopOver).Content
	assert.True(t, cont.Position().Y <= pos.Y+theme.Padding()) // window was too small so we render higher up
	assert.True(t, cont.Position().X > pos.X)                  // but X position is unaffected
}
