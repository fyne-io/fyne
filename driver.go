package fyne

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	// Create a new UI Window
	CreateWindow(string) Window
	// Get a slice containing all app windows
	AllWindows() []Window

	// Return the size required to render the given string of specified font size and style
	RenderedTextSize(string, int, TextStyle) Size
	// Close the driver and open windows then exit the application
	Quit()
}

var driver Driver

// GetDriver returns the current render driver.
// This is basically an internal detail - use with caution.
func GetDriver() Driver {
	return driver
}

// setDriver sets the application driver.
// This bridges internal modularity - do not call this method directly.
func setDriver(d Driver) {
	driver = d
}
