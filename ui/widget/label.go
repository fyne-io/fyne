package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

type Label struct {
	Size     ui.Size
	Position ui.Position

	Text string

	objects []ui.CanvasObject
	label   *canvas.TextObject
}

func (l *Label) CurrentSize() ui.Size {
	return l.Size
}

func (l *Label) Resize(size ui.Size) {
	l.Size = size
}

func (l *Label) CurrentPosition() ui.Position {
	return l.Position
}

func (l *Label) Move(pos ui.Position) {
	l.Position = pos
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
		objects: []ui.CanvasObject{
			canvas.NewRectangle(theme.BackgroundColor()),
			obj,
		},
		Text:  text,
		label: obj,
	}
}
