package gomobile

import (
	"fyne.io/fyne"
)

type window struct {
	title           string
	padded, visible bool

	canvas fyne.Canvas
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
	panic("implement me")
}

func (w *window) FixedSize() bool {
	return true
}

func (w *window) SetFixedSize(bool) {
	// no-op
}

func (w *window) CenterOnScreen() {
	// no-op
}

func (w *window) Padded() bool {
	return w.padded
}

func (w *window) SetPadded(padded bool) {
	w.padded = padded
}

func (w *window) Icon() fyne.Resource {
	panic("implement me")
}

func (w *window) SetIcon(fyne.Resource) {
	//	panic("implement me")
}

func (w *window) MainMenu() *fyne.MainMenu {
	panic("implement me")
}

func (w *window) SetMainMenu(*fyne.MainMenu) {
	//	panic("implement me")
}

func (w *window) SetOnClosed(func()) {
	panic("implement me")
}

func (w *window) Show() {
	w.visible = true
}

func (w *window) Hide() {
	w.visible = false
}

func (w *window) Close() {
	//	d := fyne.CurrentApp().Driver().(*driver)

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
	panic("implement me")
}

func (w *window) RunWithContext(f func()) {
	//	ctx, _ = e.DrawContext.(gl.Context)

	f()
}

func (w *window) RescaleContext() {
	// TODO
}

func (w *window) Context() interface{} {
	return fyne.CurrentApp().Driver().(*driver).glctx
}
