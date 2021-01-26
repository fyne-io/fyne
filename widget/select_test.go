package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelect(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, 2, len(combo.Options))
	assert.Equal(t, "", combo.Selected)
}

func TestSelect_ChangeTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(220, 220))
	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	test.Tap(combo)

	test.AssertImageMatches(t, "select/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		combo.Resize(combo.MinSize())
		combo.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "select/theme_changed.png", w.Canvas().Capture())
	})
}

func TestSelect_ClearSelected(t *testing.T) {
	const (
		opt1     = "1"
		opt2     = "2"
		optClear = ""
	)
	combo := widget.NewSelect([]string{opt1, opt2}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.SetSelected(opt1)
	assert.Equal(t, opt1, combo.Selected)

	var triggered bool
	var triggeredValue string
	combo.OnChanged = func(s string) {
		triggered = true
		triggeredValue = s
	}
	combo.ClearSelected()
	assert.Equal(t, optClear, combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, optClear, triggeredValue)
}

func TestSelect_Disable(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	sel := widget.NewSelect([]string{"Hi"}, func(string) {})
	w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	c := fyne.CurrentApp().Driver().CanvasForObject(sel)

	sel.Disable()
	test.AssertRendersToMarkup(t, "select/disabled.xml", c)
	test.Tap(sel)
	assert.Nil(t, c.Overlays().Top(), "no pop-up for disabled Select")
	test.AssertRendersToMarkup(t, "select/disabled.xml", c)
}

func TestSelect_Disabled(t *testing.T) {
	sel := widget.NewSelect([]string{"Hi"}, func(string) {})
	assert.False(t, sel.Disabled())
	sel.Disable()
	assert.True(t, sel.Disabled())
	sel.Enable()
	assert.False(t, sel.Disabled())
}

func TestSelect_Enable(t *testing.T) {
	selected := ""
	sel := widget.NewSelect([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	sel.Disable()
	require.True(t, sel.Disabled())

	sel.Enable()
	test.Tap(sel)
	c := fyne.CurrentApp().Driver().CanvasForObject(sel)
	ovl := c.Overlays().Top()
	if assert.NotNil(t, ovl, "pop-up for enabled Select") {
		test.TapCanvas(c, ovl.Position().Add(fyne.NewPos(theme.Padding()*2, theme.Padding()*2)))
		assert.Equal(t, "Hi", selected, "Radio should have been re-enabled.")
	}
}

func TestSelect_FocusRendering(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	t.Run("gain/lose focus", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(200, 150))

		c := w.Canvas()
		test.AssertRendersToMarkup(t, "select/focus_unfocused_none_selected.xml", c)
		c.Focus(sel)
		test.AssertRendersToMarkup(t, "select/focus_focused_none_selected.xml", c)
		c.Unfocus()
		test.AssertRendersToMarkup(t, "select/focus_unfocused_none_selected.xml", c)

		sel.SetSelected("Option B")
		assert.Equal(t, "Option B", sel.Selected)
		test.AssertRendersToMarkup(t, "select/focus_unfocused_b_selected.xml", c)
		c.Focus(sel)
		test.AssertRendersToMarkup(t, "select/focus_focused_b_selected.xml", c)
		c.Unfocus()
		test.AssertRendersToMarkup(t, "select/focus_unfocused_b_selected.xml", c)
	})
	t.Run("disable/enable focused", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(200, 150))

		c := w.Canvas().(test.WindowlessCanvas)
		c.FocusNext()
		test.AssertRendersToMarkup(t, "select/focus_focused_none_selected.xml", c)
		sel.Disable()
		test.AssertRendersToMarkup(t, "select/disabled.xml", c)
		sel.Enable()
		test.AssertRendersToMarkup(t, "select/focus_focused_none_selected.xml", c)
	})
}

func TestSelect_KeyboardControl(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	t.Run("activate pop-up", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(150, 200))
		c := w.Canvas()
		c.Focus(sel)

		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected_popup.xml", c)
		test.TapCanvas(c, fyne.NewPos(0, 0))

		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected_popup.xml", c)
		test.TapCanvas(c, fyne.NewPos(0, 0))

		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected_popup.xml", c)
		test.TapCanvas(c, fyne.NewPos(0, 0))

		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
	})

	t.Run("traverse options without pop-up", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(150, 200))
		c := w.Canvas()
		c.Focus(sel)
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option C", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_c_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option B", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_b_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option A", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_a_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option C", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_c_selected.xml", c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option A", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_a_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option B", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_b_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option C", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_c_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option A", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_a_selected.xml", c)
	})

	t.Run("trying to traverse empty options without pop-up", func(t *testing.T) {
		sel := widget.NewSelect([]string{}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(150, 200))
		c := w.Canvas()
		c.Focus(sel)
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, "select/kbdctrl_none_selected.xml", c)
	})
}

func TestSelect_Move(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	combo := widget.NewSelect([]string{"1", "2"}, nil)
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))

	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "select/move_initial.xml", w.Canvas())

	combo.Tapped(&fyne.PointEvent{})
	test.AssertRendersToMarkup(t, "select/move_tapped.xml", w.Canvas())

	combo.Move(fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, "select/move_moved.xml", w.Canvas())
}

func TestSelect_PlaceHolder(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.PlaceHolder = "changed!"
	assert.Equal(t, "changed!", combo.PlaceHolder)
}

func TestSelect_SelectedIndex(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, -1, combo.SelectedIndex())

	combo.SetSelected("2")
	assert.Equal(t, 1, combo.SelectedIndex())

	combo.Selected = "invalid"
	assert.Equal(t, -1, combo.SelectedIndex())
}

func TestSelect_SetSelected(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var triggered bool
	var triggeredValue string
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
		triggeredValue = s
	})
	w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), combo))
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))

	c := w.Canvas()
	test.AssertRendersToMarkup(t, "select/set_selected_none_selected.xml", c)
	combo.SetSelected("2")
	test.AssertRendersToMarkup(t, "select/set_selected_2nd_selected.xml", c)
	assert.Equal(t, "2", combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, "2", triggeredValue)
}

func TestSelect_SetSelected_Callback(t *testing.T) {
	selected := ""
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		selected = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", selected)
}

func TestSelect_SetSelected_Invalid(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("3")

	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected_InvalidNoCallback(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {
		triggered = true
	})
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_SetSelected_InvalidReplace(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")
	combo.SetSelected("3")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_NoChangeOnEmpty(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(string) { triggered = true })
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_SetSelectedIndex(t *testing.T) {
	var triggered bool
	var triggeredValue string
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
		triggeredValue = s
	})
	combo.SetSelectedIndex(1)

	assert.Equal(t, "2", combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, "2", triggeredValue)
}

func TestSelect_SetSelectedIndex_Invalid(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
	})
	combo.SetSelectedIndex(-1)
	assert.Equal(t, -1, combo.SelectedIndex())
	assert.False(t, triggered)
	combo.SetSelectedIndex(2)
	assert.Equal(t, -1, combo.SelectedIndex())
	assert.False(t, triggered)
}

func TestSelect_Tapped(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())

	test.Tap(combo)
	canvas := fyne.CurrentApp().Driver().CanvasForObject(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	test.AssertRendersToMarkup(t, "select/tapped.xml", w.Canvas())
}

func TestSelect_Tapped_Constrained(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())

	canvas := w.Canvas()
	combo.Move(fyne.NewPos(canvas.Size().Width-10, canvas.Size().Height-10))
	test.Tap(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	test.AssertRendersToMarkup(t, "select/tapped_constrained.xml", w.Canvas())
}

func TestSelect_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		placeholder string
		options     []string
		selected    string
		expanded    bool
	}{
		"empty": {},
		"empty_placeholder": {
			placeholder: "(Pick 1)",
		},
		"empty_expanded": {
			expanded: true,
		},
		"empty_expanded_placeholder": {
			placeholder: "(Pick 1)",
			expanded:    true,
		},
		"single": {
			options: []string{"Test"},
		},
		"single_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
		},
		"single_selected": {
			options:  []string{"Test"},
			selected: "Test",
		},
		"single_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			selected:    "Test",
		},
		"single_expanded": {
			options:  []string{"Test"},
			expanded: true,
		},
		"single_expanded_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			expanded:    true,
		},
		"single_expanded_selected": {
			options:  []string{"Test"},
			selected: "Test",
			expanded: true,
		},
		"single_expanded_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			selected:    "Test",
			expanded:    true,
		},
		"multiple": {
			options: []string{"Foo", "Bar"},
		},
		"multiple_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
		},
		"multiple_selected": {
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
		},
		"multiple_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			selected:    "Foo",
		},
		"multiple_expanded": {
			options:  []string{"Foo", "Bar"},
			expanded: true,
		},
		"multiple_expanded_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			expanded:    true,
		},
		"multiple_expanded_selected": {
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
			expanded: true,
		},
		"multiple_expanded_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			selected:    "Foo",
			expanded:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			combo := &widget.Select{
				PlaceHolder: tt.placeholder,
				Options:     tt.options,
				Selected:    tt.selected,
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), combo))
			if tt.expanded {
				test.Tap(combo)
			}
			window.Resize(combo.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, "select/layout_"+name+".xml", window.Canvas())

			window.Close()
		})
	}
}
