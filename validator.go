package fyne

import (
	"errors"
	"regexp"
)

// Validator is a common interface for validating inputs
type Validator interface {
	Validate(string) error
}

// RegexpValidator is a regexp implementation of fyne.Validator
type RegexpValidator struct {
	Reason string
	Regexp *regexp.Regexp
}

// Validate will validate the provided string with a regular expression
// Returns nil if valid, otherwise returns an error with a reason text
func (r RegexpValidator) Validate(text string) error {
	if r.Regexp != nil && !r.Regexp.MatchString(text) {
		return errors.New(r.Reason)
	}

	return nil // Nothing to validate with, same as having no validator.
}
