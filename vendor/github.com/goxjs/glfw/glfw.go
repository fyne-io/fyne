// Package glfw experimentally provides a glfw-like API
// with desktop (via glfw) and browser (via HTML5 canvas) backends.
//
// It is used for creating a GL context and receiving events.
//
// Note: This package is currently in development. The API is incomplete and may change.
package glfw

// ContextWatcher is a general mechanism for being notified when context is made current or detached.
type ContextWatcher interface {
	// OnMakeCurrent is called after a context is made current.
	// context is is a platform-specific representation of the context, or nil if unavailable.
	OnMakeCurrent(context interface{})

	// OnDetach is called after the current context is detached.
	OnDetach()
}

// VidMode describes a single video mode.
type VidMode struct {
	Width       int // The width, in pixels, of the video mode.
	Height      int // The height, in pixels, of the video mode.
	RedBits     int // The bit depth of the red channel of the video mode.
	GreenBits   int // The bit depth of the green channel of the video mode.
	BlueBits    int // The bit depth of the blue channel of the video mode.
	RefreshRate int // The refresh rate, in Hz, of the video mode.
}
