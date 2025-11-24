package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

// ScrollDirection represents the directions in which a Scroll can scroll its child content.
type ScrollDirection = fyne.ScrollDirection

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

	// what fraction of the page to scroll when tapping on the scroll bar area
	pageScrollFraction = float32(0.95)
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
	th := theme.CurrentForWidget(r.scrollBar)
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.background.FillColor = th.Color(theme.ColorNameScrollBar, v)
	r.background.CornerRadius = th.Size(theme.SizeNameScrollBarRadius)
	r.background.Refresh()
}

var (
	_ desktop.Hoverable = (*scrollBar)(nil)
	_ fyne.Draggable    = (*scrollBar)(nil)
)

type scrollBar struct {
	Base
	area            *scrollBarArea
	draggedDistance float32
	dragStart       float32
	orientation     scrollBarOrientation
}

func (b *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	th := theme.CurrentForWidget(b)
	v := fyne.CurrentApp().Settings().ThemeVariant()

	background := canvas.NewRectangle(th.Color(theme.ColorNameScrollBar, v))
	background.CornerRadius = th.Size(theme.SizeNameScrollBarRadius)
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
	b.area.isDragging = false

	if fyne.CurrentDevice().IsMobile() {
		b.area.MouseOut()
		return
	}
	b.area.Refresh()
}

func (b *scrollBar) Dragged(e *fyne.DragEvent) {
	if !b.area.isDragging {
		b.area.isDragging = true
		b.area.MouseIn(nil)

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

func (a *scrollBarArea) isLarge() bool {
	return a.isMouseIn || a.isDragging
}

type scrollBarAreaRenderer struct {
	BaseRenderer
	area       *scrollBarArea
	bar        *scrollBar
	background *canvas.Rectangle
}

func (r *scrollBarAreaRenderer) Layout(size fyne.Size) {
	r.layoutWithTheme(theme.CurrentForWidget(r.area), size)
}

func (r *scrollBarAreaRenderer) layoutWithTheme(th fyne.Theme, size fyne.Size) {
	var barHeight, barWidth, barX, barY float32
	var bkgHeight, bkgWidth, bkgX, bkgY float32
	switch r.area.orientation {
	case scrollBarOrientationHorizontal:
		barWidth, barHeight, barX, barY = r.barSizeAndOffset(th, r.area.scroll.Offset.X, r.area.scroll.Content.Size().Width, r.area.scroll.Size().Width)
		r.area.barLeadingEdge = barX
		r.area.barTrailingEdge = barX + barWidth
		bkgWidth, bkgHeight, bkgX, bkgY = size.Width, barHeight, 0, barY
	default:
		barHeight, barWidth, barY, barX = r.barSizeAndOffset(th, r.area.scroll.Offset.Y, r.area.scroll.Content.Size().Height, r.area.scroll.Size().Height)
		r.area.barLeadingEdge = barY
		r.area.barTrailingEdge = barY + barHeight
		bkgWidth, bkgHeight, bkgX, bkgY = barWidth, size.Height, barX, 0
	}
	r.bar.Move(fyne.NewPos(barX, barY))
	r.bar.Resize(fyne.NewSize(barWidth, barHeight))
	r.background.Move(fyne.NewPos(bkgX, bkgY))
	r.background.Resize(fyne.NewSize(bkgWidth, bkgHeight))
}

func (r *scrollBarAreaRenderer) MinSize() fyne.Size {
	th := theme.CurrentForWidget(r.area)

	barSize := th.Size(theme.SizeNameScrollBar)
	min := barSize
	if !r.area.isLarge() {
		min = th.Size(theme.SizeNameScrollBarSmall) * 2
	}
	switch r.area.orientation {
	case scrollBarOrientationHorizontal:
		return fyne.NewSize(barSize, min)
	default:
		return fyne.NewSize(min, barSize)
	}
}

func (r *scrollBarAreaRenderer) Refresh() {
	th := theme.CurrentForWidget(r.area)
	r.bar.Refresh()
	r.background.FillColor = th.Color(theme.ColorNameScrollBarBackground, fyne.CurrentApp().Settings().ThemeVariant())
	r.background.Hidden = !r.area.isLarge()
	r.layoutWithTheme(th, r.area.Size())
	canvas.Refresh(r.bar)
	canvas.Refresh(r.background)
}

func (r *scrollBarAreaRenderer) barSizeAndOffset(th fyne.Theme, contentOffset, contentLength, scrollLength float32) (length, width, lengthOffset, widthOffset float32) {
	scrollBarSize := th.Size(theme.SizeNameScrollBar)
	if scrollLength < contentLength {
		portion := scrollLength / contentLength
		length = float32(int(scrollLength)) * portion
		length = fyne.Max(length, scrollBarSize)
	} else {
		length = scrollLength
	}
	if contentOffset != 0 {
		lengthOffset = (scrollLength - length) * (contentOffset / (contentLength - scrollLength))
	}
	if r.area.isLarge() {
		width = scrollBarSize
	} else {
		widthOffset = th.Size(theme.SizeNameScrollBarSmall)
		width = widthOffset
	}
	return length, width, lengthOffset, widthOffset
}

var (
	_ desktop.Hoverable = (*scrollBarArea)(nil)
	_ fyne.Tappable     = (*scrollBarArea)(nil)
)

type scrollBarArea struct {
	Base

	isDragging  bool
	isMouseIn   bool
	scroll      *Scroll
	bar         *scrollBar
	orientation scrollBarOrientation

	// updated from renderer Layout
	// coordinates Y in vertical orientation, X in horizontal
	barLeadingEdge  float32
	barTrailingEdge float32
}

func (a *scrollBarArea) CreateRenderer() fyne.WidgetRenderer {
	th := theme.CurrentForWidget(a)
	v := fyne.CurrentApp().Settings().ThemeVariant()
	a.bar = newScrollBar(a)
	background := canvas.NewRectangle(th.Color(theme.ColorNameScrollBarBackground, v))
	background.Hidden = !a.isLarge()
	return &scrollBarAreaRenderer{BaseRenderer: NewBaseRenderer([]fyne.CanvasObject{background, a.bar}), area: a, bar: a.bar, background: background}
}

func (a *scrollBarArea) Tapped(e *fyne.PointEvent) {
	if isScrollerPageOnTap() {
		a.scrollFullPageOnTap(e)
		return
	}

	// scroll to tapped position
	barSize := a.bar.Size()
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		if e.Position.X < a.barLeadingEdge || e.Position.X > a.barTrailingEdge {
			a.moveBar(fyne.Max(0, e.Position.X-barSize.Width/2), barSize)
		}
	case scrollBarOrientationVertical:
		if e.Position.Y < a.barLeadingEdge || e.Position.Y > a.barTrailingEdge {
			a.moveBar(fyne.Max(0, e.Position.Y-barSize.Height/2), a.bar.Size())
		}
	}
}

func (a *scrollBarArea) scrollFullPageOnTap(e *fyne.PointEvent) {
	// when tapping above/below or left/right of the bar, scroll the content
	// nearly a full page (pageScrollFraction) up/down or left/right, respectively
	newOffset := a.scroll.Offset
	switch a.orientation {
	case scrollBarOrientationHorizontal:
		if e.Position.X < a.barLeadingEdge {
			newOffset.X = fyne.Max(0, newOffset.X-a.scroll.Size().Width*pageScrollFraction)
		} else if e.Position.X > a.barTrailingEdge {
			viewWid := a.scroll.Size().Width
			newOffset.X = fyne.Min(a.scroll.Content.Size().Width-viewWid, newOffset.X+viewWid*pageScrollFraction)
		}
	default:
		if e.Position.Y < a.barLeadingEdge {
			newOffset.Y = fyne.Max(0, newOffset.Y-a.scroll.Size().Height*pageScrollFraction)
		} else if e.Position.Y > a.barTrailingEdge {
			viewHt := a.scroll.Size().Height
			newOffset.Y = fyne.Min(a.scroll.Content.Size().Height-viewHt, newOffset.Y+viewHt*pageScrollFraction)
		}
	}
	if newOffset == a.scroll.Offset {
		return
	}

	a.scroll.Offset = newOffset
	if f := a.scroll.OnScrolled; f != nil {
		f(a.scroll.Offset)
	}
	a.scroll.refreshWithoutOffsetUpdate()
}

func (a *scrollBarArea) MouseIn(*desktop.MouseEvent) {
	a.isMouseIn = true
	a.scroll.refreshBars()
}

func (a *scrollBarArea) MouseMoved(*desktop.MouseEvent) {
}

func (a *scrollBarArea) MouseOut() {
	a.isMouseIn = false
	if a.isDragging {
		return
	}

	a.scroll.refreshBars()
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
	scrollerSize := r.scroll.Size()
	if r.scroll.Direction == ScrollVerticalOnly || r.scroll.Direction == ScrollBoth {
		r.vertArea.Resize(fyne.NewSize(r.vertArea.MinSize().Width, size.Height))
		r.vertArea.Move(fyne.NewPos(scrollerSize.Width-r.vertArea.Size().Width, 0))
		r.topShadow.Resize(fyne.NewSize(size.Width, 0))
		r.bottomShadow.Resize(fyne.NewSize(size.Width, 0))
		r.bottomShadow.Move(fyne.NewPos(0, scrollerSize.Height))
	}

	if r.scroll.Direction == ScrollHorizontalOnly || r.scroll.Direction == ScrollBoth {
		r.horizArea.Resize(fyne.NewSize(size.Width, r.horizArea.MinSize().Height))
		r.horizArea.Move(fyne.NewPos(0, scrollerSize.Height-r.horizArea.Size().Height))
		r.leftShadow.Resize(fyne.NewSize(0, size.Height))
		r.rightShadow.Resize(fyne.NewSize(0, size.Height))
		r.rightShadow.Move(fyne.NewPos(scrollerSize.Width, 0))
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
	r.horizArea.Refresh()
	r.vertArea.Refresh()
	r.leftShadow.Refresh()
	r.topShadow.Refresh()
	r.rightShadow.Refresh()
	r.bottomShadow.Refresh()

	if len(r.BaseRenderer.Objects()) == 0 || r.BaseRenderer.Objects()[0] != r.scroll.Content {
		// push updated content object to baseRenderer
		r.BaseRenderer.Objects()[0] = r.scroll.Content
	}
	size := r.scroll.Size()
	newMin := r.scroll.Content.MinSize()
	if r.oldMinSize == newMin && r.oldMinSize == r.scroll.Content.Size() &&
		(size.Width <= r.oldMinSize.Width && size.Height <= r.oldMinSize.Height) {
		r.layoutBars(size)
		return
	}

	r.oldMinSize = newMin
	r.Layout(size)
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
		cache.Renderer(r.vertArea).Layout(scrollSize)
	} else {
		r.vertArea.Hide()
		r.topShadow.Hide()
		r.bottomShadow.Hide()
	}
	if r.scroll.Direction == ScrollHorizontalOnly || r.scroll.Direction == ScrollBoth {
		r.handleAreaVisibility(contentSize.Width, scrollSize.Width, r.horizArea)
		r.handleShadowVisibility(r.scroll.Offset.X, contentSize.Width, scrollSize.Width, r.leftShadow, r.rightShadow)
		cache.Renderer(r.horizArea).Layout(scrollSize)
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
	OnScrolled func(fyne.Position) `json:"-"`
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
	s.refreshBars()
}

// ScrollToTop will scroll content to container top
func (s *Scroll) ScrollToTop() {
	s.ScrollToOffset(fyne.Position{})
	s.refreshBars()
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
	s.refreshBars()

	if s.Content != nil {
		s.Content.Refresh()
	}
}

// Resize is called when this scroller should change size. We refresh to ensure the scroll bars are updated.
func (s *Scroll) Resize(sz fyne.Size) {
	if sz == s.Size() {
		return
	}

	s.Base.Resize(sz)
	s.refreshBars()
}

// ScrollToOffset will update the location of the content of this scroll container.
//
// Since: 2.6
func (s *Scroll) ScrollToOffset(p fyne.Position) {
	if s.Offset == p {
		return
	}

	s.Offset = p
	s.refreshBars()
}

func (s *Scroll) refreshWithoutOffsetUpdate() {
	s.Base.Refresh()
}

// Scrolled is called when an input device triggers a scroll event
func (s *Scroll) Scrolled(ev *fyne.ScrollEvent) {
	if s.Direction != ScrollNone {
		s.scrollBy(ev.Scrolled.DX, ev.Scrolled.DY)
	}
}

func (s *Scroll) refreshBars() {
	s.updateOffset(0, 0)
	s.refreshWithoutOffsetUpdate()
}

func (s *Scroll) scrollBy(dx, dy float32) {
	min := s.Content.MinSize()
	size := s.Size()
	if size.Width < min.Width && size.Height >= min.Height && dx == 0 {
		dx, dy = dy, dx
	}
	if s.updateOffset(dx, dy) {
		s.refreshWithoutOffsetUpdate()
	}
}

func (s *Scroll) updateOffset(deltaX, deltaY float32) bool {
	size := s.Size()
	contentSize := s.Content.Size()
	if contentSize.Width <= size.Width && contentSize.Height <= size.Height {
		if s.Offset.X != 0 || s.Offset.Y != 0 {
			s.Offset.X = 0
			s.Offset.Y = 0
			return true
		}
		return false
	}
	oldX := s.Offset.X
	oldY := s.Offset.Y
	min := s.Content.MinSize()
	s.Offset.X = computeOffset(s.Offset.X, -deltaX, size.Width, min.Width)
	s.Offset.Y = computeOffset(s.Offset.Y, -deltaY, size.Height, min.Height)

	moved := s.Offset.X != oldX || s.Offset.Y != oldY
	if f := s.OnScrolled; f != nil && moved {
		f(s.Offset)
	}
	return moved
}

func computeOffset(start, delta, outerWidth, innerWidth float32) float32 {
	offset := start + delta
	if offset+outerWidth >= innerWidth {
		offset = innerWidth - outerWidth
	}

	return fyne.Max(offset, 0)
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
