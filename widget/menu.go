package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	m := NewMenuWidget(menu, false)
	pop := newPopUp(m, c)
	pop.WithoutPadding = true
	focused := c.Focused()
	m.DismissAction = func() {
		if c.Focused() == nil {
			c.Focus(focused)
		}
		pop.Hide()
	}
	pop.ShowAtPosition(pos)
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}

func newMenuItemWidget(item *fyne.MenuItem, parent *MenuWidget) *MenuItemWidget {
	ret := &MenuItemWidget{Item: item, parent: parent}
	ret.ExtendBaseWidget(ret)
	return ret
}

func newSeparator() fyne.CanvasObject {
	s := canvas.NewRectangle(theme.DisabledTextColor())
	s.SetMinSize(fyne.NewSize(1, 2))
	return s
}

// NewMenuWidget creates a menu widget populated with menu items from the passed menu structure.
func NewMenuWidget(menu *fyne.Menu, showShadow bool) *MenuWidget {
	items := make([]fyne.CanvasObject, len(menu.Items))
	w := &MenuWidget{Items: items, ShowShadow: showShadow}
	w.ExtendBaseWidget(w)
	for i, item := range menu.Items {
		if item.IsSeparator {
			items[i] = newSeparator()
		} else {
			items[i] = newMenuItemWidget(item, w)
		}
	}
	return w
}

// NewMenuBarWidget creates a menu bar widget populated with menu items from the passed main menu structure.
func NewMenuBarWidget(mainMenu *fyne.MainMenu) *MenuBarWidget {
	items := make([]fyne.CanvasObject, len(mainMenu.Items))
	w := &MenuBarWidget{MenuWidget: MenuWidget{Items: items, ShowShadow: true}}
	w.ExtendBaseWidget(w)
	for i, menu := range mainMenu.Items {
		item := fyne.NewMenuItem(menu.Label, nil)
		item.ChildMenu = menu
		iw := &MenuBarItemWidget{
			MenuItemWidget{Item: item, parent: &w.MenuWidget},
			w,
		}
		iw.ExtendBaseWidget(iw)
		items[i] = iw
	}
	return w
}

type MenuWidget struct {
	BaseWidget
	DismissAction func()
	Items         []fyne.CanvasObject
	ShowShadow    bool

	activeChild *MenuWidget
}

func (w *MenuWidget) CreateRenderer() fyne.WidgetRenderer {
	shadowLevel := baseLevel
	if w.ShowShadow {
		shadowLevel = menuLevel
	}

	box := NewVBox(w.Items...)
	box.background = color.Transparent
	return &menuWidgetRenderer{
		newShadowingRenderer([]fyne.CanvasObject{box}, shadowLevel),
		theme.BackgroundColor,
		box,
		func() fyne.Size { return fyne.NewSize(0, theme.Padding()*2) },
		w,
	}
}

func (w *MenuWidget) dismiss() {
	if w.activeChild != nil {
		defer w.activeChild.dismiss()
		w.activeChild.Hide()
		w.activeChild = nil
	}
	if w.DismissAction != nil {
		w.DismissAction()
	}
}

type menuWidgetRenderer struct {
	*shadowingRenderer

	bgColor func() color.Color
	box     *Box
	padding func() fyne.Size
	w       *MenuWidget
}

func (r *menuWidgetRenderer) BackgroundColor() color.Color {
	return r.bgColor()
}

func (r *menuWidgetRenderer) Layout(size fyne.Size) {
	r.layoutShadow(size, fyne.NewPos(0, 0))
	padding := r.padding()
	r.box.Resize(size.Subtract(padding))
	r.box.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuWidgetRenderer) MinSize() fyne.Size {
	return r.box.MinSize().Add(r.padding())
}

func (r *menuWidgetRenderer) Refresh() {
	canvas.Refresh(r.w)
}

type MenuBarWidget struct {
	MenuWidget
	ActivateAction func()

	active bool
}

func (w *MenuBarWidget) CreateRenderer() fyne.WidgetRenderer {
	shadowLevel := baseLevel
	if w.ShowShadow {
		shadowLevel = menuLevel
	}

	box := NewHBox(w.Items...)
	box.background = color.Transparent

	return &menuWidgetRenderer{
		newShadowingRenderer([]fyne.CanvasObject{box}, shadowLevel),
		theme.ButtonColor,
		box,
		func() fyne.Size { return fyne.NewSize(theme.Padding()*2, 0) },
		&w.MenuWidget,
	}
}

func (w *MenuBarWidget) Deactivate() {
	if !w.active {
		return
	}

	w.active = false
	w.dismiss()
}

func (w *MenuBarWidget) activate() {
	if w.active {
		return
	}

	w.active = true
	if w.ActivateAction != nil {
		w.ActivateAction()
	}
}

type MenuItemWidget struct {
	BaseWidget
	Item *fyne.MenuItem

	child   *MenuWidget
	hovered bool
	parent  *MenuWidget
}

func (t *MenuItemWidget) Tapped(*fyne.PointEvent) {
	if t.Item.Action == nil {
		return
	}

	t.Item.Action()
	t.parent.dismiss()
}

func (t *MenuItemWidget) activateChild() {
	if t.child != nil && t.child.activeChild != nil {
		t.child.activeChild.Hide()
		t.child.activeChild = nil
	}

	if t.parent.activeChild == t.child {
		return
	}

	if t.parent.activeChild != nil {
		t.parent.activeChild.Hide()
	}
	t.parent.activeChild = t.child
	if t.child != nil {
		t.child.Show()
	}
}

func (t *MenuItemWidget) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(t.Item.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	var icon *Icon
	if t.Item.ChildMenu != nil {
		icon = NewIcon(theme.MenuExpandIcon())
		objects = append(objects, icon)
		t.initChildWidget()
		objects = append(objects, t.child)
	}

	return &menuItemWidgetRenderer{
		baseRenderer{objects},
		t.child,
		icon,
		false,
		text,
		t,
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *MenuItemWidget) MouseIn(*desktop.MouseEvent) {
	t.hovered = true
	t.activateChild()
	t.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (t *MenuItemWidget) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (t *MenuItemWidget) MouseOut() {
	t.hovered = false
	t.Refresh()
}

func (t *MenuItemWidget) initChildWidget() {
	if t.child != nil {
		return
	}
	t.child = NewMenuWidget(t.Item.ChildMenu, t.parent.ShowShadow)
	t.child.Hide()
	t.child.DismissAction = func() { t.parent.dismiss() }
}

type menuItemWidgetRenderer struct {
	baseRenderer
	child          *MenuWidget
	icon           *Icon
	showChildBelow bool
	text           *canvas.Text
	w              *MenuItemWidget
}

func (r *menuItemWidgetRenderer) Layout(size fyne.Size) {
	padding := r.padding()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))

	if r.icon != nil {
		r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(fyne.NewPos(size.Width-theme.IconInlineSize(), (size.Height-theme.IconInlineSize())/2))
	}

	if r.child != nil {
		r.child.Resize(r.child.MinSize())
		if r.showChildBelow {
			r.child.Move(fyne.NewPos(0, r.w.Size().Height))
		} else {
			r.child.Move(fyne.NewPos(r.w.Size().Width, 0))
		}
	}
}

func (r *menuItemWidgetRenderer) MinSize() fyne.Size {
	s := r.text.MinSize().Add(r.padding())
	if r.icon != nil {
		s = s.Add(fyne.NewSize(theme.Padding(), 0))
		s = s.Add(fyne.NewSize(theme.IconInlineSize(), 0))
	}
	return s
}

func (r *menuItemWidgetRenderer) Refresh() {
	if r.text.TextSize != theme.TextSize() {
		defer r.Layout(r.w.Size())
	}
	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	canvas.Refresh(r.text)
}

func (r *menuItemWidgetRenderer) BackgroundColor() color.Color {
	if r.w.hovered {
		return theme.HoverColor()
	}

	return color.Transparent
}

func (r *menuItemWidgetRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}

type MenuBarItemWidget struct {
	MenuItemWidget
	parent *MenuBarWidget
}

func (t *MenuBarItemWidget) Tapped(*fyne.PointEvent) {
	if t.parent.active {
		t.parent.Deactivate()
	} else {
		t.parent.activate()
		t.activateChild()
	}
	t.Refresh()
}

func (t *MenuBarItemWidget) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(t.Item.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	if t.Item.ChildMenu != nil {
		t.initChildWidget()
		objects = append(objects, t.child)
	}

	return &menuItemWidgetRenderer{
		baseRenderer{objects},
		t.child,
		nil,
		true,
		text,
		&t.MenuItemWidget,
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *MenuBarItemWidget) MouseIn(e *desktop.MouseEvent) {
	if t.parent.active {
		t.MenuItemWidget.MouseIn(e)
	} else {
		t.hovered = true
		t.Refresh()
	}
}

func (t *MenuBarItemWidget) initChildWidget() {
	t.MenuItemWidget.initChildWidget()
	t.child.DismissAction = func() { t.parent.Deactivate() }
}
