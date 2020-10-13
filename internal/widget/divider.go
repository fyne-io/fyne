package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*Divider)(nil)

// Divider is a widget for displaying a divider with themeable color.
type Divider struct {
	Base
}

// NewDivider creates a new divider.
func NewDivider() *Divider {
	return &Divider{}
}

// CreateRenderer returns a new renderer for the divider.
//
// Implements: fyne.Widget
func (d *Divider) CreateRenderer() fyne.WidgetRenderer {
	bar := canvas.NewRectangle(theme.DisabledTextColor())
	objects := []fyne.CanvasObject{bar}
	return &dividerRenderer{
		BaseRenderer: NewBaseRenderer(objects),
		bar:          bar,
		d:            d,
	}
}

// Hide hides the divider.
//
// Implements: fyne.Widget
func (d *Divider) Hide() {
	HideWidget(&d.Base, d)
}

// MinSize returns the minimal size of the divider.
//
// Implements: fyne.Widget
func (d *Divider) MinSize() fyne.Size {
	return MinSizeOf(d)
}

// Move sets the position of the divider relative to its parent.
//
// Implements: fyne.Widget
func (d *Divider) Move(pos fyne.Position) {
	MoveWidget(&d.Base, d, pos)
}

// Refresh triggers a redraw of the divider.
//
// Implements: fyne.Widget
func (d *Divider) Refresh() {
	RefreshWidget(d)
}

// Resize changes the size of the divider.
//
// Implements: fyne.Widget
func (d *Divider) Resize(size fyne.Size) {
	ResizeWidget(&d.Base, d, size)
}

// Show makes the divider visible.
//
// Implements: fyne.Widget
func (d *Divider) Show() {
	ShowWidget(&d.Base, d)
}

var _ fyne.WidgetRenderer = (*dividerRenderer)(nil)

type dividerRenderer struct {
	BaseRenderer
	bar *canvas.Rectangle
	d   *Divider
}

func (r *dividerRenderer) Layout(size fyne.Size) {
	r.bar.Resize(size)
}

func (r *dividerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (r *dividerRenderer) Refresh() {
	r.bar.FillColor = theme.DisabledTextColor()
	canvas.Refresh(r.d)
}
