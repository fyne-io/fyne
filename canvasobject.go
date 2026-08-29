package fyne

import (
	"image"
)

// CanvasObject describes any graphical object that can be added to a canvas.
// Objects have a size and position that can be controlled through this API.
// MinSize is used to determine the minimum size which this object should be displayed.
// An object will be visible by default but can be hidden with Hide() and re-shown with Show().
//
// Note: If this object is controlled as part of a Layout you should not call
// Resize(Size) or Move(Position).
type CanvasObject interface {
	// geometry

	// MinSize returns the minimum size this object needs to be drawn.
	MinSize() Size
	// Move moves this object to the given position relative to its parent.
	// This should only be called if your object is not in a container with a layout manager.
	Move(Position)
	// Position returns the current position of the object relative to its parent.
	Position() Position
	// Resize resizes this object to the given size.
	// This should only be called if your object is not in a container with a layout manager.
	Resize(Size)
	// Size returns the current size of this object.
	Size() Size

	// visibility

	// Hide hides this object.
	Hide()
	// Visible returns whether this object is visible or not.
	Visible() bool
	// Show shows this object.
	Show()

	// Refresh must be called if this object should be redrawn because its inner state changed.
	Refresh()
}

// Disableable describes any [CanvasObject] that can be disabled.
// This is primarily used with objects that also implement the Tappable interface.
type Disableable interface {
	Enable()
	Disable()
	Disabled() bool
}

// DoubleTappable describes any [CanvasObject] that can also be double tapped.
type DoubleTappable interface {
	DoubleTapped(*PointEvent)
}

// Draggable indicates that a [CanvasObject] can be dragged.
// This is used for any item that the user has indicated should be moved across the screen.
type Draggable interface {
	Dragged(*DragEvent)
	DragEnd()
}

// Focusable describes any [CanvasObject] that can respond to being focused.
// It will receive the FocusGained and FocusLost events appropriately.
// When focused it will also have TypedRune called as text is input and
// TypedKey called when other keys are pressed.
//
// Note: You must not change canvas state (including overlays or focus) in FocusGained or FocusLost
// or you would end up with a dead-lock.
type Focusable interface {
	// FocusGained is a hook called by the focus handling logic after this object gained the focus.
	FocusGained()
	// FocusLost is a hook called by the focus handling logic after this object lost the focus.
	FocusLost()

	// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
	TypedRune(rune)
	// TypedKey is a hook called by the input handling logic on key events if this object is focused.
	TypedKey(*KeyEvent)
}

// Scrollable describes any [CanvasObject] that can also be scrolled.
// This is mostly used to implement the widget.ScrollContainer.
type Scrollable interface {
	Scrolled(*ScrollEvent)
}

// SecondaryTappable describes a [CanvasObject] that can be right-clicked or long-tapped.
type SecondaryTappable interface {
	TappedSecondary(*PointEvent)
}

// Shortcutable describes any [CanvasObject] that can respond to shortcut commands (quit, cut, copy, and paste).
type Shortcutable interface {
	TypedShortcut(Shortcut)
}

// Tabbable describes any object that needs to accept the Tab key presses.
//
// Since: 2.1
type Tabbable interface {
	// AcceptsTab is a hook called by the key press handling logic.
	// If it returns true then the Tab key events will be sent using TypedKey.
	AcceptsTab() bool
}

// Tappable describes any [CanvasObject] that can also be tapped.
// This should be implemented by buttons etc that wish to handle pointer interactions.
type Tappable interface {
	Tapped(*PointEvent)
}

type DragItem interface {
	MimeType() string
	URI() URI
}

type DragItemCursor interface {
	// Image returns the image for the given cursor, or nil if none should be shown.
	// It also returns the x and y pixels that should act as the hot-spot (measured from top left corner).
	Image() (image.Image, int, int)
}

// DragSource is implemented by CanvasObjects to indicate that they can be dragged and dropped
// onto another CanvasObject which implements DropTarget.
type DragSource interface {
	// DragSourceBegin is invoked when a drag and drop is initiated on this DragSource.
	// It should return the DragItems it is sourcing, and optionally a non-nil DragItemCursor
	// which will be used to visually represent the sourced content during the in-progress drag.
	//
	// Once a drag begins, either DropAccepted or DropCanceled
	// will be invoked to notify this DragSource of the outcome of the drag.
	// This is guaranteed to happen before another invocation of DragSourceBegin may be called.
	DragSourceBegin(*PointEvent) ([]DragItem, DragItemCursor)

	// DragSourceDragged is invoked when a drag and drop continues within this DragSource
	// and has not yet entered over a potential DropTarget.
	DragSourceDragged(*PointEvent)

	// DropPending is invoked the first time an in progress drag and drop enters
	// over a specific potential DropTarget that may accept the drop.
	// The items that the drop target can accept are passed to this invocation,
	// which MAY be a subset of the items returned by DragSourceBegin.
	// Implementations may wish to use this hook to update rendering to
	// highlight the item(s) which may be about to be dropped.
	DropPending([]DragItem)

	// DropUnpending is invoked when an in progress drag and drop which had previously
	// entered a potential DropTarget (and having invoked DropPending()) leaves
	// that DropTarget, but the drag still remains in progress.
	// Implementations may wish to cancel any visual changes made in response
	// to the corresponding DropPending invocation.
	DropUnpending()

	// DropCanceled is invoked when a drag and drop ends without being
	// accepted by a DropTarget.
	DropCanceled()

	// DropCanceled is invoked when a drag and drop ends with the content
	// being accepted by a DropTarget. Implementations may wish to update
	// their visual rendering and backend state to reflect that the dragged
	// content has been moved.
	DropAccepted([]DragItem)
}

// DropTarget is implemented by CanvasObjects to signify that they can accept items dropped
// from another CanvasObject implementing DragSource.
type DropTarget interface {
	// DropBegin is invoked when a drag and drop beginning from a DragSource
	// enters over this DropTarget. Implementations should return the (sub)-set of items that
	// the drop can accept, (possibly based on the content types of the DragItems),
	// and optionally a non-nil cursor to visually represent the content which may be
	// dropped onto this DropTarget.
	// If a nil or zero-length DragItem slice is returned, this DropTarget is considered
	// to not accept the drag and drop, and no further callbacks will be invoked for
	// the in-progress drag and drop operation.
	// If both the DragSource and this DropTarget return a cursor, the DropTarget
	// cursor will be used to represent the drag and drop while it is dragged over this target.
	DropBegin(*PointEvent, []DragItem) (accept []DragItem, cursor DragItemCursor)

	// DropDragged is invoked when a drag and drop that had entered this DropTarget
	// (resulting in a DropBegin invocation) continues to be dragged within this DropTarget.
	// Implementations may wish to update rendering to show where the dragged content
	// would land if dropped.
	DropDragged(*PointEvent)

	// DropEnd is invoked when a drag and drop that had entered this DropTarget
	// leaves without being dropped.
	DropEnd()

	// Dropped is invoked when a drag and drop ends by being dropped on this DropTarget.
	Dropped(*PointEvent)
}
