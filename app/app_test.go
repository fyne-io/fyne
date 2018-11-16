package app

import (
	"github.com/fyne-io/fyne"
	"testing"
)

func TestDummyApp(t *testing.T) {
	app := NewAppWithDriver(new(dummyDriver))

	app.Quit()
}

type dummyDriver struct {
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
