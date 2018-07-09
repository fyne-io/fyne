package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"

// List widget is a simple list where the child elements are arranged in a single column.
type List struct {
	baseWidget
}

// Prepend inserts a new CanvasObject at the top of the list
func (l *List) Prepend(object fyne.CanvasObject) {
	l.objects = append([]fyne.CanvasObject{object}, l.objects...)
	l.Layout(l.CurrentSize())

	fyne.GetCanvas(l).Refresh(l)
}

// Append adds a new CanvasObject to the end of the list
func (l *List) Append(object fyne.CanvasObject) {
	l.objects = append(l.objects, object)
	l.Layout(l.CurrentSize())

	fyne.GetCanvas(l).Refresh(l)
}

// ApplyTheme is called when the List may need to update it's look
func (l *List) ApplyTheme() {
}

// NewList creates a new list widget with the specified list of child objects
func NewList(children ...fyne.CanvasObject) *List {
	l := &List{
		baseWidget{
			objects: children,
			layout: layout.NewGridLayout(1),
		},
	}

	l.Layout(l.MinSize())
	return l
}
