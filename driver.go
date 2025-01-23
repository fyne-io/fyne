package fyne

import "time"

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	// CreateWindow creates a new UI Window for a certain implementation.
	// Developers should use [App.NewWindow].
	CreateWindow(string) Window
	// AllWindows returns a slice containing all app windows.
	AllWindows() []Window

	// RenderedTextSize returns the size required to render the given string of specified
	// font size and style. It also returns the height to text baseline, measured from the top.
	// If the source is specified it will be used, otherwise the current theme will be asked for the font.
	RenderedTextSize(text string, fontSize float32, style TextStyle, source Resource) (size Size, baseline float32)

	// CanvasForObject returns the canvas that is associated with a given [CanvasObject].
	CanvasForObject(CanvasObject) Canvas
	// AbsolutePositionForObject returns the position of a given [CanvasObject] relative to the top/left of a canvas.
	AbsolutePositionForObject(CanvasObject) Position

	// Device returns the device that the application is currently running on.
	Device() Device
	// Run starts the main event loop of the driver.
	Run()
	// Quit closes the driver and open windows, then exit the application.
	// On some operating systems this does nothing, for example iOS and Android.
	Quit()

	// StartAnimation registers a new animation with this driver and requests it be started.
	// Developers should use the [Animation.Start] function.
	StartAnimation(*Animation)
	// StopAnimation stops an animation and unregisters from this driver.
	// Developers should use the [Animation.Stop] function.
	StopAnimation(*Animation)

	// DoubleTapDelay returns the maximum duration where a second tap after a first one
	// will be considered a [DoubleTap] instead of two distinct [Tap] events.
	//
	// Since: 2.5
	DoubleTapDelay() time.Duration

	// SetDisableScreenBlanking allows an app to ask the device not to sleep/lock/blank displays
	//
	// Since: 2.5
	SetDisableScreenBlanking(bool)

	// DoFromGoroutine provides a way to queue a function `fn` that is running on a goroutine back to
	// the central thread for Fyne updates, waiting for it to return if `wait` is true.
	// The driver provides the implementation normally accessed through [fyne.Do].
	// This is required when background tasks want to execute code safely in the graphical context.
	//
	// Since: 2.6
	DoFromGoroutine(fn func(), wait bool)
}
