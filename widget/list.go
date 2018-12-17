package widget

import (
	"image/color"
	"log"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/theme"
)

// List widget is a simple list where the child elements are arranged in a single column.
type List struct {
	baseWidget

	Children []fyne.CanvasObject
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *List) Resize(size fyne.Size) {
	l.resize(size, l)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *List) Move(pos fyne.Position) {
	l.move(pos, l)
}

// MinSize returns the smallest size this widget can shrink to
func (l *List) MinSize() fyne.Size {
	return l.minSize(l)
}

// Show this widget, if it was previously hidden
func (l *List) Show() {
	l.show(l)
}

// Hide this widget, if it was previously visible
func (l *List) Hide() {
	l.hide(l)
}

// Prepend inserts a new CanvasObject at the top of the list
func (l *List) Prepend(object fyne.CanvasObject) {
	l.Children = append([]fyne.CanvasObject{object}, l.Children...)

	Renderer(l).Refresh()
}

// Append adds a new CanvasObject to the end of the list
func (l *List) Append(object fyne.CanvasObject) {
	l.Children = append(l.Children, object)

	Renderer(l).Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (l *List) CreateRenderer() fyne.WidgetRenderer {
	log.Println("Deprecated: The List widget is replaced by VBox")
	return &listRenderer{objects: l.Children, layout: layout.NewVBoxLayout(), list: l}
}

// NewList creates a new list widget with the specified list of child objects
// Deprecated: NewList has been replaced with NewVBox
func NewList(children ...fyne.CanvasObject) *List {
	list := &List{baseWidget{}, children}

	Renderer(list).Layout(list.MinSize())
	return list
}

type listRenderer struct {
	layout fyne.Layout

	objects []fyne.CanvasObject
	list    *List
}

func (l *listRenderer) MinSize() fyne.Size {
	return l.layout.MinSize(l.objects)
}

func (l *listRenderer) Layout(size fyne.Size) {
	l.layout.Layout(l.objects, size)
}

// ApplyTheme is a fallback method that applies the new theme to all contained
// objects. Widgets that override this should consider doing similarly.
func (l *listRenderer) ApplyTheme() {
	l.Refresh()
}

func (l *listRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (l *listRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

func (l *listRenderer) Refresh() {
	l.objects = l.list.Children
	l.Layout(l.list.Size())

	canvas.Refresh(l.list)
}
