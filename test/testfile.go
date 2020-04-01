package test

import (
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne"
)

type file struct {
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

func (d *testDriver) FileFromURI(uri string) fyne.File {
	if len(uri) < 8 || uri[:7] != "file://" {
		return nil
	}

	return &file{path: uri[7:]}
}
