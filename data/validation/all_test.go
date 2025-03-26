package validation_test

import (
	"testing"

	"fyne.io/fyne/v2/data/validation"
	"github.com/stretchr/testify/require"
)

func TestNewAllStrings(t *testing.T) {
	validate := validation.NewAllStrings(
		validation.NewTime("2006-01-02"),
		validation.NewRegexp(`1\d{3}-\d{2}-\d{2}`, "Only years before 2000 are allowed"),
	)

	require.NoError(t, validate("1988-02-21"))
	require.Error(t, validate("2005-02-21"))
}
