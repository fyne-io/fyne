package theme

import (
	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewThemeChangesIconColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	ic1 := fyne.CurrentApp().Settings().Theme().IconColor()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	ic2 := fyne.CurrentApp().Settings().Theme().IconColor()
	assert.NotEqual(t, ic1, ic2)
}
