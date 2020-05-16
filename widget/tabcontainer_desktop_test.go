// +build !mobile

package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
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

func TestTabContainer_SetTabLocation(t *testing.T) {
	tab1 := NewTabItem("Test1", NewLabel("Test1"))
	tab2 := NewTabItem("Test2", NewLabel("Test2"))
	tab3 := NewTabItem("Test3", NewLabel("Test3"))
	tabs := NewTabContainer(tab1, tab2, tab3)
	r := test.WidgetRenderer(tabs).(*tabContainerRenderer)

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
			r := test.WidgetRenderer(tabs).(*tabContainerRenderer)
			require.Len(t, r.tabBar.Objects, 1)
			tabs.SetTabLocation(tt.location)
			r.Layout(r.MinSize())

			b := r.tabBar.Objects[0].(*tabButton)
			assert.Equal(t, tt.wantButtonSize, b.Size())
			br := test.WidgetRenderer(b).(*tabButtonRenderer)
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
