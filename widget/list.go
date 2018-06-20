package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"

// List widget is a simple list where the child elements are arranged in a single column.
type List struct {
	baseWidget
}

// MinSize calculates the minimum size of a list.
// This is based on the contained children with a standard amount of padding added.
func (l *List) MinSize() fyne.Size {
	return layout.NewGridLayout(1).MinSize(l.objects)
}

// Layout the components of the list widget
func (l *List) Layout(size fyne.Size) []fyne.CanvasObject {
	layout.NewGridLayout(1).Layout(l.objects, size)

	return l.objects
}

// Prepend inserts a new CanvasObject at the top of the list
func (l *List) Prepend(object fyne.CanvasObject) {
	l.objects = append([]fyne.CanvasObject{object}, l.objects...)
	fyne.GetCanvas(l).Refresh(l)
}

// Append adds a new CanvasObject to the end of the list
func (l *List) Append(object fyne.CanvasObject) {
	l.objects = append(l.objects, object)
	fyne.GetCanvas(l).Refresh(l)
}

// ApplyTheme is called when the List may need to update it's look
func (l *List) ApplyTheme() {
}

// NewList creates a new list widget with the specified list of child objects
func NewList(children ...fyne.CanvasObject) *List {
	return &List{
		baseWidget{
			objects: children,
		},
	}
}
