package fyne

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	// Create a new UI Window.
	CreateWindow(string) Window
	// Get a slice containing all app windows.
	AllWindows() []Window

	// Return the size required to render the given string of specified
	// font size and style.
	RenderedTextSize(string, int, TextStyle) Size

	// Get the canvas that is associated with a given CanvasObject.
	CanvasForObject(CanvasObject) Canvas

	// Start the main event loop of the driver.
	Run()
	// Close the driver and open windows then exit the application.
	Quit()
}
