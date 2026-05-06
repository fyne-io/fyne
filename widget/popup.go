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

	overlay *widget.OverlayContainer
	modal   bool
}

// Hide this widget, if it was previously visible
func (p *PopUp) Hide() {
	if p.overlay != nil {
		p.Canvas.Overlays().Remove(p.overlay)
		p.overlay = nil
	}
	p.BaseWidget.Hide()
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

// AccessibilityRole returns the role used to describe this popup to
// assistive technologies.
//
// Since: 2.8
func (p *PopUp) AccessibilityRole() fyne.AccessibleRole {
	return fyne.AccessibleRoleContainer
}

// AccessibilityLabel returns the text used by assistive technologies as the
// name of this popup.
//
// Since: 2.8
func (p *PopUp) AccessibilityLabel() string {
	return ""
}

// AccessibilityChildren returns the popup content so it can be navigated by
// assistive technologies.
//
// Since: 2.8
func (p *PopUp) AccessibilityChildren() []fyne.CanvasObject {
	if p.Content == nil {
		return nil
	}
	return []fyne.CanvasObject{p.Content}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *PopUp) CreateRenderer() fyne.WidgetRenderer {
	th := p.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	p.ExtendBaseWidget(p)
	background := canvas.NewRectangle(th.Color(theme.ColorNameOverlayBackground, v))
	if p.modal {
		objects := []fyne.CanvasObject{background, p.Content}
		return &modalPopUpRenderer{
			widget.NewShadowingRenderer(objects, widget.DialogLevel),
			popUpBaseRenderer{popUp: p, background: background},
		}
	}
	objects := []fyne.CanvasObject{background, p.Content}
	return &popUpRenderer{
		widget.NewShadowingRenderer(objects, widget.PopUpLevel),
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

func (r *popUpBaseRenderer) offset() fyne.Position {
	th := r.popUp.Theme()
	return fyne.NewSquareOffsetPos(th.Size(theme.SizeNameInnerPadding) / 2)
}

type popUpRenderer struct {
	*widget.ShadowingRenderer
	popUpBaseRenderer
}

func (r *popUpRenderer) Layout(s fyne.Size) {
	innerPos := r.popUp.Content.Position()
	padding := r.padding()
	size := r.popUp.Size().Max(r.popUp.Content.MinSize().Add(padding))
	innerSize := size.Subtract(padding)

	canvasSize := r.popUp.Canvas.Size()
	innerSize = innerSize.Max(r.popUp.Content.MinSize())
	if !canvasSize.IsZero() {
		innerSize = innerSize.Min(canvasSize)
	}
	r.popUp.Content.Resize(innerSize)

	if innerPos.X+innerSize.Width > r.popUp.Canvas.Size().Width {
		innerPos.X = r.popUp.Canvas.Size().Width - innerSize.Width
		if innerPos.X < 0 {
			innerPos.X = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if innerPos.Y+innerSize.Height > r.popUp.Canvas.Size().Height {
		innerPos.Y = r.popUp.Canvas.Size().Height - innerSize.Height
		if innerPos.Y < 0 {
			innerPos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}

	r.background.Resize(innerSize)
	r.background.Move(innerPos)
	r.LayoutShadow(innerSize, innerPos)
}

func (r *popUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize().Add(r.padding())
}

func (r *popUpRenderer) Refresh() {
	innerPos := r.popUp.Content.Position()
	innerSize := r.popUp.Content.Size()

	th := r.popUp.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	r.background.FillColor = th.Color(theme.ColorNameOverlayBackground, v)
	expectedContentSize := innerSize.Max(r.popUp.MinSize()).Subtract(r.padding())
	shouldRelayout := r.popUp.Content.Size() != expectedContentSize

	if r.background.Size() != innerSize || r.background.Position() != innerPos || shouldRelayout {
		r.Layout(r.popUp.Size())
	}
	r.popUp.Content.Refresh()
	r.background.Refresh()
	r.ShadowingRenderer.RefreshShadow()
}

type modalPopUpRenderer struct {
	*widget.ShadowingRenderer
	popUpBaseRenderer
}

func (r *modalPopUpRenderer) Layout(_ fyne.Size) {
	canvasSize := r.popUp.Canvas.Size()
	padding := r.padding()
	innerSize := r.popUp.Size().Max(r.popUp.Content.MinSize().Add(padding))
	if !canvasSize.IsZero() {
		innerSize = innerSize.Min(canvasSize)
	}

	size := innerSize.Subtract(padding)
	if !canvasSize.IsZero() {
		size = size.Min(canvasSize.Subtract(padding))
	}
	pos := r.popUp.Position()
	if pos.IsZero() {
		pos = fyne.NewPos((canvasSize.Width-size.Width)/2, (canvasSize.Height-size.Height)/2)
		r.popUp.Content.Move(pos)
	}
	r.popUp.Content.Resize(size)

	innerPos := pos.Subtract(r.offset())
	r.background.Move(innerPos)
	r.background.Resize(size.Add(padding))
	r.LayoutShadow(innerSize, innerPos)
}

func (r *modalPopUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize().Add(r.padding())
}

func (r *modalPopUpRenderer) Refresh() {
	th := r.popUp.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.background.FillColor = th.Color(theme.ColorNameOverlayBackground, v)
	r.popUp.Content.Refresh()
	r.background.Refresh()
	r.ShadowingRenderer.RefreshShadow()
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
