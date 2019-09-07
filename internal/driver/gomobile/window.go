package gomobile

import (
	"fyne.io/fyne"
)

type window struct {
	title    string
	visible  bool
	onClosed func()

	clipboard fyne.Clipboard
	canvas    *canvas
	icon      fyne.Resource
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title
}

func (w *window) FullScreen() bool {
	return true
}

func (w *window) SetFullScreen(bool) {
	// no-op
}

func (w *window) Resize(fyne.Size) {
	// no-op
}

func (w *window) RequestFocus() {
	// no-op - we cannot change which window is focused
}

func (w *window) FixedSize() bool {
	return true
}

func (w *window) SetFixedSize(bool) {
	// no-op - all windows are fixed size
}

func (w *window) CenterOnScreen() {
	// no-op
}

func (w *window) Padded() bool {
	return w.canvas.padded
}

func (w *window) SetPadded(padded bool) {
	w.canvas.padded = padded
}

func (w *window) Icon() fyne.Resource {
	if w.icon == nil {
		return fyne.CurrentApp().Icon()
	}

	return w.icon
}

func (w *window) SetIcon(icon fyne.Resource) {
	w.icon = icon
}

func (w *window) MainMenu() *fyne.MainMenu {
	// TODO add mainmenu support for mobile (burger and sidebar?)
	return nil
}

func (w *window) SetMainMenu(*fyne.MainMenu) {
	// TODO add mainmenu support for mobile (burger and sidebar?)
}

func (w *window) SetOnClosed(callback func()) {
	w.onClosed = callback
}

func (w *window) Show() {
	w.visible = true
}

func (w *window) Hide() {
	w.visible = false
}

func (w *window) Close() {
	if w.onClosed != nil {
		w.onClosed()
	}
	//	d := fyne.CurrentApp().Driver().(*mobileDriver)

	// TODO remove from d.windows
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.CurrentApp().Driver().Run()
}

func (w *window) Content() fyne.CanvasObject {
	return w.canvas.Content()
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.canvas.SetContent(content)
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) Clipboard() fyne.Clipboard {
	//if w.clipboard == nil {
	//	w.clipboard = &mobileClipboard{window: w.viewport}
	//}
	// TODO add clipboard support
	return w.clipboard
}

func (w *window) RunWithContext(f func()) {
	//	ctx, _ = e.DrawContext.(gl.Context)

	f()
}

func (w *window) RescaleContext() {
	// TODO
}

func (w *window) Context() interface{} {
	return fyne.CurrentApp().Driver().(*mobileDriver).glctx
}
