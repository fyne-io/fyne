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
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{icon}),
		icon:         icon,
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
	if r.entry.Text == "" {
		r.icon.Hide()
		canvas.Refresh(r.icon)
		return
	}

	if r.entry.validationError == nil {
		r.icon.Resource = theme.ConfirmIcon()
		r.icon.Show()
	} else if !r.entry.focused {
		r.icon.Resource = theme.NewErrorThemedResource(theme.ErrorIcon())
		r.icon.Show()
	} else {
		r.icon.Hide()
	}

	canvas.Refresh(r.icon)
}
