package gomobile

import (
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
)

type file struct {
	uri string
	done func()
}

func (f *file) Open() (io.ReadCloser, error) {
	return nativeFileOpen(f)
}

func (f *file) Save() (io.WriteCloser, error) {
	return nativeFileSave(f)
}

func (f *file) ReadOnly() bool {
	if len(f.uri) < 8 || f.uri[:7] != "file://" {
		return true
	}

	return false
}

func (f *file) Name() string {
	u, err := url.Parse(f.uri)
	if err != nil {
		return "unknown"
	}

	return filepath.Base(u.Path)
}

func (f *file) URI() string {
	return f.uri
}

func (d *mobileDriver) FileFromURI(uriOrPath string) fyne.File {
	if strings.Index(uriOrPath, "://") == -1 {
		return &file{uri: "file://" + uriOrPath}
	}

	return &file{uri: uriOrPath}
}

type hasPicker interface {
	ShowFileOpenPicker(callback func(string, func()))
}

// ShowFileOpenPicker loads the native file open dialog and returns the chosen file path via the callback func.
func ShowFileOpenPicker(callback func(fyne.File)) {
	drv := fyne.CurrentApp().Driver().(*mobileDriver)
	if a, ok := drv.app.(hasPicker); ok {
		a.ShowFileOpenPicker(func(uri string, closer func()) {
			log.Print("FILE uri", uri)
			f := drv.FileFromURI(uri)
			f.(*file).done = closer
			log.Println("FILE=", f)
			callback(f)
		})
	}
}
