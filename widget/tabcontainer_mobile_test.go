// +build mobile

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
	// the same as TabLocationBottom for mobile
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

	tabs.SetTabLocation(TabLocationBottom)
	tabs.Resize(r.MinSize())
	assert.Equal(t, fyne.NewPos(0, content.MinSize().Height+theme.Padding()), r.tabBar.Position())
	assert.Equal(t, fyne.NewPos(0, 0), content.Position())
	assert.Equal(t, fyne.NewPos(0, content.Size().Height), r.line.Position())
	assert.Equal(t, fyne.NewSize(tabs.MinSize().Width, theme.Padding()), r.line.Size())
	x = 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(x, 0), button.Position())
		x += button.Size().Width
		x += theme.Padding()
	}

	tabs.SetTabLocation(TabLocationTrailing)
	// the same as TabLocationBottom for mobile
	tabs.Resize(r.MinSize())
	assert.Equal(t, fyne.NewPos(0, content.MinSize().Height+theme.Padding()), r.tabBar.Position())
	assert.Equal(t, fyne.NewPos(0, 0), content.Position())
	assert.Equal(t, fyne.NewPos(0, content.Size().Height), r.line.Position())
	assert.Equal(t, fyne.NewSize(tabs.MinSize().Width, theme.Padding()), r.line.Size())
	x = 0
	for _, button := range buttons {
		assert.Equal(t, fyne.NewPos(x, 0), button.Position())
		x += button.Size().Width
		x += theme.Padding()
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
	iconSize := fyne.NewSize(2*theme.IconInlineSize(), 2*theme.IconInlineSize())
	mixedWidth := fyne.Max(iconSize.Width, textWidth)
	mixedIconOffset := 0
	mixedTextOffset := 0
	if mixedWidth > iconSize.Width {
		mixedIconOffset = (mixedWidth - iconSize.Width) / 2
	} else {
		mixedTextOffset = (mixedWidth - textWidth) / 2
	}

	// Leading and trailing location expectations are the same as bottom for mobile.
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
			wantButtonSize: fyne.NewSize(theme.Padding()*2+mixedWidth, theme.Padding()*3+iconSize.Height+textHeight),
			wantIconPos:    fyne.NewPos(theme.Padding()+mixedIconOffset, theme.Padding()),
			wantIconSize:   iconSize,
			wantTextPos:    fyne.NewPos(theme.Padding()+mixedTextOffset, 2*theme.Padding()+iconSize.Height),
			wantTextSize:   textSize,
		},
		{
			name:           "top: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTop,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+textWidth, theme.Padding()*2+textHeight),
			wantTextPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantTextSize:   textSize,
		},
		{
			name:           "top: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTop,
			wantButtonSize: iconSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   iconSize,
		},
		{
			name:           "bottom: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationBottom,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+mixedWidth, theme.Padding()*3+iconSize.Height+textHeight),
			wantIconPos:    fyne.NewPos(theme.Padding()+mixedIconOffset, theme.Padding()),
			wantIconSize:   iconSize,
			wantTextPos:    fyne.NewPos(theme.Padding()+mixedTextOffset, 2*theme.Padding()+iconSize.Height),
			wantTextSize:   textSize,
		},
		{
			name:           "bottom: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationBottom,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+textWidth, theme.Padding()*2+textHeight),
			wantTextPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantTextSize:   textSize,
		},
		{
			name:           "bottom: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationBottom,
			wantButtonSize: iconSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   iconSize,
		},
		{
			name:           "leading: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationLeading,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+mixedWidth, theme.Padding()*3+iconSize.Height+textHeight),
			wantIconPos:    fyne.NewPos(theme.Padding()+mixedIconOffset, theme.Padding()),
			wantIconSize:   iconSize,
			wantTextPos:    fyne.NewPos(theme.Padding()+mixedTextOffset, 2*theme.Padding()+iconSize.Height),
			wantTextSize:   textSize,
		},
		{
			name:           "leading: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationLeading,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+textWidth, theme.Padding()*2+textHeight),
			wantTextPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantTextSize:   textSize,
		},
		{
			name:           "leading: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationLeading,
			wantButtonSize: iconSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   iconSize,
		},
		{
			name:           "trailing: tab with icon and text",
			item:           NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTrailing,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+mixedWidth, theme.Padding()*3+iconSize.Height+textHeight),
			wantIconPos:    fyne.NewPos(theme.Padding()+mixedIconOffset, theme.Padding()),
			wantIconSize:   iconSize,
			wantTextPos:    fyne.NewPos(theme.Padding()+mixedTextOffset, 2*theme.Padding()+iconSize.Height),
			wantTextSize:   textSize,
		},
		{
			name:           "trailing: tab with text only",
			item:           NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTrailing,
			wantButtonSize: fyne.NewSize(theme.Padding()*2+textWidth, theme.Padding()*2+textHeight),
			wantTextPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantTextSize:   textSize,
		},
		{
			name:           "trailing: tab with icon only",
			item:           NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location:       TabLocationTrailing,
			wantButtonSize: iconSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)),
			wantIconPos:    fyne.NewPos(theme.Padding(), theme.Padding()),
			wantIconSize:   iconSize,
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
