package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
// Deprecated: Use widget.NewFileIcon instead
func NewFileIcon(uri fyne.URI) *widget.FileIcon {
	return widget.NewFileIcon(uri)
}
