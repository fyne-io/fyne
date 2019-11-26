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

func (s *scrollBar) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

func (s *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	return &scrollBarRenderer{scrollBar: s}
}

func (s *scrollBar) DragEnd() {
}

func (s *scrollBar) Dragged(e *fyne.DragEvent) {
	if !s.isDragged {
		s.isDragged = true
		switch s.orientation {
		case scrollBarOrientationHorizontal:
			s.dragStartHoriz = s.Position().X
		case scrollBarOrientationVertical:
			s.dragStartVert = s.Position().Y
		}
		s.draggedDistanceHoriz = 0
		s.draggedDistanceVert = 0
	}

	switch s.orientation {
	case scrollBarOrientationHorizontal:
		s.draggedDistanceHoriz += e.DraggedX
		s.area.moveHorizontalBar(s.draggedDistanceHoriz + s.dragStartHoriz)
	case scrollBarOrientationVertical:
		s.draggedDistanceVert += e.DraggedY
		s.area.moveVerticalBar(s.draggedDistanceVert + s.dragStartVert)
	}
}

func (s *scrollBar) MouseIn(e *desktop.MouseEvent) {
	s.area.MouseIn(e)
}

func (s *scrollBar) MouseMoved(*desktop.MouseEvent) {
}

func (s *scrollBar) MouseOut() {
	s.area.MouseOut()
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

func (s *scrollBarAreaRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s *scrollBarAreaRenderer) Destroy() {
}

func (s *scrollBarAreaRenderer) Layout(size fyne.Size) {
	switch s.orientation {
	case scrollBarOrientationHorizontal:
		s.updateHorizontalBarPosition()
	case scrollBarOrientationVertical:
		s.updateVerticalBarPosition()
	}
}

func (s *scrollBarAreaRenderer) MinSize() fyne.Size {
	var min int
	min = theme.ScrollBarSize()
	switch s.orientation {
	case scrollBarOrientationHorizontal:
		if !s.area.isTall {
			min = theme.ScrollBarSmallSize() * 2
		}
		return fyne.NewSize(theme.ScrollBarSize(), min)
	default:
		if !s.area.isWide {
			min = theme.ScrollBarSmallSize() * 2
		}
		return fyne.NewSize(min, theme.ScrollBarSize())
	}
}

func (s *scrollBarAreaRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *scrollBarAreaRenderer) Refresh() {
	switch s.orientation {
	case scrollBarOrientationHorizontal:
		s.updateHorizontalBarPosition()
	default:
		s.updateVerticalBarPosition()
	}
	canvas.Refresh(s.bar)
}

func (s *scrollBarAreaRenderer) updateHorizontalBarPosition() {
	barWidth := s.horizontalBarWidth()
	barRatio := float32(0.0)
	if s.area.scroll.Offset.X != 0 {
		barRatio = float32(s.area.scroll.Offset.X) / float32(s.area.scroll.Content.Size().Width-s.area.scroll.Size().Width)
	}
	barX := int(float32(s.area.scroll.size.Width-barWidth) * barRatio)

	var barY, barHeight int
	if s.area.isTall {
		barHeight = theme.ScrollBarSize()
	} else {
		barY = theme.ScrollBarSmallSize()
		barHeight = theme.ScrollBarSmallSize()
	}

	s.bar.Resize(fyne.NewSize(barWidth, barHeight))
	s.bar.Move(fyne.NewPos(barX, barY))
}

func (s *scrollBarAreaRenderer) updateVerticalBarPosition() {
	barHeight := s.verticalBarHeight()
	barRatio := float32(0.0)
	if s.area.scroll.Offset.Y != 0 {
		barRatio = float32(s.area.scroll.Offset.Y) / float32(s.area.scroll.Content.Size().Height-s.area.scroll.Size().Height)
	}
	barY := int(float32(s.area.scroll.size.Height-barHeight) * barRatio)

	var barX, barWidth int
	if s.area.isWide {
		barWidth = theme.ScrollBarSize()
	} else {
		barX = theme.ScrollBarSmallSize()
		barWidth = theme.ScrollBarSmallSize()
	}

	s.bar.Resize(fyne.NewSize(barWidth, barHeight))
	s.bar.Move(fyne.NewPos(barX, barY))
}

func (s *scrollBarAreaRenderer) horizontalBarWidth() int {
	portion := float32(s.area.size.Width) / float32(s.area.scroll.Content.Size().Width)
	if portion > 1.0 {
		portion = 1.0
	}

	return int(float32(s.area.size.Width) * portion)
}

func (s *scrollBarAreaRenderer) verticalBarHeight() int {
	portion := float32(s.area.size.Height) / float32(s.area.scroll.Content.Size().Height)
	if portion > 1.0 {
		portion = 1.0
	}

	return int(float32(s.area.size.Height) * portion)
}

var _ desktop.Hoverable = (*scrollBarArea)(nil)

type scrollBarArea struct {
	BaseWidget

	isWide      bool
	isTall      bool
	scroll      *ScrollContainer
	orientation scrollBarOrientation
}

func (s *scrollBarArea) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	bar := newScrollBar(s)
	return &scrollBarAreaRenderer{area: s, bar: bar, orientation: s.orientation, objects: []fyne.CanvasObject{bar}}
}

// MinSize returns the size that this widget should not shrink below
func (s *scrollBarArea) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

func (s *scrollBarArea) MouseIn(*desktop.MouseEvent) {
	switch s.orientation {
	case scrollBarOrientationHorizontal:
		s.isTall = true
	case scrollBarOrientationVertical:
		s.isWide = true
	}
	Refresh(s.scroll)
}

func (s *scrollBarArea) MouseMoved(*desktop.MouseEvent) {
}

func (s *scrollBarArea) MouseOut() {
	switch s.orientation {
	case scrollBarOrientationHorizontal:
		s.isTall = false
	case scrollBarOrientationVertical:
		s.isWide = false
	}
	Refresh(s.scroll)
}

func (s *scrollBarArea) moveHorizontalBar(x int) {
	render := Renderer(s).(*scrollBarAreaRenderer)
	barWidth := render.horizontalBarWidth()
	scrollWidth := s.scroll.Size().Width
	maxX := scrollWidth - barWidth

	if x < 0 {
		x = 0
	} else if x > maxX {
		x = maxX
	}

	ratio := float32(x) / float32(maxX)
	s.scroll.Offset.X = int(ratio * float32(s.scroll.Content.Size().Width-scrollWidth))

	Refresh(s.scroll)
}

func (s *scrollBarArea) moveVerticalBar(y int) {
	render := Renderer(s).(*scrollBarAreaRenderer)
	barHeight := render.verticalBarHeight()
	scrollHeight := s.scroll.Size().Height
	maxY := scrollHeight - barHeight

	if y < 0 {
		y = 0
	} else if y > maxY {
		y = maxY
	}

	ratio := float32(y) / float32(maxY)
	s.scroll.Offset.Y = int(ratio * float32(s.scroll.Content.Size().Height-scrollHeight))

	Refresh(s.scroll)
}

func newScrollBarArea(scroll *ScrollContainer, orientation scrollBarOrientation) *scrollBarArea {
	return &scrollBarArea{scroll: scroll, orientation: orientation}
}

type scrollRenderer struct {
	scroll                  *ScrollContainer
	vertArea                *scrollBarArea
	horizArea               *scrollBarArea
	leftShadow, rightShadow fyne.CanvasObject
	topShadow, bottomShadow fyne.CanvasObject

	objects []fyne.CanvasObject
}

func (s *scrollRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (s *scrollRenderer) Destroy() {
}

func (s *scrollRenderer) Layout(size fyne.Size) {
	// The scroll bar needs to be resized and moved on the far right
	s.horizArea.Resize(fyne.NewSize(size.Width, s.horizArea.MinSize().Height))
	s.vertArea.Resize(fyne.NewSize(s.vertArea.MinSize().Width, size.Height))

	s.horizArea.Move(fyne.NewPos(0, s.scroll.Size().Height-s.horizArea.Size().Height))
	s.vertArea.Move(fyne.NewPos(s.scroll.Size().Width-s.vertArea.Size().Width, 0))

	s.leftShadow.Resize(fyne.NewSize(0, size.Height))
	s.rightShadow.Resize(fyne.NewSize(0, size.Height))
	s.topShadow.Resize(fyne.NewSize(size.Width, 0))
	s.bottomShadow.Resize(fyne.NewSize(size.Width, 0))

	s.rightShadow.Move(fyne.NewPos(s.scroll.size.Width, 0))
	s.bottomShadow.Move(fyne.NewPos(0, s.scroll.size.Height))

	c := s.scroll.Content
	c.Resize(c.MinSize().Union(size))

	s.updatePosition()
}

func (s *scrollRenderer) MinSize() fyne.Size {
	return fyne.NewSize(25, 25) // TODO consider the smallest useful scroll view?
}

func (s *scrollRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *scrollRenderer) Refresh() {
	s.Layout(s.scroll.Size())
}

func (s *scrollRenderer) calculateScrollPosition(contentSize int, scrollSize int, offset int, area *scrollBarArea) {
	if contentSize <= scrollSize {
		if area == s.horizArea {
			s.scroll.Offset.X = 0
		} else {
			s.scroll.Offset.Y = 0
		}
		area.Hide()
	} else if s.scroll.Visible() {
		area.Show()
		if contentSize-offset < scrollSize {
			if area == s.horizArea {
				s.scroll.Offset.X = contentSize - scrollSize
			} else {
				s.scroll.Offset.Y = contentSize - scrollSize
			}
		}
	}
}

func (s *scrollRenderer) calculateShadows(offset int, contentSize int, scrollSize int, shadowStart fyne.CanvasObject, shadowEnd fyne.CanvasObject) {
	if !s.scroll.Visible() {
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

func (s *scrollRenderer) updatePosition() {
	scrollWidth := s.scroll.Size().Width
	contentWidth := s.scroll.Content.Size().Width
	scrollHeight := s.scroll.Size().Height
	contentHeight := s.scroll.Content.Size().Height
	s.calculateScrollPosition(contentWidth, scrollWidth, s.scroll.Offset.X, s.horizArea)
	s.calculateScrollPosition(contentHeight, scrollHeight, s.scroll.Offset.Y, s.vertArea)

	s.scroll.Content.Move(fyne.NewPos(-s.scroll.Offset.X, -s.scroll.Offset.Y))
	canvas.Refresh(s.scroll.Content)

	s.calculateShadows(s.scroll.Offset.X, contentWidth, scrollWidth, s.leftShadow, s.rightShadow)
	s.calculateShadows(s.scroll.Offset.Y, contentHeight, scrollHeight, s.topShadow, s.bottomShadow)

	Renderer(s.vertArea).Layout(s.scroll.size)
	Renderer(s.horizArea).Layout(s.scroll.size)
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
	s.impl = s
	s.hbar = newScrollBarArea(s, scrollBarOrientationHorizontal)
	s.vbar = newScrollBarArea(s, scrollBarOrientationVertical)
	leftShadow := newShadow(shadowRight, theme.Padding()*2)
	rightShadow := newShadow(shadowLeft, theme.Padding()*2)
	topShadow := newShadow(shadowBottom, theme.Padding()*2)
	bottomShadow := newShadow(shadowTop, theme.Padding()*2)
	return &scrollRenderer{
		objects:      []fyne.CanvasObject{s.Content, s.hbar, s.vbar, topShadow, bottomShadow},
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
	if s.Content.Size().Width <= s.Size().Width &&
		s.Content.Size().Height <= s.Size().Height {
		return
	}
	if s.hbar.Visible() && !s.vbar.Visible() {
		s.Offset.X -= ev.DeltaY
		s.Offset.Y -= ev.DeltaX
	} else {
		s.Offset.X -= ev.DeltaX
		s.Offset.Y -= ev.DeltaY
	}
	if s.Offset.X < 0 {
		s.Offset.X = 0
	} else if s.Offset.X+s.Size().Width >= s.Content.Size().Width {
		s.Offset.X = s.Content.Size().Width - s.Size().Width
	}
	if s.Offset.Y < 0 {
		s.Offset.Y = 0
	} else if s.Offset.Y+s.Size().Height >= s.Content.Size().Height {
		s.Offset.Y = s.Content.Size().Height - s.Size().Height
	}

	Refresh(s)
}

// NewScrollContainer creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed objects.
func NewScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	s := &ScrollContainer{Content: content}
	s.ExtendBaseWidget(s)
	return s
}
