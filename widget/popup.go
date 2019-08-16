package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// PopUp is a widget that can float above the user interface.
// It wraps any standard elements with padding and a shadow.
// If it is modal then the shadow will cover the entire canvas it hovers over and block interactions.
type PopUp struct {
	baseWidget

	Content fyne.CanvasObject
	Canvas  fyne.Canvas

	modal bool
}

// Hide this widget, if it was previously visible
func (p *PopUp) Hide() {
	p.hide(p)
	p.Canvas.SetOverlay(nil)
}

// MinSize returns the smallest size this widget can shrink to
func (p *PopUp) MinSize() fyne.Size {
	return p.minSize(p)
}

// Move the widget to a new position. A PopUp position is absolute to the top, left of its canvas.
// For PopUp this actually moves the content so checking Position() will not return the same value as is set here.
func (p *PopUp) Move(pos fyne.Position) {
	if p.modal {
		return
	}

	if pos.Y+p.Content.Size().Height > p.Canvas.Size().Height-theme.Padding()*2 {
		pos.Y = p.Canvas.Size().Height - p.Content.Size().Height - theme.Padding()*2
		if pos.Y < 0 {
			pos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}

	p.Content.Move(pos.Add(fyne.NewPos(theme.Padding(), theme.Padding())))
	Renderer(p).Layout(p.Size())
}

// Resize sets a new size for a widget. Most PopUp widgets are shown at MinSize.
func (p *PopUp) Resize(size fyne.Size) {
	p.resize(size, p)
}

// Show this widget, if it was previously hidden
func (p *PopUp) Show() {
	p.show(p)
	p.Canvas.SetOverlay(p)
}

// Tapped is called when the user taps the popUp background - if not modal then dismiss this widget
func (p *PopUp) Tapped(_ *fyne.PointEvent) {
	if !p.modal {
		p.Hide()
	}
}

// TappedSecondary is called when the user right/alt taps the background - if not modal then dismiss this widget
func (p *PopUp) TappedSecondary(_ *fyne.PointEvent) {
	if !p.modal {
		p.Hide()
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *PopUp) CreateRenderer() fyne.WidgetRenderer {
	if p.modal {
		bg := canvas.NewRectangle(theme.BackgroundColor())
		return &modalPopUpRenderer{center: layout.NewCenterLayout(), popUp: p, bg: bg,
			objects: []fyne.CanvasObject{bg, p.Content}}
	}

	shadow := newShadow(shadowAround, theme.Padding()*2)
	bg := canvas.NewRectangle(theme.BackgroundColor())
	objects := []fyne.CanvasObject{shadow, bg, p.Content}
	return &popUpRenderer{popUp: p, shadow: shadow, bg: bg, objects: objects}
}

// NewPopUp creates a new popUp for the specified content and displays it on the passed canvas.
func NewPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	ret := &PopUp{Content: content, Canvas: canvas, modal: false}
	ret.Show()
	return ret
}

// NewModalPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// A modal PopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func NewModalPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	ret := &PopUp{Content: content, Canvas: canvas, modal: true}
	ret.Show()
	return ret
}

type popUpRenderer struct {
	popUp   *PopUp
	shadow  fyne.CanvasObject
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *popUpRenderer) Layout(size fyne.Size) {
	pos := r.popUp.Content.Position().Subtract(fyne.NewPos(theme.Padding(), theme.Padding()))
	innerSize := r.popUp.Content.MinSize()
	r.popUp.Content.Resize(innerSize)
	r.popUp.Content.Move(pos.Add(fyne.NewPos(theme.Padding(), theme.Padding())))

	size = innerSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	r.bg.Resize(size)
	r.bg.Move(pos)
	r.shadow.Resize(size)
	r.shadow.Move(pos)
}

func (r *popUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (r *popUpRenderer) Refresh() {
}

func (r *popUpRenderer) ApplyTheme() {
	r.bg.FillColor = theme.BackgroundColor()
}

func (r *popUpRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *popUpRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *popUpRenderer) Destroy() {
}

type modalPopUpRenderer struct {
	center  fyne.Layout
	popUp   *PopUp
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *modalPopUpRenderer) Layout(size fyne.Size) {
	r.center.Layout(r.objects, size)

	r.bg.Move(r.popUp.Content.Position().Subtract(fyne.NewPos(theme.Padding(), theme.Padding())))
	r.bg.Resize(r.MinSize())
}

func (r *modalPopUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (r *modalPopUpRenderer) Refresh() {
}

func (r *modalPopUpRenderer) ApplyTheme() {
	r.bg.FillColor = theme.BackgroundColor()
}

func (r *modalPopUpRenderer) BackgroundColor() color.Color {
	return theme.ShadowColor()
}

func (r *modalPopUpRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *modalPopUpRenderer) Destroy() {
}
