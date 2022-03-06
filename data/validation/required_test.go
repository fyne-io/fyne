package validation_test

import (
	"testing"

	"fyne.io/fyne/v2/data/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidator_RequiredSelection(t *testing.T) {
	rs := validation.NewRequiredSelection()

	assert.NoError(t, rs(2))
	assert.NoError(t, rs(3))
	errorStr := "nothing was selected"
	assert.EqualError(t, rs(0), errorStr)
}

func TestValidator_RequiredString(t *testing.T) {
	rs := validation.NewRequiredString()

	assert.NoError(t, rs("typed"))
	assert.NoError(t, rs("not empty"))
	errorStr := "no text has been entered"
	assert.EqualError(t, rs(""), errorStr)
}
