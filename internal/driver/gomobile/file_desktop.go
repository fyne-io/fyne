// +build !ios,!android

package gomobile

import (
	"io"

	intRepo "fyne.io/fyne/internal/repository"
	"fyne.io/fyne/storage/repository"
)

func nativeFileOpen(*fileOpen) (io.ReadCloser, error) {
	// no-op as we use the internal FileRepository
	return nil, nil
}

func registerRepository(d *mobileDriver) {
	repository.Register("file", intRepo.NewFileRepository())
}
