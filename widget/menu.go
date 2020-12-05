package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*Menu)(nil)
var _ fyne.Tappable = (*Menu)(nil)

// Menu is a widget for displaying a fyne.Menu.
type Menu struct {
	widget.Base
	Items       []fyne.CanvasObject
	OnDismiss   func()
	activeItem  *menuItem
	customSized bool
}

// NewMenu creates a new Menu.
func NewMenu(menu *fyne.Menu) *Menu {
	items := make([]fyne.CanvasObject, len(menu.Items))
	m := &Menu{Items: items}
	for i, item := range menu.Items {
		if item.IsSeparator {
			items[i] = NewSeparator()
		} else {
			items[i] = newMenuItem(item, m, m.activateChild)
		}
	}
	return m
}

// CreateRenderer returns a new renderer for the menu.
//
// Implements: fyne.Widget
func (m *Menu) CreateRenderer() fyne.WidgetRenderer {
	box := newMenuBox(m.Items)
	scroll := NewVScrollContainer(box)
	scroll.SetMinSize(box.MinSize())
	objects := []fyne.CanvasObject{scroll}
	for _, i := range m.Items {
		if item, ok := i.(*menuItem); ok && item.Child() != nil {
			objects = append(objects, item.Child())
		}
	}

	return &menuRenderer{
		widget.NewShadowingRenderer(objects, widget.MenuLevel),
		box,
		m,
		scroll,
	}
}

// DeactivateChild deactivates the active child menu.
func (m *Menu) DeactivateChild() {
	if m.activeItem != nil {
		m.activeItem.Child().Hide()
		m.activeItem = nil
	}
}

// Hide hides the menu.
//
// Implements: fyne.Widget
func (m *Menu) Hide() {
	widget.HideWidget(&m.Base, m)
}

// MinSize returns the minimal size of the menu.
//
// Implements: fyne.Widget
func (m *Menu) MinSize() fyne.Size {
	return widget.MinSizeOf(m)
}

// Move sets the position of the widget relative to its parent.
//
// Implements: fyne.Widget
func (m *Menu) Move(pos fyne.Position) {
	widget.MoveWidget(&m.Base, m, pos)
}

// Refresh triggers a redraw of the menu.
//
// Implements: fyne.Widget
func (m *Menu) Refresh() {
	widget.RefreshWidget(m)
}

// Resize has no effect because menus are always displayed with their minimal size.
//
// Implements: fyne.Widget
func (m *Menu) Resize(size fyne.Size) {
	widget.ResizeWidget(&m.Base, m, size)
}

// Show makes the menu visible.
//
// Implements: fyne.Widget
func (m *Menu) Show() {
	widget.ShowWidget(&m.Base, m)
}

// Tapped catches taps on separators and the menu background. It doesn’t perform any action.
//
// Implements: fyne.Tappable
func (m *Menu) Tapped(*fyne.PointEvent) {
	// Hit a separator or padding -> do nothing.
}

// Dismiss dismisses the menu by dismissing and hiding the active child and performing OnDismiss.
func (m *Menu) Dismiss() {
	if m.activeItem != nil {
		defer m.activeItem.Child().Dismiss()
		m.DeactivateChild()
	}
	if m.OnDismiss != nil {
		m.OnDismiss()
	}
}

func (m *Menu) activateChild(item *menuItem) {
	if item.Child() != nil {
		item.Child().DeactivateChild()
	}
	if m.activeItem == item {
		return
	}

	m.DeactivateChild()
	if item.Child() == nil {
		return
	}

	m.activeItem = item
	item.Child().Show()
	m.Refresh()
}

type menuRenderer struct {
	*widget.ShadowingRenderer
	box    *menuBox
	m      *Menu
	scroll *ScrollContainer
}

func (r *menuRenderer) Layout(s fyne.Size) {
	minSize := r.MinSize()
	var boxSize fyne.Size
	if r.m.customSized {
		boxSize = minSize.Max(s)
	} else {
		boxSize = minSize
	}
	scrollSize := boxSize
	if c := fyne.CurrentApp().Driver().CanvasForObject(r.m); c != nil {
		ap := fyne.CurrentApp().Driver().AbsolutePositionForObject(r.m)
		pos, size := c.InteractiveArea()
		bottomPad := c.Size().Height - pos.Y - size.Height
		if ah := c.Size().Height - bottomPad - ap.Y; ah < boxSize.Height {
			scrollSize = fyne.NewSize(boxSize.Width, ah)
		}
	}
	if scrollSize != r.m.Size() {
		r.m.Resize(scrollSize)
		return
	}

	r.LayoutShadow(scrollSize, fyne.NewPos(0, 0))
	r.scroll.Resize(scrollSize)
	r.box.Resize(boxSize)
	r.layoutActiveChild()
}

func (r *menuRenderer) MinSize() fyne.Size {
	return r.box.MinSize()
}

func (r *menuRenderer) Refresh() {
	r.layoutActiveChild()
	canvas.Refresh(r.m)
}

func (r *menuRenderer) layoutActiveChild() {
	item := r.m.activeItem
	if item == nil {
		return
	}

	if item.Child().Size().IsZero() {
		item.Child().Resize(item.Child().MinSize())
	}

	itemSize := item.Size()
	cp := fyne.NewPos(itemSize.Width, item.Position().Y)
	d := fyne.CurrentApp().Driver()
	c := d.CanvasForObject(item)
	if c != nil {
		absPos := d.AbsolutePositionForObject(item)
		childSize := item.Child().Size()
		if absPos.X+itemSize.Width+childSize.Width > c.Size().Width {
			if absPos.X-childSize.Width >= 0 {
				cp.X = -childSize.Width
			} else {
				cp.X = c.Size().Width - absPos.X - childSize.Width
			}
		}
		requiredHeight := childSize.Height - theme.Padding()
		availableHeight := c.Size().Height - absPos.Y
		missingHeight := requiredHeight - availableHeight
		if missingHeight > 0 {
			cp.Y -= missingHeight
		}
	}
	item.Child().Move(cp)
}

type menuBox struct {
	BaseWidget
	items []fyne.CanvasObject
}

var _ fyne.Widget = (*menuBox)(nil)

func newMenuBox(items []fyne.CanvasObject) *menuBox {
	b := &menuBox{items: items}
	b.ExtendBaseWidget(b)
	return b
}

func (b *menuBox) CreateRenderer() fyne.WidgetRenderer {
	cont := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), b.items...)
	return &menuBoxRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{cont}),
		b:            b,
		cont:         cont,
	}
}

type menuBoxRenderer struct {
	widget.BaseRenderer
	b    *menuBox
	cont *fyne.Container
}

var _ fyne.WidgetRenderer = (*menuBoxRenderer)(nil)

func (r *menuBoxRenderer) Layout(size fyne.Size) {
	r.cont.Resize(fyne.NewSize(size.Width, size.Height+2*theme.Padding()))
	r.cont.Move(fyne.NewPos(0, theme.Padding()))
}

func (r *menuBoxRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(fyne.NewSize(0, 2*theme.Padding()))
}

func (r *menuBoxRenderer) Refresh() {
	canvas.Refresh(r.b)
}
