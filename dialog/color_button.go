package dialog

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	internalwidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*colorButton)(nil)
var _ desktop.Hoverable = (*colorButton)(nil)

// colorButton displays a color and triggers the callback when tapped.
type colorButton struct {
	widget.BaseWidget
	color   color.Color
	onTap   func(color.Color)
	hovered bool
}

// newColorButton creates a colorButton with the given color and callback.
func newColorButton(color color.Color, onTap func(color.Color)) *colorButton {
	b := &colorButton{
		color: color,
		onTap: onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *colorButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := newCheckeredBackground(false)
	rectangle := &canvas.Rectangle{
		FillColor: b.color,
	}
	return &colorButtonRenderer{
		BaseRenderer: internalwidget.NewBaseRenderer([]fyne.CanvasObject{background, rectangle}),
		button:       b,
		background:   background,
		rectangle:    rectangle,
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (b *colorButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (b *colorButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *colorButton) MouseMoved(*desktop.MouseEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (b *colorButton) MinSize() fyne.Size {
	return b.BaseWidget.MinSize()
}

// SetColor updates the color selected in this color widget
func (b *colorButton) SetColor(color color.Color) {
	if b.color == color {
		return
	}
	b.color = color
	b.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (b *colorButton) Tapped(*fyne.PointEvent) {
	if f := b.onTap; f != nil {
		f(b.color)
	}
}

type colorButtonRenderer struct {
	internalwidget.BaseRenderer
	button     *colorButton
	background *canvas.Raster
	rectangle  *canvas.Rectangle
}

func (r *colorButtonRenderer) Layout(size fyne.Size) {
	r.rectangle.Move(fyne.NewPos(0, 0))
	r.rectangle.Resize(size)
	r.background.Resize(size)
}

func (r *colorButtonRenderer) MinSize() fyne.Size {
	return r.rectangle.MinSize().Max(fyne.NewSize(32, 32))
}

func (r *colorButtonRenderer) Refresh() {
	if r.button.hovered {
		r.rectangle.StrokeColor = theme.HoverColor()
		r.rectangle.StrokeWidth = theme.Padding()
	} else {
		r.rectangle.StrokeWidth = 0
	}
	r.rectangle.FillColor = r.button.color
	r.background.Refresh()
	canvas.Refresh(r.button)
}
