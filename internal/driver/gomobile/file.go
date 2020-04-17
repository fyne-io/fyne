package gomobile

import (
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

//func (f *fileOpen) Save() (io.WriteCloser, error) {
//	return nativeFileSave(f)
//}

func (f *fileOpen) Name() string {
	return nameFromURI(f.uri)
}

func (f *fileOpen) URI() string {
	return f.uri
}

func (d *mobileDriver) FileReaderForURI(uriOrPath string) (fyne.FileReader, error) {
	uri := uriOrPath
	if strings.Index(uriOrPath, "://") == -1 {
		uri = "file://" + uriOrPath
	}

	file := &fileOpen{uri: uri}
	read, err := nativeFileOpen(file)
	file.ReadCloser = read
	return file, err
}

type fileSave struct {
	io.WriteCloser
	uri string
}

func (f *fileSave) Name() string {
	return nameFromURI(f.uri)
}

func (f *fileSave) URI() string {
	return f.uri
}

func (d *mobileDriver) FileWriterForURI(uriOrPath string) (fyne.FileWriter, error) {
	uri := uriOrPath
	if strings.Index(uriOrPath, "://") == -1 {
		uri = "file://" + uriOrPath
	}

	file := &fileSave{uri: uri}
	write, err := nativeFileSave(file)
	file.WriteCloser = write
	return file, err
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
func ShowFileOpenPicker(callback func(fyne.FileReader, error)) {
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasPicker); ok {
		a.ShowFileOpenPicker(func(uri string, closer func()) {
			f, err := drv.FileReaderForURI(uri)
			f.(*fileOpen).done = closer
			callback(f, err)
		})
	}
}
