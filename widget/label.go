package widget

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/theme"

// Label widget is a basic text component with appropriate padding and layout.
type Label struct {
	baseWidget

	Text      string         // The content of the label
	Alignment fyne.TextAlign // The alignment of the Text

	label      *canvas.Text
	background *canvas.Rectangle
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (l *Label) MinSize() fyne.Size {
	return l.label.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

// SetText updates the text of the label widget
func (l *Label) SetText(text string) {
	l.Text = text
	l.label.Text = text
	fyne.GetCanvas(l).Refresh(l)
}

// ApplyTheme is called when the Label may need to update it's look
func (l *Label) ApplyTheme() {
	l.label.Color = theme.TextColor()
	l.background.FillColor = theme.BackgroundColor()
	l.label.Alignment = l.Alignment
}

// NewLabel creates a new layout widget with the set text content
func NewLabel(text string) *Label {
	obj := canvas.NewText(text, theme.TextColor())
	bg := canvas.NewRectangle(theme.BackgroundColor())

	l := &Label{
		baseWidget{
			objects: []fyne.CanvasObject{
				bg,
				obj,
			},
			layout: layout.NewMaxLayout(),
		},
		text,
		fyne.TextAlignLeading,
		obj,
		bg,
	}

	l.Layout(l.MinSize())
	return l
}
