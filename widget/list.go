package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"

// List widget is a simple list where the child elements are arranged in a single column.
type List struct {
	baseWidget

	Children []fyne.CanvasObject
}

// Prepend inserts a new CanvasObject at the top of the list
func (l *List) Prepend(object fyne.CanvasObject) {
	l.Children = append([]fyne.CanvasObject{object}, l.Children...)

	render := l.Renderer().(*listRenderer)
	render.objects = l.Children
	render.Layout(l.CurrentSize())

	fyne.GetCanvas(l).Refresh(l)
}

// Append adds a new CanvasObject to the end of the list
func (l *List) Append(object fyne.CanvasObject) {
	l.Children = append(l.Children, object)

	render := l.Renderer().(*listRenderer)

	render.Refresh()
}

func (l *List) createRenderer() fyne.WidgetRenderer {
	return &listRenderer{objects: l.Children, layout: layout.NewListLayout(), list: l}
}

// Renderer is a private method to Fyne which links this widget to it's renderer
func (l *List) Renderer() fyne.WidgetRenderer {
	if l.renderer == nil {
		l.renderer = l.createRenderer()
	}

	return l.renderer
}

// NewList creates a new list widget with the specified list of child objects
func NewList(children ...fyne.CanvasObject) *List {
	list := &List{baseWidget{}, children}

	list.Renderer().Layout(list.MinSize())
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
	for _, child := range l.objects {
		switch themed := child.(type) {
		case fyne.ThemedObject:
			themed.ApplyTheme()
		}
	}
}

func (l *listRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

func (l *listRenderer) Refresh() {
	l.objects = l.list.Children
	l.Layout(l.list.CurrentSize())

	fyne.GetCanvas(l.list).Refresh(l.list)
}
