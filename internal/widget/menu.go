package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*Menu)(nil)
var _ fyne.Tappable = (*Menu)(nil)

// Menu is a widget for displaying a fyne.Menu.
type Menu struct {
	base
	menuBase
	Items []fyne.CanvasObject
}

// NewMenu creates a new Menu.
func NewMenu(menu *fyne.Menu) *Menu {
	items := make([]fyne.CanvasObject, len(menu.Items))
	m := &Menu{Items: items}
	for i, item := range menu.Items {
		if item.IsSeparator {
			items[i] = NewMenuItemSeparator()
		} else {
			items[i] = NewMenuItem(item, m)
		}
	}
	return m
}

// CreateRenderer satisfies the fyne.Widget interface.
func (m *Menu) CreateRenderer() fyne.WidgetRenderer {
	cont := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), m.Items...)
	return &menuRenderer{
		NewShadowingRenderer([]fyne.CanvasObject{cont}, MenuLevel),
		cont,
		m,
	}
}

// Hide satisfies the fyne.Widget interface.
func (m *Menu) Hide() {
	m.hide(m)
}

// MinSize satisfies the fyne.Widget interface.
func (m *Menu) MinSize() fyne.Size {
	return m.minSize(m)
}

// Refresh satisfies the fyne.Widget interface.
func (m *Menu) Refresh() {
	m.refresh(m)
}

// Resize satisfies the fyne.Widget interface.
func (m *Menu) Resize(size fyne.Size) {
	m.resize(size, m)
}

// Show satisfies the fyne.Widget interface.
func (m *Menu) Show() {
	m.show(m)
}

// Tapped satisfies the fyne.Tappable interface.
func (m *Menu) Tapped(*fyne.PointEvent) {
	// Hit a separator or padding -> do nothing.
}

type menuRenderer struct {
	*ShadowingRenderer
	cont *fyne.Container
	m    *Menu
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *menuRenderer) Layout(size fyne.Size) {
	r.LayoutShadow(size, fyne.NewPos(0, 0))
	padding := r.padding()
	r.cont.Resize(size.Subtract(padding))
	r.cont.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(r.padding())
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *menuRenderer) Refresh() {
	canvas.Refresh(r.m)
}

func (r *menuRenderer) padding() fyne.Size {
	return fyne.NewSize(0, theme.Padding()*2)
}
