package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestAppTabs_Selected(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := NewAppTabs(tab1, tab2)

	assert.Len(t, tabs.Items, 2)
	assert.Equal(t, tab1, tabs.Selected())
}

func TestAppTabs_SelectedIndex(t *testing.T) {
	tabs := NewAppTabs(&TabItem{Text: "Test", Content: widget.NewLabel("Test")})

	assert.Len(t, tabs.Items, 1)
	assert.Equal(t, 0, tabs.SelectedIndex())
}

func TestAppTabs_Empty(t *testing.T) {
	tabs := NewAppTabs()
	assert.Empty(t, tabs.Items)
	assert.Equal(t, -1, tabs.SelectedIndex())
	assert.Nil(t, tabs.Selected())
	min := tabs.MinSize()
	assert.Equal(t, float32(0), min.Width)
	assert.Equal(t, theme.Padding(), min.Height)

	tabs = &AppTabs{}
	assert.Empty(t, tabs.Items)
	assert.Nil(t, tabs.Selected())
	assert.NotNil(t, test.TempWidgetRenderer(t, tabs)) // doesn't crash
}

func TestAppTabs_Hidden_AsChild(t *testing.T) {
	c1 := widget.NewLabel("Tab 1 content")
	c2 := widget.NewLabel("Tab 2 content\nTab 2 content\nTab 2 content")
	ti1 := NewTabItem("Tab 1", c1)
	ti2 := NewTabItem("Tab 2", c2)
	tabs := NewAppTabs(ti1, ti2)
	tabs.Refresh()

	assert.True(t, c1.Visible())
	assert.False(t, c2.Visible())

	tabs.SelectIndex(1)
	assert.False(t, c1.Visible())
	assert.True(t, c2.Visible())
}

func TestAppTabs_Resize_Empty(t *testing.T) {
	tabs := NewAppTabs()
	tabs.Resize(fyne.NewSize(10, 10))
	size := tabs.Size()
	assert.Equal(t, float32(10), size.Height)
	assert.Equal(t, float32(10), size.Width)
}

func TestAppTabs_Select(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := NewAppTabs(tab1, tab2)

	assert.Len(t, tabs.Items, 2)
	assert.Equal(t, tab1, tabs.Selected())

	var selectedTab *TabItem
	tabs.OnSelected = func(tab *TabItem) {
		selectedTab = tab
	}
	var unselectedTab *TabItem
	tabs.OnUnselected = func(tab *TabItem) {
		unselectedTab = tab
	}
	tabs.Select(tab2)
	assert.Equal(t, tab2, tabs.Selected())
	assert.Equal(t, tab2, selectedTab)
	assert.Equal(t, tab1, unselectedTab)

	tabs.OnSelected = func(tab *TabItem) {
		assert.Fail(t, "unexpected tab selected")
	}
	tabs.OnUnselected = func(tab *TabItem) {
		assert.Fail(t, "unexpected tab unselected")
	}
	tabs.Select(NewTabItem("Test3", widget.NewLabel("Test3")))
	assert.Equal(t, tab2, tabs.Selected())
}

func TestAppTabs_SelectFocus(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewEntry()}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewEntry()}
	tabs := NewAppTabs(tab1, tab2)
	w := test.NewTempWindow(t, tabs)

	tabs.OnSelected = func(t *TabItem) {
		w.Canvas().Focus(t.Content.(*widget.Entry))
	}

	tabs.Select(tab2)
	assert.Equal(t, tab2.Content, w.Canvas().Focused())
}

func TestAppTabs_SelectIndex(t *testing.T) {
	tabs := NewAppTabs(&TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	assert.Len(t, tabs.Items, 2)
	assert.Equal(t, 0, tabs.SelectedIndex())

	var selectedTab *TabItem
	tabs.OnSelected = func(tab *TabItem) {
		selectedTab = tab
	}
	var unselectedTab *TabItem
	tabs.OnUnselected = func(tab *TabItem) {
		unselectedTab = tab
	}
	tabs.SelectIndex(1)
	assert.Equal(t, 1, tabs.SelectedIndex())
	assert.Equal(t, tabs.Items[1], selectedTab)
	assert.Equal(t, tabs.Items[0], unselectedTab)
}

func TestAppTabs_RemoveIndex(t *testing.T) {
	tabs := NewAppTabs(&TabItem{Text: "Test1", Content: widget.NewLabel("Test1")},
		&TabItem{Text: "Test2", Content: widget.NewLabel("Test2")})

	tabs.SelectIndex(1)
	tabs.RemoveIndex(1)
	assert.Equal(t, 0, tabs.SelectedIndex()) // check max item selected and no panic

	tabs.SelectIndex(0)
	tabs.RemoveIndex(0)
	assert.Equal(t, -1, tabs.SelectedIndex()) // check deselected and no panic
}

func TestAppTabs_EnableItem(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := NewAppTabs(tab1, tab2)
	tabs.DisableItem(tab1)

	assert.True(t, tab1.Disabled())

	tabs.EnableItem(tab1)
	assert.False(t, tab1.Disabled())
}

func TestAppTabs_EnableIndex(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := NewAppTabs(tab1, tab2)
	tabs.DisableItem(tab1)

	assert.True(t, tab1.Disabled())

	tabs.EnableIndex(0)
	assert.False(t, tab1.Disabled())
}

func TestAppTabs_DisableItem(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := NewAppTabs(tab1, tab2)

	assert.False(t, tab1.Disabled())

	tabs.DisableItem(tab1)
	assert.True(t, tab1.Disabled())

	assert.Equal(t, 1, tabs.SelectedIndex())
}

func TestAppTabs_DisableIndex(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: widget.NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: widget.NewLabel("Test2")}
	tabs := NewAppTabs(tab1, tab2)

	assert.False(t, tab1.Disabled())

	tabs.DisableIndex(0)
	assert.True(t, tab1.Disabled())

	assert.Equal(t, 1, tabs.SelectedIndex())
}

func TestAppTabs_ShowAfterAdd(t *testing.T) {
	tabs := NewAppTabs()
	renderer := test.TempWidgetRenderer(t, tabs).(*appTabsRenderer)

	assert.True(t, renderer.indicator.Hidden)

	tabs.Append(NewTabItem("test", widget.NewLabel("Test")))

	assert.False(t, renderer.indicator.Hidden)
}
