package fyne

// ScrollDirection represents the directions in which a scrollable container or widget can scroll its child content.
type ScrollDirection int

// Constants for valid values of ScrollDirection used in containers and widgets.
const (
	// ScrollBoth supports horizontal and vertical scrolling.
	ScrollBoth ScrollDirection = iota
	// ScrollHorizontalOnly specifies the scrolling should only happen left to right.
	ScrollHorizontalOnly
	// ScrollVerticalOnly specifies the scrolling should only happen top to bottom.
	ScrollVerticalOnly
	// ScrollNone turns off scrolling for this container.
	//
	// Since: 2.6
	ScrollNone
)
