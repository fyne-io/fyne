package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelect(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, 2, len(combo.Options))
	assert.Equal(t, "", combo.Selected)
}

func TestSelect_PlaceHolder(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.PlaceHolder = "changed!"
	assert.Equal(t, "changed!", combo.PlaceHolder)
}

func TestSelect_ClearSelected(t *testing.T) {
	const (
		opt1     = "1"
		opt2     = "2"
		optClear = ""
	)
	combo := NewSelect([]string{opt1, opt2}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.SetSelected(opt1)
	assert.Equal(t, opt1, combo.Selected)

	var triggered bool
	var triggeredValue string
	combo.OnChanged = func(s string) {
		triggered = true
		triggeredValue = s
	}
	combo.ClearSelected()
	assert.Equal(t, optClear, combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, optClear, triggeredValue)
}

func TestSelect_SetSelected(t *testing.T) {
	var triggered bool
	var triggeredValue string
	combo := NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
		triggeredValue = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, "2", triggeredValue)
}

func TestSelect_SetSelected_NoChangeOnEmpty(t *testing.T) {
	var triggered bool
	combo := NewSelect([]string{"1", "2"}, func(string) { triggered = true })
	combo.SetSelected("")

	assert.False(t, triggered)
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

func TestSelect_SetSelected_InvalidNoCallback(t *testing.T) {
	var triggered bool
	combo := NewSelect([]string{"1", "2"}, func(string) {
		triggered = true
	})
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_updateSelected(t *testing.T) {
	const (
		opt1 = "1"
		opt2 = "2"
	)

	combo := NewSelect([]string{opt1, opt2}, func(string) {})

	combo.updateSelected(opt2)

	assert.Equal(t, opt2, combo.Selected)
}

func TestSelect_updateSelected_Callback(t *testing.T) {
	const (
		opt1 = "1"
		opt2 = "2"
	)
	var (
		triggered      bool
		triggeredValue string
	)

	combo := NewSelect([]string{opt1, opt2}, func(s string) {
		triggered = true
		triggeredValue = s
	})

	combo.updateSelected(opt2)

	assert.True(t, triggered)
	assert.Equal(t, opt2, triggeredValue)
}

func TestSelect_Tapped(t *testing.T) {
	combo := NewSelect([]string{"1", "2"}, func(s string) {})
	test.Tap(combo)

	canvas := fyne.CurrentApp().Driver().CanvasForObject(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	over := canvas.Overlays().Top()
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(over)

	cont := over.(*PopUp).Content
	assert.Equal(t, cont.Position().X, pos.X)
	assert.Equal(t, cont.Position().Y, pos.Y)

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
	canvas.SetContent(combo)

	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(canvas.Size().Width-10, canvas.Size().Height-10))
	test.Tap(combo)

	comboPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	cont := canvas.Overlays().Top().(*PopUp).Content
	assert.Less(t, cont.Position().Y, comboPos.Y, "window too small so we render higher up")
	assert.Less(t, cont.Position().X, comboPos.X, "window too small so we render to the left")
}

func TestSelectRenderer_ApplyTheme(t *testing.T) {
	sel := &Select{}
	render := test.WidgetRenderer(sel).(*selectRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	withTestTheme(func() {
		render.Refresh()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func TestSelect_Move(t *testing.T) {
	// fresh app for this test
	app := test.NewApp()
	// don't let our app hang around for too long
	defer test.NewApp()

	combo := NewSelect([]string{"1", "2"}, nil)
	canvas := app.Driver().CanvasForObject(combo)
	canvas.(test.WindowlessCanvas).Resize(fyne.NewSize(100, 100))
	canvas.SetContent(combo)

	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	require.Equal(t, fyne.NewPos(10, 10), combo.Position())

	combo.Tapped(&fyne.PointEvent{})
	require.Equal(t, fyne.NewPos(10, 39), combo.popUp.innerPos)

	combo.Move(fyne.NewPos(20, 20))
	assert.Equal(t, fyne.NewPos(20, 49), combo.popUp.innerPos)
}
