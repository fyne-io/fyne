package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestNewSelect(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, 2, len(combo.Options))
	assert.Equal(t, "", combo.Selected)
}

func TestSelect_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindowWithPainter(combo, software.NewPainter())
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	test.Tap(combo)

	test.AssertImageMatches(t, "select_theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		combo.Refresh()
		time.Sleep(100 * time.Millisecond)
		// Looks weird but the test infrastructure does not adjust min-sizes.
		test.AssertImageMatches(t, "select_theme_changed.png", w.Canvas().Capture())
	})
}

func TestSelect_ClearSelected(t *testing.T) {
	const (
		opt1     = "1"
		opt2     = "2"
		optClear = ""
	)
	combo := widget.NewSelect([]string{opt1, opt2}, func(string) {})
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

func TestSelect_Move(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, nil)
	w := test.NewWindowWithPainter(combo, software.NewPainter())
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))

	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "select_move_initial.png", w.Canvas().Capture())

	combo.Tapped(&fyne.PointEvent{})
	test.AssertImageMatches(t, "select_move_tapped.png", w.Canvas().Capture())

	combo.Move(fyne.NewPos(20, 20))
	test.AssertImageMatches(t, "select_move_moved.png", w.Canvas().Capture())
}

func TestSelect_PlaceHolder(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.PlaceHolder = "changed!"
	assert.Equal(t, "changed!", combo.PlaceHolder)
}

func TestSelect_SetSelected(t *testing.T) {
	var triggered bool
	var triggeredValue string
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
		triggeredValue = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, "2", triggeredValue)
}

func TestSelect_SetSelected_Callback(t *testing.T) {
	selected := ""
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		selected = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", selected)
}

func TestSelect_SetSelected_Invalid(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("3")

	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected_InvalidNoCallback(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {
		triggered = true
	})
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_SetSelected_InvalidReplace(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")
	combo.SetSelected("3")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_NoChangeOnEmpty(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(string) { triggered = true })
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_Tapped(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindowWithPainter(combo, software.NewPainter())
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())

	test.Tap(combo)
	canvas := fyne.CurrentApp().Driver().CanvasForObject(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	test.AssertImageMatches(t, "select_tapped.png", w.Canvas().Capture())
}

func TestSelect_Tapped_Constrained(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindowWithPainter(combo, software.NewPainter())
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())

	canvas := w.Canvas()
	combo.Move(fyne.NewPos(canvas.Size().Width-10, canvas.Size().Height-10))
	test.Tap(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	test.AssertImageMatches(t, "select_tapped_constrained.png", w.Canvas().Capture())
}
