//go:build !mobile && !android && !ios

package widget

import "fyne.io/fyne/v2"

var _ fyne.Draggable = (*selectable)(nil)

func (s *selectable) DragEnd() {
	if s.cursorColumn == s.selectColumn && s.cursorRow == s.selectRow {
		s.selecting = false
	}

	shouldRefresh := !s.selecting
	if shouldRefresh {
		s.Refresh()
	}
	s.selectEnded = true
}

func (s *selectable) Dragged(d *fyne.DragEvent) {
	s.dragSelect(d)
}
