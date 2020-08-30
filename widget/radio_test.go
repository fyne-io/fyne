package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestRadio_Layout(t *testing.T) {
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
			radio := &widget.Radio{
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
