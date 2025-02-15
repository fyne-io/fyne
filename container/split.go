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

	// to communicate to the renderer that the next refresh
	// is just an offset update (ie a resize and move only)
	// cleared by renderer in Refresh()
	offsetUpdated bool
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
	s.offsetUpdated = true
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
		dividerSize.Width = dividerThickness(r.divider)
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
		dividerSize.Height = dividerThickness(r.divider)
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
	if r.split.offsetUpdated {
		r.Layout(r.split.Size())
		r.split.offsetUpdated = false
		return
	}

	r.objects[0] = r.split.Leading
	// [1] is divider which doesn't change
	r.objects[2] = r.split.Trailing
	r.Layout(r.split.Size())

	r.split.Leading.Refresh()
	r.divider.Refresh()
	r.split.Trailing.Refresh()
	canvas.Refresh(r.split)
}

func (r *splitContainerRenderer) computeSplitLengths(total, lMin, tMin float32) (float32, float32) {
	available := float64(total - dividerThickness(r.divider))
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
	th := d.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	background := canvas.NewRectangle(th.Color(theme.ColorNameShadow, v))
	foreground := canvas.NewRectangle(th.Color(theme.ColorNameForeground, v))
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
		widthFree := float64(d.split.Size().Width - dividerThickness(d))
		leadingRatio = float64(d.split.Leading.MinSize().Width) / widthFree
		trailingRatio = 1. - (float64(d.split.Trailing.MinSize().Width) / widthFree)
		offset = float64(x-d.startDragOff.X) / widthFree
	} else {
		heightFree := float64(d.split.Size().Height - dividerThickness(d))
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
		x = (dividerThickness(r.divider) - handleThickness(r.divider)) / 2
		y = (size.Height - handleLength(r.divider)) / 2
		w = handleThickness(r.divider)
		h = handleLength(r.divider)
	} else {
		x = (size.Width - handleLength(r.divider)) / 2
		y = (dividerThickness(r.divider) - handleThickness(r.divider)) / 2
		w = handleLength(r.divider)
		h = handleThickness(r.divider)
	}
	r.foreground.Move(fyne.NewPos(x, y))
	r.foreground.Resize(fyne.NewSize(w, h))
}

func (r *dividerRenderer) MinSize() fyne.Size {
	if r.divider.split.Horizontal {
		return fyne.NewSize(dividerThickness(r.divider), dividerLength(r.divider))
	}
	return fyne.NewSize(dividerLength(r.divider), dividerThickness(r.divider))
}

func (r *dividerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *dividerRenderer) Refresh() {
	th := r.divider.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	if r.divider.hovered {
		r.background.FillColor = th.Color(theme.ColorNameHover, v)
	} else {
		r.background.FillColor = th.Color(theme.ColorNameShadow, v)
	}
	r.background.Refresh()
	r.foreground.FillColor = th.Color(theme.ColorNameForeground, v)
	r.foreground.Refresh()
	r.Layout(r.divider.Size())
}

func dividerTheme(d *divider) fyne.Theme {
	if d == nil {
		return theme.Current()
	}

	return d.Theme()
}

func dividerThickness(d *divider) float32 {
	th := dividerTheme(d)
	return th.Size(theme.SizeNamePadding) * 2
}

func dividerLength(d *divider) float32 {
	th := dividerTheme(d)
	return th.Size(theme.SizeNamePadding) * 6
}

func handleThickness(d *divider) float32 {
	th := dividerTheme(d)
	return th.Size(theme.SizeNamePadding) / 2

}

func handleLength(d *divider) float32 {
	th := dividerTheme(d)
	return th.Size(theme.SizeNamePadding) * 4
}
