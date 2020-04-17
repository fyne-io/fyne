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

	// FileReaderForURI opens a file reader for the given resource indicator.
	// This may refer to a filesystem (typical on desktop) or data from another application.
	FileReaderForURI(string) (FileReader, error)

	// FileWriterForURI opens a file writer for the given resource indicator.
	// This should refer to a filesystem resource as external data will not be writable.
	FileWriterForURI(string) (FileWriter, error)

	// Get the canvas that is associated with a given CanvasObject.
	CanvasForObject(CanvasObject) Canvas
	// Get the position of a given CanvasObject relative to the top/left of a canvas.
	AbsolutePositionForObject(CanvasObject) Position

	// Get the device that the application is currently running on
	Device() Device
	// Start the main event loop of the driver.
	Run()
	// Close the driver and open windows then exit the application.
	Quit()
}
