package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type extendedTabContainer struct {
	TabContainer
}

func newExtendedTabContainer(items ...*TabItem) *extendedTabContainer {
	ret := &extendedTabContainer{}
	ret.ExtendBaseWidget(ret)

	ret.Items = items
	return ret
}

func TestTabContainer_Extended_Tapped(t *testing.T) {
	tabs := newExtendedTabContainer(
		NewTabItem("Test1", NewLabel("Test1")),
		NewTabItem("Test2", NewLabel("Test2")),
	)
	r := test.WidgetRenderer(tabs).(*tabContainerRenderer)

	tab1 := r.tabBar.Objects[0].(*tabButton)
	tab2 := r.tabBar.Objects[1].(*tabButton)
	require.Equal(t, 0, tabs.CurrentTabIndex())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)

	tab2.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	require.Equal(t, theme.TextColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab2).(*tabButtonRenderer).label.Color)
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())

	tab2.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	require.Equal(t, theme.TextColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab2).(*tabButtonRenderer).label.Color)
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())

	tab1.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 0, tabs.CurrentTabIndex())
	require.Equal(t, theme.PrimaryColor(), test.WidgetRenderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.TextColor(), test.WidgetRenderer(tab2).(*tabButtonRenderer).label.Color)
	assert.True(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
}
