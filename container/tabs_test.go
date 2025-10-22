package container

import (
	"testing"

	internalTest "fyne.io/fyne/v2/internal/test"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
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
	a := test.NewTempApp(t)
	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))

	tabs := NewAppTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")))
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

func TestTab_RemoveCallsHandlers(t *testing.T) {
	tabs := NewDocTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")),
	)

	selectedText := ""
	unSelectedText := ""
	tabs.OnSelected = func(ti *TabItem) {
		selectedText = ti.Text
	}
	tabs.OnUnselected = func(ti *TabItem) {
		unSelectedText = ti.Text
	}

	tabs.RemoveIndex(0)

	assert.Equal(t, "b", selectedText)
	assert.Equal(t, "a", unSelectedText)
}

func TestTab_RemoveNonSelectedTab(t *testing.T) {
	tabs := NewDocTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")),
	)

	selectedText := ""
	unSelectedText := ""
	tabs.OnSelected = func(ti *TabItem) {
		selectedText = ti.Text
	}
	tabs.OnUnselected = func(ti *TabItem) {
		unSelectedText = ti.Text
	}

	tabs.RemoveIndex(1)

	assert.Equal(t, "a", selectedText)
	assert.Equal(t, "", unSelectedText)
}

func TestTab_RemoveLastTab(t *testing.T) {
	tabs := NewDocTabs(
		NewTabItem("a", widget.NewLabel("a")),
	)
	selectedText := ""
	tabs.OnUnselected = func(ti *TabItem) {
		selectedText = ti.Text
	}

	tabs.RemoveIndex(0)

	assert.Equal(t, "a", selectedText)
}
