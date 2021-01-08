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
	popUp *widget.PopUpMenu
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
			bar: &tabBar{
				buttons: &fyne.Container{},
			},
			buttons:   make(map[*TabItem]*tabButton),
			divider:   canvas.NewRectangle(theme.ShadowColor()),
			indicator: canvas.NewRectangle(theme.PrimaryColor()),
		},
		appTabs: t,
	}
	r.updateTabs()
	return r
}

// Hide hides the select.
//
// Implements: fyne.Widget
func (t *AppTabs) Hide() {
	if t.popUp != nil {
		t.popUp.Hide()
		t.popUp = nil
	}
	t.BaseWidget.Hide()
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
	barMin := r.bar.MinSize()

	var (
		barPos, dividerPos, contentPos    fyne.Position
		barSize, dividerSize, contentSize fyne.Size
	)

	switch r.appTabs.tabLocation {
	case TabLocationTop:
		barHeight := barMin.Height
		barPos = fyne.NewPos(0, 0)
		barSize = fyne.NewSize(size.Width, barHeight)
		dividerPos = fyne.NewPos(0, barHeight)
		dividerSize = fyne.NewSize(size.Width, theme.Padding())
		contentPos = fyne.NewPos(0, barHeight+theme.Padding())
		contentSize = fyne.NewSize(size.Width, size.Height-barHeight-theme.Padding())
	case TabLocationLeading:
		barWidth := barMin.Width
		barPos = fyne.NewPos(0, 0)
		barSize = fyne.NewSize(barWidth, size.Height)
		dividerPos = fyne.NewPos(barWidth, 0)
		dividerSize = fyne.NewSize(theme.Padding(), size.Height)
		contentPos = fyne.NewPos(barWidth+theme.Padding(), 0)
		contentSize = fyne.NewSize(size.Width-barWidth-theme.Padding(), size.Height)
	case TabLocationBottom:
		barHeight := barMin.Height
		barPos = fyne.NewPos(0, size.Height-barHeight)
		barSize = fyne.NewSize(size.Width, barHeight)
		dividerPos = fyne.NewPos(0, size.Height-barHeight-theme.Padding())
		dividerSize = fyne.NewSize(size.Width, theme.Padding())
		contentPos = fyne.NewPos(0, 0)
		contentSize = fyne.NewSize(size.Width, size.Height-barHeight-theme.Padding())
	case TabLocationTrailing:
		barWidth := barMin.Width
		barPos = fyne.NewPos(size.Width-barWidth, 0)
		barSize = fyne.NewSize(barWidth, size.Height)
		dividerPos = fyne.NewPos(size.Width-barWidth-theme.Padding(), 0)
		dividerSize = fyne.NewSize(theme.Padding(), size.Height)
		contentPos = fyne.NewPos(0, 0)
		contentSize = fyne.NewSize(size.Width-barWidth-theme.Padding(), size.Height)
	}

	r.bar.Move(barPos)
	r.bar.Resize(barSize)
	r.divider.Move(dividerPos)
	r.divider.Resize(dividerSize)
	if r.appTabs.current >= 0 && r.appTabs.current < len(r.appTabs.Items) {
		content := r.appTabs.Items[r.appTabs.current].Content
		content.Move(contentPos)
		content.Resize(contentSize)
	}
}

func (r *appTabsRenderer) MinSize() (min fyne.Size) {
	barMin := r.bar.MinSize()

	contentMin := fyne.NewSize(0, 0)
	for _, content := range r.appTabs.Items {
		contentMin = contentMin.Max(content.Content.MinSize())
	}

	switch r.appTabs.tabLocation {
	case TabLocationLeading, TabLocationTrailing:
		return fyne.NewSize(barMin.Width+contentMin.Width+theme.Padding(),
			fyne.Max(barMin.Height, contentMin.Height))
	default:
		return fyne.NewSize(fyne.Max(barMin.Width, contentMin.Width),
			barMin.Height+contentMin.Height+theme.Padding())
	}
}

func (r *appTabsRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bar, r.divider, r.indicator}
	if i, is := r.appTabs.current, r.appTabs.Items; i >= 0 && i < len(is) {
		objects = append(objects, is[i].Content)
	}
	return objects
}

func (r *appTabsRenderer) Refresh() {
	r.updateTabs()
	r.divider.FillColor = theme.ShadowColor()
	r.divider.Refresh()
	r.indicator.FillColor = theme.PrimaryColor()
	r.Layout(r.appTabs.Size())
	canvas.Refresh(r.appTabs)
}

func (r *appTabsRenderer) buildOverflow() (overflow *widget.Button) {
	overflow = widget.NewButtonWithIcon("", theme.MenuExpandIcon() /* TODO OverflowIcon() */, func() {
		// Show pop up containing all tabs which did not fit in the tab bar

		var items []*fyne.MenuItem
		for i := len(r.bar.buttons.Objects); i < len(r.appTabs.Items); i++ {
			item := r.appTabs.Items[i]
			// FIXME MenuItem doesn't support icons
			items = append(items, fyne.NewMenuItem(item.Text, func() {
				r.appTabs.Select(item)
				r.appTabs.popUp = nil
			}))
		}
		d := fyne.CurrentApp().Driver()
		c := d.CanvasForObject(overflow)
		r.appTabs.popUp = widget.NewPopUpMenu(fyne.NewMenu("", items...), c)
		buttonPos := d.AbsolutePositionForObject(overflow)
		buttonMin := overflow.MinSize()
		popUpMin := r.appTabs.popUp.MinSize()
		var popUpPos fyne.Position
		switch r.appTabs.tabLocation {
		case TabLocationLeading:
			popUpPos.X = buttonPos.X + buttonMin.Width
			popUpPos.Y = buttonPos.Y + buttonMin.Height - popUpMin.Height
		case TabLocationTrailing:
			popUpPos.X = buttonPos.X - popUpMin.Width
			popUpPos.Y = buttonPos.Y + buttonMin.Height - popUpMin.Height
		case TabLocationTop:
			popUpPos.X = buttonPos.X + buttonMin.Width - popUpMin.Width
			popUpPos.Y = buttonPos.Y + buttonMin.Height
		case TabLocationBottom:
			popUpPos.X = buttonPos.X + buttonMin.Width - popUpMin.Width
			popUpPos.Y = buttonPos.Y - popUpMin.Height
		}
		r.appTabs.popUp.ShowAtPosition(popUpPos)
	})
	return
}

func (r *appTabsRenderer) updateTabs() {
	tabCount := len(r.appTabs.Items)

	// Set overflow action
	if tabCount < MAX_APP_TABS {
		r.bar.action = nil
		r.bar.layout = layout.NewMaxLayout()
	} else if r.bar.action == nil {
		tabCount = MAX_APP_TABS
		r.bar.action = r.buildOverflow()
		// Set layout of tab bar containing tab buttons and overflow action
		if r.appTabs.tabLocation == TabLocationLeading || r.appTabs.tabLocation == TabLocationTrailing {
			r.bar.layout = layout.NewBorderLayout(nil, r.bar.action, nil, nil)
		} else {
			r.bar.layout = layout.NewBorderLayout(nil, nil, nil, r.bar.action)
		}
	}

	// Set tab buttons
	var iconPos buttonIconPosition
	if fyne.CurrentDevice().IsMobile() {
		cells := tabCount
		if cells == 0 {
			cells = 1
		}
		if cells >= MAX_APP_TABS {
			cells = MAX_APP_TABS
		}
		r.bar.buttons.Layout = layout.NewGridLayout(cells)
		iconPos = buttonIconTop
	} else if r.appTabs.tabLocation == TabLocationLeading || r.appTabs.tabLocation == TabLocationTrailing {
		r.bar.buttons.Layout = layout.NewVBoxLayout()
		iconPos = buttonIconTop
	} else {
		r.bar.buttons.Layout = layout.NewHBoxLayout()
		iconPos = buttonIconInline
	}
	r.bar.buttons.Objects = nil
	for i := 0; i < tabCount; i++ {
		item := r.appTabs.Items[i]
		button, ok := r.buttons[item]
		if !ok {
			button = &tabButton{
				OnTap: func() { r.appTabs.Select(item) },
			}
			r.buttons[item] = button
		}
		button.Text = item.Text
		button.Icon = item.Icon
		button.IconPosition = iconPos
		if i == r.appTabs.current {
			button.Importance = widget.HighImportance
		} else {
			button.Importance = widget.MediumImportance
		}
		r.bar.buttons.Objects = append(r.bar.buttons.Objects, button)
	}
	r.bar.buttons.Refresh()
	r.moveIndicator(r.appTabs.tabLocation, r.appTabs.current)
}
