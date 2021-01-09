package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const MAX_APP_TABS = 7

// Declare conformity with Widget interface.
var _ fyne.Widget = (*AppTabs)(nil)

// AppTabs container is used to split your application into various different areas identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 2.0.0
type AppTabs struct {
	baseTabs
}

// NewAppTabs creates a new tab container that allows the user to choose between different areas of an app.
//
// Since: 2.0.0
func NewAppTabs(items ...*TabItem) *AppTabs {
	tabs := &AppTabs{
		baseTabs: baseTabs{
			BaseWidget: widget.BaseWidget{},
			current:    -1,
		},
	}
	tabs.ExtendBaseWidget(tabs)
	tabs.SetItems(items)
	return tabs
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *AppTabs) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &appTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			bar:         &fyne.Container{},
			buttonCache: make(map[*TabItem]*tabButton),
			divider:     canvas.NewRectangle(theme.ShadowColor()),
			indicator:   canvas.NewRectangle(theme.PrimaryColor()),
		},
		appTabs: t,
	}
	// Initially setup the tab bar to only show one tab, all others will be in overflow.
	// When the widget is laid out, and we know the size, the tab bar will be updated to show as many as can fit.
	r.updateTabs(1)
	r.moveIndicator()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (t *AppTabs) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// SetTabLocation sets the location of the tab bar
func (t *AppTabs) SetTabLocation(l TabLocation) {
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

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*appTabsRenderer)(nil)

type appTabsRenderer struct {
	baseTabsRenderer
	appTabs *AppTabs
}

func (r *appTabsRenderer) Layout(size fyne.Size) {
	// Try render as many tabs as will fit, others will appear in the overflow
	for i := MAX_APP_TABS; i > 0; i-- {
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
	r.moveIndicator()
}

func (r *appTabsRenderer) MinSize() fyne.Size {
	return r.minSize(&r.appTabs.baseTabs)
}

func (r *appTabsRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseTabsRenderer.Objects()
	if i, is := r.appTabs.current, r.appTabs.Items; i >= 0 && i < len(is) {
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
		for i := len(r.bar.Objects[0].(*fyne.Container).Objects); i < len(r.appTabs.Items); i++ {
			item := r.appTabs.Items[i]
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

func (r *appTabsRenderer) moveIndicator() {
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

	r.animateIndicator(indicatorPos, indicatorSize)
}

func (r *appTabsRenderer) updateTabs(max int) {
	tabCount := len(r.appTabs.Items)

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

	buttons := r.buildTabButtons(&r.appTabs.baseTabs, tabCount)

	r.bar.Objects = []fyne.CanvasObject{buttons}
	if a := r.action; a != nil {
		r.bar.Objects = append(r.bar.Objects, a)
	}

	r.bar.Refresh()
}
