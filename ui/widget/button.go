package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

type Button struct {
	baseWidget
	Style    ButtonStyle

	OnClicked func(*event.MouseEvent)

	label      *canvas.TextObject
	background *canvas.RectangleObject
}

type ButtonStyle int

const (
	DefaultButton ButtonStyle = iota
	PrimaryButton
)

func (b *Button) MinSize() ui.Size {
	return b.label.MinSize().Add(ui.NewSize(theme.Padding()*4, theme.Padding()*2))
}

func (b *Button) Layout() []ui.CanvasObject {
	if b.Style == PrimaryButton {
		b.background.Color = theme.PrimaryColor()
	}
	layout.NewMaxLayout().Layout(b.objects, b.Size)

	return b.objects
}

func NewButton(label string, clicked func(*event.MouseEvent)) *Button {
	text := canvas.NewText(label)
	bg := canvas.NewRectangle(theme.ButtonColor())

	return &Button{
		baseWidget{
			objects: []ui.CanvasObject{
				bg,
				text,
			},
		},
		DefaultButton,
		clicked,
		text,
		bg,
	}
}
