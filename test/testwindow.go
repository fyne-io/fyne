package test

import "github.com/fyne-io/fyne"

type testWindow struct {
	title      string
	fullscreen bool

	canvas fyne.Canvas
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

func (w *testWindow) Close() {
	i := 0
	for _, window := range windows {
		if window == w {
			break
		}
		i++
	}

	windows = append(windows[:i], windows[i+1:]...)
}

func (w *testWindow) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *testWindow) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *testWindow) Canvas() fyne.Canvas {
	return w.canvas
}

// NewTestWindow creates and registers a new window for test purposes
func NewTestWindow(content fyne.CanvasObject) fyne.Window {
	canvas := NewTestCanvas()
	canvas.SetContent(content)
	window := &testWindow{canvas: canvas}

	windows = append(windows, window)
	return window
}
