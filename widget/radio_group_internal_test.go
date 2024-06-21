package widget

import (
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRadioGroup_MinSize(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)
	min := radio.MinSize()

	assert.True(t, min.Width > theme.InnerPadding())
	assert.True(t, min.Height > theme.InnerPadding())

	radio2 := NewRadioGroup([]string{"Hi", "H"}, nil)
	min2 := radio2.MinSize()

	assert.Equal(t, min.Width, min2.Width)
	assert.Greater(t, min2.Height, min.Height)
}

func TestRadioGroup_MinSize_Horizontal(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)
	min := radio.MinSize()

	radio2 := NewRadioGroup([]string{"Hi", "He"}, nil)
	radio2.Horizontal = true
	min2 := radio2.MinSize()

	assert.True(t, min2.Width > min.Width)
	assert.Equal(t, min.Height, min2.Height)
}

func TestRadioGroup_Selected(t *testing.T) {
	selected := ""
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radioGroupTestTapItem(t, radio, 0)

	assert.Equal(t, "Hi", selected)
}

func TestRadioGroup_Unselected(t *testing.T) {
	selected := "Hi"
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radioGroupTestTapItem(t, radio, 0)

	assert.Equal(t, "", selected)
}

func TestRadioGroup_DisableWhenSelected(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)
	radio.SetSelected("Hi")
	render := radioGroupTestItemRenderer(t, radio, 0)
	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "primary_"))

	radio.Disable()
	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "disabled_"))
}

func TestRadioGroup_DisableWhenNotSelected(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)
	render := radioGroupTestItemRenderer(t, radio, 0)

	radio.Disable()
	resName := render.over.Resource.Name()
	assert.True(t, strings.HasPrefix(resName, "disabled_"))
}

func TestRadioGroup_SelectedOther(t *testing.T) {
	selected := "Hi"
	radio := NewRadioGroup([]string{"Hi", "Hi2"}, func(sel string) {
		selected = sel
	})
	radioGroupTestTapItem(t, radio, 1)

	assert.Equal(t, "Hi2", selected)
}

func TestRadioGroup_Append(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.TempWidgetRenderer(t, radio).(*radioGroupRenderer).items))

	radio.Options = append(radio.Options, "Another")
	radio.Refresh()

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.TempWidgetRenderer(t, radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_Remove(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.TempWidgetRenderer(t, radio).(*radioGroupRenderer).items))

	radio.Options = radio.Options[:1]
	radio.Refresh()

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.TempWidgetRenderer(t, radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_SetSelected(t *testing.T) {
	changed := false

	radio := NewRadioGroup([]string{"Hi", "Another"}, func(_ string) {
		changed = true
	})

	radio.SetSelected("Another")

	assert.Equal(t, "Another", radio.Selected)
	assert.Equal(t, true, changed)
}

func TestRadioGroup_SetSelectedWithSameOption(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi", "Another"}, nil)

	radio.Selected = "Another"
	radio.Refresh()

	radio.SetSelected("Another")

	assert.Equal(t, "Another", radio.Selected)
}

func TestRadioGroup_SetSelectedEmpty(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi", "Another"}, nil)

	radio.Selected = "Another"
	radio.Refresh()

	radio.SetSelected("")

	assert.Equal(t, "", radio.Selected)
}

func TestRadioGroup_DuplicatedOptions(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi", "Hi", "Hi", "Another", "Another"}, nil)

	assert.Equal(t, 5, len(radio.Options))
	assert.Equal(t, 5, len(test.TempWidgetRenderer(t, radio).(*radioGroupRenderer).items))

	radioGroupTestTapItem(t, radio, 1)
	assert.Equal(t, "Hi", radio.Selected)
	assert.Equal(t, 1, radio.selectedIndex())
	item0 := test.TempWidgetRenderer(t, radio).Objects()[0].(*radioItem)
	assert.Equal(t, false, item0.focused)
	item1 := test.TempWidgetRenderer(t, radio).Objects()[0].(*radioItem)
	assert.Equal(t, false, item1.focused)
}

func TestRadioGroup_AppendDuplicate(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)

	radio.Append("Hi")

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.TempWidgetRenderer(t, radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_Disable(t *testing.T) {
	selected := ""
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radioGroupTestTapItem(t, radio, 0)

	assert.Equal(t, "", selected, "RadioGroup should have been disabled.")
}

func TestRadioGroup_Enable(t *testing.T) {
	selected := ""
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radioGroupTestTapItem(t, radio, 0)
	assert.Equal(t, "", selected, "Radio should have been disabled.")

	radio.Enable()
	radioGroupTestTapItem(t, radio, 0)
	assert.Equal(t, "Hi", selected, "Radio should have been re-enabled.")
}

func TestRadioGroup_Disabled(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, func(string) {})
	assert.False(t, radio.Disabled())
	radio.Disable()
	assert.True(t, radio.Disabled())
	radio.Enable()
	assert.False(t, radio.Disabled())
}

func TestRadioGroup_Hovered(t *testing.T) {

	tests := []struct {
		name         string
		options      []string
		isHorizontal bool
	}{
		{
			name:         "Horizontal",
			options:      []string{"Hi", "Another"},
			isHorizontal: true,
		},
		{
			name:         "Vertical",
			options:      []string{"Hi", "Another"},
			isHorizontal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			radio := NewRadioGroup(tt.options, nil)
			radio.Horizontal = tt.isHorizontal
			item1 := test.TempWidgetRenderer(t, radio).Objects()[0].(*radioItem)
			render1 := radioGroupTestItemRenderer(t, radio, 0)
			render2 := radioGroupTestItemRenderer(t, radio, 1)

			assert.False(t, item1.hovered)
			assert.Equal(t, color.Transparent, render1.focusIndicator.FillColor)
			assert.Equal(t, color.Transparent, render2.focusIndicator.FillColor)

			radio.SetSelected("Hi")
			assert.Equal(t, color.Transparent, render1.focusIndicator.FillColor)
			assert.Equal(t, color.Transparent, render2.focusIndicator.FillColor)

			radio.SetSelected("Another")
			assert.Equal(t, color.Transparent, render1.focusIndicator.FillColor)
			assert.Equal(t, color.Transparent, render2.focusIndicator.FillColor)

			item1.MouseIn(&desktop.MouseEvent{
				PointEvent: fyne.PointEvent{
					Position: fyne.NewPos(theme.Padding(), theme.Padding()),
				},
			})
			assert.True(t, item1.hovered)
			assert.Equal(t, theme.Color(theme.ColorNameHover), render1.focusIndicator.FillColor)
			assert.Equal(t, color.Transparent, render2.focusIndicator.FillColor)

			item1.MouseOut()
			assert.False(t, item1.hovered)
			assert.Equal(t, color.Transparent, render1.focusIndicator.FillColor)
			assert.Equal(t, color.Transparent, render2.focusIndicator.FillColor)
		})
	}
}

func TestRadioGroup_Required(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi", "There"}, func(string) {})
	radio.Required = true
	assert.True(t, radio.Required)
	assert.Equal(t, "", radio.Selected, "the developer should select the default value if “none” is not wanted")

	radio = NewRadioGroup([]string{"Hi", "There"}, func(string) {})
	radio.SetSelected("There")
	radio.Required = true
	assert.True(t, radio.Required)
	assert.Equal(t, "There", radio.Selected, "radio becoming required does not affect a valid selection")

	radio.SetSelected("")
	assert.True(t, radio.Required)
	assert.Equal(t, "", radio.Selected, "the developer should select the default value if “none” is not wanted")

	radio = NewRadioGroup([]string{"Hi", "There"}, func(string) {})
	radio.Required = true
	radio.Resize(radio.MinSize())
	radio.SetSelected("Hi")
	require.Equal(t, "Hi", radio.Selected)
	radioGroupTestTapItem(t, radio, 0)
	assert.Equal(t, "Hi", radio.Selected, "tapping selected option of required radio does nothing")
	radioGroupTestTapItem(t, radio, 1)
	assert.Equal(t, "There", radio.Selected)
}

func TestRadioGroupRenderer_ApplyTheme(t *testing.T) {
	radio := NewRadioGroup([]string{"Test"}, func(string) {})
	render := radioGroupTestItemRenderer(t, radio, 0)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.Refresh()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func radioGroupTestTapItem(t *testing.T, radio *RadioGroup, item int) {
	t.Helper()
	renderer := test.TempWidgetRenderer(t, radio)
	radioItem := renderer.Objects()[item].(*radioItem)
	radioItem.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
}

func radioGroupTestItemRenderer(t *testing.T, radio *RadioGroup, item int) *radioItemRenderer {
	t.Cleanup(func() { cache.DestroyRenderer(radio) })
	return cache.Renderer(test.TempWidgetRenderer(t, radio).Objects()[item].(fyne.Widget)).(*radioItemRenderer)
}
