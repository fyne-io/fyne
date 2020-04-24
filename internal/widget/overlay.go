package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
)

var _ fyne.Widget = (*Overlay)(nil)
var _ fyne.Tappable = (*Overlay)(nil)

// Overlay is a transparent widget containing one fyne.CanvasObject and meant to be used as overlay.
type Overlay struct {
	base
	Content fyne.CanvasObject

	canvas        fyne.Canvas
	dismissAction func()
	shown         bool
}

// NewOverlay creates an Overlay.
func NewOverlay(c fyne.CanvasObject, canvas fyne.Canvas, dismissAction func()) *Overlay {
	return &Overlay{canvas: canvas, Content: c, dismissAction: dismissAction}
}

// CreateRenderer satisfies the fyne.Widget interface.
func (o *Overlay) CreateRenderer() fyne.WidgetRenderer {
	return &overlayRenderer{BaseRenderer{[]fyne.CanvasObject{o.Content}}, o}
}

// Hide satisfies the fyne.Widget interface.
func (o *Overlay) Hide() {
	if o.shown {
		o.canvas.Overlays().Remove(o)
		o.shown = false
	}
	o.hide(o)
}

// MinSize satisfies the fyne.Widget interface.
func (o *Overlay) MinSize() fyne.Size {
	return o.minSize(o)
}

// MouseIn satisfies the desktop.Hoverable interface.
func (o *Overlay) MouseIn(*desktop.MouseEvent) {
}

// MouseOut satisfies the desktop.Hoverable interface.
func (o *Overlay) MouseOut() {
}

// MouseMoved satisfies the desktop.Hoverable interface.
func (o *Overlay) MouseMoved(*desktop.MouseEvent) {
}

// Refresh satisfies the fyne.Widget interface.
func (o *Overlay) Refresh() {
	o.refresh(o)
}

// Resize satisfies the fyne.Widget interface.
func (o *Overlay) Resize(size fyne.Size) {
	o.resize(size, o)
}

// Show satisfies the fyne.Widget interface.
func (o *Overlay) Show() {
	if !o.shown {
		o.canvas.Overlays().Add(o)
		o.shown = true
	}
	o.show(o)
}

// Tapped satisfies the fyne.Tappable interface.
func (o *Overlay) Tapped(*fyne.PointEvent) {
	if o.dismissAction != nil {
		o.dismissAction()
	}
}

type overlayRenderer struct {
	BaseRenderer
	o *Overlay
}

var _ fyne.WidgetRenderer = (*overlayRenderer)(nil)

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *overlayRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *overlayRenderer) Layout(fyne.Size) {
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *overlayRenderer) MinSize() fyne.Size {
	return r.o.canvas.Size()
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *overlayRenderer) Refresh() {
}
