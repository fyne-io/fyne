package fyneapp

import "github.com/fyne-io/fyne/ui"

// Driver defines an abstract concept of a Fyne render driver.
// Any implementation must provide at least these methods.
type Driver interface {
	Run()                          // Start the driver
	Quit()                         // Cleanly exit the driver
	CreateWindow(string) ui.Window // Create a new UI Window
}
