package fyne

// Validatable is an interface for specifying if a widget is validatable. Implementation of this interface
// is not sufficient to validate a widget within a Form widget.
//
// Since: 1.4
type Validatable interface {
	Validate() error

	// SetOnValidationChanged is used to set the callback that will be triggered when the validation state changes.
	// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
	SetOnValidationChanged(func(error))
}

// FormValidatable is an interface for specifying if a widget can be validated within a Form widget.
//
// Since: 2.7
type FormValidatable interface {
	Validatable

	// GetValidationError retrieves the widget's validation error. This will be null if validation passes.
	GetValidationError() error

	// GetValidator retrieves the validator function.
	GetValidator() StringValidator

	// HasFocus indicates that the widget has focus.
	HasFocus() bool

	// IsDirty indicates that the widget's contents have changed.
	IsDirty() bool

	// SetValidator sets the validator function. This function should be called from within the Validate function.
	SetValidator(StringValidator)

	// SetOnFocusChanged sets the function that should be called from within the widget's FocusGained and FocusLost
	// methods.
	SetOnFocusChanged(func(bool))

	// SetValidation error sets the widget's validation error.
	SetValidationError(error)
}

// StringValidator is a function signature for validating string inputs.
//
// Since: 1.4
type StringValidator func(string) error
