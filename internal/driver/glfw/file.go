package glfw

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne"
)

type file struct {
	*os.File
	path string
}

func (d *gLDriver) FileReaderForURI(uri string) (fyne.FileReadCloser, error) {
	return openFile(uri, false)
}

func (f *file) Name() string {
	return filepath.Base(f.path)
}

func (f *file) URI() string {
	return "file://" + f.path
}

type fileWriter struct {
	*os.File
	path string
}

func (d *gLDriver) FileWriterForURI(uri string) (fyne.FileWriteCloser, error) {
	return openFile(uri, true)
}

func (f *fileWriter) Name() string {
	return filepath.Base(f.path)
}

func (f *fileWriter) URI() string {
	return "file://" + f.path
}

func openFile(uri string, create bool) (*file, error) {
	if len(uri) < 8 || uri[:7] != "file://" {
		return nil, fmt.Errorf("invalid URI for file: %s", uri)
	}

	path := uri[7:]
	f, err := os.Open(path)
	if err != nil && create {
		f, err = os.Create(path)
	}

	return &file{File: f, path: path}, err
}
