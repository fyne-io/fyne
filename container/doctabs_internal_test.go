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

func TestDocTabs_tabButtonRenderer_Remove(t *testing.T) {
	items := []*TabItem{
		{Text: "1", Content: widget.NewLabel("Content1")},
		{Text: "2", Content: widget.NewLabel("Content2")},
		{Text: "3", Content: widget.NewLabel("Content3")},
	}
	tabs := NewDocTabs(items...)
	tabs.Resize(fyne.NewSize(160, 160))
	tabRenderer := cache.Renderer(tabs).(*docTabsRenderer)

	tabs.SelectIndex(1)
	pos := tabRenderer.indicator.Position()
	tabs.RemoveIndex(0)
	assert.Equal(t, 0, tabs.SelectedIndex())

	assert.Less(t, tabRenderer.indicator.Position().X, pos.X)
}

func TestDocTabs_AccessibilityRoleAndLabel(t *testing.T) {
	tabs := NewDocTabs(&TabItem{Text: "One", Content: widget.NewLabel("One")})

	assert.Equal(t, fyne.AccessibleRoleTabList, tabs.AccessibilityRole())
	assert.Equal(t, "", tabs.AccessibilityLabel())
}

func TestDocTabs_AccessibilityChildren(t *testing.T) {
	tabs := NewDocTabs(
		&TabItem{Text: "One", Content: widget.NewLabel("One")},
		&TabItem{Text: "Two", Content: widget.NewLabel("Two")},
	)
	tabs.Resize(fyne.NewSize(300, 200))

	r := cache.Renderer(tabs).(*docTabsRenderer)
	assert.Equal(t, r.Objects(), tabs.AccessibilityChildren())
}

func TestDocTabs_AccessibilityChildren_NoRendererYet(t *testing.T) {
	tabs := NewDocTabs(&TabItem{Text: "One", Content: widget.NewLabel("One")})

	assert.Nil(t, tabs.AccessibilityChildren())
}

func TestDocTabs_TabButton_AccessibilitySelectsTab(t *testing.T) {
	tabs := NewDocTabs(
		&TabItem{Text: "One", Content: widget.NewLabel("One")},
		&TabItem{Text: "Two", Content: widget.NewLabel("Two")},
	)
	tabs.Resize(fyne.NewSize(300, 200))

	r := cache.Renderer(tabs).(*docTabsRenderer)
	buttons := r.bar.Objects[0].(*Scroll).Content.(*fyne.Container).Objects
	second := buttons[1].(*tabButton)

	assert.Equal(t, fyne.AccessibleRoleTab, second.AccessibilityRole())
	assert.True(t, second.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 1, tabs.SelectedIndex())
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateSelected}, second.AccessibilityStates())
}
