package gomobile

import (
	"errors"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
)

type fileOpen struct {
	io.ReadCloser
	uri  string
	done func()
}

func (f *fileOpen) Name() string {
	return nameFromURI(f.uri)
}

func (f *fileOpen) URI() string {
	return f.uri
}

func (d *mobileDriver) FileReaderForURI(uriOrPath string) (fyne.FileReadCloser, error) {
	uri := uriOrPath
	if strings.Index(uriOrPath, "://") == -1 {
		uri = "file://" + uriOrPath
	}

	file := &fileOpen{uri: uri}
	read, err := nativeFileOpen(file)
	file.ReadCloser = read
	return file, err
}

func (d *mobileDriver) FileWriterForURI(uriOrPath string) (fyne.FileWriteCloser, error) {
	return nil, errors.New("file writing is not supported on mobile")
}

func nameFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return "unknown"
	}

	return filepath.Base(u.Path)
}

type hasPicker interface {
	ShowFileOpenPicker(callback func(string, func()))
}

// ShowFileOpenPicker loads the native file open dialog and returns the chosen file path via the callback func.
func ShowFileOpenPicker(callback func(fyne.FileReadCloser, error)) {
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasPicker); ok {
		a.ShowFileOpenPicker(func(uri string, closer func()) {
			f, err := drv.FileReaderForURI(uri)
			f.(*fileOpen).done = closer
			callback(f, err)
		})
	}
}
