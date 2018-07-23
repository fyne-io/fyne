// Package widget defines the UI widgets within the Fyne toolkit
package widget

import "github.com/fyne-io/fyne"

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
	if w.layout == nil {
		return fyne.NewSize(1, 1)
	}
	return w.layout.MinSize(w.objects)
}

func (w *baseWidget) Layout(size fyne.Size) {
	if w.layout == nil {
		return
	}
	w.layout.Layout(w.objects, size)
}

// ApplyTheme is a fallback method that applies the new theme to all contained
// objects. Widgets that override this should consider doing similarly.
func (w *baseWidget) ApplyTheme() {
	for _, child := range w.objects {
		switch themed := child.(type) {
		case fyne.ThemedObject:
			themed.ApplyTheme()
		}
	}
}

func (w *baseWidget) CanvasObjects() []fyne.CanvasObject {
	return w.objects
}
