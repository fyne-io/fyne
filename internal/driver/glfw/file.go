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

func (d *gLDriver) FileReaderForURI(uri fyne.URI) (fyne.FileReadCloser, error) {
	return openFile(uri, false)
}

func (f *file) Name() string {
	return filepath.Base(f.path)
}

func (f *file) URI() fyne.URI {
	return storage.NewURI("file://" + f.path)
}

type fileWriter struct {
	*os.File
	path string
}

func (d *gLDriver) FileWriterForURI(uri fyne.URI) (fyne.FileWriteCloser, error) {
	return openFile(uri, true)
}

func (f *fileWriter) Name() string {
	return filepath.Base(f.path)
}

func (f *fileWriter) URI() fyne.URI {
	return storage.NewURI("file://" + f.path)
}

func openFile(uri fyne.URI, create bool) (*file, error) {
	if uri.Scheme() != "file" {
		return nil, fmt.Errorf("invalid URI for file: %s", uri)
	}

	path := uri.String()[7:]
	f, err := os.Open(path)
	if err != nil && create {
		f, err = os.Create(path)
	}

	return &file{File: f, path: path}, err
}
