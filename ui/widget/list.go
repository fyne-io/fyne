package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/layout"

// List widget is a simple list where the child elements are arranged in a single column.
type List struct {
	baseWidget
}

// MinSize calculates the minimum size of a list.
// This is based on the contained children with a standard amount of padding added.
func (l *List) MinSize() ui.Size {
	return layout.NewGridLayout(1).MinSize(l.objects)
}

// Layout the components of the list widget
func (l *List) Layout(size ui.Size) []ui.CanvasObject {
	layout.NewGridLayout(1).Layout(l.objects, size)

	return l.objects
}

// Prepend inserts a new CanvasObject at the top of the list
func (l *List) Prepend(object ui.CanvasObject) {
	l.objects = append([]ui.CanvasObject{object}, l.objects...)
	ui.GetCanvas(l).Refresh(l)
}

// Append adds a new CanvasObject to the end of the list
func (l *List) Append(object ui.CanvasObject) {
	l.objects = append(l.objects, object)
	ui.GetCanvas(l).Refresh(l)
}

func (l *List) ApplyTheme() {
}

// NewList creates a new list widget with the specified list of child objects
func NewList(children ...ui.CanvasObject) *List {
	return &List{
		baseWidget{
			objects: children,
		},
	}
}
