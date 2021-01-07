package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

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
func (c *AppTabs) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	r := &appTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			divider:   canvas.NewRectangle(theme.ShadowColor()),
			indicator: canvas.NewRectangle(theme.PrimaryColor()),
		},
		appTabs: c,
	}
	// TODO r.updateTabs()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (c *AppTabs) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

// SetTabLocation sets the location of the tab bar
func (c *AppTabs) SetTabLocation(l TabLocation) {
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
	c.baseTabs.SetTabLocation(l)
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*appTabsRenderer)(nil)

type appTabsRenderer struct {
	baseTabsRenderer
	appTabs *AppTabs
}

func (r *appTabsRenderer) Layout(size fyne.Size) {
}

func (r *appTabsRenderer) MinSize() (min fyne.Size) {
	return
}

func (r *appTabsRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bar, r.divider, r.indicator}
	if i, is := r.appTabs.current, r.appTabs.Items; i >= 0 && i < len(is) {
		objects = append(objects, is[i].Content)
	}
	return objects
}

func (r *appTabsRenderer) Refresh() {
}
