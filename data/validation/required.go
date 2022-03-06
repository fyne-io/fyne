package validation

import (
	"errors"

	"fyne.io/fyne/v2"
)

// NewRequiredSelection returns a fyne.SelectionValidator that requires
// that one selection (or more) has been made.
func NewRequiredSelection() fyne.SelectionValidator {
	err := errors.New("nothing was selected")
	
	return func(got int) error {
		if got < 1 {
			return err
		}

		return nil
	}
}


// NewRequiredString returns a fyne.StringValidator that requires
// the given string to not be empty.
func NewRequiredString() fyne.StringValidator {
	err := errors.New("no text has been entered")
	
	return func(text string) error {
		if text == "" {
			return err
		}

		return nil
	}
}
