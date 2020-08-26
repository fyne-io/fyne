package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// SetValidationError manually updates the validation status until the next input change
func (e *Entry) SetValidationError(err error) {
	e.validationError = err

	if e.Validator != nil {
		e.validationStatus.Refresh()
	}
}

var _ fyne.Widget = (*validationStatus)(nil)

type validationStatus struct {
	BaseWidget
	entry *Entry
	icon  *canvas.Image
}

func newValidationStatus(e *Entry) *validationStatus {
	rs := &validationStatus{
		icon:  canvas.NewImageFromResource(theme.ConfirmIcon()),
		entry: e,
	}

	rs.ExtendBaseWidget(rs)
	return rs
}

func (r *validationStatus) CreateRenderer() fyne.WidgetRenderer {
	return &validationStatusRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{r.icon}),
		icon:         r.icon,
		entry:        r.entry,
	}
}

var _ fyne.WidgetRenderer = (*validationStatusRenderer)(nil)

type validationStatusRenderer struct {
	widget.BaseRenderer
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
	if r.entry.validationError == nil {
		r.icon.Resource = theme.ConfirmIcon()
		r.icon.Show()
	} else {
		r.icon.Hide()
	}

	if !r.entry.Focused() && r.entry.Text != "" {
		if r.entry.validationError != nil {
			r.icon.Resource = theme.NewErrorThemedResource(theme.ErrorIcon())
		} else {
			r.icon.Resource = theme.ConfirmIcon()
		}

		r.icon.Show()
	}

	canvas.Refresh(r.icon)
}
