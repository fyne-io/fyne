package fyne

import "regexp"

// Validator is a handler for validating inputs in widgets
type Validator struct {
	Regex *regexp.Regexp

	// Custom specifies a custom validator that takes any input value and returns a bool for the validation status.
	// If defined, Custom() will be preferred over other validation methods.
	Custom func(value interface{}) bool
}

// Validate runs validation on a given set of input using an available validator
func (v *Validator) Validate(input interface{}) bool {
	if v.Custom != nil {
		return v.Custom(input)
	} else if text, ok := input.(string); v.Regex != nil && ok {
		return v.Regex.MatchString(text)
	}

	return true // Nothing to validate with, the same as having no validator.
}
