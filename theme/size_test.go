package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func Test_TextSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNameText), TextSize(), "wrong text size")
}

func Test_Padding(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNamePadding), Padding(), "wrong padding")
}

func Test_IconInlineSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNameInlineIcon), IconInlineSize(), "wrong inline icon size")
}

func Test_ScrollBarSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Size(SizeNameScrollBar), ScrollBarSize(), "wrong inline icon size")
}
