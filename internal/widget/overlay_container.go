package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var (
	_ fyne.Widget       = (*OverlayContainer)(nil)
	_ fyne.Tappable     = (*OverlayContainer)(nil)
	_ desktop.Hoverable = (*OverlayContainer)(nil)
)

// OverlayContainer is a transparent widget containing one fyne.CanvasObject and meant to be used as overlay.
type OverlayContainer struct {
	Base
	Content fyne.CanvasObject

	canvas    fyne.Canvas
	onDismiss func()
	shown     bool
}

// NewOverlayContainer creates an OverlayContainer.
func NewOverlayContainer(c fyne.CanvasObject, canvas fyne.Canvas, onDismiss func()) *OverlayContainer {
	o := &OverlayContainer{canvas: canvas, Content: c, onDismiss: onDismiss}
	o.ExtendBaseWidget(o)
	return o
}

// CreateRenderer returns a new renderer for the overlay container.
func (o *OverlayContainer) CreateRenderer() fyne.WidgetRenderer {
	return &overlayRenderer{BaseRenderer{[]fyne.CanvasObject{o.Content}}, o}
}

// Hide hides the overlay container.
func (o *OverlayContainer) Hide() {
	if o.shown {
		o.canvas.Overlays().Remove(o)
		o.shown = false
	}
	o.Base.Hide()
}

// MouseIn catches mouse-in events not handled by the container’s content. It does nothing.
func (o *OverlayContainer) MouseIn(*desktop.MouseEvent) {
}

// MouseMoved catches mouse-moved events not handled by the container’s content. It does nothing.
func (o *OverlayContainer) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut catches mouse-out events not handled by the container’s content. It does nothing.
func (o *OverlayContainer) MouseOut() {
}

// SetCanvas allows an overlay container to be re-used on a different canvas.
//
// Since: 2.8
func (o *OverlayContainer) SetCanvas(c fyne.Canvas) {
	o.canvas.Overlays().Remove(o)
	o.canvas = c
	o.canvas.Overlays().Add(o)
}

// Show makes the overlay container visible.
func (o *OverlayContainer) Show() {
	if !o.shown {
		o.canvas.Overlays().Add(o)
		o.shown = true
	}
	o.Base.Show()
}

// Tapped catches tap events not handled by the container’s content.
// It performs the overlay container’s dismiss action.
func (o *OverlayContainer) Tapped(*fyne.PointEvent) {
	if o.onDismiss != nil {
		o.onDismiss()
	}
}

type overlayRenderer struct {
	BaseRenderer
	o *OverlayContainer
}

var _ fyne.WidgetRenderer = (*overlayRenderer)(nil)

func (r *overlayRenderer) Layout(fyne.Size) {
}

func (r *overlayRenderer) MinSize() fyne.Size {
	return r.o.canvas.Size()
}

func (r *overlayRenderer) Refresh() {
}
