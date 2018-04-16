package ui

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	Run()                       // Start the driver
	Quit()                      // Cleanly exit the driver
	CreateWindow(string) Window // Create a new UI Window
	AllWindows() []Window       // Get a slice containing all app windows
}

var driver Driver

// GetDriver returns the current render driver.
// This is basically an internal detail - use with caution.
func GetDriver() Driver {
	return driver
}

// SetDriver sets the application driver.
// This bridges internal modularity - do not call this method directly.
func SetDriver(d Driver) {
	driver = d
}
