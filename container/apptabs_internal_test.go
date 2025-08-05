package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestAppTabs_tabButtonRenderer_SetText(t *testing.T) {
	item := &TabItem{Text: "Test", Content: widget.NewLabel("Content")}
	tabs := NewAppTabs(item)
	tabRenderer := cache.Renderer(tabs).(*appTabsRenderer)
	button := tabRenderer.bar.Objects[0].(*fyne.Container).Objects[0].(*tabButton)
	renderer := cache.Renderer(button).(*tabButtonRenderer)

	assert.Equal(t, "Test", renderer.label.Text)

	button.text = "Temp"
	button.Refresh()
	assert.Equal(t, "Temp", renderer.label.Text)

	item.Text = "Replace"
	tabs.Refresh()
	button = tabRenderer.bar.Objects[0].(*fyne.Container).Objects[0].(*tabButton)
	renderer = cache.Renderer(button).(*tabButtonRenderer)
	assert.Equal(t, "Replace", renderer.label.Text)
}

func Test_tabButtonRenderer_DeleteAdd(t *testing.T) {
	item1 := &TabItem{Text: "Test", Content: widget.NewLabel("Content")}
	item2 := &TabItem{Text: "Delete", Content: widget.NewLabel("Delete")}
	tabs := NewAppTabs(item1, item2)
	tabRenderer := cache.Renderer(tabs).(*appTabsRenderer)
	indicator := tabRenderer.indicator

	pos := indicator.Position()
	tabs.SelectTab(item2)
	assert.NotEqual(t, pos, indicator.Position())
	pos = indicator.Position()

	tabs.Remove(item2)
	tabs.Append(item2)
	tabs.SelectTab(item2)
	assert.Equal(t, pos, indicator.Position())
}

func Test_tabButtonRenderer_EmptyDeleteAdd(t *testing.T) {
	item1 := &TabItem{Text: "Test", Content: widget.NewLabel("Content")}
	tabs := NewAppTabs()

	// ensure enough space for buttons to be created.
	tabs.Resize(fyne.NewSize(300, 200))

	tabRenderer := cache.Renderer(tabs).(*appTabsRenderer)
	assert.Equal(t, 0, len(tabRenderer.bar.Objects[0].(*fyne.Container).Objects))

	tabs.Append(item1)
	assert.Equal(t, 1, len(tabRenderer.bar.Objects[0].(*fyne.Container).Objects))

	tabs.Remove(item1)
	assert.Equal(t, 0, len(tabRenderer.bar.Objects[0].(*fyne.Container).Objects))
}

func TestAppTabs_ShowAfterAdd(t *testing.T) {
	tabs := NewAppTabs()
	renderer := cache.Renderer(tabs).(*appTabsRenderer)

	assert.True(t, renderer.indicator.Hidden)

	tabs.Append(NewTabItem("test", widget.NewLabel("Test")))

	assert.False(t, renderer.indicator.Hidden)
}
