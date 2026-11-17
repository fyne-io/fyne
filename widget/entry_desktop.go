//go:build !mobile && !android && !ios

package widget

import "fyne.io/fyne/v2"

// Dragging over text selects it on desktop. On mobile a drag is reserved for
// scrolling, so the content is not Draggable there and the gesture passes
// through to the entry's scroller or the nearest enclosing scroll container.
var _ fyne.Draggable = (*entryContent)(nil)

// DragEnd is called at end of a drag event.
func (e *entryContent) DragEnd() {
	// we need to propagate the focus, top level widget handles focus APIs
	e.entry.requestFocus()

	e.entry.DragEnd()
}

// Dragged is called when the pointer moves while a button is held down.
// It updates the selection accordingly.
func (e *entryContent) Dragged(d *fyne.DragEvent) {
	e.entry.Dragged(d)
}
