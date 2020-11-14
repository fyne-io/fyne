package widget

import (
	"math"
	
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var _ fyne.Draggable = (*TickerPopUp)(nil)

type ringBuffer struct {
	data []byte
	start int // Start of ringBuffer can be anywhere in array.
	bound int
	forward bool
}

func (rb *ringBuffer) Init(start int, data []byte) {
	rb.data = data
	rb.start = start
}

// Turn - rotates the ringbuffer by appropriate offset.  -offset is left, +offset is right.
func (rb *ringBuffer) Turn(offset int) {
	if rb.forward {
		rb.start = rb.start - offset
	} else {
		rb.start = rb.start + offset
	}
	if rb.start > len(rb.data) {
		rb.start = rb.start % len(rb.data)
	} else if rb.start < 0 {
		rb.start = len(rb.data) + rb.start
	}
}

// Turn - rotates the ringbuffer by appropriate offset.  -offset is left, +offset is right.
func (rb *ringBuffer) Seek(position int) {
	if position >= rb.start {
		rb.Turn(position - rb.start)
	} else {
		rb.Turn(-(rb.start - position))
	}
}

// Data - returns current data at current turn, read circularly
func (rb *ringBuffer) Data() []byte {
	var data []byte
	if rb.start == 0 {
		data = rb.data
	} else {
		data = append(rb.data[rb.start:], rb.data[0: rb.start - 1]...)
	}

	if rb.bound > 0 {
		return data[0:rb.bound - 1]
	} else {
		return data
	}
}

func (rb *ringBuffer) Length() int {
	return len(rb.data)
}

// TickerPopUp is a widget that can float above the user interface.
// It wraps any standard elements with padding and a shadow.
// If it is modal then the shadow will cover the entire canvas it hovers over and block interactions.
type TickerPopUp struct {
	BaseWidget

	Content fyne.CanvasObject
	rb      ringBuffer // backing the content with a ringBuffer.
	Canvas  fyne.Canvas

	innerPos     fyne.Position
	innerSize    fyne.Size
	modal        bool
	overlayShown bool
	draggedX     fyne.Position
}

// Hide this widget, if it was previously visible
func (p *TickerPopUp) Hide() {
	if p.overlayShown {
		p.Canvas.Overlays().Remove(p)
		p.overlayShown = false
	}
	p.BaseWidget.Hide()
}

// Move the widget to a new position. A TickerPopUp position is absolute to the top, left of its canvas.
// For TickerPopUp this actually moves the content so checking Position() will not return the same value as is set here.
func (p *TickerPopUp) Move(pos fyne.Position) {
	if p.modal {
		return
	}
	p.innerPos = pos
	p.Refresh()
}

// Resize changes the size of the TickerPopUp.
// TickerPopUps always have the size of their canvas.
// However, Resize changes the size of the TickerPopUp's content.
//
// Implements: fyne.Widget
func (p *TickerPopUp) Resize(size fyne.Size) {
	p.innerSize = size
	p.BaseWidget.Resize(p.Canvas.Size())
	// The canvas size might not have changed and therefore the Resize won't trigger a layout.
	// Until we have a widget.Relayout() or similar, the renderer's refresh will do the re-layout.
	p.Refresh()
}

// Show this pop-up as overlay if not already shown.
func (p *TickerPopUp) Show() {
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
func (p *TickerPopUp) ShowAtPosition(pos fyne.Position) {
	p.Move(pos)
	p.Show()
}


// DragEnd function.
func (p *TickerPopUp) DragEnd() {
}

// Tapped is called when the user taps the tickerPopUp background - if not modal then dismiss this widget
func (p *TickerPopUp) Tapped(e *fyne.PointEvent) {
	// TODO: Calculate ratio..
	// Determine proportional character index.
	// Choose appropriate selectable.
//	if !p.modal {
//		if label, ok := p.Content.(*Label); ok {

//		}
//		p.Hide()
//	}
}

func (p *TickerPopUp) endOffset() int {
	return p.innerPos.X + theme.Padding()
}

func (p *TickerPopUp) getRatio(posDiff *fyne.Position) float64 {
	pad := p.endOffset()

	x := posDiff.X

	if x > p.innerPos.X + p.innerSize.Width {
		return 1.0
	} else if x < pad {
		return 0.0
	} else {
		return float64(x-pad) / float64(p.innerPos.X + p.innerSize.Width-pad*2)
	}
	
	return 0.0
}

func (p *TickerPopUp) Dragged(e *fyne.DragEvent) {
	//ratio := p.getRatio(&(e.PointEvent.Position))

	if p.draggedX.X == 0 {
		p.draggedX = e.Position
		return
	}
	var diffPosition fyne.Position
	diffPosition.X = e.Position.X - p.draggedX.X
	diffPosition.Y = e.Position.Y- p.draggedX.Y
	ratio := p.getRatio(&diffPosition)

	if label, ok := p.Content.(*Label); ok {
		writeArea := p.rb.Length()
		if p.rb.bound != 0 {
			writeArea = p.rb.bound
		}
		offset := math.Round(ratio * float64(writeArea))
		p.rb.Seek(int(offset))
		label.Text = string(p.rb.Data())
		label.Refresh() // Sweet!
	}
}

// TappedSecondary is called when the user right/alt taps the background - if not modal then dismiss this widget
func (p *TickerPopUp) TappedSecondary(_ *fyne.PointEvent) {
	if !p.modal {
//		p.Hide()
	}
}

// MinSize returns the size that this widget should not shrink below
func (p *TickerPopUp) MinSize() fyne.Size {
	p.ExtendBaseWidget(p)
	return p.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (p *TickerPopUp) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	bg := canvas.NewRectangle(theme.BackgroundColor())
	objects := []fyne.CanvasObject{bg, p.Content}
	if p.modal {
		return &modalTickerPopUpRenderer{
			widget.NewShadowingRenderer(objects, widget.DialogLevel),
			tickerPopupBaseRenderer{tickerPopUp: p, bg: bg},
		}
	}

	return &tickerPopUpRenderer{
		widget.NewShadowingRenderer(objects, widget.PopUpLevel),
		tickerPopupBaseRenderer{tickerPopUp: p, bg: bg},
	}
}

// NewTickerPopUpAtPosition creates a new tickerPopUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
//
// Deprecated: Use ShowTickerPopUpAtPosition() instead.
func NewTickerPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, pos fyne.Position, size fyne.Size) *TickerPopUp {
	p := newTickerPopUp(content, canvas, size)
	p.ShowAtPosition(pos)
	return p
}

// ShowTickerPopUpAtPosition creates a new tickerPopUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
func ShowTickerPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, pos fyne.Position, size fyne.Size) {
	newTickerPopUp(content, canvas, size).ShowAtPosition(pos)
}

func newTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas, size fyne.Size) *TickerPopUp {
	if label, ok := content.(*Label); ok {

		// Calculate area based on size and reset text.
		characterSize := fyne.MeasureText("M", 8, label.TextStyle)
		numChars := size.Width / characterSize.Width

		rb := ringBuffer{ data: []byte(label.Text), start: 0, bound: numChars, forward: true }
		label.Text = string(rb.Data())

		ret := &TickerPopUp{Content: content, rb: rb, Canvas: canvas, modal: false}
		ret.ExtendBaseWidget(ret)
		ret.Resize(size)
		return ret
	}

	return nil // non-label tickers not supported yet.  Consider what this would mean anyways.
}

// NewTickerPopUp creates a new tickerPopUp for the specified content and displays it on the passed canvas.
//
// Deprecated: This will no longer show the pop-up in 2.0. Use ShowTickerPopUp() instead.
func NewTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas, size fyne.Size) *TickerPopUp {
	return NewTickerPopUpAtPosition(content, canvas, fyne.NewPos(0, 0), size)
}

// ShowTickerPopUp creates a new tickerPopUp for the specified content and displays it on the passed canvas.
func ShowTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas, size fyne.Size) {
	newTickerPopUp(content, canvas, size).Show()
}

func newModalTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *TickerPopUp {
	p := &TickerPopUp{Content: content, Canvas: canvas, modal: true}
	p.ExtendBaseWidget(p)
	return p
}

// NewModalTickerPopUp creates a new tickerPopUp for the specified content and displays it on the passed canvas.
// A modal TickerPopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
//
// Deprecated: This will no longer show the pop-up in 2.0. Use ShowModalTickerPopUp instead.
func NewModalTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *TickerPopUp {
	p := newModalTickerPopUp(content, canvas)
	p.Show()
	return p
}

// ShowModalTickerPopUp creates a new tickerPopUp for the specified content and displays it on the passed canvas.
// A modal TickerPopUp blocks interactions with underlying elements, covered with a semi-transparent overlay.
func ShowModalTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas) {
	p := newModalTickerPopUp(content, canvas)
	p.Show()
}

type tickerPopupBaseRenderer struct {
	tickerPopUp *TickerPopUp
	bg    *canvas.Rectangle
}

func (r *tickerPopupBaseRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}

func (r *tickerPopupBaseRenderer) offset() fyne.Position {
	return fyne.NewPos(theme.Padding(), theme.Padding())
}

type tickerPopUpRenderer struct {
	*widget.ShadowingRenderer
	tickerPopupBaseRenderer
}

func (r *tickerPopUpRenderer) Layout(_ fyne.Size) {
	r.tickerPopUp.Content.Resize(r.tickerPopUp.innerSize.Subtract(r.padding()))

	innerPos := r.tickerPopUp.innerPos
	if innerPos.X+r.tickerPopUp.innerSize.Width > r.tickerPopUp.Canvas.Size().Width {
		innerPos.X = r.tickerPopUp.Canvas.Size().Width - r.tickerPopUp.innerSize.Width
		if innerPos.X < 0 {
			innerPos.X = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if innerPos.Y+r.tickerPopUp.innerSize.Height > r.tickerPopUp.Canvas.Size().Height {
		innerPos.Y = r.tickerPopUp.Canvas.Size().Height - r.tickerPopUp.innerSize.Height
		if innerPos.Y < 0 {
			innerPos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	r.tickerPopUp.Content.Move(innerPos.Add(r.offset()))

	r.bg.Resize(r.tickerPopUp.innerSize)
	r.bg.Move(innerPos)
	r.LayoutShadow(r.tickerPopUp.innerSize, innerPos)
}

func (r *tickerPopUpRenderer) MinSize() fyne.Size {
	return r.tickerPopUp.Content.MinSize().Add(r.padding())
}

func (r *tickerPopUpRenderer) Refresh() {
	r.bg.FillColor = theme.BackgroundColor()
	if r.bg.Size() != r.tickerPopUp.innerSize || r.bg.Position() != r.tickerPopUp.innerPos {
		r.Layout(r.tickerPopUp.Size())
	}
}

func (r *tickerPopUpRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

type modalTickerPopUpRenderer struct {
	*widget.ShadowingRenderer
	tickerPopupBaseRenderer
}

func (r *modalTickerPopUpRenderer) Layout(canvasSize fyne.Size) {
	padding := r.padding()
	requestedSize := r.tickerPopUp.innerSize.Subtract(padding)
	size := r.tickerPopUp.Content.MinSize().Max(requestedSize)
	size = size.Min(canvasSize.Subtract(padding))
	pos := fyne.NewPos((canvasSize.Width-size.Width)/2, (canvasSize.Height-size.Height)/2)
	r.tickerPopUp.Content.Move(pos)
	r.tickerPopUp.Content.Resize(size)

	innerPos := pos.Subtract(r.offset())
	r.bg.Move(innerPos)
	r.bg.Resize(size.Add(padding))
	r.LayoutShadow(r.tickerPopUp.innerSize, innerPos)
}

func (r *modalTickerPopUpRenderer) MinSize() fyne.Size {
	return r.tickerPopUp.Content.MinSize().Add(r.padding())
}

func (r *modalTickerPopUpRenderer) Refresh() {
	r.bg.FillColor = theme.BackgroundColor()
	if r.bg.Size() != r.tickerPopUp.innerSize {
		r.Layout(r.tickerPopUp.Size())
	}
}

func (r *modalTickerPopUpRenderer) BackgroundColor() color.Color {
	return theme.ShadowColor()
}
