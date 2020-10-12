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

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	cmdName := "xdg-user-dir"
	fallback, err := storage.Child(homeURI, name)
	var fallbackErr error
	if err != nil {
		fallbackErr = fmt.Errorf("couldn't get fallback: %s", err)
	}

	if _, err := exec.LookPath(cmdName); err != nil {
		if fallbackErr != nil {
			return nil, fallbackErr
		}
		return fallback, fmt.Errorf("%s not found in PATH. using fallback paths", cmdName)
	}

	cmd := exec.Command(cmdName, name)
	stdout := bytes.NewBufferString("")

	cmd.Stdout = stdout
	err = cmd.Run()
	if err != nil {
		if fallbackErr != nil {
			return nil, fallbackErr
		}
		return fallback, err
	}

	loc := stdout.String()
	// Remove \n at the end
	loc = loc[:len(loc)-1]
	loc = "file://" + loc

	if loc == homeURI.String() {
		if fallbackErr != nil {
			return nil, fallbackErr
		}
		return fallback, fmt.Errorf("this computer does not have a %s folder", name)
	}

	return storage.NewURI(loc), nil
}

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeURI := storage.NewURI("file://" + homeDir)

	favoriteNames := getFavoriteOrder()
	arguments := map[string]string{
		"Documents": "DOCUMENTS",
		"Downloads": "DOWNLOADS",
	}

	favoriteLocations := make(map[string]fyne.ListableURI)
	for _, favName := range favoriteNames {
		var uri fyne.URI
		var err1 error
		if favName == "Home" {
			uri = homeURI
		} else {
			favURI, err2 := getFavoriteLocation(homeURI, arguments[favName])
			if err2 != nil {
				err1 = err2
			}
			uri = favURI
		}
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
