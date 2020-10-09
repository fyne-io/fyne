package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// AppTabs container is used to split your application into various different areas identified by tabs.
// The tabs contain text and/or an icon and allow the user to switch between the content specified in each TabItem.
// Each item is represented by a button at the edge of the container.
type AppTabs = widget.TabContainer

// TabItem represents a single view in a TabContainer.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
type TabItem = widget.TabItem

// TabLocation is the location where the tabs of a tab container should be rendered
type TabLocation = widget.TabLocation

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeading
	TabLocationBottom
	TabLocationTrailing
)

// NewAppTabs creates a new tab container that allows the user to choose between different areas of an app.
func NewAppTabs(items ...*TabItem) *AppTabs {
	return widget.NewTabContainer(items...)
}

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return widget.NewTabItem(text, content)
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
func NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject) *TabItem {
	return widget.NewTabItemWithIcon(text, icon, content)
}

// TODO move the implementation into here in 2.0 when we delete the old API.
// we cannot do that right now due to Scroll dependency order.
