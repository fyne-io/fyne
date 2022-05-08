package widget

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Validatable = (*Entry)(nil)

// Validate validates the current text in the widget
func (e *Entry) Validate() error {
	if e.Validator == nil {
		return nil
	}

	err := e.Validator(e.Text)
	e.SetValidationError(err)
	return err
}

// SetOnValidationChanged is intended for parent widgets or containers to hook into the validation.
// The function might be overwritten by a parent that cares about child validation (e.g. widget.Form).
func (e *Entry) SetOnValidationChanged(callback func(error)) {
	e.onValidationChanged = callback
}

// SetValidationError manually updates the validation status until the next input change
func (e *Entry) SetValidationError(err error) {
	if e.Validator == nil {
		return
	}
	if err == nil && e.validationError == nil {
		return
	}

	if !errors.Is(err, e.validationError) {
		e.validationError = err

		if e.onValidationChanged != nil {
			e.onValidationChanged(err)
		}

		e.Refresh()
	}
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
	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

func (r *validationStatusRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *validationStatusRenderer) Refresh() {
	r.entry.propertyLock.RLock()
	defer r.entry.propertyLock.RUnlock()
	if r.entry.disabled {
		r.icon.Hide()
		return
	}

	if r.entry.validationError == nil && r.entry.Text != "" {
		r.icon.Resource = theme.ConfirmIcon()
		r.icon.Show()
	} else if r.entry.validationError != nil && !r.entry.focused && r.entry.dirty {
		r.icon.Resource = theme.NewErrorThemedResource(theme.ErrorIcon())
		r.icon.Show()
	} else {
		r.icon.Hide()
	}

	canvas.Refresh(r.icon)
}
