// Package validation provides validation for data inside widgets
package validation

import (
	"errors"
	"regexp"

	"fyne.io/fyne/v2"
)

// NewRegexp creates a new validator that uses regular expression parsing.
// The validator will return nil if valid, otherwise returns an error with a reason text.
//
// Since: 1.4
func NewRegexp(regexpstr, reason string) fyne.StringValidator {
	expression, err := regexp.Compile(regexpstr)
	if err != nil {
		fyne.LogError("Regexp did not compile", err)
		return nil
	}

	return func(text string) error {
		if expression != nil && !expression.MatchString(text) {
			return errors.New(reason)
		}

		return nil // Nothing to validate with, same as having no validator.
	}
}
