// +build linux openbsd freebsd netbsd
// +build !android

package dialog

import (
	"os"
	"strings"
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
			t.Errorf("missing favorite location: %s", name)
			continue
		}

		if subdir == "" {
			assert.Equal(t, homeURI.String(), fav.String())
		} else {
			ok, err := storage.Exists(fav)
			assert.Nil(t, err)
			if ok {
				// folder found
				continue
			}
			fallbackName := fav.Name()
			if !strings.EqualFold(subdir, fallbackName) {
				t.Errorf("%s should equal fold with %s", subdir, fallbackName)
			}
		}
	}
}
