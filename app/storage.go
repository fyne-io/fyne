package app

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/storage"
)

type store struct {
	*internal.Docs
	a *fyneApp
}

func (s *store) RootURI() fyne.URI {
	if s.a.UniqueID() == "" {
		fyne.LogError("Storage API requires a unique ID, use app.NewWithID()", nil)
		return storage.NewFileURI(os.TempDir())
	}

	return storage.NewFileURI(s.a.storageRoot())
}

func (s *store) docRootURI() (fyne.URI, error) {
	return storage.Child(s.RootURI(), "Documents")
}
