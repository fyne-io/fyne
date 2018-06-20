package fyne

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	CreateWindow(string) Window // Create a new UI Window
	AllWindows() []Window       // Get a slice containing all app windows

	RenderedTextSize(string, int) Size // Return the size required to render a string of specified font size
	Quit()                             // Close the driver and open windows then exit the application
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
