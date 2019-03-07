package test

import (
	"sync"

	"fyne.io/fyne"
)

type testWindow struct {
	title      string
	fullScreen bool
	fixedSize  bool
	padded     bool
	focused    bool
	onClosed   func()

	canvas    fyne.Canvas
	clipboard fyne.Clipboard
}

var windows = make([]fyne.Window, 0)
var windowsMutex = sync.RWMutex{}

func (w *testWindow) Title() string {
	return w.title
}

func (w *testWindow) SetTitle(title string) {
	w.title = title
}

func (w *testWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *testWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *testWindow) CenterOnScreen() {
	// no-op
}

func (w *testWindow) Resize(fyne.Size) {
	// no-op
}

func (w *testWindow) RequestFocus() {
	for _, win := range windows {
		win.(*testWindow).focused = false
	}

	w.focused = true
}

func (w *testWindow) FixedSize() bool {
	return w.fixedSize
}

func (w *testWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *testWindow) Padded() bool {
	return w.padded
}

func (w *testWindow) SetPadded(padded bool) {
	w.padded = padded
}

func (w *testWindow) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

func (w *testWindow) SetIcon(icon fyne.Resource) {
	// no-op
}

func (w *testWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *testWindow) Show() {
	w.RequestFocus()
}

func (w *testWindow) Clipboard() fyne.Clipboard {
	return w.clipboard
}

func (w *testWindow) Hide() {
	w.focused = false
}

func (w *testWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false

	windowsMutex.Lock()
	i := 0
	for _, window := range windows {
		if window == w {
			break
		}
		i++
	}

	windows = append(windows[:i], windows[i+1:]...)
	windowsMutex.Unlock()
}

func (w *testWindow) ShowAndRun() {
	w.Show()
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

// NewWindow creates and registers a new window for test purposes
func NewWindow(content fyne.CanvasObject) fyne.Window {
	canvas := NewCanvas()
	canvas.SetContent(content)
	window := &testWindow{canvas: canvas}
	window.clipboard = &testClipboard{}

	windowsMutex.Lock()
	windows = append(windows, window)
	windowsMutex.Unlock()
	return window
}
