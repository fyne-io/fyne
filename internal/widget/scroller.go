package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

// ScrollDirection represents the directions in which a Scroll can scroll its child content.
type ScrollDirection int

// Constants for valid values of ScrollDirection.
const (
	// ScrollBoth supports horizontal and vertical scrolling.
	ScrollBoth ScrollDirection = iota
	// ScrollHorizontalOnly specifies the scrolling should only happen left to right.
	ScrollHorizontalOnly
	// ScrollVerticalOnly specifies the scrolling should only happen top to bottom.
	ScrollVerticalOnly
	// ScrollNone turns off scrolling for this container.
	//
	// Since: 2.0
	ScrollNone
)

type scrollBarOrientation int

// We default to vertical as 0 due to that being the original orientation offered
const (
	scrollBarOrientationVertical   scrollBarOrientation = 0
	scrollBarOrientationHorizontal scrollBarOrientation = 1
	scrollContainerMinSize                              = float32(32) // TODO consider the smallest useful scroll view?
)

type scrollBarRenderer struct {
	BaseRenderer
	scrollBar  *scrollBar
	background *canvas.Rectangle
	minSize    fyne.Size
}

func (r *scrollBarRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
}

func (r *scrollBarRenderer) MinSize() fyne.Size {
	return r.minSize
}

func (r *scrollBarRenderer) Refresh() {
	r.background.FillColor = theme.ScrollBarColor()
	r.background.Refresh()
}

var _ desktop.Hoverable = (*scrollBar)(nil)
var _ fyne.Draggable = (*scrollBar)(nil)

type scrollBar struct {
	Base
	area            *scrollBarArea
	draggedDistance float32
	dragStart       float32
	isDragged       bool
	orientation     scrollBarOrientation
}

func (b *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.ScrollBarColor())
	r := &scrollBarRenderer{
		scrollBar:  b,
		background: background,
	}
	r.SetObjects([]fyne.CanvasObject{background})
	return r
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
		b.draggedDistance += e.Dragged.DX
	case scrollBarOrientationVertical:
		b.draggedDistance += e.Dragged.DY
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
	BaseRenderer
	area *scrollBarArea
	bar  *scrollBar
}

func (r *scrollBarAreaRenderer) Layout(_ fyne.Size) {
	var barHeight, barWidth, barX, barY float32
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
	min := theme.ScrollBarSize()
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

func (r *scrollBarAreaRenderer) barSizeAndOffset(contentOffset, contentLength, scrollLength float32) (length, width, lengthOffset, widthOffset float32) {
	if scrollLength < contentLength {
		portion := scrollLength / contentLength
		length = float32(int(scrollLength)) * portion
		if length < theme.ScrollBarSize() {
			length = theme.ScrollBarSize()
		}
	} else {
		length = scrollLength
	}
	if contentOffset != 0 {
		lengthOffset = (scrollLength - length) * (contentOffset / (contentLength - scrollLength))
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
	Base

	isLarge     bool
	scroll      *Scroll
	orientation scrollBarOrientation
}

func (a *scrollBarArea) CreateRenderer() fyne.WidgetRenderer {
	bar := newScrollBar(a)
	return &scrollBarAreaRenderer{BaseRenderer: NewBaseRenderer([]fyne.CanvasObject{bar}), area: a, bar: bar}
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

func (a *scrollBarArea) moveBar(offset float32, barSize fyne.Size) {
	oldX := a.scroll.Offset.X
	oldY := a.scroll.Offset.Y
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		a.scroll.Offset.X = a.computeScrollOffset(barSize.Width, offset, a.scroll.Size().Width, a.scroll.Content.Size().Width)
	default:
		a.scroll.Offset.Y = a.computeScrollOffset(barSize.Height, offset, a.scroll.Size().Height, a.scroll.Content.Size().Height)
	}
	if f := a.scroll.OnScrolled; f != nil && (a.scroll.Offset.X != oldX || a.scroll.Offset.Y != oldY) {
		f(a.scroll.Offset)
	}
	a.scroll.refreshWithoutOffsetUpdate()
}

func (a *scrollBarArea) computeScrollOffset(length, offset, scrollLength, contentLength float32) float32 {
	maxOffset := scrollLength - length
	if offset < 0 {
		offset = 0
	} else if offset > maxOffset {
		offset = maxOffset
	}
	ratio := offset / maxOffset
	scrollOffset := ratio * (contentLength - scrollLength)
	return scrollOffset
}

func newScrollBarArea(scroll *Scroll, orientation scrollBarOrientation) *scrollBarArea {
	a := &scrollBarArea{scroll: scroll, orientation: orientation}
	a.ExtendBaseWidget(a)
	return a
}

type scrollContainerRenderer struct {
	BaseRenderer
	scroll                  *Scroll
	vertArea                *scrollBarArea
	horizArea               *scrollBarArea
	leftShadow, rightShadow *Shadow
	topShadow, bottomShadow *Shadow
	oldMinSize              fyne.Size
}

func (r *scrollContainerRenderer) layoutBars(size fyne.Size) {
	if r.scroll.Direction == ScrollVerticalOnly || r.scroll.Direction == ScrollBoth {
		r.vertArea.Resize(fyne.NewSize(r.vertArea.MinSize().Width, size.Height))
		r.vertArea.Move(fyne.NewPos(r.scroll.Size().Width-r.vertArea.Size().Width, 0))
		r.topShadow.Resize(fyne.NewSize(size.Width, 0))
		r.bottomShadow.Resize(fyne.NewSize(size.Width, 0))
		r.bottomShadow.Move(fyne.NewPos(0, r.scroll.size.Height))
	}

	if r.scroll.Direction == ScrollHorizontalOnly || r.scroll.Direction == ScrollBoth {
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
		r.BaseRenderer.Objects()[0] = r.scroll.Content
	}
	if r.oldMinSize == r.scroll.Content.MinSize() && r.oldMinSize == r.scroll.Content.Size() &&
		(r.scroll.Size().Width <= r.oldMinSize.Width && r.scroll.Size().Height <= r.oldMinSize.Height) {
		r.layoutBars(r.scroll.Size())
		return
	}

	r.oldMinSize = r.scroll.Content.MinSize()
	r.Layout(r.scroll.Size())
}

func (r *scrollContainerRenderer) handleAreaVisibility(contentSize, scrollSize float32, area *scrollBarArea) {
	if contentSize <= scrollSize {
		area.Hide()
	} else if r.scroll.Visible() {
		area.Show()
	}
}

func (r *scrollContainerRenderer) handleShadowVisibility(offset, contentSize, scrollSize float32, shadowStart fyne.CanvasObject, shadowEnd fyne.CanvasObject) {
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
	if r.scroll.Content == nil {
		return
	}
	scrollSize := r.scroll.Size()
	contentSize := r.scroll.Content.Size()

	r.scroll.Content.Move(fyne.NewPos(-r.scroll.Offset.X, -r.scroll.Offset.Y))

	if r.scroll.Direction == ScrollVerticalOnly || r.scroll.Direction == ScrollBoth {
		r.handleAreaVisibility(contentSize.Height, scrollSize.Height, r.vertArea)
		r.handleShadowVisibility(r.scroll.Offset.Y, contentSize.Height, scrollSize.Height, r.topShadow, r.bottomShadow)
		cache.Renderer(r.vertArea).Layout(r.scroll.size)
	} else {
		r.vertArea.Hide()
		r.topShadow.Hide()
		r.bottomShadow.Hide()
	}
	if r.scroll.Direction == ScrollHorizontalOnly || r.scroll.Direction == ScrollBoth {
		r.handleAreaVisibility(contentSize.Width, scrollSize.Width, r.horizArea)
		r.handleShadowVisibility(r.scroll.Offset.X, contentSize.Width, scrollSize.Width, r.leftShadow, r.rightShadow)
		cache.Renderer(r.horizArea).Layout(r.scroll.size)
	} else {
		r.horizArea.Hide()
		r.leftShadow.Hide()
		r.rightShadow.Hide()
	}

	if r.scroll.Direction != ScrollHorizontalOnly {
		canvas.Refresh(r.vertArea) // this is required to force the canvas to update, we have no "Redraw()"
	} else {
		canvas.Refresh(r.horizArea) // this is required like above but if we are horizontal
	}
}

// Scroll defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
type Scroll struct {
	Base
	minSize   fyne.Size
	Direction ScrollDirection
	Content   fyne.CanvasObject
	Offset    fyne.Position
	// OnScrolled can be set to be notified when the Scroll has changed position.
	// You should not update the Scroll.Offset from this method.
	//
	// Since: 2.0
	OnScrolled func(fyne.Position)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Scroll) CreateRenderer() fyne.WidgetRenderer {
	scr := &scrollContainerRenderer{
		BaseRenderer: NewBaseRenderer([]fyne.CanvasObject{s.Content}),
		scroll:       s,
	}
	scr.vertArea = newScrollBarArea(s, scrollBarOrientationVertical)
	scr.topShadow = NewShadow(ShadowBottom, SubmergedContentLevel)
	scr.bottomShadow = NewShadow(ShadowTop, SubmergedContentLevel)
	scr.horizArea = newScrollBarArea(s, scrollBarOrientationHorizontal)
	scr.leftShadow = NewShadow(ShadowRight, SubmergedContentLevel)
	scr.rightShadow = NewShadow(ShadowLeft, SubmergedContentLevel)
	scr.SetObjects(append(scr.Objects(), scr.topShadow, scr.bottomShadow, scr.leftShadow, scr.rightShadow,
		scr.vertArea, scr.horizArea))
	scr.updatePosition()

	return scr
}

// ScrollToBottom will scroll content to container bottom - to show latest info which end user just added
func (s *Scroll) ScrollToBottom() {
	s.scrollBy(0, -1*(s.Content.MinSize().Height-s.Size().Height-s.Offset.Y))
	s.Refresh()
}

// ScrollToTop will scroll content to container top
func (s *Scroll) ScrollToTop() {
	s.scrollBy(0, -s.Offset.Y)
}

// DragEnd will stop scrolling on mobile has stopped
func (s *Scroll) DragEnd() {
}

// Dragged will scroll on any drag - bar or otherwise - for mobile
func (s *Scroll) Dragged(e *fyne.DragEvent) {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}

	if s.updateOffset(e.Dragged.DX, e.Dragged.DY) {
		s.refreshWithoutOffsetUpdate()
	}
}

// MinSize returns the smallest size this widget can shrink to
func (s *Scroll) MinSize() fyne.Size {
	min := fyne.NewSize(scrollContainerMinSize, scrollContainerMinSize).Max(s.minSize)
	switch s.Direction {
	case ScrollHorizontalOnly:
		min.Height = fyne.Max(min.Height, s.Content.MinSize().Height)
	case ScrollVerticalOnly:
		min.Width = fyne.Max(min.Width, s.Content.MinSize().Width)
	case ScrollNone:
		return s.Content.MinSize()
	}
	return min
}

// SetMinSize specifies a minimum size for this scroll container.
// If the specified size is larger than the content size then scrolling will not be enabled
// This can be helpful to appear larger than default if the layout is collapsing this widget.
func (s *Scroll) SetMinSize(size fyne.Size) {
	s.minSize = size
}

// Refresh causes this widget to be redrawn in it's current state
func (s *Scroll) Refresh() {
	s.updateOffset(0, 0)
	s.refreshWithoutOffsetUpdate()
}

// Resize is called when this scroller should change size. We refresh to ensure the scroll bars are updated.
func (s *Scroll) Resize(sz fyne.Size) {
	if sz == s.size {
		return
	}

	s.Base.Resize(sz)
	s.Refresh()
}

func (s *Scroll) refreshWithoutOffsetUpdate() {
	s.Base.Refresh()
}

// Scrolled is called when an input device triggers a scroll event
func (s *Scroll) Scrolled(ev *fyne.ScrollEvent) {
	s.scrollBy(ev.Scrolled.DX, ev.Scrolled.DY)
}

func (s *Scroll) scrollBy(dx, dy float32) {
	if s.Size().Width < s.Content.MinSize().Width && s.Size().Height >= s.Content.MinSize().Height && dx == 0 {
		dx, dy = dy, dx
	}
	if s.updateOffset(dx, dy) {
		s.refreshWithoutOffsetUpdate()
	}
}

func (s *Scroll) updateOffset(deltaX, deltaY float32) bool {
	if s.Content.Size().Width <= s.Size().Width && s.Content.Size().Height <= s.Size().Height {
		if s.Offset.X != 0 || s.Offset.Y != 0 {
			s.Offset.X = 0
			s.Offset.Y = 0
			return true
		}
		return false
	}
	oldX := s.Offset.X
	oldY := s.Offset.Y
	s.Offset.X = computeOffset(s.Offset.X, -deltaX, s.Size().Width, s.Content.MinSize().Width)
	s.Offset.Y = computeOffset(s.Offset.Y, -deltaY, s.Size().Height, s.Content.MinSize().Height)
	if f := s.OnScrolled; f != nil && (s.Offset.X != oldX || s.Offset.Y != oldY) {
		f(s.Offset)
	}
	return true
}

func computeOffset(start, delta, outerWidth, innerWidth float32) float32 {
	offset := start + delta
	if offset+outerWidth >= innerWidth {
		offset = innerWidth - outerWidth
	}
	if offset < 0 {
		offset = 0
	}
	return offset
}

// NewScroll creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed object.
func NewScroll(content fyne.CanvasObject) *Scroll {
	s := newScrollContainerWithDirection(ScrollBoth, content)
	s.ExtendBaseWidget(s)
	return s
}

// NewHScroll create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Width to be smaller than that of the passed object.
func NewHScroll(content fyne.CanvasObject) *Scroll {
	s := newScrollContainerWithDirection(ScrollHorizontalOnly, content)
	s.ExtendBaseWidget(s)
	return s
}

// NewVScroll create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Height to be smaller than that of the passed object.
func NewVScroll(content fyne.CanvasObject) *Scroll {
	s := newScrollContainerWithDirection(ScrollVerticalOnly, content)
	s.ExtendBaseWidget(s)
	return s
}

func newScrollContainerWithDirection(direction ScrollDirection, content fyne.CanvasObject) *Scroll {
	s := &Scroll{
		Direction: direction,
		Content:   content,
	}
	s.ExtendBaseWidget(s)
	return s
}
