package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type window struct {
	title              string
	visible            bool
	onClosed           func()
	onCloseIntercepted func()
	isChild            bool

	canvas *canvas
	icon   fyne.Resource
	menu   *fyne.MainMenu
	handle uintptr // the window handle - currently just Android
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

func (w *window) Resize(size fyne.Size) {
	w.Canvas().(*canvas).Resize(size)
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

func (w *window) SetMaster() {
	// no-op on mobile
}

func (w *window) MainMenu() *fyne.MainMenu {
	return w.menu
}

func (w *window) SetMainMenu(menu *fyne.MainMenu) {
	w.menu = menu
}

func (w *window) SetOnClosed(callback func()) {
	w.onClosed = callback
}

func (w *window) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *window) SetOnDropped(dropped func(fyne.Position, []fyne.URI)) {
	// not implemented yet
}

func (w *window) Show() {
	menu := fyne.CurrentApp().Driver().(*driver).findMenu(w)
	menuButton := w.newMenuButton(menu)
	if menu == nil {
		menuButton.Hide()
	}

	if w.isChild {
		exit := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			w.tryClose()
		})
		title := widget.NewLabel(w.title)
		title.Alignment = fyne.TextAlignCenter
		w.canvas.setWindowHead(container.NewHBox(menuButton,
			layout.NewSpacer(), title, layout.NewSpacer(), exit))
		w.canvas.Resize(w.canvas.size)
	} else {
		w.canvas.setWindowHead(container.NewHBox(menuButton))
	}
	w.visible = true

	if w.Content() != nil {
		w.Content().Refresh()
		w.Content().Show()
	}
}

func (w *window) Hide() {
	w.visible = false

	if w.Content() != nil {
		w.Content().Hide()
	}
}

func (w *window) tryClose() {
	if w.onCloseIntercepted != nil {
		w.onCloseIntercepted()
		return
	}

	w.Close()
}

func (w *window) Close() {
	d := fyne.CurrentApp().Driver().(*driver)
	pos := -1
	for i, win := range d.windows {
		if win == w {
			pos = i
		}
	}
	if pos != -1 {
		d.windows = append(d.windows[:pos], d.windows[pos+1:]...)
	}

	cache.RangeTexturesFor(w.canvas, w.canvas.Painter().Free)

	w.canvas.WalkTrees(nil, func(node *common.RenderCacheNode, _ fyne.Position) {
		if wid, ok := node.Obj().(fyne.Widget); ok {
			cache.DestroyRenderer(wid)
		}
	})
	cache.CleanCanvas(w.canvas)

	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.CurrentApp().Run()
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
	return NewClipboard()
}

func (w *window) RunWithContext(f func()) {
	//	ctx, _ = e.DrawContext.(gl.Context)

	f()
}

func (w *window) RescaleContext() {
	// TODO
}

func (w *window) Context() any {
	return fyne.CurrentApp().Driver().(*driver).glctx
}
