package widget_test

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestCheckGroup_FocusRendering(t *testing.T) {
	t.Run("gain/lose focus", func(t *testing.T) {
		check := widget.NewCheckGroup([]string{"Option A", "Option B", "Option C"}, nil)
		window := test.NewWindow(check)
		defer window.Close()
		window.Resize(check.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		test.AssertRendersToMarkup(t, "check_group/focus_none_focused_none_selected.xml", canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, "check_group/focus_a_focused_none_selected.xml", canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, "check_group/focus_b_focused_none_selected.xml", canvas)
		canvas.Unfocus()
		test.AssertRendersToMarkup(t, "check_group/focus_none_focused_none_selected.xml", canvas)

		check.SetSelected([]string{"Option B", "Option C"})
		test.AssertRendersToMarkup(t, "check_group/focus_none_focused_b_selected.xml", canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, "check_group/focus_a_focused_b_selected.xml", canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, "check_group/focus_b_focused_b_selected.xml", canvas)
		canvas.Unfocus()
		test.AssertRendersToMarkup(t, "check_group/focus_none_focused_b_selected.xml", canvas)
	})

	t.Run("disable/enable focused", func(t *testing.T) {
		check := &widget.CheckGroup{Options: []string{"Option A", "Option B", "Option C"}}
		window := test.NewWindow(check)
		defer window.Close()
		window.Resize(check.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		canvas.FocusNext()
		check.Disable()
		test.AssertRendersToMarkup(t, "check_group/focus_disabled_none_selected.xml", canvas)
		check.Enable()
		test.AssertRendersToMarkup(t, "check_group/focus_a_focused_none_selected.xml", canvas)
	})

	t.Run("append disabled", func(t *testing.T) {
		check := &widget.CheckGroup{Options: []string{"Option A", "Option B"}}
		window := test.NewWindow(check)
		defer window.Close()
		window.Resize(check.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		check.Disable()
		test.AssertRendersToMarkup(t, "check_group/disabled_none_selected.xml", canvas)
		check.Append("Option C")
		test.AssertRendersToMarkup(t, "check_group/disabled_append_none_selected.xml", canvas)
	})
}

func TestCheckGroup_Layout(t *testing.T) {
	for name, tt := range map[string]struct {
		disabled   bool
		horizontal bool
		options    []string
		selected   []string
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
			selected: []string{"Test"},
		},
		"single_selected_disabled": {
			disabled: true,
			options:  []string{"Test"},
			selected: []string{"Test"},
		},
		"single_selected_horizontal": {
			horizontal: true,
			options:    []string{"Test"},
			selected:   []string{"Test"},
		},
		"single_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Test"},
			selected:   []string{"Test"},
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
			options:    []string{"Foo", "Barley"},
		},
		"multiple_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Foo", "Barley"},
		},
		"multiple_selected": {
			options:  []string{"Foo", "Bar"},
			selected: []string{"Foo", "Bar"},
		},
		"multiple_selected_disabled": {
			disabled: true,
			options:  []string{"Foo", "Bar"},
			selected: []string{"Foo", "Bar"},
		},
		"multiple_selected_horizontal": {
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			selected:   []string{"Foo", "Bar"},
		},
		"multiple_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			selected:   []string{"Foo", "Bar"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			check := &widget.CheckGroup{
				Horizontal: tt.horizontal,
				Options:    tt.options,
				Selected:   tt.selected,
			}
			if tt.disabled {
				check.Disable()
			}

			window := test.NewTempWindow(t, check)
			window.Resize(check.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, "check_group/layout_"+name+".xml", window.Canvas())
		})
	}
}

func TestCheckGroup_ToggleSelectionWithSpaceKey(t *testing.T) {
	check := &widget.CheckGroup{Options: []string{"Option A", "Option B", "Option C"}}
	window := test.NewWindow(check)
	defer window.Close()

	var empty []string
	assert.Equal(t, empty, check.Selected)

	canvas := window.Canvas().(test.WindowlessCanvas)
	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, []string{"Option B"}, check.Selected)

	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, []string{"Option B", "Option A"}, check.Selected)

	test.Type(canvas.Focused(), " ")
	assert.Equal(t, []string{"Option B"}, check.Selected)

	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, []string{"Option B", "Option C"}, check.Selected)

	test.Type(canvas.Focused(), " ")
	canvas.FocusPrevious()
	check.Required = true
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, []string{"Option B"}, check.Selected, "cannot unselect required check")
}

func TestCheckGroup_ManipulateOptions(t *testing.T) {
	check := &widget.CheckGroup{Options: []string{}}
	assert.Empty(t, check.Options)

	check.Append("test1")
	assert.Len(t, check.Options, 1)
	check.SetSelected([]string{"test1"})
	assert.Len(t, check.Selected, 1)

	check.Append("test2")
	assert.Len(t, check.Options, 2)

	removed := check.Remove("nope")
	assert.False(t, removed)
	assert.Len(t, check.Options, 2)

	removed = check.Remove("test1")
	assert.True(t, removed)
	assert.Len(t, check.Options, 1)
	assert.Empty(t, check.Selected)
}

func TestCheckGroup_LayoutColumns(t *testing.T) {
	lotsofoptions := []string{}
	for i := 0; i < 50; i++ {
		lotsofoptions = append(lotsofoptions, fmt.Sprintf("Test %d", i))
	}

	testmap := map[string]struct {
		disabled   bool
		horizontal bool
		options    []string
		selected   []string
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
			selected: []string{"Test"},
		},
		"single_selected_disabled": {
			disabled: true,
			options:  []string{"Test"},
			selected: []string{"Test"},
		},
		"single_selected_horizontal": {
			horizontal: true,
			options:    []string{"Test"},
			selected:   []string{"Test"},
		},
		"single_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Test"},
			selected:   []string{"Test"},
		},
		"multiple": {
			options: lotsofoptions,
		},
		"multiple_disabled": {
			disabled: true,
			options:  lotsofoptions,
		},
		"multiple_horizontal": {
			horizontal: true,
			options:    lotsofoptions,
		},
		"multiple_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    lotsofoptions,
		},
		"multiple_selected": {
			options:  lotsofoptions,
			selected: []string{lotsofoptions[0], lotsofoptions[5]},
		},
		"multiple_selected_disabled": {
			disabled: true,
			options:  lotsofoptions,
			selected: []string{lotsofoptions[0], lotsofoptions[5]},
		},
		"multiple_selected_horizontal": {
			horizontal: true,
			options:    lotsofoptions,
			selected:   []string{lotsofoptions[0], lotsofoptions[5]},
		},
		"multiple_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    lotsofoptions,
			selected:   []string{lotsofoptions[0], lotsofoptions[5]},
		},
	}

	for _, i := range []int{-2, 2, 5} {
		for name, tt := range testmap {
			t.Run(name, func(t *testing.T) {
				check := &widget.CheckGroup{
					Horizontal: tt.horizontal,
					Options:    tt.options,
					Selected:   tt.selected,
				}
				check.SetColumns(i)
				if tt.disabled {
					check.Disable()
				}

				window := test.NewTempWindow(t, check)
				window.Resize(check.MinSize().Max(fyne.NewSize(1500, 1500)))

				test.AssertRendersToMarkup(t, fmt.Sprintf("check_group/layout_columns_%d_%s.xml", i, name), window.Canvas())
			})
		}
	}
}
