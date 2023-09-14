package fyne

// ExpandedValidator expands on the Validatable interface to allow the widget to expose more context to the parent.
//
// Since: 2.5
type ExpandedValidator interface {
	Validatable

	// HasValidator makes it possible for a parent that cares about child validation (e.g. widget.Form)
	// to know if the developer has attached a validator or not.
	HasValidator() bool

	// SetOnFocusChanged is intended for parent widgets or containers to hook into the change of focus.
	// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
	SetOnFocusChanged(func(bool))

	// SetValidationError manually updates the validation status until the next input change.
	SetValidationError(error)
}

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
