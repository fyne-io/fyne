package fyne

// Validatable is an interface for specifying if a widget is validatable.
//
// Since: 1.4
type Validatable interface {
	Validate() error

	// SetOnValidationChanged is used to set the callback that will be triggered when the validation state changes.
	// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
	SetOnValidationChanged(func(error))
}

// StringValidator is a function signature for validating string inputs.
//
// Since: 1.4
type StringValidator func(string) error

// IntValidator is a function sinature for validating integers.
//
// Since: 2.2
type IntValidator func(int) error
