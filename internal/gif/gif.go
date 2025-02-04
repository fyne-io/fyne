package gif

import (
	"bytes"
	"fyne.io/fyne/v2"
	"path/filepath"
	"strings"
)

func IsFileGIF(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".gif")
}

// IsResourceGIF checks if the resource is a GIF or not.
func IsResourceGIF(res fyne.Resource) bool {
	if IsFileGIF(res.Name()) {
		return true
	}

	var content []byte = res.Content()
	// Check if content starts with GIF header "GIF87a" or "GIF89a"
	if len(content) >= 6 {
		if bytes.Equal(content[:6], []byte("GIF87a")) || bytes.Equal(content[:6], []byte("GIF89a")) {
			return true
		}
	}
	return false
}
