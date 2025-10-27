package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with Widget interface
var _ fyne.Widget = (*Clip)(nil)

// Clip describes a rectangular region  that will clip anything outside its bounds.
//
// Since: 2.7
type Clip struct {
	widget.BaseWidget
	Content fyne.CanvasObject
}

// NewClip returns a new rectangular clipping object.
//
// Since: 2.7
func NewClip(content fyne.CanvasObject) *Clip {
	return &Clip{Content: content}
}

func (c *Clip) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	return newClipRenderer(c)
}

// MinSize for a Clip simply returns Size{1, 1} as there is no
// explicit content
func (c *Clip) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return fyne.NewSize(1, 1)
}

type clipRenderer struct {
	c       *Clip
	objects []fyne.CanvasObject
}

func newClipRenderer(c *Clip) *clipRenderer {
	return &clipRenderer{c: c, objects: []fyne.CanvasObject{c.Content}}
}

func (r *clipRenderer) Destroy() {
}

func (r *clipRenderer) Layout(s fyne.Size) {
	o := r.objects[0]
	o.Resize(s.Max(o.MinSize()))
}

func (r *clipRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

func (r *clipRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *clipRenderer) Refresh() {
	r.objects[0] = r.c.Content
	r.Layout(r.c.Size())
	r.objects[0].Refresh()
}

// IsClip marks this widget as clipping. It is on the renderer to avoid a public API addition.
func (r *clipRenderer) IsClip() {}
