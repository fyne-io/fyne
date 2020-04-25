package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
)

var _ fyne.Widget = (*OverlayContainer)(nil)
var _ fyne.Tappable = (*OverlayContainer)(nil)

// OverlayContainer is a transparent widget containing one fyne.CanvasObject and meant to be used as overlay.
type OverlayContainer struct {
	base
	Content fyne.CanvasObject

	canvas        fyne.Canvas
	dismissAction func()
	shown         bool
}

// NewOverlayContainer creates an OverlayContainer.
func NewOverlayContainer(c fyne.CanvasObject, canvas fyne.Canvas, dismissAction func()) *OverlayContainer {
	return &OverlayContainer{canvas: canvas, Content: c, dismissAction: dismissAction}
}

// CreateRenderer satisfies the fyne.Widget interface.
func (o *OverlayContainer) CreateRenderer() fyne.WidgetRenderer {
	return &overlayRenderer{BaseRenderer{[]fyne.CanvasObject{o.Content}}, o}
}

// Hide satisfies the fyne.Widget interface.
func (o *OverlayContainer) Hide() {
	if o.shown {
		o.canvas.Overlays().Remove(o)
		o.shown = false
	}
	o.hide(o)
}

// MinSize satisfies the fyne.Widget interface.
func (o *OverlayContainer) MinSize() fyne.Size {
	return o.minSize(o)
}

// MouseIn satisfies the desktop.Hoverable interface.
func (o *OverlayContainer) MouseIn(*desktop.MouseEvent) {
}

// MouseOut satisfies the desktop.Hoverable interface.
func (o *OverlayContainer) MouseOut() {
}

// MouseMoved satisfies the desktop.Hoverable interface.
func (o *OverlayContainer) MouseMoved(*desktop.MouseEvent) {
}

// Refresh satisfies the fyne.Widget interface.
func (o *OverlayContainer) Refresh() {
	o.refresh(o)
}

// Resize satisfies the fyne.Widget interface.
func (o *OverlayContainer) Resize(size fyne.Size) {
	o.resize(size, o)
}

// Show satisfies the fyne.Widget interface.
func (o *OverlayContainer) Show() {
	if !o.shown {
		o.canvas.Overlays().Add(o)
		o.shown = true
	}
	o.show(o)
}

// Tapped satisfies the fyne.Tappable interface.
func (o *OverlayContainer) Tapped(*fyne.PointEvent) {
	if o.dismissAction != nil {
		o.dismissAction()
	}
}

type overlayRenderer struct {
	BaseRenderer
	o *OverlayContainer
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
