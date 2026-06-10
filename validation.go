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

// Requireable is implemented by any widgets that want to support the
// [Required] field of a [FormItem]
//
// Since: 2.8
type Requireable interface {
	HasValue() bool

	// SetOnRequiredChanged is used to set the callback that will be triggered when the required state changes.
	// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
	SetOnRequiredChanged(func(bool))
}

// StringValidator is a function signature for validating string inputs.
//
// Since: 1.4
type StringValidator func(string) error
