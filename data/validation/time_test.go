package validation_test

import (
	"testing"

	"fyne.io/fyne/v2/data/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	validate := validation.NewTime("2006-01-02")

	require.NoError(t, validate("2021-02-21"))
	require.Error(t, validate("21/02-2021"))

	validate = validation.NewTime("15:04")

	require.NoError(t, validate("12:30"))
	assert.Error(t, validate("25:77"))
}
