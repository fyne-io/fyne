package test

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/storage"
)

type testStorage struct {
	*internal.Docs
}

func (s *testStorage) RootURI() fyne.URI {
	return storage.NewFileURI(os.TempDir())
}

func (s *testStorage) docRootURI() (fyne.URI, error) {
	return storage.Child(s.RootURI(), "Documents")
}
