package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
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

func TestAppTabs_AccessibilityRoleAndLabel(t *testing.T) {
	tabs := NewAppTabs(&TabItem{Text: "One", Content: widget.NewLabel("One")})

	assert.Equal(t, fyne.AccessibleRoleTabList, tabs.AccessibilityRole())
	assert.Equal(t, "", tabs.AccessibilityLabel())
}

func TestAppTabs_AccessibilityChildren(t *testing.T) {
	one := &TabItem{Text: "One", Content: widget.NewLabel("One")}
	two := &TabItem{Text: "Two", Content: widget.NewLabel("Two")}
	tabs := NewAppTabs(one, two)
	tabs.Resize(fyne.NewSize(300, 200))

	children := tabs.AccessibilityChildren()

	r := cache.Renderer(tabs).(*appTabsRenderer)
	assert.Equal(t, r.Objects(), children)
}

func TestAppTabs_AccessibilityChildren_NoRendererYet(t *testing.T) {
	tabs := NewAppTabs(&TabItem{Text: "One", Content: widget.NewLabel("One")})

	assert.Nil(t, tabs.AccessibilityChildren())
}

func Test_tabButton_AccessibilityMetadata(t *testing.T) {
	one := &TabItem{Text: "One", Content: widget.NewLabel("One")}
	two := &TabItem{Text: "Two", Content: widget.NewLabel("Two")}
	tabs := NewAppTabs(one, two)
	tabs.Resize(fyne.NewSize(300, 200))
	tabs.SelectIndex(1)

	r := cache.Renderer(tabs).(*appTabsRenderer)
	buttons := r.bar.Objects[0].(*fyne.Container).Objects
	first := buttons[0].(*tabButton)
	second := buttons[1].(*tabButton)

	assert.Equal(t, fyne.AccessibleRoleTab, first.AccessibilityRole())
	assert.Equal(t, "One", first.AccessibilityLabel())
	assert.Equal(t,
		[]fyne.AccessibleAction{fyne.AccessibleActionPress, fyne.AccessibleActionSelect},
		first.AccessibilityActions())

	assert.Empty(t, first.AccessibilityStates())
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateSelected}, second.AccessibilityStates())

	one.disable()
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateDisabled}, first.AccessibilityStates())
}

func Test_tabButton_AccessibilityPerformAction_SelectsTab(t *testing.T) {
	one := &TabItem{Text: "One", Content: widget.NewLabel("One")}
	two := &TabItem{Text: "Two", Content: widget.NewLabel("Two")}
	tabs := NewAppTabs(one, two)
	tabs.Resize(fyne.NewSize(300, 200))

	r := cache.Renderer(tabs).(*appTabsRenderer)
	second := r.bar.Objects[0].(*fyne.Container).Objects[1].(*tabButton)

	assert.True(t, second.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 1, tabs.SelectedIndex())

	tabs.SelectIndex(0)
	assert.True(t, second.AccessibilityPerformAction(fyne.AccessibleActionSelect))
	assert.Equal(t, 1, tabs.SelectedIndex())

	assert.False(t, second.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
}

func Test_tabButton_AccessibilityPerformAction_DisabledIsNoop(t *testing.T) {
	one := &TabItem{Text: "One", Content: widget.NewLabel("One")}
	two := &TabItem{Text: "Two", Content: widget.NewLabel("Two")}
	tabs := NewAppTabs(one, two)
	tabs.Resize(fyne.NewSize(300, 200))

	r := cache.Renderer(tabs).(*appTabsRenderer)
	second := r.bar.Objects[0].(*fyne.Container).Objects[1].(*tabButton)
	two.disable()

	assert.False(t, second.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 0, tabs.SelectedIndex())
}

func Test_tabButton_AccessibilityLabel_FallsBackToIcon(t *testing.T) {
	icon := theme.HomeIcon()
	tabs := NewAppTabs(&TabItem{Icon: icon, Content: widget.NewLabel("Home")})
	tabs.Resize(fyne.NewSize(300, 200))

	r := cache.Renderer(tabs).(*appTabsRenderer)
	first := r.bar.Objects[0].(*fyne.Container).Objects[0].(*tabButton)

	assert.Equal(t, icon.Name(), first.AccessibilityLabel())
}
