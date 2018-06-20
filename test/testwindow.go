package test

import "github.com/fyne-io/fyne"

type testWindow struct {
	title      string
	fullscreen bool
}

var windows = make([]fyne.Window, 0)

func (w *testWindow) Title() string {
	return w.title
}

func (w *testWindow) SetTitle(title string) {
	w.title = title
}

func (w *testWindow) Fullscreen() bool {
	return w.fullscreen
}

func (w *testWindow) SetFullscreen(fullscreen bool) {
	w.fullscreen = fullscreen
}

func (w *testWindow) Show() {}

func (w *testWindow) Hide() {}

func (w *testWindow) Close() {}

func (w *testWindow) Canvas() fyne.Canvas {
	return GetTestCanvas()
}

// NewTestWindow creates and registers a new window for test purposes
func NewTestWindow() fyne.Window {
	window := &testWindow{}
	windows = append(windows, window)

	return window
}
