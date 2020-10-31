// +build linux openbsd freebsd netbsd
// +build !android

package dialog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func getFavoriteLocation(homeURI fyne.URI, name, fallbackName string) (fyne.URI, error) {
	cmdName := "xdg-user-dir"
	if _, err := exec.LookPath(cmdName); err != nil {
		return storage.Child(homeURI, fallbackName) // no lookup possible
	}

	cmd := exec.Command(cmdName, name)
	stdout := bytes.NewBufferString("")

	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return storage.Child(homeURI, fallbackName)
	}

	loc := stdout.String()
	// Remove \n at the end
	loc = loc[:len(loc)-1]
	locURI := storage.NewFileURI(loc)

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
	}

	favoriteLocations := make(map[string]fyne.ListableURI)
	for _, favName := range favoriteNames {
		var uri fyne.URI
		if favName == "Home" {
			uri = homeURI
		} else {
			uri, err = getFavoriteLocation(homeURI, arguments[favName], favName)
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
