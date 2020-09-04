package validation_test

import (
	"regexp"
	"testing"

	"fyne.io/fyne/data/validation"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Regexp(t *testing.T) {
	r := validation.Regexp{Reason: "Input is not valid"}

	assert.NoError(t, r.Validate("nothing")) // Nothing to validate with, same as no validator.

	r.Regexp = regexp.MustCompile(`^\d{2}-\w{4}$`)

	assert.NoError(t, r.Validate("10-fyne"))
	assert.EqualError(t, r.Validate("100-fynedesk"), r.Reason)
}
