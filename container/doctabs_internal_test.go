package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestDocTabs_tabButtonRenderer_SetText(t *testing.T) {
	item := &TabItem{Text: "Test", Content: widget.NewLabel("Content")}
	tabs := NewDocTabs(item)
	tabRenderer := cache.Renderer(tabs).(*docTabsRenderer)
	buttons := tabRenderer.bar.Objects[0].(*Scroll).Content.(*fyne.Container).Objects
	button := buttons[0].(*tabButton)
	renderer := cache.Renderer(button).(*tabButtonRenderer)

	assert.Equal(t, "Test", renderer.label.Text)

	button.text = "Temp"
	button.Refresh()
	assert.Equal(t, "Temp", renderer.label.Text)

	item.Text = "Replace"
	tabs.Refresh()
	button = buttons[0].(*tabButton)
	renderer = cache.Renderer(button).(*tabButtonRenderer)
	assert.Equal(t, "Replace", renderer.label.Text)
}
