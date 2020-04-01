package glfw

import (
	"fmt"
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
	return false // TODO can we actually check the read/write status?
}

func (f *file) Name() string {
	return filepath.Base(f.path)
}

func (f *file) URI() string {
	return "file://" + f.path
}

func fileWithPath(path string) fyne.File {
	return &file{path: path}
}

func (d *gLDriver) FileFromURI(uri string) fyne.File {
	if len(uri) < 8 || uri[:7] != "file://" {
		fyne.LogError(fmt.Sprintf("Invalid URI for file: %s", uri), nil)
		return nil
	}

	return fileWithPath(uri[7:])
}
