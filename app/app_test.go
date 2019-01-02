package app

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestDummyApp(t *testing.T) {
	app := New()

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := New()

	assert.Equal(t, app, fyne.CurrentApp())
}

// TODO figure how we can re-instate this
//func TestThemeFromEnv(t *testing.T) {
//	settingsCache = nil
//	os.Setenv("FYNE_THEME", "light")

//	assert.Equal(t, "light", fyne.GlobalSettings().Theme())
//}
