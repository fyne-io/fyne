package test

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

type testStorage struct {
}

func (s *testStorage) RootURI() fyne.URI {
	return storage.NewURI("file://" + os.TempDir())
}
