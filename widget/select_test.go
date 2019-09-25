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

	cont := over.(*PopUp).Content
	assert.Equal(t, cont.Position().X, pos.X+theme.Padding())
	assert.True(t, cont.Position().Y > pos.Y)

	items := cont.(*Box).Children
	assert.Equal(t, 2, len(items))
}

func TestSelect_Tapped_Constrained(t *testing.T) {
	// fresh app for this test
	test.NewApp()
	// don't let our app hang around for too long
	defer test.NewApp()

	combo := NewSelect([]string{"1", "2"}, func(s string) {})
	canvas := fyne.CurrentApp().Driver().CanvasForObject(combo)
	canvas.(test.WindowlessCanvas).Resize(fyne.NewSize(100, 100))

	combo.Move(fyne.NewPos(canvas.Size().Width-10, canvas.Size().Height-10))
	test.Tap(combo)

	comboPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(combo)
	cont := canvas.Overlay().(*PopUp).Content
	assert.Less(t, cont.Position().Y, comboPos.Y, "window too small so we render higher up")
	assert.Less(t, cont.Position().X, comboPos.X, "window too small so we render to the left")
}

func TestSelectRenderer_ApplyTheme(t *testing.T) {
	sel := &Select{}
	render := Renderer(sel).(*selectRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	withTestTheme(func() {
		render.ApplyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
