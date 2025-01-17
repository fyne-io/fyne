package test

import (
	"testing"

	"fyne.io/fyne/v2"
)

type window struct {
	title              string
	fullScreen         bool
	fixedSize          bool
	focused            bool
	onClosed           func()
	onCloseIntercepted func()

	canvas *canvas
	driver *driver
	menu   *fyne.MainMenu
}

// NewTempWindow creates and registers a new window for test purposes.
// This window will get removed automatically once the running test ends.
//
// Since: 2.5
func NewTempWindow(t testing.TB, content fyne.CanvasObject) fyne.Window {
	window := NewWindow(content)
	t.Cleanup(window.Close)
	return window
}

// NewWindow creates and registers a new window for test purposes
func NewWindow(content fyne.CanvasObject) fyne.Window {
	window := fyne.CurrentApp().NewWindow("")
	window.SetContent(content)
	return window
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) CenterOnScreen() {
	// no-op
}

func (w *window) Clipboard() fyne.Clipboard {
	return NewClipboard()
}

func (w *window) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

func (w *window) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

func (w *window) Hide() {
	w.focused = false
}

func (w *window) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

func (w *window) MainMenu() *fyne.MainMenu {
	return w.menu
}

func (w *window) Padded() bool {
	return w.canvas.Padded()
}

func (w *window) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*window).focused = false
	}

	w.focused = true
}

func (w *window) Resize(size fyne.Size) {
	w.canvas.Resize(size)
}

func (w *window) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *window) SetIcon(_ fyne.Resource) {
	// no-op
}

func (w *window) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *window) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

func (w *window) SetMaster() {
	// no-op
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *window) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *window) SetOnDropped(dropped func(fyne.Position, []fyne.URI)) {

}

func (w *window) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

func (w *window) SetTitle(title string) {
	w.title = title
}

func (w *window) Show() {
	w.RequestFocus()
}

func (w *window) ShowAndRun() {
	w.Show()
}

func (w *window) Title() string {
	return w.title
}
