package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// ScrollDirection represents the directions in which a ScrollContainer can scroll its child content.
//
// Deprecated: use container.ScrollDirection instead.
type ScrollDirection int

// Constants for valid values of ScrollDirection.
const (
	// Deprecated: use container.ScrollBoth instead
	ScrollBoth ScrollDirection = iota
	// Deprecated: use container.ScrollHorizontalOnly instead
	ScrollHorizontalOnly
	// Deprecated: use container.ScrollVerticalOnly instead
	ScrollVerticalOnly
)

type scrollBarOrientation int

// We default to vertical as 0 due to that being the original orientation offered
const (
	scrollBarOrientationVertical   scrollBarOrientation = 0
	scrollBarOrientationHorizontal scrollBarOrientation = 1
	scrollContainerMinSize                              = 32 // TODO consider the smallest useful scroll view?
)

type scrollBarRenderer struct {
	widget.BaseRenderer
	scrollBar *scrollBar
	minSize   fyne.Size
}

func (r *scrollBarRenderer) BackgroundColor() color.Color {
	return theme.ScrollBarColor()
}

func (r *scrollBarRenderer) Layout(_ fyne.Size) {
}

func (r *scrollBarRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *scrollBarRenderer) Refresh() {
}

var _ desktop.Hoverable = (*scrollBar)(nil)
var _ fyne.Draggable = (*scrollBar)(nil)

type scrollBar struct {
	BaseWidget
	area            *scrollBarArea
	draggedDistance int
	dragStart       int
	isDragged       bool
	orientation     scrollBarOrientation
}

func (b *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	return &scrollBarRenderer{scrollBar: b}
}

func (b *scrollBar) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (b *scrollBar) DragEnd() {
	b.isDragged = false
}

func (b *scrollBar) Dragged(e *fyne.DragEvent) {
	if !b.isDragged {
		b.isDragged = true
		switch b.orientation {
		case scrollBarOrientationHorizontal:
			b.dragStart = b.Position().X
		case scrollBarOrientationVertical:
			b.dragStart = b.Position().Y
		}
		b.draggedDistance = 0
	}

	switch b.orientation {
	case scrollBarOrientationHorizontal:
		b.draggedDistance += e.DraggedX
	case scrollBarOrientationVertical:
		b.draggedDistance += e.DraggedY
	}
	b.area.moveBar(b.draggedDistance+b.dragStart, b.Size())
}

func (b *scrollBar) MouseIn(e *desktop.MouseEvent) {
	b.area.MouseIn(e)
}

func (b *scrollBar) MouseMoved(*desktop.MouseEvent) {
}

func (b *scrollBar) MouseOut() {
	b.area.MouseOut()
}

func newScrollBar(area *scrollBarArea) *scrollBar {
	b := &scrollBar{area: area, orientation: area.orientation}
	b.ExtendBaseWidget(b)
	return b
}

type scrollBarAreaRenderer struct {
	widget.BaseRenderer
	area *scrollBarArea
	bar  *scrollBar
}

func (r *scrollBarAreaRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *scrollBarAreaRenderer) Layout(_ fyne.Size) {
	var barHeight, barWidth, barX, barY int
	switch r.area.orientation {
	case scrollBarOrientationHorizontal:
		barWidth, barHeight, barX, barY = r.barSizeAndOffset(r.area.scroll.Offset.X, r.area.scroll.Content.Size().Width, r.area.scroll.Size().Width)
	default:
		barHeight, barWidth, barY, barX = r.barSizeAndOffset(r.area.scroll.Offset.Y, r.area.scroll.Content.Size().Height, r.area.scroll.Size().Height)
	}
	r.bar.Move(fyne.NewPos(barX, barY))
	r.bar.Resize(fyne.NewSize(barWidth, barHeight))
}

func (r *scrollBarAreaRenderer) MinSize() fyne.Size {
	var min int
	min = theme.ScrollBarSize()
	if !r.area.isLarge {
		min = theme.ScrollBarSmallSize() * 2
	}
	switch r.area.orientation {
	case scrollBarOrientationHorizontal:
		return fyne.NewSize(theme.ScrollBarSize(), min)
	default:
		return fyne.NewSize(min, theme.ScrollBarSize())
	}
}

func (r *scrollBarAreaRenderer) Refresh() {
	r.Layout(r.area.Size())
	canvas.Refresh(r.bar)
}

func (r *scrollBarAreaRenderer) barSizeAndOffset(contentOffset, contentLength, scrollLength int) (length, width, lengthOffset, widthOffset int) {
	if scrollLength < contentLength {
		portion := float64(scrollLength) / float64(contentLength)
		length = int(float64(scrollLength) * portion)
		if length < theme.ScrollBarSize() {
			length = theme.ScrollBarSize()
		}
	} else {
		length = scrollLength
	}
	if contentOffset != 0 {
		lengthOffset = int(float64(scrollLength-length) * (float64(contentOffset) / float64(contentLength-scrollLength)))
	}
	if r.area.isLarge {
		width = theme.ScrollBarSize()
	} else {
		widthOffset = theme.ScrollBarSmallSize()
		width = theme.ScrollBarSmallSize()
	}
	return
}

var _ desktop.Hoverable = (*scrollBarArea)(nil)

type scrollBarArea struct {
	BaseWidget

	isLarge     bool
	scroll      *ScrollContainer
	orientation scrollBarOrientation
}

func (a *scrollBarArea) CreateRenderer() fyne.WidgetRenderer {
	bar := newScrollBar(a)
	return &scrollBarAreaRenderer{BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{bar}), area: a, bar: bar}
}

func (a *scrollBarArea) MouseIn(*desktop.MouseEvent) {
	a.isLarge = true
	a.scroll.Refresh()
}

func (a *scrollBarArea) MouseMoved(*desktop.MouseEvent) {
}

func (a *scrollBarArea) MouseOut() {
	a.isLarge = false
	a.scroll.Refresh()
}

func (a *scrollBarArea) moveBar(offset int, barSize fyne.Size) {
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		a.scroll.Offset.X = a.computeScrollOffset(barSize.Width, offset, a.scroll.Size().Width, a.scroll.Content.Size().Width)
	default:
		a.scroll.Offset.Y = a.computeScrollOffset(barSize.Height, offset, a.scroll.Size().Height, a.scroll.Content.Size().Height)
	}
	if f := a.scroll.onOffsetChanged; f != nil {
		f()
	}
	a.scroll.refreshWithoutOffsetUpdate()
}

func (a *scrollBarArea) computeScrollOffset(length, offset, scrollLength, contentLength int) int {
	maxOffset := scrollLength - length
	if offset < 0 {
		offset = 0
	} else if offset > maxOffset {
		offset = maxOffset
	}
	ratio := float32(offset) / float32(maxOffset)
	scrollOffset := int(ratio * float32(contentLength-scrollLength))
	return scrollOffset
}

func newScrollBarArea(scroll *ScrollContainer, orientation scrollBarOrientation) *scrollBarArea {
	a := &scrollBarArea{scroll: scroll, orientation: orientation}
	a.ExtendBaseWidget(a)
	return a
}

type scrollContainerRenderer struct {
	widget.BaseRenderer
	scroll                  *ScrollContainer
	vertArea                *scrollBarArea
	horizArea               *scrollBarArea
	leftShadow, rightShadow *widget.Shadow
	topShadow, bottomShadow *widget.Shadow
	oldMinSize              fyne.Size
}

func (r *scrollContainerRenderer) layoutBars(size fyne.Size) {
	if r.scroll.Direction != ScrollHorizontalOnly {
		r.vertArea.Resize(fyne.NewSize(r.vertArea.MinSize().Width, size.Height))
		r.vertArea.Move(fyne.NewPos(r.scroll.Size().Width-r.vertArea.Size().Width, 0))
		r.topShadow.Resize(fyne.NewSize(size.Width, 0))
		r.bottomShadow.Resize(fyne.NewSize(size.Width, 0))
		r.bottomShadow.Move(fyne.NewPos(0, r.scroll.size.Height))
	}

	if r.scroll.Direction != ScrollVerticalOnly {
		r.horizArea.Resize(fyne.NewSize(size.Width, r.horizArea.MinSize().Height))
		r.horizArea.Move(fyne.NewPos(0, r.scroll.Size().Height-r.horizArea.Size().Height))
		r.leftShadow.Resize(fyne.NewSize(0, size.Height))
		r.rightShadow.Resize(fyne.NewSize(0, size.Height))
		r.rightShadow.Move(fyne.NewPos(r.scroll.size.Width, 0))
	}

	r.updatePosition()
}

func (r *scrollContainerRenderer) Layout(size fyne.Size) {
	c := r.scroll.Content
	c.Resize(c.MinSize().Max(size))

	r.layoutBars(size)
}

func (r *scrollContainerRenderer) MinSize() fyne.Size {
	return r.scroll.MinSize()
}

func (r *scrollContainerRenderer) Refresh() {
	if len(r.BaseRenderer.Objects()) == 0 || r.BaseRenderer.Objects()[0] != r.scroll.Content {
		// push updated content object to baseRenderer
		r.BaseRenderer.SetObjects([]fyne.CanvasObject{r.scroll.Content})
	}
	if r.oldMinSize == r.scroll.Content.MinSize() && r.oldMinSize == r.scroll.Content.Size() &&
		(r.scroll.Size().Width <= r.oldMinSize.Width && r.scroll.Size().Height <= r.oldMinSize.Height) {
		r.layoutBars(r.scroll.Size())
		return
	}

	r.oldMinSize = r.scroll.Content.MinSize()
	r.Layout(r.scroll.Size())
}

func (r *scrollContainerRenderer) handleAreaVisibility(contentSize int, scrollSize int, area *scrollBarArea) {
	if contentSize <= scrollSize {
		area.Hide()
	} else if r.scroll.Visible() {
		area.Show()
	}
}

func (r *scrollContainerRenderer) handleShadowVisibility(offset int, contentSize int, scrollSize int, shadowStart fyne.CanvasObject, shadowEnd fyne.CanvasObject) {
	if !r.scroll.Visible() {
		return
	}
	if offset > 0 {
		shadowStart.Show()
	} else {
		shadowStart.Hide()
	}
	if offset < contentSize-scrollSize {
		shadowEnd.Show()
	} else {
		shadowEnd.Hide()
	}
}

func (r *scrollContainerRenderer) updatePosition() {
	scrollSize := r.scroll.Size()
	contentSize := r.scroll.Content.Size()

	r.scroll.Content.Move(fyne.NewPos(-r.scroll.Offset.X, -r.scroll.Offset.Y))

	if r.scroll.Direction != ScrollHorizontalOnly {
		r.handleAreaVisibility(contentSize.Height, scrollSize.Height, r.vertArea)
		r.handleShadowVisibility(r.scroll.Offset.Y, contentSize.Height, scrollSize.Height, r.topShadow, r.bottomShadow)
		Renderer(r.vertArea).Layout(r.scroll.size)
	}
	if r.scroll.Direction != ScrollVerticalOnly {
		r.handleAreaVisibility(contentSize.Width, scrollSize.Width, r.horizArea)
		r.handleShadowVisibility(r.scroll.Offset.X, contentSize.Width, scrollSize.Width, r.leftShadow, r.rightShadow)
		Renderer(r.horizArea).Layout(r.scroll.size)
	}

	if r.scroll.Direction != ScrollHorizontalOnly {
		canvas.Refresh(r.vertArea) // this is required to force the canvas to update, we have no "Redraw()"
	} else {
		canvas.Refresh(r.horizArea) // this is required like above but if we are horizontal
	}
}

// ScrollContainer defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
//
// Deprecated: use container.Scroll instead.
type ScrollContainer struct {
	BaseWidget
	minSize         fyne.Size
	Direction       ScrollDirection
	Content         fyne.CanvasObject
	Offset          fyne.Position
	onOffsetChanged func()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *ScrollContainer) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	scr := &scrollContainerRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{s.Content}),
		scroll:       s,
	}
	if s.Direction != ScrollHorizontalOnly {
		scr.vertArea = newScrollBarArea(s, scrollBarOrientationVertical)
		scr.topShadow = widget.NewShadow(widget.ShadowBottom, widget.SubmergedContentLevel)
		scr.bottomShadow = widget.NewShadow(widget.ShadowTop, widget.SubmergedContentLevel)
		scr.SetObjects(append(scr.Objects(), scr.vertArea, scr.topShadow, scr.bottomShadow))
	}
	if s.Direction != ScrollVerticalOnly {
		scr.horizArea = newScrollBarArea(s, scrollBarOrientationHorizontal)
		scr.leftShadow = widget.NewShadow(widget.ShadowRight, widget.SubmergedContentLevel)
		scr.rightShadow = widget.NewShadow(widget.ShadowLeft, widget.SubmergedContentLevel)
		scr.SetObjects(append(scr.Objects(), scr.horizArea, scr.leftShadow, scr.rightShadow))
	}
	return scr
}

//ScrollToBottom will scroll content to container bottom - to show latest info which end user just added
func (s *ScrollContainer) ScrollToBottom() {
	s.Offset.Y = s.Content.Size().Height - s.Size().Height
	s.Refresh()
}

//ScrollToTop will scroll content to container top
func (s *ScrollContainer) ScrollToTop() {
	s.Offset.Y = 0
	s.Refresh()
}

// DragEnd will stop scrolling on mobile has stopped
func (s *ScrollContainer) DragEnd() {
}

// Dragged will scroll on any drag - bar or otherwise - for mobile
func (s *ScrollContainer) Dragged(e *fyne.DragEvent) {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}

	if s.updateOffset(e.DraggedX, e.DraggedY) {
		s.refreshWithoutOffsetUpdate()
	}
}

// MinSize returns the smallest size this widget can shrink to
func (s *ScrollContainer) MinSize() fyne.Size {
	min := fyne.NewSize(scrollContainerMinSize, scrollContainerMinSize).Max(s.minSize)
	switch s.Direction {
	case ScrollHorizontalOnly:
		min.Height = fyne.Max(min.Height, s.Content.MinSize().Height)
	case ScrollVerticalOnly:
		min.Width = fyne.Max(min.Width, s.Content.MinSize().Width)
	}
	return min
}

// SetMinSize specifies a minimum size for this scroll container.
// If the specified size is larger than the content size then scrolling will not be enabled
// This can be helpful to appear larger than default if the layout is collapsing this widget.
func (s *ScrollContainer) SetMinSize(size fyne.Size) {
	s.minSize = size
}

// Refresh causes this widget to be redrawn in it's current state
func (s *ScrollContainer) Refresh() {
	s.updateOffset(0, 0)
	s.refreshWithoutOffsetUpdate()
}

func (s *ScrollContainer) refreshWithoutOffsetUpdate() {
	s.BaseWidget.Refresh()
}

// Resize sets a new size for the scroll container.
func (s *ScrollContainer) Resize(size fyne.Size) {
	if size != s.size {
		s.size = size
		s.Refresh()
	}
}

// Scrolled is called when an input device triggers a scroll event
func (s *ScrollContainer) Scrolled(ev *fyne.ScrollEvent) {
	dx, dy := ev.DeltaX, ev.DeltaY
	if s.Size().Width < s.Content.MinSize().Width && s.Size().Height >= s.Content.MinSize().Height && dx == 0 {
		dx, dy = dy, dx
	}
	if s.updateOffset(dx, dy) {
		s.refreshWithoutOffsetUpdate()
	}
}

func (s *ScrollContainer) updateOffset(deltaX, deltaY int) bool {
	if s.Content.Size().Width <= s.Size().Width && s.Content.Size().Height <= s.Size().Height {
		if s.Offset.X != 0 || s.Offset.Y != 0 {
			s.Offset.X = 0
			s.Offset.Y = 0
			return true
		}
		return false
	}
	s.Offset.X = computeOffset(s.Offset.X, -deltaX, s.Size().Width, s.Content.MinSize().Width)
	s.Offset.Y = computeOffset(s.Offset.Y, -deltaY, s.Size().Height, s.Content.MinSize().Height)
	if f := s.onOffsetChanged; f != nil {
		f()
	}
	return true
}

func computeOffset(start, delta, outerWidth, innerWidth int) int {
	offset := start + delta
	if offset+outerWidth >= innerWidth {
		offset = innerWidth - outerWidth
	}
	if offset < 0 {
		offset = 0
	}
	return offset
}

// NewScrollContainer creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed object.
//
// Deprecated: use container.NewScroll instead.
func NewScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	return newScrollContainerWithDirection(ScrollBoth, content)
}

// NewHScrollContainer create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Width to be smaller than that of the passed object.
//
// Deprecated: use container.NewHScroll instead.
func NewHScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	return newScrollContainerWithDirection(ScrollHorizontalOnly, content)
}

// NewVScrollContainer create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Height to be smaller than that of the passed object.
//
// Deprecated: use container.NewVScroll instead.
func NewVScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	return newScrollContainerWithDirection(ScrollVerticalOnly, content)
}

func newScrollContainerWithDirection(direction ScrollDirection, content fyne.CanvasObject) *ScrollContainer {
	s := &ScrollContainer{
		Direction: direction,
		Content:   content,
	}
	s.ExtendBaseWidget(s)
	return s
}
