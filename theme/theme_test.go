package theme

import "github.com/fyne-io/fyne"

import "testing"

import "github.com/stretchr/testify/assert"

func TestThemeChange(t *testing.T) {
	fyne.GlobalSettings().SetTheme("dark")
	bg := BackgroundColor()

	fyne.GlobalSettings().SetTheme("light")
	assert.NotEqual(t, bg, BackgroundColor())
}
