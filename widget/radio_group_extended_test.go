package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

type extendedRadioGroup struct {
	RadioGroup
}

func newextendedRadioGroup(opts []string, f func(string)) *extendedRadioGroup {
	ret := &extendedRadioGroup{}
	ret.Options = opts
	ret.OnChanged = f
	ret.ExtendBaseWidget(ret)
	ret.update() // Not needed for extending Radio but for the tests to be able to access items without creating a renderer first.

	return ret
}

func TestRadioGroup_Extended_Selected(t *testing.T) {
	selected := ""
	radio := newextendedRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "Hi", selected)
}

func TestRadioGroup_Extended_Unselected(t *testing.T) {
	selected := "Hi"
	radio := newextendedRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected)
}

func TestRadioGroup_Extended_Append(t *testing.T) {
	radio := newextendedRadioGroup([]string{"Hi"}, nil)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(cache.Renderer(radio).(*radioGroupRenderer).items))

	radio.Options = append(radio.Options, "Another")
	radio.Refresh()

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(cache.Renderer(radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_Extended_Remove(t *testing.T) {
	radio := newextendedRadioGroup([]string{"Hi", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(cache.Renderer(radio).(*radioGroupRenderer).items))

	radio.Options = radio.Options[:1]
	radio.Refresh()

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(cache.Renderer(radio).(*radioGroupRenderer).items))
}

func TestRadioGroup_Extended_SetSelected(t *testing.T) {
	radio := newextendedRadioGroup([]string{"Hi", "Another"}, nil)

	radio.SetSelected("Another")

	assert.Equal(t, "Another", radio.Selected)
}

func TestRadioGroup_Extended_Disable(t *testing.T) {
	selected := ""
	radio := newextendedRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected, "RadioGroup should have been disabled.")
}

func TestRadioGroup_Extended_Enable(t *testing.T) {
	selected := ""
	radio := newextendedRadioGroup([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "", selected, "Radio should have been disabled.")

	radio.Enable()
	radio.items[0].Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "Hi", selected, "Radio should have been re-enabled.")
}

func TestRadioGroup_Extended_Hovered(t *testing.T) {
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
			radio := newextendedRadioGroup(tt.options, nil)
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

func TestRadioGroupRenderer_Extended_ApplyTheme(t *testing.T) {
	radio := newextendedRadioGroup([]string{"Test"}, func(string) {})
	render := cache.Renderer(radio.items[0]).(*radioItemRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.Refresh()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
