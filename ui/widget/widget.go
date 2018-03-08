package widget

import "github.com/fyne-io/fyne/ui"

type Widget interface {
	CurrentSize() ui.Size
	Resize(ui.Size)
	CurrentPosition() ui.Position
	Move(ui.Position)

	MinSize() ui.Size

	// TODO should this move to a widget impl?... (private)
	Layout() []ui.CanvasObject
}


type baseWidget struct {
	Size     ui.Size
	Position ui.Position

	objects []ui.CanvasObject
}

func (w *baseWidget) CurrentSize() ui.Size {
	return w.Size
}

func (w *baseWidget) Resize(size ui.Size) {
	w.Size = size
}

func (w *baseWidget) CurrentPosition() ui.Position {
	return w.Position
}

func (w *baseWidget) Move(pos ui.Position) {
	w.Position = pos
}
