package software

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
)

type softwareWindow struct {
	common.Window

	title              string
	fullScreen         bool
	fixedSize          bool
	focused            bool
	onClosed           func()
	onCloseIntercepted func()

	canvas    *softwareCanvas
	clipboard fyne.Clipboard
	driver    *softwareDriver
	menu      *fyne.MainMenu
}

// NewWindow creates and registers a new window for test purposes
func NewWindow(content fyne.CanvasObject) fyne.Window {
	window := fyne.CurrentApp().NewWindow("")
	window.SetContent(content)
	return window
}

func (w *softwareWindow) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *softwareWindow) CenterOnScreen() {
	// no-op
}

func (w *softwareWindow) Clipboard() fyne.Clipboard {
	return w.clipboard
}

func (w *softwareWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

func (w *softwareWindow) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *softwareWindow) FixedSize() bool {
	return w.fixedSize
}

func (w *softwareWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *softwareWindow) Hide() {
	w.focused = false
}

func (w *softwareWindow) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

func (w *softwareWindow) MainMenu() *fyne.MainMenu {
	return w.menu
}

func (w *softwareWindow) Padded() bool {
	return w.canvas.Padded()
}

func (w *softwareWindow) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*softwareWindow).focused = false
	}

	w.focused = true
}

func (w *softwareWindow) Resize(size fyne.Size) {
	w.canvas.Resize(size)
}

func (w *softwareWindow) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *softwareWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *softwareWindow) SetIcon(_ fyne.Resource) {
	// no-op
}

func (w *softwareWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *softwareWindow) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

func (w *softwareWindow) SetMaster() {
	// no-op
}

func (w *softwareWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *softwareWindow) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *softwareWindow) SetOnDropped(dropped func(fyne.Position, []fyne.URI)) {

}

func (w *softwareWindow) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

func (w *softwareWindow) SetTitle(title string) {
	w.title = title
}

func (w *softwareWindow) Show() {
	w.RequestFocus()
}

func (w *softwareWindow) ShowAndRun() {
	w.Show()
	w.driver.Run()
}

func (w *softwareWindow) Title() string {
	return w.title
}
