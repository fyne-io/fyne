package fyne

import "testing"

func TestDummyApp(t *testing.T) {
	app := NewAppWithDriver(new(dummyDriver))

	app.Quit()
}

type dummyDriver struct {
}

func (d *dummyDriver) CreateWindow(string) Window {
	return nil
}

func (d *dummyDriver) AllWindows() []Window {
	return nil
}

func (d *dummyDriver) RenderedTextSize(text string, size int, _ TextStyle) Size {
	return NewSize(len(text)*size, size)
}

func (d *dummyDriver) Quit() {
	// no-op
}
