package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/widget"
)

// Scroll defines a container that is smaller than the Content.
// The Offset is used to determine the position of the child widgets within the container.
//
// Since: 1.4
type Scroll = widget.Scroll

// ScrollDirection represents the directions in which a Scroll container can scroll its child content.
//
// Since: 1.4
type ScrollDirection = widget.ScrollDirection

// Constants for valid values of ScrollDirection.
const (
	// ScrollBoth supports horizontal and vertical scrolling.
	ScrollBoth ScrollDirection = widget.ScrollBoth
	// ScrollHorizontalOnly specifies the scrolling should only happen left to right.
	ScrollHorizontalOnly = widget.ScrollHorizontalOnly
	// ScrollVerticalOnly specifies the scrolling should only happen top to bottom.
	ScrollVerticalOnly = widget.ScrollVerticalOnly
	// ScrollNone turns off scrolling for this container.
	//
	// Since: 2.1
	ScrollNone = widget.ScrollNone
)

// NewScroll creates a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize to be smaller than that of the passed object.
//
// Since: 1.4
func NewScroll(content fyne.CanvasObject) *Scroll {
	return widget.NewScroll(content)
}

// NewHScroll create a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Width to be smaller than that of the passed object.
//
// Since: 1.4
func NewHScroll(content fyne.CanvasObject) *Scroll {
	return widget.NewHScroll(content)
}

// NewVScroll a scrollable parent wrapping the specified content.
// Note that this may cause the MinSize.Height to be smaller than that of the passed object.
//
// Since: 1.4
func NewVScroll(content fyne.CanvasObject) *Scroll {
	return widget.NewVScroll(content)
}
