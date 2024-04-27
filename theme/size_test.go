package theme_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func Test_IconInlineSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	assert.Equal(t, theme.DarkTheme().Size(theme.SizeNameInlineIcon), theme.IconInlineSize(), "wrong inline icon size")
}

func Test_Padding(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	assert.Equal(t, theme.DarkTheme().Size(theme.SizeNamePadding), theme.Padding(), "wrong padding")
}

func Test_ScrollBarSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	assert.Equal(t, theme.DarkTheme().Size(theme.SizeNameScrollBar), theme.ScrollBarSize(), "wrong inline icon size")
}

func Test_TextSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	assert.Equal(t, theme.DarkTheme().Size(theme.SizeNameText), theme.TextSize(), "wrong text size")
}
