package widget_test

import (
	"errors"
	"testing"

	"fyne.io/fyne/data/validation"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

var validator = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "Input is not a valid date")

func TestEntry_Validated(t *testing.T) {
	entry, window := setupImageTest(false)
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Validator = validator
	test.AssertImageMatches(t, "entry/validate_initial.png", c.Capture())

	test.Type(entry, "2020-02")
	assert.Error(t, entry.Validator(entry.Text))
	entry.FocusLost()
	test.AssertImageMatches(t, "entry/validate_invalid.png", c.Capture())

	test.Type(entry, "-12")
	assert.NoError(t, entry.Validator(entry.Text))
	test.AssertImageMatches(t, "entry/validate_valid.png", c.Capture())
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

func TestEntry_SetValidationError(t *testing.T) {
	entry, window := setupImageTest(false)
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

	entry.SetText("2020")
	assert.True(t, modified)

	modified = false
	test.Type(entry, "-01-01")
	assert.True(t, modified)
}
