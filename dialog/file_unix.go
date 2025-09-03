//go:build !windows && !android && !ios && !wasm && !js

package dialog

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
)

func (f *fileDialog) getPlaces() []favoriteItem {
	lister, err := storage.ListerForURI(storage.NewFileURI("/"))
	if err != nil {
		fyne.LogError("could not create lister for /", err)
		return []favoriteItem{}
	}
	return []favoriteItem{{
		"Computer",
		theme.ComputerIcon(),
		lister,
	}}
}

func isHidden(file fyne.URI) bool {
	if file.Scheme() != "file" {
		fyne.LogError("Cannot check if non file is hidden", nil)
		return false
	}

	name := filepath.Base(file.Path())
	return name == "" || name[0] == '.'
}

func hideFile(_ string) error {
	return nil
}
