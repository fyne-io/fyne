package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTabContainer_Empty(t *testing.T) {
	tabs := NewTabContainer()
	assert.Equal(t, 0, len(tabs.Items))
	assert.Equal(t, -1, tabs.CurrentTabIndex())
	assert.Nil(t, tabs.CurrentTab())
	min := tabs.MinSize()
	assert.Equal(t, 0, min.Width)
	assert.Equal(t, theme.Padding(), min.Height)
}

func TestTabContainer_Resize_Empty(t *testing.T) {
	tabs := NewTabContainer()
	tabs.Resize(fyne.NewSize(10, 10))
	size := tabs.Size()
	assert.Equal(t, 10, size.Height)
	assert.Equal(t, 10, size.Width)
}

func TestTabContainer_CurrentTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test", Content: NewLabel("Test")})

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())
}

func TestTabContainer_CurrentTab(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())
}

func TestTabContainer_SelectTab(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())

	var selectedTab *TabItem
	tabs.OnChanged = func(tab *TabItem) {
		selectedTab = tab
	}
	tabs.SelectTab(tab2)
	assert.Equal(t, tab2, tabs.CurrentTab())
	assert.Equal(t, tab2, selectedTab)

	tabs.OnChanged = func(tab *TabItem) {
		assert.Fail(t, "unexpected tab selection")
	}
	tabs.SelectTab(NewTabItem("Test3", NewLabel("Test3")))
	assert.Equal(t, tab2, tabs.CurrentTab())
}

func TestTabContainer_SelectTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")},
		&TabItem{Text: "Test2", Content: NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())

	var selectedTab *TabItem
	tabs.OnChanged = func(tab *TabItem) {
		selectedTab = tab
	}
	tabs.SelectTabIndex(1)
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	assert.Equal(t, tabs.Items[1], selectedTab)
}

func TestTabItem_Content(t *testing.T) {
	item1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	item2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(item1, item2)
	r := test.WidgetRenderer(tabs).(*tabContainerRenderer)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, item1.Content, r.objects[0])

	newContent := NewLabel("Test3")
	item1.Content = newContent
	r.Refresh()

	assert.Equal(t, newContent, r.objects[0])
	assert.True(t, newContent.Visible())

	newUnselected := NewLabel("Test4")
	item2.Content = newUnselected
	r.Refresh()

	assert.Equal(t, newUnselected, r.objects[1])
	assert.False(t, newUnselected.Visible())
}

func Test_tabContainer_Tapped(t *testing.T) {
	tabs := NewTabContainer(
		NewTabItem("Test1", NewLabel("Test1")),
		NewTabItem("Test2", NewLabel("Test2")),
		NewTabItem("Test3", NewLabel("Test3")),
	)
	r := test.WidgetRenderer(tabs).(*tabContainerRenderer)

	tab1 := r.tabBar.Objects[0].(*tabButton)
	tab2 := r.tabBar.Objects[1].(*tabButton)
	tab3 := r.tabBar.Objects[2].(*tabButton)
	require.Equal(t, 0, tabs.CurrentTabIndex())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab1).BackgroundColor())

	tab2.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	require.Equal(t, theme.BackgroundColor(), test.WidgetRenderer(tab1).BackgroundColor())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab2).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), test.WidgetRenderer(tab3).BackgroundColor())
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())
	assert.False(t, tabs.Items[2].Content.Visible())

	tab3.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 2, tabs.CurrentTabIndex())
	require.Equal(t, theme.BackgroundColor(), test.WidgetRenderer(tab1).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), test.WidgetRenderer(tab2).BackgroundColor())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab3).BackgroundColor())
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
	assert.True(t, tabs.Items[2].Content.Visible())

	tab1.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 0, tabs.CurrentTabIndex())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab1).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), test.WidgetRenderer(tab2).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), test.WidgetRenderer(tab3).BackgroundColor())
	assert.True(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
	assert.False(t, tabs.Items[2].Content.Visible())
}

func TestTabContainer_Hidden_AsChild(t *testing.T) {
	c1 := NewLabel("Tab 1 content")
	c2 := NewLabel("Tab 2 content\nTab 2 content\nTab 2 content")
	ti1 := NewTabItem("Tab 1", c1)
	ti2 := NewTabItem("Tab 2", c2)
	tabs := NewTabContainer(ti1, ti2)
	tabs.Refresh()

	assert.True(t, c1.Visible())
	assert.False(t, c2.Visible())

	tabs.SelectTabIndex(1)
	assert.False(t, c1.Visible())
	assert.True(t, c2.Visible())
}

func TestTabContainerRenderer_ApplyTheme(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")})
	underline := test.WidgetRenderer(tabs).(*tabContainerRenderer).line
	barColor := underline.FillColor

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	tabs.Refresh()
	assert.NotEqual(t, barColor, underline.FillColor)
}

func TestTabContainer_DynamicTabs(t *testing.T) {
	item := NewTabItem("Text1", canvas.NewCircle(theme.BackgroundColor()))
	tabs := NewTabContainer(item)
	r := test.WidgetRenderer(tabs).(*tabContainerRenderer)

	appendedItem := NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor()))

	tabs.Append(appendedItem)
	assert.Equal(t, len(tabs.Items), 2)
	assert.Equal(t, len(r.tabBar.Objects), 2)
	assert.Equal(t, tabs.Items[1].Text, appendedItem.Text)

	tabs.RemoveIndex(1)
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, len(r.tabBar.Objects), 1)
	assert.Equal(t, tabs.Items[0].Text, item.Text)

	tabs.Append(appendedItem)
	tabs.Remove(tabs.Items[0])
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, len(r.tabBar.Objects), 1)
	assert.Equal(t, tabs.Items[0].Text, appendedItem.Text)

	appendedItem3 := NewTabItem("Text3", canvas.NewCircle(theme.BackgroundColor()))
	appendedItem4 := NewTabItem("Text4", canvas.NewCircle(theme.BackgroundColor()))
	appendedItem5 := NewTabItem("Text5", canvas.NewCircle(theme.BackgroundColor()))
	tabs.Append(appendedItem3)
	tabs.Append(appendedItem4)
	tabs.Append(appendedItem5)
	assert.Equal(t, 4, len(tabs.Items))
	assert.Equal(t, 4, len(r.tabBar.Objects))
	assert.Equal(t, "Text3", tabs.Items[1].Text)
	assert.Equal(t, "Text4", tabs.Items[2].Text)
	assert.Equal(t, "Text5", tabs.Items[3].Text)
}

func Test_tabButton_Hovered(t *testing.T) {
	b := &tabButton{}
	r := test.WidgetRenderer(b)
	require.Equal(t, theme.BackgroundColor(), r.BackgroundColor())

	b.MouseIn(&desktop.MouseEvent{})
	assert.Equal(t, theme.HoverColor(), r.BackgroundColor())
	b.MouseOut()
	assert.Equal(t, theme.BackgroundColor(), r.BackgroundColor())
}

func Test_tabButtonRenderer_BackgroundColor(t *testing.T) {
	r := test.WidgetRenderer(&tabButton{})
	assert.Equal(t, theme.BackgroundColor(), r.BackgroundColor())
}

func TestTabButtonRenderer_ApplyTheme(t *testing.T) {
	item := &TabItem{Text: "Test", Content: NewLabel("Content")}
	tabs := NewTabContainer(item)
	tabButton := tabs.makeButton(item)
	render := test.WidgetRenderer(tabButton).(*tabButtonRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.Refresh()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func Test_tabButtonRenderer_SetText(t *testing.T) {
	item := &TabItem{Text: "Test", Content: NewLabel("Content")}
	tabs := NewTabContainer(item)
	tabRenderer := test.WidgetRenderer(tabs).(*tabContainerRenderer)
	tabButton := tabRenderer.tabBar.Objects[0].(*tabButton)
	renderer := test.WidgetRenderer(tabButton).(*tabButtonRenderer)

	assert.Equal(t, "Test", renderer.label.Text)

	tabButton.setText("Temp")
	assert.Equal(t, "Temp", renderer.label.Text)

	item.Text = "Replace"
	tabs.Refresh()
	assert.Equal(t, "Replace", renderer.label.Text)
}
