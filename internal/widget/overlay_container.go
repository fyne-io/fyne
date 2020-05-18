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

// CreateRenderer returns a new renderer for the overlay container.
// Implements: fyne.Widget
func (o *OverlayContainer) CreateRenderer() fyne.WidgetRenderer {
	return &overlayRenderer{BaseRenderer{[]fyne.CanvasObject{o.Content}}, o}
}

// Hide hides the overlay container.
// Implements: fyne.Widget
func (o *OverlayContainer) Hide() {
	if o.shown {
		o.canvas.Overlays().Remove(o)
		o.shown = false
	}
	o.hide(o)
}

// MinSize returns the minimal size of the overlay container.
// This is the same as the size of the canvas.
// Implements: fyne.Widget
func (o *OverlayContainer) MinSize() fyne.Size {
	return o.minSize(o)
}

// MouseIn catches mouse-in events not handled by the container’s content. It does nothing.
// Implements: desktop.Hoverable
func (o *OverlayContainer) MouseIn(*desktop.MouseEvent) {
}

// MouseMoved catches mouse-moved events not handled by the container’s content. It does nothing.
// Implements: desktop.Hoverable
func (o *OverlayContainer) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut catches mouse-out events not handled by the container’s content. It does nothing.
// Implements: desktop.Hoverable
func (o *OverlayContainer) MouseOut() {
}

// Refresh triggers a redraw of the overlay container.
// Implements: fyne.Widget
func (o *OverlayContainer) Refresh() {
	o.refresh(o)
}

// Resize changes the size of the overlay container.
// This is normally called by the canvas overlay management.
// Implements: fyne.Widget
func (o *OverlayContainer) Resize(size fyne.Size) {
	o.resize(size, o)
}

// Show makes the overlay container visible.
// Implements: fyne.Widget
func (o *OverlayContainer) Show() {
	if !o.shown {
		o.canvas.Overlays().Add(o)
		o.shown = true
	}
	o.show(o)
}

// Tapped catches tap events not handled by the container’s content.
// It performs the overlay container’s dismiss action.
// Implements: fyne.Tappable
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

func (r *overlayRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *overlayRenderer) Layout(fyne.Size) {
}

func (r *overlayRenderer) MinSize() fyne.Size {
	return r.o.canvas.Size()
}

func (r *overlayRenderer) Refresh() {
}
