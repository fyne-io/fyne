package app

import (
	"testing"

	"github.com/fyne-io/fyne"
)

func TestDummyApp(t *testing.T) {
	app := NewAppWithDriver(new(dummyDriver))

	app.Quit()
}

// TODO figure how we can re-instate this
//func TestThemeFromEnv(t *testing.T) {
//	settingsCache = nil
//	os.Setenv("FYNE_THEME", "light")

//	assert.Equal(t, "light", fyne.GlobalSettings().Theme())
//}

type dummyDriver struct {
}

func (d *dummyDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	return nil
}

func (d *dummyDriver) CreateWindow(string) fyne.Window {
	return nil
}

func (d *dummyDriver) AllWindows() []fyne.Window {
	return nil
}

func (d *dummyDriver) RenderedTextSize(text string, size int, _ fyne.TextStyle) fyne.Size {
	return fyne.NewSize(len(text)*size, size)
}

func (d *dummyDriver) Run() {
	// no-op
}

func (d *dummyDriver) Quit() {
	// no-op
}
