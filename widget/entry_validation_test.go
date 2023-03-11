package widget_test

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

var validator = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "Input is not a valid date")

func TestEntry_DisabledHideValidation(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Validator = validator
	entry.SetText("invalid text")
	entry.Disable()

	test.AssertImageMatches(t, "entry/validation_disabled.png", c.Capture())
}

func TestEntry_ValidatedEntry(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	r := validation.NewRegexp(`^\d{4}-\d{2}-\d{2}`, "Input is not a valid date")
	entry.Validator = r
	test.AssertRendersToMarkup(t, "entry/validate_initial.xml", c)

	test.Type(entry, "2020-02")
	assert.Error(t, r(entry.Text))
	entry.FocusLost()
	test.AssertRendersToMarkup(t, "entry/validate_invalid.xml", c)

	test.Type(entry, "-12")
	assert.NoError(t, r(entry.Text))
	test.AssertRendersToMarkup(t, "entry/validate_valid.xml", c)
}

func TestEntry_Validate(t *testing.T) {
	entry := widget.NewEntry()
	entry.Validator = validator

	test.Type(entry, "2020-02")
	assert.Error(t, entry.Validate())
	assert.Equal(t, entry.Validate(), entry.Validator(entry.Text))

	test.Type(entry, "-12")
	assert.NoError(t, entry.Validate())
	assert.Equal(t, entry.Validate(), entry.Validator(entry.Text))

	entry.SetText("incorrect")
	assert.Error(t, entry.Validate())
	assert.Equal(t, entry.Validate(), entry.Validator(entry.Text))
}

func TestEntry_NotEmptyValidator(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	entry := widget.NewEntry()
	entry.Validator = func(s string) error {
		if s == "" {
			return errors.New("should not be empty")
		}
		return nil
	}
	w := test.NewWindow(entry)
	defer w.Close()

	test.AssertRendersToMarkup(t, "entry/validator_not_empty_initial.xml", w.Canvas())

	w.Canvas().Focus(entry)

	test.AssertRendersToMarkup(t, "entry/validator_not_empty_focused.xml", w.Canvas())

	w.Canvas().Focus(nil)

	test.AssertRendersToMarkup(t, "entry/validator_not_empty_unfocused.xml", w.Canvas())
}

func TestEntry_SetValidationError(t *testing.T) {
	entry, window := setupImageTest(t, false)
	test.ApplyTheme(t, theme.LightTheme())
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Validator = validator

	entry.SetText("2020-30-30")
	entry.SetValidationError(errors.New("set invalid"))
	test.AssertImageMatches(t, "entry/validation_set_invalid.png", c.Capture())

	entry.SetText("set valid")
	entry.SetValidationError(nil)
	test.AssertImageMatches(t, "entry/validation_set_valid.png", c.Capture())
}

func TestEntry_SetOnValidationChanged(t *testing.T) {
	entry := widget.NewEntry()
	entry.Validator = validator

	modified := false
	entry.SetOnValidationChanged(func(err error) {
		assert.Equal(t, err, entry.Validator(entry.Text))
		modified = true
	})

	test.Type(entry, "2020")
	assert.True(t, modified)

	modified = false
	test.Type(entry, "-01-01")
	assert.True(t, modified)

	modified = false
	entry.SetOnValidationChanged(nil)
	test.Type(entry, "invalid")
	assert.False(t, modified)
}
