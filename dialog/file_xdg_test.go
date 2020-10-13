// +build linux openbsd freebsd netbsd
// +build !android

package dialog

import (
	"os"
	"testing"

	"fyne.io/fyne/storage"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteLocations(t *testing.T) {
	favoriteLocations, _ := getFavoriteLocations()

	// This is a map containing the base name
	// of every key that should be in favoriteLocations
	expected := map[string]string{
		"Home":      "",
		"Documents": "Documents",
		"Downloads": "Downloads",
	}

	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)
	homeURI := storage.NewFileURI(homeDir)

	for name, subdir := range expected {
		fav, ok := favoriteLocations[name]
		if !ok {
			continue
		}

		if subdir == "" {
			assert.Equal(t, homeURI.String(), fav.String())
		}
	}
}
