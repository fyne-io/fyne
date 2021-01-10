// +build !mobile

package widget_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestTabContainer_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(
		widget.NewTabContainer(&widget.TabItem{Text: "Test", Content: widget.NewLabel("Text")}),
	)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertImageMatches(t, "tabcontainer/desktop/single_initial.png", c.Capture())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "tabcontainer/desktop/single_custom_theme.png", c.Capture())
}

func TestTabContainer_ChangeItemContent(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_content_initial.xml", c)

	item1.Content = widget.NewLabel("Text3")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_content_change_visible.xml", c)

	item2.Content = widget.NewLabel("Text4")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_content_change_hidden.xml", c)
}

func TestTabContainer_ChangeItemIcon(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Icon: theme.CancelIcon(), Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Icon: theme.ConfirmIcon(), Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_icon_initial.xml", c)

	item1.Icon = theme.InfoIcon()
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_icon_change_selected.xml", c)

	item2.Icon = theme.ContentAddIcon()
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_icon_change_unselected.xml", c)
}

func TestTabContainer_ChangeItemText(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_label_initial.xml", c)

	item1.Text = "New 1"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_label_change_selected.xml", c)

	item2.Text = "New 2"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/change_label_change_unselected.xml", c)
}

func TestTabContainer_DynamicTabs(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	tabs := widget.NewTabContainer(item1)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "tabcontainer/desktop/dynamic_initial.xml", c)

	appendedItem := widget.NewTabItem("Test2", widget.NewLabel("Text 2"))
	tabs.Append(appendedItem)
	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, "Test2", tabs.Items[1].Text)
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/dynamic_appended.xml", c)

	tabs.RemoveIndex(1)
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, "Test1", tabs.Items[0].Text)
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/dynamic_initial.xml", c)

	tabs.Append(appendedItem)
	tabs.Remove(tabs.Items[0])
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, "Test2", tabs.Items[0].Text)
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/dynamic_appended_and_removed.xml", c)

	tabs.Append(widget.NewTabItem("Test3", canvas.NewCircle(theme.BackgroundColor())))
	tabs.Append(widget.NewTabItem("Test4", canvas.NewCircle(theme.BackgroundColor())))
	tabs.Append(widget.NewTabItem("Test5", canvas.NewCircle(theme.BackgroundColor())))
	assert.Equal(t, 4, len(tabs.Items))
	assert.Equal(t, "Test3", tabs.Items[1].Text)
	assert.Equal(t, "Test4", tabs.Items[2].Text)
	assert.Equal(t, "Test5", tabs.Items[3].Text)
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/dynamic_appended_another_three.xml", c)

	tabs.SetItems([]*widget.TabItem{
		widget.NewTabItem("Test6", widget.NewLabel("Text 6")),
		widget.NewTabItem("Test7", widget.NewLabel("Text 7")),
		widget.NewTabItem("Test8", widget.NewLabel("Text 8")),
	})
	assert.Equal(t, 3, len(tabs.Items))
	assert.Equal(t, "Test6", tabs.Items[0].Text)
	assert.Equal(t, "Test7", tabs.Items[1].Text)
	assert.Equal(t, "Test8", tabs.Items[2].Text)
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/dynamic_replaced_completely.xml", c)
}

func TestTabContainer_HoverButtons(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "tabcontainer/desktop/hover_none.xml", c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/hover_first.xml", c)

	test.MoveMouse(c, fyne.NewPos(75, 10))
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/hover_second.xml", c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/hover_first.xml", c)
}

func TestTabContainer_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(nil)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	tests := []struct {
		name     string
		item     *widget.TabItem
		location widget.TabLocation
		want     string
	}{
		{
			name:     "top: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTop,
			want:     "tabcontainer/desktop/layout_top_icon_and_text.xml",
		},
		{
			name:     "top: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTop,
			want:     "tabcontainer/desktop/layout_top_text.xml",
		},
		{
			name:     "top: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTop,
			want:     "tabcontainer/desktop/layout_top_icon.xml",
		},
		{
			name:     "bottom: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want:     "tabcontainer/desktop/layout_bottom_icon_and_text.xml",
		},
		{
			name:     "bottom: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want:     "tabcontainer/desktop/layout_bottom_text.xml",
		},
		{
			name:     "bottom: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want:     "tabcontainer/desktop/layout_bottom_icon.xml",
		},
		{
			name:     "leading: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want:     "tabcontainer/desktop/layout_leading_icon_and_text.xml",
		},
		{
			name:     "leading: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want:     "tabcontainer/desktop/layout_leading_text.xml",
		},
		{
			name:     "leading: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want:     "tabcontainer/desktop/layout_leading_icon.xml",
		},
		{
			name:     "trailing: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want:     "tabcontainer/desktop/layout_trailing_icon_and_text.xml",
		},
		{
			name:     "trailing: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want:     "tabcontainer/desktop/layout_trailing_text.xml",
		},
		{
			name:     "trailing: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want:     "tabcontainer/desktop/layout_trailing_icon.xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := widget.NewTabContainer(tt.item)
			tabs.SetTabLocation(tt.location)
			w.SetContent(tabs)
			w.Resize(fyne.NewSize(150, 150))

			test.AssertRendersToMarkup(t, tt.want, c)
		})
	}
}

func TestTabContainer_SetTabLocation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &widget.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := widget.NewTabContainer(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tab_location_top.xml", c)

	tabs.SetTabLocation(widget.TabLocationLeading)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tab_location_leading.xml", c)

	tabs.SetTabLocation(widget.TabLocationBottom)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tab_location_bottom.xml", c)

	tabs.SetTabLocation(widget.TabLocationTrailing)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tab_location_trailing.xml", c)

	tabs.SetTabLocation(widget.TabLocationTop)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tab_location_top.xml", c)
}

func TestTabContainer_Tapped(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &widget.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := widget.NewTabContainer(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(200, 100))
	c := w.Canvas()

	require.Equal(t, 0, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tapped_first_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(75, 10))
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tapped_second_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(150, 10))
	assert.Equal(t, 2, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tapped_third_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(10, 10))
	require.Equal(t, 0, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, "tabcontainer/desktop/tapped_first_selected.xml", c)
}
