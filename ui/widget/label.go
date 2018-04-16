package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

// Label widget is a basic text component with appropriate padding and layout.
type Label struct {
	baseWidget

	Text string

	label *canvas.Text
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (l *Label) MinSize() ui.Size {
	return l.label.MinSize().Add(ui.NewSize(theme.Padding()*2, theme.Padding()*2))
}

// SetText updates the text of the label widget
func (l *Label) SetText(text string) {
	l.Text = text
	l.label.Text = text
}

// Layout the components of the label widget
func (l *Label) Layout(size ui.Size) []ui.CanvasObject {
	layout.NewMaxLayout().Layout(l.objects, size)

	return l.objects
}

// NewLabel creates a new layout widget with the set text content
func NewLabel(text string) *Label {
	obj := canvas.NewText(text)
	return &Label{
		baseWidget{
			objects: []ui.CanvasObject{
				canvas.NewRectangle(theme.BackgroundColor()),
				obj,
			},
		},
		text,
		obj,
	}
}
