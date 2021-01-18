// Package storage provides storage access and management functionality.
package storage

import (
	"errors"

	"fyne.io/fyne/v2"
)

// OpenFileFromURI loads a file read stream from a resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
//
// Deprecated: this has been replaced by storage.Reader(URI)
func OpenFileFromURI(uri fyne.URI) (fyne.URIReadCloser, error) {
	return Reader(uri)
}

// SaveFileToURI loads a file write stream to a resource identifier.
// This is mostly provided so that file references can be saved using their URI and written to again later.
//
// Deprecated: this has been replaced by storage.Writer(URI)
func SaveFileToURI(uri fyne.URI) (fyne.URIWriteCloser, error) {
	return Writer(uri)
}

// ListerForURI will attempt to use the application's driver to convert a
// standard URI into a listable URI.
//
// Since: 1.4
func ListerForURI(uri fyne.URI) (fyne.ListableURI, error) {
	listable, err := CanList(uri)
	if err != nil {
		return nil, err
	}
	if !listable {
		return nil, errors.New("uri is not listable")
	}

	return &legacyListable{uri}, nil
}

type legacyListable struct {
	fyne.URI
}

func (l *legacyListable) List() ([]fyne.URI, error) {
	return List(l.URI)
}
