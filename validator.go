package fyne

import "regexp"

// Validator is a handler for validating inputs in widgets
type Validator struct {
	Regex *regexp.Regexp
}

// RegexStringValidator simplifies verification of strings with regex.
func (v *Validator) RegexStringValidator(text string) bool {
	if v.Regex == nil {
		return true
	}

	return v.Regex.MatchString(text)
}
