package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type extendedAppTabs struct {
	AppTabs
}

func newExtendedAppTabs(items ...*TabItem) *extendedAppTabs {
	ret := &extendedAppTabs{}
	ret.ExtendBaseWidget(ret)

	ret.Items = items
	return ret
}

func TestAppTabs_Extended_Tapped(t *testing.T) {
	tabs := newExtendedAppTabs(
		NewTabItem("Test1", widget.NewLabel("Test1")),
		NewTabItem("Test2", widget.NewLabel("Test2")),
	)
	tabs.Resize(fyne.NewSize(150, 150)) // Ensure AppTabs is big enough to show both tab buttons
	r := cache.Renderer(tabs).(*appTabsRenderer)

	tab1 := r.bar.Objects[0].(*fyne.Container).Objects[0].(*tabButton)
	tab2 := r.bar.Objects[0].(*fyne.Container).Objects[1].(*tabButton)
	require.Equal(t, 0, tabs.SelectedIndex())
	require.Equal(t, theme.Color(theme.ColorNamePrimary), cache.Renderer(tab1).(*tabButtonRenderer).label.Color)

	tab2.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 1, tabs.SelectedIndex())
	require.Equal(t, theme.Color(theme.ColorNameForeground), cache.Renderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.Color(theme.ColorNamePrimary), cache.Renderer(tab2).(*tabButtonRenderer).label.Color)
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())

	tab2.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 1, tabs.SelectedIndex())
	require.Equal(t, theme.Color(theme.ColorNameForeground), cache.Renderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.Color(theme.ColorNamePrimary), cache.Renderer(tab2).(*tabButtonRenderer).label.Color)
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())

	tab1.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 0, tabs.SelectedIndex())
	require.Equal(t, theme.Color(theme.ColorNamePrimary), cache.Renderer(tab1).(*tabButtonRenderer).label.Color)
	require.Equal(t, theme.Color(theme.ColorNameForeground), cache.Renderer(tab2).(*tabButtonRenderer).label.Color)
	assert.True(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
}
