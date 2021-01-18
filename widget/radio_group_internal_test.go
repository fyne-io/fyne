package widget

import (
	"fmt"
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

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

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
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "Hi", selected)
}

func TestRadioGroup_Unselected(t *testing.T) {
	selected := "Hi"
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected)
}

func TestRadioGroup_DisableWhenSelected(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)
	radio.SetSelected("Hi")
	render := test.WidgetRenderer(radio.items[0]).(*radioItemRenderer)
	resName := render.icon.Resource.Name()

	assert.Equal(t, resName, theme.RadioButtonCheckedIcon().Name())

	radio.Disable()
	resName = render.icon.Resource.Name()
	assert.Equal(t, resName, fmt.Sprintf("disabled_%v", theme.RadioButtonCheckedIcon().Name()))
}

func TestRadioGroup_DisableWhenNotSelected(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)
	render := test.WidgetRenderer(radio.items[0]).(*radioItemRenderer)
	resName := render.icon.Resource.Name()

	assert.Equal(t, resName, theme.RadioButtonIcon().Name())

	radio.Disable()
	resName = render.icon.Resource.Name()
	assert.Equal(t, resName, fmt.Sprintf("disabled_%v", theme.RadioButtonIcon().Name()))
}

func TestRadioGroup_SelectedOther(t *testing.T) {
	selected := "Hi"
	radio := NewRadioGroup([]string{"Hi", "Hi2"}, func(sel string) {
		selected = sel
	})
	radio.items[1].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), radio.MinSize().Height-theme.Padding())})

	assert.Equal(t, "Hi2", selected)
}

func TestRadioGroup_Append(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.WidgetRenderer(radio).(*radioGroupRenderer).items))

	radio.Options = append(radio.Options, "Another")
	radio.Refresh()

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.WidgetRenderer(radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_Remove(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.WidgetRenderer(radio).(*radioGroupRenderer).items))

	radio.Options = radio.Options[:1]
	radio.Refresh()

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.WidgetRenderer(radio).(*radioGroupRenderer).items))
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

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.WidgetRenderer(radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_AppendDuplicate(t *testing.T) {
	radio := NewRadioGroup([]string{"Hi"}, nil)

	radio.Append("Hi")

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.WidgetRenderer(radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_Disable(t *testing.T) {
	selected := ""
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected, "RadioGroup should have been disabled.")
}

func TestRadioGroup_Enable(t *testing.T) {
	selected := ""
	radio := NewRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "", selected, "Radio should have been disabled.")

	radio.Enable()
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
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
			item1 := radio.items[0]
			render1 := cache.Renderer(item1).(*radioItemRenderer)
			render2 := cache.Renderer(radio.items[1]).(*radioItemRenderer)

			assert.False(t, item1.hovered)
			assert.Equal(t, theme.BackgroundColor(), render1.focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render2.focusIndicator.FillColor)

			radio.SetSelected("Hi")
			assert.Equal(t, theme.BackgroundColor(), render1.focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render2.focusIndicator.FillColor)

			radio.SetSelected("Another")
			assert.Equal(t, theme.BackgroundColor(), render1.focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render2.focusIndicator.FillColor)

			item1.MouseIn(&desktop.MouseEvent{
				PointEvent: fyne.PointEvent{
					Position: fyne.NewPos(theme.Padding(), theme.Padding()),
				},
			})
			assert.True(t, item1.hovered)
			assert.Equal(t, theme.HoverColor(), render1.focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render2.focusIndicator.FillColor)

			item1.MouseOut()
			assert.False(t, item1.hovered)
			assert.Equal(t, theme.BackgroundColor(), render1.focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render2.focusIndicator.FillColor)
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
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "Hi", radio.Selected, "tapping selected option of required radio does nothing")
	radio.items[1].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), radio.Size().Height-theme.Padding())})
	assert.Equal(t, "There", radio.Selected)
}

func TestRadioGroupRenderer_ApplyTheme(t *testing.T) {
	radio := NewRadioGroup([]string{"Test"}, func(string) {})
	render := cache.Renderer(radio.items[0]).(*radioItemRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.Refresh()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
