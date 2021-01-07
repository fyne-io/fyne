package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Declare conformity with Widget interface.
var _ fyne.Widget = (*DocTabs)(nil)

// DocTabs container is used to display various pieces of content identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
//
// Since: 2.0.0
type DocTabs struct {
	baseTabs
}

// NewDocTabs creates a new tab container that allows the user to choose between various pieces of content.
//
// Since: 2.0.0
func NewDocTabs(items ...*TabItem) *DocTabs {
	tabs := &DocTabs{
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
func (t *DocTabs) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &docTabsRenderer{
		baseTabsRenderer: baseTabsRenderer{
			divider:   canvas.NewRectangle(theme.ShadowColor()),
			indicator: canvas.NewRectangle(theme.PrimaryColor()),
		},
		docTabs: t,
	}
	// TODO r.updateTabs()
	return r
}

// MinSize returns the size that this widget should not shrink below
func (t *DocTabs) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*docTabsRenderer)(nil)

type docTabsRenderer struct {
	baseTabsRenderer
	docTabs *DocTabs
}

func (r *docTabsRenderer) Layout(size fyne.Size) {
}

func (r *docTabsRenderer) MinSize() (min fyne.Size) {
	return
}

func (r *docTabsRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bar, r.divider, r.indicator}
	if i, is := r.docTabs.current, r.docTabs.Items; i >= 0 && i < len(is) {
		objects = append(objects, is[i].Content)
	}
	return objects
}

func (r *docTabsRenderer) Refresh() {
}
