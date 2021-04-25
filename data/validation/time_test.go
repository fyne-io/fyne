package validation_test

import (
	"testing"

	"fyne.io/fyne/v2/data/validation"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	time := validation.NewTime("2006-01-02")

	assert.NoError(t, time("2021-02-21"))
	assert.Error(t, time("21/02-2021"))

	time = validation.NewTime("15:04")

	assert.NoError(t, time("12:30"))
	assert.Error(t, time("25:77"))
}
