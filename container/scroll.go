package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// Scroll defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
//
// Since: 1.4
type Scroll = widget.ScrollContainer

// ScrollDirection represents the directions in which a Scroll container can scroll its child content.
//
// Since: 1.4
type ScrollDirection = widget.ScrollDirection

// Constants for valid values of ScrollDirection.
const (
	ScrollBoth ScrollDirection = iota
	ScrollHorizontalOnly
	ScrollVerticalOnly
)

// NewScroll creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed object.
//
// Since: 1.4
func NewScroll(content fyne.CanvasObject) *Scroll {
	return widget.NewScrollContainer(content)
}

// NewHScroll create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Width to be smaller than that of the passed object.
//
// Since: 1.4
func NewHScroll(content fyne.CanvasObject) *Scroll {
	return widget.NewHScrollContainer(content)
}

// NewVScroll a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Height to be smaller than that of the passed object.
//
// Since: 1.4
func NewVScroll(content fyne.CanvasObject) *Scroll {
	return widget.NewVScrollContainer(content)
}

// TODO move the implementation into internal/scroll.go in 2.0 when we delete the old API.
// we cannot do that right now due to the fact that this package depends on widgets and they depend on us.
// Once moving the bulk of scroller will go to an internal package, then this and the widget package can both depend.
