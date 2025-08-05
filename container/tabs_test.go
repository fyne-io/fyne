package container_test

import (
	"testing"

	"fyne.io/fyne/v2/container"
	internalTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestTab_ThemeChange(t *testing.T) {
	a := test.NewTempApp(t)
	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))

	tabs := container.NewAppTabs(
		container.NewTabItem("a", widget.NewLabel("a")),
		container.NewTabItem("b", widget.NewLabel("b")))
	w := test.NewTempWindow(t, tabs)
	w.Resize(fyne.NewSize(180, 120))

	initial := w.Canvas().Capture()

	a.Settings().SetTheme(internalTest.DarkTheme(theme.DefaultTheme()))
	tabs.SelectIndex(1)
	second := w.Canvas().Capture()
	assert.NotEqual(t, initial, second)

	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))
	tabs.SelectIndex(0)
	assert.Equal(t, initial, w.Canvas().Capture())
}
