package fyne

// Validator is a common interface for validating inputs
type Validator interface {
	Validate(string) error
}
