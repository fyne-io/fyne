package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type scrollRenderer struct {
	scroll  *ScrollContainer
	vertBar *canvas.Rectangle

	objects []fyne.CanvasObject
}

func (s *scrollRenderer) updatePosition() {
	if s.scroll.Content.Size().Height-s.scroll.Offset.Y < s.scroll.Size().Height {
		if s.scroll.Content.Size().Height <= s.scroll.Size().Height {
			s.scroll.Offset.Y = 0
		} else {
			s.scroll.Offset.Y = s.scroll.Content.Size().Height - s.scroll.Size().Height
		}
	}
	s.scroll.Content.Move(fyne.NewPos(-s.scroll.Offset.X, -s.scroll.Offset.Y))
	canvas.Refresh(s.scroll.Content)

	s.updateBarPosition()
}

func (s *scrollRenderer) updateBarPosition() {
	barSize := s.barSizeVertical()
	barRatio := float32(0.0)
	if s.scroll.Offset.Y != 0 {
		barRatio = float32(s.scroll.Offset.Y) / float32(s.scroll.Content.Size().Height-s.scroll.Size().Height)
	}
	barOff := int(float32(s.scroll.size.Height-barSize.Height) * barRatio)

	s.vertBar.Resize(barSize)
	s.vertBar.Move(fyne.NewPos(s.scroll.size.Width-theme.ScrollBarSize(), barOff))
}

func (s *scrollRenderer) barSizeVertical() fyne.Size {
	portion := float32(s.scroll.size.Height) / float32(s.scroll.Content.Size().Height)
	if portion > 1.0 {
		portion = 1.0
	}

	barHeight := int(float32(s.scroll.size.Height) * portion)
	return fyne.NewSize(theme.ScrollBarSize(), barHeight)
}

func (s *scrollRenderer) Layout(size fyne.Size) {
	c := s.scroll.Content
	c.Resize(c.MinSize().Union(size))

	s.vertBar.Resize(fyne.NewSize(theme.ScrollBarSize(), size.Height))
	s.vertBar.Move(fyne.NewPos(size.Width-theme.ScrollBarSize(), 0))

	s.updatePosition()
}

func (s *scrollRenderer) MinSize() fyne.Size {
	// TODO determine if width or height should be resepected based on a which-way-to-scroll flag
	return fyne.NewSize(s.scroll.Content.MinSize().Width, 25) // TODO consider the smallest useful scroll view?
}

func (s *scrollRenderer) Refresh() {
	s.updatePosition()
}

func (s *scrollRenderer) ApplyTheme() {
	s.vertBar.FillColor = theme.ScrollBarColor()
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

// Move the widget to a new position, relative to it's parent.
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

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (s *ScrollContainer) CreateRenderer() fyne.WidgetRenderer {
	bar := canvas.NewRectangle(theme.ScrollBarColor())
	return &scrollRenderer{scroll: s, vertBar: bar, objects: []fyne.CanvasObject{s.Content, bar}}
}

// NewScrollContainer creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed objects.
func NewScrollContainer(content fyne.CanvasObject) *ScrollContainer {
	return &ScrollContainer{Content: content}
}
