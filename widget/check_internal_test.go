package widget

import (
	"image/color"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestCheckSize(t *testing.T) {
	check := NewCheck("Hi", nil)
	min := check.MinSize()

	assert.Greater(t, min.Width, theme.InnerPadding())
	assert.Greater(t, min.Height, theme.InnerPadding())
}

func TestCheckChecked(t *testing.T) {
	checked := false
	check := NewCheck("Hi", func(on bool) {
		checked = on
	})
	test.Tap(check)

	assert.True(t, checked)
}

func TestCheckUnChecked(t *testing.T) {
	checked := true
	check := NewCheck("Hi", func(on bool) {
		checked = on
	})
	check.Checked = checked
	test.Tap(check)

	assert.False(t, checked)
}

func TestCheck_DisabledWhenChecked(t *testing.T) {
	check := NewCheck("Hi", nil)
	check.SetChecked(true)
	render := test.TempWidgetRenderer(t, check).(*checkRenderer)

	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "primary_"))

	check.Disable()
	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "disabled_"))
}

func TestCheck_DisabledWhenUnchecked(t *testing.T) {
	check := NewCheck("Hi", nil)
	render := test.TempWidgetRenderer(t, check).(*checkRenderer)
	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "inputBorder_"))

	check.Disable()
	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "disabled_"))
}

func TestCheckIsDisabledByDefault(t *testing.T) {
	checkedStateFromCallback := false
	NewCheck("", func(on bool) {
		checkedStateFromCallback = on
	})

	assert.False(t, checkedStateFromCallback)
}

func TestCheckIsEnabledAfterUpdating(t *testing.T) {
	checkedStateFromCallback := false
	check := NewCheck("", func(on bool) {
		checkedStateFromCallback = on
	})

	check.SetChecked(true)

	assert.True(t, checkedStateFromCallback)
}

func TestCheckStateIsCorrectAfterMultipleUpdates(t *testing.T) {
	checkedStateFromCallback := false
	check := NewCheck("", func(on bool) {
		checkedStateFromCallback = on
	})

	expectedCheckedState := false
	for i := 0; i < 5; i++ {
		check.SetChecked(expectedCheckedState)
		assert.Equal(t, expectedCheckedState, checkedStateFromCallback)

		expectedCheckedState = !expectedCheckedState
	}
}

func TestCheck_Disable(t *testing.T) {
	checked := false
	check := NewCheck("Test", func(on bool) {
		checked = on
	})

	check.Disable()
	check.Checked = checked
	test.Tap(check)
	assert.False(t, checked, "Check box should have been disabled.")
}

func TestCheck_Enable(t *testing.T) {
	checked := false
	check := NewCheck("Test", func(on bool) {
		checked = on
	})

	check.Disable()
	check.Checked = checked
	test.Tap(check)
	assert.False(t, checked, "Check box should have been disabled.")

	check.Enable()
	test.Tap(check)
	assert.True(t, checked, "Check box should have been re-enabled.")
}

func TestCheck_Focused(t *testing.T) {
	check := NewCheck("Test", func(on bool) {})
	w := test.NewWindow(check)
	defer w.Close()
	render := test.TempWidgetRenderer(t, check).(*checkRenderer)

	assert.False(t, check.focused)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)

	check.SetChecked(true)
	assert.False(t, check.focused)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)

	test.Tap(check)
	if fyne.CurrentDevice().IsMobile() {
		assert.False(t, check.focused)
	} else {
		assert.True(t, check.focused)
		assert.Equal(t, theme.Color(theme.ColorNameFocus), render.focusIndicator.FillColor)
	}

	check.Disable()
	assert.True(t, check.Disabled())
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)

	check.Enable()
	if fyne.CurrentDevice().IsMobile() {
		assert.False(t, check.focused)
	} else {
		assert.True(t, check.focused)
	}

	check.Hide()
	assert.False(t, check.focused)

	check.Show()
	assert.False(t, check.focused)
}

func TestCheck_Hovered(t *testing.T) {
	check := NewCheck("Test", func(on bool) {})
	w := test.NewWindow(check)
	defer w.Close()
	render := test.TempWidgetRenderer(t, check).(*checkRenderer)

	check.SetChecked(true)
	assert.False(t, check.hovered)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)

	check.MouseIn(&desktop.MouseEvent{})
	assert.True(t, check.hovered)
	assert.Equal(t, theme.Color(theme.ColorNameHover), render.focusIndicator.FillColor)

	test.Tap(check)
	assert.True(t, check.hovered)
	if !fyne.CurrentDevice().IsMobile() {
		assert.Equal(t, theme.Color(theme.ColorNameFocus), render.focusIndicator.FillColor)
	}

	check.Disable()
	assert.True(t, check.Disabled())
	assert.True(t, check.hovered)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)

	check.Enable()
	assert.True(t, check.hovered)
	if fyne.CurrentDevice().IsMobile() {
		assert.Equal(t, theme.Color(theme.ColorNameHover), render.focusIndicator.FillColor)
	} else {
		assert.Equal(t, theme.Color(theme.ColorNameFocus), render.focusIndicator.FillColor)
	}

	check.MouseOut()
	assert.False(t, check.hovered)
	if fyne.CurrentDevice().IsMobile() {
		assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)
	} else {
		assert.Equal(t, theme.Color(theme.ColorNameFocus), render.focusIndicator.FillColor)
	}

	check.FocusLost()
	assert.False(t, check.hovered)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)
}

func TestCheck_HoveredOutsideActiveArea(t *testing.T) {
	check := NewCheck("Test", func(on bool) {})
	w := test.NewWindow(check)
	defer w.Close()
	render := test.TempWidgetRenderer(t, check).(*checkRenderer)

	check.SetChecked(true)
	assert.False(t, check.hovered)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)

	ms := check.MinSize()
	check.MouseIn(&desktop.MouseEvent{PointEvent: fyne.PointEvent{
		Position: fyne.NewPos(ms.Width+2, 1),
	}})
	assert.False(t, check.hovered)
	assert.Equal(t, color.Transparent, render.focusIndicator.FillColor)
}

func TestCheck_TappedOutsideActiveArea(t *testing.T) {
	check := NewCheck("Test", func(on bool) {})
	w := test.NewWindow(check)
	defer w.Close()

	check.SetChecked(true)

	ms := check.MinSize()
	check.Tapped(&fyne.PointEvent{
		Position: fyne.NewPos(ms.Width+2, 1),
	})
	assert.True(t, check.Checked)
}

func TestCheck_TypedRune(t *testing.T) {
	check := NewCheck("Test", func(on bool) {})
	w := test.NewWindow(check)
	defer w.Close()
	assert.False(t, check.Checked)

	test.Tap(check)
	if !fyne.CurrentDevice().IsMobile() {
		assert.True(t, check.focused)
	}
	assert.True(t, check.Checked)

	test.Type(check, " ")
	assert.False(t, check.Checked)

	check.Disable()
	test.Type(check, " ")
	assert.False(t, check.Checked)
}

func TestCheck_Disabled(t *testing.T) {
	check := NewCheck("Test", func(bool) {})
	assert.False(t, check.Disabled())
	check.Disable()
	assert.True(t, check.Disabled())
	check.Enable()
	assert.False(t, check.Disabled())
}
