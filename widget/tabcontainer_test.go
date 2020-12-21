package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestTabContainer_CurrentTab(t *testing.T) {
	tab1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := widget.NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())
}

func TestTabContainer_CurrentTabIndex(t *testing.T) {
	tabs := widget.NewTabContainer(&widget.TabItem{Text: "Test", Content: widget.NewLabel("Test")})

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())
}

func TestTabContainer_Empty(t *testing.T) {
	tabs := widget.NewTabContainer()
	assert.Equal(t, 0, len(tabs.Items))
	assert.Equal(t, -1, tabs.CurrentTabIndex())
	assert.Nil(t, tabs.CurrentTab())
	min := tabs.MinSize()
	assert.Equal(t, 0, min.Width)
	assert.Equal(t, theme.Padding(), min.Height)
}

func TestTabContainer_Hidden_AsChild(t *testing.T) {
	c1 := widget.NewLabel("Tab 1 content")
	c2 := widget.NewLabel("Tab 2 content\nTab 2 content\nTab 2 content")
	ti1 := widget.NewTabItem("Tab 1", c1)
	ti2 := widget.NewTabItem("Tab 2", c2)
	tabs := widget.NewTabContainer(ti1, ti2)
	tabs.Refresh()

	assert.True(t, c1.Visible())
	assert.False(t, c2.Visible())

	tabs.SelectTabIndex(1)
	assert.False(t, c1.Visible())
	assert.True(t, c2.Visible())
}

func TestTabContainer_Resize_Empty(t *testing.T) {
	tabs := widget.NewTabContainer()
	tabs.Resize(fyne.NewSize(10, 10))
	size := tabs.Size()
	assert.Equal(t, 10, size.Height)
	assert.Equal(t, 10, size.Width)
}

func TestTabContainer_SelectTab(t *testing.T) {
	tab1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := widget.NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())

	var selectedTab *widget.TabItem
	tabs.OnChanged = func(tab *widget.TabItem) {
		selectedTab = tab
	}
	tabs.SelectTab(tab2)
	assert.Equal(t, tab2, tabs.CurrentTab())
	assert.Equal(t, tab2, selectedTab)

	tabs.OnChanged = func(tab *widget.TabItem) {
		assert.Fail(t, "unexpected tab selection")
	}
	tabs.SelectTab(widget.NewTabItem("Test3", widget.NewLabel("Test3")))
	assert.Equal(t, tab2, tabs.CurrentTab())
}

func TestTabContainer_SelectTabIndex(t *testing.T) {
	tabs := widget.NewTabContainer(&widget.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&widget.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())

	var selectedTab *widget.TabItem
	tabs.OnChanged = func(tab *widget.TabItem) {
		selectedTab = tab
	}
	tabs.SelectTabIndex(1)
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	assert.Equal(t, tabs.Items[1], selectedTab)
}

func TestTabContainer_RemoveIndex(t *testing.T) {
	tabs := widget.NewTabContainer(&widget.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&widget.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	tabs.SelectTabIndex(1)
	tabs.RemoveIndex(1)
	assert.Equal(t, 0, tabs.CurrentTabIndex()) // check max item selection and no panic

	tabs.SelectTabIndex(0)
	tabs.RemoveIndex(0)
	assert.Equal(t, -1, tabs.CurrentTabIndex()) // check deselection and no panic
}
