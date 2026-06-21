package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*PopUp)(nil)

// PopUp is a widget that can float above the user interface.
// It wraps any standard elements with padding and a shadow.
// If it is modal then the shadow will cover the entire canvas it hovers over and block interactions.
type PopUp struct {
	BaseWidget

	Content fyne.CanvasObject
	Canvas  fyne.Canvas

	overlay       *widget.OverlayContainer
	modal, manual bool
}

// Hide this widget, if it was previously visible
func (p *PopUp) Hide() {
	if p.overlay != nil {
		p.Canvas.Overlays().Remove(p.overlay)
		p.overlay = nil
	}

	p.BaseWidget.Hide()
	p.Move(fyne.Position{}) // reset so that Show or ShowAtPosition next time is respected
}

// Refresh the background for a modal popup and the content of this popup.
func (p *PopUp) Refresh() {
	if p.modal && p.overlay != nil {
		th := p.Theme()
		v := fyne.CurrentApp().Settings().ThemeVariant()

		bg := p.overlay.Background.(*fyne.Container).Objects[1].(*canvas.Rectangle)
		bg.FillColor = th.Color(theme.ColorNameShadow, v)
	}

	p.BaseWidget.Refresh()
}

// Show this pop-up as overlay if not already shown.
func (p *PopUp) Show() {
	if p.overlay == nil {
		dismiss := p.Hide
		if p.modal {
			dismiss = nil
		}
		p.overlay = widget.NewOverlayContainer(p.super(), p.Canvas, dismiss)
		if p.modal {
			th := p.Theme()
			v := fyne.CurrentApp().Settings().ThemeVariant()

			p.overlay.Background = &fyne.Container{
				Layout: layout.NewStackLayout(),
				Objects: []fyne.CanvasObject{
					canvas.NewBlur(th.Size(theme.SizeNameModalBlurRadius)),
					canvas.NewRectangle(th.Color(theme.ColorNameShadow, v)),
				},
			}
		}
		p.Canvas.Overlays().Add(p.overlay)
	}
	p.Refresh()
	p.BaseWidget.Show()
}

// ShowAtPosition shows this pop-up at the given position.
func (p *PopUp) ShowAtPosition(pos fyne.Position) {
	p.manual = true
	p.Move(pos)
	p.Show()
}

// ShowAtRelativePosition shows this pop-up at the given position relative to stated object.
//
// Since 2.4
func (p *PopUp) ShowAtRelativePosition(rel fyne.Position, to fyne.CanvasObject) {
	withRelativePosition(rel, to, p.ShowAtPosition)
}

// Tapped is called when the user taps the popUp.
// This is not called when tapping the background, but non-modal popups will dismiss when tapped outside.
func (p *PopUp) Tapped(*fyne.PointEvent) {
}

// TappedSecondary is called when the user right/alt taps the popUp.
// This is not called when tapping the background, but non-modal popups will dismiss when tapped outside.
func (p *PopUp) TappedSecondary(*fyne.PointEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (p *PopUp) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *PopUp) CreateRenderer() fyne.WidgetRenderer {
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	p.ExtendBaseWidget(p)
	background := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, v))
	widget.ApplyShadowForLevel(&background.Shadow, widget.PopUpLevel, th.Color(theme.ColorNameShadow, v))
	background.CornerRadius = th.Size(theme.SizeNamePopupRadius)
	objects := []fyne.CanvasObject{background, p.Content}
	return &popUpRenderer{
		widget.NewBaseRenderer(objects),
		popUpBaseRenderer{popUp: p, background: background},
	}
}

// ShowPopUpAtPosition creates a new popUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
func ShowPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, pos fyne.Position) {
	newPopUp(content, canvas).ShowAtPosition(pos)
}

// ShowPopUpAtRelativePosition shows a new popUp for the specified content at the given position relative to stated object.
// It will then display the popup on the passed canvas.
//
// Since 2.4
func ShowPopUpAtRelativePosition(content fyne.CanvasObject, canvas fyne.Canvas, rel fyne.Position, to fyne.CanvasObject) {
	withRelativePosition(rel, to, func(pos fyne.Position) {
		ShowPopUpAtPosition(content, canvas, pos)
	})
}

func newPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	ret := &PopUp{Content: content, Canvas: canvas, modal: false}
	ret.ExtendBaseWidget(ret)
	return ret
}

// NewPopUp creates a new popUp for the specified content and displays it on the passed canvas.
func NewPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	return newPopUp(content, canvas)
}

// ShowPopUp creates a new popUp for the specified content and displays it on the passed canvas.
func ShowPopUp(content fyne.CanvasObject, canvas fyne.Canvas) {
	newPopUp(content, canvas).Show()
}

func newModalPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	p := &PopUp{Content: content, Canvas: canvas, modal: true}
	p.ExtendBaseWidget(p)
	return p
}

// NewModalPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// A modal PopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func NewModalPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	return newModalPopUp(content, canvas)
}

// ShowModalPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// A modal PopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func ShowModalPopUp(content fyne.CanvasObject, canvas fyne.Canvas) {
	p := newModalPopUp(content, canvas)
	p.Show()
}

type popUpBaseRenderer struct {
	popUp      *PopUp
	background *canvas.Rectangle
}

func (r *popUpBaseRenderer) padding() fyne.Size {
	th := r.popUp.Theme()
	return fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding))
}

type popUpRenderer struct {
	widget.BaseRenderer
	popUpBaseRenderer
}

func (r *popUpRenderer) Layout(s fyne.Size) {
	size := s.Max(r.popUp.Content.MinSize())
	if r.popUp.Canvas != nil {
		canvasSize := r.popUp.Canvas.Size()
		if !canvasSize.IsZero() {
			size = size.Min(r.popUp.Canvas.Size())
		}
	}
	r.popUp.Content.Resize(size)

	r.background.Resize(size)
}

func (r *popUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize()
}

func (r *popUpRenderer) Refresh() {
	innerPos := r.popUp.Content.Position()
	innerSize := r.popUp.Content.Size()

	th := r.popUp.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	r.background.FillColor = th.Color(theme.ColorNameOverlayBackground, v)
	r.background.Shadow.Color = th.Color(theme.ColorNameShadow, v)
	r.background.CornerRadius = th.Size(theme.SizeNamePopupRadius)
	expectedContentSize := innerSize.Max(r.popUp.MinSize()).Subtract(r.padding())
	shouldRelayout := r.popUp.Content.Size() != expectedContentSize

	if r.background.Size() != innerSize || r.background.Position() != innerPos || shouldRelayout {
		r.Layout(r.popUp.Size())
	}
	r.popUp.Content.Refresh()
	r.background.Refresh()
}

func withRelativePosition(rel fyne.Position, to fyne.CanvasObject, f func(position fyne.Position)) {
	d := fyne.CurrentApp().Driver()
	c := d.CanvasForObject(to)
	if c == nil {
		fyne.LogError("Could not locate parent object to display relative to", nil)
		f(rel)
		return
	}

	pos := d.AbsolutePositionForObject(to).Add(rel)
	f(pos)
}
