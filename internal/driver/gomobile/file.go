package gomobile

import (
	"errors"
	"io"
	"net/url"
	"path/filepath"

	"github.com/fyne-io/mobile/app"

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

func (d *mobileDriver) FileReaderForURI(u fyne.URI) (fyne.URIReadCloser, error) {
	file := &fileOpen{uri: u}
	read, err := nativeFileOpen(file)
	if read == nil {
		return nil, err
	}
	file.ReadCloser = read
	return file, err
}

func (d *mobileDriver) FileWriterForURI(u fyne.URI) (fyne.URIWriteCloser, error) {
	return nil, errors.New("file writing is not supported on mobile")
}

func mobileFilter(filter storage.FileFilter) *app.FileFilter {
	mobile := &app.FileFilter{}

	if f, ok := filter.(*storage.MimeTypeFileFilter); ok {
		mobile.MimeTypes = f.MimeTypes
	} else if f, ok := filter.(*storage.ExtensionFileFilter); ok {
		mobile.Extensions = f.Extensions
	} else {
		fyne.LogError("Custom filter types not supported on mobile", nil)
	}

	return mobile
}

func nameFromURI(uri fyne.URI) string {
	u, err := url.Parse(uri.String())
	if err != nil {
		return "unknown"
	}

	return filepath.Base(u.Path)
}

type hasPicker interface {
	ShowFileOpenPicker(func(string, func()), *app.FileFilter)
}

// ShowFileOpenPicker loads the native file open dialog and returns the chosen file path via the callback func.
func ShowFileOpenPicker(callback func(fyne.URIReadCloser, error), filter storage.FileFilter) {
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasPicker); ok {
		a.ShowFileOpenPicker(func(uri string, closer func()) {
			f, err := drv.FileReaderForURI(storage.NewURI(uri))
			if f != nil {
				f.(*fileOpen).done = closer
			}
			callback(f, err)
		}, mobileFilter(filter))
	}
}
