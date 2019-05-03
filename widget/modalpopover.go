package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// ModalPopOver is a widget that floats above the user interface blocking interaction with any other elements.
// It wraps any standard content and centers it within the space over a background that indicates inactivity below.
type ModalPopOver struct {
	baseWidget

	Content fyne.CanvasObject
	Canvas  fyne.Canvas
}

// Hide this widget, if it was previously visible
func (m *ModalPopOver) Hide() {
	m.hide(m)
	m.Canvas.SetOverlay(nil)
}

// MinSize returns the smallest size this widget can shrink to
func (m *ModalPopOver) MinSize() fyne.Size {
	return m.minSize(m)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (m *ModalPopOver) Move(pos fyne.Position) {
	m.move(pos, m)
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (m *ModalPopOver) Resize(size fyne.Size) {
	m.resize(size, m)
}

// Show this widget, if it was previously hidden
func (m *ModalPopOver) Show() {
	m.Canvas.SetOverlay(m)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (m *ModalPopOver) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(theme.BackgroundColor())
	return &modalPopoverRenderer{center: layout.NewCenterLayout(), popover: m, bg: bg,
		objects: []fyne.CanvasObject{bg, m.Content}}
}

// NewModalPopOver creates a new modal popover for the specified content and displays it on the passed canvas.
func NewModalPopOver(content fyne.CanvasObject, canvas fyne.Canvas) *ModalPopOver {
	ret := &ModalPopOver{Content: content, Canvas: canvas}
	ret.Show()
	return ret
}

type modalPopoverRenderer struct {
	center  fyne.Layout
	popover *ModalPopOver
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *modalPopoverRenderer) Layout(size fyne.Size) {
	r.center.Layout(r.objects, size)

	r.bg.Move(r.popover.Content.Position().Subtract(fyne.NewPos(theme.Padding(), theme.Padding())))
	r.bg.Resize(r.MinSize())
}

func (r *modalPopoverRenderer) MinSize() fyne.Size {
	return r.popover.Content.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (r *modalPopoverRenderer) Refresh() {
}

func (r *modalPopoverRenderer) ApplyTheme() {
	r.bg.FillColor = theme.BackgroundColor()
}

func (r *modalPopoverRenderer) BackgroundColor() color.Color {
	return theme.ShadowColor()
}

func (r *modalPopoverRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *modalPopoverRenderer) Destroy() {
}
