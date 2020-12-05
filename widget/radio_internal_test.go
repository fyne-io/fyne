package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRadio_MinSize(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	min := radio.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)

	radio2 := NewRadio([]string{"Hi", "H"}, nil)
	min2 := radio2.MinSize()

	assert.Equal(t, min.Width, min2.Width)
	assert.Greater(t, min2.Height, min.Height)
}

func TestRadio_BackgroundStyle(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	bg := test.WidgetRenderer(radio).BackgroundColor()

	assert.Equal(t, bg, theme.BackgroundColor())
}

func TestRadio_MinSize_Horizontal(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	min := radio.MinSize()

	radio2 := NewRadio([]string{"Hi", "He"}, nil)
	radio2.Horizontal = true
	min2 := radio2.MinSize()

	assert.True(t, min2.Width > min.Width)
	assert.Equal(t, min.Height, min2.Height)
}

func TestRadio_Selected(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "Hi", selected)
}

func TestRadio_Unselected(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected)
}

func TestRadio_DisableWhenSelected(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	radio.SetSelected("Hi")
	render := test.WidgetRenderer(radio).(*radioRenderer)
	resName := render.items[0].icon.Resource.Name()

	assert.Equal(t, resName, theme.RadioButtonCheckedIcon().Name())

	radio.Disable()
	resName = render.items[0].icon.Resource.Name()
	assert.Equal(t, resName, fmt.Sprintf("disabled_%v", theme.RadioButtonCheckedIcon().Name()))
}

func TestRadio_DisableWhenNotSelected(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)
	render := test.WidgetRenderer(radio).(*radioRenderer)
	resName := render.items[0].icon.Resource.Name()

	assert.Equal(t, resName, theme.RadioButtonIcon().Name())

	radio.Disable()
	resName = render.items[0].icon.Resource.Name()
	assert.Equal(t, resName, fmt.Sprintf("disabled_%v", theme.RadioButtonIcon().Name()))
}

func TestRadio_SelectedOther(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi", "Hi2"}, func(sel string) {
		selected = sel
	})
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), radio.MinSize().Height-theme.Padding())})

	assert.Equal(t, "Hi2", selected)
}

func TestRadio_SelectedHorizontal(t *testing.T) {
	selected := "Hi"
	radio := NewRadio([]string{"Hi", "Hi2"}, func(sel string) {
		selected = sel
	})
	radio.Horizontal = true
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(radio.MinSize().Width-theme.Padding(), theme.Padding())})

	assert.Equal(t, "Hi2", selected)
}

func TestRadio_SelectedNone(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(0, -2)})
	assert.Equal(t, "", selected)

	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(0, radio.size.Height-2)})
	assert.Equal(t, "", selected)
}

func TestRadio_Append(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.WidgetRenderer(radio).(*radioRenderer).items))

	radio.Options = append(radio.Options, "Another")
	radio.Refresh()

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.WidgetRenderer(radio).(*radioRenderer).items))
}

func TestRadio_Remove(t *testing.T) {
	radio := NewRadio([]string{"Hi", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.WidgetRenderer(radio).(*radioRenderer).items))

	radio.Options = radio.Options[:1]
	radio.Refresh()

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.WidgetRenderer(radio).(*radioRenderer).items))
}

func TestRadio_SetSelected(t *testing.T) {
	changed := false

	radio := NewRadio([]string{"Hi", "Another"}, func(_ string) {
		changed = true
	})

	radio.SetSelected("Another")

	assert.Equal(t, "Another", radio.Selected)
	assert.Equal(t, true, changed)
}

func TestRadio_SetSelectedWithSameOption(t *testing.T) {
	radio := NewRadio([]string{"Hi", "Another"}, nil)

	radio.Selected = "Another"
	radio.Refresh()

	radio.SetSelected("Another")

	assert.Equal(t, "Another", radio.Selected)
}

func TestRadio_SetSelectedEmpty(t *testing.T) {
	radio := NewRadio([]string{"Hi", "Another"}, nil)

	radio.Selected = "Another"
	radio.Refresh()

	radio.SetSelected("")

	assert.Equal(t, "", radio.Selected)
}

func TestRadio_DuplicatedOptions(t *testing.T) {
	radio := NewRadio([]string{"Hi", "Hi", "Hi", "Another", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(test.WidgetRenderer(radio).(*radioRenderer).items))
}

func TestRadio_AppendDuplicate(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, nil)

	radio.Append("Hi")

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(test.WidgetRenderer(radio).(*radioRenderer).items))
}

func TestRadio_Disable(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected, "Radio should have been disabled.")
}

func TestRadio_Enable(t *testing.T) {
	selected := ""
	radio := NewRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "", selected, "Radio should have been disabled.")

	radio.Enable()
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "Hi", selected, "Radio should have been re-enabled.")
}

func TestRadio_Disabled(t *testing.T) {
	radio := NewRadio([]string{"Hi"}, func(string) {})
	assert.False(t, radio.Disabled())
	radio.Disable()
	assert.True(t, radio.Disabled())
	radio.Enable()
	assert.False(t, radio.Disabled())
}

func TestRadio_Hovered(t *testing.T) {

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
			radio := NewRadio(tt.options, nil)
			radio.Horizontal = tt.isHorizontal
			render := test.WidgetRenderer(radio).(*radioRenderer)

			assert.Equal(t, false, radio.hovered)
			assert.Equal(t, theme.BackgroundColor(), render.items[0].focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render.items[1].focusIndicator.FillColor)

			radio.SetSelected("Hi")
			assert.Equal(t, theme.BackgroundColor(), render.items[0].focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render.items[1].focusIndicator.FillColor)

			radio.SetSelected("Another")
			assert.Equal(t, theme.BackgroundColor(), render.items[0].focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render.items[1].focusIndicator.FillColor)

			radio.MouseIn(&desktop.MouseEvent{
				PointEvent: fyne.PointEvent{
					Position: fyne.NewPos(theme.Padding(), theme.Padding()),
				},
			})
			assert.Equal(t, 0, radio.hoveredItemIndex)
			assert.Equal(t, theme.HoverColor(), render.items[0].focusIndicator.FillColor)
			assert.Equal(t, theme.BackgroundColor(), render.items[1].focusIndicator.FillColor)

			var mouseMove fyne.Position
			if tt.isHorizontal {
				mouseMove = fyne.NewPos(radio.MinSize().Width-theme.Padding(), theme.Padding())
			} else {
				mouseMove = fyne.NewPos(theme.Padding(), radio.MinSize().Height-theme.Padding())
			}

			radio.MouseMoved(&desktop.MouseEvent{
				PointEvent: fyne.PointEvent{
					Position: mouseMove,
				},
			})
			assert.Equal(t, 1, radio.hoveredItemIndex)
			assert.Equal(t, true, radio.hovered)
			assert.Equal(t, theme.BackgroundColor(), render.items[0].focusIndicator.FillColor)
			assert.Equal(t, theme.HoverColor(), render.items[1].focusIndicator.FillColor)
		})
	}
}

func TestRadio_FocusIndicator_Centered_Vertically(t *testing.T) {
	focusIndicatorSize := theme.IconInlineSize() + theme.Padding()*2

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
			radio := NewRadio(tt.options, nil)
			radio.Horizontal = tt.isHorizontal
			render := test.WidgetRenderer(radio).(*radioRenderer)
			render.Layout(radio.MinSize())

			heightCenterOffset := (radio.itemHeight() - focusIndicatorSize) / 2

			for i, item := range render.items {
				x, y := 0, heightCenterOffset

				if tt.isHorizontal {
					x = i * radio.itemWidth()
				} else {
					y = i*radio.itemHeight() + heightCenterOffset
				}

				assert.Equal(t, fyne.NewPos(x, y), item.focusIndicator.Position1)
			}
		})
	}
}

func TestRadio_Required(t *testing.T) {
	radio := NewRadio([]string{"Hi", "There"}, func(string) {})
	radio.Required = true
	assert.True(t, radio.Required)
	assert.Equal(t, "", radio.Selected, "the developer should select the default value if “none” is not wanted")

	radio = NewRadio([]string{"Hi", "There"}, func(string) {})
	radio.SetSelected("There")
	radio.Required = true
	assert.True(t, radio.Required)
	assert.Equal(t, "There", radio.Selected, "radio becoming required does not affect a valid selection")

	radio.SetSelected("")
	assert.True(t, radio.Required)
	assert.Equal(t, "", radio.Selected, "the developer should select the default value if “none” is not wanted")

	radio = NewRadio([]string{"Hi", "There"}, func(string) {})
	radio.Required = true
	radio.Resize(radio.MinSize())
	radio.SetSelected("Hi")
	require.Equal(t, "Hi", radio.Selected)
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "Hi", radio.Selected, "tapping selected option of required radio does nothing")
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), radio.Size().Height-theme.Padding())})
	assert.Equal(t, "There", radio.Selected)
}

func TestRadioRenderer_ApplyTheme(t *testing.T) {
	radio := NewRadio([]string{"Test"}, func(string) {})
	render := test.WidgetRenderer(radio).(*radioRenderer)

	item := render.items[0]
	textSize := item.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.applyTheme()
		customTextSize = item.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
