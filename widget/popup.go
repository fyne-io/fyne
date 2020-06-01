package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// PopUp is a widget that can float above the user interface.
// It wraps any standard elements with padding and a shadow.
// If it is modal then the shadow will cover the entire canvas it hovers over and block interactions.
type PopUp struct {
	BaseWidget

	Content fyne.CanvasObject
	Canvas  fyne.Canvas

	innerPos     fyne.Position
	innerSize    fyne.Size
	modal        bool
	overlayShown bool
}

// Hide this widget, if it was previously visible
func (p *PopUp) Hide() {
	if p.overlayShown {
		p.Canvas.Overlays().Remove(p)
		p.overlayShown = false
	}
	p.BaseWidget.Hide()
}

// Move the widget to a new position. A PopUp position is absolute to the top, left of its canvas.
// For PopUp this actually moves the content so checking Position() will not return the same value as is set here.
func (p *PopUp) Move(pos fyne.Position) {
	if p.modal {
		return
	}
	p.innerPos = pos
	p.Refresh()
}

// Resize changes the size of the PopUp.
// PopUps always have the size of their canvas.
// However, Resize changes the size of the PopUp's content.
// Implements: fyne.Widget
func (p *PopUp) Resize(size fyne.Size) {
	p.innerSize = size
	p.BaseWidget.Resize(p.Canvas.Size())
	// The canvas size might not have changed and therefore the Resize won't trigger a layout.
	// Until we have a widget.Relayout() or similar, the renderer's refresh will do the re-layout.
	p.Refresh()
}

// Show this pop-up as overlay if not already shown.
func (p *PopUp) Show() {
	if !p.overlayShown {
		if p.Size().IsZero() {
			p.Resize(p.MinSize())
		}
		p.Canvas.Overlays().Add(p)
		p.overlayShown = true
	}
	p.BaseWidget.Show()
}

// ShowAtPosition shows this pop-up at the given position.
func (p *PopUp) ShowAtPosition(pos fyne.Position) {
	p.Move(pos)
	p.Show()
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

// MinSize returns the size that this widget should not shrink below
func (p *PopUp) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *PopUp) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	bg := canvas.NewRectangle(theme.BackgroundColor())
	objects := []fyne.CanvasObject{bg, p.Content}
	if p.modal {
		return &modalPopUpRenderer{
			widget.NewBaseRenderer(objects),
			popUpBaseRenderer{popUp: p, bg: bg},
		}
	}

	return &popUpRenderer{
		widget.NewShadowingRenderer(objects, widget.PopUpLevel),
		popUpBaseRenderer{popUp: p, bg: bg},
	}
}

// NewPopUpAtPosition creates a new popUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
// Deprecated: Use ShowPopUpAtPosition() instead.
func NewPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, pos fyne.Position) *PopUp {
	p := newPopUp(content, canvas)
	p.ShowAtPosition(pos)
	return p
}

// ShowPopUpAtPosition creates a new popUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
func ShowPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, pos fyne.Position) {
	newPopUp(content, canvas).ShowAtPosition(pos)
}

func newPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	ret := &PopUp{Content: content, Canvas: canvas, modal: false}
	ret.ExtendBaseWidget(ret)
	return ret
}

// NewPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// Deprecated: This will no longer show the pop-up in 2.0. Use ShowPopUp() instead.
func NewPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	return NewPopUpAtPosition(content, canvas, fyne.NewPos(0, 0))
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
// Deprecated: This will no longer show the pop-up in 2.0. Use ShowModalPopUp instead.
func NewModalPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	p := newModalPopUp(content, canvas)
	p.Show()
	return p
}

// ShowModalPopUp creates a new popUp for the specified content and displays it on the passed canvas.
// A modal PopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func ShowModalPopUp(content fyne.CanvasObject, canvas fyne.Canvas) {
	p := newModalPopUp(content, canvas)
	p.Show()
}

type popUpBaseRenderer struct {
	popUp *PopUp
	bg    *canvas.Rectangle
}

func (r *popUpBaseRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}

func (r *popUpBaseRenderer) offset() fyne.Position {
	return fyne.NewPos(theme.Padding(), theme.Padding())
}

type popUpRenderer struct {
	*widget.ShadowingRenderer
	popUpBaseRenderer
}

func (r *popUpRenderer) Layout(_ fyne.Size) {
	r.popUp.Content.Resize(r.popUp.innerSize.Subtract(r.padding()))

	innerPos := r.popUp.innerPos
	if innerPos.X+r.popUp.innerSize.Width > r.popUp.Canvas.Size().Width {
		innerPos.X = r.popUp.Canvas.Size().Width - r.popUp.innerSize.Width
		if innerPos.X < 0 {
			innerPos.X = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if innerPos.Y+r.popUp.innerSize.Height > r.popUp.Canvas.Size().Height {
		innerPos.Y = r.popUp.Canvas.Size().Height - r.popUp.innerSize.Height
		if innerPos.Y < 0 {
			innerPos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	r.popUp.Content.Move(innerPos.Add(r.offset()))

	r.bg.Resize(r.popUp.innerSize)
	r.bg.Move(innerPos)
	r.LayoutShadow(r.popUp.innerSize, innerPos)
}

func (r *popUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize().Add(r.padding())
}

func (r *popUpRenderer) Refresh() {
	r.bg.FillColor = theme.BackgroundColor()
	if r.bg.Size() != r.popUp.innerSize || r.bg.Position() != r.popUp.innerPos {
		r.Layout(r.popUp.Size())
	}
}

func (r *popUpRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

type modalPopUpRenderer struct {
	widget.BaseRenderer
	popUpBaseRenderer
}

func (r *modalPopUpRenderer) Layout(canvasSize fyne.Size) {
	padding := r.padding()
	requestedSize := r.popUp.innerSize.Subtract(padding)
	size := r.popUp.Content.MinSize().Max(requestedSize)
	size = size.Min(canvasSize.Subtract(padding))
	pos := fyne.NewPos((canvasSize.Width-size.Width)/2, (canvasSize.Height-size.Height)/2)
	r.popUp.Content.Move(pos)
	r.popUp.Content.Resize(size)

	r.bg.Move(pos.Subtract(r.offset()))
	r.bg.Resize(size.Add(padding))
}

func (r *modalPopUpRenderer) MinSize() fyne.Size {
	return r.popUp.Content.MinSize().Add(r.padding())
}

func (r *modalPopUpRenderer) Refresh() {
	r.bg.FillColor = theme.BackgroundColor()
	if r.bg.Size() != r.popUp.innerSize {
		r.Layout(r.popUp.Size())
	}
}

func (r *modalPopUpRenderer) BackgroundColor() color.Color {
	return theme.ShadowColor()
}
