package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTabContainer_Empty(t *testing.T) {
	tabs := NewTabContainer()
	assert.Equal(t, 0, len(tabs.Items))
	assert.Equal(t, -1, tabs.CurrentTabIndex())
	assert.Nil(t, tabs.CurrentTab())
	min := tabs.MinSize()
	assert.Equal(t, 4, min.Height)
	assert.Equal(t, 0, min.Width)
}

func TestTabContainer_Resize_Empty(t *testing.T) {
	tabs := NewTabContainer()
	tabs.Resize(fyne.NewSize(10, 10))
	size := tabs.Size()
	assert.Equal(t, 10, size.Height)
	assert.Equal(t, 10, size.Width)
}

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

func TestTabItem_Content(t *testing.T) {
	item1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	item2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(item1, item2)
	r := Renderer(tabs).(*tabContainerRenderer)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, item1.Content, r.objects[0])

	newContent := NewLabel("Test3")
	item1.Content = newContent
	r.Refresh()

	assert.Equal(t, newContent, r.objects[0])
	assert.True(t, newContent.Visible())

	newUnselected := NewLabel("Test4")
	item2.Content = newUnselected
	r.Refresh()

	assert.Equal(t, newUnselected, r.objects[1])
	assert.False(t, newUnselected.Visible())
}

func TestTabContainer_SetTabLocation(t *testing.T) {
	tab1 := NewTabItem("Test1", NewLabel("Test1"))
	tab2 := NewTabItem("Test2", NewLabel("Test2"))
	tab3 := NewTabItem("Test3", NewLabel("Test3"))
	tabs := NewTabContainer(tab1, tab2, tab3)
	r := cache.Renderer(tabs).(*tabContainerRenderer)

	buttons := r.tabBar.Objects
	require.Len(t, buttons, 3)
	content := tabs.Items[0].Content

	tabs.SetTabLocation(TabLocationLeading)
	tabs.Resize(r.MinSize())
	assert.Equal(t, fyne.NewPos(0, 0), r.tabBar.Position())
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
	tabs.Resize(r.MinSize())
	assert.Equal(t, fyne.NewPos(0, content.MinSize().Height+theme.Padding()), r.tabBar.Position())
	assert.Equal(t, fyne.NewPos(0, 0), content.Position())
	assert.Equal(t, fyne.NewPos(0, content.Size().Height), r.line.Position())
	assert.Equal(t, fyne.NewSize(tabs.MinSize().Width, theme.Padding()), r.line.Size())
	x := 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(x, 0), button.Position())
		x += button.Size().Width
		x += theme.Padding()
	}

	tabs.SetTabLocation(TabLocationTrailing)
	tabs.Resize(r.MinSize())
	assert.Equal(t, fyne.NewPos(content.Size().Width+theme.Padding(), 0), r.tabBar.Position())
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
	tabs.Resize(r.MinSize())
	assert.Equal(t, fyne.NewPos(0, 0), r.tabBar.Position())
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

func Test_tabContainer_Tapped(t *testing.T) {
	tabs := NewTabContainer(
		NewTabItem("Test1", NewLabel("Test1")),
		NewTabItem("Test2", NewLabel("Test2")),
		NewTabItem("Test3", NewLabel("Test3")),
	)
	r := Renderer(tabs).(*tabContainerRenderer)

	tab1 := r.tabBar.Objects[0].(*tabButton)
	tab2 := r.tabBar.Objects[1].(*tabButton)
	tab3 := r.tabBar.Objects[2].(*tabButton)
	require.Equal(t, 0, tabs.CurrentTabIndex())
	require.Equal(t, theme.PrimaryColor(), Renderer(tab1).BackgroundColor())

	tab2.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	require.Equal(t, theme.BackgroundColor(), Renderer(tab1).BackgroundColor())
	require.Equal(t, theme.PrimaryColor(), Renderer(tab2).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), Renderer(tab3).BackgroundColor())
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.True(t, tabs.Items[1].Content.Visible())
	assert.False(t, tabs.Items[2].Content.Visible())

	tab3.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 2, tabs.CurrentTabIndex())
	require.Equal(t, theme.BackgroundColor(), Renderer(tab1).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), Renderer(tab2).BackgroundColor())
	require.Equal(t, theme.PrimaryColor(), Renderer(tab3).BackgroundColor())
	assert.False(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
	assert.True(t, tabs.Items[2].Content.Visible())

	tab1.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 0, tabs.CurrentTabIndex())
	require.Equal(t, theme.PrimaryColor(), Renderer(tab1).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), Renderer(tab2).BackgroundColor())
	require.Equal(t, theme.BackgroundColor(), Renderer(tab3).BackgroundColor())
	assert.True(t, tabs.Items[0].Content.Visible())
	assert.False(t, tabs.Items[1].Content.Visible())
	assert.False(t, tabs.Items[2].Content.Visible())
}

func TestTabContainer_Hidden_AsChild(t *testing.T) {
	c1 := NewLabel("Tab 1 content")
	c2 := NewLabel("Tab 2 content\nTab 2 content\nTab 2 content")
	ti1 := NewTabItem("Tab 1", c1)
	ti2 := NewTabItem("Tab 2", c2)
	tabs := NewTabContainer(ti1, ti2)
	Renderer(tabs)

	assert.True(t, c1.Visible())
	assert.False(t, c2.Visible())

	tabs.SelectTabIndex(1)
	assert.False(t, c1.Visible())
	assert.True(t, c2.Visible())
}

func TestTabContainerRenderer_ApplyTheme(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")})
	underline := Renderer(tabs).(*tabContainerRenderer).line
	barColor := underline.FillColor

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	Renderer(tabs).Refresh()
	assert.NotEqual(t, barColor, underline.FillColor)
}

func TestTabContainerRenderer_Layout(t *testing.T) {
	textSize := canvas.NewText("Text0", theme.TextColor()).MinSize()
	textWidth := textSize.Width
	textHeight := textSize.Height
	horizontalContentHeight := fyne.Max(theme.IconInlineSize(), textHeight)
	horizontalIconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	horizontalTextSize := fyne.NewSize(textWidth, horizontalContentHeight)
	verticalIconSize := fyne.NewSize(2*theme.IconInlineSize(), 2*theme.IconInlineSize())
	verticalTextSize := fyne.NewSize(textWidth, textHeight)
	verticalMixedWidth := fyne.Max(verticalIconSize.Width, textWidth)
	verticalMixedIconOffset := 0
	verticalMixedTextOffset := 0
	if verticalMixedWidth > verticalIconSize.Width {
		verticalMixedIconOffset = (verticalMixedWidth - verticalIconSize.Width) / 2
	} else {
		verticalMixedTextOffset = (verticalMixedWidth - textWidth) / 2
	}

	tests := []struct {
		name           string
		item           *TabItem
		location       TabLocation
		wantButtonSize fyne.Size
		wantIconPos    fyne.Position
		wantIconSize   fyne.Size
		wantTextPos    fyne.Position
		wantTextSize   fyne.Size
	}{
		{
			name:           "top: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTop,
			wantButtonSize: fyne.NewSize(theme.Padding()*5+theme.IconInlineSize()+textWidth, theme.Padding()*2+horizontalContentHeight),
			wantIconPos:    fyne.NewPos(2*theme.Padding(), theme.Padding()),
			wantIconSize:   horizontalIconSize,
			wantTextPos:    fyne.NewPos(3*theme.Padding()+theme.IconInlineSize(), theme.Padding()),
			wantTextSize:   horizontalTextSize,
		},
		{
			name:           "top: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTop,
			wantButtonSize: fyne.NewSize(theme.Padding()*4+textWidth, theme.Padding()*2+horizontalContentHeight),
			wantTextPos:    fyne.NewPos(2*theme.Padding(), theme.Padding()),
			wantTextSize:   horizontalTextSize,
		},
		{
			name:           "top: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTop,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+theme.IconInlineSize(), theme.Padding()*2+horizontalContentHeight),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   horizontalIconSize,
		},
		{
			name:           "bottom: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationBottom,
			wantButtonSize: fyne.NewSize(theme.Padding()*5+theme.IconInlineSize()+textWidth, theme.Padding()*2+horizontalContentHeight),
			wantIconPos:    fyne.NewPos(2*theme.Padding(), theme.Padding()),
			wantIconSize:   horizontalIconSize,
			wantTextPos:    fyne.NewPos(3*theme.Padding()+theme.IconInlineSize(), theme.Padding()),
			wantTextSize:   horizontalTextSize,
		},
		{
			name:           "bottom: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationBottom,
			wantButtonSize: fyne.NewSize(theme.Padding()*4+textWidth, theme.Padding()*2+horizontalContentHeight),
			wantTextPos:    fyne.NewPos(2*theme.Padding(), theme.Padding()),
			wantTextSize:   horizontalTextSize,
		},
		{
			name:           "bottom: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationBottom,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+theme.IconInlineSize(), theme.Padding()*2+horizontalContentHeight),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   horizontalIconSize,
		},
		{
			name:           "leading: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationLeading,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+verticalMixedWidth, theme.Padding()*3+verticalIconSize.Height+textHeight),
			wantIconPos:    fyne.NewPos(theme.Padding()+verticalMixedIconOffset, theme.Padding()),
			wantIconSize:   verticalIconSize,
			wantTextPos:    fyne.NewPos(theme.Padding()+verticalMixedTextOffset, 2*theme.Padding()+verticalIconSize.Height),
			wantTextSize:   verticalTextSize,
		},
		{
			name:           "leading: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationLeading,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+textWidth, theme.Padding()*2+textHeight),
			wantTextPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantTextSize:   verticalTextSize,
		},
		{
			name:           "leading: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationLeading,
			wantButtonSize: verticalIconSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   verticalIconSize,
		},
		{
			name:           "trailing: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTrailing,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+verticalMixedWidth, theme.Padding()*3+verticalIconSize.Height+textHeight),
			wantIconPos:    fyne.NewPos(theme.Padding()+verticalMixedIconOffset, theme.Padding()),
			wantIconSize:   verticalIconSize,
			wantTextPos:    fyne.NewPos(theme.Padding()+verticalMixedTextOffset, theme.Padding()*2+verticalIconSize.Height),
			wantTextSize:   verticalTextSize,
		},
		{
			name:           "trailing: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTrailing,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+textWidth, theme.Padding()*2+textHeight),
			wantTextPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantTextSize:   verticalTextSize,
		},
		{
			name:           "trailing: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTrailing,
			wantButtonSize: verticalIconSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   verticalIconSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := NewTabContainer(tt.item)
			r := Renderer(tabs).(*tabContainerRenderer)
			require.Len(t, r.tabBar.Objects, 1)
			tabs.SetTabLocation(tt.location)
			r.Layout(r.MinSize())

			b := r.tabBar.Objects[0].(*tabButton)
			assert.Equal(t, tt.wantButtonSize, b.Size())
			br := Renderer(b).(*tabButtonRenderer)
			if tt.item.Icon != nil {
				assert.Equal(t, tt.item.Icon, br.icon.Resource)
				assert.Equal(t, tt.wantIconSize, br.icon.Size(), "icon size")
				assert.Equal(t, tt.wantIconPos, br.icon.Position(), "icon position")
			} else {
				assert.Nil(t, br.icon)
			}
			assert.Equal(t, tt.item.Text, br.label.Text)
			assert.Equal(t, tt.wantTextSize, br.label.Size(), "label size")
			assert.Equal(t, tt.wantTextPos, br.label.Position(), "label position")
		})
	}
}

func TestTabContainer_DynamicTabs(t *testing.T) {
	item := NewTabItem("Text1", canvas.NewCircle(theme.BackgroundColor()))
	tabs := NewTabContainer(item)
	r := Renderer(tabs).(*tabContainerRenderer)

	appendedItem := NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor()))

	tabs.Append(appendedItem)
	assert.Equal(t, len(tabs.Items), 2)
	assert.Equal(t, len(r.tabBar.Objects), 2)
	assert.Equal(t, tabs.Items[1].Text, appendedItem.Text)

	tabs.RemoveIndex(1)
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, len(r.tabBar.Objects), 1)
	assert.Equal(t, tabs.Items[0].Text, item.Text)

	tabs.Append(appendedItem)
	tabs.Remove(tabs.Items[0])
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, len(r.tabBar.Objects), 1)
	assert.Equal(t, tabs.Items[0].Text, appendedItem.Text)

	appendedItem3 := NewTabItem("Text3", canvas.NewCircle(theme.BackgroundColor()))
	appendedItem4 := NewTabItem("Text4", canvas.NewCircle(theme.BackgroundColor()))
	appendedItem5 := NewTabItem("Text5", canvas.NewCircle(theme.BackgroundColor()))
	tabs.Append(appendedItem3)
	tabs.Append(appendedItem4)
	tabs.Append(appendedItem5)
	assert.Equal(t, 4, len(tabs.Items))
	assert.Equal(t, 4, len(r.tabBar.Objects))
	assert.Equal(t, "Text3", tabs.Items[1].Text)
	assert.Equal(t, "Text4", tabs.Items[2].Text)
	assert.Equal(t, "Text5", tabs.Items[3].Text)
}

func Test_tabButton_Hovered(t *testing.T) {
	b := &tabButton{}
	r := Renderer(b)
	require.Equal(t, theme.BackgroundColor(), r.BackgroundColor())

	b.MouseIn(&desktop.MouseEvent{})
	assert.Equal(t, theme.HoverColor(), r.BackgroundColor())
	b.MouseOut()
	assert.Equal(t, theme.BackgroundColor(), r.BackgroundColor())
}

func Test_tabButtonRenderer_BackgroundColor(t *testing.T) {
	r := Renderer(&tabButton{})
	assert.Equal(t, theme.BackgroundColor(), r.BackgroundColor())
}

func TestTabButtonRenderer_ApplyTheme(t *testing.T) {
	item := &TabItem{Text: "Test", Content: NewLabel("Content")}
	tabs := NewTabContainer(item)
	tabButton := tabs.makeButton(item)
	render := Renderer(tabButton).(*tabButtonRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	withTestTheme(func() {
		render.Refresh()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func Test_tabButtonRenderer_SetText(t *testing.T) {
	item := &TabItem{Text: "Test", Content: NewLabel("Content")}
	tabs := NewTabContainer(item)
	tabRenderer := test.WidgetRenderer(tabs).(*tabContainerRenderer)
	tabButton := tabRenderer.tabBar.Objects[0].(*tabButton)
	renderer := test.WidgetRenderer(tabButton).(*tabButtonRenderer)

	assert.Equal(t, "Test", renderer.label.Text)

	tabButton.setText("Temp")
	assert.Equal(t, "Temp", renderer.label.Text)

	item.Text = "Replace"
	tabs.Refresh()
	assert.Equal(t, "Replace", renderer.label.Text)
}
