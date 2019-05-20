package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// PopOver is a widget that can float above the user interface.
// It wraps any standard elements with padding and a shadow.
// If it is modal then the shadow will cover the entire canvas it hovers over and block interactions.
type PopOver struct {
	baseWidget

	Content fyne.CanvasObject
	Canvas  fyne.Canvas

	modal bool
}

// Hide this widget, if it was previously visible
func (p *PopOver) Hide() {
	p.hide(p)
	p.Canvas.SetOverlay(nil)
}

// MinSize returns the smallest size this widget can shrink to
func (p *PopOver) MinSize() fyne.Size {
	return p.minSize(p)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (p *PopOver) Move(pos fyne.Position) {
	if p.modal {
		return
	}

	p.Content.Move(pos.Add(fyne.NewPos(theme.Padding(), theme.Padding())))
	Renderer(p).Layout(p.Size())
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (p *PopOver) Resize(size fyne.Size) {
	p.resize(size, p)
}

// Show this widget, if it was previously hidden
func (p *PopOver) Show() {
	p.show(p)
	p.Canvas.SetOverlay(p)
}

// Tapped is called when the user taps the popover background - if not modal then dismiss this widget
func (p *PopOver) Tapped(_ *fyne.PointEvent) {
	if !p.modal {
		p.Hide()
	}
}

// TappedSecondary is called when the user right/alt taps the background - ignore
func (p *PopOver) TappedSecondary(_ *fyne.PointEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *PopOver) CreateRenderer() fyne.WidgetRenderer {
	if p.modal {
		bg := canvas.NewRectangle(theme.BackgroundColor())
		return &modalPopoverRenderer{center: layout.NewCenterLayout(), popover: p, bg: bg,
			objects: []fyne.CanvasObject{bg, p.Content}}
	}

	shadow := canvas.NewRectangle(theme.ShadowColor())
	bg := canvas.NewRectangle(theme.BackgroundColor())
	objects := []fyne.CanvasObject{shadow, bg, p.Content}
	return &popoverRenderer{popover: p, shadow: shadow, bg: bg, objects: objects}
}

// NewPopOver creates a new popover for the specified content and displays it on the passed canvas.
func NewPopOver(content fyne.CanvasObject, canvas fyne.Canvas) *PopOver {
	ret := &PopOver{Content: content, Canvas: canvas, modal: false}
	ret.Show()
	return ret
}

// NewModalPopOver creates a new popover for the specified content and displays it on the passed canvas.
// A modal PopOver blocks interactions with underlying elements, covered with a semi-transparent overlay.
func NewModalPopOver(content fyne.CanvasObject, canvas fyne.Canvas) *PopOver {
	ret := &PopOver{Content: content, Canvas: canvas, modal: true}
	ret.Show()
	return ret
}

type popoverRenderer struct {
	popover    *PopOver
	shadow, bg *canvas.Rectangle
	objects    []fyne.CanvasObject
}

func (r *popoverRenderer) Layout(size fyne.Size) {
	pos := r.popover.Content.Position()
	innerSize := r.popover.Content.MinSize()
	r.popover.Content.Resize(innerSize)
	r.popover.Content.Move(pos.Add(fyne.NewPos(theme.Padding(), theme.Padding())))

	size = innerSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	r.shadow.Resize(size.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	r.shadow.Move(pos.Subtract(fyne.NewPos(theme.Padding(), theme.Padding())))

	r.bg.Resize(size)
	r.bg.Move(pos)
}

func (r *popoverRenderer) MinSize() fyne.Size {
	return r.popover.Content.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (r *popoverRenderer) Refresh() {
}

func (r *popoverRenderer) ApplyTheme() {
}

func (r *popoverRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *popoverRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *popoverRenderer) Destroy() {
}

type modalPopoverRenderer struct {
	center  fyne.Layout
	popover *PopOver
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
