//go:build !ios && !android && !wasm && !js

package dialog

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeURI := storage.NewFileURI(homeDir)
	home, _ := storage.ListerForURI(homeURI)

	favoriteLocations := map[string]fyne.ListableURI{"Home": home}
	for _, favName := range getFavoritesOrder() {
		uri, err1 := storage.Child(homeURI, favName)
		if err1 != nil {
			err = err1
			continue
		}

		listURI, err1 := storage.ListerForURI(uri)
		if err1 != nil {
			err = err1
			continue
		}
		favoriteLocations[favName] = listURI
	}

	return favoriteLocations, err
}
