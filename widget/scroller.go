package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

type scrollBarOrientation int

// We default to vertical as 0 due to that being the original orientation offered
const (
	scrollBarOrientationVertical   scrollBarOrientation = 0
	scrollBarOrientationHorizontal scrollBarOrientation = 1
)

type scrollBarRenderer struct {
	scrollBar *scrollBar

	minSize fyne.Size
}

func (r *scrollBarRenderer) BackgroundColor() color.Color {
	return theme.ScrollBarColor()
}

func (r *scrollBarRenderer) Destroy() {
}

func (r *scrollBarRenderer) Layout(size fyne.Size) {
}

func (r *scrollBarRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *scrollBarRenderer) Objects() []fyne.CanvasObject {
	return nil
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

func (b *scrollBar) DragEnd() {
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
	area *scrollBarArea
	bar  *scrollBar

	objects []fyne.CanvasObject
}

func (r *scrollBarAreaRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *scrollBarAreaRenderer) Destroy() {
}

func (r *scrollBarAreaRenderer) Layout(size fyne.Size) {
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

func (r *scrollBarAreaRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *scrollBarAreaRenderer) Refresh() {
	r.Layout(r.area.Size())
	canvas.Refresh(r.bar)
}

func (r *scrollBarAreaRenderer) barSizeAndOffset(contentOffset, contentLength, scrollLength int) (length, width, lengthOffset, widthOffset int) {
	if scrollLength < contentLength {
		portion := float64(scrollLength) / float64(contentLength)
		length = int(float64(scrollLength) * portion)
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
	return &scrollBarAreaRenderer{area: a, bar: bar, objects: []fyne.CanvasObject{bar}}
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
	scroll                  *ScrollContainer
	vertArea                *scrollBarArea
	horizArea               *scrollBarArea
	leftShadow, rightShadow *shadow
	topShadow, bottomShadow *shadow

	objects []fyne.CanvasObject
}

func (r *scrollContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *scrollContainerRenderer) Destroy() {
}

func (r *scrollContainerRenderer) Layout(size fyne.Size) {
	// The scroll bar needs to be resized and moved on the far right
	r.horizArea.Resize(fyne.NewSize(size.Width, r.horizArea.MinSize().Height))
	r.vertArea.Resize(fyne.NewSize(r.vertArea.MinSize().Width, size.Height))

	r.horizArea.Move(fyne.NewPos(0, r.scroll.Size().Height-r.horizArea.Size().Height))
	r.vertArea.Move(fyne.NewPos(r.scroll.Size().Width-r.vertArea.Size().Width, 0))

	r.leftShadow.Resize(fyne.NewSize(0, size.Height))
	r.rightShadow.Resize(fyne.NewSize(0, size.Height))
	r.topShadow.Resize(fyne.NewSize(size.Width, 0))
	r.bottomShadow.Resize(fyne.NewSize(size.Width, 0))

	r.rightShadow.Move(fyne.NewPos(r.scroll.size.Width, 0))
	r.bottomShadow.Move(fyne.NewPos(0, r.scroll.size.Height))

	c := r.scroll.Content
	c.Resize(c.MinSize().Union(size))

	r.updatePosition()
}

func (r *scrollContainerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(32, 32) // TODO consider the smallest useful scroll view?
}

func (r *scrollContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *scrollContainerRenderer) Refresh() {
	r.leftShadow.depth = theme.Padding() * 2
	r.rightShadow.depth = theme.Padding() * 2
	r.topShadow.depth = theme.Padding() * 2
	r.bottomShadow.depth = theme.Padding() * 2

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
	scrollWidth := r.scroll.Size().Width
	contentWidth := r.scroll.Content.Size().Width
	scrollHeight := r.scroll.Size().Height
	contentHeight := r.scroll.Content.Size().Height
	r.handleAreaVisibility(contentWidth, scrollWidth, r.horizArea)
	r.handleAreaVisibility(contentHeight, scrollHeight, r.vertArea)

	r.scroll.Content.Move(fyne.NewPos(-r.scroll.Offset.X, -r.scroll.Offset.Y))
	r.handleShadowVisibility(r.scroll.Offset.X, contentWidth, scrollWidth, r.leftShadow, r.rightShadow)
	r.handleShadowVisibility(r.scroll.Offset.Y, contentHeight, scrollHeight, r.topShadow, r.bottomShadow)

	Renderer(r.vertArea).Layout(r.scroll.size)
	Renderer(r.horizArea).Layout(r.scroll.size)

	canvas.Refresh(r.vertArea)  // this is required to force the canvas to update, we have no "Redraw()"
	canvas.Refresh(r.horizArea) // this is required like above but if we are horizontal
}

// ScrollContainer defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
type ScrollContainer struct {
	BaseWidget

	Content fyne.CanvasObject
	Offset  fyne.Position
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *ScrollContainer) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	hbar := newScrollBarArea(s, scrollBarOrientationHorizontal)
	vbar := newScrollBarArea(s, scrollBarOrientationVertical)
	leftShadow := newShadow(shadowRight, theme.Padding()*2)
	rightShadow := newShadow(shadowLeft, theme.Padding()*2)
	topShadow := newShadow(shadowBottom, theme.Padding()*2)
	bottomShadow := newShadow(shadowTop, theme.Padding()*2)
	return &scrollContainerRenderer{
		objects:      []fyne.CanvasObject{s.Content, hbar, vbar, topShadow, bottomShadow, leftShadow, rightShadow},
		scroll:       s,
		horizArea:    hbar,
		vertArea:     vbar,
		leftShadow:   leftShadow,
		rightShadow:  rightShadow,
		topShadow:    topShadow,
		bottomShadow: bottomShadow,
	}
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
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
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
// Note that this may cause the MinSize to be smaller than that of the passed objects.
func NewScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	s := &ScrollContainer{Content: content}
	s.ExtendBaseWidget(s)
	return s
}
