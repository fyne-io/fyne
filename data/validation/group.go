package validation

import (
	"fyne.io/fyne/v2"
)

// NewGroup creates a validator that groups together multiple validators into one.
// All validators need to pass for the result to be valid.
//
// Since: 2.2
func NewGroup(validators ...fyne.StringValidator) fyne.StringValidator {
	return func(text string) error {
		for _, validator := range validators {
			if err := validator(text); err != nil {
				return err
			}
		}
		return nil
	}
}
