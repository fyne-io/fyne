package theme

import (
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestIconThemeChangeName(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	cancel := CancelIcon()
	name := cancel.Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, name, cancel.Name())
}

func TestIconThemeChangeContent(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	cancel := ConfirmIcon()
	content := cancel.Content()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, cancel.Content())
}
