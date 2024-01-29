package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*AppTabs)(nil)

// AppTabs container is used to split your application into various different areas identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 1.4
type AppTabs struct {
	widget.BaseWidget

	Items []*TabItem

	// Deprecated: Use `OnSelected func(*TabItem)` instead.
	OnChanged    func(*TabItem)
	OnSelected   func(*TabItem)
	OnUnselected func(*TabItem)

	current         int
	location        TabLocation
	isTransitioning bool

	popUpMenu *widget.PopUpMenu
}

// NewAppTabs creates a new tab container that allows the user to choose between different areas of an app.
//
// Since: 1.4
func NewAppTabs(items ...*TabItem) *AppTabs {
	tabs := &AppTabs{}
	tabs.BaseWidget.ExtendBaseWidget(tabs)
	tabs.SetItems(items)
	return tabs
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
//
// Implements: fyne.Widget
func (t *AppTabs) CreateRenderer() fyne.WidgetRenderer {
	t.BaseWidget.ExtendBaseWidget(t)
	r := &appTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			bar:       &fyne.Container{},
			divider:   canvas.NewRectangle(theme.ShadowColor()),
			indicator: canvas.NewRectangle(theme.PrimaryColor()),
		},
		appTabs: t,
	}
	r.action = r.buildOverflowTabsButton()
	r.tabs = t

	// Initially setup the tab bar to only show one tab, all others will be in overflow.
	// When the widget is laid out, and we know the size, the tab bar will be updated to show as many as can fit.
	r.updateTabs(1)
	r.updateIndicator(false)
	r.applyTheme(t)
	return r
}

// Append adds a new TabItem to the end of the tab bar.
func (t *AppTabs) Append(item *TabItem) {
	t.SetItems(append(t.Items, item))
}

// CurrentTab returns the currently selected TabItem.
//
// Deprecated: Use `AppTabs.Selected() *TabItem` instead.
func (t *AppTabs) CurrentTab() *TabItem {
	if t.current < 0 || t.current >= len(t.Items) {
		return nil
	}
	return t.Items[t.current]
}

// CurrentTabIndex returns the index of the currently selected TabItem.
//
// Deprecated: Use `AppTabs.SelectedIndex() int` instead.
func (t *AppTabs) CurrentTabIndex() int {
	return t.current
}

// DisableIndex disables the TabItem at the specified index.
//
// Since: 2.3
func (t *AppTabs) DisableIndex(i int) {
	disableIndex(t, i)
}

// DisableItem disables the specified TabItem.
//
// Since: 2.3
func (t *AppTabs) DisableItem(item *TabItem) {
	disableItem(t, item)
}

// EnableIndex enables the TabItem at the specified index.
//
// Since: 2.3
func (t *AppTabs) EnableIndex(i int) {
	enableIndex(t, i)
}

// EnableItem enables the specified TabItem.
//
// Since: 2.3
func (t *AppTabs) EnableItem(item *TabItem) {
	enableItem(t, item)
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
//
// Deprecated: Support for extending containers is being removed
func (t *AppTabs) ExtendBaseWidget(wid fyne.Widget) {
	t.BaseWidget.ExtendBaseWidget(wid)
}

// Hide hides the widget.
//
// Implements: fyne.CanvasObject
func (t *AppTabs) Hide() {
	if t.popUpMenu != nil {
		t.popUpMenu.Hide()
		t.popUpMenu = nil
	}
	t.BaseWidget.Hide()
}

// MinSize returns the size that this widget should not shrink below
//
// Implements: fyne.CanvasObject
func (t *AppTabs) MinSize() fyne.Size {
	t.BaseWidget.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Remove tab by value.
func (t *AppTabs) Remove(item *TabItem) {
	removeItem(t, item)
	t.Refresh()
}

// RemoveIndex removes tab by index.
func (t *AppTabs) RemoveIndex(index int) {
	removeIndex(t, index)
	t.Refresh()
}

// Select sets the specified TabItem to be selected and its content visible.
func (t *AppTabs) Select(item *TabItem) {
	selectItem(t, item)
	t.Refresh()
}

// SelectIndex sets the TabItem at the specific index to be selected and its content visible.
func (t *AppTabs) SelectIndex(index int) {
	selectIndex(t, index)
	t.Refresh()
}

// SelectTab sets the specified TabItem to be selected and its content visible.
//
// Deprecated: Use `AppTabs.Select(*TabItem)` instead.
func (t *AppTabs) SelectTab(item *TabItem) {
	for i, child := range t.Items {
		if child == item {
			t.SelectTabIndex(i)
			return
		}
	}
}

// SelectTabIndex sets the TabItem at the specific index to be selected and its content visible.
//
// Deprecated: Use `AppTabs.SelectIndex(int)` instead.
func (t *AppTabs) SelectTabIndex(index int) {
	if index < 0 || index >= len(t.Items) || t.current == index {
		return
	}
	t.current = index
	t.Refresh()

	if t.OnChanged != nil {
		t.OnChanged(t.Items[t.current])
	}
}

// Selected returns the currently selected TabItem.
func (t *AppTabs) Selected() *TabItem {
	return selected(t)
}

// SelectedIndex returns the index of the currently selected TabItem.
func (t *AppTabs) SelectedIndex() int {
	return t.current
}

// SetItems sets the containers items and refreshes.
func (t *AppTabs) SetItems(items []*TabItem) {
	setItems(t, items)
	t.Refresh()
}

// SetTabLocation sets the location of the tab bar
func (t *AppTabs) SetTabLocation(l TabLocation) {
	t.location = tabsAdjustedLocation(l)
	t.Refresh()
}

// Show this widget, if it was previously hidden
//
// Implements: fyne.CanvasObject
func (t *AppTabs) Show() {
	t.BaseWidget.Show()
	t.SelectIndex(t.current)
}

func (t *AppTabs) onUnselected() func(*TabItem) {
	return t.OnUnselected
}

func (t *AppTabs) onSelected() func(*TabItem) {
	return func(tab *TabItem) {
		if f := t.OnChanged; f != nil {
			f(tab)
		}
		if f := t.OnSelected; f != nil {
			f(tab)
		}
	}
}

func (t *AppTabs) items() []*TabItem {
	return t.Items
}

func (t *AppTabs) selected() int {
	return t.current
}

func (t *AppTabs) setItems(items []*TabItem) {
	t.Items = items
}

func (t *AppTabs) setSelected(selected int) {
	t.current = selected
}

func (t *AppTabs) setTransitioning(transitioning bool) {
	t.isTransitioning = transitioning
}

func (t *AppTabs) tabLocation() TabLocation {
	return t.location
}

func (t *AppTabs) transitioning() bool {
	return t.isTransitioning
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*appTabsRenderer)(nil)

type appTabsRenderer struct {
	baseTabsRenderer
	appTabs *AppTabs
}

func (r *appTabsRenderer) Layout(size fyne.Size) {
	// Try render as many tabs as will fit, others will appear in the overflow
	if len(r.appTabs.Items) == 0 {
		r.updateTabs(0)
	} else {
		for i := len(r.appTabs.Items); i > 0; i-- {
			r.updateTabs(i)
			barMin := r.bar.MinSize()
			if r.appTabs.location == TabLocationLeading || r.appTabs.location == TabLocationTrailing {
				if barMin.Height <= size.Height {
					// Tab bar is short enough to fit
					break
				}
			} else {
				if barMin.Width <= size.Width {
					// Tab bar is thin enough to fit
					break
				}
			}
		}
	}

	r.layout(r.appTabs, size)
	r.updateIndicator(r.appTabs.transitioning())
	if r.appTabs.transitioning() {
		r.appTabs.setTransitioning(false)
	}
}

func (r *appTabsRenderer) MinSize() fyne.Size {
	return r.minSize(r.appTabs)
}

func (r *appTabsRenderer) Objects() []fyne.CanvasObject {
	return r.objects(r.appTabs)
}

func (r *appTabsRenderer) Refresh() {
	r.Layout(r.appTabs.Size())

	r.refresh(r.appTabs)

	canvas.Refresh(r.appTabs)
}

func (r *appTabsRenderer) buildOverflowTabsButton() (overflow *widget.Button) {
	overflow = &widget.Button{Icon: moreIcon(r.appTabs), Importance: widget.LowImportance, OnTapped: func() {
		// Show pop up containing all tabs which did not fit in the tab bar

		itemLen, objLen := len(r.appTabs.Items), len(r.bar.Objects[0].(*fyne.Container).Objects)
		items := make([]*fyne.MenuItem, 0, itemLen-objLen)
		for i := objLen; i < itemLen; i++ {
			index := i // capture
			// FIXME MenuItem doesn't support icons (#1752)
			// FIXME MenuItem can't show if it is the currently selected tab (#1753)
			items = append(items, fyne.NewMenuItem(r.appTabs.Items[i].Text, func() {
				r.appTabs.SelectIndex(index)
				if r.appTabs.popUpMenu != nil {
					r.appTabs.popUpMenu.Hide()
					r.appTabs.popUpMenu = nil
				}
			}))
		}

		r.appTabs.popUpMenu = buildPopUpMenu(r.appTabs, overflow, items)
	}}

	return overflow
}

func (r *appTabsRenderer) buildTabButtons(count int) *fyne.Container {
	buttons := &fyne.Container{}

	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() {
		cells := count
		if cells == 0 {
			cells = 1
		}
		if r.appTabs.location == TabLocationTop || r.appTabs.location == TabLocationBottom {
			buttons.Layout = layout.NewGridLayoutWithColumns(cells)
		} else {
			buttons.Layout = layout.NewGridLayoutWithRows(cells)
		}
		iconPos = buttonIconTop
	} else if r.appTabs.location == TabLocationLeading || r.appTabs.location == TabLocationTrailing {
		buttons.Layout = layout.NewVBoxLayout()
		iconPos = buttonIconTop
	} else {
		buttons.Layout = layout.NewHBoxLayout()
		iconPos = buttonIconInline
	}

	for i := 0; i < count; i++ {
		item := r.appTabs.Items[i]
		if item.button == nil {
			item.button = &tabButton{
				onTapped: func() { r.appTabs.Select(item) },
			}
		}
		button := item.button
		button.icon = item.Icon
		button.iconPosition = iconPos
		if i == r.appTabs.current {
			button.importance = widget.HighImportance
		} else {
			button.importance = widget.MediumImportance
		}
		button.text = item.Text
		button.textAlignment = fyne.TextAlignCenter
		button.Refresh()
		buttons.Objects = append(buttons.Objects, button)
	}
	return buttons
}

func (r *appTabsRenderer) updateIndicator(animate bool) {
	if r.appTabs.current < 0 {
		r.indicator.Hide()
		return
	}

	var selectedPos fyne.Position
	var selectedSize fyne.Size

	buttons := r.bar.Objects[0].(*fyne.Container).Objects
	if r.appTabs.current >= len(buttons) {
		if a := r.action; a != nil {
			selectedPos = a.Position()
			selectedSize = a.Size()
		}
	} else {
		selected := buttons[r.appTabs.current]
		selectedPos = selected.Position()
		selectedSize = selected.Size()
	}

	var indicatorPos fyne.Position
	var indicatorSize fyne.Size

	switch r.appTabs.location {
	case TabLocationTop:
		indicatorPos = fyne.NewPos(selectedPos.X, r.bar.MinSize().Height)
		indicatorSize = fyne.NewSize(selectedSize.Width, theme.Padding())
	case TabLocationLeading:
		indicatorPos = fyne.NewPos(r.bar.MinSize().Width, selectedPos.Y)
		indicatorSize = fyne.NewSize(theme.Padding(), selectedSize.Height)
	case TabLocationBottom:
		indicatorPos = fyne.NewPos(selectedPos.X, r.bar.Position().Y-theme.Padding())
		indicatorSize = fyne.NewSize(selectedSize.Width, theme.Padding())
	case TabLocationTrailing:
		indicatorPos = fyne.NewPos(r.bar.Position().X-theme.Padding(), selectedPos.Y)
		indicatorSize = fyne.NewSize(theme.Padding(), selectedSize.Height)
	}

	r.moveIndicator(indicatorPos, indicatorSize, animate)
}

func (r *appTabsRenderer) updateTabs(max int) {
	tabCount := len(r.appTabs.Items)

	// Set overflow action
	if tabCount <= max {
		r.action.Hide()
		r.bar.Layout = layout.NewStackLayout()
	} else {
		tabCount = max
		r.action.Show()

		// Set layout of tab bar containing tab buttons and overflow action
		if r.appTabs.location == TabLocationLeading || r.appTabs.location == TabLocationTrailing {
			r.bar.Layout = layout.NewBorderLayout(nil, r.action, nil, nil)
		} else {
			r.bar.Layout = layout.NewBorderLayout(nil, nil, nil, r.action)
		}
	}

	buttons := r.buildTabButtons(tabCount)

	r.bar.Objects = []fyne.CanvasObject{buttons}
	if a := r.action; a != nil {
		r.bar.Objects = append(r.bar.Objects, a)
	}

	r.bar.Refresh()
}
