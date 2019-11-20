package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

type scrollBarRenderer struct {
	scrollBar *scrollBar

	color   color.Color
	minSize fyne.Size
	objects []fyne.CanvasObject
}

func (r *scrollBarRenderer) ApplyTheme() {
	r.color = theme.ScrollBarColor()
}

func (r *scrollBarRenderer) BackgroundColor() color.Color {
	return r.color
}

func (r *scrollBarRenderer) Destroy() {
}

func (r *scrollBarRenderer) Layout(size fyne.Size) {
}

func (r *scrollBarRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *scrollBarRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *scrollBarRenderer) Refresh() {
}

var _ desktop.Hoverable = (*scrollBar)(nil)
var _ fyne.Draggable = (*scrollBar)(nil)

type scrollBar struct {
	baseWidget
	area            *scrollBarArea
	draggedDistance int
	dragStart       int
	isDragged       bool
}

func (s *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	r := &scrollBarRenderer{scrollBar: s}
	r.ApplyTheme()
	return r
}

func (s *scrollBar) DragEnd() {
}

func (s *scrollBar) Dragged(e *fyne.DragEvent) {
	if !s.isDragged {
		s.isDragged = true
		s.dragStart = s.Position().Y
		s.draggedDistance = 0
	}
	s.draggedDistance += e.DraggedY
	s.area.moveBar(s.draggedDistance + s.dragStart)
}

func (s *scrollBar) Hide() {
	s.hide(s)
}

func (s *scrollBar) MinSize() fyne.Size {
	return s.minSize(s)
}

func (s *scrollBar) MouseIn(e *desktop.MouseEvent) {
	s.area.MouseIn(e)
}

func (s *scrollBar) MouseMoved(*desktop.MouseEvent) {
}

func (s *scrollBar) MouseOut() {
	s.area.MouseOut()
}

func (s *scrollBar) Move(pos fyne.Position) {
	s.move(pos, s)
}

func (s *scrollBar) Resize(size fyne.Size) {
	s.resize(size, s)
}

func (s *scrollBar) Show() {
	s.show(s)
}

func newScrollBar(area *scrollBarArea) *scrollBar {
	return &scrollBar{area: area}
}

type scrollBarAreaRenderer struct {
	area *scrollBarArea
	bar  *scrollBar

	objects []fyne.CanvasObject
}

func (s *scrollBarAreaRenderer) ApplyTheme() {
}

func (s *scrollBarAreaRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s *scrollBarAreaRenderer) Destroy() {
}

func (s *scrollBarAreaRenderer) Layout(size fyne.Size) {
	s.updateBarPosition()
}

func (s *scrollBarAreaRenderer) MinSize() fyne.Size {
	var minWidth int
	if s.area.isWide {
		minWidth = theme.ScrollBarSize()
	} else {
		minWidth = theme.ScrollBarSmallSize() * 2
	}
	return fyne.NewSize(minWidth, theme.ScrollBarSize())
}

func (s *scrollBarAreaRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *scrollBarAreaRenderer) Refresh() {
	s.updateBarPosition()
	canvas.Refresh(s.bar)
}

func (s *scrollBarAreaRenderer) updateBarPosition() {
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

func (s *scrollBarAreaRenderer) verticalBarHeight() int {
	portion := float32(s.area.size.Height) / float32(s.area.scroll.Content.Size().Height)
	if portion > 1.0 {
		portion = 1.0
	}

	return int(float32(s.area.size.Height) * portion)
}

var _ desktop.Hoverable = (*scrollBarArea)(nil)

type scrollBarArea struct {
	baseWidget

	isWide bool
	scroll *ScrollContainer
}

func (s *scrollBarArea) CreateRenderer() fyne.WidgetRenderer {
	bar := newScrollBar(s)
	return &scrollBarAreaRenderer{area: s, bar: bar, objects: []fyne.CanvasObject{bar}}
}

func (s *scrollBarArea) Hide() {
	s.hide(s)
}

func (s *scrollBarArea) MinSize() fyne.Size {
	return s.minSize(s)
}

func (s *scrollBarArea) MouseIn(*desktop.MouseEvent) {
	s.isWide = true
	Refresh(s.scroll)
}

func (s *scrollBarArea) MouseMoved(*desktop.MouseEvent) {
}

func (s *scrollBarArea) MouseOut() {
	s.isWide = false
	Refresh(s.scroll)
}

func (s *scrollBarArea) Move(pos fyne.Position) {
	s.move(pos, s)

}

func (s *scrollBarArea) moveBar(y int) {
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

func (s *scrollBarArea) Resize(size fyne.Size) {
	s.resize(size, s)
}

func (s *scrollBarArea) Show() {
	s.show(s)
}

func newScrollBarArea(scroll *ScrollContainer) *scrollBarArea {
	return &scrollBarArea{scroll: scroll}
}

type scrollRenderer struct {
	scroll                  *ScrollContainer
	vertArea                *scrollBarArea
	topShadow, bottomShadow fyne.CanvasObject

	objects []fyne.CanvasObject
}

func (s *scrollRenderer) ApplyTheme() {
}

func (s *scrollRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (s *scrollRenderer) Destroy() {
}

func (s *scrollRenderer) Layout(size fyne.Size) {
	// The scroll bar needs to be resized and moved on the far right
	scrollBarArea := s.vertArea
	scrollBarArea.Resize(fyne.NewSize(scrollBarArea.MinSize().Width, size.Height))
	scrollBarArea.Move(fyne.NewPos(s.scroll.Size().Width-scrollBarArea.Size().Width, 0))
	s.topShadow.Resize(fyne.NewSize(size.Width, 0))
	s.bottomShadow.Resize(fyne.NewSize(size.Width, 0))
	s.bottomShadow.Move(fyne.NewPos(0, s.scroll.size.Height))

	c := s.scroll.Content
	c.Resize(c.MinSize().Union(size))

	s.updatePosition()
}

func (s *scrollRenderer) MinSize() fyne.Size {
	// TODO determine if width or height should be respected based on a which-way-to-scroll flag
	return fyne.NewSize(s.scroll.Content.MinSize().Width, 25) // TODO consider the smallest useful scroll view?
}

func (s *scrollRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *scrollRenderer) Refresh() {
	s.Layout(s.scroll.Size())
}

func (s *scrollRenderer) updatePosition() {
	scrollHeight := s.scroll.Size().Height
	contentHeight := s.scroll.Content.Size().Height
	if contentHeight <= scrollHeight {
		s.scroll.Offset.Y = 0
		s.vertArea.Hide()
	} else if s.scroll.Visible() {
		s.vertArea.Show()
		if contentHeight-s.scroll.Offset.Y < scrollHeight {
			s.scroll.Offset.Y = contentHeight - scrollHeight
		}
	}
	s.scroll.Content.Move(fyne.NewPos(-s.scroll.Offset.X, -s.scroll.Offset.Y))
	canvas.Refresh(s.scroll.Content)

	if s.scroll.Offset.Y > 0 && s.scroll.Visible() {
		s.topShadow.Show()
	} else {
		s.topShadow.Hide()
	}
	if s.scroll.Offset.Y < contentHeight-scrollHeight && s.scroll.Visible() {
		s.bottomShadow.Show()
	} else {
		s.bottomShadow.Hide()
	}

	Renderer(s.vertArea).Layout(s.scroll.size)
}

// ScrollContainer defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
type ScrollContainer struct {
	baseWidget

	Content fyne.CanvasObject
	Offset  fyne.Position
	draggedDistance int
	vbar *scrollBarArea
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *ScrollContainer) CreateRenderer() fyne.WidgetRenderer {
	s.vbar = newScrollBarArea(s)
	topShadow := newShadow(shadowBottom, theme.Padding()*2)
	bottomShadow := newShadow(shadowTop, theme.Padding()*2)
	return &scrollRenderer{
		objects:      []fyne.CanvasObject{s.Content, s.vbar, topShadow, bottomShadow},
		scroll:       s,
		vertArea:     s.vbar,
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

	render := Renderer(s.vbar).(*scrollBarAreaRenderer)
	barHeight := render.verticalBarHeight()
	scrollHeight := s.Size().Height
	maxY := scrollHeight - barHeight

	s.draggedDistance += e.DraggedY
	if s.draggedDistance > maxY {
		s.draggedDistance = maxY
	}
	if s.draggedDistance < 0 {
		s.draggedDistance = 0
	}
	s.vbar.moveBar(s.vbar.position.Y+s.draggedDistance)
}

// Hide this widget, if it was previously visible
func (s *ScrollContainer) Hide() {
	s.hide(s)
}

// MinSize returns the smallest size this widget can shrink to
func (s *ScrollContainer) MinSize() fyne.Size {
	return s.minSize(s)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *ScrollContainer) Move(pos fyne.Position) {
	s.move(pos, s)
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *ScrollContainer) Resize(size fyne.Size) {
	s.resize(size, s)
}

// Scrolled is called when an input device triggers a scroll event
func (s *ScrollContainer) Scrolled(ev *fyne.ScrollEvent) {
	if s.Content.Size().Height <= s.Size().Height {
		return
	}

	s.Offset.Y -= ev.DeltaY

	if s.Offset.Y < 0 {
		s.Offset.Y = 0
	} else if s.Offset.Y+s.Size().Height >= s.Content.Size().Height {
		s.Offset.Y = s.Content.Size().Height - s.Size().Height
	}

	Refresh(s)
}

// Show this widget, if it was previously hidden
func (s *ScrollContainer) Show() {
	s.show(s)
}

// NewScrollContainer creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed objects.
func NewScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	return &ScrollContainer{Content: content}
}
