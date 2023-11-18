package widget

import (
	"fyne.io/fyne/v2"
)

// Container widget helps to use a container as a widget.
type ContainerWidget struct {
	BaseWidget
	*fyne.Container
}

type ContainerWidgetRenderer struct {
	*fyne.Container
}

func (r *ContainerWidgetRenderer) Layout(size fyne.Size) {
	r.Container.Layout.Layout(r.Container.Objects, size)
}
func (r *ContainerWidgetRenderer) MinSize() fyne.Size {
	return r.Container.MinSize()
}
func (r *ContainerWidgetRenderer) Refresh() {
	r.Container.Refresh()
}
func (r *ContainerWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.Container.Objects
}
func (r *ContainerWidgetRenderer) Destroy() {}

var _ fyne.CanvasObject = (*ContainerWidget)(nil)
var _ fyne.WidgetRenderer = (*ContainerWidgetRenderer)(nil)

// implement CanvasObject interface
func (t *ContainerWidget) CreateRenderer() fyne.WidgetRenderer {
	return &ContainerWidgetRenderer{
		Container: t.Container,
	}
}

func (t *ContainerWidget) Hide() {
	t.Container.Hide()
}

func (t *ContainerWidget) MinSize() fyne.Size {
	return t.Container.MinSize()
}

func (t *ContainerWidget) Move(pos fyne.Position) {
	t.Container.Move(pos)
}

func (t *ContainerWidget) Position() fyne.Position {
	return t.Container.Position()
}

func (t *ContainerWidget) Refresh() {
	t.Container.Refresh()
}

func (t *ContainerWidget) Resize(size fyne.Size) {
	t.Container.Resize(size)
}

func (t *ContainerWidget) Show() {
	t.Container.Show()
}

func (t *ContainerWidget) Size() fyne.Size {
	return t.Container.Size()
}

func (t *ContainerWidget) Visible() bool {
	return t.Container.Visible()
}
