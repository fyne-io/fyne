package app

import (
	"testing"
	"time"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
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

func TestDoubleRun(t *testing.T) {
	app := New()

	go func() {
		time.Sleep(100 * time.Millisecond)

		app.Run() // this run is erroneous, but let's avoid breakage if users do it
		app.Quit()
	}()

	app.NewWindow("testing") // just a dummy window to ensure that any driver contexts are initialised
	app.Run()
}
