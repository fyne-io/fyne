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
	half := float64(halfDividerThickness())

	var min, max float64
	offset := r.split.Offset
	if r.split.Horizontal {
		sw := float64(size.Width)
		lw := float64(r.split.Leading.MinSize().Width)
		tw := float64(r.split.Trailing.MinSize().Width)
		max = 1.0 - ((tw + half) / sw)
		min = (lw + half) / sw
	} else {
		sh := float64(size.Height)
		lh := float64(r.split.Leading.MinSize().Height)
		th := float64(r.split.Trailing.MinSize().Height)
		max = 1.0 - ((th + half) / sh)
		min = (lh + half) / sh
	}
	if offset < min {
		offset = min
	}
	if offset > max {
		offset = max
	}

	if r.split.Horizontal {
		x := 0
		leadingPos.X = x
		leadingSize.Width = int(offset*float64(size.Width) - half)
		leadingSize.Height = size.Height
		x += leadingSize.Width
		dividerPos.X = x
		dividerSize.Width = dividerThickness()
		dividerSize.Height = size.Height
		x += dividerSize.Width
		trailingPos.X = x
		trailingSize.Width = int((1.0-offset)*float64(size.Width) - half)
		trailingSize.Height = size.Height
	} else {
		y := 0
		leadingPos.Y = y
		leadingSize.Width = size.Width
		leadingSize.Height = int(offset*float64(size.Height) - half)
		y += leadingSize.Height
		dividerPos.Y = y
		dividerSize.Width = size.Width
		dividerSize.Height = dividerThickness()
		y += dividerSize.Height
		trailingPos.Y = y
		trailingSize.Width = size.Width
		trailingSize.Height = int((1.0-offset)*float64(size.Height) - half)
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

func halfDividerThickness() int {
	return theme.Padding()
}

func dividerThickness() int {
	return halfDividerThickness() * 2
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
