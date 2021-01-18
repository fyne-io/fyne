package test

import (
	"fyne.io/fyne/v2"
)

type testWindow struct {
	title              string
	fullScreen         bool
	fixedSize          bool
	focused            bool
	onClosed           func()
	onCloseIntercepted func()

	canvas    *testCanvas
	clipboard fyne.Clipboard
	driver    *testDriver
	menu      *fyne.MainMenu
}

// NewWindow creates and registers a new window for test purposes
func NewWindow(content fyne.CanvasObject) fyne.Window {
	window := fyne.CurrentApp().NewWindow("")
	window.SetContent(content)
	return window
}

func (w *testWindow) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *testWindow) CenterOnScreen() {
	// no-op
}

func (w *testWindow) Clipboard() fyne.Clipboard {
	return w.clipboard
}

func (w *testWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

func (w *testWindow) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *testWindow) FixedSize() bool {
	return w.fixedSize
}

func (w *testWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *testWindow) Hide() {
	w.focused = false
}

func (w *testWindow) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

func (w *testWindow) MainMenu() *fyne.MainMenu {
	return w.menu
}

func (w *testWindow) Padded() bool {
	return w.canvas.Padded()
}

func (w *testWindow) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*testWindow).focused = false
	}

	w.focused = true
}

func (w *testWindow) Resize(size fyne.Size) {
	w.canvas.Resize(size)
}

func (w *testWindow) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *testWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *testWindow) SetIcon(_ fyne.Resource) {
	// no-op
}

func (w *testWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *testWindow) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

func (w *testWindow) SetMaster() {
	// no-op
}

func (w *testWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *testWindow) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *testWindow) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

func (w *testWindow) SetTitle(title string) {
	w.title = title
}

func (w *testWindow) Show() {
	w.RequestFocus()
}

func (w *testWindow) ShowAndRun() {
	w.Show()
}

func (w *testWindow) Title() string {
	return w.title
}
