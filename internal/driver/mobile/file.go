package mobile

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/storage"
)

type fileOpen struct {
	io.ReadCloser
	uri  fyne.URI
	done func()
}

func (f *fileOpen) URI() fyne.URI {
	return f.uri
}

func fileReaderForURI(u fyne.URI) (fyne.URIReadCloser, error) {
	file := &fileOpen{uri: u}
	read, err := nativeFileOpen(file)
	if read == nil {
		return nil, err
	}
	file.ReadCloser = read
	return file, err
}

func mobileFilter(filter storage.FileFilter) *app.FileFilter {
	mobile := &app.FileFilter{}

	if f, ok := filter.(*storage.MimeTypeFileFilter); ok {
		mobile.MimeTypes = f.MimeTypes
	} else if f, ok := filter.(*storage.ExtensionFileFilter); ok {
		mobile.Extensions = f.Extensions
	} else if filter != nil {
		fyne.LogError("Custom filter types not supported on mobile", nil)
	}

	return mobile
}

type hasOpenPicker interface {
	ShowFileOpenPicker(func(string, func()), *app.FileFilter)
}

// ShowFileOpenPicker loads the native file open dialog and returns the chosen file path via the callback func.
func ShowFileOpenPicker(callback func(fyne.URIReadCloser, error), filter storage.FileFilter) {
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasOpenPicker); ok {
		a.ShowFileOpenPicker(func(uri string, closer func()) {
			if uri == "" {
				callback(nil, nil)
				return
			}
			f, err := fileReaderForURI(nativeURI(uri))
			if f != nil {
				f.(*fileOpen).done = closer
			}
			callback(f, err)
		}, mobileFilter(filter))
	}
}

// ShowFolderOpenPicker loads the native folder open dialog and calls back the chosen directory path as a ListableURI.
func ShowFolderOpenPicker(callback func(fyne.ListableURI, error)) {
	filter := storage.NewMimeTypeFileFilter([]string{"application/x-directory"})
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasOpenPicker); ok {
		a.ShowFileOpenPicker(func(path string, _ func()) {
			if path == "" {
				callback(nil, nil)
				return
			}

			uri, err := storage.ParseURI(path)
			if err != nil {
				callback(nil, err)
				return
			}

			callback(listerForURI(uri))
		}, mobileFilter(filter))
	}
}

type fileSave struct {
	io.WriteCloser
	uri  fyne.URI
	done func()
}

func (f *fileSave) URI() fyne.URI {
	return f.uri
}

func fileWriterForURI(u fyne.URI) (fyne.URIWriteCloser, error) {
	file := &fileSave{uri: u}
	write, err := nativeFileSave(file)
	if write == nil {
		return nil, err
	}
	file.WriteCloser = write
	return file, err
}

type hasSavePicker interface {
	ShowFileSavePicker(func(string, func()), *app.FileFilter, string)
}

// ShowFileSavePicker loads the native file save dialog and returns the chosen file path via the callback func.
func ShowFileSavePicker(callback func(fyne.URIWriteCloser, error), filter storage.FileFilter, filename string) {
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasSavePicker); ok {
		a.ShowFileSavePicker(func(path string, closer func()) {
			if path == "" {
				callback(nil, nil)
				return
			}

			uri, err := storage.ParseURI(path)
			if err != nil {
				callback(nil, err)
				return
			}

			f, err := fileWriterForURI(uri)
			if f != nil {
				f.(*fileSave).done = closer
			}
			callback(f, err)
		}, mobileFilter(filter), filename)
	}
}
