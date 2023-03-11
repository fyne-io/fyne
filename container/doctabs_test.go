package container_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestDocTabs_Selected(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewDocTabs(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.Selected())
}

func TestDocTabs_SelectedIndex(t *testing.T) {
	tabs := container.NewDocTabs(&container.TabItem{Text: "Test", Content: widget.NewLabel("Test")})

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.SelectedIndex())
}

func TestDocTabs_Empty(t *testing.T) {
	tabs := container.NewDocTabs()
	assert.Equal(t, 0, len(tabs.Items))
	assert.Equal(t, -1, tabs.SelectedIndex())
	assert.Nil(t, tabs.Selected())
	min := tabs.MinSize()
	assert.Equal(t, float32(0), min.Width)
	assert.Equal(t, 4*theme.Padding()+theme.IconInlineSize(), min.Height)

	tabs = &container.DocTabs{}
	assert.Equal(t, 0, len(tabs.Items))
	assert.Nil(t, tabs.Selected())
	assert.NotNil(t, test.WidgetRenderer(tabs)) // doesn't crash
}

func TestDocTabs_Hidden_AsChild(t *testing.T) {
	c1 := widget.NewLabel("Tab 1 content")
	c2 := widget.NewLabel("Tab 2 content\nTab 2 content\nTab 2 content")
	ti1 := container.NewTabItem("Tab 1", c1)
	ti2 := container.NewTabItem("Tab 2", c2)
	tabs := container.NewDocTabs(ti1, ti2)
	tabs.Refresh()

	assert.True(t, c1.Visible())
	assert.False(t, c2.Visible())

	tabs.SelectIndex(1)
	assert.False(t, c1.Visible())
	assert.True(t, c2.Visible())
}

func TestDocTabs_Resize_Empty(t *testing.T) {
	tabs := container.NewDocTabs()
	tabs.Resize(fyne.NewSize(10, 10))
	size := tabs.Size()
	assert.Equal(t, float32(10), size.Height)
	assert.Equal(t, float32(10), size.Width)
}

func TestDocTabs_Select(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewDocTabs(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.Selected())

	var selectedTab *container.TabItem
	tabs.OnSelected = func(tab *container.TabItem) {
		selectedTab = tab
	}
	var unselectedTab *container.TabItem
	tabs.OnUnselected = func(tab *container.TabItem) {
		unselectedTab = tab
	}
	tabs.Select(tab2)
	assert.Equal(t, tab2, tabs.Selected())
	assert.Equal(t, tab2, selectedTab)
	assert.Equal(t, tab1, unselectedTab)

	tabs.OnSelected = func(tab *container.TabItem) {
		assert.Fail(t, "unexpected tab selected")
	}
	tabs.OnUnselected = func(tab *container.TabItem) {
		assert.Fail(t, "unexpected tab unselected")
	}
	tabs.Select(container.NewTabItem("Test3", widget.NewLabel("Test3")))
	assert.Equal(t, tab2, tabs.Selected())
}

func TestDocTabs_SelectIndex(t *testing.T) {
	tabs := container.NewDocTabs(&container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.SelectedIndex())

	var selectedTab *container.TabItem
	tabs.OnSelected = func(tab *container.TabItem) {
		selectedTab = tab
	}
	var unselectedTab *container.TabItem
	tabs.OnUnselected = func(tab *container.TabItem) {
		unselectedTab = tab
	}
	tabs.SelectIndex(1)
	assert.Equal(t, 1, tabs.SelectedIndex())
	assert.Equal(t, tabs.Items[1], selectedTab)
	assert.Equal(t, tabs.Items[0], unselectedTab)
}

func TestDocTabs_RemoveIndex(t *testing.T) {
	tabs := container.NewDocTabs(&container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	tabs.SelectIndex(1)
	tabs.RemoveIndex(1)
	assert.Equal(t, 0, tabs.SelectedIndex()) // check max item selected and no panic

	tabs.SelectIndex(0)
	tabs.RemoveIndex(0)
	assert.Equal(t, -1, tabs.SelectedIndex()) // check deselected and no panic
}

func TestDocTabs_EnableItem(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewDocTabs(tab1, tab2)
	tabs.DisableItem(tab1)

	assert.True(t, tab1.Disabled())

	tabs.EnableItem(tab1)
	assert.False(t, tab1.Disabled())
}

func TestDocTabs_EnableIndex(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewDocTabs(tab1, tab2)
	tabs.DisableItem(tab1)

	assert.True(t, tab1.Disabled())

	tabs.EnableIndex(0)
	assert.False(t, tab1.Disabled())
}

func TestDocTabs_DisableItem(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewDocTabs(tab1, tab2)

	assert.False(t, tab1.Disabled())

	tabs.DisableItem(tab1)
	assert.True(t, tab1.Disabled())

	assert.Equal(t, 1, tabs.SelectedIndex())
}

func TestDocTabs_DisableIndex(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewDocTabs(tab1, tab2)

	assert.False(t, tab1.Disabled())

	tabs.DisableIndex(0)
	assert.True(t, tab1.Disabled())

	assert.Equal(t, 1, tabs.SelectedIndex())
}
