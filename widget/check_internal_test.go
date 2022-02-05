package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestCheckSize(t *testing.T) {
	check := NewCheck("Hi", nil)
	min := check.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
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
	render := test.WidgetRenderer(check).(*checkRenderer)

	assert.Equal(t, "primary-"+theme.CheckButtonCheckedIcon().Name(), render.icon.Resource.Name())

	check.Disable()
	assert.Equal(t, fmt.Sprintf("disabled_%v", theme.CheckButtonCheckedIcon().Name()), render.icon.Resource.Name())
}

func TestCheck_DisabledWhenUnchecked(t *testing.T) {
	check := NewCheck("Hi", nil)
	render := test.WidgetRenderer(check).(*checkRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CheckButtonIcon().Name())

	check.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CheckButtonIcon().Name()))
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
		assert.True(t, checkedStateFromCallback == expectedCheckedState)

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
	render := test.WidgetRenderer(check).(*checkRenderer)

	assert.False(t, check.focused)
	assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)

	check.SetChecked(true)
	assert.False(t, check.focused)
	assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)

	test.Tap(check)
	if fyne.CurrentDevice().IsMobile() {
		assert.False(t, check.focused)
	} else {
		assert.True(t, check.focused)
		assert.Equal(t, theme.FocusColor(), render.focusIndicator.FillColor)
	}

	check.Disable()
	assert.True(t, check.disabled)
	assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)

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
	render := test.WidgetRenderer(check).(*checkRenderer)

	check.SetChecked(true)
	assert.False(t, check.hovered)
	assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)

	check.MouseIn(&desktop.MouseEvent{})
	assert.True(t, check.hovered)
	assert.Equal(t, theme.HoverColor(), render.focusIndicator.FillColor)

	test.Tap(check)
	assert.True(t, check.hovered)
	if !fyne.CurrentDevice().IsMobile() {
		assert.Equal(t, theme.FocusColor(), render.focusIndicator.FillColor)
	}

	check.Disable()
	assert.True(t, check.disabled)
	assert.True(t, check.hovered)
	assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)

	check.Enable()
	assert.True(t, check.hovered)
	if fyne.CurrentDevice().IsMobile() {
		assert.Equal(t, theme.HoverColor(), render.focusIndicator.FillColor)
	} else {
		assert.Equal(t, theme.FocusColor(), render.focusIndicator.FillColor)
	}

	check.MouseOut()
	assert.False(t, check.hovered)
	if fyne.CurrentDevice().IsMobile() {
		assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)
	} else {
		assert.Equal(t, theme.FocusColor(), render.focusIndicator.FillColor)
	}

	check.FocusLost()
	assert.False(t, check.hovered)
	assert.Equal(t, theme.BackgroundColor(), render.focusIndicator.FillColor)
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

func TestCheckRenderer_ApplyTheme(t *testing.T) {
	check := &Check{}
	render := test.WidgetRenderer(check).(*checkRenderer)

	textSize := render.label.TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		render.applyTheme()
		customTextSize = render.label.TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}
