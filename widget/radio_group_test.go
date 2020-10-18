package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestRadioGroup_FocusRendering(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	t.Run("gain/lose focus", func(t *testing.T) {
		radio := widget.NewRadioGroup([]string{"Option A", "Option B", "Option C"}, nil)
		window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), radio))
		defer window.Close()
		window.Resize(radio.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		test.AssertImageMatches(t, "radio/focus_none_focused_none_selected.png", canvas.Capture())
		canvas.FocusNext()
		test.AssertImageMatches(t, "radio/focus_a_focused_none_selected.png", canvas.Capture())
		canvas.FocusNext()
		test.AssertImageMatches(t, "radio/focus_b_focused_none_selected.png", canvas.Capture())
		canvas.Unfocus()
		test.AssertImageMatches(t, "radio/focus_none_focused_none_selected.png", canvas.Capture())

		radio.SetSelected("Option B")
		test.AssertImageMatches(t, "radio/focus_none_focused_b_selected.png", canvas.Capture())
		canvas.FocusNext()
		test.AssertImageMatches(t, "radio/focus_a_focused_b_selected.png", canvas.Capture())
		canvas.FocusNext()
		test.AssertImageMatches(t, "radio/focus_b_focused_b_selected.png", canvas.Capture())
		canvas.Unfocus()
		test.AssertImageMatches(t, "radio/focus_none_focused_b_selected.png", canvas.Capture())
	})

	t.Run("disable/enable focused", func(t *testing.T) {
		radio := &widget.RadioGroup{Options: []string{"Option A", "Option B", "Option C"}}
		window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), radio))
		defer window.Close()
		window.Resize(radio.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		canvas.FocusNext()
		radio.Disable()
		test.AssertImageMatches(t, "radio/focus_disabled_none_selected.png", canvas.Capture())
		radio.Enable()
		test.AssertImageMatches(t, "radio/focus_a_focused_none_selected.png", canvas.Capture())
	})
}

func TestRadioGroup_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		disabled   bool
		horizontal bool
		options    []string
		selected   string
	}{
		"single": {
			options: []string{"Test"},
		},
		"single_disabled": {
			disabled: true,
			options:  []string{"Test"},
		},
		"single_horizontal": {
			horizontal: true,
			options:    []string{"Test"},
		},
		"single_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Test"},
		},
		"single_selected": {
			options:  []string{"Test"},
			selected: "Test",
		},
		"single_selected_disabled": {
			disabled: true,
			options:  []string{"Test"},
			selected: "Test",
		},
		"single_selected_horizontal": {
			horizontal: true,
			options:    []string{"Test"},
			selected:   "Test",
		},
		"single_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Test"},
			selected:   "Test",
		},
		"multiple": {
			options: []string{"Foo", "Bar"},
		},
		"multiple_disabled": {
			disabled: true,
			options:  []string{"Foo", "Bar"},
		},
		"multiple_horizontal": {
			horizontal: true,
			options:    []string{"Foo", "Bar"},
		},
		"multiple_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Foo", "Bar"},
		},
		"multiple_selected": {
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
		},
		"multiple_selected_disabled": {
			disabled: true,
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
		},
		"multiple_selected_horizontal": {
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			selected:   "Foo",
		},
		"multiple_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			selected:   "Foo",
		},
	} {
		t.Run(name, func(t *testing.T) {
			radio := &widget.RadioGroup{
				Horizontal: tt.horizontal,
				Options:    tt.options,
				Selected:   tt.selected,
			}
			if tt.disabled {
				radio.Disable()
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), radio))
			window.Resize(radio.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertImageMatches(t, "radio/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

func TestRadioGroup_ToggleSelectionWithSpaceKey(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	radio := &widget.RadioGroup{Options: []string{"Option A", "Option B", "Option C"}}
	window := test.NewWindow(radio)
	defer window.Close()

	assert.Equal(t, "", radio.Selected)

	canvas := window.Canvas().(test.WindowlessCanvas)
	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option B", radio.Selected)

	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option A", radio.Selected)

	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "", radio.Selected)

	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option C", radio.Selected)

	radio.Required = true
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option C", radio.Selected, "cannot unselect required radio")
}
