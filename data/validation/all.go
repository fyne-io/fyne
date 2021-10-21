package validation

import (
	"fyne.io/fyne/v2"
)

// NewAll creates a validator that requires all of the passed validators to pass.
// In short, it combines multiple validators into one.
//
// Since: 2.2
func NewAll(validators ...fyne.StringValidator) fyne.StringValidator {
	return func(text string) error {
		for _, validator := range validators {
			if err := validator(text); err != nil {
				return err
			}
		}
		return nil
	}
}
