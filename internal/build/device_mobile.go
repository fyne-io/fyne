//go:build mobile
// +build mobile

package build

// IsMobile returns true if running on a mobile device.
func IsMobile() bool {
	return true
}
