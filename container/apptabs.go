package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// maxAppTabs is the maximum number of visible tabs in an AppTab.
// Any addition tabs are moved to the overflow menu.
const maxAppTabs = 7

// Declare conformity with Widget interface.
var _ fyne.Widget = (*AppTabs)(nil)

// AppTabs container is used to split your application into various different areas identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 1.4
type AppTabs struct {
	baseTabs
	// Deprecated: 'Items' has been moved into an embedded type.
	// Most devs using AppTabs will not be affected when this is removed as they can still assign to Items field and call the same public APIs.
	// The only devs affected would be ones using struct init (a := &AppTabs{Items: ...}) who'll need to change to (a := &AppTabs{}; a.Items = ...)
	Items []*TabItem
	// Deprecated: Use `OnSelected func(*TabItem)` instead.
	OnChanged func(tab *TabItem)
}

// NewAppTabs creates a new tab container that allows the user to choose between different areas of an app.
//
// Since: 1.4
func NewAppTabs(items ...*TabItem) *AppTabs {
	tabs := &AppTabs{
		baseTabs: baseTabs{
			BaseWidget: widget.BaseWidget{},
			current:    -1,
		},
	}
	tabs.SetItems(items)
	tabs.BaseWidget.ExtendBaseWidget(tabs)
	return tabs
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *AppTabs) CreateRenderer() fyne.WidgetRenderer {
	t.baseTabs.Items = t.Items
	t.BaseWidget.ExtendBaseWidget(t)
	r := &appTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			bar:         &fyne.Container{},
			divider:     canvas.NewRectangle(theme.ShadowColor()),
			indicator:   canvas.NewRectangle(theme.PrimaryColor()),
			buttonCache: make(map[*TabItem]*tabButton),
		},
		appTabs: t,
	}
	// Initially setup the tab bar to only show one tab, all others will be in overflow.
	// When the widget is laid out, and we know the size, the tab bar will be updated to show as many as can fit (capped to maxAppTab).
	r.updateTabs(1)
	r.updateIndicator()
	return r
}

// Append adds a new TabItem to the end of the tab bar.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) Append(item *TabItem) {
	t.baseTabs.Items = t.Items
	t.baseTabs.Append(item)
	t.Items = t.baseTabs.Items
}

// CurrentTab returns the currently selected TabItem.
//
// Deprecated: Use `AppTabs.Selection() *TabItem` instead.
func (t *AppTabs) CurrentTab() *TabItem {
	if t.current < 0 || t.current >= len(t.Items) {
		return nil
	}
	return t.Items[t.current]
}

// CurrentTabIndex returns the index of the currently selected TabItem.
//
// Deprecated: Use `AppTabs.SelectionIndex() int` instead.
func (t *AppTabs) CurrentTabIndex() int {
	return t.current
}

// Hide hides the widget.
//
// Implements: fyne.Widget
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) Hide() {
	t.baseTabs.Items = t.Items
	t.baseTabs.Hide()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
//
// Deprecated: Support for extending containers is being removed
func (t *AppTabs) ExtendBaseWidget(wid fyne.Widget) {
	t.BaseWidget.ExtendBaseWidget(wid)
}

// MinSize returns the size that this widget should not shrink below
func (t *AppTabs) MinSize() fyne.Size {
	t.baseTabs.Items = t.Items
	t.BaseWidget.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Remove tab by value.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) Remove(item *TabItem) {
	t.baseTabs.Items = t.Items
	t.baseTabs.Remove(item)
	t.Items = t.baseTabs.Items
}

// RemoveIndex removes tab by index.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) RemoveIndex(index int) {
	t.baseTabs.Items = t.Items
	t.baseTabs.RemoveIndex(index)
	t.Items = t.baseTabs.Items
}

// Select sets the specified TabItem to be selected and its content visible.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) Select(item *TabItem) {
	t.baseTabs.Items = t.Items
	t.baseTabs.Select(item)
}

// SelectIndex sets the TabItem at the specific index to be selected and its content visible.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) SelectIndex(index int) {
	t.baseTabs.Items = t.Items
	t.baseTabs.SelectIndex(index)
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

// Selection returns the currently selected TabItem.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) Selection() *TabItem {
	t.baseTabs.Items = t.Items
	return t.baseTabs.Selection()
}

// SelectionIndex returns the index of the currently selected TabItem.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) SelectionIndex() int {
	t.baseTabs.Items = t.Items
	return t.baseTabs.SelectionIndex()
}

// SetItems sets the container’s items and refreshes.
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) SetItems(items []*TabItem) {
	t.baseTabs.SetItems(items)
	t.Items = t.baseTabs.Items
}

// SetTabLocation sets the location of the tab bar
func (t *AppTabs) SetTabLocation(l TabLocation) {
	t.baseTabs.Items = t.Items
	// Mobile has limited screen space, so don't put app tab bar on long edges
	if d := fyne.CurrentDevice(); d.IsMobile() {
		if o := d.Orientation(); fyne.IsVertical(o) {
			if l == TabLocationLeading || l == TabLocationTrailing {
				l = TabLocationBottom
			}
		} else {
			if l == TabLocationTop || l == TabLocationBottom {
				l = TabLocationLeading
			}
		}
	}
	t.baseTabs.SetTabLocation(l)
}

// Show this widget, if it was previously hidden
//
// Deprecated: Moved into an embedded type.
func (t *AppTabs) Show() {
	t.baseTabs.Items = t.Items
	t.baseTabs.Show()
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*appTabsRenderer)(nil)

type appTabsRenderer struct {
	baseTabsRenderer
	appTabs *AppTabs
}

func (r *appTabsRenderer) Layout(size fyne.Size) {
	// Try render as many tabs as will fit (capped to maxAppTab), others will appear in the overflow
	for i := maxAppTabs; i > 0; i-- {
		r.updateTabs(i)
		barMin := r.bar.MinSize()
		if r.appTabs.tabLocation == TabLocationLeading || r.appTabs.tabLocation == TabLocationTrailing {
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

	r.layout(&r.appTabs.baseTabs, size)
	r.updateIndicator()
}

func (r *appTabsRenderer) MinSize() fyne.Size {
	return r.minSize(&r.appTabs.baseTabs)
}

func (r *appTabsRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseTabsRenderer.Objects()
	if i, is := r.appTabs.current, r.appTabs.baseTabs.Items; i >= 0 && i < len(is) {
		objects = append(objects, is[i].Content)
	}
	return objects
}

func (r *appTabsRenderer) Refresh() {
	r.Layout(r.appTabs.Size())

	r.refresh(&r.appTabs.baseTabs)

	canvas.Refresh(r.appTabs)
}

func (r *appTabsRenderer) buildOverflowTabsButton() (overflow *widget.Button) {
	overflow = widget.NewButton("", func() {
		// Show pop up containing all tabs which did not fit in the tab bar

		var items []*fyne.MenuItem
		for i := len(r.bar.Objects[0].(*fyne.Container).Objects); i < len(r.appTabs.baseTabs.Items); i++ {
			item := r.appTabs.baseTabs.Items[i]
			// FIXME MenuItem doesn't support icons (#1752)
			// FIXME MenuItem can't show if it is the currently selected tab (#1753)
			items = append(items, fyne.NewMenuItem(item.Text, func() {
				r.appTabs.Select(item)
				r.appTabs.popUp = nil
			}))
		}

		r.appTabs.showPopUp(overflow, items)
	})
	overflow.Importance = widget.LowImportance
	return
}

func (r *appTabsRenderer) buildTabButtons(count int) *fyne.Container {
	buttons := &fyne.Container{}

	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() {
		cells := count
		if cells == 0 {
			cells = 1
		}
		buttons.Layout = layout.NewGridLayout(cells)
		iconPos = buttonIconTop
	} else if r.appTabs.tabLocation == TabLocationLeading || r.appTabs.tabLocation == TabLocationTrailing {
		buttons.Layout = layout.NewVBoxLayout()
		iconPos = buttonIconTop
	} else {
		buttons.Layout = layout.NewHBoxLayout()
		iconPos = buttonIconInline
	}

	for i := 0; i < count; i++ {
		item := r.appTabs.baseTabs.Items[i]
		button, ok := r.buttonCache[item]
		if !ok {
			button = &tabButton{
				onTapped: func() { r.appTabs.Select(item) },
			}
			r.buttonCache[item] = button
		}
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

func (r *appTabsRenderer) updateIndicator() {
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
	} else if r.appTabs.current >= 0 {
		selected := buttons[r.appTabs.current]
		selectedPos = selected.Position()
		selectedSize = selected.Size()
	}

	var indicatorPos fyne.Position
	var indicatorSize fyne.Size

	switch r.appTabs.tabLocation {
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

	r.moveIndicator(indicatorPos, indicatorSize, true)
}

func (r *appTabsRenderer) updateTabs(max int) {
	tabCount := len(r.appTabs.baseTabs.Items)

	// Set overflow action
	if tabCount < max {
		r.action = nil
		r.bar.Layout = layout.NewMaxLayout()
	} else {
		tabCount = max
		if r.action == nil {
			r.action = r.buildOverflowTabsButton()
		}
		// Set layout of tab bar containing tab buttons and overflow action
		if r.appTabs.tabLocation == TabLocationLeading || r.appTabs.tabLocation == TabLocationTrailing {
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
