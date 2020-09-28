package glfw

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
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

func (d *directory) List() ([]fyne.URI, error) {
	if d.Scheme() != "file" {
		return nil, fmt.Errorf("unsupported URL protocol")
	}

	path := d.String()[len(d.Scheme())+3 : len(d.String())]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	urilist := []fyne.URI{}

	for _, f := range files {
		uri := storage.NewURI("file://" + filepath.Join(path, f.Name()))
		urilist = append(urilist, uri)
	}

	return urilist, nil
}

func (d *gLDriver) ListerForURI(uri fyne.URI) (fyne.ListableURI, error) {
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

func (d *gLDriver) FileReaderForURI(uri fyne.URI) (fyne.URIReadCloser, error) {
	return openFile(uri, false)
}

func (f *file) Name() string {
	return f.URI().Name()
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
