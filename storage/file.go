// Package storage provides storage access and management functionality.
package storage

import (
	"fmt"

	"fyne.io/fyne"
)

// OpenFileFromURI loads a file read stream from a resource identifier.
// This is mostly provided so that file references can be saved using their URI and loaded again later.
func OpenFileFromURI(uri fyne.URI) (fyne.URIReadCloser, error) {
	return fyne.CurrentApp().Driver().FileReaderForURI(uri)
}

// SaveFileToURI loads a file write stream to a resource identifier.
// This is mostly provided so that file references can be saved using their URI and written to again later.
func SaveFileToURI(uri fyne.URI) (fyne.URIWriteCloser, error) {
	return fyne.CurrentApp().Driver().FileWriterForURI(uri)
}

// FileHandler implements URIHandler
var _ fyne.URIHandler = (*FileHandler)(nil)

// XXX: we might be better off implementing this inside of the driver
type FileHandler struct {
}

func (f *FileHandler) Validate(u fyne.URI) error {
	if u.Scheme() != "file" {
		return fmt.Errorf("FileHandler does not implement scheme '%s'", u.Scheme())
	}

	return nil
}

func (f *FileHandler) Create(u fyne.URI) error {

	// TODO: this may need some support from the driver??
	return nil
}

func (f *FileHandler) WriterTo(u fyne.URI) (fyne.URIWriteCloser, error) {
	return fyne.CurrentApp().Driver().FileWriterForURI(u)
}

func (f *FileHandler) ReaderFrom(u fyne.URI) (fyne.URIReadCloser, error) {
	return fyne.CurrentApp().Driver().FileReaderForURI(u)
}

func (f *FileHandler) ListerOf(u fyne.URI) (fyne.ListableURI, error) {

	// TODO: may also need some driver support
	return nil, nil
}
