package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne"
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

func (f *file) URI() string {
	return "file://" + f.path
}

func openFile(uri string, create bool) (*file, error) {
	if len(uri) < 8 || uri[:7] != "file://" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := uri[7:]
	f, err := os.Open(path)
	if err != nil && create {
		f, err = os.Create(path)
	}
	return &file{File: f, path: path}, err
}

func (d *testDriver) FileReaderForURI(uri string) (fyne.FileReader, error) {
	return openFile(uri, false)
}

func (d *testDriver) FileWriterForURI(uri string) (fyne.FileWriter, error) {
	return openFile(uri, true)
}
