package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	intTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestTabButton_Icon_Change(t *testing.T) {
	b := &tabButton{icon: theme.CancelIcon()}
	r := cache.Renderer(b)
	icon := r.Objects()[3].(*canvas.Image)
	oldResource := icon.Resource

	b.icon = theme.ConfirmIcon()
	b.Refresh()
	assert.NotEqual(t, oldResource, icon.Resource)
}

func TestTab_ThemeChange(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	a.Settings().SetTheme(test.Theme())

	tabs := NewAppTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")))
	w := test.NewWindow(tabs)
	w.Resize(fyne.NewSquareSize(150))

	initial := w.Canvas().Capture()

	a.Settings().SetTheme(intTest.DarkTheme())
	tabs.SelectIndex(1)
	second := w.Canvas().Capture()
	assert.NotEqual(t, initial, second)

	a.Settings().SetTheme(test.Theme())
	tabs.SelectIndex(0)
	assert.Equal(t, initial, w.Canvas().Capture())
}
