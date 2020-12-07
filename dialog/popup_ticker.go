package dialog

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/internal/widget"
	commonwidget "fyne.io/fyne/widget"
)

var _ fyne.Draggable = (*TickerPopUp)(nil)

type TextWidget interface {
	GetText() string
	SetText(string)
	SetSelected(string)
	TextStyle() fyne.TextStyle
}

type ringBuffer struct {
	data []rune
	start int // Start of ringBuffer can be anywhere in array.
	bufferWidth int
	width int                     // Width of draw area
	SeparatorRune rune            // Separator character
	entryPadding int              // Any padding around entries
	labelFontSize int             // Font size in label
	labelTextStyle fyne.TextStyle // font text style
	forward bool                  // Dragg direction
}

type PopupTickerListener interface {
	TapCallback(fyne.Tappable, *fyne.PointEvent)
}

func (rb *ringBuffer) Init(start int, data []rune) {
	rb.data = data
	rb.start = start
}

// Turn - rotates the ringbuffer by appropriate offset.  -offset is left, +offset is right.
func (rb *ringBuffer) Turn(offset int) {
	if len(rb.data) == 0 {
		return
	}
	if offset > 0 {
		offset = 1
	} else if offset < 0 {
		offset = -1
	}

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
	rb.Turn(position)
}

// GetSelected -- given pixel offset, returns selected text.
func (rb *ringBuffer) GetSelected(popupTickerPosX int, selectedPosX int) string {
	currentData := rb.Data(false)
	if len(currentData) == 0 {
		return ""
	}

	// Seek the offset by character widths.
	width := popupTickerPosX
	nearestIndex := 0
	for i := 0; i < len(currentData); i++ {
		if (currentData[i] == rb.SeparatorRune) {
			width = width + rb.entryPadding
		}
		charWidth := fyne.MeasureText(string(currentData[i]), rb.labelFontSize, rb.labelTextStyle).Width
		width = width + charWidth
		if width >= selectedPosX {
			nearestIndex = i
			break
		}
		if (currentData[i] == rb.SeparatorRune) {
			width = width + rb.entryPadding
		}
	}

	if currentData[nearestIndex] == rb.SeparatorRune {
		// Don't start on a separator.
		nearestIndex = nearestIndex + 1
	}

	nearestSeparatorFound := false
	for i := nearestIndex; i >= 0; i-- {
		if currentData[i] == rb.SeparatorRune {
			nearestSeparatorFound = true
			break
		}
		nearestIndex = i
	}

	endIndex := nearestIndex
	farthestSeparatorFound := false

	for i := nearestIndex; i < len(currentData); i++ {
		if currentData[i] == rb.SeparatorRune {
			farthestSeparatorFound = true
			endIndex = i
			break
		}
	}
	var result string

	if !nearestSeparatorFound {
		fullData := rb.Data(true)

		for i := len(fullData) - 1; i >= 0; i-- {
			if fullData[i] == rb.SeparatorRune {
				nearestSeparatorFound = true
				break
			}
			nearestIndex = i
		}
		if nearestIndex > endIndex {
			result = string(fullData[nearestIndex:]) + string(fullData[0:endIndex])
		} else {
			result = string(rb.data[nearestIndex:endIndex])
		}
	} else if !farthestSeparatorFound {
		fullData := rb.Data(true)

		for i := nearestIndex; i < len(fullData); i++ {
			if fullData[i] == rb.SeparatorRune {
				farthestSeparatorFound = true
				endIndex = i
				break
			}
		}
		result = string(fullData[nearestIndex:endIndex])
	} else {
		result = string(currentData[nearestIndex:endIndex])
	}

	return result
}

// Data - returns current data at current turn, read circularly
func (rb *ringBuffer) Data(complete bool) []rune {
	var data []rune
	if rb.start == 0 {
		data = rb.data
	} else {
		data = append(rb.data[rb.start:], rb.data[0: rb.start]...)
	}

	if !complete && rb.width > 0 {
		width := 0
		boundIndex := 0
		for i := 0; i < len(data); i++ {
			if (data[i] == rb.SeparatorRune) {
				width = width + rb.entryPadding
			}
			charWidth := fyne.MeasureText(string(data[i]), rb.labelFontSize, rb.labelTextStyle).Width
			width = width + charWidth
			if width >= rb.width {
				boundIndex = i
				break
			}
			if (data[i] == rb.SeparatorRune) {
				width = width + rb.entryPadding
			}
		}

		return data[0:boundIndex]
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
	commonwidget.BaseWidget

	Content fyne.CanvasObject
	Canvas  fyne.Canvas
	CurrentSelection string

	popupTickerListener  PopupTickerListener
	rb      ringBuffer // backing the content with a ringBuffer.
	innerPos     fyne.Position
	innerSize    fyne.Size
	modal        bool
	overlayShown bool
	draggedX     fyne.Position
	dragging     int
	dragScale    int
	dsCount      int
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
	p.dragging = 1
}

// Tapped is called when the user taps the tickerPopUp background - if not modal then dismiss this widget
func (p *TickerPopUp) Tapped(e *fyne.PointEvent) {
	if e.AbsolutePosition.X < p.innerPos.X || e.AbsolutePosition.Y < p.innerPos.Y || e.AbsolutePosition.X > (p.innerPos.X + p.innerSize.Width) || e.AbsolutePosition.Y > (p.innerPos.Y + p.innerSize.Height) {
		p.Hide()
		return
	}
	if p.dragging > 0 {
		p.dragging = 0
		return
	}
	if p.popupTickerListener != nil {
		p.GetSelected(&e.AbsolutePosition)
		p.popupTickerListener.TapCallback(p, e)
	}
}

func (p *TickerPopUp) endOffset() int {
	return p.innerPos.X + theme.Padding()
}

func (p *TickerPopUp) getRatio(pos *fyne.Position) float64 {
	x := pos.X - p.innerPos.X

	tickerWidth := p.rb.width

	if x > p.innerPos.X + tickerWidth {
		return 1.0
	} else if pos.X < p.innerPos.X {
		return 0.0
	} else {
		return float64(x) / float64(tickerWidth - (2 * theme.Padding()))
	}
}

func (p *TickerPopUp) GetSelected(pos *fyne.Position) string {
	p.CurrentSelection = p.rb.GetSelected(p.Content.Position().X + theme.Padding(), pos.X)
	return p.CurrentSelection
}

func (p *TickerPopUp) SetText(data []rune) {
	p.rb.data = data
	switch p.Content.(type) {
	case *commonwidget.Label:
		p.Content.(*commonwidget.Label).Text = string(p.rb.Data(false))
		break
	case *commonwidget.Entry:
		p.Content.(*commonwidget.Entry).Text = string(p.rb.Data(false))
		break
	case TextWidget:
		p.Content.(TextWidget).SetText(string(p.rb.Data(false)))
	}
}

func (p *TickerPopUp) Dragged(e *fyne.DragEvent) {
	p.dragging = 2
	if p.draggedX.X == 0 {
		p.draggedX = e.Position
		return
	}

	if e.AbsolutePosition.X < p.innerPos.X || e.AbsolutePosition.Y < p.innerPos.Y || e.AbsolutePosition.X > (p.innerPos.X + p.innerSize.Width) || e.AbsolutePosition.Y > (p.innerPos.Y + p.innerSize.Height) {
		p.Hide()
		return
	}

	var diffPosition fyne.Position
	diffPosition.X = e.Position.X - p.draggedX.X
	diffPosition.Y = e.Position.Y- p.draggedX.Y
	p.draggedX = e.Position

	if diffPosition.X <= p.dragScale && diffPosition.X >= -p.dragScale {
		p.dsCount = (p.dsCount + 1) % p.dragScale
	} else {
		p.dsCount = 0
	}

	if p.dsCount != 0 {
		p.rb.Seek(int(diffPosition.X))
		switch p.Content.(type) {
		case *commonwidget.Label:
			p.Content.(*commonwidget.Label).Text = string(p.rb.Data(false))
			break
		case *commonwidget.Entry:
			p.Content.(*commonwidget.Entry).Text = string(p.rb.Data(false))
			break
		case TextWidget:
			p.Content.(TextWidget).SetText(string(p.rb.Data(false)))
		}
		p.Content.Refresh()
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
func (p *TickerPopUp) GetSeparatorRune() rune {
	return p.rb.SeparatorRune
}
// NewTickerPopUpAtPosition creates a new tickerPopUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
//
// Deprecated: Use ShowTickerPopUpAtPosition() instead.
func NewTickerPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, popupTickerListener PopupTickerListener, pos fyne.Position, size fyne.Size, fontSize int, separator rune, entryPadding int) *TickerPopUp {
	p := newTickerPopUp(content, canvas, popupTickerListener, size, fontSize, separator, entryPadding)
	p.ShowAtPosition(pos)
	return p
}

// ShowTickerPopUpAtPosition creates a new tickerPopUp for the specified content at the specified absolute position.
// It will then display the popup on the passed canvas.
func ShowTickerPopUpAtPosition(content fyne.CanvasObject, canvas fyne.Canvas, pos fyne.Position, popupTickerListener PopupTickerListener, size fyne.Size, fontSize int, separator rune, entryPadding int) {
	newTickerPopUp(content, canvas, popupTickerListener, size, fontSize, separator, entryPadding).ShowAtPosition(pos)
}

func newTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas, popupTickerListener PopupTickerListener, size fyne.Size, fontSize int, separator rune, entryPadding int) *TickerPopUp {
	rb := ringBuffer{ start: 0, labelFontSize: fontSize, width: size.Width, forward: true, SeparatorRune: separator, entryPadding: entryPadding }

	// TODO: would be nice if Label and Entry implemented GetText() and TextStyle().  Then remove switch and access directly.
	switch content.(type) {
	case *commonwidget.Label:
		rb.data = []rune(content.(*commonwidget.Label).Text)
		rb.labelTextStyle = content.(*commonwidget.Label).TextStyle
		content.(*commonwidget.Label).Text = string(rb.Data(false))
		break
	case *commonwidget.Entry:
		rb.data = []rune(content.(*commonwidget.Entry).Text)
		rb.labelTextStyle = fyne.TextStyle{} // Entry should provide.
		content.(*commonwidget.Entry).Text = string(rb.Data(false))
		break
	case TextWidget:
		rb.data = []rune(content.(TextWidget).GetText())
		rb.labelTextStyle = content.(TextWidget).TextStyle()
		content.(TextWidget).SetText(string(rb.Data(false)))
	}

	ret := &TickerPopUp{Content: content, rb: rb, Canvas: canvas, popupTickerListener: popupTickerListener, modal: false, dragScale: 5}
	ret.ExtendBaseWidget(ret)
	ret.Resize(size)
	return ret
}

// NewTickerPopUp creates a new tickerPopUp for the specified content and displays it on the passed canvas.
//
// Deprecated: This will no longer show the pop-up in 2.0. Use ShowTickerPopUp() instead.
func NewTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas, popupTickerListener PopupTickerListener, size fyne.Size, fontSize int, separator rune, entryPadding int) *TickerPopUp {
	return NewTickerPopUpAtPosition(content, canvas, popupTickerListener, fyne.NewPos(0, 0), size, fontSize, separator, entryPadding)
}

// ShowTickerPopUp creates a new tickerPopUp for the specified content and displays it on the passed canvas.
func ShowTickerPopUp(content fyne.CanvasObject, canvas fyne.Canvas, popupTickerListener PopupTickerListener, size fyne.Size, fontSize int, separator rune, entryPadding int) {
	newTickerPopUp(content, canvas, popupTickerListener, size, fontSize, separator, entryPadding).Show()
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
