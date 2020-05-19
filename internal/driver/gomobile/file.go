package gomobile

import (
	"errors"
	"io"
	"net/url"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

type fileOpen struct {
	io.ReadCloser
	uri  fyne.URI
	done func()
}

func (f *fileOpen) Name() string {
	return nameFromURI(f.uri)
}

func (f *fileOpen) URI() fyne.URI {
	return f.uri
}

func (d *mobileDriver) FileReaderForURI(u fyne.URI) (fyne.FileReadCloser, error) {
	file := &fileOpen{uri: u}
	read, err := nativeFileOpen(file)
	file.ReadCloser = read
	return file, err
}

func (d *mobileDriver) FileWriterForURI(u fyne.URI) (fyne.FileWriteCloser, error) {
	return nil, errors.New("file writing is not supported on mobile")
}

func nameFromURI(uri fyne.URI) string {
	u, err := url.Parse(uri.String())
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
			f, err := drv.FileReaderForURI(storage.NewURI(uri))
			f.(*fileOpen).done = closer
			callback(f, err)
		})
	}
}
