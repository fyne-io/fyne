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
	bar       *canvas.Rectangle

	objects []fyne.CanvasObject
}

func (s *scrollBarRenderer) Layout(size fyne.Size) {
	s.updateBarPosition()
	canvas.Refresh(s.bar)
}

func (s *scrollBarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(s.minWidth(), theme.ScrollBarSize())
}

func (s *scrollBarRenderer) minWidth() int {
	return s.scrollBar.minWidth
}

func (s *scrollBarRenderer) Refresh() {
	s.updateBarPosition()
}

func (s *scrollBarRenderer) ApplyTheme() {
	s.bar.FillColor = theme.ScrollBarColor()
}

func (s *scrollBarRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (s *scrollBarRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *scrollBarRenderer) Destroy() {

}

func (s *scrollBarRenderer) barSizeVertical() fyne.Size {
	portion := float32(s.scrollBar.size.Height) / float32(s.scrollBar.scroll.Content.Size().Height)
	if portion > 1.0 {
		portion = 1.0
	}

	barHeight := int(float32(s.scrollBar.size.Height) * portion)
	return fyne.NewSize(s.minWidth(), barHeight)
}

func (s *scrollBarRenderer) updateBarPosition() {
	barSize := s.barSizeVertical()
	barRatio := float32(0.0)
	if s.scrollBar.scroll.Offset.Y != 0 {
		barRatio = float32(s.scrollBar.scroll.Offset.Y) / float32(s.scrollBar.scroll.Content.Size().Height-s.scrollBar.scroll.Size().Height)
	}
	barOff := int(float32(s.scrollBar.scroll.size.Height-barSize.Height) * barRatio)

	s.bar.Resize(barSize)
	s.bar.Move(fyne.NewPos(0, barOff))
}

var _ desktop.Hoverable = (*scrollBar)(nil)

type scrollBar struct {
	baseWidget

	minWidth int
	scroll   *ScrollContainer
}

func (s *scrollBar) Dragged(ev *fyne.DragEvent) {
	render := Renderer(s).(*scrollBarRenderer)
	barHeight := render.barSizeVertical().Height
	barTop := render.bar.Position().Y
	barBottom := barTop + barHeight

	// The point clicked is outside the bar rectangle
	if ev.Position.Y < barTop || ev.Position.Y > barBottom {
		return
	}

	dragRatio := float32(ev.DraggedY) / float32(s.scroll.size.Height-barHeight)
	addiotionalOffset := int(dragRatio * float32(s.scroll.Content.Size().Height-s.scroll.Size().Height))
	s.scroll.Offset.Y = s.scroll.Offset.Y + addiotionalOffset
	if s.scroll.Offset.Y < 0 {
		s.scroll.Offset.Y = 0
	}

	Refresh(s.scroll)
}

func (s *scrollBar) Resize(size fyne.Size) {
	s.resize(size, s)
}

func (s *scrollBar) Move(pos fyne.Position) {
	s.move(pos, s)
}

func (s *scrollBar) MinSize() fyne.Size {
	return s.minSize(s)
}

func (s *scrollBar) Show() {
	s.show(s)
}

func (s *scrollBar) Hide() {
	s.hide(s)
}

func (s *scrollBar) CreateRenderer() fyne.WidgetRenderer {
	bar := canvas.NewRectangle(theme.ScrollBarColor())
	return &scrollBarRenderer{scrollBar: s, bar: bar, objects: []fyne.CanvasObject{bar}}
}

func (s *scrollBar) MouseIn(*desktop.MouseEvent) {
	s.minWidth = theme.ScrollBarSize()
	Refresh(s.scroll)
}

func (s *scrollBar) MouseMoved(*desktop.MouseEvent) {
}

func (s *scrollBar) MouseOut() {
	s.minWidth = theme.ScrollBarSmallSize()
	Refresh(s.scroll)
}

func newScrollBar(scroll *ScrollContainer) *scrollBar {
	return &scrollBar{scroll: scroll, minWidth: theme.ScrollBarSmallSize()}
}

type scrollRenderer struct {
	scroll                  *ScrollContainer
	vertBar                 *scrollBar
	topShadow, bottomShadow fyne.CanvasObject

	objects []fyne.CanvasObject
}

func (s *scrollRenderer) updatePosition() {
	scrollHeight := s.scroll.Size().Height
	contentHeight := s.scroll.Content.Size().Height
	if contentHeight <= scrollHeight {
		s.scroll.Offset.Y = 0
		s.vertBar.Hide()
	} else {
		s.vertBar.Show()
		if contentHeight-s.scroll.Offset.Y < scrollHeight {
			s.scroll.Offset.Y = contentHeight - scrollHeight
		}
	}
	s.scroll.Content.Move(fyne.NewPos(-s.scroll.Offset.X, -s.scroll.Offset.Y))
	canvas.Refresh(s.scroll.Content)

	if s.scroll.Offset.Y > 0 {
		s.topShadow.Show()
	} else {
		s.topShadow.Hide()
	}
	if s.scroll.Offset.Y < contentHeight-scrollHeight {
		s.bottomShadow.Show()
	} else {
		s.bottomShadow.Hide()
	}

	Renderer(s.vertBar).Layout(s.scroll.size)
}

func (s *scrollRenderer) Layout(size fyne.Size) {
	// The scroll bar needs to be resized and moved on the far right
	scrollBar := s.vertBar
	scrollBar.Resize(fyne.NewSize(scrollBar.MinSize().Width, size.Height))
	scrollBar.Move(fyne.NewPos(s.scroll.Size().Width-scrollBar.Size().Width, 0))
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

func (s *scrollRenderer) Refresh() {
	s.Layout(s.scroll.Size())
}

func (s *scrollRenderer) ApplyTheme() {
}

func (s *scrollRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (s *scrollRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *scrollRenderer) Destroy() {
}

// ScrollContainer defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
type ScrollContainer struct {
	baseWidget

	Content fyne.CanvasObject
	Offset  fyne.Position
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

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *ScrollContainer) Resize(size fyne.Size) {
	s.resize(size, s)
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *ScrollContainer) Move(pos fyne.Position) {
	s.move(pos, s)
}

// MinSize returns the smallest size this widget can shrink to
func (s *ScrollContainer) MinSize() fyne.Size {
	return s.minSize(s)
}

// Show this widget, if it was previously hidden
func (s *ScrollContainer) Show() {
	s.show(s)
}

// Hide this widget, if it was previously visible
func (s *ScrollContainer) Hide() {
	s.hide(s)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *ScrollContainer) CreateRenderer() fyne.WidgetRenderer {
	bar := newScrollBar(s)
	topShadow := newShadow(shadowBottom, theme.Padding()*2)
	bottomShadow := newShadow(shadowTop, theme.Padding()*2)
	return &scrollRenderer{
		objects:      []fyne.CanvasObject{s.Content, bar, topShadow, bottomShadow},
		scroll:       s,
		vertBar:      bar,
		topShadow:    topShadow,
		bottomShadow: bottomShadow,
	}
}

// NewScrollContainer creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed objects.
func NewScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	return &ScrollContainer{Content: content}
}
