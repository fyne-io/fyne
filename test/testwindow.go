package test

import (
	"fyne.io/fyne"
)

type testWindow struct {
	title      string
	fullScreen bool
	fixedSize  bool
	focused    bool
	onClosed   func()

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

// Canvas satisfies the fyne.Window interface.
func (w *testWindow) Canvas() fyne.Canvas {
	return w.canvas
}

// CenterOnScreen satisfies the fyne.Window interface.
func (w *testWindow) CenterOnScreen() {
	// no-op
}

// Clipboard satisfies the fyne.Window interface.
func (w *testWindow) Clipboard() fyne.Clipboard {
	return w.clipboard
}

// Close satisfies the fyne.Window interface.
func (w *testWindow) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	w.focused = false
	w.driver.removeWindow(w)
}

// Content satisfies the fyne.Window interface.
func (w *testWindow) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

// FixedSize satisfies the fyne.Window interface.
func (w *testWindow) FixedSize() bool {
	return w.fixedSize
}

// FullScreen satisfies the fyne.Window interface.
func (w *testWindow) FullScreen() bool {
	return w.fullScreen
}

// Hide satisfies the fyne.Window interface.
func (w *testWindow) Hide() {
	w.focused = false
}

// Icon satisfies the fyne.Window interface.
func (w *testWindow) Icon() fyne.Resource {
	return fyne.CurrentApp().Icon()
}

// MainMenu satisfies the fyne.Window interface.
func (w *testWindow) MainMenu() *fyne.MainMenu {
	return w.menu
}

// Padded satisfies the fyne.Window interface.
func (w *testWindow) Padded() bool {
	return w.canvas.Padded()
}

// RequestFocus satisfies the fyne.Window interface.
func (w *testWindow) RequestFocus() {
	for _, win := range w.driver.AllWindows() {
		win.(*testWindow).focused = false
	}

	w.focused = true
}

// Resize satisfies the fyne.Window interface.
func (w *testWindow) Resize(size fyne.Size) {
	w.canvas.Resize(size)
}

// SetContent satisfies the fyne.Window interface.
func (w *testWindow) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

// SetFixedSize satisfies the fyne.Window interface.
func (w *testWindow) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

// SetIcon satisfies the fyne.Window interface.
func (w *testWindow) SetIcon(icon fyne.Resource) {
	// no-op
}

// SetFullScreen satisfies the fyne.Window interface.
func (w *testWindow) SetFullScreen(fullScreen bool) {
	w.fullScreen = fullScreen
}

// SetMainMenu satisfies the fyne.Window interface.
func (w *testWindow) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

// SetMaster satisfies the fyne.Window interface.
func (w *testWindow) SetMaster() {
	// no-op
}

// SetOnClosed satisfies the fyne.Window interface.
func (w *testWindow) SetOnClosed(closed func()) {
	w.onClosed = closed
}

// SetPadded satisfies the fyne.Window interface.
func (w *testWindow) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)
}

// SetTitle satisfies the fyne.Window interface.
func (w *testWindow) SetTitle(title string) {
	w.title = title
}

// Show satisfies the fyne.Window interface.
func (w *testWindow) Show() {
	w.RequestFocus()
}

// ShowAndRun satisfies the fyne.Window interface.
func (w *testWindow) ShowAndRun() {
	w.Show()
}

// Title satisfies the fyne.Window interface.
func (w *testWindow) Title() string {
	return w.title
}
