package fyne

import "strings"

// Locale represents a user's locale (language, region and script)
//
// Since: 2.5
type Locale string

// LanguageString returns a version of the local without the script portion.
// For example "en" or "fr-FR".
func (l Locale) LanguageString() string {
	count := strings.Count(string(l), "-")
	if count < 2 {
		return string(l)
	}

	pos := strings.LastIndex(string(l), "-")
	return string(l)[:pos]
}

// String returns the complete locale as a standard string.
func (l Locale) String() string {
	return string(l)
}
