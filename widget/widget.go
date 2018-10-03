// Package widget defines the UI widgets within the Fyne toolkit
package widget

import "github.com/fyne-io/fyne"

// A base widget class to define the standard widget behaviours.
type baseWidget struct {
	Size     fyne.Size
	Position fyne.Position

	renderer fyne.WidgetRenderer
}

// Get the current size of this widget.
func (w *baseWidget) CurrentSize() fyne.Size {
	return w.Size
}

// Set a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (w *baseWidget) Resize(size fyne.Size) {
	w.Size = size

	if w.renderer != nil {
		w.renderer.Layout(size)
	}
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
	if w.renderer == nil {
		return fyne.NewSize(0, 0)
	}
	return w.renderer.MinSize()
}

func (w *baseWidget) ApplyTheme() {
	if w.renderer == nil {
		return
	}
	w.renderer.ApplyTheme()
}
