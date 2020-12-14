package container

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// TabItem represents a single view in a TabContainer.
// The Text and Icon are used for the tab button and the Content is shown when the corresponding tab is active.
//
// Since: 1.4
type TabItem = widget.TabItem

// TabLocation is the location where the tabs of a tab container should be rendered
//
// Since: 1.4
type TabLocation = widget.TabLocation

// TabLocation values
const (
	TabLocationTop TabLocation = iota
	TabLocationLeading
	TabLocationBottom
	TabLocationTrailing
)

// NewTabItem creates a new item for a tabbed widget - each item specifies the content and a label for its tab.
//
// Since: 1.4
func NewTabItem(text string, content fyne.CanvasObject) *TabItem {
	return widget.NewTabItem(text, content)
}

// NewTabItemWithIcon creates a new item for a tabbed widget - each item specifies the content and a label with an icon for its tab.
//
// Since: 1.4
func NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject) *TabItem {
	return widget.NewTabItemWithIcon(text, icon, content)
}

// TODO move the implementation into here in 2.0 when we delete the old API.
// we cannot do that right now due to Scroll dependency order.

type baseTabs struct {
	widget.BaseWidget

	Items              []*TabItem
	OnSelectionChanged func(tab *TabItem)

	current     int
	tabLocation TabLocation
}

// Append adds a new TabItem to the end of the tab panel
func (c *baseTabs) Append(item *TabItem) {
	c.SetItems(append(c.Items, item))
}

// Remove tab by value
func (c *baseTabs) Remove(item *TabItem) {
	for index, existingItem := range c.Items {
		if existingItem == item {
			c.RemoveIndex(index)
			break
		}
	}
}

// RemoveIndex removes tab by index
func (c *baseTabs) RemoveIndex(index int) {
	if index < 0 || index >= len(c.Items) {
		return
	}
	c.SetItems(append(c.Items[:index], c.Items[index+1:]...))
}

// Select sets the specified TabItem to be selected and its content visible.
func (c *baseTabs) Select(item *TabItem) {
	for i, child := range c.Items {
		if child == item {
			c.SelectIndex(i)
			return
		}
	}
}

// SelectIndex sets the TabItem at the specific index to be selected and its content visible.
func (c *baseTabs) SelectIndex(index int) {
	if index < 0 || index >= len(c.Items) || c.current == index {
		return
	}

	c.current = index
	c.Refresh()

	if f := c.OnSelectionChanged; f != nil {
		f(c.Items[c.current])
	}
}

// Selection returns the currently selected TabItem.
func (c *baseTabs) Selection() *TabItem {
	if c.current < 0 || c.current >= len(c.Items) {
		return nil
	}
	return c.Items[c.current]
}

// SelectionIndex returns the index of the currently selected TabItem.
func (c *baseTabs) SelectionIndex() int {
	return c.current
}

// SetItems sets the container’s items and refreshes.
func (c *baseTabs) SetItems(items []*TabItem) {
	if mismatchedTabItems(items) {
		internal.LogHint("Tab items should all have the same type of content (text, icons or both)")
	}
	c.Items = items
	if len(items) == 0 {
		// No items available to be current
		c.current = -1
	} else if c.current < 0 {
		// Current is first tab item
		c.current = 0
	}
	c.Refresh()
}

// SetTabLocation sets the location of the tab bar
func (c *baseTabs) SetTabLocation(l TabLocation) {
	c.tabLocation = l
	c.Refresh()
}

// Show this widget, if it was previously hidden
func (c *baseTabs) Show() {
	c.BaseWidget.Show()
	c.SelectIndex(c.current)
	c.Refresh()
}

type baseTabsRenderer struct {
	bar                *fyne.Container
	divider, indicator *canvas.Rectangle
}

func (r *baseTabsRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *baseTabsRenderer) Destroy() {
}

func mismatchedTabItems(items []*TabItem) bool {
	var hasText, hasIcon bool
	for _, tab := range items {
		hasText = hasText || tab.Text != ""
		hasIcon = hasIcon || tab.Icon != nil
	}

	mismatch := false
	for _, tab := range items {
		if (hasText && tab.Text == "") || (hasIcon && tab.Icon == nil) {
			mismatch = true
			break
		}
	}

	return mismatch
}
