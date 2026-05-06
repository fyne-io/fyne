package widget_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestButton_Accessibility(t *testing.T) {
	tapped := 0
	b := widget.NewButton("Click", func() { tapped++ })

	assert.Equal(t, fyne.AccessibleRoleButton, b.AccessibilityRole())
	assert.Equal(t, "Click", b.AccessibilityLabel())
	assert.Equal(t, []fyne.AccessibleAction{fyne.AccessibleActionPress}, b.AccessibilityActions())

	assert.Empty(t, b.AccessibilityStates())
	b.Disable()
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateDisabled}, b.AccessibilityStates())
	assert.False(t, b.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 0, tapped)

	b.Enable()
	assert.True(t, b.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 1, tapped)

	assert.False(t, b.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
}

func TestHyperlink_Accessibility(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	hl := widget.NewHyperlink("Link", u)

	assert.Equal(t, fyne.AccessibleRoleLink, hl.AccessibilityRole())
	assert.Equal(t, "Link", hl.AccessibilityLabel())
	assert.Equal(t, []fyne.AccessibleAction{fyne.AccessibleActionPress}, hl.AccessibilityActions())

	tapped := 0
	hl.OnTapped = func() { tapped++ }
	assert.True(t, hl.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 1, tapped)
}

func TestCheck_Accessibility(t *testing.T) {
	changed := 0
	c := widget.NewCheck("Agree", func(bool) { changed++ })

	assert.Equal(t, fyne.AccessibleRoleCheckbox, c.AccessibilityRole())
	assert.Equal(t, "Agree", c.AccessibilityLabel())
	assert.Empty(t, c.AccessibilityStates())

	assert.True(t, c.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.True(t, c.Checked)
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateChecked}, c.AccessibilityStates())
	assert.Equal(t, 1, changed)

	c.Disable()
	assert.False(t, c.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.True(t, c.Checked) // unchanged after disable
	assert.Contains(t, c.AccessibilityStates(), fyne.AccessibleStateDisabled)
}

func TestEntry_Accessibility(t *testing.T) {
	e := widget.NewEntry()
	e.PlaceHolder = "Email"
	e.Text = "user@example.com"

	assert.Equal(t, fyne.AccessibleRoleTextField, e.AccessibilityRole())
	assert.Equal(t, "Email", e.AccessibilityLabel())
	assert.Equal(t, "user@example.com", e.AccessibilityValue())

	assert.True(t, e.AccessibilitySetValue("new@example.com"))
	assert.Equal(t, "new@example.com", e.Text)

	e.Disable()
	assert.Contains(t, e.AccessibilityStates(), fyne.AccessibleStateDisabled)
	assert.False(t, e.AccessibilitySetValue("ignored"))
	assert.Equal(t, "new@example.com", e.Text)

	e.Enable()
	e.Validator = func(string) error { return nil }
	assert.Contains(t, e.AccessibilityStates(), fyne.AccessibleStateRequired)
}

func TestSelect_Accessibility(t *testing.T) {
	test.NewTempApp(t)

	s := widget.NewSelect([]string{"A", "B"}, nil)
	s.PlaceHolder = "Pick one"

	assert.Equal(t, fyne.AccessibleRoleButton, s.AccessibilityRole())
	assert.Equal(t, "Pick one", s.AccessibilityLabel())
	assert.Equal(t, "", s.AccessibilityValue())
	assert.Equal(t,
		[]fyne.AccessibleAction{fyne.AccessibleActionPress, fyne.AccessibleActionShowMenu},
		s.AccessibilityActions())

	s.SetSelected("A")
	assert.Equal(t, "A", s.AccessibilityValue())

	s.Disable()
	assert.Contains(t, s.AccessibilityStates(), fyne.AccessibleStateDisabled)
	assert.False(t, s.AccessibilityPerformAction(fyne.AccessibleActionShowMenu))
}

func TestSlider_Accessibility(t *testing.T) {
	s := widget.NewSlider(0, 10)
	s.Step = 2
	s.SetValue(4)

	assert.Equal(t, fyne.AccessibleRoleSlider, s.AccessibilityRole())
	assert.Equal(t, "4", s.AccessibilityValue())
	assert.Equal(t,
		[]fyne.AccessibleAction{
			fyne.AccessibleActionIncrement,
			fyne.AccessibleActionDecrement,
			fyne.AccessibleActionSetValue,
		},
		s.AccessibilityActions())

	assert.True(t, s.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
	assert.Equal(t, 6.0, s.Value)
	assert.True(t, s.AccessibilityPerformAction(fyne.AccessibleActionDecrement))
	assert.Equal(t, 4.0, s.Value)

	assert.True(t, s.AccessibilitySetValue("8"))
	assert.Equal(t, 8.0, s.Value)
	assert.False(t, s.AccessibilitySetValue("not-a-number"))

	s.Disable()
	assert.Contains(t, s.AccessibilityStates(), fyne.AccessibleStateDisabled)
	assert.False(t, s.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
}

func TestProgressBar_Accessibility(t *testing.T) {
	p := widget.NewProgressBar()
	p.Min = 0
	p.Max = 200
	p.SetValue(50)

	assert.Equal(t, fyne.AccessibleRoleProgressBar, p.AccessibilityRole())
	assert.Equal(t, "25%", p.AccessibilityValue())

	p.TextFormatter = func() string { return "halfway" }
	p.SetValue(100)
	assert.Equal(t, "halfway", p.AccessibilityValue())
}

func TestProgressBarInfinite_Accessibility(t *testing.T) {
	p := widget.NewProgressBarInfinite()
	defer p.Stop()

	assert.Equal(t, fyne.AccessibleRoleProgressBar, p.AccessibilityRole())
	assert.Equal(t, "indeterminate", p.AccessibilityValue())
}

func TestRadioGroup_Accessibility(t *testing.T) {
	rg := widget.NewRadioGroup([]string{"A", "B"}, nil)

	assert.Equal(t, fyne.AccessibleRoleList, rg.AccessibilityRole())
	assert.Empty(t, rg.AccessibilityChildren()) // not yet rendered

	test.NewTempWindow(t, rg)
	assert.Len(t, rg.AccessibilityChildren(), 2)
}

func TestCheckGroup_Accessibility(t *testing.T) {
	cg := widget.NewCheckGroup([]string{"A", "B"}, nil)

	assert.Equal(t, fyne.AccessibleRoleList, cg.AccessibilityRole())
	assert.Empty(t, cg.AccessibilityChildren())

	test.NewTempWindow(t, cg)
	assert.Len(t, cg.AccessibilityChildren(), 2)
}

func TestForm_Accessibility(t *testing.T) {
	form := widget.NewForm(
		widget.NewFormItem("Name", widget.NewEntry()),
	)
	form.OnSubmit = func() {}

	assert.Equal(t, fyne.AccessibleRoleContainer, form.AccessibilityRole())

	test.NewTempWindow(t, form)
	children := form.AccessibilityChildren()
	assert.NotEmpty(t, children)
}
