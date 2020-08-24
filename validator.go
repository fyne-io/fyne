package fyne

import (
	"errors"
	"regexp"
)

// Validator is a common interface for validating inputs
type Validator interface {
	Validate(string) error
}

// RegexValidater is a regex implementation of fyne.Validator
type RegexValidator struct {
	Reason string
	Regex  *regexp.Regexp
}

// Validate will validate the provided string with a regular expression
// Returns nil if valid, otherwise returns an error with a set reason.
func (r RegexValidator) Validate(text string) error {
	if r.Regex != nil && !r.Regex.MatchString(text) {
		return errors.New(r.Reason)
	}

	return nil // Nothing to validate with, same as having no validator.
}
