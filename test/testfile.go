package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

type file struct {
	*os.File
	path string
}

type directory struct {
	fyne.URI
}

// Declare conformity to the ListableURI interface
var _ fyne.ListableURI = (*directory)(nil)

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
	return storage.NewFileURI(f.path)
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

func (d *testDriver) ListerForURI(uri fyne.URI) (fyne.ListableURI, error) {
	if uri.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := uri.String()[len(uri.Scheme())+3 : len(uri.String())]
	s, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !s.IsDir() {
		return nil, fmt.Errorf("path '%s' is not a directory, cannot convert to listable URI", path)
	}

	return &directory{URI: uri}, nil
}

func (d *directory) List() ([]fyne.URI, error) {
	if d.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := d.String()[len(d.Scheme())+3 : len(d.String())]
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	urilist := []fyne.URI{}

	for _, f := range files {
		uri := storage.NewFileURI(filepath.Join(path, f.Name()))
		urilist = append(urilist, uri)
	}

	return urilist, nil
}
