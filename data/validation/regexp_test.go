package validation_test

import (
	"testing"

	"fyne.io/fyne/data/validation"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Regexp(t *testing.T) {
	r := validation.NewRegexp(``, "Input is not valid")

	assert.NoError(t, r("nothing")) // Nothing to validate with, same as no validator.

	r = validation.NewRegexp(`^\d{2}-\w{4}$`, "Input is not valid")

	assert.NoError(t, r("10-fyne"))
	assert.EqualError(t, r("100-fynedesk"), "Input is not valid")
}
