package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

type extendedRadio struct {
	Radio
}

func newExtendedRadio(opts []string, f func(string)) *extendedRadio {
	ret := &extendedRadio{}
	ret.Options = opts
	ret.OnChanged = f
	ret.ExtendBaseWidget(ret)

	return ret
}

func TestRadio_Extended_BackgroundStyle(t *testing.T) {
	radio := newExtendedRadio([]string{"Hi"}, nil)
	bg := cache.Renderer(radio).BackgroundColor()

	assert.Equal(t, bg, theme.BackgroundColor())
}

func TestRadio_Extended_Selected(t *testing.T) {
	selected := ""
	radio := newExtendedRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "Hi", selected)
}

func TestRadio_Extended_Unselected(t *testing.T) {
	selected := "Hi"
	radio := newExtendedRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	radio.Selected = selected
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected)
}

func TestRadio_Extended_Append(t *testing.T) {
	radio := newExtendedRadio([]string{"Hi"}, nil)

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(cache.Renderer(radio).(*radioRenderer).items))

	radio.Options = append(radio.Options, "Another")
	radio.Refresh()

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(cache.Renderer(radio).(*radioRenderer).items))
}

func TestRadio_Extended_Remove(t *testing.T) {
	radio := newExtendedRadio([]string{"Hi", "Another"}, nil)

	assert.Equal(t, 2, len(radio.Options))
	assert.Equal(t, 2, len(cache.Renderer(radio).(*radioRenderer).items))

	radio.Options = radio.Options[:1]
	radio.Refresh()

	assert.Equal(t, 1, len(radio.Options))
	assert.Equal(t, 1, len(cache.Renderer(radio).(*radioRenderer).items))
}

func TestRadio_Extended_SetSelected(t *testing.T) {
	radio := newExtendedRadio([]string{"Hi", "Another"}, nil)

	radio.SetSelected("Another")

	assert.Equal(t, "Another", radio.Selected)
}

func TestRadio_Extended_Disable(t *testing.T) {
	selected := ""
	radio := newExtendedRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})

	assert.Equal(t, "", selected, "Radio should have been disabled.")
}

func TestRadio_Extended_Enable(t *testing.T) {
	selected := ""
	radio := newExtendedRadio([]string{"Hi"}, func(sel string) {
		selected = sel
	})

	radio.Disable()
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "", selected, "Radio should have been disabled.")

	radio.Enable()
	radio.Tapped(&fyne.PointEvent{Position: fyne.NewPos(theme.Padding(), theme.Padding())})
	assert.Equal(t, "Hi", selected, "Radio should have been re-enabled.")
}

func TestRadio_Extended_Hovered(t *testing.T) {
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
			radio := newExtendedRadio(tt.options, nil)
			radio.Horizontal = tt.isHorizontal
			render := cache.Renderer(radio).(*radioRenderer)

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
			assert.Equal(t, true, radio.hovered)
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
			assert.Equal(t, theme.BackgroundColor(), render.items[0].focusIndicator.FillColor)
			assert.Equal(t, theme.HoverColor(), render.items[1].focusIndicator.FillColor)
		})
	}
}

func TestRadioRenderer_Extended_ApplyTheme(t *testing.T) {
	radio := newExtendedRadio([]string{"Test"}, func(string) {})
	render := cache.Renderer(radio).(*radioRenderer)

	item := render.items[0]
	textSize := item.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.applyTheme()
		customTextSize = item.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
