//go:build (linux || openbsd || freebsd || netbsd) && !android && !wasm && !js
// +build linux openbsd freebsd netbsd
// +build !android
// +build !wasm
// +build !js

package dialog

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"

	"golang.org/x/sys/execabs"
)

func getFavoriteLocation(homeURI fyne.URI, name, fallbackName string) (fyne.URI, error) {
	cmdName := "xdg-user-dir"
	if _, err := execabs.LookPath(cmdName); err != nil {
		return storage.Child(homeURI, fallbackName) // no lookup possible
	}

	cmd := execabs.Command(cmdName, name)
	loc, err := cmd.Output()
	if err != nil {
		return storage.Child(homeURI, fallbackName)
	}

	// Remove \n at the end
	loc = loc[:len(loc)-1]
	locURI := storage.NewFileURI(string(loc))

	if locURI.String() == homeURI.String() {
		fallback, _ := storage.Child(homeURI, fallbackName)
		return fallback, fmt.Errorf("this computer does not define a %s folder", name)
	}

	return locURI, nil
}

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeURI := storage.NewFileURI(homeDir)

	favoriteNames := getFavoriteOrder()
	arguments := map[string]string{
		"Documents": "DOCUMENTS",
		"Downloads": "DOWNLOAD",
		"Music":     "MUSIC",
		"Pictures":  "PICTURES",
		"Videos":    "VIDEOS",
	}

	home, _ := storage.ListerForURI(homeURI)
	favoriteLocations := map[string]fyne.ListableURI{
		"Home": home,
	}
	for _, favName := range favoriteNames {
		var uri fyne.URI
		uri, err = getFavoriteLocation(homeURI, arguments[favName], favName)

		listURI, err1 := storage.ListerForURI(uri)
		if err1 != nil {
			err = err1
			continue
		}
		favoriteLocations[favName] = listURI
	}

	return favoriteLocations, err
}
