package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Flexible defines an object with flexible properties.
type Flexible struct {
	BaseWidget
	child fyne.CanvasObject
	flex  int
}

// NewFlexible creates a new flexible widget.
func NewFlexible(flex int, w fyne.CanvasObject) *Flexible {
	child := w
	if w == nil {
		child = canvas.NewRectangle(color.Transparent)
	}
	f := &Flexible{flex: flex, child: child}
	f.ExtendBaseWidget(f)
	return f
}

// NewExpanded creates a new flexible with flex factor equals to 1.
func NewExpanded(w fyne.CanvasObject) *Flexible {
	return NewFlexible(1, w)
}

// DistanceToTextBaseline calculates the distance from the top of the widget until its text baseline.
func (f *Flexible) DistanceToTextBaseline() float32 {
	type baseliner interface{ DistanceToTextBaseline() float32 }
	distance := float32(0)
	if bl, ok := f.child.(baseliner); ok {
		distance = bl.DistanceToTextBaseline()
	}
	return distance
}

// Flex returns the flex factor.
func (f *Flexible) Flex() int {
	return f.flex
}

// ===============================================================
// Renderer
// ===============================================================

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (f *Flexible) CreateRenderer() fyne.WidgetRenderer {
	objects := []fyne.CanvasObject{f.child}
	return &flexibleRenderer{objects, f}
}

type flexibleRenderer struct {
	objects []fyne.CanvasObject
	f       *Flexible
}

func (r *flexibleRenderer) Destroy() {}

func (r *flexibleRenderer) MinSize() fyne.Size {
	return r.f.child.MinSize()
}

func (r *flexibleRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *flexibleRenderer) Layout(size fyne.Size) {
	r.f.child.Move(fyne.NewPos(0, 0))
	r.f.child.Resize(size)
}

func (r *flexibleRenderer) Refresh() {
	r.f.child.Refresh()
	canvas.Refresh(r.f.super())
}
