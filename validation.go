package fyne

// Validatable is an interface for specifying if a widget is validatable
type Validatable interface {
	Validate() error

	// SetOnValidationChanged is a medthod for a parent to hook into the child validation
	SetOnValidationChanged(func(error))
}

// StringValidator is a function signature for validating string inputs
type StringValidator func(string) error
