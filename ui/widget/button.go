package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

type Button struct {
	Size     ui.Size
	Position ui.Position
	Style    ButtonStyle

	OnClicked func(*event.MouseEvent)

	objects    []ui.CanvasObject
	label      *canvas.TextObject
	background *canvas.RectangleObject
}

type ButtonStyle int

const (
	DefaultButton ButtonStyle = iota
	PrimaryButton
)

func (b *Button) CurrentSize() ui.Size {
	return b.Size
}

func (b *Button) Resize(size ui.Size) {
	b.Size = size
}

func (b *Button) CurrentPosition() ui.Position {
	return b.Position
}

func (b *Button) Move(pos ui.Position) {
	b.Position = pos
}

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
		OnClicked: clicked,
		objects: []ui.CanvasObject{
			bg,
			text,
		},
		label:      text,
		background: bg,
	}
}
