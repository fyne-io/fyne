package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*ColorButton)(nil)
var _ desktop.Hoverable = (*ColorButton)(nil)

// ColorButton displays a color and triggers the callback when tapped.
type ColorButton struct {
	BaseWidget
	color   color.Color
	minSize fyne.Size
	onTap   func(color.Color)
	hovered bool
}

// NewColorButton creates a ColorButton with the given color and callback.
func NewColorButton(color color.Color, onTap func(color.Color)) *ColorButton {
	b := &ColorButton{
		color: color,
		onTap: onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *ColorButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	background := newCheckeredBackground()
	rectangle := &canvas.Rectangle{
		FillColor: b.color,
	}
	return &colorButtonRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{background, rectangle}),
		button:       b,
		background:   background,
		rectangle:    rectangle,
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (b *ColorButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (b *ColorButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (b *ColorButton) MouseMoved(*desktop.MouseEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (b *ColorButton) MinSize() (min fyne.Size) {
	b.ExtendBaseWidget(b)
	b.propertyLock.RLock()
	min = b.minSize
	b.propertyLock.RUnlock()
	if min.IsZero() {
		min = b.BaseWidget.MinSize()
	}
	return
}

// SetMinSize specifies the smallest size this object should be
func (b *ColorButton) SetMinSize(size fyne.Size) {
	b.propertyLock.Lock()
	defer b.propertyLock.Unlock()

	b.minSize = size
}

// SetColor updates the color selected in this color widget
func (b *ColorButton) SetColor(color color.Color) {
	if b.color == color {
		return
	}
	b.color = color
	b.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (b *ColorButton) Tapped(*fyne.PointEvent) {
	writeRecentColor(colorToString(b.color))
	if f := b.onTap; f != nil {
		f(b.color)
	}
}

type colorButtonRenderer struct {
	widget.BaseRenderer
	button     *ColorButton
	background *canvas.Raster
	rectangle  *canvas.Rectangle
}

func (r *colorButtonRenderer) BackgroundColor() color.Color {
	if r.button.hovered {
		return theme.HoverColor()
	}
	return r.BaseRenderer.BackgroundColor()
}

func (r *colorButtonRenderer) Layout(size fyne.Size) {
	r.rectangle.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	r.rectangle.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height-2*theme.Padding()))
}

func (r *colorButtonRenderer) MinSize() (min fyne.Size) {
	min = r.rectangle.MinSize()
	size := theme.IconInlineSize()
	min = min.Max(fyne.NewSize(size, size))
	min = min.Add(fyne.NewSize(2*theme.Padding(), 2*theme.Padding()))
	return
}

func (r *colorButtonRenderer) Refresh() {
	r.rectangle.FillColor = r.button.color
	r.rectangle.Refresh()
	canvas.Refresh(r.button.super())
}
