package validation

import (
	"fyne.io/fyne/v2"
	"time"
)

// NewTime creates a new validator that verifies times using time.Parse.
// The reference time for the format: Mon Jan 2 15:04:05 -0700 MST 2006.
// The validator will return nil if valid, otherwise returns an error with a reason text.
//
// Since: 2.1
func NewTime(format string) fyne.StringValidator {
	return func(text string) error {
		if _, err := time.Parse(format, text); err != nil {
			return err
		}

		return nil
	}
}
