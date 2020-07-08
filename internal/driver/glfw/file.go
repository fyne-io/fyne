package glfw

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

type file struct {
	*os.File
	path string
}

func (d *gLDriver) FileReaderForURI(uri fyne.URI) (fyne.URIReadCloser, error) {
	return openFile(uri, false)
}

func (f *file) Name() string {
	return filepath.Base(f.path)
}

func (f *file) URI() fyne.URI {
	return storage.NewURI("file://" + f.path)
}

func (d *gLDriver) FileWriterForURI(uri fyne.URI) (fyne.URIWriteCloser, error) {
	return openFile(uri, true)
}

func openFile(uri fyne.URI, create bool) (*file, error) {
	if uri.Scheme() != "file" {
		return nil, fmt.Errorf("invalid URI for file: %s", uri)
	}

	path := uri.String()[7:]
	var f *os.File
	var err error
	if create {
		f, err = os.Create(path) // If it exists this will truncate which is what we wanted
	} else {
		f, err = os.Open(path)
	}
	return &file{File: f, path: path}, err
}
