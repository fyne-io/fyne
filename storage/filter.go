package storage

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/repository/mime"
)

// FileFilter is an interface that can be implemented to provide a filter to a file dialog.
type FileFilter interface {
	Matches(fyne.URI) bool
}

// ExtensionFileFilter represents a file filter based on the the ending of file names,
// for example ".txt" and ".png".
type ExtensionFileFilter struct {
	Extensions []string
}

// MimeTypeFileFilter represents a file filter based on the files mime type,
// for example "image/*", "audio/mp3".
type MimeTypeFileFilter struct {
	MimeTypes []string
}

// Matches returns true if a file URI has one of the filtered extensions.
func (e *ExtensionFileFilter) Matches(uri fyne.URI) bool {
	extension := uri.Extension()
	for _, ext := range e.Extensions {
		if strings.EqualFold(extension, ext) {
			return true
		}
	}
	return false
}

// NewExtensionFileFilter takes a string slice of extensions with a leading . and creates a filter for the file dialog.
// Example: .jpg, .mp3, .txt, .sh
func NewExtensionFileFilter(extensions []string) FileFilter {
	return &ExtensionFileFilter{Extensions: extensions}
}

// Matches returns true if a file URI has one of the filtered mimetypes.
func (mt *MimeTypeFileFilter) Matches(uri fyne.URI) bool {
	mimeType, mimeSubType := mime.Split(uri.MimeType())
	for _, mimeTypeFull := range mt.MimeTypes {
		mType, mSubType := mime.Split(mimeTypeFull)
		if mType == "" || mSubType == "" {
			continue
		}

		// Replace with strings.Cut() when Go 1.18 is our new base version.
		subTypeSeparatorIndex := strings.IndexByte(mSubType, ';')
		if subTypeSeparatorIndex != -1 {
			mSubType = mSubType[:subTypeSeparatorIndex]
		}

		if mType == mimeType && (mSubType == mimeSubType || mSubType == "*") {
			return true
		}
	}
	return false
}

// NewMimeTypeFileFilter takes a string slice of mimetypes, including globs, and creates a filter for the file dialog.
// Example: image/*, audio/mp3, text/plain, application/*
func NewMimeTypeFileFilter(mimeTypes []string) FileFilter {
	return &MimeTypeFileFilter{MimeTypes: mimeTypes}
}
