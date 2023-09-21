package container

import (
	"fyne.io/fyne/v2"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/widget"
)

// MultipleWindows is a container that handles multiple `InnerWindow` containers.
// Each inner window can be dragged, resized and the stacking will change when the title bar is tapped.
//
// Since: 2.5
type MultipleWindows struct {
	widget.BaseWidget

	Windows []*InnerWindow

	content *fyne.Container
}

// NewMultipleWindows creates a new `MultipleWindows` container to manage many inner windows.
// The initial window list is passed optionally to this constructor function.
// You can add new more windows to this container by calling `Add` or updating the `Windows`
// field and calling `Refresh`.
//
// Since: 2.5
func NewMultipleWindows(wins ...*InnerWindow) *MultipleWindows {
	m := &MultipleWindows{Windows: wins}
	m.ExtendBaseWidget(m)
	return m
}

func (m *MultipleWindows) Add(w *InnerWindow) {
	m.Windows = append(m.Windows, w)
	m.refreshChildren()
}

func (m *MultipleWindows) CreateRenderer() fyne.WidgetRenderer {
	m.content = New(&multiWinLayout{})
	m.refreshChildren()
	return widget.NewSimpleRenderer(intWidget.NewScroll(m.content))
}

func (m *MultipleWindows) Refresh() {
	m.refreshChildren()
	//	m.BaseWidget.Refresh()
}

func (m *MultipleWindows) raise(w *InnerWindow) {
	id := -1
	for i, ww := range m.Windows {
		if ww == w {
			id = i
			break
		}
	}
	if id == -1 {
		return
	}

	windows := append(m.Windows[:id], m.Windows[id+1:]...)
	m.Windows = append(windows, w)
	m.refreshChildren()
}

func (m *MultipleWindows) refreshChildren() {
	if m.content == nil {
		return
	}

	objs := make([]fyne.CanvasObject, len(m.Windows))
	for i, w := range m.Windows {
		objs[i] = w

		m.setupChild(w)
	}
	m.content.Objects = objs
	m.content.Refresh()
}

func (m *MultipleWindows) setupChild(w *InnerWindow) {
	w.OnDragged = func(ev *fyne.DragEvent) {
		w.Move(w.Position().Add(ev.Dragged))
	}
	w.OnResized = func(ev *fyne.DragEvent) {
		size := w.Size().Add(ev.Dragged)
		w.Resize(size.Max(w.MinSize()))
	}
	w.OnTappedBar = func() {
		m.raise(w)
	}
}

type multiWinLayout struct {
}

func (m *multiWinLayout) Layout(objects []fyne.CanvasObject, _ fyne.Size) {
	for _, w := range objects { // update the windows so they have real size
		w.Resize(w.MinSize().Max(w.Size()))
	}
}

func (m *multiWinLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{}
}
