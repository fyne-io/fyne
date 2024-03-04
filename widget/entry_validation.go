package widget

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Validatable = (*Entry)(nil)

// Validate validates the current text in the widget.
func (e *Entry) Validate() error {
	if e.Validator == nil {
		return nil
	}

	err := e.Validator(e.Text)
	e.SetValidationError(err)
	return err
}

// validate works like Validate but only updates the internal state and does not refresh.
func (e *Entry) validate() {
	if e.Validator == nil {
		return
	}

	err := e.Validator(e.Text)
	e.setValidationError(err)
}

// SetOnValidationChanged is intended for parent widgets or containers to hook into the validation.
// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
func (e *Entry) SetOnValidationChanged(callback func(error)) {
	e.onValidationChanged = callback
}

// SetValidationError manually updates the validation status until the next input change.
func (e *Entry) SetValidationError(err error) {
	if e.Validator == nil {
		return
	}

	if !e.setValidationError(err) {
		return
	}

	e.Refresh()
}

// setValidationError sets the validation error and returns a bool to indicate if it changes.
// It assumes that the widget has a validator.
func (e *Entry) setValidationError(err error) bool {
	if err == nil && e.validationError == nil {
		return false
	}
	if errors.Is(err, e.validationError) {
		return false
	}

	e.validationError = err

	if e.onValidationChanged != nil {
		e.onValidationChanged(err)
	}

	return true
}

var _ fyne.Widget = (*validationStatus)(nil)

type validationStatus struct {
	BaseWidget
	entry *Entry
}

func newValidationStatus(e *Entry) *validationStatus {
	rs := &validationStatus{
		entry: e,
	}

	rs.ExtendBaseWidget(rs)
	return rs
}

func (r *validationStatus) CreateRenderer() fyne.WidgetRenderer {
	icon := &canvas.Image{}
	icon.Hide()
	return &validationStatusRenderer{
		WidgetRenderer: NewSimpleRenderer(icon),
		icon:           icon,
		entry:          r.entry,
	}
}

var _ fyne.WidgetRenderer = (*validationStatusRenderer)(nil)

type validationStatusRenderer struct {
	fyne.WidgetRenderer
	entry *Entry
	icon  *canvas.Image
}

func (r *validationStatusRenderer) Layout(size fyne.Size) {
	iconSize := r.entry.Theme().Size(theme.SizeNameInlineIcon)
	r.icon.Resize(fyne.NewSquareSize(iconSize))
	r.icon.Move(fyne.NewPos((size.Width-iconSize)/2, (size.Height-iconSize)/2))
}

func (r *validationStatusRenderer) MinSize() fyne.Size {
	iconSize := r.entry.Theme().Size(theme.SizeNameInlineIcon)
	return fyne.NewSquareSize(iconSize)
}

func (r *validationStatusRenderer) Refresh() {
	th := r.entry.Theme()
	r.entry.propertyLock.RLock()
	defer r.entry.propertyLock.RUnlock()
	if r.entry.disabled.Load() {
		r.icon.Hide()
		return
	}

	if r.entry.validationError == nil && r.entry.Text != "" {
		r.icon.Resource = th.Icon(theme.IconNameConfirm)
		r.icon.Show()
	} else if r.entry.validationError != nil && !r.entry.focused && r.entry.dirty {
		r.icon.Resource = theme.NewErrorThemedResource(th.Icon(theme.IconNameError))
		r.icon.Show()
	} else {
		r.icon.Hide()
	}

	r.icon.Refresh()
}
