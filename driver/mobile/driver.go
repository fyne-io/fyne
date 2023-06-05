// Package mobile provides desktop specific mobile functionality.
package mobile

// Driver represents the extended capabilities of a mobile driver
//
// Since: 2.4
type Driver interface {
	// GoBack asks the OS to go to the previous app / activity, where supported
	GoBack()
}
