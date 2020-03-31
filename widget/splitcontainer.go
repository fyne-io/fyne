package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

func dividerThickness() int {
	return theme.Padding() * 2
}

func dividerLength() int {
	return theme.Padding() * 6
}

func handleThickness() int {
	return theme.Padding() / 2
}

func handleLength() int {
	return theme.Padding() * 4
}

// Declare conformity with interfaces
var _ fyne.CanvasObject = (*divider)(nil)
var _ fyne.Draggable = (*divider)(nil)
var _ desktop.Hoverable = (*divider)(nil)

type divider struct {
	BaseWidget
	split   *SplitContainer
	hovered bool
}

func (d *divider) DragEnd() {
}

func (d *divider) Dragged(event *fyne.DragEvent) {
	offset := d.split.offset
	if d.split.Horizontal {
		offset += event.DraggedX
	} else {
		offset += event.DraggedY
	}
	d.split.updateOffset(offset)
}

func (d *divider) MouseIn(event *desktop.MouseEvent) {
	d.hovered = true
	d.split.Refresh()
}

func (d *divider) MouseOut() {
	d.hovered = false
	d.split.Refresh()
}

func (d *divider) MouseMoved(event *desktop.MouseEvent) {}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (d *divider) CreateRenderer() fyne.WidgetRenderer {
	d.ExtendBaseWidget(d)
	r := canvas.NewRectangle(theme.IconColor())
	return &dividerRenderer{
		divider:   d,
		rectangle: r,
		objects:   []fyne.CanvasObject{r},
	}
}

func newDivider(split *SplitContainer) *divider {
	d := &divider{
		split: split,
	}
	d.ExtendBaseWidget(d)
	return d
}

type dividerRenderer struct {
	divider   *divider
	rectangle *canvas.Rectangle
	objects   []fyne.CanvasObject
}

func (r *dividerRenderer) Layout(size fyne.Size) {
	var x, y, w, h int
	if r.divider.split.Horizontal {
		x = (dividerThickness() - handleThickness()) / 2
		y = (size.Height - handleLength()) / 2
		w = handleThickness()
		h = handleLength()
	} else {
		x = (size.Width - handleLength()) / 2
		y = (dividerThickness() - handleThickness()) / 2
		w = handleLength()
		h = handleThickness()
	}
	r.rectangle.Move(fyne.NewPos(x, y))
	r.rectangle.Resize(fyne.NewSize(w, h))
}

func (r *dividerRenderer) MinSize() fyne.Size {
	if r.divider.split.Horizontal {
		return fyne.NewSize(dividerThickness(), dividerLength())
	}
	return fyne.NewSize(dividerLength(), dividerThickness())
}

func (r *dividerRenderer) Refresh() {
	r.rectangle.FillColor = theme.IconColor()
	r.Layout(r.divider.Size())
}

func (r *dividerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *dividerRenderer) BackgroundColor() color.Color {
	if r.divider.hovered {
		return theme.HoverColor()
	}
	return theme.ButtonColor()
}

func (r *dividerRenderer) Destroy() {
}

type splitContainerRenderer struct {
	split   *SplitContainer
	divider *divider
	objects []fyne.CanvasObject
}

func (r *splitContainerRenderer) Layout(size fyne.Size) {
	var dividerPos, childAPos, childBPos fyne.Position
	var dividerSize, childASize, childBSize fyne.Size
	if r.split.Horizontal {
		x := 0
		half := (size.Width - dividerThickness()) / 2
		childAPos.X = x
		childASize.Width = half + r.split.offset
		childASize.Height = size.Height
		x += childASize.Width
		dividerPos.X = x
		dividerSize.Width = dividerThickness()
		dividerSize.Height = size.Height
		x += dividerSize.Width
		childBPos.X = x
		childBSize.Width = half - r.split.offset
		childBSize.Height = size.Height
	} else {
		y := 0
		half := (size.Height - dividerThickness()) / 2
		childAPos.Y = y
		childASize.Width = size.Width
		childASize.Height = half + r.split.offset
		y += childASize.Height
		dividerPos.Y = y
		dividerSize.Width = size.Width
		dividerSize.Height = dividerThickness()
		y += dividerSize.Height
		childBPos.Y = y
		childBSize.Width = size.Width
		childBSize.Height = half - r.split.offset
	}
	r.divider.Move(dividerPos)
	r.divider.Resize(dividerSize)
	r.split.ChildA.Move(childAPos)
	r.split.ChildA.Resize(childASize)
	r.split.ChildB.Move(childBPos)
	r.split.ChildB.Resize(childBSize)
	canvas.Refresh(r.divider)
	canvas.Refresh(r.split.ChildA)
	canvas.Refresh(r.split.ChildB)
}

func (r *splitContainerRenderer) MinSize() fyne.Size {
	s := fyne.NewSize(0, 0)
	for _, o := range r.objects {
		min := o.MinSize()
		if r.split.Horizontal {
			s.Width += min.Width
			s.Height = fyne.Max(s.Height, min.Height)
		} else {
			s.Width = fyne.Max(s.Width, min.Width)
			s.Height += min.Height
		}
	}
	return s
}

func (r *splitContainerRenderer) Refresh() {
	r.Layout(r.split.Size())
}

func (r *splitContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *splitContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *splitContainerRenderer) Destroy() {
}

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*SplitContainer)(nil)

// SplitContainer defines a container whose size is split between two children.
type SplitContainer struct {
	BaseWidget
	Horizontal     bool
	ChildA, ChildB fyne.CanvasObject
	offset         int // Adjusts how the size is split between the children, positive favours the first while negative favours the second.
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *SplitContainer) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	d := newDivider(s)
	return &splitContainerRenderer{
		split:   s,
		divider: d,
		objects: []fyne.CanvasObject{s.ChildA, d, s.ChildB},
	}
}

func (s *SplitContainer) updateOffset(offset int) {
	var positiveLimit int
	var negativeLimit int
	if s.Horizontal {
		half := (s.size.Width - dividerThickness()) / 2
		positiveLimit = half - s.ChildB.MinSize().Width
		negativeLimit = s.ChildA.MinSize().Width - half
	} else {
		half := (s.size.Height - dividerThickness()) / 2
		positiveLimit = half - s.ChildB.MinSize().Height
		negativeLimit = s.ChildA.MinSize().Height - half
	}
	if offset < negativeLimit {
		offset = negativeLimit
	}
	if offset > positiveLimit {
		offset = positiveLimit
	}
	s.offset = offset
	s.Refresh()
}

// NewHorizontalSplitContainer create a splitable parent wrapping the specified children.
func NewHorizontalSplitContainer(left, right fyne.CanvasObject) *SplitContainer {
	return newSplitContainer(true, left, right)
}

// NewVerticalSplitContainer create a splitable parent wrapping the specified children.
func NewVerticalSplitContainer(top, bottom fyne.CanvasObject) *SplitContainer {
	return newSplitContainer(false, top, bottom)
}

func newSplitContainer(horizontal bool, a, b fyne.CanvasObject) *SplitContainer {
	s := &SplitContainer{
		Horizontal: horizontal,
		ChildA:     a,
		ChildB:     b,
	}
	s.ExtendBaseWidget(s)
	return s
}
