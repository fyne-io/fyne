package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var _ desktop.Cursorable = (*passwordRevealer)(nil)
var _ fyne.Tappable = (*passwordRevealer)(nil)
var _ fyne.Widget = (*passwordRevealer)(nil)

type passwordRevealer struct {
	BaseWidget

	icon  *canvas.Image
	entry *Entry
}

func newPasswordRevealer(e *Entry) *passwordRevealer {
	pr := &passwordRevealer{
		icon:  canvas.NewImageFromResource(theme.VisibilityOffIcon()),
		entry: e,
	}
	pr.ExtendBaseWidget(pr)
	return pr
}

func (r *passwordRevealer) CreateRenderer() fyne.WidgetRenderer {
	return &passwordRevealerRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{r.icon}),
		icon:         r.icon,
		entry:        r.entry,
	}
}

func (r *passwordRevealer) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (r *passwordRevealer) Tapped(*fyne.PointEvent) {
	r.entry.setFieldsAndRefresh(func() {
		r.entry.Password = !r.entry.Password
	})
	fyne.CurrentApp().Driver().CanvasForObject(r).Focus(r.entry)
}

var _ fyne.WidgetRenderer = (*passwordRevealerRenderer)(nil)

type passwordRevealerRenderer struct {
	widget.BaseRenderer
	entry *Entry
	icon  *canvas.Image
}

func (r *passwordRevealerRenderer) Layout(size fyne.Size) {
	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

func (r *passwordRevealerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *passwordRevealerRenderer) Refresh() {
	r.entry.propertyLock.RLock()
	defer r.entry.propertyLock.RUnlock()
	if !r.entry.Password {
		r.icon.Resource = theme.VisibilityIcon()
	} else {
		r.icon.Resource = theme.VisibilityOffIcon()
	}
	canvas.Refresh(r.icon)
}
