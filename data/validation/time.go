package validation

import (
	"time"

	"fyne.io/fyne/v2"
)

// NewTime creates a new validator that verifies times using time.Parse.
// The validator will return nil if valid, otherwise returns an error with a reason text.
// The reference time for the format: Mon Jan 2 15:04:05 -0700 MST 2006.
// See time.Parse() for more information about the reference time: https://pkg.go.dev/time#Parse
//
// Since: 2.1
func NewTime(format string) fyne.StringValidator {
	return func(text string) error {
		_, err := time.Parse(format, text)
		return err
	}
}
