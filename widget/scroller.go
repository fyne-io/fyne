package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

type scrollRenderer struct {
	scroll *ScrollContainer

	objects []fyne.CanvasObject
}

func (s *scrollRenderer) updatePosition() {
	s.scroll.Content.Move(fyne.NewPos(-s.scroll.Offset.X, -s.scroll.Offset.Y))
	canvas.Refresh(s.scroll.Content)
}

func (s *scrollRenderer) Layout(size fyne.Size) {
	c := s.scroll.Content
	c.Resize(c.MinSize())

	s.updatePosition()
}

func (s *scrollRenderer) MinSize() fyne.Size {
	return fyne.NewSize(25, 25) // TODO consider the smallest useful scroll view?
}

func (s *scrollRenderer) Refresh() {
	s.updatePosition()
}

func (s *scrollRenderer) ApplyTheme() {
}

func (s *scrollRenderer) BackgroundColor() color.Color {
	return color.White //theme.BackgroundColor()
}

func (s *scrollRenderer) Objects() []fyne.CanvasObject {
	return s.objects
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
	return &scrollRenderer{scroll: s, objects: []fyne.CanvasObject{s.Content}}
}

// NewScroller creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed objects.
func NewScroller(content fyne.CanvasObject) *ScrollContainer {
	return &ScrollContainer{Content: content}
}
