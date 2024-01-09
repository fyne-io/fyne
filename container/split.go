package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Split)(nil)

// Split defines a container whose size is split between two children.
//
// Since: 1.4
type Split struct {
	widget.BaseWidget
	Offset     float64
	Horizontal bool
	Leading    fyne.CanvasObject
	Trailing   fyne.CanvasObject
}

// NewHSplit creates a horizontally arranged container with the specified leading and trailing elements.
// A vertical split bar that can be dragged will be added between the elements.
//
// Since: 1.4
func NewHSplit(leading, trailing fyne.CanvasObject) *Split {
	return newSplitContainer(true, leading, trailing)
}

// NewVSplit creates a vertically arranged container with the specified top and bottom elements.
// A horizontal split bar that can be dragged will be added between the elements.
//
// Since: 1.4
func NewVSplit(top, bottom fyne.CanvasObject) *Split {
	return newSplitContainer(false, top, bottom)
}

func newSplitContainer(horizontal bool, leading, trailing fyne.CanvasObject) *Split {
	s := &Split{
		Offset:     0.5, // Sensible default, can be overridden with SetOffset
		Horizontal: horizontal,
		Leading:    leading,
		Trailing:   trailing,
	}
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Split) CreateRenderer() fyne.WidgetRenderer {
	s.BaseWidget.ExtendBaseWidget(s)
	d := newDivider(s)
	return &splitContainerRenderer{
		split:   s,
		divider: d,
		objects: []fyne.CanvasObject{s.Leading, d, s.Trailing},
	}
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
//
// Deprecated: Support for extending containers is being removed
func (s *Split) ExtendBaseWidget(wid fyne.Widget) {
	s.BaseWidget.ExtendBaseWidget(wid)
}

// SetOffset sets the offset (0.0 to 1.0) of the Split divider.
// 0.0 - Leading is min size, Trailing uses all remaining space.
// 0.5 - Leading & Trailing equally share the available space.
// 1.0 - Trailing is min size, Leading uses all remaining space.
func (s *Split) SetOffset(offset float64) {
	if s.Offset == offset {
		return
	}
	s.Offset = offset
	s.Refresh()
}

var _ fyne.WidgetRenderer = (*splitContainerRenderer)(nil)

type splitContainerRenderer struct {
	split   *Split
	divider *divider
	objects []fyne.CanvasObject
}

func (r *splitContainerRenderer) Destroy() {
}

func (r *splitContainerRenderer) Layout(size fyne.Size) {
	var dividerPos, leadingPos, trailingPos fyne.Position
	var dividerSize, leadingSize, trailingSize fyne.Size

	if r.split.Horizontal {
		lw, tw := r.computeSplitLengths(size.Width, r.minLeadingWidth(), r.minTrailingWidth())
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
		lh, th := r.computeSplitLengths(size.Height, r.minLeadingHeight(), r.minTrailingHeight())
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

func (r *splitContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *splitContainerRenderer) Refresh() {
	r.objects[0] = r.split.Leading
	// [1] is divider which doesn't change
	r.objects[2] = r.split.Trailing
	r.Layout(r.split.Size())
	canvas.Refresh(r.split)
}

func (r *splitContainerRenderer) computeSplitLengths(total, lMin, tMin float32) (float32, float32) {
	available := float64(total - dividerThickness())
	if available <= 0 {
		return 0, 0
	}
	ld := float64(lMin)
	tr := float64(tMin)
	offset := r.split.Offset

	min := ld / available
	max := 1 - tr/available
	if min <= max {
		if offset < min {
			offset = min
		}
		if offset > max {
			offset = max
		}
	} else {
		offset = ld / (ld + tr)
	}

	ld = offset * available
	tr = available - ld
	return float32(ld), float32(tr)
}

func (r *splitContainerRenderer) minLeadingWidth() float32 {
	if r.split.Leading.Visible() {
		return r.split.Leading.MinSize().Width
	}
	return 0
}

func (r *splitContainerRenderer) minLeadingHeight() float32 {
	if r.split.Leading.Visible() {
		return r.split.Leading.MinSize().Height
	}
	return 0
}

func (r *splitContainerRenderer) minTrailingWidth() float32 {
	if r.split.Trailing.Visible() {
		return r.split.Trailing.MinSize().Width
	}
	return 0
}

func (r *splitContainerRenderer) minTrailingHeight() float32 {
	if r.split.Trailing.Visible() {
		return r.split.Trailing.MinSize().Height
	}
	return 0
}

// Declare conformity with interfaces
var _ fyne.CanvasObject = (*divider)(nil)
var _ fyne.Draggable = (*divider)(nil)
var _ desktop.Cursorable = (*divider)(nil)
var _ desktop.Hoverable = (*divider)(nil)

type divider struct {
	widget.BaseWidget
	split          *Split
	hovered        bool
	startDragOff   *fyne.Position
	currentDragPos fyne.Position
}

func newDivider(split *Split) *divider {
	d := &divider{
		split: split,
	}
	d.ExtendBaseWidget(d)
	return d
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (d *divider) CreateRenderer() fyne.WidgetRenderer {
	d.ExtendBaseWidget(d)
	background := canvas.NewRectangle(theme.ShadowColor())
	foreground := canvas.NewRectangle(theme.ForegroundColor())
	return &dividerRenderer{
		divider:    d,
		background: background,
		foreground: foreground,
		objects:    []fyne.CanvasObject{background, foreground},
	}
}

func (d *divider) Cursor() desktop.Cursor {
	if d.split.Horizontal {
		return desktop.HResizeCursor
	}
	return desktop.VResizeCursor
}

func (d *divider) DragEnd() {
	d.startDragOff = nil
}

func (d *divider) Dragged(e *fyne.DragEvent) {
	if d.startDragOff == nil {
		d.currentDragPos = d.Position().Add(e.Position)
		start := e.Position.Subtract(e.Dragged)
		d.startDragOff = &start
	} else {
		d.currentDragPos = d.currentDragPos.Add(e.Dragged)
	}

	x, y := d.currentDragPos.Components()
	var offset, leadingRatio, trailingRatio float64
	if d.split.Horizontal {
		widthFree := float64(d.split.Size().Width - dividerThickness())
		leadingRatio = float64(d.split.Leading.MinSize().Width) / widthFree
		trailingRatio = 1. - (float64(d.split.Trailing.MinSize().Width) / widthFree)
		offset = float64(x-d.startDragOff.X) / widthFree
	} else {
		heightFree := float64(d.split.Size().Height - dividerThickness())
		leadingRatio = float64(d.split.Leading.MinSize().Height) / heightFree
		trailingRatio = 1. - (float64(d.split.Trailing.MinSize().Height) / heightFree)
		offset = float64(y-d.startDragOff.Y) / heightFree
	}

	if offset < leadingRatio {
		offset = leadingRatio
	}
	if offset > trailingRatio {
		offset = trailingRatio
	}
	d.split.SetOffset(offset)
}

func (d *divider) MouseIn(event *desktop.MouseEvent) {
	d.hovered = true
	d.split.Refresh()
}

func (d *divider) MouseMoved(event *desktop.MouseEvent) {}

func (d *divider) MouseOut() {
	d.hovered = false
	d.split.Refresh()
}

var _ fyne.WidgetRenderer = (*dividerRenderer)(nil)

type dividerRenderer struct {
	divider    *divider
	background *canvas.Rectangle
	foreground *canvas.Rectangle
	objects    []fyne.CanvasObject
}

func (r *dividerRenderer) Destroy() {
}

func (r *dividerRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	var x, y, w, h float32
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
	r.foreground.Move(fyne.NewPos(x, y))
	r.foreground.Resize(fyne.NewSize(w, h))
}

func (r *dividerRenderer) MinSize() fyne.Size {
	if r.divider.split.Horizontal {
		return fyne.NewSize(dividerThickness(), dividerLength())
	}
	return fyne.NewSize(dividerLength(), dividerThickness())
}

func (r *dividerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *dividerRenderer) Refresh() {
	if r.divider.hovered {
		r.background.FillColor = theme.HoverColor()
	} else {
		r.background.FillColor = theme.ShadowColor()
	}
	r.background.Refresh()
	r.foreground.FillColor = theme.ForegroundColor()
	r.foreground.Refresh()
	r.Layout(r.divider.Size())
}

func dividerThickness() float32 {
	return theme.Padding() * 2
}

func dividerLength() float32 {
	return theme.Padding() * 6
}

func handleThickness() float32 {
	return theme.Padding() / 2
}

func handleLength() float32 {
	return theme.Padding() * 4
}
