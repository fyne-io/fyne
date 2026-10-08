// Package mouse defines an event for mouse input.
//
// Mouse input is only reported when running a mobile application in the
// simulator on a desktop machine (Linux, BSD or macOS), devices with a touch
// panel only report touch events.
package mouse // import "fyne.io/fyne/v2/internal/driver/mobile/event/mouse"

// ScrollEvent is a mouse scroll event.
type ScrollEvent struct {
	// X and Y are the mouse location, in pixels from the top-left of the screen.
	X, Y float32

	// ScrollY is how far the scroll wheel moved at that location, in lines,
	// positive moves the view up, horizontal scrolling is not supported.
	ScrollY float32
}
