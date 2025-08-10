package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

var (
	_ desktop.Cursorable = (*passwordRevealer)(nil)
	_ fyne.Tappable      = (*passwordRevealer)(nil)
	_ fyne.Widget        = (*passwordRevealer)(nil)
)

type passwordRevealer struct {
	BaseWidget

	icon  *canvas.Image
	entry *Entry
}

func newPasswordRevealer(e *Entry) *passwordRevealer {
	pr := &passwordRevealer{
		icon:  canvas.NewImageFromResource(theme.IconForWidget(theme.IconNameVisibilityOff, e)),
		entry: e,
	}
	pr.ExtendBaseWidget(pr)
	return pr
}

func (r *passwordRevealer) CreateRenderer() fyne.WidgetRenderer {
	return &passwordRevealerRenderer{
		WidgetRenderer: NewSimpleRenderer(r.icon),
		icon:           r.icon,
		entry:          r.entry,
	}
}

func (r *passwordRevealer) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (r *passwordRevealer) Tapped(*fyne.PointEvent) {
	if r.entry.Disabled() {
		return
	}

	r.entry.setFieldsAndRefresh(func() {
		r.entry.Password = !r.entry.Password
	})
	fyne.CurrentApp().Driver().CanvasForObject(r).Focus(r.entry.super().(fyne.Focusable))
}

var _ fyne.WidgetRenderer = (*passwordRevealerRenderer)(nil)

type passwordRevealerRenderer struct {
	fyne.WidgetRenderer
	entry *Entry
	icon  *canvas.Image
}

func (r *passwordRevealerRenderer) Layout(size fyne.Size) {
	iconSize := theme.SizeForWidget(theme.SizeNameInlineIcon, r.entry)
	r.icon.Resize(fyne.NewSquareSize(iconSize))
	r.icon.Move(fyne.NewPos((size.Width-iconSize)/2, (size.Height-iconSize)/2))
}

func (r *passwordRevealerRenderer) MinSize() fyne.Size {
	iconSize := theme.SizeForWidget(theme.SizeNameInlineIcon, r.entry)
	return fyne.NewSquareSize(iconSize + theme.SizeForWidget(theme.SizeNameInnerPadding, r.entry)*2)
}

func (r *passwordRevealerRenderer) Refresh() {
	if !r.entry.Password {
		r.icon.Resource = theme.IconForWidget(theme.IconNameVisibility, r.entry)
	} else {
		r.icon.Resource = theme.IconForWidget(theme.IconNameVisibilityOff, r.entry)
	}

	if r.entry.Disabled() {
		r.icon.Resource = theme.NewDisabledResource(r.icon.Resource)
	}
	r.icon.Refresh()
}
