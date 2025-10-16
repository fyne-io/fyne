package embedded

import (
	"fyne.io/fyne/v2"
)

type noosWindow struct {
	c *touchCanvas
	d *noosDriver

	title string
}

func (w *noosWindow) Title() string {
	return w.title
}

func (w *noosWindow) SetTitle(s string) {
	w.title = s
}

func (w *noosWindow) FullScreen() bool {
	return true
}

func (w *noosWindow) SetFullScreen(_ bool) {
}

func (w *noosWindow) Resize(s fyne.Size) {
	w.c.Resize(s)
}

func (w *noosWindow) RequestFocus() {
	// TODO implement me
	panic("implement me")
}

func (w *noosWindow) FixedSize() bool {
	return true
}

func (w *noosWindow) SetFixedSize(bool) {}

func (w *noosWindow) CenterOnScreen() {}

func (w *noosWindow) Padded() bool {
	return w.c.Padded()
}

func (w *noosWindow) SetPadded(pad bool) {
	w.c.SetPadded(pad)
}

func (w *noosWindow) Icon() fyne.Resource {
	// TODO implement me
	return nil
}

func (w *noosWindow) SetIcon(fyne.Resource) {
	// TODO implement me
}

func (w *noosWindow) SetMaster() {
	// TODO implement me
}

func (w *noosWindow) MainMenu() *fyne.MainMenu {
	// TODO implement me
	return nil
}

func (w *noosWindow) SetMainMenu(menu *fyne.MainMenu) {
	// TODO implement me
}

func (w *noosWindow) SetOnClosed(f func()) {
	// TODO implement me
}

func (w *noosWindow) SetCloseIntercept(f func()) {
	// TODO implement me
}

func (w *noosWindow) SetOnDropped(func(fyne.Position, []fyne.URI)) {}

func (w *noosWindow) Show() {
	w.d.renderWindow(w)
}

func (w *noosWindow) Hide() {}

func (w *noosWindow) Close() {
	i := -1
	for _, win := range w.d.wins {
		if win == w {
			break
		}
		i++
	}
	if i == -1 {
		return
	}

	copy(w.d.wins[i:], w.d.wins[i+1:])
	w.d.wins[len(w.d.wins)-1] = nil // Allow the garbage collector to reclaim the memory.
	w.d.wins = w.d.wins[:len(w.d.wins)-1]

	if w.d.current > 0 {
		w.d.current--
	}
}

func (w *noosWindow) ShowAndRun() {
	w.Show()
	w.d.Run()
}

func (w *noosWindow) Content() fyne.CanvasObject {
	return w.c.Content()
}

func (w *noosWindow) SetContent(object fyne.CanvasObject) {
	w.c.SetContent(object)
}

func (w *noosWindow) Canvas() fyne.Canvas {
	return w.c
}

func (w *noosWindow) Clipboard() fyne.Clipboard {
	// TODO implement me
	return nil
}

func newWindow(d *noosDriver) fyne.Window {
	return &noosWindow{c: newTouchCanvas(), d: d}
}
