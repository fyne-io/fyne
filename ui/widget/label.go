package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

type Label struct {
	baseWidget

	Text string

	label   *canvas.TextObject
}

func (l *Label) MinSize() ui.Size {
	return l.label.MinSize().Add(ui.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (l *Label) SetText(text string) {
	l.Text = text
	l.label.Text = text
}

func (l *Label) Layout() []ui.CanvasObject {
	layout.NewMaxLayout().Layout(l.objects, l.Size)

	return l.objects
}

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
