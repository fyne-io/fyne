package fyne

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_RegexValidator(t *testing.T) {
	r := RegexpValidator{Reason: "Input is not valid"}

	assert.Equal(t, nil, r.Validate("nothing")) // Nothing to validate with, same as no validator.

	r.Regexp = regexp.MustCompile(`^\d{2}-\w{4}$`)

	assert.NoError(t, r.Validate("10-fyne"))
	assert.EqualError(t, r.Validate("100-fynedesk"), r.Reason)
}
