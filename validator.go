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

// NewRegexpValidator creates a new validator that uses regular expression parsing
func NewRegexpValidator(regexpstr, reason string) *RegexpValidator {
	r := &RegexpValidator{Reason: reason}
	expression, err := regexp.Compile(regexpstr)
	if err != nil {
		LogError("Regexp did not compile", err)
		return nil
	}

	r.Regexp = expression
	return r
}

// Validate will validate the provided string with a regular expression
// Returns nil if valid, otherwise returns an error with a reason text
func (r RegexpValidator) Validate(text string) error {
	if r.Regexp != nil && !r.Regexp.MatchString(text) {
		return errors.New(r.Reason)
	}

	return nil // Nothing to validate with, same as having no validator.
}
