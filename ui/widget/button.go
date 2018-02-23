package widget

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

type Button struct {
	Size     ui.Size
	Position ui.Position

	objects   []ui.CanvasObject
	label     *ui.TextObject
	onClicked func()
}

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
	return b.label.MinSize().Add(ui.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (b *Button) OnClicked() {
	if b.onClicked != nil {
		b.onClicked()
	}
}

func (b *Button) Layout() []ui.CanvasObject {
	layout.NewMaxLayout().Layout(b.objects, b.Size)

	return b.objects
}

func NewButton(label string, clicked func()) *Button {
	text := ui.NewText(label)
	button := &Button{
		onClicked: clicked,
		objects: []ui.CanvasObject{
			ui.NewRectangle(theme.ButtonColor()),
			text,
		},
		label: text,
	}

	return button
}
