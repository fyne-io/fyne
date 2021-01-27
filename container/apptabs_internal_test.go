package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func Test_tabButtonRenderer_SetText(t *testing.T) {
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
