package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

// Box widget is a simple list where the child elements are arranged in a single column
// for vertical or a single row for horizontal arrangement.
// Deprecated: Use container.NewVBox() or container.NewHBox().
type Box struct {
	BaseWidget
	background color.Color

	Horizontal bool
	Children   []fyne.CanvasObject
}

// Refresh updates this box to match the current theme
func (b *Box) Refresh() {
	if b.background != nil {
		b.background = theme.BackgroundColor()
	}

	b.BaseWidget.Refresh()
}

// Prepend inserts a new CanvasObject at the top/left of the box
func (b *Box) Prepend(object fyne.CanvasObject) {
	b.Children = append([]fyne.CanvasObject{object}, b.Children...)

	b.Refresh()
}

// Append adds a new CanvasObject to the end/right of the box
func (b *Box) Append(object fyne.CanvasObject) {
	b.Children = append(b.Children, object)

	b.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (b *Box) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *Box) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	var lay fyne.Layout
	if b.Horizontal {
		lay = layout.NewHBoxLayout()
	} else {
		lay = layout.NewVBoxLayout()
	}

	return &boxRenderer{BaseRenderer: widget.NewBaseRenderer(b.Children), layout: lay, box: b}
}

// NewHBox creates a new horizontally aligned box widget with the specified list of child objects.
// Deprecated: Use container.NewHBox() instead.
func NewHBox(children ...fyne.CanvasObject) *Box {
	return &Box{BaseWidget: BaseWidget{}, Horizontal: true, Children: children}
}

// NewVBox creates a new vertically aligned box widget with the specified list of child objects.
// Deprecated: Use container.NewVBox() instead.
func NewVBox(children ...fyne.CanvasObject) *Box {
	return &Box{BaseWidget: BaseWidget{}, Horizontal: false, Children: children}
}

type boxRenderer struct {
	widget.BaseRenderer
	layout fyne.Layout
	box    *Box
}

func (b *boxRenderer) MinSize() fyne.Size {
	return b.layout.MinSize(b.Objects())
}

func (b *boxRenderer) Layout(size fyne.Size) {
	b.layout.Layout(b.Objects(), size)
}

func (b *boxRenderer) BackgroundColor() color.Color {
	if b.box.background == nil {
		return theme.BackgroundColor()
	}

	return b.box.background
}

func (b *boxRenderer) Refresh() {
	b.SetObjects(b.box.Children)
	for _, child := range b.Objects() {
		child.Refresh()
	}

	canvas.Refresh(b.box)
}
