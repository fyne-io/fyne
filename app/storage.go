package app

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

type store struct {
	a *fyneApp
}

func (s *store) RootURI() fyne.URI {
	if s.a.uniqueID == "" {
		fyne.LogError("Storage API requires a unique ID, use app.NewWithID()", nil)
		return storage.NewURI("file://" + os.TempDir())
	}

	return storage.NewURI("file://" + s.a.storageRoot())
}
