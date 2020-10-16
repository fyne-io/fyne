package util

import "strings"

// IsMobile returns true if the given os parameter represents a platform handled by gomobile.
func IsMobile(os string) bool {
	return os == "ios" || strings.HasPrefix(os, "android")
}
