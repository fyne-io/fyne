package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*SplitContainer)(nil)

// SplitContainer defines a container whose size is split between two children.
//
// Deprecated: use container.Split instead.
type SplitContainer struct {
	BaseWidget
	Offset     float64
	Horizontal bool
	Leading    fyne.CanvasObject
	Trailing   fyne.CanvasObject
}

// NewHSplitContainer creates a horizontally arranged container with the specified leading and trailing elements.
// A vertical split bar that can be dragged will be added between the elements.
//
// Deprecated: use container.NewHSplit instead.
func NewHSplitContainer(leading, trailing fyne.CanvasObject) *SplitContainer {
	return newSplitContainer(true, leading, trailing)
}

// NewVSplitContainer creates a vertically arranged container with the specified top and bottom elements.
// A horizontal split bar that can be dragged will be added between the elements.
//
// Deprecated: use container.NewVSplit instead.
func NewVSplitContainer(top, bottom fyne.CanvasObject) *SplitContainer {
	return newSplitContainer(false, top, bottom)
}

func newSplitContainer(horizontal bool, leading, trailing fyne.CanvasObject) *SplitContainer {
	s := &SplitContainer{
		Offset:     0.5, // Sensible default, can be overridden with SetOffset
		Horizontal: horizontal,
		Leading:    leading,
		Trailing:   trailing,
	}
	s.ExtendBaseWidget(s)
	return s
}

// SetOffset sets the offset (0.0 to 1.0) of the SplitContainer divider.
// 0.0 - Leading is min size, Trailing uses all remaining space.
// 0.5 - Leading & Trailing equally share the available space.
// 1.0 - Trailing is min size, Leading uses all remaining space.
func (s *SplitContainer) SetOffset(offset float64) {
	if s.Offset == offset {
		return
	}
	s.Offset = offset
	s.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *SplitContainer) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	d := newDivider(s)
	return &splitContainerRenderer{
		split:   s,
		divider: d,
		objects: []fyne.CanvasObject{s.Leading, d, s.Trailing},
	}
}

type splitContainerRenderer struct {
	split   *SplitContainer
	divider *divider
	objects []fyne.CanvasObject
}

func (r *splitContainerRenderer) Layout(size fyne.Size) {
	var dividerPos, leadingPos, trailingPos fyne.Position
	var dividerSize, leadingSize, trailingSize fyne.Size

	if r.split.Horizontal {
		lw, tw := r.computeSplitLengths(size.Width, r.split.Leading.MinSize().Width, r.split.Trailing.MinSize().Width)
		leadingPos.X = 0
		leadingSize.Width = lw
		leadingSize.Height = size.Height
		dividerPos.X = lw
		dividerSize.Width = dividerThickness()
		dividerSize.Height = size.Height
		trailingPos.X = lw + dividerSize.Width
		trailingSize.Width = tw
		trailingSize.Height = size.Height
	} else {
		lh, th := r.computeSplitLengths(size.Height, r.split.Leading.MinSize().Height, r.split.Trailing.MinSize().Height)
		leadingPos.Y = 0
		leadingSize.Width = size.Width
		leadingSize.Height = lh
		dividerPos.Y = lh
		dividerSize.Width = size.Width
		dividerSize.Height = dividerThickness()
		trailingPos.Y = lh + dividerSize.Height
		trailingSize.Width = size.Width
		trailingSize.Height = th
	}

	r.divider.Move(dividerPos)
	r.divider.Resize(dividerSize)
	r.split.Leading.Move(leadingPos)
	r.split.Leading.Resize(leadingSize)
	r.split.Trailing.Move(trailingPos)
	r.split.Trailing.Resize(trailingSize)
	canvas.Refresh(r.divider)
	canvas.Refresh(r.split.Leading)
	canvas.Refresh(r.split.Trailing)
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
	canvas.Refresh(r.split)
}

func (r *splitContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *splitContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *splitContainerRenderer) Destroy() {
}

func (r *splitContainerRenderer) computeSplitLengths(total, lMin, tMin int) (int, int) {
	available := float64(total - dividerThickness())
	ld := float64(lMin)
	tr := float64(tMin)
	offset := r.split.Offset

	min := ld / available
	max := 1 - tr/available
	if offset < min {
		offset = min
	}
	if offset > max {
		offset = max
	}

	ld = offset * available
	tr = available - ld
	return int(ld), int(tr)
}

// Declare conformity with interfaces
var _ fyne.CanvasObject = (*divider)(nil)
var _ fyne.Draggable = (*divider)(nil)
var _ desktop.Cursorable = (*divider)(nil)
var _ desktop.Hoverable = (*divider)(nil)

type divider struct {
	BaseWidget
	split   *SplitContainer
	hovered bool
}

func newDivider(split *SplitContainer) *divider {
	d := &divider{
		split: split,
	}
	d.ExtendBaseWidget(d)
	return d
}

func (d *divider) Cursor() desktop.Cursor {
	if d.split.Horizontal {
		return desktop.HResizeCursor
	}
	return desktop.VResizeCursor
}

func (d *divider) DragEnd() {
}

func (d *divider) Dragged(event *fyne.DragEvent) {
	offset := d.split.Offset
	if d.split.Horizontal {
		offset += float64(event.DraggedX) / float64(d.split.Size().Width)
	} else {
		offset += float64(event.DraggedY) / float64(d.split.Size().Height)
	}
	d.split.SetOffset(offset)
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
	return theme.ShadowColor()
}

func (r *dividerRenderer) Destroy() {
}

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
