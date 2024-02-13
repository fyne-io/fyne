package container

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*DocTabs)(nil)

// DocTabs container is used to display various pieces of content identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 2.1
type DocTabs struct {
	widget.BaseWidget

	Items []*TabItem

	CreateTab      func() *TabItem
	CloseIntercept func(*TabItem)
	OnClosed       func(*TabItem)
	OnSelected     func(*TabItem)
	OnUnselected   func(*TabItem)

	current         int
	location        TabLocation
	isTransitioning bool

	popUpMenu *widget.PopUpMenu
}

// NewDocTabs creates a new tab container that allows the user to choose between various pieces of content.
//
// Since: 2.1
func NewDocTabs(items ...*TabItem) *DocTabs {
	tabs := &DocTabs{}
	tabs.ExtendBaseWidget(tabs)
	tabs.SetItems(items)
	return tabs
}

// Append adds a new TabItem to the end of the tab bar.
func (t *DocTabs) Append(item *TabItem) {
	t.SetItems(append(t.Items, item))
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
//
// Implements: fyne.Widget
func (t *DocTabs) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &docTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			bar:       &fyne.Container{},
			divider:   canvas.NewRectangle(theme.ShadowColor()),
			indicator: canvas.NewRectangle(theme.PrimaryColor()),
		},
		docTabs:  t,
		scroller: NewScroll(&fyne.Container{}),
	}
	r.action = r.buildAllTabsButton()
	r.create = r.buildCreateTabsButton()
	r.tabs = t

	r.box = NewHBox(r.create, r.action)
	r.scroller.OnScrolled = func(offset fyne.Position) {
		r.updateIndicator(false)
	}
	r.updateAllTabs()
	r.updateCreateTab()
	r.updateTabs()
	r.updateIndicator(false)
	r.applyTheme(t)
	return r
}

// DisableIndex disables the TabItem at the specified index.
//
// Since: 2.3
func (t *DocTabs) DisableIndex(i int) {
	disableIndex(t, i)
}

// DisableItem disables the specified TabItem.
//
// Since: 2.3
func (t *DocTabs) DisableItem(item *TabItem) {
	disableItem(t, item)
}

// EnableIndex enables the TabItem at the specified index.
//
// Since: 2.3
func (t *DocTabs) EnableIndex(i int) {
	enableIndex(t, i)
}

// EnableItem enables the specified TabItem.
//
// Since: 2.3
func (t *DocTabs) EnableItem(item *TabItem) {
	enableItem(t, item)
}

// Hide hides the widget.
//
// Implements: fyne.CanvasObject
func (t *DocTabs) Hide() {
	if t.popUpMenu != nil {
		t.popUpMenu.Hide()
		t.popUpMenu = nil
	}
	t.BaseWidget.Hide()
}

// MinSize returns the size that this widget should not shrink below
//
// Implements: fyne.CanvasObject
func (t *DocTabs) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Remove tab by value.
func (t *DocTabs) Remove(item *TabItem) {
	removeItem(t, item)
	t.Refresh()
}

// RemoveIndex removes tab by index.
func (t *DocTabs) RemoveIndex(index int) {
	removeIndex(t, index)
	t.Refresh()
}

// Select sets the specified TabItem to be selected and its content visible.
func (t *DocTabs) Select(item *TabItem) {
	selectItem(t, item)
	t.Refresh()
}

// SelectIndex sets the TabItem at the specific index to be selected and its content visible.
func (t *DocTabs) SelectIndex(index int) {
	selectIndex(t, index)
	t.Refresh()
}

// Selected returns the currently selected TabItem.
func (t *DocTabs) Selected() *TabItem {
	return selected(t)
}

// SelectedIndex returns the index of the currently selected TabItem.
func (t *DocTabs) SelectedIndex() int {
	return t.current
}

// SetItems sets the containers items and refreshes.
func (t *DocTabs) SetItems(items []*TabItem) {
	setItems(t, items)
	t.Refresh()
}

// SetTabLocation sets the location of the tab bar
func (t *DocTabs) SetTabLocation(l TabLocation) {
	t.location = tabsAdjustedLocation(l)
	t.Refresh()
}

// Show this widget, if it was previously hidden
//
// Implements: fyne.CanvasObject
func (t *DocTabs) Show() {
	t.BaseWidget.Show()
	t.SelectIndex(t.current)
}

func (t *DocTabs) close(item *TabItem) {
	if f := t.CloseIntercept; f != nil {
		f(item)
	} else {
		t.Remove(item)
		if f := t.OnClosed; f != nil {
			f(item)
		}
	}
}

func (t *DocTabs) onUnselected() func(*TabItem) {
	return t.OnUnselected
}

func (t *DocTabs) onSelected() func(*TabItem) {
	return t.OnSelected
}

func (t *DocTabs) items() []*TabItem {
	return t.Items
}

func (t *DocTabs) selected() int {
	return t.current
}

func (t *DocTabs) setItems(items []*TabItem) {
	t.Items = items
}

func (t *DocTabs) setSelected(selected int) {
	t.current = selected
}

func (t *DocTabs) setTransitioning(transitioning bool) {
	t.isTransitioning = transitioning
}

func (t *DocTabs) tabLocation() TabLocation {
	return t.location
}

func (t *DocTabs) transitioning() bool {
	return t.isTransitioning
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*docTabsRenderer)(nil)

type docTabsRenderer struct {
	baseTabsRenderer
	docTabs      *DocTabs
	scroller     *Scroll
	box          *fyne.Container
	create       *widget.Button
	lastSelected int
}

func (r *docTabsRenderer) Layout(size fyne.Size) {
	r.updateAllTabs()
	r.updateCreateTab()
	r.updateTabs()
	r.layout(r.docTabs, size)

	// lay out buttons before updating indicator, which is relative to their position
	buttons := r.scroller.Content.(*fyne.Container)
	buttons.Layout.Layout(buttons.Objects, buttons.Size())
	r.updateIndicator(r.docTabs.transitioning())

	if r.docTabs.transitioning() {
		r.docTabs.setTransitioning(false)
	}
}

func (r *docTabsRenderer) MinSize() fyne.Size {
	return r.minSize(r.docTabs)
}

func (r *docTabsRenderer) Objects() []fyne.CanvasObject {
	return r.objects(r.docTabs)
}

func (r *docTabsRenderer) Refresh() {
	r.Layout(r.docTabs.Size())

	if c := r.docTabs.current; c != r.lastSelected {
		if c >= 0 && c < len(r.docTabs.Items) {
			r.scrollToSelected()
		}
		r.lastSelected = c
	}

	r.refresh(r.docTabs)

	canvas.Refresh(r.docTabs)
}

func (r *docTabsRenderer) buildAllTabsButton() (all *widget.Button) {
	all = &widget.Button{Importance: widget.LowImportance, OnTapped: func() {
		// Show pop up containing all tabs

		items := make([]*fyne.MenuItem, len(r.docTabs.Items))
		for i := 0; i < len(r.docTabs.Items); i++ {
			index := i // capture
			// FIXME MenuItem doesn't support icons (#1752)
			items[i] = fyne.NewMenuItem(r.docTabs.Items[i].Text, func() {
				r.docTabs.SelectIndex(index)
				if r.docTabs.popUpMenu != nil {
					r.docTabs.popUpMenu.Hide()
					r.docTabs.popUpMenu = nil
				}
			})
			items[i].Checked = index == r.docTabs.current
		}

		r.docTabs.popUpMenu = buildPopUpMenu(r.docTabs, all, items)
	}}

	return all
}

func (r *docTabsRenderer) buildCreateTabsButton() *widget.Button {
	create := widget.NewButton("", func() {
		if f := r.docTabs.CreateTab; f != nil {
			if tab := f(); tab != nil {
				r.docTabs.Append(tab)
				r.docTabs.SelectIndex(len(r.docTabs.Items) - 1)
			}
		}
	})
	create.Importance = widget.LowImportance
	return create
}

func (r *docTabsRenderer) buildTabButtons(count int, buttons *fyne.Container) {
	buttons.Objects = nil

	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() {
		cells := count
		if cells == 0 {
			cells = 1
		}
		if r.docTabs.location == TabLocationTop || r.docTabs.location == TabLocationBottom {
			buttons.Layout = layout.NewGridLayoutWithColumns(cells)
		} else {
			buttons.Layout = layout.NewGridLayoutWithRows(cells)
		}
		iconPos = buttonIconTop
	} else if r.docTabs.location == TabLocationLeading || r.docTabs.location == TabLocationTrailing {
		buttons.Layout = layout.NewVBoxLayout()
		iconPos = buttonIconTop
	} else {
		buttons.Layout = layout.NewHBoxLayout()
		iconPos = buttonIconInline
	}

	for i := 0; i < count; i++ {
		item := r.docTabs.Items[i]
		if item.button == nil {
			item.button = &tabButton{
				onTapped: func() { r.docTabs.Select(item) },
				onClosed: func() { r.docTabs.close(item) },
			}
		}
		button := item.button
		button.icon = item.Icon
		button.iconPosition = iconPos
		if i == r.docTabs.current {
			button.importance = widget.HighImportance
		} else {
			button.importance = widget.MediumImportance
		}
		button.text = item.Text
		button.textAlignment = fyne.TextAlignLeading
		button.Refresh()
		buttons.Objects = append(buttons.Objects, button)
	}
}

func (r *docTabsRenderer) scrollToSelected() {
	buttons := r.scroller.Content.(*fyne.Container)

	// https://github.com/fyne-io/fyne/issues/3909
	// very dirty temporary fix to this crash!
	if r.docTabs.current < 0 || r.docTabs.current >= len(buttons.Objects) {
		return
	}

	button := buttons.Objects[r.docTabs.current]
	pos := button.Position()
	size := button.Size()
	offset := r.scroller.Offset
	viewport := r.scroller.Size()
	if r.docTabs.location == TabLocationLeading || r.docTabs.location == TabLocationTrailing {
		if pos.Y < offset.Y {
			offset.Y = pos.Y
		} else if pos.Y+size.Height > offset.Y+viewport.Height {
			offset.Y = pos.Y + size.Height - viewport.Height
		}
	} else {
		if pos.X < offset.X {
			offset.X = pos.X
		} else if pos.X+size.Width > offset.X+viewport.Width {
			offset.X = pos.X + size.Width - viewport.Width
		}
	}
	r.scroller.Offset = offset
	r.updateIndicator(false)
}

func (r *docTabsRenderer) updateIndicator(animate bool) {
	if r.docTabs.current < 0 {
		r.indicator.FillColor = color.Transparent
		r.moveIndicator(fyne.NewPos(0, 0), fyne.NewSize(0, 0), animate)
		return
	}

	var selectedPos fyne.Position
	var selectedSize fyne.Size

	buttons := r.scroller.Content.(*fyne.Container).Objects

	if r.docTabs.current >= len(buttons) {
		if a := r.action; a != nil {
			selectedPos = a.Position()
			selectedSize = a.Size()
			minSize := a.MinSize()
			if minSize.Width > selectedSize.Width {
				selectedSize = minSize
			}
		}
	} else {
		selected := buttons[r.docTabs.current]
		selectedPos = selected.Position()
		selectedSize = selected.Size()
		minSize := selected.MinSize()
		if minSize.Width > selectedSize.Width {
			selectedSize = minSize
		}
	}

	scrollOffset := r.scroller.Offset
	scrollSize := r.scroller.Size()

	var indicatorPos fyne.Position
	var indicatorSize fyne.Size

	switch r.docTabs.location {
	case TabLocationTop:
		indicatorPos = fyne.NewPos(selectedPos.X-scrollOffset.X, r.bar.MinSize().Height)
		indicatorSize = fyne.NewSize(fyne.Min(selectedSize.Width, scrollSize.Width-indicatorPos.X), theme.Padding())
	case TabLocationLeading:
		indicatorPos = fyne.NewPos(r.bar.MinSize().Width, selectedPos.Y-scrollOffset.Y)
		indicatorSize = fyne.NewSize(theme.Padding(), fyne.Min(selectedSize.Height, scrollSize.Height-indicatorPos.Y))
	case TabLocationBottom:
		indicatorPos = fyne.NewPos(selectedPos.X-scrollOffset.X, r.bar.Position().Y-theme.Padding())
		indicatorSize = fyne.NewSize(fyne.Min(selectedSize.Width, scrollSize.Width-indicatorPos.X), theme.Padding())
	case TabLocationTrailing:
		indicatorPos = fyne.NewPos(r.bar.Position().X-theme.Padding(), selectedPos.Y-scrollOffset.Y)
		indicatorSize = fyne.NewSize(theme.Padding(), fyne.Min(selectedSize.Height, scrollSize.Height-indicatorPos.Y))
	}

	if indicatorPos.X < 0 {
		indicatorSize.Width = indicatorSize.Width + indicatorPos.X
		indicatorPos.X = 0
	}
	if indicatorPos.Y < 0 {
		indicatorSize.Height = indicatorSize.Height + indicatorPos.Y
		indicatorPos.Y = 0
	}
	if indicatorSize.Width < 0 || indicatorSize.Height < 0 {
		r.indicator.FillColor = color.Transparent
		r.indicator.Refresh()
		return
	}

	r.moveIndicator(indicatorPos, indicatorSize, animate)
}

func (r *docTabsRenderer) updateAllTabs() {
	if len(r.docTabs.Items) > 0 {
		r.action.Show()
	} else {
		r.action.Hide()
	}
}

func (r *docTabsRenderer) updateCreateTab() {
	if r.docTabs.CreateTab != nil {
		r.create.SetIcon(theme.ContentAddIcon())
		r.create.Show()
	} else {
		r.create.Hide()
	}
}

func (r *docTabsRenderer) updateTabs() {
	tabCount := len(r.docTabs.Items)
	r.buildTabButtons(tabCount, r.scroller.Content.(*fyne.Container))

	// Set layout of tab bar containing tab buttons and overflow action
	if r.docTabs.location == TabLocationLeading || r.docTabs.location == TabLocationTrailing {
		r.bar.Layout = layout.NewBorderLayout(nil, r.box, nil, nil)
		r.scroller.Direction = ScrollVerticalOnly
	} else {
		r.bar.Layout = layout.NewBorderLayout(nil, nil, nil, r.box)
		r.scroller.Direction = ScrollHorizontalOnly
	}

	r.bar.Objects = []fyne.CanvasObject{r.scroller, r.box}
	r.bar.Refresh()
}
