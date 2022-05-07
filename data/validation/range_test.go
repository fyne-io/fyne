package validation_test

import (
	"testing"

	"fyne.io/fyne/v2/data/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Range(t *testing.T) {
	s := validation.NewRange(2, 3)

	assert.NoError(t, s(2))
	assert.NoError(t, s(3))

	assert.Error(t, s(1))
	assert.Error(t, s(4))
	assert.Equal(t, s(1), s(4))
}
