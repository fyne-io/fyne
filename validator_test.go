package fyne

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {
	v := &Validator{}

	assert.Equal(t, true, v.Validate(105)) // Nothing to validate with, same as no validator.

	v.Regex = regexp.MustCompile(`^\d{2}-\w{4}$`)
	assert.Equal(t, true, v.Validate("10-fyne"))
	assert.Equal(t, false, v.Validate("100-fynedesk"))

	v.Custom = func(input interface{}) bool {
		if value, ok := input.(int); ok {
			return value < 10
		}

		return false
	}
	assert.Equal(t, true, v.Validate(5))
	assert.Equal(t, false, v.Validate(15))
	assert.Equal(t, false, v.Validate("15"))
}
