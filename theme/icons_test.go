package theme

import "github.com/fyne-io/fyne"

import "testing"

import "github.com/stretchr/testify/assert"

func TestIconThemeChangeName(t *testing.T) {
	fyne.GetSettings().SetTheme("dark")
	cancel := CancelIcon()
	name := cancel.Name()

	fyne.GetSettings().SetTheme("light")
	assert.NotEqual(t, name, cancel.Name())
}

func TestIconThemeChangeContent(t *testing.T) {
	fyne.GetSettings().SetTheme("dark")
	cancel := ConfirmIcon()
	content := cancel.Content()

	fyne.GetSettings().SetTheme("light")
	assert.NotEqual(t, content, cancel.Content())
}

func TestIconThemeChangePath(t *testing.T) {
	fyne.GetSettings().SetTheme("dark")
	checked := CheckedIcon()
	path := checked.CachePath()

	fyne.GetSettings().SetTheme("light")
	assert.NotEqual(t, path, checked.CachePath())
}
