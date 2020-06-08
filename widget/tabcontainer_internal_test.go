package widget

import (
	"testing"

	"fyne.io/fyne/internal/cache"

	"github.com/stretchr/testify/assert"
)

func Test_tabButtonRenderer_SetText(t *testing.T) {
	item := &TabItem{Text: "Test", Content: NewLabel("Content")}
	tabs := NewTabContainer(item)
	tabRenderer := cache.Renderer(tabs).(*tabContainerRenderer)
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
