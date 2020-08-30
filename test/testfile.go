package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

type file struct {
	*os.File
	path string
}

func (f *file) Open() (io.ReadCloser, error) {
	return os.Open(f.path)
}

func (f *file) Save() (io.WriteCloser, error) {
	return os.Open(f.path)
}

func (f *file) ReadOnly() bool {
	return true
}

func (f *file) Name() string {
	return filepath.Base(f.path)
}

func (f *file) URI() fyne.URI {
	return storage.NewURI("file://" + f.path)
}

func openFile(uri fyne.URI, create bool) (*file, error) {
	if uri.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := uri.String()[7:]
	f, err := os.Open(path)
	if err != nil && create {
		f, err = os.Create(path)
	}
	return &file{File: f, path: path}, err
}

func (d *testDriver) FileReaderForURI(uri fyne.URI) (fyne.URIReadCloser, error) {
	return openFile(uri, false)
}

func (d *testDriver) FileWriterForURI(uri fyne.URI) (fyne.URIWriteCloser, error) {
	return openFile(uri, true)
}
