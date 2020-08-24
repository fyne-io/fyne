package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var erroricon = theme.NewErrorThemedResource(theme.ErrorIcon()) // Avoid having to parse XML on each status refresh

var _ fyne.Widget = (*validationStatus)(nil)

type validationStatus struct {
	BaseWidget
	entry *Entry
	icon  *canvas.Image
}

func newValidationStatus(e *Entry) *validationStatus {
	rs := &validationStatus{
		icon:  canvas.NewImageFromResource(erroricon),
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
	if r.entry.validInput {
		r.icon.Resource = theme.ConfirmIcon()
		r.icon.Show()
	} else {
		r.icon.Hide()
	}

	if !r.entry.Focused() && r.entry.Text != "" {
		if !r.entry.validInput {
			r.icon.Resource = erroricon
		} else {
			r.icon.Resource = theme.ConfirmIcon()
		}

		r.icon.Show()
	}

	canvas.Refresh(r.icon)
}
