package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*Menu)(nil)
var _ fyne.Tappable = (*Menu)(nil)

// Menu is a widget for displaying a fyne.Menu.
type Menu struct {
	BaseWidget
	alignment     fyne.TextAlign
	Items         []fyne.CanvasObject
	OnDismiss     func()
	activeItem    *menuItem
	customSized   bool
	containsCheck bool
}

// NewMenu creates a new Menu.
func NewMenu(menu *fyne.Menu) *Menu {
	m := &Menu{}
	m.ExtendBaseWidget(m)
	m.setMenu(menu)
	return m
}

// ActivateLastSubmenu finds the last active menu item traversing through the open submenus
// and activates its submenu if any.
// It returns `true` if there was a submenu and it was activated and `false` elsewhere.
// Activating a submenu does show it and activate its first item.
func (m *Menu) ActivateLastSubmenu() bool {
	if m.activeItem == nil {
		return false
	}
	if !m.activeItem.activateLastSubmenu() {
		return false
	}
	m.Refresh()
	return true
}

// ActivateNext activates the menu item following the currently active menu item.
// If there is no menu item active, it activates the first menu item.
// If there is no menu item after the current active one, it does nothing.
// If a submenu is open, it delegates the activation to this submenu.
func (m *Menu) ActivateNext() {
	if m.activeItem != nil && m.activeItem.isSubmenuOpen() {
		m.activeItem.Child().ActivateNext()
		return
	}

	found := m.activeItem == nil
	for _, item := range m.Items {
		if mItem, ok := item.(*menuItem); ok {
			if found {
				m.activateItem(mItem)
				return
			}
			if mItem == m.activeItem {
				found = true
			}
		}
	}
}

// ActivatePrevious activates the menu item preceding the currently active menu item.
// If there is no menu item active, it activates the last menu item.
// If there is no menu item before the current active one, it does nothing.
// If a submenu is open, it delegates the activation to this submenu.
func (m *Menu) ActivatePrevious() {
	if m.activeItem != nil && m.activeItem.isSubmenuOpen() {
		m.activeItem.Child().ActivatePrevious()
		return
	}

	found := m.activeItem == nil
	for i := len(m.Items) - 1; i >= 0; i-- {
		item := m.Items[i]
		if mItem, ok := item.(*menuItem); ok {
			if found {
				m.activateItem(mItem)
				return
			}
			if mItem == m.activeItem {
				found = true
			}
		}
	}
}

// CreateRenderer returns a new renderer for the menu.
//
// Implements: fyne.Widget
func (m *Menu) CreateRenderer() fyne.WidgetRenderer {
	m.ExtendBaseWidget(m)
	box := newMenuBox(m.Items)
	scroll := widget.NewVScroll(box)
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

// DeactivateChild deactivates the active menu item and hides its submenu if any.
func (m *Menu) DeactivateChild() {
	if m.activeItem != nil {
		defer m.activeItem.Refresh()
		if c := m.activeItem.Child(); c != nil {
			c.Hide()
		}
		m.activeItem = nil
	}
}

// DeactivateLastSubmenu finds the last open submenu traversing through the open submenus,
// deactivates its active item and hides it.
// This also deactivates any submenus of the deactivated submenu.
// It returns `true` if there was a submenu open and closed and `false` elsewhere.
func (m *Menu) DeactivateLastSubmenu() bool {
	if m.activeItem == nil {
		return false
	}
	return m.activeItem.deactivateLastSubmenu()
}

// MinSize returns the minimal size of the menu.
//
// Implements: fyne.Widget
func (m *Menu) MinSize() fyne.Size {
	m.ExtendBaseWidget(m)
	return m.BaseWidget.MinSize()
}

// Refresh updates the menu to reflect changes in the data.
//
// Implements: fyne.Widget
func (m *Menu) Refresh() {
	for _, item := range m.Items {
		item.Refresh()
	}
	m.BaseWidget.Refresh()
}

func (m *Menu) getContainsCheck() bool {
	for _, item := range m.Items {
		if mi, ok := item.(*menuItem); ok && mi.Item.Checked {
			return true
		}
	}
	return false
}

// Tapped catches taps on separators and the menu background. It doesn't perform any action.
//
// Implements: fyne.Tappable
func (m *Menu) Tapped(*fyne.PointEvent) {
	// Hit a separator or padding -> do nothing.
}

// TriggerLast finds the last active menu item traversing through the open submenus and triggers it.
func (m *Menu) TriggerLast() {
	if m.activeItem == nil {
		m.Dismiss()
		return
	}
	m.activeItem.triggerLast()
}

// Dismiss dismisses the menu by dismissing and hiding the active child and performing OnDismiss.
func (m *Menu) Dismiss() {
	if m.activeItem != nil {
		if m.activeItem.Child() != nil {
			defer m.activeItem.Child().Dismiss()
		}
		m.DeactivateChild()
	}
	if m.OnDismiss != nil {
		m.OnDismiss()
	}
}

func (m *Menu) activateItem(item *menuItem) {
	if item.Child() != nil {
		item.Child().DeactivateChild()
	}
	if m.activeItem == item {
		return
	}

	m.DeactivateChild()
	m.activeItem = item
	m.activeItem.Refresh()
	if m.activeItem.child != nil {
		m.Refresh()
	}
}

func (m *Menu) setMenu(menu *fyne.Menu) {
	m.Items = make([]fyne.CanvasObject, len(menu.Items))
	for i, item := range menu.Items {
		if item.IsSeparator {
			m.Items[i] = NewSeparator()
		} else {
			m.Items[i] = newMenuItem(item, m)
		}
	}
	m.containsCheck = m.getContainsCheck()
}

type menuRenderer struct {
	*widget.ShadowingRenderer
	box    *menuBox
	m      *Menu
	scroll *widget.Scroll
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
	if c := fyne.CurrentApp().Driver().CanvasForObject(r.m.super()); c != nil {
		ap := fyne.CurrentApp().Driver().AbsolutePositionForObject(r.m.super())
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
	r.ShadowingRenderer.RefreshShadow()

	for _, i := range r.m.Items {
		if txt, ok := i.(*menuItem); ok {
			txt.alignment = r.m.alignment
			txt.Refresh()
		}
	}

	canvas.Refresh(r.m)
}

func (r *menuRenderer) layoutActiveChild() {
	item := r.m.activeItem
	if item == nil || item.Child() == nil {
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
	background := canvas.NewRectangle(theme.MenuBackgroundColor())
	cont := &fyne.Container{Layout: layout.NewVBoxLayout(), Objects: b.items}
	return &menuBoxRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{background, cont}),
		b:            b,
		background:   background,
		cont:         cont,
	}
}

type menuBoxRenderer struct {
	widget.BaseRenderer
	b          *menuBox
	background *canvas.Rectangle
	cont       *fyne.Container
}

var _ fyne.WidgetRenderer = (*menuBoxRenderer)(nil)

func (r *menuBoxRenderer) Layout(size fyne.Size) {
	s := fyne.NewSize(size.Width, size.Height)
	r.background.Resize(s)
	r.cont.Resize(s)
}

func (r *menuBoxRenderer) MinSize() fyne.Size {
	return r.cont.MinSize()
}

func (r *menuBoxRenderer) Refresh() {
	r.background.FillColor = theme.MenuBackgroundColor()
	r.background.Refresh()
	canvas.Refresh(r.b)
}
