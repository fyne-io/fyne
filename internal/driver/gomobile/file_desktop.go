// +build !ios,!android

package gomobile

import (
	"errors"
	"io"
	"os"

	intRepo "fyne.io/fyne/internal/repository"
	"fyne.io/fyne/storage/repository"
)

func nativeFileOpen(f *fileOpen) (io.ReadCloser, error) {
	if f.uri.Scheme() != "file" {
		return nil, errors.New("mobile simulator mode only supports file:// URIs")
	}

	return os.Open(f.uri.String()[7:])
}

func registerRepository(d *mobileDriver) {
	repository.Register("file", intRepo.NewFileRepository("file"))
}
