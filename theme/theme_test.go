package theme

import "github.com/fyne-io/fyne"

import "testing"

import "github.com/stretchr/testify/assert"

func TestThemeChange(t *testing.T) {
	bg := BackgroundColor()

	fyne.GetSettings().SetTheme("light")
	assert.NotEqual(t, bg, BackgroundColor())
}
