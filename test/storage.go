package test

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

type testStorage struct {
}

func (s *testStorage) RootURI() fyne.URI {
	return storage.NewURI("file://" + os.TempDir())
}
