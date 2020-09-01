package fyne

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	// CreateWindow creates a new UI Window.
	CreateWindow(string) Window
	// AllWindows returns a slice containing all app windows.
	AllWindows() []Window

	// RenderedTextSize returns the size required to render the given string of specified
	// font size and style.
	RenderedTextSize(string, int, TextStyle) Size

	// FileReaderForURI opens a file reader for the given resource indicator.
	// This may refer to a filesystem (typical on desktop) or data from another application.
	FileReaderForURI(URI) (URIReadCloser, error)

	// FileWriterForURI opens a file writer for the given resource indicator.
	// This should refer to a filesystem resource as external data will not be writable.
	FileWriterForURI(URI) (URIWriteCloser, error)

	// ListerForURI converts a URI to a listable URI, if it is possible to do so.
	ListerForURI(URI) (ListableURI, error)

	// CanvasForObject returns the canvas that is associated with a given CanvasObject.
	CanvasForObject(CanvasObject) Canvas
	// AbsolutePositionForObject returns the position of a given CanvasObject relative to the top/left of a canvas.
	AbsolutePositionForObject(CanvasObject) Position

	// Device returns the device that the application is currently running on.
	Device() Device
	// Run starts the main event loop of the driver.
	Run()
	// Quit closes the driver and open windows, then exit the application.
	Quit()
}
