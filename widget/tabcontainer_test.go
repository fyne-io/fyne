package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	tabs.SelectTab(tab2)
	assert.Equal(t, tab2, tabs.CurrentTab())
}

func TestTabContainer_SelectTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")},
		&TabItem{Text: "Test2", Content: NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())

	tabs.SelectTabIndex(1)
	assert.Equal(t, 1, tabs.CurrentTabIndex())
}

func TestTabContainer_SetTabLocation(t *testing.T) {
	tab1 := NewTabItem("Test1", NewLabel("Test1"))
	tab2 := NewTabItem("Test2", NewLabel("Test2"))
	tab3 := NewTabItem("Test3", NewLabel("Test3"))
	tabs := NewTabContainer(tab1, tab2, tab3)

	r := Renderer(tabs).(*tabContainerRenderer)

	buttons := r.tabBar.Children
	require.Len(t, buttons, 3)
	content := tabs.Items[0].Content

	tabs.SetTabLocation(TabLocationLeft)
	assert.Equal(t, fyne.NewPos(0, 0), r.tabBar.Position())
	assert.False(t, r.tabBar.Horizontal)
	assert.Equal(t, fyne.NewPos(r.tabBar.MinSize().Width+theme.Padding(), 0), content.Position())
	assert.Equal(t, fyne.NewPos(r.tabBar.MinSize().Width, 0), r.line.Position())
	assert.Equal(t, fyne.NewSize(theme.Padding(), tabs.MinSize().Height), r.line.Size())
	y := 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(0, y), button.Position())
		y += button.Size().Height
		y += theme.Padding()
	}

	tabs.SetTabLocation(TabLocationBottom)
	assert.Equal(t, fyne.NewPos(0, content.MinSize().Height+theme.Padding()), r.tabBar.Position())
	assert.True(t, r.tabBar.Horizontal)
	assert.Equal(t, fyne.NewPos(0, 0), content.Position())
	assert.Equal(t, fyne.NewPos(0, content.Size().Height), r.line.Position())
	assert.Equal(t, fyne.NewSize(tabs.MinSize().Width, theme.Padding()), r.line.Size())
	x := 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(x, 0), button.Position())
		x += button.Size().Width
		x += theme.Padding()
	}

	tabs.SetTabLocation(TabLocationRight)
	assert.Equal(t, fyne.NewPos(content.Size().Width+theme.Padding(), 0), r.tabBar.Position())
	assert.False(t, r.tabBar.Horizontal)
	assert.Equal(t, fyne.NewPos(0, 0), content.Position())
	assert.Equal(t, fyne.NewPos(content.Size().Width, 0), r.line.Position())
	assert.Equal(t, fyne.NewSize(theme.Padding(), tabs.MinSize().Height), r.line.Size())
	y = 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(0, y), button.Position())
		y += button.Size().Height
		y += theme.Padding()
	}

	tabs.SetTabLocation(TabLocationTop)
	assert.Equal(t, fyne.NewPos(0, 0), r.tabBar.Position())
	assert.True(t, r.tabBar.Horizontal)
	assert.Equal(t, fyne.NewPos(0, r.tabBar.MinSize().Height+theme.Padding()), content.Position())
	assert.Equal(t, fyne.NewPos(0, r.tabBar.MinSize().Height), r.line.Position())
	assert.Equal(t, fyne.NewSize(tabs.MinSize().Width, theme.Padding()), r.line.Size())
	x = 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(x, 0), button.Position())
		x += button.Size().Width
		x += theme.Padding()
	}
}

func TestTabContainerRenderer_ApplyTheme(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")})
	underline := Renderer(tabs).(*tabContainerRenderer).line
	barColor := underline.FillColor

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	Renderer(tabs).ApplyTheme()
	assert.NotEqual(t, barColor, underline.FillColor)
}

func TestTabContainerRenderer_Layout(t *testing.T) {
	tabs := NewTabContainer(
		NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
		NewTabItemWithIcon("Text2", theme.ConfirmIcon(), canvas.NewCircle(theme.BackgroundColor())),
	)

	r := Renderer(tabs).(*tabContainerRenderer)
	r.Layout(fyne.NewSize(100, 100))

	require.Len(t, r.tabBar.Children, 2)
	assert.Equal(t, theme.CancelIcon(), r.tabBar.Children[0].(*Button).Icon)
	assert.Equal(t, theme.ConfirmIcon(), r.tabBar.Children[1].(*Button).Icon)
}
