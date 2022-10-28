//go:build !mobile
// +build !mobile

package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestDocTabs_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(
		container.NewDocTabs(&container.TabItem{Text: "Test", Content: widget.NewLabel("Text")}),
	)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToImage(t, "doctabs/desktop/single_initial.png", c)

	test.ApplyTheme(t, test.NewTheme())
	test.AssertRendersToImage(t, "doctabs/desktop/single_custom_theme.png", c)
}

func TestDocTabs_ChangeItemContent(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := container.NewDocTabs(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "doctabs/desktop/change_content_initial.xml", c)

	item1.Content = widget.NewLabel("Text3")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_content_change_visible.xml", c)

	item2.Content = widget.NewLabel("Text4")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_content_change_hidden.xml", c)
}

func TestDocTabs_ChangeItemIcon(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Icon: theme.CancelIcon(), Content: widget.NewLabel("Text1")}
	item2 := &container.TabItem{Icon: theme.ConfirmIcon(), Content: widget.NewLabel("Text2")}
	tabs := container.NewDocTabs(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "doctabs/desktop/change_icon_initial.xml", c)

	item1.Icon = theme.InfoIcon()
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_icon_change_selected.xml", c)

	item2.Icon = theme.ContentAddIcon()
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_icon_change_unselected.xml", c)
}

func TestDocTabs_ChangeItemText(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := container.NewDocTabs(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "doctabs/desktop/change_label_initial.xml", c)

	item1.Text = "New 1"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_label_change_selected.xml", c)

	item2.Text = "New 2"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_label_change_unselected.xml", c)

	// use bigger window to have indicator be based on tab size
	w.Resize(fyne.NewSize(500, 150))
	item1.Text = "New longer text 1"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, "doctabs/desktop/change_label_to_longer_text_selected.xml", c)
}

func TestDocTabs_DynamicTabs(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	tabs := container.NewDocTabs(item1)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "doctabs/desktop/dynamic_initial.xml", c)

	appendedItem := container.NewTabItem("Test2", widget.NewLabel("Text 2"))
	tabs.Append(appendedItem)
	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, "Test2", tabs.Items[1].Text)
	test.AssertRendersToMarkup(t, "doctabs/desktop/dynamic_appended.xml", c)

	tabs.RemoveIndex(1)
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, "Test1", tabs.Items[0].Text)
	test.AssertRendersToMarkup(t, "doctabs/desktop/dynamic_initial.xml", c)

	tabs.Append(appendedItem)
	tabs.Remove(tabs.Items[0])
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, "Test2", tabs.Items[0].Text)
	test.AssertRendersToMarkup(t, "doctabs/desktop/dynamic_appended_and_removed.xml", c)

	tabs.Append(container.NewTabItem("Test3", canvas.NewCircle(theme.BackgroundColor())))
	tabs.Append(container.NewTabItem("Test4", canvas.NewCircle(theme.BackgroundColor())))
	tabs.Append(container.NewTabItem("Test5", canvas.NewCircle(theme.BackgroundColor())))
	assert.Equal(t, 4, len(tabs.Items))
	assert.Equal(t, "Test3", tabs.Items[1].Text)
	assert.Equal(t, "Test4", tabs.Items[2].Text)
	assert.Equal(t, "Test5", tabs.Items[3].Text)
	test.AssertRendersToMarkup(t, "doctabs/desktop/dynamic_appended_another_three.xml", c)

	tabs.SetItems([]*container.TabItem{
		container.NewTabItem("Test6", widget.NewLabel("Text 6")),
		container.NewTabItem("Test7", widget.NewLabel("Text 7")),
		container.NewTabItem("Test8", widget.NewLabel("Text 8")),
	})
	assert.Equal(t, 3, len(tabs.Items))
	assert.Equal(t, "Test6", tabs.Items[0].Text)
	assert.Equal(t, "Test7", tabs.Items[1].Text)
	assert.Equal(t, "Test8", tabs.Items[2].Text)
	test.AssertRendersToMarkup(t, "doctabs/desktop/dynamic_replaced_completely.xml", c)
}

func TestDocTabs_HoverButtons(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := container.NewDocTabs(item1, item2)
	tabs.CreateTab = func() *container.TabItem {
		return &container.TabItem{Text: "Another", Content: widget.NewLabel("Another Tab")}
	}
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_none.xml", c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_first.xml", c)

	test.MoveMouse(c, fyne.NewPos(75, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_first_close.xml", c)

	test.MoveMouse(c, fyne.NewPos(90, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_second.xml", c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_first.xml", c)

	test.MoveMouse(c, fyne.NewPos(136, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_all_tabs.xml", c)

	test.MoveMouse(c, fyne.NewPos(104, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/hover_create_tab.xml", c)
}

func TestDocTabs_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(nil)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	tests := []struct {
		name     string
		item     *container.TabItem
		location container.TabLocation
		want     string
	}{
		{
			name:     "top: tab with icon and text",
			item:     container.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationTop,
			want:     "doctabs/desktop/layout_top_icon_and_text.xml",
		},
		{
			name:     "top: tab with text only",
			item:     container.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationTop,
			want:     "doctabs/desktop/layout_top_text.xml",
		},
		{
			name:     "top: tab with icon only",
			item:     container.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationTop,
			want:     "doctabs/desktop/layout_top_icon.xml",
		},
		{
			name:     "bottom: tab with icon and text",
			item:     container.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationBottom,
			want:     "doctabs/desktop/layout_bottom_icon_and_text.xml",
		},
		{
			name:     "bottom: tab with text only",
			item:     container.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationBottom,
			want:     "doctabs/desktop/layout_bottom_text.xml",
		},
		{
			name:     "bottom: tab with icon only",
			item:     container.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationBottom,
			want:     "doctabs/desktop/layout_bottom_icon.xml",
		},
		{
			name:     "leading: tab with icon and text",
			item:     container.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationLeading,
			want:     "doctabs/desktop/layout_leading_icon_and_text.xml",
		},
		{
			name:     "leading: tab with text only",
			item:     container.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationLeading,
			want:     "doctabs/desktop/layout_leading_text.xml",
		},
		{
			name:     "leading: tab with icon only",
			item:     container.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationLeading,
			want:     "doctabs/desktop/layout_leading_icon.xml",
		},
		{
			name:     "trailing: tab with icon and text",
			item:     container.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationTrailing,
			want:     "doctabs/desktop/layout_trailing_icon_and_text.xml",
		},
		{
			name:     "trailing: tab with text only",
			item:     container.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationTrailing,
			want:     "doctabs/desktop/layout_trailing_text.xml",
		},
		{
			name:     "trailing: tab with icon only",
			item:     container.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: container.TabLocationTrailing,
			want:     "doctabs/desktop/layout_trailing_icon.xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := container.NewDocTabs(tt.item)
			tabs.SetTabLocation(tt.location)
			w.SetContent(tabs)
			w.Resize(fyne.NewSize(150, 150))

			test.AssertRendersToMarkup(t, tt.want, c)
		})
	}
}

func TestDocTabs_SetTabLocation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &container.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := container.NewDocTabs(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tab_location_top.xml", c)

	tabs.SetTabLocation(container.TabLocationLeading)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tab_location_leading.xml", c)

	tabs.SetTabLocation(container.TabLocationBottom)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tab_location_bottom.xml", c)

	tabs.SetTabLocation(container.TabLocationTrailing)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tab_location_trailing.xml", c)

	tabs.SetTabLocation(container.TabLocationTop)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tab_location_top.xml", c)
}

func TestDocTabs_Tapped(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &container.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := container.NewDocTabs(item1, item2, item3)
	tabs.CreateTab = func() *container.TabItem {
		return &container.TabItem{Text: "Another", Content: widget.NewLabel("Another Tab")}
	}
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 100))
	c := w.Canvas()

	require.Equal(t, 0, tabs.SelectedIndex())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tapped_first_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(90, 10))
	assert.Equal(t, 1, tabs.SelectedIndex())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tapped_second_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(180, 10))
	assert.Equal(t, 2, tabs.SelectedIndex())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tapped_third_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(10, 10))
	require.Equal(t, 0, tabs.SelectedIndex())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tapped_first_selected.xml", c)

	test.TapCanvas(c, fyne.NewPos(254, 10))
	require.Equal(t, 3, tabs.SelectedIndex())
	test.AssertRendersToMarkup(t, "doctabs/desktop/tapped_create_tab.xml", c)

	test.TapCanvas(c, fyne.NewPos(286, 10))
	test.AssertRendersToMarkup(t, "doctabs/desktop/tapped_all_tabs.xml", c)
}
