package container

import (
	"testing"

	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func Test_tabButtonRenderer_SetText(t *testing.T) {
	item := &TabItem{Text: "Test", Content: widget.NewLabel("Content")}
	tabs := NewAppTabs(item)
	tabRenderer := cache.Renderer(tabs).(*appTabsRenderer)
	button := tabRenderer.tabBar.Objects[0].(*tabButton)
	renderer := cache.Renderer(button).(*tabButtonRenderer)

	assert.Equal(t, "Test", renderer.label.Text)

	button.setText("Temp")
	assert.Equal(t, "Temp", renderer.label.Text)

	item.Text = "Replace"
	tabs.Refresh()
	button = tabRenderer.tabBar.Objects[0].(*tabButton)
	renderer = cache.Renderer(button).(*tabButtonRenderer)
	assert.Equal(t, "Replace", renderer.label.Text)
}

func Test_tabButtonRenderer_DeleteAdd(t *testing.T) {
	item1 := &TabItem{Text: "Test", Content: widget.NewLabel("Content")}
	item2 := &TabItem{Text: "Delete", Content: widget.NewLabel("Delete")}
	tabs := NewAppTabs(item1, item2)
	tabRenderer := cache.Renderer(tabs).(*appTabsRenderer)
	underline := tabRenderer.underline

	pos := underline.Position()
	tabs.SelectTab(item2)
	assert.NotEqual(t, pos, underline.Position())
	pos = underline.Position()

	tabs.Remove(item2)
	tabs.Append(item2)
	tabs.SelectTab(item2)
	assert.Equal(t, pos, underline.Position())
}
