// Package validation provides validation for data inside widgets
package validation

import (
	"errors"
	"regexp"

	"fyne.io/fyne"
)

// Regexp is a regexp implementation of fyne.Validator
type Regexp struct {
	Reason string
	Regexp *regexp.Regexp
}

// NewRegexp creates a new validator that uses regular expression parsing
func NewRegexp(regexpstr, reason string) *Regexp {
	r := &Regexp{Reason: reason}
	expression, err := regexp.Compile(regexpstr)
	if err != nil {
		fyne.LogError("Regexp did not compile", err)
		return nil
	}

	r.Regexp = expression
	return r
}

// Validate will validate the provided string with a regular expression
// Returns nil if valid, otherwise returns an error with a reason text
func (r Regexp) Validate(text string) error {
	if r.Regexp != nil && !r.Regexp.MatchString(text) {
		return errors.New(r.Reason)
	}

	return nil // Nothing to validate with, same as having no validator.
}
