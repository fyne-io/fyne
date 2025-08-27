//go:build (linux || openbsd || freebsd || netbsd) && !android && !wasm && !js

package dialog

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	const cmdName = "xdg-user-dir"
	if _, err := exec.LookPath(cmdName); err != nil {
		return storage.Child(homeURI, name) // no lookup possible
	}

	lookupName := strings.ToUpper(name)
	cmd := exec.Command(cmdName, lookupName)
	loc, err := cmd.Output()
	if err != nil {
		return storage.Child(homeURI, name)
	}

	// Remove \n at the end
	loc = loc[:len(loc)-1]
	locURI := storage.NewFileURI(string(loc))

	if strings.TrimRight(locURI.String(), "/") == strings.TrimRight(homeURI.String(), "/") {
		fallback, _ := storage.Child(homeURI, name)
		return fallback, fmt.Errorf("this computer does not define a %s folder", lookupName)
	}

	return locURI, nil
}

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeURI := storage.NewFileURI(homeDir)
	home, _ := storage.ListerForURI(homeURI)

	favoriteLocations := map[string]fyne.ListableURI{"Home": home}
	for _, favName := range getFavoritesOrder() {
		uri, err1 := getFavoriteLocation(homeURI, favName)
		if err != nil {
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
