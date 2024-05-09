package software

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
)

type SoftwareWindow struct {
	common.Window

	title              string
	fullScreen         bool
	fixedSize          bool
	focused            bool
	onClosed           func()
	onCloseIntercepted func()

	canvas    *SoftwareCanvas
	clipboard fyne.Clipboard
	driver    *SoftwareDriver
	menu      *fyne.MainMenu
}

// NewWindow creates and registers a new window for test purposes
func NewWindow(content fyne.CanvasObject) fyne.Window {
	window := fyne.CurrentApp().NewWindow("")
	window.SetContent(content)
	window.(*SoftwareWindow).clipboard = NewClipboard()
	return window
}

func (w *SoftwareWindow) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *SoftwareWindow) CenterOnScreen() {
	// no-op
}

func (w *SoftwareWindow) Clipboard() fyne.Clipboard {
	return w.clipboard
}

func (w *SoftwareWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

func (w *SoftwareWindow) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *SoftwareWindow) FixedSize() bool {
	return w.fixedSize
}

func (w *SoftwareWindow) FullScreen() bool {
	return w.fullScreen
}

func (w *SoftwareWindow) Hide() {
	w.focused = false
}

func (w *SoftwareWindow) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

func (w *SoftwareWindow) MainMenu() *fyne.MainMenu {
	return w.menu
}

func (w *SoftwareWindow) Padded() bool {
	return w.canvas.Padded()
}

func (w *SoftwareWindow) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*SoftwareWindow).focused = false
	}

	w.focused = true
}

func (w *SoftwareWindow) Resize(size fyne.Size) {
	w.canvas.Resize(size)
}

func (w *SoftwareWindow) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *SoftwareWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *SoftwareWindow) SetIcon(_ fyne.Resource) {
	// no-op
}

func (w *SoftwareWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

func (w *SoftwareWindow) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

func (w *SoftwareWindow) SetMaster() {
	// no-op
}

func (w *SoftwareWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *SoftwareWindow) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *SoftwareWindow) SetOnDropped(dropped func(fyne.Position, []fyne.URI)) {

}

func (w *SoftwareWindow) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

func (w *SoftwareWindow) SetTitle(title string) {
	w.title = title
}

func (w *SoftwareWindow) Show() {
	w.RequestFocus()
}

func (w *SoftwareWindow) ShowAndRun() {
	w.Show()
	w.driver.Run()
}

func (w *SoftwareWindow) Title() string {
	return w.title
}
