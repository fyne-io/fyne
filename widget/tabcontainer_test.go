package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
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

	buttons := r.tabBar.Objects
	require.Len(t, buttons, 3)
	content := tabs.Items[0].Content

	tabs.SetTabLocation(TabLocationLeading)
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

func TestTabContainerRenderer_ApplyTheme(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")})
	underline := Renderer(tabs).(*tabContainerRenderer).line
	barColor := underline.FillColor

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	Renderer(tabs).ApplyTheme()
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
