package mobile

import (
	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type menuButton struct {
	widget.BaseWidget
	win  *window
	menu *fyne.MainMenu
}

func (w *window) newMenuButton(menu *fyne.MainMenu) *menuButton {
	b := &menuButton{win: w, menu: menu}
	b.ExtendBaseWidget(b)
	return b
}

func (m *menuButton) CreateRenderer() fyne.WidgetRenderer {
	return &menuButtonRenderer{btn: widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		m.win.canvas.showMenu(m.menu)
	}), bg: fynecanvas.NewRectangle(theme.Color(theme.ColorNameBackground))}
}

type menuButtonRenderer struct {
	btn *widget.Button
	bg  *fynecanvas.Rectangle
}

func (m *menuButtonRenderer) Destroy() {
}

func (m *menuButtonRenderer) Layout(size fyne.Size) {
	m.bg.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))
	m.bg.Resize(size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding())))
	m.btn.Resize(size)
}

func (m *menuButtonRenderer) MinSize() fyne.Size {
	return m.btn.MinSize()
}

func (m *menuButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{m.bg, m.btn}
}

func (m *menuButtonRenderer) Refresh() {
	m.bg.FillColor = theme.Color(theme.ColorNameBackground)
}
