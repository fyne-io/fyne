// Package widget defines the UI widgets within the Fyne toolkit
package widget

import "github.com/fyne-io/fyne"

// Widget defines the standard behaviours of any widget. This extends the
// CanvasObject - a widget behaves in the same basic way but will encapsulate
// many child objects to create the rendered widget.
type Widget interface {
	CurrentSize() fyne.Size
	Resize(fyne.Size)
	CurrentPosition() fyne.Position
	Move(fyne.Position)

	MinSize() fyne.Size
	ApplyTheme()
	CanvasObjects() []fyne.CanvasObject
}

// A base widget class to define the standard widget behaviours.
type baseWidget struct {
	Size     fyne.Size
	Position fyne.Position

	objects []fyne.CanvasObject
	layout  fyne.Layout
}

// Get the current size of this widget.
func (w *baseWidget) CurrentSize() fyne.Size {
	return w.Size
}

// Set a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) Resize(size fyne.Size) {
	w.Size = size

	w.Layout(size)
}

// Get the current position of this widget, relative to it's parent.
func (w *baseWidget) CurrentPosition() fyne.Position {
	return w.Position
}

// Move the widget to a new position, relative to it's parent.
// Note this hould not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) Move(pos fyne.Position) {
	w.Position = pos
}

func (w *baseWidget) MinSize() fyne.Size {
if w.layout == nil {return fyne.NewSize(1, 1)}
	return w.layout.MinSize(w.objects)
}

func (w *baseWidget) Layout(size fyne.Size) {
if w.layout == nil {return}
	w.layout.Layout(w.objects, size)
}

func (w *baseWidget) CanvasObjects() []fyne.CanvasObject {
	return w.objects
}
