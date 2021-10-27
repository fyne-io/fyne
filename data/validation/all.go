package validation

import (
	"fyne.io/fyne/v2"
)

// NewAllString creates a validator that requires all of the passed string validators to pass.
// In short, it combines multiple string validators into one.
//
// Since: 2.2
func NewAllString(validators ...fyne.StringValidator) fyne.StringValidator {
	return func(text string) error {
		for _, validator := range validators {
			if err := validator(text); err != nil {
				return err
			}
		}
		return nil
	}
}
