package container_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestAppTabs_Selection(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewAppTabs(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.Selection())
}

func TestAppTabs_SelectionIndex(t *testing.T) {
	tabs := container.NewAppTabs(&container.TabItem{Text: "Test", Content: widget.NewLabel("Test")})

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.SelectionIndex())
}

func TestAppTabs_Empty(t *testing.T) {
	tabs := container.NewAppTabs()
	assert.Equal(t, 0, len(tabs.Items))
	assert.Equal(t, -1, tabs.SelectionIndex())
	assert.Nil(t, tabs.Selection())
	min := tabs.MinSize()
	assert.Equal(t, float32(0), min.Width)
	assert.Equal(t, theme.Padding(), min.Height)
}

func TestAppTabs_Hidden_AsChild(t *testing.T) {
	c1 := widget.NewLabel("Tab 1 content")
	c2 := widget.NewLabel("Tab 2 content\nTab 2 content\nTab 2 content")
	ti1 := container.NewTabItem("Tab 1", c1)
	ti2 := container.NewTabItem("Tab 2", c2)
	tabs := container.NewAppTabs(ti1, ti2)
	tabs.Refresh()

	assert.True(t, c1.Visible())
	assert.False(t, c2.Visible())

	tabs.SelectIndex(1)
	assert.False(t, c1.Visible())
	assert.True(t, c2.Visible())
}

func TestAppTabs_Resize_Empty(t *testing.T) {
	tabs := container.NewAppTabs()
	tabs.Resize(fyne.NewSize(10, 10))
	size := tabs.Size()
	assert.Equal(t, float32(10), size.Height)
	assert.Equal(t, float32(10), size.Width)
}

func TestAppTabs_Select(t *testing.T) {
	tab1 := &container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := container.NewAppTabs(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.Selection())

	var selectedTab *container.TabItem
	tabs.OnSelected = func(tab *container.TabItem) {
		selectedTab = tab
	}
	var unselectedTab *container.TabItem
	tabs.OnUnselected = func(tab *container.TabItem) {
		unselectedTab = tab
	}
	tabs.Select(tab2)
	assert.Equal(t, tab2, tabs.Selection())
	assert.Equal(t, tab2, selectedTab)
	assert.Equal(t, tab1, unselectedTab)

	tabs.OnSelected = func(tab *container.TabItem) {
		assert.Fail(t, "unexpected tab selection")
	}
	tabs.OnUnselected = func(tab *container.TabItem) {
		assert.Fail(t, "unexpected tab unselection")
	}
	tabs.Select(container.NewTabItem("Test3", widget.NewLabel("Test3")))
	assert.Equal(t, tab2, tabs.Selection())
}

func TestAppTabs_SelectIndex(t *testing.T) {
	tabs := container.NewAppTabs(&container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.SelectionIndex())

	var selectedTab *container.TabItem
	tabs.OnSelected = func(tab *container.TabItem) {
		selectedTab = tab
	}
	var unselectedTab *container.TabItem
	tabs.OnUnselected = func(tab *container.TabItem) {
		unselectedTab = tab
	}
	tabs.SelectIndex(1)
	assert.Equal(t, 1, tabs.SelectionIndex())
	assert.Equal(t, tabs.Items[1], selectedTab)
	assert.Equal(t, tabs.Items[0], unselectedTab)
}

func TestAppTabs_RemoveIndex(t *testing.T) {
	tabs := container.NewAppTabs(&container.TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&container.TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	tabs.SelectIndex(1)
	tabs.RemoveIndex(1)
	assert.Equal(t, 0, tabs.SelectionIndex()) // check max item selection and no panic

	tabs.SelectIndex(0)
	tabs.RemoveIndex(0)
	assert.Equal(t, -1, tabs.SelectionIndex()) // check deselection and no panic
}
