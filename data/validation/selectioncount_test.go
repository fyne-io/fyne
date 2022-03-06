package validation_test

import (
	"testing"

	"fyne.io/fyne/v2/data/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidator_SelectionCount(t *testing.T) {
	s := validation.NewSelectionCount(2)

	assert.NoError(t, s(2))
	assert.NoError(t, s(3))
	errorStr := "too few selections, expected 2 or more"
	assert.EqualError(t, s(0), errorStr)
	assert.EqualError(t, s(1), errorStr)
}
