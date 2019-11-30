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
	area                 *scrollBarArea
	draggedDistanceHoriz int
	draggedDistanceVert  int
	dragStartHoriz       int
	dragStartVert        int
	isDragged            bool
	orientation          scrollBarOrientation
}

func (b *scrollBar) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	return &scrollBarRenderer{scrollBar: b}
}

func (b *scrollBar) DragEnd() {
}

func (b *scrollBar) Dragged(e *fyne.DragEvent) {
	if !b.isDragged {
		b.isDragged = true
		switch b.orientation {
		case scrollBarOrientationHorizontal:
			b.dragStartHoriz = b.Position().X
		case scrollBarOrientationVertical:
			b.dragStartVert = b.Position().Y
		}
		b.draggedDistanceHoriz = 0
		b.draggedDistanceVert = 0
	}

	switch b.orientation {
	case scrollBarOrientationHorizontal:
		b.draggedDistanceHoriz += e.DraggedX
		b.area.moveHorizontalBar(b.draggedDistanceHoriz + b.dragStartHoriz)
	case scrollBarOrientationVertical:
		b.draggedDistanceVert += e.DraggedY
		b.area.moveVerticalBar(b.draggedDistanceVert + b.dragStartVert)
	}
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
	return &scrollBar{area: area, orientation: area.orientation}
}

type scrollBarAreaRenderer struct {
	area        *scrollBarArea
	bar         *scrollBar
	orientation scrollBarOrientation

	objects []fyne.CanvasObject
}

func (r *scrollBarAreaRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *scrollBarAreaRenderer) Destroy() {
}

func (r *scrollBarAreaRenderer) Layout(size fyne.Size) {
	switch r.orientation {
	case scrollBarOrientationHorizontal:
		r.updateHorizontalBarPosition()
	case scrollBarOrientationVertical:
		r.updateVerticalBarPosition()
	}
}

func (r *scrollBarAreaRenderer) MinSize() fyne.Size {
	var min int
	min = theme.ScrollBarSize()
	switch r.orientation {
	case scrollBarOrientationHorizontal:
		if !r.area.isTall {
			min = theme.ScrollBarSmallSize() * 2
		}
		return fyne.NewSize(theme.ScrollBarSize(), min)
	default:
		if !r.area.isWide {
			min = theme.ScrollBarSmallSize() * 2
		}
		return fyne.NewSize(min, theme.ScrollBarSize())
	}
}

func (r *scrollBarAreaRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *scrollBarAreaRenderer) Refresh() {
	switch r.orientation {
	case scrollBarOrientationHorizontal:
		r.updateHorizontalBarPosition()
	default:
		r.updateVerticalBarPosition()
	}
	canvas.Refresh(r.bar)
}

func (r *scrollBarAreaRenderer) updateHorizontalBarPosition() {
	barWidth := r.horizontalBarWidth()
	barRatio := float32(0.0)
	if r.area.scroll.Offset.X != 0 {
		barRatio = float32(r.area.scroll.Offset.X) / float32(r.area.scroll.Content.Size().Width-r.area.scroll.Size().Width)
	}
	barX := int(float32(r.area.scroll.size.Width-barWidth) * barRatio)

	var barY, barHeight int
	if r.area.isTall {
		barHeight = theme.ScrollBarSize()
	} else {
		barY = theme.ScrollBarSmallSize()
		barHeight = theme.ScrollBarSmallSize()
	}

	r.bar.Resize(fyne.NewSize(barWidth, barHeight))
	r.bar.Move(fyne.NewPos(barX, barY))
}

func (r *scrollBarAreaRenderer) updateVerticalBarPosition() {
	barHeight := r.verticalBarHeight()
	barRatio := float32(0.0)
	if r.area.scroll.Offset.Y != 0 {
		barRatio = float32(r.area.scroll.Offset.Y) / float32(r.area.scroll.Content.Size().Height-r.area.scroll.Size().Height)
	}
	barY := int(float32(r.area.scroll.size.Height-barHeight) * barRatio)

	var barX, barWidth int
	if r.area.isWide {
		barWidth = theme.ScrollBarSize()
	} else {
		barX = theme.ScrollBarSmallSize()
		barWidth = theme.ScrollBarSmallSize()
	}

	r.bar.Resize(fyne.NewSize(barWidth, barHeight))
	r.bar.Move(fyne.NewPos(barX, barY))
}

func (r *scrollBarAreaRenderer) horizontalBarWidth() int {
	portion := float32(r.area.size.Width) / float32(r.area.scroll.Content.Size().Width)
	if portion > 1.0 {
		portion = 1.0
	}

	return int(float32(r.area.size.Width) * portion)
}

func (r *scrollBarAreaRenderer) verticalBarHeight() int {
	portion := float32(r.area.size.Height) / float32(r.area.scroll.Content.Size().Height)
	if portion > 1.0 {
		portion = 1.0
	}

	return int(float32(r.area.size.Height) * portion)
}

var _ desktop.Hoverable = (*scrollBarArea)(nil)

type scrollBarArea struct {
	BaseWidget

	isWide      bool
	isTall      bool
	scroll      *ScrollContainer
	orientation scrollBarOrientation
}

func (a *scrollBarArea) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	bar := newScrollBar(a)
	return &scrollBarAreaRenderer{area: a, bar: bar, orientation: a.orientation, objects: []fyne.CanvasObject{bar}}
}

// MinSize returns the size that this widget should not shrink below
func (a *scrollBarArea) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

func (a *scrollBarArea) MouseIn(*desktop.MouseEvent) {
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		a.isTall = true
	case scrollBarOrientationVertical:
		a.isWide = true
	}
	a.scroll.Refresh()
}

func (a *scrollBarArea) MouseMoved(*desktop.MouseEvent) {
}

func (a *scrollBarArea) MouseOut() {
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		a.isTall = false
	case scrollBarOrientationVertical:
		a.isWide = false
	}
	a.scroll.Refresh()
}

func (a *scrollBarArea) moveHorizontalBar(x int) {
	render := Renderer(a).(*scrollBarAreaRenderer)
	barWidth := render.horizontalBarWidth()
	scrollWidth := a.scroll.Size().Width
	maxX := scrollWidth - barWidth

	if x < 0 {
		x = 0
	} else if x > maxX {
		x = maxX
	}

	ratio := float32(x) / float32(maxX)
	a.scroll.Offset.X = int(ratio * float32(a.scroll.Content.Size().Width-scrollWidth))

	Refresh(a.scroll)
}

func (a *scrollBarArea) moveVerticalBar(y int) {
	render := Renderer(a).(*scrollBarAreaRenderer)
	barHeight := render.verticalBarHeight()
	scrollHeight := a.scroll.Size().Height
	maxY := scrollHeight - barHeight

	if y < 0 {
		y = 0
	} else if y > maxY {
		y = maxY
	}

	ratio := float32(y) / float32(maxY)
	a.scroll.Offset.Y = int(ratio * float32(a.scroll.Content.Size().Height-scrollHeight))

	Refresh(a.scroll)
}

func newScrollBarArea(scroll *ScrollContainer, orientation scrollBarOrientation) *scrollBarArea {
	return &scrollBarArea{scroll: scroll, orientation: orientation}
}

type scrollContainerRenderer struct {
	scroll                  *ScrollContainer
	vertArea                *scrollBarArea
	horizArea               *scrollBarArea
	leftShadow, rightShadow fyne.CanvasObject
	topShadow, bottomShadow fyne.CanvasObject

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
	return fyne.NewSize(25, 25) // TODO consider the smallest useful scroll view?
}

func (r *scrollContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *scrollContainerRenderer) Refresh() {
	r.Layout(r.scroll.Size())
}

func (r *scrollContainerRenderer) calculateScrollPosition(contentSize int, scrollSize int, offset int, area *scrollBarArea) {
	if contentSize <= scrollSize {
		if area == r.horizArea {
			r.scroll.Offset.X = 0
		} else {
			r.scroll.Offset.Y = 0
		}
		area.Hide()
	} else if r.scroll.Visible() {
		area.Show()
		if contentSize-offset < scrollSize {
			if area == r.horizArea {
				r.scroll.Offset.X = contentSize - scrollSize
			} else {
				r.scroll.Offset.Y = contentSize - scrollSize
			}
		}
	}
}

func (r *scrollContainerRenderer) calculateShadows(offset int, contentSize int, scrollSize int, shadowStart fyne.CanvasObject, shadowEnd fyne.CanvasObject) {
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
	r.calculateScrollPosition(contentWidth, scrollWidth, r.scroll.Offset.X, r.horizArea)
	r.calculateScrollPosition(contentHeight, scrollHeight, r.scroll.Offset.Y, r.vertArea)

	r.scroll.Content.Move(fyne.NewPos(-r.scroll.Offset.X, -r.scroll.Offset.Y))
	canvas.Refresh(r.scroll.Content)

	r.calculateShadows(r.scroll.Offset.X, contentWidth, scrollWidth, r.leftShadow, r.rightShadow)
	r.calculateShadows(r.scroll.Offset.Y, contentHeight, scrollHeight, r.topShadow, r.bottomShadow)

	Renderer(r.vertArea).Layout(r.scroll.size)
	Renderer(r.horizArea).Layout(r.scroll.size)
}

// ScrollContainer defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
type ScrollContainer struct {
	BaseWidget

	Content                   fyne.CanvasObject
	Offset                    fyne.Position
	horizontalDraggedDistance int
	verticalDraggedDistance   int
	hbar                      *scrollBarArea
	vbar                      *scrollBarArea
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *ScrollContainer) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	s.hbar = newScrollBarArea(s, scrollBarOrientationHorizontal)
	s.vbar = newScrollBarArea(s, scrollBarOrientationVertical)
	leftShadow := newShadow(shadowRight, theme.Padding()*2)
	rightShadow := newShadow(shadowLeft, theme.Padding()*2)
	topShadow := newShadow(shadowBottom, theme.Padding()*2)
	bottomShadow := newShadow(shadowTop, theme.Padding()*2)
	return &scrollContainerRenderer{
		objects:      []fyne.CanvasObject{s.Content, s.hbar, s.vbar, topShadow, bottomShadow, leftShadow, rightShadow},
		scroll:       s,
		horizArea:    s.hbar,
		vertArea:     s.vbar,
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
	maxX := s.Content.Size().Width
	maxY := s.Content.Size().Height

	s.horizontalDraggedDistance -= e.DraggedX
	s.verticalDraggedDistance -= e.DraggedY

	if s.horizontalDraggedDistance > maxX {
		s.horizontalDraggedDistance = maxX
	} else if s.horizontalDraggedDistance < 0 {
		s.horizontalDraggedDistance = 0
	}

	if s.verticalDraggedDistance > maxY {
		s.verticalDraggedDistance = maxY
	} else if s.verticalDraggedDistance < 0 {
		s.verticalDraggedDistance = 0
	}
	s.Offset.X = s.horizontalDraggedDistance
	s.Offset.Y = s.verticalDraggedDistance
	s.Refresh()

}

// MinSize returns the smallest size this widget can shrink to
func (s *ScrollContainer) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// Scrolled is called when an input device triggers a scroll event
func (s *ScrollContainer) Scrolled(ev *fyne.ScrollEvent) {
	if s.Content.Size().FitsInto(s.Size()) {
		return
	}

	var deltaX, deltaY int
	if s.hbar.Visible() && !s.vbar.Visible() && ev.DeltaX == 0 {
		deltaX = ev.DeltaY
		deltaY = ev.DeltaX
	} else {
		deltaX = ev.DeltaX
		deltaY = ev.DeltaY
	}
	s.Offset.X = computeOffset(s.Offset.X, -deltaX, s.Size().Width, s.Content.Size().Width)
	s.Offset.Y = computeOffset(s.Offset.Y, -deltaY, s.Size().Height, s.Content.Size().Height)

	s.Refresh()
}

func computeOffset(start, delta, outerWidth, innerWidth int) int {
	offset := start + delta
	if offset < 0 {
		offset = 0
	} else if offset+outerWidth >= innerWidth {
		offset = innerWidth - outerWidth
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
